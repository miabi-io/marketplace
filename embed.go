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

// Package marketplace embeds the catalog content so the service is stateless: git is the database.
package marketplace

import "embed"

// Content holds <source>/<slug>/<version>/template.yaml plus optional metadata.yaml and README.md per slug.
//
//go:embed all:official all:community
var Content embed.FS

// TemplateSchema is embedded so the service can serve it at the URL its own $id advertises.
//
//go:embed schema/template.schema.json
var TemplateSchema []byte

// Freshness is the generated upstream-tag report; it may be empty before the first `outdated --write`.
//
//go:embed registry/freshness.json
var Freshness []byte
