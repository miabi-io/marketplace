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

// Package registry is a minimal read-only OCI distribution client: list tags and resolve a tag's digest.
package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const dockerHub = "registry-1.docker.io"

var ErrUnauthorized = errors.New("registry refused access")

var ErrNotFound = errors.New("not found")

// Accepting index types first makes the digest the multi-arch one, which every platform can pull.
var manifestAccept = strings.Join([]string{
	"application/vnd.oci.image.index.v1+json",
	"application/vnd.docker.distribution.manifest.list.v2+json",
	"application/vnd.oci.image.manifest.v1+json",
	"application/vnd.docker.distribution.manifest.v2+json",
}, ", ")

type Client struct {
	http *http.Client
	ua   string

	mu     sync.Mutex
	tokens map[string]string
}

// New reads DOCKERHUB_USERNAME/DOCKERHUB_TOKEN when set, to lift Docker Hub's anonymous rate limit.
func New() *Client {
	return &Client{http: &http.Client{Timeout: 60 * time.Second}, ua: "miabi-marketplace", tokens: map[string]string{}}
}

func Split(image string) (host, repo string) {
	first, rest, ok := strings.Cut(image, "/")
	if ok && (strings.ContainsAny(first, ".:") || first == "localhost") {
		host, repo = first, rest
	} else {
		host, repo = dockerHub, image
		if !strings.Contains(image, "/") {
			repo = "library/" + image
		}
	}
	if host == "docker.io" || host == "index.docker.io" {
		host = dockerHub
	}
	return host, repo
}

func (c *Client) Tags(ctx context.Context, image string) ([]string, error) {
	host, repo := Split(image)
	next := fmt.Sprintf("https://%s/v2/%s/tags/list?n=1000", host, repo)
	var out []string
	for next != "" {
		resp, err := c.do(ctx, http.MethodGet, next, host, repo, "")
		if err != nil {
			return nil, err
		}
		var body struct {
			Tags []string `json:"tags"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: decode tags: %w", image, err)
		}
		out = append(out, body.Tags...)
		next = nextLink(resp.Header.Get("Link"), host)
	}
	return out, nil
}

// Digest prefers HEAD, which Docker Hub does not count against its pull rate limit.
func (c *Client) Digest(ctx context.Context, image, tag string) (string, error) {
	host, repo := Split(image)
	u := fmt.Sprintf("https://%s/v2/%s/manifests/%s", host, repo, url.PathEscape(tag))
	resp, err := c.do(ctx, http.MethodHead, u, host, repo, manifestAccept)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	d := resp.Header.Get("Docker-Content-Digest")
	if d == "" {
		// Some registries omit the header on HEAD; a GET always carries it.
		resp, err = c.do(ctx, http.MethodGet, u, host, repo, manifestAccept)
		if err != nil {
			return "", err
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		d = resp.Header.Get("Docker-Content-Digest")
	}
	if d == "" {
		return "", fmt.Errorf("%s:%s: registry returned no digest", image, tag)
	}
	return d, nil
}

func (c *Client) do(ctx context.Context, method, u, host, repo, accept string) (*http.Response, error) {
	key := host + "/" + repo
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", c.ua)
		if accept != "" {
			req.Header.Set("Accept", accept)
		}
		c.mu.Lock()
		tok := c.tokens[key]
		c.mu.Unlock()
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		switch {
		case resp.StatusCode == http.StatusUnauthorized && attempt == 0:
			challenge := resp.Header.Get("WWW-Authenticate")
			_ = resp.Body.Close()
			tok, err := c.token(ctx, host, repo, challenge)
			if err != nil {
				return nil, err
			}
			c.mu.Lock()
			c.tokens[key] = tok
			c.mu.Unlock()
			continue
		case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
			_ = resp.Body.Close()
			return nil, fmt.Errorf("%s/%s: %w", host, repo, ErrUnauthorized)
		case resp.StatusCode == http.StatusNotFound:
			_ = resp.Body.Close()
			return nil, fmt.Errorf("%s/%s: %w", host, repo, ErrNotFound)
		case resp.StatusCode >= 300:
			_ = resp.Body.Close()
			return nil, fmt.Errorf("%s %s: HTTP %d", method, u, resp.StatusCode)
		}
		return resp, nil
	}
	return nil, fmt.Errorf("%s/%s: %w", host, repo, ErrUnauthorized)
}

var challengeParam = regexp.MustCompile(`(\w+)="([^"]*)"`)

func (c *Client) token(ctx context.Context, host, repo, challenge string) (string, error) {
	params := map[string]string{}
	for _, m := range challengeParam.FindAllStringSubmatch(challenge, -1) {
		params[m[1]] = m[2]
	}
	realm := params["realm"]
	if realm == "" {
		return "", fmt.Errorf("%s/%s: %w (no token realm)", host, repo, ErrUnauthorized)
	}
	q := url.Values{"scope": {"repository:" + repo + ":pull"}}
	if s := params["service"]; s != "" {
		q.Set("service", s)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, realm+"?"+q.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", c.ua)
	if host == dockerHub {
		if u, p := os.Getenv("DOCKERHUB_USERNAME"), os.Getenv("DOCKERHUB_TOKEN"); u != "" && p != "" {
			req.SetBasicAuth(u, p)
		}
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s/%s: %w (token HTTP %d)", host, repo, ErrUnauthorized, resp.StatusCode)
	}
	var body struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Token != "" {
		return body.Token, nil
	}
	return body.AccessToken, nil
}

var linkRe = regexp.MustCompile(`<([^>]+)>`)

func nextLink(link, host string) string {
	m := linkRe.FindStringSubmatch(link)
	if m == nil {
		return ""
	}
	if strings.HasPrefix(m[1], "http") {
		return m[1]
	}
	return "https://" + host + m[1]
}
