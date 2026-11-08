# Docmost

A collaborative wiki and documentation workspace: real-time editing, spaces for
teams, page history and per-page permissions. It is an open-source alternative
to Confluence and Notion.

## What gets created

- One application (`docmost`, image `docmost/docmost:0.95.0`) on port 3000
- A PostgreSQL 17 database (`db`) for pages, users and spaces
- A dedicated Redis 8 instance (`redis`) for real-time collaboration and queues
- One volume, `storage` → `/app/data/storage`, for uploaded files and attachments
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address people visit, e.g. `https://docs.example.com`. Invitation and password-reset links are built from it. |
| **App secret** | Signs sessions and tokens. Leave blank to auto-generate; changing it later signs everyone out. |
| **Sender address / name** | The From line on invitations and password resets. |
| **SMTP host, port, username, password, Implicit TLS** | Leave the host blank to run without mail. Use port 587 with Implicit TLS off (STARTTLS), or 465 with it on. |

## After install

Open the Public URL and sign up: the first account creates the workspace and
becomes its owner. Invite the rest of your team from the workspace settings.
Invitations need a working SMTP sender, so fill in the SMTP inputs if you plan
to add other people.

## Where data lives

Pages, users and permissions live in the PostgreSQL database. Uploaded files live
on the `storage` volume. Back up both. Docmost can also store files in S3 or
Azure Blob: switch `STORAGE_DRIVER` and add the provider's credentials in the
app's environment to move them off the node.
