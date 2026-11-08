# Paperless-ngx

A document management system that turns paper into a searchable archive. Drop in
scans or PDFs and Paperless runs OCR, extracts the text, suggests tags,
correspondents and document types, and files everything for full-text search.

## What gets created

- One application (`paperless`, image `ghcr.io/paperless-ngx/paperless-ngx:3.4.2`)
  on port 8000
- One PostgreSQL 17 database and one dedicated Redis (task queue and live
  updates)
- Four volumes:
  - `data` → `/usr/src/paperless/data`: search index and app data
  - `media` → `/usr/src/paperless/media`: original and archived documents
  - `consume` → `/usr/src/paperless/consume`: drop folder, ingested
    automatically
  - `export` → `/usr/src/paperless/export`: target of the document exporter
- A route for the Public URL, created when its domain belongs to the workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address you use. Paperless checks CSRF against it, so a wrong value breaks sign-in. |
| **Admin username / password** | Superuser created on first start. Leave the password blank to auto-generate. |
| **Secret key** | Django signing key. Leave blank to auto-generate. |
| **OCR language** | Tesseract code(s), e.g. `eng`, `deu` or `eng+fra`. |
| **Time zone** | IANA zone used for dates, e.g. `Europe/Paris`. |

## After install

The first start runs the database migrations and can take a couple of minutes.
Sign in with the admin account, then upload documents from the web UI, by email
(configure a mail account under **Mail**), or by placing files in the `consume`
volume.

## Data and backups

Your documents live in the `media` volume and their metadata in the database:
back up both. For a portable, version-independent backup, run the exporter
(`document_exporter /usr/src/paperless/export`) from a shell in the container.
