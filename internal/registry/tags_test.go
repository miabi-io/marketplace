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

package registry

import "testing"

func TestCompare(t *testing.T) {
	cases := []struct {
		name, cur       string
		tags            []string
		inMajor, latest string
		comparable      bool
	}{
		{"patch and major", "11.4.0", []string{"11.4.0", "11.6.16", "12.0.0", "13.3.0", "13.3.0-ubuntu", "11.7.0-beta1"}, "11.6.16", "13.3.0", true},
		{"suffix kept", "6.45-alpine", []string{"6.45-alpine", "6.75-alpine", "6.75", "7.0-alpine"}, "6.75-alpine", "7.0-alpine", true},
		{"zero major pins minor", "0.24.6", []string{"0.24.6", "0.24.9", "0.25.0", "2.8.0"}, "0.24.9", "2.8.0", true},
		{"v prefix", "v1.52", []string{"v1.52", "v1.56", "v1.56.1", "1.60"}, "v1.56", "v1.56", true},
		{"nothing newer", "3.1", []string{"3.0", "3.1"}, "3.1", "3.1", true},
		{"commit sha", "4511a22", []string{"4511a22", "92109a0"}, "", "", false},
		{"latest", "latest", []string{"latest"}, "", "", false},
		{"multi-number suffix", "7.0-php8.3-apache", []string{"7.1-php8.3-apache", "7.1-php8.5-apache", "8.0-php8.3-apache"}, "7.1-php8.5-apache", "8.0-php8.3-apache", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Compare(c.cur, c.tags)
			if got.Comparable != c.comparable || got.LatestInMajor != c.inMajor || got.Latest != c.latest {
				t.Fatalf("Compare(%q) = %+v, want inMajor=%q latest=%q comparable=%v", c.cur, got, c.inMajor, c.latest, c.comparable)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	cases := map[string][2]string{
		"nginx":                      {"registry-1.docker.io", "library/nginx"},
		"gitea/gitea":                {"registry-1.docker.io", "gitea/gitea"},
		"ghcr.io/goauthentik/server": {"ghcr.io", "goauthentik/server"},
		"docker.io/library/redis":    {"registry-1.docker.io", "library/redis"},
		"quay.io/keycloak/keycloak":  {"quay.io", "keycloak/keycloak"},
		"localhost:5000/team/app":    {"localhost:5000", "team/app"},
	}
	for in, want := range cases {
		if h, r := Split(in); h != want[0] || r != want[1] {
			t.Errorf("Split(%q) = %s %s, want %s %s", in, h, r, want[0], want[1])
		}
	}
}
