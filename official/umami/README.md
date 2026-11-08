# Umami

Privacy-focused, cookie-less web analytics: page views, referrers, devices,
countries and custom events, on a clean dashboard, without tracking individual
visitors.

## What gets created

- One application (`umami`, image `ghcr.io/umami-software/umami:postgresql-v2.20.2`)
  on port 3000
- One PostgreSQL 17 database, which holds every event Umami records

## Inputs

| Input | Notes |
|---|---|
| **App secret** | Signs login tokens. Leave blank to auto-generate. |

## After install

Attach a domain and sign in with Umami's built-in default account, **`admin` /
`umami`**, then **change that password immediately**: the default is public
knowledge.

Add a website under **Settings → Websites** and paste the tracking snippet Umami
gives you into your site's `<head>`.

## Data and backups

All analytics data is in the database. Back it up on a schedule from the
database's page in Miabi.
