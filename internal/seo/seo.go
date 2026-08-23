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

// Package seo serves the crawler-facing surface of the storefront: per-page head
// tags, robots.txt and sitemap.xml.
//
// The storefront is a client-rendered SPA, so a crawler that indexes the HTML it
// is served — before, or instead of, running the JavaScript — sees the same shell
// with the same generic title for every URL. That is why a template page could be
// listed as a bare URL with no title. Rather than pre-render the whole site, the
// shell carries a marked-off block of head tags (see web/index.html) that this
// package rewrites per request: title, description, canonical, Open Graph and
// JSON-LD. What the browser ends up rendering is unchanged — the SPA still boots
// from the same shell and sets the same title client-side.
package seo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/fs"
	"net/http"
	"strings"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/internal/web"
)

const DefaultBaseURL = "https://marketplace.miabi.io"

const (
	markerStart = "<!-- seo:start -->"
	markerEnd   = "<!-- seo:end -->"
)

const (
	siteName    = "Miabi Marketplace"
	homeTitle   = "Miabi Marketplace — official & community app templates"
	homeDesc    = "Browse official and community Miabi templates: one-click apps, databases and stacks for the open-source PaaS for Docker."
	ogImagePath = "/og-image.png"
	// descMax keeps a meta description inside what a result snippet shows.
	descMax = 160
)

// Server renders the shell with per-page head tags.
type Server struct {
	cat  *catalog.Catalog
	base string

	prefix, suffix string
	shell          bool
}

func New(cat *catalog.Catalog, baseURL string) *Server {
	s := &Server{cat: cat, base: strings.TrimSuffix(baseURL, "/")}
	if s.base == "" {
		s.base = DefaultBaseURL
	}
	raw, err := fs.ReadFile(web.Assets, "dist/index.html")
	if err != nil {
		return s
	}
	doc := string(raw)
	s.shell = true
	start := strings.Index(doc, markerStart)
	end := strings.Index(doc, markerEnd)
	if start < 0 || end < start {
		s.prefix, s.suffix = doc, ""
		return s
	}
	s.prefix, s.suffix = doc[:start], doc[end+len(markerEnd):]
	return s
}

// Register mounts the crawler routes. They are registered before the SPA's file
// server so they win for the paths they own; every other path still falls through
// to the shell.
func Register(app *okapi.Okapi, cat *catalog.Catalog, baseURL string) {
	s := New(cat, baseURL)

	app.Get("/robots.txt", s.Robots, okapi.DocSummary("robots.txt"), okapi.DocTag("seo"))
	app.Get("/sitemap.xml", s.Sitemap, okapi.DocSummary("XML sitemap of the catalog"), okapi.DocTag("seo"))
	app.Get("/", s.Home, okapi.DocSummary("Storefront home"), okapi.DocTag("seo"))
	app.Get("/templates", s.TemplatesIndex, okapi.DocSummary("Template list (redirects home)"), okapi.DocTag("seo"))
	app.Get("/templates/{name}", s.Template, okapi.DocSummary("Template page"), okapi.DocTag("seo"))
}

func (s *Server) Robots(c *okapi.Context) error {
	body := "User-agent: *\nAllow: /\n\n" +
		"# Machine endpoints — useful to clients, noise in an index.\n" +
		"Disallow: /v1/\nDisallow: /openapi.json\nDisallow: /schema/\n\n" +
		"Sitemap: " + s.base + "/sitemap.xml\n"
	c.SetHeader("Cache-Control", "public, max-age=3600")
	return c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(body))
}

type urlEntry struct {
	Loc string `xml:"loc"`
}

type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	NS      string     `xml:"xmlns,attr"`
	URLs    []urlEntry `xml:"url"`
}

func (s *Server) Sitemap(c *okapi.Context) error {
	templates := s.cat.Templates()
	set := urlSet{
		NS:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: make([]urlEntry, 0, len(templates)+1),
	}
	set.URLs = append(set.URLs, urlEntry{Loc: s.base + "/"})
	for i := range templates {
		set.URLs = append(set.URLs, urlEntry{Loc: s.templateURL(templates[i].Name)})
	}
	body, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}
	c.SetHeader("Cache-Control", "public, max-age=3600")
	return c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(xml.Header+string(body)+"\n"))
}

// Home serves the shell with the site-level tags.
func (s *Server) Home(c *okapi.Context) error {
	page := pageMeta{
		Title:       homeTitle,
		Description: homeDesc,
		Canonical:   s.base + "/",
		OGType:      "website",
		JSONLD: map[string]any{
			"@context":    "https://schema.org",
			"@type":       "CollectionPage",
			"name":        siteName,
			"description": homeDesc,
			"url":         s.base + "/",
			"isPartOf": map[string]any{
				"@type": "WebSite",
				"name":  siteName,
				"url":   s.base + "/",
			},
		},
	}
	return s.render(c, page)
}

// TemplatesIndex redirects the bare list path home: the SPA has no route for it,
// and a redirect keeps a crawler that guessed the URL out of a soft 404.
func (s *Server) TemplatesIndex(c *okapi.Context) error {
	c.Redirect(http.StatusMovedPermanently, s.base+"/")
	return nil
}

func (s *Server) Template(c *okapi.Context) error {
	name := c.Param("name")
	t, found := s.cat.Get(name)
	if !found {
		return s.renderStatus(c, http.StatusNotFound, pageMeta{
			Title:       "Template not found — " + siteName,
			Description: homeDesc,
			Canonical:   s.base + "/",
			OGType:      "website",
			NoIndex:     true,
		})
	}
	l := t.Listing()
	m := t.Latest().Manifest
	url := s.templateURL(t.Name)
	desc := clamp(strings.TrimSpace(l.Description), descMax)
	if desc == "" {
		desc = fmt.Sprintf("Deploy %s on Miabi, the open-source PaaS for Docker.", l.DisplayName)
	}

	app := map[string]any{
		"@context":            "https://schema.org",
		"@type":               "SoftwareApplication",
		"name":                l.DisplayName,
		"description":         desc,
		"url":                 url,
		"applicationCategory": "DeveloperApplication",
		"operatingSystem":     "Linux, Docker",
		"softwareVersion":     l.Version,
		"offers": map[string]any{
			"@type": "Offer", "price": "0", "priceCurrency": "USD",
		},
	}
	if l.Icon != "" {
		app["image"] = l.Icon
	}
	if m.Metadata.Author != nil && m.Metadata.Author.Name != "" {
		app["author"] = map[string]any{"@type": "Person", "name": m.Metadata.Author.Name}
	}

	page := pageMeta{
		Title:       l.DisplayName + " — " + siteName,
		Description: desc,
		Canonical:   url,
		OGType:      "website",
		JSONLD:      []any{app, s.breadcrumb(l)},
	}
	return s.render(c, page)
}

// breadcrumb gives the result a Home › Category › Name trail.
func (s *Server) breadcrumb(l catalog.Listing) map[string]any {
	items := []any{
		map[string]any{"@type": "ListItem", "position": 1, "name": siteName, "item": s.base + "/"},
	}
	if l.Category != "" {
		items = append(items, map[string]any{
			"@type": "ListItem", "position": 2, "name": l.Category,
			"item": s.base + "/?category=" + urlQueryEscape(l.Category),
		})
	}
	items = append(items, map[string]any{
		"@type": "ListItem", "position": len(items) + 1, "name": l.DisplayName,
		"item": s.templateURL(l.Name),
	})
	return map[string]any{
		"@context":        "https://schema.org",
		"@type":           "BreadcrumbList",
		"itemListElement": items,
	}
}

// pageMeta is everything the head block is rendered from.
type pageMeta struct {
	Title       string
	Description string
	Canonical   string
	OGType      string
	NoIndex     bool
	JSONLD      any
}

func (s *Server) render(c *okapi.Context, p pageMeta) error {
	return s.renderStatus(c, http.StatusOK, p)
}

func (s *Server) renderStatus(c *okapi.Context, status int, p pageMeta) error {
	if !s.shell {
		return c.String(http.StatusNotFound, "not found")
	}

	c.SetHeader("Cache-Control", "no-cache")
	return c.Data(status, "text/html; charset=utf-8", []byte(s.prefix+s.head(p)+s.suffix))
}

func (s *Server) head(p pageMeta) string {
	img := s.base + ogImagePath
	var b strings.Builder
	tag := func(format string, args ...any) {
		b.WriteString("    ")
		fmt.Fprintf(&b, format, args...)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	tag("<title>%s</title>", html.EscapeString(p.Title))
	tag(`<meta name="description" content="%s" />`, html.EscapeString(p.Description))
	if p.NoIndex {
		tag(`<meta name="robots" content="noindex,follow" />`)
	} else {
		tag(`<meta name="robots" content="index,follow" />`)
	}
	tag(`<link rel="canonical" href="%s" />`, html.EscapeString(p.Canonical))

	tag(`<meta property="og:type" content="%s" />`, html.EscapeString(p.OGType))
	tag(`<meta property="og:site_name" content="%s" />`, siteName)
	tag(`<meta property="og:locale" content="en_US" />`)
	tag(`<meta property="og:url" content="%s" />`, html.EscapeString(p.Canonical))
	tag(`<meta property="og:title" content="%s" />`, html.EscapeString(p.Title))
	tag(`<meta property="og:description" content="%s" />`, html.EscapeString(p.Description))
	tag(`<meta property="og:image" content="%s" />`, html.EscapeString(img))

	tag(`<meta name="twitter:card" content="summary_large_image" />`)
	tag(`<meta name="twitter:title" content="%s" />`, html.EscapeString(p.Title))
	tag(`<meta name="twitter:description" content="%s" />`, html.EscapeString(p.Description))
	tag(`<meta name="twitter:image" content="%s" />`, html.EscapeString(img))

	if p.JSONLD != nil {
		if data, err := json.Marshal(p.JSONLD); err == nil {
			tag(`<script type="application/ld+json">%s</script>`, data)
		}
	}
	b.WriteString("    ")
	return b.String()
}

func (s *Server) templateURL(name string) string { return s.base + "/templates/" + name }

// clamp trims a description to n characters on a word boundary.
func clamp(s string, n int) string {
	if len(s) <= n {
		return s
	}
	cut := s[:n]
	if i := strings.LastIndex(cut, " "); i > n/2 {
		cut = cut[:i]
	}
	return strings.TrimRight(cut, " ,.;:—-") + "…"
}

// urlQueryEscape percent-encodes a category for the storefront's filter link.
func urlQueryEscape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z', ch >= '0' && ch <= '9', ch == '-', ch == '_', ch == '.', ch == '~':
			b.WriteByte(ch)
		case ch == ' ':
			b.WriteByte('+')
		default:
			fmt.Fprintf(&b, "%%%02X", ch)
		}
	}
	return b.String()
}
