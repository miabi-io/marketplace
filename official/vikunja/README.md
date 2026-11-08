# Vikunja

A self-hosted to-do and project management app: lists, kanban boards, Gantt and
table views, labels, due dates, reminders and shared projects for teams.

## What gets created

- One application (`vikunja`, image `vikunja/vikunja:0.24.6`) on port 3456,
  serving both the API and the web frontend
- One PostgreSQL 17 database for users, projects and tasks
- One volume, `files` → `/app/vikunja/files`, for task attachments and
  avatars
- A route for the Public URL, created when its domain belongs to the workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address you use. The frontend and links are built from it. |
| **JWT secret** | Signs auth tokens. Leave blank to auto-generate; changing it later signs everyone out. |

## After install

Open the Public URL and **register** your account from the login page. Every
user signs up this way. Once your team has joined, you can disable open
registration by setting `VIKUNJA_SERVICE_ENABLEREGISTRATION` to `false` in the
app's environment.

## Data and backups

Back up the database (all tasks and projects) and the `files` volume
(attachments).
