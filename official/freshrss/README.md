# FreshRSS

A fast, self-hosted RSS and Atom feed aggregator with multiple users, keyboard
shortcuts, and a Google Reader / Fever API for mobile apps.

## What gets created

- One application (`freshrss`, image `freshrss/freshrss:1.30.4`) on port 80
- One volume, `data` → `/var/www/FreshRSS/data`, holding the SQLite database,
  user settings and caches

No database server is needed: FreshRSS uses SQLite on the volume. Feeds are
refreshed every 20 minutes by the container's built-in cron.

## Inputs

| Input | Notes |
|---|---|
| **Timezone** | IANA zone for feed timestamps and the refresh schedule (default `UTC`). |

## After install

Attach a domain to port **80** and open it. FreshRSS walks you through a short
setup wizard. Choose **SQLite** as the database and create the admin account at
the end. Then add feeds, or import an OPML file from your previous reader.

To use a mobile app, enable API access under **Settings → Authentication**, then
set an API password in your profile.

## Where data lives

Everything (feeds, articles, users and configuration) is on the `data` volume.
Back it up to keep your subscriptions and read state.
