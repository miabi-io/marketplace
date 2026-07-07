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

// Package storefront renders the public, server-side marketplace website over
// the same catalog the API serves: a paginated, searchable home grid and a
// per-template detail page. No build step — html/template, SEO-friendly.
package storefront

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/manifest"
)

var funcs = template.FuncMap{
	"isURL":     func(s string) bool { return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") },
	"provision": provision,
	"badge":     badgeClass,
}

// Each named page is the shared base layout plus a content block; every set
// defines the "page" template the renderer executes.
var pages = map[string]*template.Template{
	"home":     template.Must(template.New("home").Funcs(funcs).Parse(baseHTML + homeHTML)),
	"detail":   template.Must(template.New("detail").Funcs(funcs).Parse(baseHTML + detailHTML)),
	"notfound": template.Must(template.New("notfound").Funcs(funcs).Parse(baseHTML + notFoundHTML)),
}

// renderer adapts the storefront's html/template pages to Okapi's Renderer, so
// handlers render via c.Render rather than hand-writing HTML. We keep
// html/template (not Okapi's text/template loader) to preserve contextual
// auto-escaping of the catalog's community-contributed names, descriptions and
// icon URLs.
func renderer() okapi.Renderer {
	return okapi.RendererFunc(func(w io.Writer, name string, data interface{}, _ *okapi.Context) error {
		t, ok := pages[name]
		if !ok {
			return fmt.Errorf("storefront: unknown template %q", name)
		}
		return t.ExecuteTemplate(w, "page", data)
	})
}

// Handlers renders the storefront over a loaded catalog.
type Handlers struct{ cat *catalog.Catalog }

// Register wires the storefront routes onto the app and installs the HTML renderer.
func Register(app *okapi.Okapi, cat *catalog.Catalog) {
	app.WithRenderer(renderer())
	h := &Handlers{cat: cat}
	app.Get("/", h.Home)
	app.Get("/templates/{name}", h.Detail)
}

type homeData struct {
	Q, Source, Category string
	Result              catalog.Page
	Categories          []catalog.CategoryFacet
	PrevURL, NextURL    string
}

// Home renders the paginated, searchable card grid.
func (h *Handlers) Home(c *okapi.Context) error {
	q := c.Query("q")
	source := c.Query("source")
	category := c.Query("category")
	page := atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	res := h.cat.Search(catalog.Query{Q: q, Source: source, Category: category, Page: page, PerPage: 24})

	data := homeData{Q: q, Source: source, Category: category, Result: res, Categories: h.cat.Categories()}
	if res.Page > 1 {
		data.PrevURL = pageURL(q, source, category, res.Page-1)
	}
	if res.Page < res.TotalPages {
		data.NextURL = pageURL(q, source, category, res.Page+1)
	}
	return c.Render(http.StatusOK, "home", data)
}

type detailData struct {
	T         *catalog.Template
	M         *manifest.Manifest
	Listing   catalog.Listing
	Provision string
}

// Detail renders one template's page.
func (h *Handlers) Detail(c *okapi.Context) error {
	t, ok := h.cat.Get(c.Param("name"))
	if !ok {
		return c.Render(http.StatusNotFound, "notfound", nil)
	}
	l := t.Listing()
	return c.Render(http.StatusOK, "detail", detailData{T: t, M: t.Latest().Manifest, Listing: l, Provision: provision(l)})
}

func pageURL(q, source, category string, page int) string {
	v := url.Values{}
	if q != "" {
		v.Set("q", q)
	}
	if source != "" {
		v.Set("source", source)
	}
	if category != "" {
		v.Set("category", category)
	}
	v.Set("page", strconv.Itoa(page))
	return "/?" + v.Encode()
}

func provision(l catalog.Listing) string {
	var parts []string
	if l.Applications > 0 {
		parts = append(parts, fmt.Sprintf("%d app%s", l.Applications, plural(l.Applications)))
	}
	if l.Databases > 0 {
		parts = append(parts, fmt.Sprintf("%d database%s", l.Databases, plural(l.Databases)))
	}
	if l.Volumes > 0 {
		parts = append(parts, fmt.Sprintf("%d volume%s", l.Volumes, plural(l.Volumes)))
	}
	if len(parts) == 0 {
		return "no dependencies"
	}
	return strings.Join(parts, " · ")
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func badgeClass(source string) string {
	if source == catalog.SourceCommunity {
		return "badge-community"
	}
	return "badge-official"
}

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
