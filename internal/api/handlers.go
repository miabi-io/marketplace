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

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/manifest"
)

// Handlers serves the API over a loaded catalog.
type Handlers struct{ cat *catalog.Catalog }

func Register(app *okapi.Okapi, cat *catalog.Catalog) {
	h := &Handlers{cat: cat}

	app.Get("/v1/export", h.Export,
		okapi.DocSummary("Export the full catalog bundle"),
		okapi.DocDescription("The full bundle (index + every version's manifest, inline) — Miabi's primary sync call. Bare JSON (no envelope); ETag-conditional, returns 304 when If-None-Match matches."),
		okapi.DocTag("catalog"),
		okapi.DocResponseHeader("ETag", "string", "Catalog digest; echo back as If-None-Match for a conditional GET."),
		okapi.DocResponse(http.StatusOK, &catalog.Bundle{}),
		okapi.DocResponse(http.StatusNotModified, nil),
	)

	app.Get("/v1/index", h.Index,
		okapi.DocSummary("Lightweight machine index"),
		okapi.DocDescription("The lightweight machine index (sources, versions, digests). Bare JSON; ETag-conditional."),
		okapi.DocTag("catalog"),
		okapi.DocResponseHeader("ETag", "string", "Catalog digest."),
		okapi.DocResponse(http.StatusOK, &catalog.Index{}),
		okapi.DocResponse(http.StatusNotModified, nil),
	)

	app.Get("/v1/templates", okapi.H(h.ListTemplates),
		okapi.DocSummary("Search & list templates"),
		okapi.DocDescription("Paginated search/filter over the catalog. Offset pagination: page (1-based) + per_page (default 24, max 100)."),
		okapi.DocTag("templates"),
		okapi.DocResponse(http.StatusOK, &Envelope[catalog.Page]{}),
	).WithInput(&ListTemplatesRequest{})

	app.Get("/v1/templates/{name}", okapi.H(h.GetTemplate),
		okapi.DocSummary("Template detail"),
		okapi.DocDescription("A template's detail: listing, all versions, README, and the manifest of the requested version (?version, default latest)."),
		okapi.DocTag("templates"),
		okapi.DocResponse(http.StatusOK, &Envelope[Detail]{}),
		okapi.DocErrorResponse(http.StatusNotFound, &ErrorResponse{}),
	).WithInput(&TemplateRequest{})

	app.Get("/v1/templates/{name}/versions/{version}", okapi.H(h.GetVersion),
		okapi.DocSummary("Version metadata"),
		okapi.DocDescription("One version's metadata + digest."),
		okapi.DocTag("templates"),
		okapi.DocResponse(http.StatusOK, &Envelope[VersionDetail]{}),
		okapi.DocErrorResponse(http.StatusNotFound, &ErrorResponse{}),
	).WithInput(&VersionRequest{})

	app.Get("/v1/templates/{name}/versions/{version}/manifest", okapi.H(h.GetManifest),
		okapi.DocSummary("Raw template manifest"),
		okapi.DocDescription("The raw template.yaml to install; the version digest is the ETag. Content-Type application/yaml; ETag-conditional."),
		okapi.DocTag("templates"),
		okapi.DocResponseHeader("ETag", "string", "Version digest."),
		okapi.DocResponse(http.StatusOK, ""),
		okapi.DocResponse(http.StatusNotModified, nil),
		okapi.DocErrorResponse(http.StatusNotFound, &ErrorResponse{}),
	).WithInput(&VersionRequest{})

	app.Get("/v1/categories", h.Categories,
		okapi.DocSummary("Category facets"),
		okapi.DocDescription("The category facets (name + count)."),
		okapi.DocTag("catalog"),
		okapi.DocResponse(http.StatusOK, &Envelope[[]catalog.CategoryFacet]{}),
	)

	app.Get("/healthz", h.Health,
		okapi.DocSummary("Liveness probe"),
		okapi.DocTag("ops"),
		okapi.DocResponse(http.StatusOK, &HealthResponse{}),
	)

	app.Get("/metrics", h.Metrics,
		okapi.DocSummary("Prometheus metrics"),
		okapi.DocDescription("Prometheus text exposition of catalog gauges."),
		okapi.DocTag("ops"),
		okapi.DocResponse(http.StatusOK, ""),
	)
}

// Export returns the full bundle (Miabi's primary sync call). Bare JSON (no
// envelope) so the consumer decodes it directly; ETag-conditional.
func (h *Handlers) Export(c *okapi.Context) error {
	if h.notModified(c, h.cat.ETag()) {
		return c.String(http.StatusNotModified, "")
	}
	return c.JSON(http.StatusOK, h.cat.Bundle())
}

// Index returns the lightweight machine index (sources, versions, digests).
func (h *Handlers) Index(c *okapi.Context) error {
	if h.notModified(c, h.cat.ETag()) {
		return c.String(http.StatusNotModified, "")
	}
	return c.JSON(http.StatusOK, h.cat.Index())
}

// ListTemplatesRequest is the search/filter query for ListTemplates. Okapi binds
// it from the query string, applies defaults/validation, and documents it.
type ListTemplatesRequest struct {
	Q        string `query:"q" description:"Free-text search over name/description/tags."`
	Source   string `query:"source" enum:"official,community" description:"Filter by source."`
	Category string `query:"category" description:"Filter by category name."`
	Tag      string `query:"tag" description:"Filter by tag."`
	Sort     string `query:"sort" enum:"name,updated,popularity" description:"Sort order (updated/popularity fall back to name today)."`
	Page     int    `query:"page" default:"1" description:"1-based page number."`
	PerPage  int    `query:"per_page" default:"24" description:"Items per page (max 100; larger values are clamped)."`
}

// ListTemplates is the paginated search/filter endpoint (storefront + clients).
func (h *Handlers) ListTemplates(c *okapi.Context, in *ListTemplatesRequest) error {
	page := h.cat.Search(catalog.Query{
		Q:        in.Q,
		Source:   in.Source,
		Category: in.Category,
		Tag:      in.Tag,
		Sort:     in.Sort,
		Page:     in.Page,
		PerPage:  in.PerPage,
	})
	return ok(c, page)
}

// TemplateRequest identifies a template (name) and an optional version.
type TemplateRequest struct {
	Name    string `path:"name" required:"true" description:"Template name."`
	Version string `query:"version" description:"Specific version (default: latest)."`
}

// GetTemplate returns a template's detail: listing, all versions, README, and
// the manifest of the requested version (?version, default latest).
func (h *Handlers) GetTemplate(c *okapi.Context, in *TemplateRequest) error {
	t, found := h.cat.Get(in.Name)
	if !found {
		return fail(c, http.StatusNotFound, "TEMPLATE_NOT_FOUND", "template not found")
	}
	ver, vok := t.FindVersion(in.Version)
	if !vok {
		return fail(c, http.StatusNotFound, "VERSION_NOT_FOUND", "template version not found")
	}
	return ok(c, detailOf(t, ver))
}

// VersionRequest identifies a specific template version by name + version path.
type VersionRequest struct {
	Name    string `path:"name" required:"true" description:"Template name."`
	Version string `path:"version" required:"true" description:"Version identifier."`
}

// GetVersion returns one version's metadata + digest.
func (h *Handlers) GetVersion(c *okapi.Context, in *VersionRequest) error {
	t, found := h.cat.Get(in.Name)
	if !found {
		return fail(c, http.StatusNotFound, "TEMPLATE_NOT_FOUND", "template not found")
	}
	ver, vok := t.FindVersion(in.Version)
	if !vok {
		return fail(c, http.StatusNotFound, "VERSION_NOT_FOUND", "template version not found")
	}
	return ok(c, VersionDetail{Name: t.Name, Slug: t.Name, Source: t.Source, Version: ver.Version, Digest: ver.Digest, Metadata: ver.Manifest.Metadata})
}

// GetManifest returns the raw template.yaml to install; the digest is the ETag.
func (h *Handlers) GetManifest(c *okapi.Context, in *VersionRequest) error {
	t, found := h.cat.Get(in.Name)
	if !found {
		return fail(c, http.StatusNotFound, "TEMPLATE_NOT_FOUND", "template not found")
	}
	ver, vok := t.FindVersion(in.Version)
	if !vok {
		return fail(c, http.StatusNotFound, "VERSION_NOT_FOUND", "template version not found")
	}
	if h.notModified(c, `"`+ver.Digest+`"`) {
		return c.String(http.StatusNotModified, "")
	}
	c.SetHeader("Content-Type", "application/yaml; charset=utf-8")
	return c.String(http.StatusOK, string(ver.Raw))
}

// Categories returns the category facets.
func (h *Handlers) Categories(c *okapi.Context) error {
	return ok(c, h.cat.Categories())
}

// HealthResponse is the liveness probe body.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health is the liveness probe.
func (h *Handlers) Health(c *okapi.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

// Metrics exposes a small Prometheus text exposition (dependency-free).
func (h *Handlers) Metrics(c *okapi.Context) error {
	counts := map[string]int{}
	versions := 0
	for _, t := range h.cat.Templates() {
		counts[t.Source]++
		versions += len(t.Versions)
	}
	var b strings.Builder
	b.WriteString("# HELP marketplace_templates_total Number of templates by source.\n")
	b.WriteString("# TYPE marketplace_templates_total gauge\n")
	for _, s := range []string{catalog.SourceOfficial, catalog.SourceCommunity} {
		fmt.Fprintf(&b, "marketplace_templates_total{source=%q} %d\n", s, counts[s])
	}
	b.WriteString("# HELP marketplace_versions_total Number of published template versions.\n")
	b.WriteString("# TYPE marketplace_versions_total gauge\n")
	fmt.Fprintf(&b, "marketplace_versions_total %d\n", versions)
	c.SetHeader("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	return c.String(http.StatusOK, b.String())
}

// notModified sets caching headers and reports whether the client's
// If-None-Match already matches (a 304 short-circuit).
func (h *Handlers) notModified(c *okapi.Context, etag string) bool {
	c.SetHeader("ETag", etag)
	c.SetHeader("Cache-Control", "public, max-age=300")
	return c.Header("If-None-Match") == etag
}

// Detail is the template detail document.
type Detail struct {
	Entry    catalog.Listing    `json:"entry"`
	Source   string             `json:"source"`
	Featured bool               `json:"featured,omitempty"`
	Meta     catalog.Metadata   `json:"meta"`
	Readme   string             `json:"readme,omitempty"`
	Versions []VersionRef       `json:"versions"`
	Manifest *manifest.Manifest `json:"manifest"`
}

// VersionRef is a lightweight version pointer in a detail document.
type VersionRef struct {
	Version string `json:"version"`
	Digest  string `json:"digest"`
}

// VersionDetail is one version's metadata + digest.
type VersionDetail struct {
	Name     string            `json:"name"`
	Slug     string            `json:"slug"` // deprecated: alias of Name, for pre-rename consumers
	Source   string            `json:"source"`
	Version  string            `json:"version"`
	Digest   string            `json:"digest"`
	Metadata manifest.Metadata `json:"metadata"`
}

func detailOf(t *catalog.Template, ver catalog.Version) Detail {
	refs := make([]VersionRef, 0, len(t.Versions))
	for _, v := range t.Versions {
		refs = append(refs, VersionRef{Version: v.Version, Digest: v.Digest})
	}
	return Detail{
		Entry:    t.Listing(),
		Source:   t.Source,
		Featured: t.Meta.Featured,
		Meta:     t.Meta,
		Readme:   t.Readme,
		Versions: refs,
		Manifest: ver.Manifest,
	}
}
