# Linkwarden

A collaborative bookmark manager that archives a full copy of every page you
save (a screenshot, a PDF and a readable text version), so a link outlives the
site behind it. Organise bookmarks into collections, tag them and share them
with your team.

## What gets created

- `linkwarden` (image `ghcr.io/linkwarden/linkwarden:v2.16.3`) on port 3000
- `meilisearch` (image `getmeili/meilisearch:v1.13.3`), the full-text search
  index, reachable only by Linkwarden
- A PostgreSQL 16 database (`db`)
- Two volumes: `data` → `/data/data` (archived pages) and `meili` → `/meili_data`
  (the search index)
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address people visit, with no trailing slash and no path. Sign-in callbacks are built from it, so a wrong value breaks login. |
| **Auth secret** / **Search master key** | Leave blank to auto-generate. Changing the auth secret later signs everyone out. |
| **Close public sign-up** | Off by default so you can register. Turn it on afterwards. |
| **Enable email / Sender address / SMTP URL** | Needed for password reset and collection invitations. The SMTP URL looks like `smtp://user:password@smtp.example.com:587` and is stored as a secret. |

## After install

The first start runs every database migration before serving, so allow a couple
of minutes. Then open the Public URL and sign up: that account is yours. Then
close public sign-up by reinstalling with **Close public sign-up** on, or by
setting `NEXT_PUBLIC_DISABLE_REGISTRATION=true` on the `linkwarden` app.
Otherwise anyone who finds the page can create an account.

## Where data lives

Bookmarks and users are in PostgreSQL. Archived copies are on the `data` volume:
they grow faster than the database and can't be re-created once the original
page is gone, so back them up. The search index on `meili` can be rebuilt.

## Notes

Meilisearch is pinned to the version Linkwarden's upstream setup uses, and is held
there on purpose: Meilisearch refuses an index written by a different version, so
bumping it is a migration rather than a routine update.
