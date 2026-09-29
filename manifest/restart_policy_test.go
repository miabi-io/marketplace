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

package manifest

import "testing"

func TestRestartPolicy(t *testing.T) {
	tmpl := func(policy string) string {
		return `
apiVersion: miabi.io/v1
kind: Template
metadata: {name: a, displayName: A, version: "1.0.0"}
applications:
  - name: migrate
    image: example/migrate
    restartPolicy: ` + policy + `
`
	}
	m, err := Parse([]byte(tmpl("on-failure")))
	if err != nil {
		t.Fatal(err)
	}
	if got := m.Applications[0].RestartPolicy; got != "on-failure" {
		t.Fatalf("restartPolicy = %q, want on-failure", got)
	}
	if _, err := Parse([]byte(tmpl("sometimes"))); err == nil {
		t.Fatal("an unknown restartPolicy was accepted")
	}
}
