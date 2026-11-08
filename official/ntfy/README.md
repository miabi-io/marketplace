# ntfy

A simple HTTP-based pub-sub notification service. Send a message with a plain
`curl -d "Backup done" https://ntfy.example.com/backups` and receive it on your
phone (Android/iOS apps), desktop or browser by subscribing to the topic.

## What gets created

- One application (`ntfy`, image `binwiederhier/ntfy:v2.37`) on port 80
- Two volumes: `data` → `/var/lib/ntfy` (the user and access database) and
  `cache` → `/var/cache/ntfy` (cached messages and attachments)
- A route for the Public URL, created when its domain belongs to the workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address clients use. The web app and topic links are built from it. |
| **Default access** | What anonymous clients may do: `deny-all` (default, private server), `read-only` or `read-write`. |

## After install

With the default `deny-all`, nobody can publish or subscribe until you create a
user. Open a shell in the `ntfy` container and run:

```sh
ntfy user add --role=admin alice
```

Then sign in with that user in the web app or the mobile app (set the server to
your Public URL). Use `ntfy access` to grant other users read or write on
specific topics.

## Notes

- ntfy is told it runs behind a proxy, so rate limiting uses the client's real
  address from the Miabi gateway.
- iOS instant notifications go through the upstream `ntfy.sh` relay; see the
  ntfy docs on `upstream-base-url` if you need them.
