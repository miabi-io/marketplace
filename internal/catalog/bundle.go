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

package catalog

import (
	"sort"
	"strings"

	"github.com/miabi-io/marketplace/manifest"
)

// Bundle is the GET /v1/export document: the entire catalog (both sources, all
// versions, manifests inline) in one JSON payload. This shape is the contract
// the Miabi consumer decodes (its remote.Bundle) — keep the JSON tags in sync.
type Bundle struct {
	ETag        string           `json:"etag"`
	GeneratedAt string           `json:"generatedAt"`
	Templates   []BundleTemplate `json:"templates"`
}

type BundleTemplate struct {
	Name     string          `json:"name"`
	Source   string          `json:"source"`
	Versions []BundleVersion `json:"versions"`
}

type BundleVersion struct {
	Version  string `json:"version"`
	Digest   string `json:"digest"`
	Manifest string `json:"manifest"` // raw template.yaml
}

// Bundle builds the full export document.
func (c *Catalog) Bundle() Bundle {
	b := Bundle{ETag: c.etag, GeneratedAt: c.generatedAt, Templates: make([]BundleTemplate, 0, len(c.templates))}
	for _, t := range c.templates {
		bt := BundleTemplate{Name: t.Name, Source: t.Source}
		for _, v := range t.Versions {
			bt.Versions = append(bt.Versions, BundleVersion{Version: v.Version, Digest: v.Digest, Manifest: string(v.Raw)})
		}
		b.Templates = append(b.Templates, bt)
	}
	return b
}

// Index is the lightweight GET /v1/index document: sources, versions, and
// digests, without manifests — for cheap freshness checks.
type Index struct {
	ETag        string          `json:"etag"`
	GeneratedAt string          `json:"generatedAt"`
	Templates   []IndexTemplate `json:"templates"`
}

type IndexTemplate struct {
	Name        string         `json:"name"`
	Source      string         `json:"source"`
	DisplayName string         `json:"displayName"`
	Category    string         `json:"category,omitempty"`
	Versions    []IndexVersion `json:"versions"`
}

type IndexVersion struct {
	Version string `json:"version"`
	Digest  string `json:"digest"`
}

// Index builds the lightweight machine index.
func (c *Catalog) Index() Index {
	idx := Index{ETag: c.etag, GeneratedAt: c.generatedAt, Templates: make([]IndexTemplate, 0, len(c.templates))}
	for _, t := range c.templates {
		it := IndexTemplate{Name: t.Name, Source: t.Source, DisplayName: t.Latest().Manifest.Metadata.DisplayName, Category: t.Latest().Manifest.Metadata.Category}
		for _, v := range t.Versions {
			it.Versions = append(it.Versions, IndexVersion{Version: v.Version, Digest: v.Digest})
		}
		idx.Templates = append(idx.Templates, it)
	}
	return idx
}

// Listing is the storefront card view of a template (its latest version).
type Listing struct {
	Name         string           `json:"name"`
	DisplayName  string           `json:"displayName"`
	Description  string           `json:"description"`
	Category     string           `json:"category"`
	Icon         string           `json:"icon,omitempty"`
	Tags         []string         `json:"tags,omitempty"`
	Homepage     string           `json:"homepage,omitempty"`
	Author       *manifest.Author `json:"author,omitempty"`
	Source       string           `json:"source"`
	Featured     bool             `json:"featured,omitempty"`
	Version      string           `json:"version"`
	Versions     []string         `json:"versions"`
	Applications int              `json:"applications"`
	Databases    int              `json:"databases"`
	Volumes      int              `json:"volumes"`
	DBOnly       bool             `json:"db_only"`
}

// Listing builds the card view of a template.
func (t *Template) Listing() Listing {
	m := t.Latest().Manifest
	vers := make([]string, 0, len(t.Versions))
	for _, v := range t.Versions {
		vers = append(vers, v.Version)
	}
	return Listing{
		Name: t.Name, DisplayName: m.Metadata.DisplayName, Description: m.Metadata.Description,
		Category: m.Metadata.Category, Icon: m.Metadata.Icon, Tags: m.Metadata.Tags,
		Homepage: m.Metadata.Homepage, Author: m.Metadata.Author, Source: t.Source,
		Featured: t.Meta.Featured, Version: m.Metadata.Version, Versions: vers,
		Applications: len(m.Applications), Databases: len(m.Databases), Volumes: len(m.Volumes),
		DBOnly: m.IsDatabaseOnly(),
	}
}

// CategoryFacet is one entry of GET /v1/categories.
type CategoryFacet struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// Categories returns the category facets (by latest version), sorted by name.
func (c *Catalog) Categories() []CategoryFacet {
	counts := map[string]int{}
	for _, t := range c.templates {
		cat := t.Latest().Manifest.Metadata.Category
		if cat == "" {
			cat = "Uncategorized"
		}
		counts[cat]++
	}
	out := make([]CategoryFacet, 0, len(counts))
	for cat, n := range counts {
		out = append(out, CategoryFacet{Category: cat, Count: n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Category < out[j].Category })
	return out
}

// Query parameterizes a paginated search over the catalog.
type Query struct {
	Q        string
	Source   string
	Category string
	Tag      string
	Sort     string // name | updated | popularity (updated/popularity fall back to name today)
	Page     int
	PerPage  int
}

// Page is one page of search results.
type Page struct {
	Items      []Listing `json:"items"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
}

// Search filters, sorts, and paginates the catalog (offset pagination).
func (c *Catalog) Search(q Query) Page {
	if q.PerPage <= 0 {
		q.PerPage = 24
	}
	if q.PerPage > 100 {
		q.PerPage = 100
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	needle := strings.ToLower(strings.TrimSpace(q.Q))

	var matched []Listing
	for _, t := range c.templates {
		if q.Source != "" && t.Source != q.Source {
			continue
		}
		l := t.Listing()
		if q.Category != "" && !strings.EqualFold(l.Category, q.Category) {
			continue
		}
		if q.Tag != "" && !containsFold(l.Tags, q.Tag) {
			continue
		}
		if needle != "" && !matchesText(l, needle) {
			continue
		}
		matched = append(matched, l)
	}

	// Featured first within the (already name-ordered) set, unless sorting by name.
	if q.Sort != "name" {
		sort.SliceStable(matched, func(i, j int) bool { return matched[i].Featured && !matched[j].Featured })
	}

	total := len(matched)
	totalPages := (total + q.PerPage - 1) / q.PerPage
	if totalPages == 0 {
		totalPages = 1
	}
	start := (q.Page - 1) * q.PerPage
	if start > total {
		start = total
	}
	end := start + q.PerPage
	if end > total {
		end = total
	}
	return Page{
		Items:      append([]Listing{}, matched[start:end]...),
		Page:       q.Page,
		PerPage:    q.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func matchesText(l Listing, needle string) bool {
	if strings.Contains(strings.ToLower(l.Name), needle) ||
		strings.Contains(strings.ToLower(l.DisplayName), needle) ||
		strings.Contains(strings.ToLower(l.Description), needle) ||
		strings.Contains(strings.ToLower(l.Category), needle) {
		return true
	}
	return containsFold(l.Tags, needle)
}

func containsFold(haystack []string, needle string) bool {
	for _, h := range haystack {
		if strings.Contains(strings.ToLower(h), strings.ToLower(needle)) {
			return true
		}
	}
	return false
}
