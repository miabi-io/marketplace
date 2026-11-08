# Valkey

A high-performance key-value store and cache: the Linux Foundation's
open-source fork of Redis, wire-compatible with Redis clients and commands.

## What gets created

- One application (`valkey`, image `valkey/valkey:9-alpine`), reachable on port
  6379 over the workspace network
- One volume, `data` → `/data`. Append-only persistence is on
  (`--appendonly yes`), so data survives restarts.

## Inputs

| Input | Notes |
|---|---|
| **Password** | Clients must authenticate with it (`requirepass`). Leave blank to auto-generate. Letters, digits and `- _ . : @` only, no spaces, at least 8 characters. |

## After install

Point any Redis client at the app's in-network alias on port 6379 with the
password, for example `redis://:<password>@<valkey alias>:6379/0`.

The password is stored in two secret env vars, `VALKEY_PASSWORD` and
`VALKEY_EXTRA_FLAGS`. If you change it, update both and redeploy.

## Valkey or Redis?

The catalog also has a **Redis** template, which provisions a Miabi-managed
Redis database with backups from the Databases page. Choose Valkey when you
want the open-source fork, or run it as a regular application.
