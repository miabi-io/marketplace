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

// Package web embeds the built storefront (the Vue SPA) so one binary serves the
// catalog API and the site. `dist/` is a build artifact staged by `make build-ui`
// — only a .gitkeep placeholder is committed, so `go build` works on a clean
// checkout (the API still serves; the site 404s until the UI is built).
package web

import "embed"

// Assets holds the built SPA under a top-level "dist/" directory, served via
// Okapi's WebFS with WebConfig{Root: "dist"}. The `all:` prefix keeps dotted
// files (the .gitkeep placeholder, Vite's occasional dotfiles) in the embed.
//
//go:embed all:dist
var Assets embed.FS
