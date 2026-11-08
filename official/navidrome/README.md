# Navidrome

A lightweight, self-hosted music server. It indexes your library and streams it
to a modern web player and to any Subsonic-compatible client (DSub, Symfonium,
play:Sub, Feishin and many more).

## What gets created

- One application (`navidrome`, image `deluan/navidrome:0.63.2`) on port 4533
- Two volumes: `data` → `/data` (Navidrome's database, cache and settings) and
  `music` → `/music` (your library)

## Inputs

| Input | Notes |
|---|---|
| **Library scan schedule** | How often `/music` is rescanned: a cron expression or an interval such as `@every 24h` (the default). `0` scans only on start. |
| **Log level** | `error`, `warn`, `info` (default) or `debug`. |

## After install

Attach a domain in Miabi and open it: on first visit Navidrome asks you to
create the admin account. Do that straight away, since whoever reaches the page
first claims it.

The `music` volume **starts empty**. Copy your library into it (any folder
layout works; Navidrome reads the tags), then trigger a scan from the web UI or
wait for the next scheduled one.

## Data and backups

The library lives in the `music` volume and everything else (users,
playlists, play counts, ratings) in the `data` volume. Back up `data`; `music`
is usually a copy of files you already have elsewhere.
