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

// Command marketplace runs the Miabi marketplace registry service, and (as
// subcommands) regenerates the machine index and lints the catalog for CI.
//
//	marketplace          # serve the API + storefront
//	marketplace generate # write export.json + registry/index.json (CI drift check)
//	marketplace lint     # validate every embedded template, exit non-zero on error
//
// The generated export.json is the full bundle (every manifest inline) — the
// same document GET /v1/export serves. Committed and published via a CDN
// (jsDelivr) it lets Miabi sync the catalog from static git, no server required;
// the server stays available for self-hosting the API + storefront.
package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi/okapicli"
	"github.com/miabi-io/marketplace/internal/api"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/internal/storefront"
)

func main() {
	cli := okapicli.New(okapi.New(), "Marketplace")

	cli.Command("server", "Serve the marketplace API + storefront", func(*okapicli.Command) error {
		return serve(cli)
	})
	generateCmd := func(*okapicli.Command) error { return generate() }
	cli.Command("generate", "Write export.json + registry/index.json", generateCmd)
	cli.Command("generate-index", "Alias of generate", generateCmd)
	cli.Command("lint", "Validate every embedded template (CI drift check)", func(*okapicli.Command) error {
		return lint()
	})
	cli.DefaultCommand("server")

	if err := cli.Execute(); err != nil {
		logger.Fatal("marketplace command failed", "error", err)
	}
}

func serve(cli *okapicli.CLI) error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	port := 8088
	if v := os.Getenv("MARKETPLACE_PORT"); v != "" {
		if p, e := strconv.Atoi(v); e == nil {
			port = p
		}
	}
	app := cli.Okapi()
	app.WithPort(port)
	app.WithErrorHandler(api.ErrorHandler())
	app.WithOpenAPIDocs(okapi.OpenAPI{
		Title:       "Miabi Marketplace API",
		Version:     "1.0.0",
		Description: "Official + community Miabi application & database template registry.",
		License:     okapi.License{Name: "Apache-2.0", URL: "http://www.apache.org/licenses/LICENSE-2.0"},
	})
	api.Register(app, cat)
	storefront.Register(app, cat)
	return cli.RunServer(&okapicli.RunOptions{
		ShutdownTimeout: 15 * time.Second,
		OnStarted: func() {
			logger.Info("marketplace ready", "templates", len(cat.Templates()), "port", port, "docs", "/docs")
		},
	})
}

// generate writes the static catalog artifacts: export.json (the full bundle,
// every manifest inline — what Miabi syncs from the CDN) and registry/index.json
// (the lightweight machine index). The time-based generatedAt is omitted from
// both so they are deterministic and CI can diff them for drift.
func generate() error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}

	bundle := cat.Bundle()
	bundle.GeneratedAt = ""
	if err := writeJSON("export.json", bundle); err != nil {
		return err
	}
	logger.Info("wrote export.json", "templates", len(bundle.Templates))

	idx := cat.Index()
	idx.GeneratedAt = ""
	if err := writeJSON("registry/index.json", idx); err != nil {
		return err
	}
	logger.Info("wrote registry/index.json", "templates", len(idx.Templates))
	return nil
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// lint loads the catalog, which validates every manifest, verifies digests, and
// enforces cross-source name uniqueness — a non-zero exit fails CI.
func lint() error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	logger.Info("catalog valid", "templates", len(cat.Templates()))
	return nil
}
