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

// Package catalog loads the embedded templates into an in-memory, digest-verified catalog served by the API.
package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	marketplace "github.com/miabi-io/marketplace"
	"github.com/miabi-io/marketplace/manifest"
	"gopkg.in/yaml.v3"
)

// Source folders / labels (carried through the index + API, badged in the UI).
const (
	SourceOfficial  = "official"
	SourceCommunity = "community"
)

var nameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// Metadata is storefront-only enrichment; it never affects an install.
type Metadata struct {
	Featured    bool     `yaml:"featured" json:"featured,omitempty"`
	Screenshots []string `yaml:"screenshots" json:"screenshots,omitempty"`
	SourceRepo  string   `yaml:"sourceRepo" json:"source_repo,omitempty"`
	Maintainer  string   `yaml:"maintainer" json:"maintainer,omitempty"`
	// Deprecated retires the whole template: hidden from browsing and new installs, existing installs keep running.
	Deprecated      bool   `yaml:"deprecated" json:"deprecated,omitempty"`
	DeprecationNote string `yaml:"deprecationNote" json:"deprecation_note,omitempty"`
	// ReplacedBy names the template to suggest instead.
	ReplacedBy string `yaml:"replacedBy" json:"replaced_by,omitempty"`
	// Updates steers `outdated`/`bump` per application, keyed by template app name.
	Updates map[string]UpdatePolicy `yaml:"updates" json:"-"`
}

// UpdatePolicy records a deliberate pin so the weekly bump doesn't keep proposing it.
type UpdatePolicy struct {
	// Hold is the reason the app's tag is frozen.
	Hold string `yaml:"hold"`
	// Match restricts candidate tags to a regexp, e.g. to stay on one PHP variant.
	Match string `yaml:"match"`
}

// Version is one immutable published version of a template.
type Version struct {
	Version  string
	Digest   string // sha256 of Raw, "sha256:…"
	Raw      []byte // the raw template.yaml
	Manifest *manifest.Manifest
	Release  *Release // nil when the version has no release.yaml
}

func (v Version) Info() Release {
	if v.Release == nil {
		return Release{}
	}
	return *v.Release
}

// Template is one catalog entry with all of its versions (newest-first).
type Template struct {
	Name     string
	Source   string
	Meta     Metadata
	Readme   string
	Versions []Version
	// Freshness is the latest `outdated` report's entry, nil when none covers the latest version.
	Freshness *TemplateFreshness
}

// Latest returns the newest version (Versions is sorted newest-first).
func (t *Template) Latest() Version { return t.Versions[0] }

// Installable returns the newest version a new install should get: the latest not deprecated, else the latest.
func (t *Template) Installable() Version {
	for _, v := range t.Versions {
		if !v.Info().Deprecated {
			return v
		}
	}
	return t.Versions[0]
}

// FindVersion returns a specific version (empty = latest).
func (t *Template) FindVersion(v string) (Version, bool) {
	if v == "" {
		return t.Versions[0], true
	}
	for _, ver := range t.Versions {
		if ver.Version == v {
			return ver, true
		}
	}
	return Version{}, false
}

// Catalog is the loaded, immutable set of templates.
type Catalog struct {
	templates   []Template
	byName      map[string]*Template
	etag        string
	generatedAt string
}

// Load reads the embedded content into a catalog.
func Load() (*Catalog, error) {
	c, err := loadFS(marketplace.Content)
	if err != nil {
		return nil, err
	}
	r, err := ParseFreshness(marketplace.Freshness)
	if err != nil {
		return nil, err
	}
	c.attachFreshness(r)
	return c, nil
}

// Any malformed manifest, digest/name/version mismatch or cross-source duplicate name fails the build closed.
func loadFS(fsys fs.FS) (*Catalog, error) {
	c := &Catalog{}
	seen := map[string]string{} // name -> source, for cross-folder uniqueness
	for _, source := range []string{SourceOfficial, SourceCommunity} {
		entries, err := fs.ReadDir(fsys, source)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if other, ok := seen[name]; ok {
				return nil, fmt.Errorf("duplicate template %q in %s/ and %s/", name, other, source)
			}
			t, err := loadTemplate(fsys, source, name)
			if err != nil {
				return nil, fmt.Errorf("%s/%s: %w", source, name, err)
			}
			if t == nil {
				continue // a template dir with no versions is ignored
			}
			seen[name] = source
			c.templates = append(c.templates, *t)
		}
	}

	// Stable order: official before community, then by display name.
	sort.SliceStable(c.templates, func(i, j int) bool {
		a, b := c.templates[i], c.templates[j]
		if a.Source != b.Source {
			return a.Source == SourceOfficial
		}
		return strings.ToLower(a.Latest().Manifest.Metadata.DisplayName) < strings.ToLower(b.Latest().Manifest.Metadata.DisplayName)
	})

	c.byName = make(map[string]*Template, len(c.templates))
	for i := range c.templates {
		c.byName[c.templates[i].Name] = &c.templates[i]
	}
	if err := c.crossCheck(); err != nil {
		return nil, err
	}
	c.etag = computeETag(c.templates)
	c.generatedAt = time.Now().UTC().Format(time.RFC3339)
	return c, nil
}

func loadTemplate(fsys fs.FS, source, name string) (*Template, error) {
	if !nameRe.MatchString(name) {
		return nil, fmt.Errorf("invalid template directory %q", name)
	}
	base := path.Join(source, name)
	entries, err := fs.ReadDir(fsys, base)
	if err != nil {
		return nil, err
	}
	t := &Template{Name: name, Source: source}
	if b, err := fs.ReadFile(fsys, path.Join(base, "metadata.yaml")); err == nil {
		if err := yaml.Unmarshal(b, &t.Meta); err != nil {
			return nil, fmt.Errorf("metadata.yaml: %w", err)
		}
	}
	if b, err := fs.ReadFile(fsys, path.Join(base, "README.md")); err == nil {
		t.Readme = string(b)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		ver := e.Name()
		raw, err := fs.ReadFile(fsys, path.Join(base, ver, "template.yaml"))
		if err != nil {
			return nil, fmt.Errorf("version %s: %w", ver, err)
		}
		m, err := manifest.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("version %s: %w", ver, err)
		}
		var rel *Release
		if b, err := fs.ReadFile(fsys, path.Join(base, ver, ReleaseFile)); err == nil {
			if rel, err = parseRelease(b); err != nil {
				return nil, fmt.Errorf("version %s: %w", ver, err)
			}
		}
		if m.Metadata.Name != name {
			return nil, fmt.Errorf("version %s: manifest name %q disagrees with directory %q", ver, m.Metadata.Name, name)
		}
		if m.Metadata.Version != ver {
			return nil, fmt.Errorf("version directory %q disagrees with manifest version %q", ver, m.Metadata.Version)
		}
		t.Versions = append(t.Versions, Version{Version: ver, Digest: manifest.Digest(raw), Raw: raw, Manifest: m, Release: rel})
	}
	if len(t.Versions) == 0 {
		return nil, nil
	}
	sort.SliceStable(t.Versions, func(i, j int) bool {
		return manifest.CompareVersions(t.Versions[i].Version, t.Versions[j].Version) > 0
	})
	return t, nil
}

// crossCheck validates references that span versions or templates, which a single file can't check.
func (c *Catalog) crossCheck() error {
	for _, t := range c.templates {
		where := t.Source + "/" + t.Name
		for app, p := range t.Meta.Updates {
			if p.Match != "" {
				if _, err := regexp.Compile(p.Match); err != nil {
					return fmt.Errorf("%s: updates.%s.match: %w", where, app, err)
				}
			}
			if !hasApp(t.Latest().Manifest, app) {
				return fmt.Errorf("%s: updates names unknown application %q", where, app)
			}
		}
		if r := t.Meta.ReplacedBy; r != "" {
			if r == t.Name {
				return fmt.Errorf("%s: replacedBy names the template itself", where)
			}
			if _, ok := c.byName[r]; !ok {
				return fmt.Errorf("%s: replacedBy %q is not a template in this catalog", where, r)
			}
		}
		for _, v := range t.Versions {
			rel := v.Info()
			if rf := rel.RequiresFrom; rf != "" {
				if _, ok := t.FindVersion(rf); !ok {
					return fmt.Errorf("%s@%s: requiresFrom %q is not a published version", where, v.Version, rf)
				}
				if manifest.CompareVersions(rf, v.Version) >= 0 {
					return fmt.Errorf("%s@%s: requiresFrom %q must be older than the version itself", where, v.Version, rf)
				}
			}
			apps := map[string]bool{}
			for _, a := range v.Manifest.Applications {
				apps[a.Name] = true
			}
			for app := range rel.Images {
				if !apps[app] {
					return fmt.Errorf("%s@%s: %s pins an image for unknown application %q", where, v.Version, ReleaseFile, app)
				}
			}
		}
	}
	return nil
}

func hasApp(m *manifest.Manifest, name string) bool {
	for _, a := range m.Applications {
		if a.Name == name {
			return true
		}
	}
	return false
}

// computeETag changes iff catalog content changes, so conditional GETs on /export and /index are nearly free.
func computeETag(ts []Template) string {
	lines := make([]string, 0, len(ts))
	for _, t := range ts {
		for _, v := range t.Versions {
			line := t.Source + "/" + t.Name + "@" + v.Version + "=" + v.Digest
			if v.Release != nil {
				if b, err := MarshalRelease(v.Release); err == nil {
					sum := sha256.Sum256(b)
					line += "+" + hex.EncodeToString(sum[:8])
				}
			}
			lines = append(lines, line)
		}
		if t.Meta.Deprecated {
			lines = append(lines, t.Source+"/"+t.Name+"!deprecated:"+t.Meta.ReplacedBy+":"+t.Meta.DeprecationNote)
		}
	}
	sort.Strings(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

// Templates returns all templates (official first, then by name).
func (c *Catalog) Templates() []Template { return c.templates }

// Get returns a template by name.
func (c *Catalog) Get(name string) (*Template, bool) {
	t, ok := c.byName[name]
	return t, ok
}

// ETag returns the catalog-wide strong validator (quoted).
func (c *Catalog) ETag() string { return c.etag }

// GeneratedAt returns the load time (RFC3339).
func (c *Catalog) GeneratedAt() string { return c.generatedAt }
