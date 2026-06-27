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

// Package catalog loads the embedded official/ + community/ template content
// into an in-memory, digest-verified catalog and serves the wire views the API
// exposes (export bundle, index, listings). It is the service's source of truth;
// the service is otherwise stateless (git is the database).
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

var slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// Metadata is storefront-only enrichment from <slug>/metadata.yaml — it never
// affects an install (the manifest is authoritative for that).
type Metadata struct {
	Featured    bool     `yaml:"featured" json:"featured,omitempty"`
	Screenshots []string `yaml:"screenshots" json:"screenshots,omitempty"`
	SourceRepo  string   `yaml:"sourceRepo" json:"source_repo,omitempty"`
	Maintainer  string   `yaml:"maintainer" json:"maintainer,omitempty"`
}

// Version is one immutable published version of a template.
type Version struct {
	Version  string
	Digest   string // sha256 of Raw, "sha256:…"
	Raw      []byte // the raw template.yaml
	Manifest *manifest.Manifest
}

// Template is one catalog entry with all of its versions (newest-first).
type Template struct {
	Slug     string
	Source   string
	Meta     Metadata
	Readme   string
	Versions []Version
}

// Latest returns the newest version (Versions is sorted newest-first).
func (t *Template) Latest() Version { return t.Versions[0] }

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
	bySlug      map[string]*Template
	etag        string
	generatedAt string
}

// Load reads the embedded content into a catalog.
func Load() (*Catalog, error) { return loadFS(marketplace.Content) }

// loadFS builds a catalog from any fs (embedded in production; a fixture in
// tests). A malformed manifest, a digest/slug/version mismatch, or a slug
// duplicated across sources is a hard error — the build/CI fails closed.
func loadFS(fsys fs.FS) (*Catalog, error) {
	c := &Catalog{}
	seen := map[string]string{} // slug -> source, for cross-folder uniqueness
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
			slug := e.Name()
			if other, ok := seen[slug]; ok {
				return nil, fmt.Errorf("duplicate slug %q in %s/ and %s/", slug, other, source)
			}
			t, err := loadTemplate(fsys, source, slug)
			if err != nil {
				return nil, fmt.Errorf("%s/%s: %w", source, slug, err)
			}
			if t == nil {
				continue // a slug dir with no versions is ignored
			}
			seen[slug] = source
			c.templates = append(c.templates, *t)
		}
	}

	// Stable order: official before community, then by display name.
	sort.SliceStable(c.templates, func(i, j int) bool {
		a, b := c.templates[i], c.templates[j]
		if a.Source != b.Source {
			return a.Source == SourceOfficial
		}
		return strings.ToLower(a.Latest().Manifest.Metadata.Name) < strings.ToLower(b.Latest().Manifest.Metadata.Name)
	})

	c.bySlug = make(map[string]*Template, len(c.templates))
	for i := range c.templates {
		c.bySlug[c.templates[i].Slug] = &c.templates[i]
	}
	c.etag = computeETag(c.templates)
	c.generatedAt = time.Now().UTC().Format(time.RFC3339)
	return c, nil
}

func loadTemplate(fsys fs.FS, source, slug string) (*Template, error) {
	if !slugRe.MatchString(slug) {
		return nil, fmt.Errorf("invalid slug directory %q", slug)
	}
	base := path.Join(source, slug)
	entries, err := fs.ReadDir(fsys, base)
	if err != nil {
		return nil, err
	}
	t := &Template{Slug: slug, Source: source}
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
		if m.Metadata.Slug != slug {
			return nil, fmt.Errorf("version %s: manifest slug %q disagrees with directory %q", ver, m.Metadata.Slug, slug)
		}
		if m.Metadata.Version != ver {
			return nil, fmt.Errorf("version directory %q disagrees with manifest version %q", ver, m.Metadata.Version)
		}
		t.Versions = append(t.Versions, Version{Version: ver, Digest: manifest.Digest(raw), Raw: raw, Manifest: m})
	}
	if len(t.Versions) == 0 {
		return nil, nil
	}
	sort.SliceStable(t.Versions, func(i, j int) bool {
		return manifest.CompareVersions(t.Versions[i].Version, t.Versions[j].Version) > 0
	})
	return t, nil
}

// computeETag is a strong validator over every version digest — it changes iff
// the catalog content changes, so conditional GETs on /export and /index are
// nearly free.
func computeETag(ts []Template) string {
	lines := make([]string, 0, len(ts))
	for _, t := range ts {
		for _, v := range t.Versions {
			lines = append(lines, t.Source+"/"+t.Slug+"@"+v.Version+"="+v.Digest)
		}
	}
	sort.Strings(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

// Templates returns all templates (official first, then by name).
func (c *Catalog) Templates() []Template { return c.templates }

// Get returns a template by slug.
func (c *Catalog) Get(slug string) (*Template, bool) {
	t, ok := c.bySlug[slug]
	return t, ok
}

// ETag returns the catalog-wide strong validator (quoted).
func (c *Catalog) ETag() string { return c.etag }

// GeneratedAt returns the load time (RFC3339).
func (c *Catalog) GeneratedAt() string { return c.generatedAt }
