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

package seo

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

var ldRe = regexp.MustCompile(`(?s)<script type="application/ld\+json">(.*?)</script>`)

func TestBaseURLDefaults(t *testing.T) {
	if got := New(nil, "").base; got != DefaultBaseURL {
		t.Errorf("empty base: got %q, want %q", got, DefaultBaseURL)
	}
	if got := New(nil, "https://example.test/").base; got != "https://example.test" {
		t.Errorf("trailing slash should be trimmed, got %q", got)
	}
}

func TestBuiltShellCarriesTheMarkers(t *testing.T) {
	s := New(nil, "")
	if !s.shell {
		t.Skip("SPA not built into the binary (make build-ui)")
	}
	if s.suffix == "" {
		t.Fatal("markers missing from the built index.html: the head block would never be replaced")
	}
	if strings.Contains(s.prefix, "<title>") || strings.Contains(s.suffix, "<title>") {
		t.Error("the default <title> is outside the markers, so a page would render two titles")
	}
	if !strings.Contains(s.suffix, "</head>") {
		t.Error("expected the shell's </head> after the marked block")
	}
}

func TestHeadRendersPageTags(t *testing.T) {
	s := New(nil, "https://example.test")
	head := s.head(pageMeta{
		Title:       "GitLab — Miabi Marketplace",
		Description: `A "complete" DevOps platform`,
		Canonical:   "https://example.test/templates/gitlab",
		OGType:      "website",
		JSONLD:      map[string]any{"@context": "https://schema.org", "@type": "SoftwareApplication", "name": "GitLab"},
	})

	for _, want := range []string{
		"<title>GitLab — Miabi Marketplace</title>",
		`<link rel="canonical" href="https://example.test/templates/gitlab" />`,
		`<meta property="og:url" content="https://example.test/templates/gitlab" />`,
		`<meta property="og:image" content="https://example.test/og-image.png" />`,
		`<meta name="robots" content="index,follow" />`,
	} {
		if !strings.Contains(head, want) {
			t.Errorf("head is missing %s\n%s", want, head)
		}
	}
	// A quote in a description must not break out of the attribute.
	if strings.Contains(head, `content="A "complete"`) {
		t.Error("description was not attribute-escaped")
	}
	m := ldRe.FindStringSubmatch(head)
	if m == nil {
		t.Fatal("no JSON-LD block")
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(m[1]), &doc); err != nil {
		t.Fatalf("JSON-LD is not valid JSON: %v", err)
	}
	if doc["name"] != "GitLab" {
		t.Errorf("JSON-LD name = %v", doc["name"])
	}
}

func TestHeadNoIndex(t *testing.T) {
	head := New(nil, "https://example.test").head(pageMeta{NoIndex: true})
	if !strings.Contains(head, `content="noindex,follow"`) {
		t.Errorf("expected noindex, got %s", head)
	}
}

func TestClamp(t *testing.T) {
	if got := clamp("short", 160); got != "short" {
		t.Errorf("short description should pass through, got %q", got)
	}
	long := strings.Repeat("word ", 60)
	got := clamp(long, descMax)
	if len(got) > descMax+len("…") {
		t.Errorf("clamp exceeded the limit: %d chars", len(got))
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("a clamped description should end in an ellipsis, got %q", got)
	}
}

func TestURLQueryEscape(t *testing.T) {
	if got := urlQueryEscape("Web"); got != "Web" {
		t.Errorf("got %q", got)
	}
	if got := urlQueryEscape("Data & Analytics"); got != "Data+%26+Analytics" {
		t.Errorf("got %q", got)
	}
}
