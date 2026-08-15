# Certio

Self-signed PKI and TLS certificate management in a single binary: create roots
and intermediates (or import the CA you already have), issue certificates with
real SAN support, renew and revoke them, and publish a CRL and OCSP responder
clients can actually reach.

Everything is available three ways over one engine — a web dashboard, a REST API
and a CLI — so a certificate issued from a terminal is identical to one issued
from the browser. Private keys are AES-256-GCM encrypted at rest and every
mutation lands in an append-only audit log.

Source: **[github.com/jkaninda/certio](https://github.com/jkaninda/certio)**.

## What gets created

- One application (`certio`, image `jkaninda/certio:0.0.1`) on port 8080
- One PostgreSQL 17 database — the CAs, the issued certificates and their
  encrypted private keys

No volume: with Postgres, Certio keeps no state on its own filesystem, so the
database is the only thing that needs backing up (and the master key below,
which is not in it).

## Before you install

Have the **public URL** ready. `CERTIO_BASE_URL` is baked into the **CRL
distribution point of every certificate issued from that moment on**, so it must
be the URL your clients can reach. Changing it later fixes new certificates only
— the ones already issued keep pointing at the old address.

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The domain you attach in Miabi. See the warning above. |
| **Admin email / password** | The first administrator, created once on an empty database. Leave the password blank to auto-generate. |
| **Master key** | Encrypts every stored private key. Leave blank to auto-generate. **Back it up** — see below. |
| **JWT secret** | Signs dashboard sessions. Rotating it only signs everyone out. |
| **Organization / Country** | Defaults for the `O` and `C` fields on issued certificates. |

## After install: save the master key

This is the one step you cannot skip.

The master key encrypts every private key Certio stores. **Lose it and every
certificate's private key is unreadable — including your CA's.** A database
backup does not help: the database is what the key protects, and the key is not
in it.

Read it from the app's environment (**Application → Environment →
`CERTIO_MASTER_KEY` → reveal**) and store it somewhere that is *not* this server:
a password manager, a sealed envelope, your existing secret store. Then:

- Attach your domain, matching the Public URL you entered.
- Sign in as the admin and change the password.
- Create your root CA — and consider **name-constraining** it to the domains you
  own. A root in a trust store without constraints can mint a certificate for any
  name on the internet.

## Notes and limitations

- **Do not change the master key** on an existing install. It is not a rotation
  knob; it is the key to everything already stored.
- **Key download policy** is set to `once`: an issued private key can be
  collected once and not again. Set `CERTIO_KEY_DOWNLOAD_POLICY=always` in the
  app's environment if your workflow needs repeat downloads, understanding that a
  key downloadable forever is a copy waiting to be made.
- **ACME** (`http-01`, `dns-01`, wildcards) lets cert-manager, Traefik, Caddy,
  certbot or acme.sh renew internal certificates unattended. It is gated by
  credentials an administrator issues from the dashboard.
- **Prometheus metrics** are exposed at `/metrics` — alert on certificate and CA
  expiry rather than watching a dashboard.
- **The scheduler runs in-process** (hourly by default) for renewals and expiry
  warnings. Nothing extra to deploy.
- **Version `0.0.1`** is what is published today. Templates are immutable per
  version, so a newer Certio ships as a new template version.

## Trusting the CA

A private CA is only useful once clients trust it. Certio generates per-platform
install instructions (Debian, RHEL, macOS, Windows, Java, Node.js, Docker, curl)
for each CA it holds — start there rather than hand-rolling `update-ca-certificates`.
