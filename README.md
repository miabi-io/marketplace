# Miabi Marketplace

The standalone registry + storefront that serves **official and community**
[Miabi](https://github.com/miabi-io/miabi) application & database templates. It
is the *producer* side of the marketplace; Miabi is the *consumer* (it syncs the
catalog into a local cache and serves it in the console).

- **Stateless**: git is the database. Templates under `official/` + `community/`
  are embedded at build time; merging to `main` regenerates the catalog so it
  always matches the repo.
- **Two ways to consume it** — Miabi works with either:
  - **Hosted (recommended):** point `MIABI_MARKETPLACE_URL` at
    `https://marketplace.miabi.io` — the live API + storefront, updated on every
    merge to `main`, ETag/304-capable. Nothing to host, and you get the
    storefront and `/docs` alongside the bundle.
  - **Static or self-hosted:** CI also commits `export.json` (the full bundle),
    so you can sync straight from git — e.g.
    `https://cdn.jsdelivr.net/gh/miabi-io/marketplace@main/export.json` — or run
    the Okapi service below on your own domain and point
    `MIABI_MARKETPLACE_URL` at its base URL. Useful for air-gapped installs and
    private forks.

## Repository layout

```
official/<name>/<version>/template.yaml   curated, CODEOWNERS-protected (source of truth for official)
official/<name>/metadata.yaml             optional storefront metadata (featured, screenshots, …)
official/<name>/README.md                 optional long description (detail page)
community/<name>/...                      contributed, same shape; open PRs
export.json                               GENERATED full bundle (every manifest inline) — what /v1/export serves, also syncable from git
registry/index.json                       GENERATED lightweight machine index (CI checks for drift)
manifest/                                 the miabi.io/v1 manifest module: parse + validate + digest
schema/template.schema.json               JSON Schema for editors + CI
internal/{catalog,api,web}                the Okapi service (web embeds the built storefront)
web/                                      the storefront SPA (Vue 3 + Pinia + Vite)
cmd/marketplace                           server + generate-index + lint
```

`<name>` is the template **handle** — lowercase `^[a-z0-9][a-z0-9-]*$`, matching
the manifest's `metadata.name` and unique across `official/` + `community/`. The
manifest also carries a `metadata.displayName` (the free-text label shown in the
storefront and console) and a `metadata.version`.

## API

The service mirrors the Miabi/Posta stack (Go + Okapi, `{success,data,error}`
envelope on client routes). Machine routes (`/v1/export`, `/v1/index`) return
their raw document so a consumer can decode it directly; `/v1/export` is Miabi's
primary sync call.

The full reference — every route, parameters, and schemas — is served live and
interactive:

| Path | Purpose |
|------|---------|
| `/docs` | **Interactive API documentation** (browse and try every endpoint). |
| `/openapi.json` | The raw OpenAPI spec behind `/docs`. |
| `/healthz` · `/metrics` | Health probe + Prometheus metrics. |

## Storefront

A Vue 3 + Pinia SPA (`web/`) served by the same binary over Okapi's `WebFS`, at
`/` and `/templates/{name}`. Search is instant — results follow typing, with no
submit button — and every filter is mirrored into the URL, so any result set is
a shareable link.

`make build-ui` builds it with Vite and stages the output into
`internal/web/dist`, where `go build` embeds it; `make dev-ui` runs the Vite dev
server on :3100 with `/v1` proxied to a local `make run`. Only a `.gitkeep` is
committed under `internal/web/dist`, so a clean checkout still builds — the API
serves, and the storefront 404s until the UI is built.

<p align="center">
  <img src="docs/images/storefront.png" alt="Miabi Marketplace storefront — searchable template grid" width="100%">
</p>

## Run

```sh
go run ./cmd/marketplace            # serve on :8088 (MARKETPLACE_PORT to override)
go run ./cmd/marketplace lint       # validate every embedded template
go run ./cmd/marketplace generate   # rewrite export.json + registry/index.json (CI runs + diffs this)
```

## How Miabi consumes it

Miabi does a conditional GET (ETag) of the bundle — this service's `/v1/export`
or a static `export.json` — caches it in Redis, and merges it with its embedded
official floor and per-workspace custom imports, serving Official / Community /
Custom tabs locally (search + pagination are done client-side, so no server
round-trips are needed).

Set `MIABI_MARKETPLACE_URL` to enable the sync; empty means embedded-only
(air-gapped). The recommended value is the hosted registry:

```sh
MIABI_MARKETPLACE_URL=https://marketplace.miabi.io
```

A server base URL and a static `export.json` URL are both accepted, so a fork or
an internal mirror drops in the same way.

<p align="center">
  <img src="docs/images/miabi-marketplace.png" alt="Miabi console Marketplace — Official / Community / Custom template tabs" width="100%">
</p>

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Community templates are a fast
PR-to-live path; official templates are maintainer-gated via `CODEOWNERS`. You
can also test a template live before contributing by importing it into the Miabi
console (**Marketplace → Import**) and installing it.