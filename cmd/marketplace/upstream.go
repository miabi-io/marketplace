/*
 * Copyright 2026 Jonas Kaninda
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/jkaninda/logger"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/internal/registry"
	"github.com/miabi-io/marketplace/manifest"
)

type imageRef struct{ image, tag string }

// tagCache dedupes registry calls: several templates share an image (nginx, meilisearch).
type tagCache struct {
	reg *registry.Client
	mu  sync.Mutex
	tag map[string]tagResult
}

type tagResult struct {
	tags []string
	err  error
}

func newTagCache() *tagCache {
	return &tagCache{reg: registry.New(), tag: map[string]tagResult{}}
}

func (c *tagCache) tags(ctx context.Context, image string) ([]string, error) {
	c.mu.Lock()
	r, ok := c.tag[image]
	c.mu.Unlock()
	if ok {
		return r.tags, r.err
	}
	tags, err := c.reg.Tags(ctx, image)
	c.mu.Lock()
	c.tag[image] = tagResult{tags, err}
	c.mu.Unlock()
	return tags, err
}

// selectTemplates returns every template, or only the named ones; an unknown name is an error.
func selectTemplates(cat *catalog.Catalog, names []string) ([]catalog.Template, error) {
	if len(names) == 0 {
		var out []catalog.Template
		for _, t := range cat.Templates() {
			if !t.Meta.Deprecated {
				out = append(out, t)
			}
		}
		return out, nil
	}
	out := make([]catalog.Template, 0, len(names))
	for _, n := range names {
		t, ok := cat.Get(n)
		if !ok {
			return nil, fmt.Errorf("unknown template %q", n)
		}
		out = append(out, *t)
	}
	return out, nil
}

// candidates applies the template's update policy; held is non-empty when the app is frozen.
func candidates(ctx context.Context, cache *tagCache, t catalog.Template, a manifest.AppSpec) (c registry.Candidates, held string, err error) {
	p := t.Meta.Updates[a.Name]
	if p.Hold != "" {
		return registry.Candidates{}, p.Hold, nil
	}
	tags, err := cache.tags(ctx, a.Image)
	if err != nil {
		return registry.Candidates{}, "", err
	}
	if p.Match != "" {
		re := regexp.MustCompile(p.Match)
		kept := tags[:0:0]
		for _, tg := range tags {
			if re.MatchString(tg) {
				kept = append(kept, tg)
			}
		}
		tags = kept
	}
	return registry.Compare(a.Tag, tags), "", nil
}

func freshness(ctx context.Context, cache *tagCache, t catalog.Template) catalog.TemplateFreshness {
	v := t.Latest()
	tf := catalog.TemplateFreshness{Version: v.Version, Status: catalog.FreshCurrent, Apps: []catalog.AppFreshness{}}
	for _, a := range v.Manifest.Applications {
		af := catalog.AppFreshness{App: a.Name, Image: a.Image, Tag: a.Tag, Status: catalog.FreshCurrent}
		c, held, err := candidates(ctx, cache, t, a)
		switch {
		case err != nil:
			af.Status, af.Error = "error", err.Error()
		case held != "":
			af.Status, af.Hold = catalog.FreshHeld, held
		default:
			if !c.Comparable {
				af.Status = "unknown"
				break
			}
			af.LatestInMajor, af.Latest = c.LatestInMajor, c.Latest
			switch {
			case c.LatestInMajor != a.Tag:
				af.Status = catalog.FreshUpdate
			case c.Latest != a.Tag:
				af.Status = catalog.FreshMajor
			}
			if c.Latest != c.LatestInMajor && af.Status == catalog.FreshUpdate {
				af.Status = catalog.FreshMajor
			}
		}
		tf.Apps = append(tf.Apps, af)
		tf.Status = worse(tf.Status, af.Status)
	}
	return tf
}

var statusRank = map[string]int{catalog.FreshCurrent: 0, catalog.FreshHeld: 0, "unknown": 1, catalog.FreshUpdate: 2, catalog.FreshMajor: 3, "error": 4}

func worse(a, b string) string {
	if statusRank[b] > statusRank[a] {
		return b
	}
	return a
}

func outdated(args []string, asJSON bool, writePath string) error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	ts, err := selectTemplates(cat, args)
	if err != nil {
		return err
	}
	ctx := context.Background()
	cache := newTagCache()
	report := catalog.FreshnessReport{CheckedAt: time.Now().UTC().Format(time.DateOnly), Templates: map[string]catalog.TemplateFreshness{}}
	results := make([]catalog.TemplateFreshness, len(ts))
	parallel(len(ts), 8, func(i int) { results[i] = freshness(ctx, cache, ts[i]) })
	for i, t := range ts {
		report.Templates[t.Name] = results[i]
	}

	if writePath != "" {
		if err := writeJSON(writePath, report); err != nil {
			return err
		}
		logger.Info("wrote freshness report", "path", writePath, "templates", len(report.Templates))
	}
	if asJSON {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "TEMPLATE\tVERSION\tAPP\tIMAGE\tCURRENT\tIN MAJOR\tLATEST\tSTATUS")
	for i, t := range ts {
		for _, a := range results[i].Apps {
			if a.Status == catalog.FreshCurrent {
				continue
			}
			status := a.Status
			if a.Error != "" {
				status += ": " + a.Error
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", t.Name, results[i].Version, a.App, a.Image, a.Tag, a.LatestInMajor, a.Latest, status)
		}
	}
	return w.Flush()
}

type bumpResult struct {
	Template string   `json:"template"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Path     string   `json:"path"`
	Changes  []string `json:"changes"`
	Breaking bool     `json:"breaking"`
}

type bumpOptions struct {
	major   bool
	dryRun  bool
	asJSON  bool
	today   string
	setTags map[string]string // app -> tag, overriding the upstream lookup
}

// A major bump also bumps the template major and is marked breaking, since its notes need a human.
func bump(args []string, opt bumpOptions) error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	ts, err := selectTemplates(cat, args)
	if err != nil {
		return err
	}
	if len(opt.setTags) > 0 && len(ts) != 1 {
		return errors.New("--set needs exactly one template")
	}
	ctx := context.Background()
	cache := newTagCache()
	var done []bumpResult
	for _, t := range ts {
		r, err := bumpOne(ctx, cache, t, opt)
		if err != nil {
			return fmt.Errorf("%s: %w", t.Name, err)
		}
		if r != nil {
			done = append(done, *r)
		}
	}
	if opt.asJSON {
		if done == nil {
			done = []bumpResult{}
		}
		data, _ := json.MarshalIndent(done, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	for _, r := range done {
		verb := "created"
		if opt.dryRun {
			verb = "would create"
		}
		fmt.Printf("%s %s %s -> %s: %s\n", verb, r.Template, r.From, r.To, strings.Join(r.Changes, ", "))
	}
	if len(done) == 0 {
		fmt.Println("nothing to bump")
	}
	return nil
}

func bumpOne(ctx context.Context, cache *tagCache, t catalog.Template, opt bumpOptions) (*bumpResult, error) {
	cur := t.Latest()
	targets := map[string]string{}
	for _, a := range cur.Manifest.Applications {
		if tag, ok := opt.setTags[a.Name]; ok {
			if tag != a.Tag {
				targets[a.Name] = tag
			}
			continue
		}
		if len(opt.setTags) > 0 {
			continue
		}
		c, held, err := candidates(ctx, cache, t, a)
		if err != nil {
			logger.Warn("skipping app: tags unavailable", "template", t.Name, "app", a.Name, "error", err)
			continue
		}
		if held != "" {
			continue
		}
		next := c.LatestInMajor
		if opt.major {
			next = c.Latest
		}
		if c.Comparable && next != a.Tag {
			targets[a.Name] = next
		}
	}
	for app := range opt.setTags {
		if !appIn(cur.Manifest, app) {
			return nil, fmt.Errorf("--set: no application %q in %s", app, cur.Version)
		}
	}
	if len(targets) == 0 {
		return nil, nil
	}

	breaking := opt.major && majorChanged(cur.Manifest, targets)
	next := nextVersion(cur.Version, breaking)
	if _, exists := t.FindVersion(next); exists {
		return nil, fmt.Errorf("version %s already exists", next)
	}
	srcDir := filepath.Join(t.Source, t.Name, cur.Version)
	raw, err := os.ReadFile(filepath.Join(srcDir, "template.yaml"))
	if err != nil {
		return nil, err
	}
	text, changes, err := retag(string(raw), cur.Manifest, targets)
	if err != nil {
		return nil, err
	}
	text, err = setManifestVersion(text, cur.Version, next)
	if err != nil {
		return nil, err
	}
	m, err := manifest.Parse([]byte(text))
	if err != nil {
		return nil, fmt.Errorf("bumped manifest is invalid: %w", err)
	}

	notes := "Image updates:\n"
	for _, ch := range changes {
		notes += "- " + ch + "\n"
	}
	if breaking {
		notes += "\nMajor upstream release: review the upstream upgrade notes and describe any manual steps here before merging.\n"
	}
	rel := &catalog.Release{Released: opt.today, Notes: notes, Breaking: breaking}
	if !opt.dryRun {
		if rel.Images, err = resolveDigests(ctx, cache.reg, m); err != nil {
			return nil, err
		}
	}
	dir := filepath.Join(t.Source, t.Name, next)
	res := &bumpResult{Template: t.Name, From: cur.Version, To: next, Path: dir, Changes: changes, Breaking: breaking}
	if opt.dryRun {
		return res, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "template.yaml"), []byte(text), 0o644); err != nil {
		return nil, err
	}
	return res, writeRelease(dir, rel)
}

func appIn(m *manifest.Manifest, name string) bool {
	for _, a := range m.Applications {
		if a.Name == name {
			return true
		}
	}
	return false
}

func majorChanged(m *manifest.Manifest, targets map[string]string) bool {
	for _, a := range m.Applications {
		if next, ok := targets[a.Name]; ok && registry.Compare(a.Tag, []string{next}).LatestInMajor != next {
			return true
		}
	}
	return false
}

func nextVersion(v string, major bool) string {
	parts := strings.SplitN(v, ".", 3)
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	ma, _ := strconv.Atoi(parts[0])
	mi, _ := strconv.Atoi(parts[1])
	if major {
		return fmt.Sprintf("%d.0.0", ma+1)
	}
	return fmt.Sprintf("%d.%d.0", ma, mi+1)
}

var (
	imageLine = regexp.MustCompile(`^(\s*-?\s*)image:\s*["']?([^"'\s#]+)`)
	tagLine   = regexp.MustCompile(`^(\s*)tag:\s*["']?([^"'\s#]+)["']?(.*)$`)
	nameLine  = regexp.MustCompile(`^\s*-\s*name:\s*["']?([^"'\s#]+)`)
)

// retag edits the YAML text line by line rather than re-encoding it, so comments and layout survive.
func retag(text string, m *manifest.Manifest, targets map[string]string) (string, []string, error) {
	images := map[string]string{}
	for _, a := range m.Applications {
		images[a.Name] = a.Image
	}
	lines := strings.Split(text, "\n")
	inApps, app := false, ""
	var changes []string
	applied := map[string]bool{}
	for i, l := range lines {
		if strings.HasPrefix(l, "applications:") {
			inApps = true
			continue
		}
		if inApps && len(l) > 0 && l[0] != ' ' && l[0] != '-' && l[0] != '#' {
			inApps = false
		}
		if !inApps {
			continue
		}
		if mm := nameLine.FindStringSubmatch(l); mm != nil && strings.HasPrefix(strings.TrimLeft(l, " "), "- name:") && indent(l) <= 2 {
			app = mm[1]
			continue
		}
		mm := tagLine.FindStringSubmatch(l)
		if mm == nil || app == "" || applied[app] {
			continue
		}
		next, ok := targets[app]
		if !ok {
			continue
		}
		lines[i] = fmt.Sprintf("%stag: %q%s", mm[1], next, mm[3])
		changes = append(changes, fmt.Sprintf("`%s`: `%s` %s → %s", app, images[app], mm[2], next))
		applied[app] = true
	}
	for app := range targets {
		if !applied[app] {
			return "", nil, fmt.Errorf("could not find the tag line of application %q", app)
		}
	}
	sort.Strings(changes)
	return strings.Join(lines, "\n"), changes, nil
}

func indent(l string) int { return len(l) - len(strings.TrimLeft(l, " ")) }

func setManifestVersion(text, from, to string) (string, error) {
	re := regexp.MustCompile(`(?m)^(  version:\s*["']?)` + regexp.QuoteMeta(from) + `(["']?\s*(#.*)?)$`)
	loc := re.FindStringSubmatchIndex(text)
	if loc == nil {
		return "", fmt.Errorf("metadata.version %s not found", from)
	}
	return text[:loc[3]] + to + text[loc[3]+len(from):], nil
}

func resolveDigests(ctx context.Context, reg *registry.Client, m *manifest.Manifest) (map[string]string, error) {
	out := map[string]string{}
	for _, a := range m.Applications {
		d, err := reg.Digest(ctx, a.Image, a.Tag)
		if err != nil {
			return nil, fmt.Errorf("resolve %s:%s: %w", a.Image, a.Tag, err)
		}
		out[a.Name] = d
	}
	return out, nil
}

func writeRelease(dir string, r *catalog.Release) error {
	data, err := catalog.MarshalRelease(r)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, catalog.ReleaseFile), data, 0o644)
}

// lock backfills the latest version's release file; a git-less checkout falls back to today's date.
func lock(args []string, today string) error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	ts, err := selectTemplates(cat, args)
	if err != nil {
		return err
	}
	ctx := context.Background()
	reg := registry.New()
	for _, t := range ts {
		v := t.Latest()
		rel := v.Info()
		dir := filepath.Join(t.Source, t.Name, v.Version)
		changed := false
		if rel.Released == "" {
			rel.Released = gitAddedDate(dir, today)
			changed = true
		}
		if len(rel.Images) < len(v.Manifest.Applications) {
			images, err := resolveDigests(ctx, reg, v.Manifest)
			if err != nil {
				logger.Warn("digests not locked", "template", t.Name, "version", v.Version, "error", err)
			} else {
				rel.Images = images
				changed = true
			}
		}
		if !changed {
			continue
		}
		if err := writeRelease(dir, &rel); err != nil {
			return err
		}
		logger.Info("locked", "template", t.Name, "version", v.Version, "released", rel.Released, "images", len(rel.Images))
	}
	return nil
}

func gitAddedDate(dir, fallback string) string {
	out, err := exec.Command("git", "log", "--diff-filter=A", "--format=%as", "--", filepath.Join(dir, "template.yaml")).Output()
	if err != nil {
		return fallback
	}
	lines := strings.Fields(strings.TrimSpace(string(out)))
	if len(lines) == 0 {
		return fallback
	}
	return lines[len(lines)-1]
}

type imageCheck struct {
	Template string `json:"template"`
	Version  string `json:"version"`
	App      string `json:"app"`
	Ref      string `json:"ref"`
	Digest   string `json:"digest,omitempty"`
	Locked   string `json:"locked,omitempty"`
	Problem  string `json:"problem,omitempty"`
	Moved    bool   `json:"moved,omitempty"`
	image    string
	tag      string
}

// checkImages fails on an image that no longer resolves; a tag that moved off its lock is only reported.
func checkImages(args []string, strict bool) error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	ts, err := selectTemplates(cat, args)
	if err != nil {
		return err
	}
	ctx := context.Background()
	reg := registry.New()
	var checks []imageCheck
	for _, t := range ts {
		v := t.Installable()
		for _, a := range v.Manifest.Applications {
			checks = append(checks, imageCheck{Template: t.Name, Version: v.Version, App: a.Name, Ref: a.Image + ":" + a.Tag,
				Locked: v.Info().Images[a.Name], image: a.Image, tag: a.Tag})
		}
	}
	parallel(len(checks), 8, func(i int) {
		c := &checks[i]
		d, err := reg.Digest(ctx, c.image, c.tag)
		if err != nil {
			c.Problem = err.Error()
			return
		}
		c.Digest = d
		c.Moved = c.Locked != "" && c.Locked != d
	})
	var broken, moved int
	for _, c := range checks {
		switch {
		case c.Problem != "":
			broken++
			logger.Error("image does not resolve", "template", c.Template, "version", c.Version, "app", c.App, "ref", c.Ref, "error", c.Problem)
		case c.Moved:
			moved++
			logger.Warn("tag moved off its locked digest", "template", c.Template, "version", c.Version, "app", c.App, "ref", c.Ref, "locked", c.Locked, "now", c.Digest)
		}
	}
	logger.Info("image check finished", "images", len(checks), "broken", broken, "moved", moved)
	if broken > 0 || (strict && moved > 0) {
		return fmt.Errorf("%d image(s) do not resolve, %d moved", broken, moved)
	}
	return nil
}

func parallel(n, workers int, fn func(i int)) {
	var wg sync.WaitGroup
	ch := make(chan int)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range ch {
				fn(i)
			}
		}()
	}
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
}

// lintWarnings are gaps that don't break an install; --strict turns them into failures.
func lintWarnings(cat *catalog.Catalog) []string {
	var out []string
	for _, t := range cat.Templates() {
		if t.Meta.Deprecated {
			continue
		}
		where := t.Source + "/" + t.Name
		if t.Source == catalog.SourceOfficial && strings.TrimSpace(t.Readme) == "" {
			out = append(out, where+": no README.md (the storefront detail page has only the one-line description)")
		}
		v := t.Latest()
		rel := v.Info()
		if v.Release == nil {
			out = append(out, fmt.Sprintf("%s@%s: no %s (run `marketplace lock %s`)", where, v.Version, catalog.ReleaseFile, t.Name))
			continue
		}
		if len(rel.Images) < len(v.Manifest.Applications) {
			out = append(out, fmt.Sprintf("%s@%s: image digests not locked (run `marketplace lock %s`)", where, v.Version, t.Name))
		}
		if rel.Breaking && strings.TrimSpace(rel.Notes) == "" {
			out = append(out, fmt.Sprintf("%s@%s: breaking release without notes", where, v.Version))
		}
	}
	return out
}
