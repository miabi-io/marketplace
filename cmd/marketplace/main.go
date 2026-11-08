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

// Command marketplace serves the registry API and storefront; subcommands generate, lint and version support CI.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi/okapicli"
	"github.com/miabi-io/marketplace/internal/api"
	"github.com/miabi-io/marketplace/internal/buildinfo"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/internal/seo"
	"github.com/miabi-io/marketplace/internal/web"
)

func main() {
	cli := okapicli.New(okapi.New(), "Marketplace")

	cli.Command("server", "Serve the marketplace API + storefront", func(*okapicli.Command) error {
		return serve(cli)
	})
	generateCmd := func(*okapicli.Command) error { return generate() }
	cli.Command("generate", "Write export.json + registry/index.json", generateCmd)
	cli.Command("generate-index", "Alias of generate", generateCmd)
	cli.Command("lint", "Validate every embedded template (CI drift check)", func(cmd *okapicli.Command) error {
		return lint(cmd.GetBool("strict"))
	}).Bool("strict", "", false, "Fail on warnings (missing README, release info, image locks)")
	cli.Command("outdated", "Report templates whose images have newer upstream tags", func(cmd *okapicli.Command) error {
		return outdated(cmd.Args(), cmd.GetBool("json"), cmd.GetString("write"))
	}).Bool("json", "", false, "Print the report as JSON").
		String("write", "", "", "Also write the report to this path (registry/freshness.json is embedded)")
	cli.Command("bump", "Publish a new template version with newer image tags", func(cmd *okapicli.Command) error {
		if len(cmd.Args()) == 0 && !cmd.GetBool("all") {
			return errors.New("name the templates to bump, or pass --all")
		}
		set, err := parseSet(cmd.GetString("set"))
		if err != nil {
			return err
		}
		return bump(cmd.Args(), bumpOptions{major: cmd.GetBool("major"), dryRun: cmd.GetBool("dry-run"), asJSON: cmd.GetBool("json"),
			today: time.Now().UTC().Format(time.DateOnly), setTags: set})
	}).Bool("all", "", false, "Bump every template that has an update").
		Bool("major", "", false, "Allow major upgrades (marks the version breaking)").
		Bool("dry-run", "", false, "Report what would change without writing").
		Bool("json", "", false, "Print the bumped versions as JSON").
		String("set", "", "", "Pin explicit tags instead of looking upstream: app=tag[,app=tag]")
	cli.Command("lock", "Backfill release dates and image digests on each template's latest version", func(cmd *okapicli.Command) error {
		return lock(cmd.Args(), time.Now().UTC().Format(time.DateOnly))
	})
	cli.Command("check-images", "Verify every installable image still resolves in its registry", func(cmd *okapicli.Command) error {
		return checkImages(cmd.Args(), cmd.GetBool("strict"))
	}).Bool("strict", "", false, "Also fail when a tag has moved off its locked digest")
	cli.Command("e2e", "Install (and optionally upgrade) templates on a live Miabi, then clean up", func(cmd *okapicli.Command) error {
		inputs, err := parseSet(cmd.GetString("input"))
		if err != nil {
			return err
		}
		return e2e(cmd.Args(), e2eOptions{
			baseURL:   envOr(cmd.GetString("url"), "MIABI_E2E_URL"),
			token:     envOr(cmd.GetString("token"), "MIABI_E2E_TOKEN"),
			workspace: envOr(cmd.GetString("workspace"), "MIABI_E2E_WORKSPACE"),
			domain:    envOr(cmd.GetString("domain"), "MIABI_E2E_DOMAIN"),
			upgrade:   cmd.GetBool("upgrade"),
			keep:      cmd.GetBool("keep"),
			insecure:  cmd.GetBool("insecure"),
			timeout:   cmd.GetDuration("timeout"),
			inputs:    inputs,
		})
	}).String("url", "", "", "Miabi base URL (MIABI_E2E_URL)").
		String("token", "", "", "Account API token (MIABI_E2E_TOKEN)").
		String("workspace", "", "", "Workspace ID to install into (MIABI_E2E_WORKSPACE)").
		String("domain", "", "", "Domain for generated URL/host inputs, default example.com (MIABI_E2E_DOMAIN)").
		String("input", "", "", "Explicit input answers: key=value[,key=value]").
		Bool("upgrade", "", false, "Install the previous version, then upgrade to the latest").
		Bool("keep", "", false, "Leave the install in place for inspection").
		Bool("insecure", "", false, "Skip TLS verification (test servers with private CAs)").
		Duration("timeout", "", 20*time.Minute, "Per-template timeout")
	cli.Command("version", "Print the build version", func(*okapicli.Command) error {
		fmt.Printf("marketplace %s (commit %s, built %s)\n", buildinfo.Version, buildinfo.CommitID, buildinfo.BuildDate)
		return nil
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
		Version:     buildinfo.Version,
		Description: "Official + community Miabi application & database template registry.",
		Favicon:     "/favicon.svg",
		License:     okapi.License{Name: "Apache-2.0", URL: "http://www.apache.org/licenses/LICENSE-2.0"},
		Contact: okapi.Contact{
			Name:  "Hello Miabi",
			Email: " hello@miabi.io",
		},
		UI: okapi.ScalarUI,
	})
	api.Register(app, cat)
	// Registered before the file server so these paths reach the SEO handlers.
	seo.Register(app, cat, os.Getenv("MARKETPLACE_BASE_URL"))
	// Registered last so every API route keeps precedence; unmatched paths fall back to index.html.
	app.WebFS("/", web.Assets, okapi.WebConfig{
		Root:    "dist",
		Exclude: []string{"/v1", "/docs", "/healthz", "/metrics"},
		MaxAge:  365 * 24 * time.Hour,
	})
	return cli.RunServer(&okapicli.RunOptions{
		ShutdownTimeout: 15 * time.Second,
		OnStarted: func() {
			logger.Info("marketplace ready", "version", buildinfo.Version, "templates", len(cat.Templates()), "port", port, "docs", "/docs")
		},
	})
}

// generatedAt is omitted so the output is deterministic and CI can diff it for drift.
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

func lint(strict bool) error {
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	warnings := lintWarnings(cat)
	for _, w := range warnings {
		logger.Warn(w)
	}
	logger.Info("catalog valid", "templates", len(cat.Templates()), "warnings", len(warnings))
	if strict && len(warnings) > 0 {
		return fmt.Errorf("%d lint warning(s)", len(warnings))
	}
	return nil
}

func envOr(v, env string) string {
	if v != "" {
		return v
	}
	return os.Getenv(env)
}

func parseSet(s string) (map[string]string, error) {
	out := map[string]string{}
	for _, pair := range strings.Split(s, ",") {
		if pair = strings.TrimSpace(pair); pair == "" {
			continue
		}
		app, tag, ok := strings.Cut(pair, "=")
		if !ok || app == "" || tag == "" {
			return nil, fmt.Errorf("--set: %q is not app=tag", pair)
		}
		out[app] = tag
	}
	return out, nil
}
