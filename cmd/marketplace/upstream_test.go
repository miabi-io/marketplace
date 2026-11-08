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
	"strings"
	"testing"

	"github.com/miabi-io/marketplace/manifest"
)

const twoApps = `apiVersion: miabi.io/v1
kind: Template
metadata:
  name: demo
  displayName: Demo
  version: 1.2.0 # keep me
databases:
  - name: db
    engine: postgres
    version: "17-alpine"
applications:
  - name: server
    primary: true
    image: jkaninda/posta
    # pinned for the demo
    tag: "0.15.0"
    ports:
      - container: 8080
    env:
      TAG: "not-a-tag-line"
  - name: worker
    image: jkaninda/posta
    tag: 0.15.1
inputs:
  - key: url
    type: string
routes:
  - app: server
    host: "{{ .inputs.url }}"
`

func TestRetagKeepsLayoutAndOnlyTouchesTargets(t *testing.T) {
	m, err := manifest.Parse([]byte(twoApps))
	if err != nil {
		t.Fatal(err)
	}
	out, changes, err := retag(twoApps, m, map[string]string{"worker": "0.19.0"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `    tag: "0.15.0"`) || !strings.Contains(out, `    tag: "0.19.0"`) {
		t.Fatalf("worker not retagged or server touched:\n%s", out)
	}
	if !strings.Contains(out, "# pinned for the demo") || !strings.Contains(out, `version: "17-alpine"`) {
		t.Fatalf("comments or database version lost:\n%s", out)
	}
	if len(changes) != 1 || !strings.Contains(changes[0], "0.15.1 → 0.19.0") {
		t.Fatalf("changes = %v", changes)
	}

	out, err = setManifestVersion(out, "1.2.0", "1.3.0")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "  version: 1.3.0 # keep me") {
		t.Fatalf("version not set:\n%s", out)
	}
	if _, err := manifest.Parse([]byte(out)); err != nil {
		t.Fatalf("result does not parse: %v", err)
	}
}

func TestRetagUnknownAppFails(t *testing.T) {
	m, err := manifest.Parse([]byte(twoApps))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := retag(twoApps, m, map[string]string{"nope": "1"}); err == nil {
		t.Fatal("expected an error for an app with no tag line")
	}
}

func TestNextVersion(t *testing.T) {
	if got := nextVersion("1.4.2", false); got != "1.5.0" {
		t.Fatalf("minor: %s", got)
	}
	if got := nextVersion("1.4.2", true); got != "2.0.0" {
		t.Fatalf("major: %s", got)
	}
}
