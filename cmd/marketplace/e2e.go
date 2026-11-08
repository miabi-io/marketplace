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

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jkaninda/logger"
	"github.com/miabi-io/marketplace/internal/catalog"
	"github.com/miabi-io/marketplace/manifest"
)

type e2eOptions struct {
	baseURL   string
	token     string
	workspace string
	domain    string
	upgrade   bool
	keep      bool
	insecure  bool
	timeout   time.Duration
	inputs    map[string]string
}

type miabiClient struct {
	base, token, ws string
	http            *http.Client
}

type apiEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   any             `json:"error"`
	Message string          `json:"message"`
}

func (c *miabiClient) call(ctx context.Context, method, path string, body, out any) error {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+"/api/v1/workspaces/"+c.ws+path, rd)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: HTTP %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	if out == nil {
		return nil
	}
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	return json.Unmarshal(env.Data, out)
}

type installJob struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	Error    string   `json:"error"`
	Warnings []string `json:"warnings"`
	Result   *struct {
		InstallID uint `json:"install_id"`
		Apps      []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"apps"`
	} `json:"result"`
}

type appStatus struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// e2e imports each template under an e2e- name, since a workspace's custom template would otherwise shadow the official one.
func e2e(names []string, opt e2eOptions) error {
	if opt.baseURL == "" || opt.token == "" || opt.workspace == "" {
		return errors.New("e2e needs --url, --token and --workspace (or MIABI_E2E_URL, MIABI_E2E_TOKEN, MIABI_E2E_WORKSPACE)")
	}
	if opt.domain == "" {
		opt.domain = "example.com"
	}
	cat, err := catalog.Load()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return errors.New("name the templates to test")
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: opt.insecure} //nolint:gosec // opt-in for private test CAs
	c := &miabiClient{base: strings.TrimRight(opt.baseURL, "/"), token: opt.token, ws: opt.workspace, http: &http.Client{Timeout: 60 * time.Second, Transport: tr}}
	var failed []string
	for _, n := range names {
		t, ok := cat.Get(n)
		if !ok {
			return fmt.Errorf("unknown template %q", n)
		}
		if t.Meta.Deprecated {
			logger.Info("skipping deprecated template", "template", n)
			continue
		}
		if err := e2eOne(c, t, opt); err != nil {
			logger.Error("e2e failed", "template", n, "error", err)
			failed = append(failed, n)
			continue
		}
		logger.Info("e2e passed", "template", n)
	}
	if len(failed) > 0 {
		return fmt.Errorf("e2e failed for: %s", strings.Join(failed, ", "))
	}
	return nil
}

func e2eOne(c *miabiClient, t *catalog.Template, opt e2eOptions) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), opt.timeout)
	defer cancel()
	name := "e2e-" + t.Name
	target := t.Installable()
	from := target
	if opt.upgrade {
		for _, v := range t.Versions {
			if manifest.CompareVersions(v.Version, target.Version) < 0 && !v.Info().Deprecated {
				from = v
				break
			}
		}
		if from.Version == target.Version {
			logger.Info("no older version to upgrade from; testing a fresh install", "template", t.Name)
		}
	}

	versions := []catalog.Version{from}
	if from.Version != target.Version {
		versions = append(versions, target)
	}
	for _, v := range versions {
		raw, rerr := renamed(v.Raw, t.Name, name)
		if rerr != nil {
			return rerr
		}
		if ierr := c.call(ctx, http.MethodPost, "/marketplace/templates/import", map[string]string{"yaml": raw}, nil); ierr != nil {
			return fmt.Errorf("import %s: %w", v.Version, ierr)
		}
	}
	defer func() {
		if opt.keep {
			return
		}
		cleanup, ccancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer ccancel()
		if derr := c.call(cleanup, http.MethodDelete, "/marketplace/templates/"+url.PathEscape(name), nil, nil); derr != nil && err == nil {
			err = fmt.Errorf("delete imported template: %w", derr)
		}
	}()

	inputs := fillInputs(from.Manifest, t.Name, opt.domain, opt.inputs)
	var job installJob
	body := map[string]any{"name": name, "version": from.Version, "inputs": inputs, "display_name": name}
	if err := c.call(ctx, http.MethodPost, "/marketplace/install/jobs", body, &job); err != nil {
		return fmt.Errorf("start install: %w", err)
	}
	logger.Info("installing", "template", t.Name, "version", from.Version, "job", job.ID)
	if err := waitJob(ctx, c, &job); err != nil {
		return err
	}
	installID := job.Result.InstallID
	if !opt.keep {
		defer func() {
			cleanup, ccancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer ccancel()
			if uerr := c.call(cleanup, http.MethodDelete, fmt.Sprintf("/marketplace/installs/%d", installID), nil, nil); uerr != nil && err == nil {
				err = fmt.Errorf("uninstall: %w", uerr)
			}
		}()
	}
	appIDs := longRunning(from.Manifest, job.Result.Apps)
	if err := waitRunning(ctx, c, appIDs); err != nil {
		return fmt.Errorf("%s: %w", from.Version, err)
	}
	for _, w := range job.Warnings {
		logger.Warn("install warning", "template", t.Name, "warning", w)
	}

	if from.Version == target.Version {
		return nil
	}
	newInputs := fillInputs(target.Manifest, t.Name, opt.domain, opt.inputs)
	var res struct {
		Warnings []string `json:"warnings"`
	}
	logger.Info("upgrading", "template", t.Name, "from", from.Version, "to", target.Version)
	if err := c.call(ctx, http.MethodPost, fmt.Sprintf("/marketplace/installs/%d/upgrade", installID),
		map[string]any{"version": target.Version, "inputs": newInputs}, &res); err != nil {
		return fmt.Errorf("upgrade: %w", err)
	}
	for _, w := range res.Warnings {
		logger.Warn("upgrade warning", "template", t.Name, "warning", w)
	}
	// Give the redeploy a moment to flip the status off "running" before polling it.
	time.Sleep(5 * time.Second)
	if err := waitRunning(ctx, c, appIDs); err != nil {
		return fmt.Errorf("after upgrade to %s: %w", target.Version, err)
	}
	return nil
}

// Result apps are created in manifest order; one-shot apps (restartPolicy no/on-failure) exit by design.
func longRunning(m *manifest.Manifest, apps []struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}) []uint {
	var ids []uint
	for i, a := range apps {
		if i < len(m.Applications) {
			if p := m.Applications[i].RestartPolicy; p == "no" || p == "on-failure" {
				continue
			}
		}
		ids = append(ids, a.ID)
	}
	return ids
}

var nameField = regexp.MustCompile(`(?m)^(  name:\s*["']?)([a-z0-9-]+)(["']?\s*)$`)

func renamed(raw []byte, from, to string) (string, error) {
	s := string(raw)
	loc := nameField.FindStringSubmatchIndex(s)
	if loc == nil || s[loc[4]:loc[5]] != from {
		return "", fmt.Errorf("metadata.name %q not found", from)
	}
	return s[:loc[4]] + to + s[loc[5]:], nil
}

// fillInputs answers only what an install can't default: required inputs with no default that aren't generated.
func fillInputs(m *manifest.Manifest, name, domain string, overrides map[string]string) map[string]string {
	out := map[string]string{}
	for _, in := range m.Inputs {
		if v, ok := overrides[in.Key]; ok {
			out[in.Key] = v
			continue
		}
		if !in.Required || in.Default != "" || in.Generate {
			continue
		}
		out[in.Key] = guessInput(in, name, domain)
	}
	return out
}

func guessInput(in manifest.Input, name, domain string) string {
	key := strings.ToLower(in.Key)
	host := "e2e-" + name + "." + domain
	switch {
	case in.Type == manifest.InputSelect && len(in.Options) > 0:
		return in.Options[0]
	case in.Type == manifest.InputBool:
		return "false"
	case in.Type == manifest.InputNumber:
		return "1"
	case in.Type == manifest.InputPassword:
		return randomSecret()
	case strings.Contains(key, "email") || strings.Contains(key, "mail"):
		return "e2e@" + domain
	case strings.Contains(key, "url"):
		return "https://" + host
	case strings.Contains(key, "host") || strings.Contains(key, "domain"):
		return host
	}
	if in.Pattern != "" && strings.Contains(in.Pattern, "https?://") {
		return "https://" + host
	}
	return "e2e"
}

func randomSecret() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func waitJob(ctx context.Context, c *miabiClient, job *installJob) error {
	for {
		switch job.Status {
		case "succeeded":
			if job.Result == nil {
				return errors.New("install job succeeded without a result")
			}
			return nil
		case "failed":
			return fmt.Errorf("install job failed: %s", job.Error)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("install job %s: %w", job.ID, ctx.Err())
		case <-time.After(3 * time.Second):
		}
		if err := c.call(ctx, http.MethodGet, "/marketplace/install/jobs/"+url.PathEscape(job.ID), nil, job); err != nil {
			return err
		}
	}
}

// A deploy with an http/command healthcheck only reports running once healthy, so running is the pass signal.
func waitRunning(ctx context.Context, c *miabiClient, ids []uint) error {
	for {
		pending := 0
		for _, id := range ids {
			var a appStatus
			if err := c.call(ctx, http.MethodGet, fmt.Sprintf("/apps/%d", id), nil, &a); err != nil {
				return err
			}
			switch a.Status {
			case "running":
			case "failed", "stopped":
				return fmt.Errorf("app %s is %s", a.Name, a.Status)
			default:
				pending++
			}
		}
		if pending == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%d app(s) not running: %w", pending, ctx.Err())
		case <-time.After(5 * time.Second):
		}
	}
}
