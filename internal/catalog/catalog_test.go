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
	"strings"
	"testing"
	"testing/fstest"

	"github.com/miabi-io/marketplace/manifest"
)

func TestLoadEmbedded(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(c.Templates()) < 12 {
		t.Fatalf("expected the seeded official catalog, got %d templates", len(c.Templates()))
	}

	// The community sample is loaded and labeled.
	w, ok := c.Get("okapi-example")
	if !ok || w.Source != SourceCommunity {
		t.Fatalf("okapi-example should be a community template, got %+v ok=%v", w, ok)
	}

	// An official template resolves with a sha256 digest.
	ng, ok := c.Get("nginx")
	if !ok || ng.Source != SourceOfficial {
		t.Fatalf("nginx should be official, got %+v ok=%v", ng, ok)
	}
	if !strings.HasPrefix(ng.Latest().Digest, "sha256:") {
		t.Fatalf("digest format: %q", ng.Latest().Digest)
	}

	if c.Bundle().ETag != c.ETag() || c.ETag() == "" {
		t.Fatal("bundle ETag should match the catalog ETag and be non-empty")
	}
}

// TestBundleDigestRoundTrip is the exact integrity check the Miabi consumer
// runs: every inlined manifest's digest must equal sha256 of its bytes.
func TestBundleDigestRoundTrip(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	for _, bt := range c.Bundle().Templates {
		for _, v := range bt.Versions {
			if got := manifest.Digest([]byte(v.Manifest)); got != v.Digest {
				t.Fatalf("%s@%s: digest %s != %s", bt.Slug, v.Version, got, v.Digest)
			}
		}
	}
}

func TestSearchFilterAndPaginate(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if comm := c.Search(Query{Source: SourceCommunity}); comm.Total != 1 {
		t.Fatalf("expected 1 community template, got %d", comm.Total)
	}
	p := c.Search(Query{PerPage: 5, Page: 1})
	if p.PerPage != 5 || len(p.Items) > 5 {
		t.Fatalf("pagination: per_page=%d items=%d", p.PerPage, len(p.Items))
	}
	if p.TotalPages < 2 {
		t.Fatalf("expected multiple pages for %d templates at 5/page", p.Total)
	}
	if hit := c.Search(Query{Q: "nginx"}); hit.Total != 1 || hit.Items[0].Slug != "nginx" {
		t.Fatalf("search q=nginx: %+v", hit)
	}
}

const fixtureManifest = `apiVersion: miabi.io/v1
kind: Template
metadata:
  name: dup
  displayName: Dup
  version: 1.0.0
  category: Web
applications:
  - name: app
    image: nginx
    tag: latest
`

// TestDuplicateSlugAcrossSourcesFails ensures the loader fails closed when the
// same slug appears in both official/ and community/.
func TestDuplicateSlugAcrossSourcesFails(t *testing.T) {
	fsys := fstest.MapFS{
		"official/dup/1.0.0/template.yaml":  {Data: []byte(fixtureManifest)},
		"community/dup/1.0.0/template.yaml": {Data: []byte(fixtureManifest)},
	}
	if _, err := loadFS(fsys); err == nil || !strings.Contains(err.Error(), "duplicate slug") {
		t.Fatalf("expected a duplicate-slug error, got %v", err)
	}
}
