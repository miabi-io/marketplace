# Goma Gateway

A lightweight, declarative API gateway and reverse proxy: route by host, path,
WebSocket, gRPC or raw TCP/UDP, and put middleware — rate limiting, basic/JWT/
LDAP/OIDC auth, CORS, access policies, caching — in front of any backend.

This template runs a **second, application-level gateway inside your workspace**.
It does not replace the Goma Gateway that Miabi itself runs at the edge: that one
attaches domains and issues certificates for your apps. Use this one when you
want your own routing layer — an API surface stitched together from several
internal services, per-route auth, canary weights between two versions of a
backend — behind a single Miabi domain.

## What gets created

- One application (`goma-gateway`, image `jkaninda/goma-gateway:0.14.0`) on ports
  8080 (HTTP) and 8443 (HTTPS)
- Two volumes: `data` → `/etc/goma`, `acme` → `/etc/letsencrypt`
- One config (`goma-gateway-config`) holding `goma.yml`, mounted at
  `/etc/goma/goma.yml`

## Inputs

| Input | Notes |
|---|---|
| **Log level** | `debug`, `info`, `warn`, `error` or `off` (default `info`). |
| **Prometheus metrics** | Exposes `/metrics` on the gateway (default on). It is not authenticated by default — add a `basicAuth` middleware to `monitoring.middleware.metrics` before routing that path anywhere public. |

## After install

The gateway starts with **no routes**: `/healthz` answers 200 and everything else
returns 404 until you add one. There are two places to do that, and the
difference is whether the change needs a restart.

**Static config — `goma.yml`.** Open **Configs** in the workspace sidebar, edit
`goma.yml`, and save: the config is versioned and every app mounting it is
redeployed, which is the restart this file needs anyway. Entry points, timeouts,
logging, TLS, providers and the certificate manager live here.

```yaml
gateway:
  routes:
    - name: api
      path: /
      target: http://my-api:8080     # a sibling app's network alias
```

**Dynamic config — `/etc/goma/extra`.** Anything in that directory (on the `data`
volume) is watched and hot-reloaded, so routes and middlewares added there take
effect **without a redeploy**. It is the right home for route definitions you
change often; `goma.yml` is the right home for how the gateway itself runs.

Then attach a Miabi domain to port **8080**. TLS terminates at the platform
gateway, so the container only needs to speak HTTP — port 8443 is Goma's own TLS
listener, for traffic that reaches this app directly.

## Notes and limitations

- **Backends are addressed by app name.** Other applications in the workspace are
  reachable at their network alias (`http://my-api:8080`), so routes point at
  service names, not IPs.
- **ACME is off, but its storage is ready.** No certificate provider is enabled by
  default. If you turn one on under `certManager`, leave its `storageFile` under
  `/etc/letsencrypt` — that is the `acme` volume, so the account key and issued
  certificates survive a redeploy instead of re-issuing into Let's Encrypt's rate
  limits. Certificates you bring yourself are a different thing and belong in
  `/etc/goma/certs`, on the `data` volume.
- **Behind the platform gateway, enable proxy trust.** Set `gateway.proxy.enabled:
  true` with the platform network in `trustedProxies` if any middleware depends on
  the real client IP (rate limiting by IP, access policies); otherwise every
  request looks like it came from the gateway.
- **Rate limiting and caching are per-instance** unless you point `gateway.redis`
  at a Redis instance — install one from the marketplace and reference it there if
  you scale the gateway beyond a single replica.
- **`/healthz` drives the app healthcheck.** Leaving `monitoring.enableLiveness`
  on, and `monitoring.host` empty, is what keeps the container reported healthy.

## Documentation

- [Goma Gateway documentation](https://goma.jkaninda.dev/)
- [Source](https://github.com/jkaninda/goma-gateway)
