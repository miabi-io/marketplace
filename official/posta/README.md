# Posta

A self-hosted transactional email platform: send through an API or SMTP,
manage templates, domains and webhooks, receive inbound mail, verify addresses,
and follow every message from a dashboard.

Source: **[goposta.dev](https://goposta.dev)**.

## What gets created

A stack of two applications sharing one configuration:

- `server` (image `jkaninda/posta:0.15.1`) on port 9000: the API and the
  embedded dashboard
- `worker` (same image, `worker` command): delivers queued mail with 10
  concurrent jobs and up to 5 retries, so it scales independently of the API
- One PostgreSQL 17 database and one dedicated Redis (queue and caches)
- One volume, `data` → `/data`, shared by both apps
- A route for the Public URL, created when its domain belongs to the workspace

The server runs no embedded worker, so each delivery is attempted exactly once
per retry.

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address you use. It drives links, OAuth callbacks and CORS. |
| **Admin email / password** | The first account. The password must be at least 12 characters; leave it blank to auto-generate. |
| **JWT secret** | Signs auth tokens and tracking links. Leave blank to auto-generate. |
| **Encryption key** | Encrypts stored SMTP passwords. Leave blank to auto-generate, and keep it: changing it makes stored SMTP passwords unreadable. |

## After install

Sign in with the admin account, add your sending domain, and publish the DNS
records Posta shows for it (SPF, DKIM, DMARC) before sending real mail.

The public form-ingest endpoint (`/api/v1/f/{key}`) is enabled and
rate-limited per IP, form, sender and workspace. The limits are env vars on the
`server` app if you need to tune them.

## Data and backups

Back up the database (accounts, templates, message history) and the `data`
volume. Keep a copy of the encryption key with the backup.
