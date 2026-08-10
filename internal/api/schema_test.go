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
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	marketplace "github.com/miabi-io/marketplace"
)

// The schema is only useful to an editor if the URL it advertises is the URL it
// is served from, so the route and the $id have to stay in agreement.
func TestSchemaIDMatchesItsRoute(t *testing.T) {
	const route = "/schema/template.schema.json"

	var doc struct {
		ID     string `json:"$id"`
		Schema string `json:"$schema"`
		Title  string `json:"title"`
	}
	if err := json.Unmarshal(marketplace.TemplateSchema, &doc); err != nil {
		t.Fatalf("embedded schema is not valid JSON: %v", err)
	}
	if doc.Schema == "" || doc.Title == "" {
		t.Errorf("schema is missing $schema/title: %+v", doc)
	}
	u, err := url.Parse(doc.ID)
	if err != nil {
		t.Fatalf("$id %q is not a URL: %v", doc.ID, err)
	}
	if u.Path != route {
		t.Errorf("$id path = %q, but the service serves %q — an editor pointed at $id would 404", u.Path, route)
	}
}

// The ETag is a strong validator over the served bytes, so a schema change
// invalidates the copy an editor cached.
func TestSchemaETag(t *testing.T) {
	if !strings.HasPrefix(schemaETag, `"`) || !strings.HasSuffix(schemaETag, `"`) {
		t.Errorf("etag %q should be quoted", schemaETag)
	}
	if len(schemaETag) != 66 {
		t.Errorf("etag %q is not a quoted sha256", schemaETag)
	}
}
