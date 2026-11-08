# libSQL

A managed [libSQL](https://github.com/tursodatabase/libsql) server (`sqld`), the
open-source fork of SQLite that adds a network protocol, so applications can
use SQLite over HTTP and WebSockets.

## What gets created

- One dedicated libSQL database instance (`libsql`), with no applications

The database is managed by Miabi like any other: it appears under **Databases**,
with its connection details, backups and lifecycle actions.

## After install

Open the database in Miabi to copy its connection URL and auth token. Connect
from an app in the same workspace using a libSQL client (`@libsql/client`, the
Go or Rust driver, and so on) and its network alias.

## Notes

- The template follows libSQL's `latest` image, so it isn't pinned to a release.
- Use a dedicated instance per application: libSQL has no logical databases to
  share between tenants.
