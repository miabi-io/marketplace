# Vaultwarden

A lightweight server compatible with the Bitwarden clients: browser
extensions, desktop and mobile apps, and the CLI. Store passwords, passkeys,
cards and notes, share them through organizations, and sync everywhere.

Vaultwarden is an independent implementation, not the official Bitwarden
server.

## What gets created

- One application (`vaultwarden`, image `vaultwarden/server:1.37.5-alpine`) on
  port 80
- One volume, `data` → `/data`, holding the SQLite database, attachments, sends
  and the server's keys
- A route for the Public URL, created when its domain belongs to the workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | Must be the exact HTTPS address clients use. WebAuthn/passkeys and email links depend on it. |
| **Admin token** | Guards the `/admin` panel. Leave blank to auto-generate. It is stored as a secret env var. |
| **Allow signups** | `true` (default) lets anyone who reaches the server register. |

## After install

1. Open the Public URL and create your account(s).
2. **Turn signups off**: set `SIGNUPS_ALLOWED` to `false` in the app's
   environment and redeploy, or change it in `/admin`. Invite further users
   from an organization instead.
3. In the Bitwarden apps, choose **Self-hosted** and enter the Public URL as
   the server.

The `/admin` panel (SMTP, signups, users, diagnostics) asks for the admin token.

## Data and backups

Everything, including the vault database and the keys that protect it, lives in
the `data` volume. **Back it up regularly and test a restore.** Losing it means
losing every vault on the server.
