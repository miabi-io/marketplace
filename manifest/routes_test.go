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
 */package manifest

import (
	"strings"
	"testing"
)

func routesTemplate(routes string) string {
	return `
apiVersion: miabi.io/v1
kind: Template
metadata: {name: a, displayName: A, version: "1.0.0"}
inputs:
  - key: public_url
    label: URL
    type: string
applications:
  - name: web
    image: nginx
    ports:
      - container: 8080
        scheme: http
  - name: worker
    image: nginx
routes:
` + routes
}

func TestRoutesValidation(t *testing.T) {
	cases := map[string]string{
		"":                `  - app: web` + "\n" + `    host: "{{ .inputs.public_url }}"`,
		"unknown app":     `  - app: nope` + "\n" + `    host: "{{ .inputs.public_url }}"`,
		"no port":         `  - app: worker` + "\n" + `    host: "{{ .inputs.public_url }}"`,
		"undeclared port": `  - app: web` + "\n" + `    host: "{{ .inputs.public_url }}"` + "\n" + `    port: 9999`,
		"literal host":    `  - app: web` + "\n" + `    host: status.example.com`,
	}
	for want, routes := range cases {
		_, err := Parse([]byte(routesTemplate(routes + "\n")))
		if want == "" {
			if err != nil {
				t.Errorf("valid route rejected: %v", err)
			}
			continue
		}
		if err == nil {
			t.Errorf("%s: accepted", want)
		}
	}
}

func TestRouteHost(t *testing.T) {
	cases := map[string]string{
		"https://Status.Example.com/":       "status.example.com",
		"https://status.example.com:8443/x": "status.example.com",
		"status.example.com":                "status.example.com",
		"status.example.com.":               "status.example.com",
		"":                                  "",
		"http://localhost:3000":             "",
		"http://10.0.0.5":                   "",
		"https://[::1]:443":                 "",
		"intranet":                          "",
	}
	for in, want := range cases {
		if got := RouteHost(in); got != want {
			t.Errorf("RouteHost(%q) = %q, want %q", in, got, want)
		}
	}
	if strings.Contains(RouteHost("https://a.example.com/path?q=1"), "/") {
		t.Error("path leaked into the host")
	}
}
