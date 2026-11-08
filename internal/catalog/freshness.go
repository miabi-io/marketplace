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
	"encoding/json"
	"fmt"
)

// FreshnessFile is written by `marketplace outdated --write` and embedded at build time.
const FreshnessFile = "registry/freshness.json"

// Freshness statuses, worst last.
const (
	FreshCurrent = "current" // every image is on its newest upstream tag
	FreshUpdate  = "update"  // a newer tag exists within the pinned major
	FreshMajor   = "major"   // upstream has a newer major than the one pinned
	FreshHeld    = "held"    // frozen on purpose by metadata.yaml updates.<app>.hold
)

type FreshnessReport struct {
	CheckedAt string                       `json:"checked_at"`
	Templates map[string]TemplateFreshness `json:"templates"`
}

// TemplateFreshness is one template's entry; Status is the worst of its apps.
type TemplateFreshness struct {
	Version   string         `json:"version"`
	Status    string         `json:"status"`
	CheckedAt string         `json:"checked_at,omitempty"`
	Apps      []AppFreshness `json:"apps"`
}

// LatestInMajor equals Tag when there is nothing to bump.
type AppFreshness struct {
	App           string `json:"app"`
	Image         string `json:"image"`
	Tag           string `json:"tag"`
	LatestInMajor string `json:"latest_in_major,omitempty"`
	Latest        string `json:"latest,omitempty"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
	Hold          string `json:"hold,omitempty"`
}

func ParseFreshness(data []byte) (*FreshnessReport, error) {
	r := &FreshnessReport{Templates: map[string]TemplateFreshness{}}
	if len(data) == 0 {
		return r, nil
	}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, fmt.Errorf("%s: %w", FreshnessFile, err)
	}
	if r.Templates == nil {
		r.Templates = map[string]TemplateFreshness{}
	}
	return r, nil
}

// A report entry for an older version is stale, so it is dropped rather than shown against the new one.
func (c *Catalog) attachFreshness(r *FreshnessReport) {
	for i := range c.templates {
		t := &c.templates[i]
		f, ok := r.Templates[t.Name]
		if !ok || f.Version != t.Latest().Version {
			continue
		}
		f.CheckedAt = r.CheckedAt
		t.Freshness = &f
	}
}
