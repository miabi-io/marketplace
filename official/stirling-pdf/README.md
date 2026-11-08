# Stirling PDF

A self-hosted PDF toolkit in the browser: merge, split, rotate, compress,
convert to and from other formats, OCR, add signatures and watermarks, and
more. Files are processed on your server, not a third-party service.

## What gets created

- One application (`stirling-pdf`, image `stirlingtools/stirling-pdf:2.14.3`)
  on port 8080
- One volume, `configs` → `/configs`, holding settings and the user database

## Inputs

| Input | Notes |
|---|---|
| **Require login** | `true` (default) protects the instance with accounts. Use `false` only on a private network. |
| **Admin username / password** | The account created on first start when login is required. Leave the password blank to auto-generate. |

## After install

The first start takes a minute or two; the healthcheck allows 90 seconds.
Attach a domain and sign in with the admin account. Stirling PDF asks you to
change the password on first login. Add users for others from the settings
instead of sharing the admin account.

## Notes

Large OCR and conversion jobs are CPU- and memory-heavy. Give the app room if
your team processes big documents.
