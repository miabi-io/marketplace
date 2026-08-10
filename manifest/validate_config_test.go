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

package manifest_test

import (
	"strings"
	"testing"

	"github.com/miabi-io/marketplace/manifest"
)

// head is the manifest preamble every case below shares.
const head = `apiVersion: miabi.io/v1
kind: Template
metadata: { name: t, displayName: T, version: 1.0.0 }
`

// anApp satisfies the "at least one application or database" rule, so a case
// exercising configs alone still reaches the config checks.
const anApp = "applications:\n  - name: a\n    image: i\n"

func parse(t *testing.T, body string) error {
	t.Helper()
	_, err := manifest.Parse([]byte(head + body))
	return err
}

func TestConfigTemplateValid(t *testing.T) {
	err := parse(t, `configs:
  - name: provisioning
    mode: "0640"
    sensitive: true
    delimiters: ["<<", ">>"]
    files:
      datasources/ds.yml: "apiVersion: 1"
applications:
  - name: grafana
    image: grafana/grafana
    mounts:
      - config: provisioning
        path: /etc/grafana/provisioning
      - config: provisioning
        key: datasources/ds.yml
        path: /etc/grafana/provisioning/datasources/ds.yml
        mode: "0444"
`)
	if err != nil {
		t.Fatalf("valid config template rejected: %v", err)
	}
}

// The linter must reject here whatever the platform would reject at install —
// a template that lints but fails to install is the failure mode this guards.
func TestConfigTemplateRejects(t *testing.T) {
	tests := []struct {
		name, body, want string
	}{
		{"no files", "configs:\n  - name: c\n    files: {}\n" + anApp, "at least one entry"},
		{"bad name", "configs:\n  - name: Bad_Name\n    files: { a.txt: x }\n" + anApp, "must match"},
		{"duplicate name", "configs:\n  - name: c\n    files: { a.txt: x }\n  - name: c\n    files: { b.txt: y }\n" + anApp, "duplicate config"},
		{"absolute key", "configs:\n  - name: c\n    files: { /etc/a.txt: x }\n" + anApp, "relative path"},
		{"traversal key", "configs:\n  - name: c\n    files: { \"../a.txt\": x }\n" + anApp, "relative path"},
		{"setuid mode", "configs:\n  - name: c\n    mode: \"4755\"\n    files: { a.txt: x }\n" + anApp, "setuid"},
		{"one delimiter", "configs:\n  - name: c\n    delimiters: [\"<<\"]\n    files: { a.txt: x }\n" + anApp, "two distinct"},
		{"same delimiters", "configs:\n  - name: c\n    delimiters: [\"<<\", \"<<\"]\n    files: { a.txt: x }\n" + anApp, "two distinct"},
		{
			"unknown config",
			"applications:\n  - name: a\n    image: i\n    mounts: [{ config: nope, path: /x }]\n",
			"unknown config",
		},
		{
			"both sources",
			"volumes: [{ name: v }]\nconfigs:\n  - name: c\n    files: { a.txt: x }\napplications:\n  - name: a\n    image: i\n    mounts: [{ volume: v, config: c, path: /x }]\n",
			"both",
		},
		{
			"neither source",
			"applications:\n  - name: a\n    image: i\n    mounts: [{ path: /x }]\n",
			"exactly one",
		},
		{
			"key on a volume mount",
			"volumes: [{ name: v }]\napplications:\n  - name: a\n    image: i\n    mounts: [{ volume: v, key: a.txt, path: /x }]\n",
			"only valid with a config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parse(t, tt.body)
			if err == nil {
				t.Fatalf("expected rejection, got none")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %q does not mention %q", err, tt.want)
			}
		})
	}
}
