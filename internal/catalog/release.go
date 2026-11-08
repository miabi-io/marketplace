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
	"bytes"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

// ReleaseFile is the sidecar next to a version's template.yaml.
const ReleaseFile = "release.yaml"

var digestRe = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)

// Release lives beside template.yaml because Miabi rejects unknown manifest fields, so a new one would break older installs' sync.
type Release struct {
	Released string `yaml:"released,omitempty" json:"released,omitempty"`
	Notes    string `yaml:"notes,omitempty" json:"notes,omitempty"`
	Breaking bool   `yaml:"breaking,omitempty" json:"breaking,omitempty"`
	// RequiresFrom is the oldest version that may upgrade straight here (e.g. Nextcloud can't skip a major).
	RequiresFrom    string `yaml:"requiresFrom,omitempty" json:"requires_from,omitempty"`
	Deprecated      bool   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	DeprecationNote string `yaml:"deprecationNote,omitempty" json:"deprecation_note,omitempty"`
	// Images maps a template app name to the digest its tag resolved to when the version was published.
	Images map[string]string `yaml:"images,omitempty" json:"images,omitempty"`
}

func parseRelease(data []byte) (*Release, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	var r Release
	if err := dec.Decode(&r); err != nil {
		return nil, fmt.Errorf("%s: %w", ReleaseFile, err)
	}
	if r.Released != "" {
		if _, err := time.Parse(time.DateOnly, r.Released); err != nil {
			return nil, fmt.Errorf("%s: released %q is not a YYYY-MM-DD date", ReleaseFile, r.Released)
		}
	}
	for app, d := range r.Images {
		if !digestRe.MatchString(d) {
			return nil, fmt.Errorf("%s: images.%s: %q is not a sha256 digest", ReleaseFile, app, d)
		}
	}
	return &r, nil
}

func MarshalRelease(r *Release) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
