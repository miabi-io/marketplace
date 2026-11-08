# authentik

A self-hosted identity provider: single sign-on with OAuth2/OpenID Connect, SAML
and LDAP, plus flows for enrolment, recovery and multi-factor authentication.

## What gets created

A stack of two applications sharing one secret key and one database:

- `server` (image `ghcr.io/goauthentik/server:2026.8.5`, command `server`) on
  port 9000: the web UI and the protocol endpoints
- `worker` (same image, command `worker`): background tasks, outposts and
  certificate handling. It serves no HTTP.
- A PostgreSQL 17 database (`db`). authentik dropped Redis in 2025.10, so none is
  provisioned.
- One volume, `data` → `/data`, mounted in both apps, for media, uploaded icons
  and certificates

## Inputs

| Input | Notes |
|---|---|
| **Secret key** | Signs sessions and tokens. Leave blank to auto-generate, and keep it: changing it invalidates every session. |
| **Initial admin email** | Email of the `akadmin` superuser created on first launch. |
| **Initial admin password** | Password for `akadmin`. Leave blank to auto-generate; the install summary shows it. |

## After install

Attach a domain to the `server` app's port **9000** and sign in as `akadmin` with
the password from the install summary. Then create your own admin user and add
applications and providers under **Admin interface → Applications**.

## Upgrades

authentik does not support skipping major releases, and since 2025.12 its
storage lives under `/data` rather than `/media`. When upgrading an older
install, step through the releases one at a time and move the media by hand.
Don't jump straight to the newest version.
