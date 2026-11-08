# Meilisearch

A fast, typo-tolerant search engine with a REST API: index JSON documents and get
relevant, instant results with filtering, facets and sorting.

## What gets created

- One application (`meilisearch`, image `getmeili/meilisearch:v1.56`) on port
  7700, running in production mode
- One volume, `data` → `/meili_data`, holding the indexes

## Inputs

| Input | Notes |
|---|---|
| **Master key** | Root key that protects every endpoint. Leave blank to auto-generate; the install summary shows it. |

## After install

Apps in the same workspace reach Meilisearch at its network alias on port
**7700**. Use the master key only to create scoped API keys (`/keys`), and give
each client the narrowest one it needs: a search-only key for browsers, and an
admin key for the indexer.

Attach a domain only if browsers need to query it directly. In production mode
every request needs a key, and the built-in search preview is turned off.

## Upgrades

The template sets `MEILI_EXPERIMENTAL_DUMPLESS_UPGRADE=true`, so a newer
Meilisearch opens a database written by an older one instead of refusing to
start. Back up the `data` volume (or create a dump) before upgrading anyway: the
upgrade migrates indexes in place.
