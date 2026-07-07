# Guestbook

A tiny, self-hosted **guestbook** — visitors leave a name and a message, and the
wall shows every signature. It's the reference example for deploying an app to
[Miabi](https://github.com/miabi-io/miabi): an [Okapi](https://github.com/jkaninda/okapi)
(Go) API with an embedded UI, backed by a managed **PostgreSQL** database.

## What gets created

- One **PostgreSQL** database (`db`), provisioned and attached automatically.
- One application (`miabi/guestbook`) listening on port **8080** (HTTP).

## Features

- **Live updates over SSE** — new/removed signatures stream to every open tab.
- **Online-user count** — a real-time "N online" indicator.
- **Live server-time card** — shows the serving version and host, so a
  **canary rollout** is obvious at a glance.

## Notes

- `DATABASE_URL` is wired from the managed database and stored encrypted — you
  never set it by hand.
- Set the **Display name** input to customize the header; leave it for the
  default *Miabi Guestbook*.
- Health is checked at `GET /healthz`, which stays `503` until the database is
  reachable.
