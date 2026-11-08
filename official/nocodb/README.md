# NocoDB

An open-source Airtable alternative: build tables with grid, gallery, kanban and
form views, share them with your team, and use the REST API NocoDB generates
for every table.

## What gets created

- One application (`nocodb`, image `nocodb/nocodb:2026.10.1`) on port 8080
- One PostgreSQL 17 database, which holds NocoDB's metadata and the tables you
  create in it
- One volume, `data` → `/usr/app/data`, for uploaded attachments
- A route for the Public URL, created when its domain belongs to the workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address people use, e.g. `https://nocodb.example.com`. Invitation and share links are built from it. |
| **JWT secret** | Signs auth tokens. Leave blank to auto-generate; changing it later signs everyone out. |

## After install

Open the Public URL. The first account you sign up with becomes the
**super admin**, so create it right away. Invite the rest of your team from
there.

## Data and backups

Back up the database (everything you build) and the `data` volume
(attachments).
