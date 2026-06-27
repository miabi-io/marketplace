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

// Package marketplace embeds the catalog content (official/ + community/) so the
// service ships it and is stateless — git is the database. The same directories
// are the canonical, forkable source contributors edit; merge to main rebuilds
// the image and the live API matches the repo.
package marketplace

import "embed"

// Content holds every template's files: <source>/<slug>/<version>/template.yaml
// plus optional <source>/<slug>/metadata.yaml and README.md.
//
//go:embed all:official all:community
var Content embed.FS
