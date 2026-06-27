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

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Parse decodes a template manifest from YAML and validates it. Unknown fields
// are rejected (community manifests are untrusted input — see Validate). The
// returned manifest is normalized (defaults filled).
func Parse(data []byte) (*Manifest, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // reject unknown/misspelled fields rather than ignoring them

	var m Manifest
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	m.normalize()
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

// normalize fills defaults so downstream code never has to special-case empties.
func (m *Manifest) normalize() {
	for i := range m.Databases {
		if m.Databases[i].Placement == "" {
			m.Databases[i].Placement = PlacementAuto
		}
	}
	for i := range m.Applications {
		for j := range m.Applications[i].Ports {
			if m.Applications[i].Ports[j].Scheme == "" {
				m.Applications[i].Ports[j].Scheme = "http"
			}
		}
	}
	// A single application is implicitly primary.
	if len(m.Applications) == 1 {
		m.Applications[0].Primary = true
	}
}

// Digest returns the sha256 of data, prefixed "sha256:", matching the index.yaml
// format.
func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}
