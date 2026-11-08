# Jellyfin

A free, self-hosted media server: stream your movies, shows, music and photos to
browsers, phones, TVs and set-top boxes, with no account or subscription.

## What gets created

- One application (`jellyfin`, image `jellyfin/jellyfin:10.11.11`) on port 8096
- Two volumes: `config` → `/config` (library database, users, metadata and
  settings) and `cache` → `/cache` (transcoding and image cache)

## After install

Attach a domain to port **8096** and open it. The setup wizard asks for the
display language and creates the administrator account, then lets you add
libraries.

## Adding media

The template doesn't mount a media library: templates cannot bind host
directories. Attach a volume holding your media to the `jellyfin` app after
install (for example at `/media`), then point a library at that path in
**Dashboard → Libraries**.

## Notes and limitations

- **No hardware transcoding.** The container has no access to the host's GPU,
  so transcoding runs on the CPU. Prefer clients that can direct-play your
  formats.
- **Local network discovery doesn't work** behind the gateway. Add the server to
  apps by its URL.

## Upgrades

Jellyfin 10.11 migrates the library database on its first start. Back up the
`config` volume before upgrading an existing install, and let the first boot
finish without interrupting it.
