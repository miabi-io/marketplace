# Mattermost

Self-hosted team chat: channels, threads, direct messages, file sharing, calls,
playbooks and a large integration ecosystem. This template runs the Team Edition.

## What gets created

- One application (`mattermost`, image `mattermost/mattermost-team-edition:11.11.4`)
  on port 8065
- A PostgreSQL 17 database (`db`)
- Five volumes: `config`, `data` (uploaded files), `logs`, `plugins` and
  `client-plugins`, mounted under `/mattermost`
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address people visit, e.g. `https://chat.example.com`. Desktop and mobile clients, invitation links and integrations all use it. |

## After install

Open the Public URL. **The first account created becomes the system admin.** Then
create a team and invite people. Configure an SMTP server under
**System Console → Environment → SMTP** so invitations and notifications can be
emailed.

## Notes and limitations

- **No container healthcheck.** The image ships without a shell, so Miabi can't
  run a probe inside it; watch the logs on first start.
- Full-text search uses the built-in Bleve index (`/mattermost/bleve-indexes`).

## Where data lives

Messages, users and channels are in PostgreSQL. Uploaded files are on the `data`
volume and server settings on `config`. Back up the database and both volumes.
