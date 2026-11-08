# PostgreSQL

A standalone, managed PostgreSQL 17 instance, for applications you deploy
yourself or for connecting from outside the workspace.

## What gets created

- One **dedicated** PostgreSQL instance (`postgres:17-alpine`) with a database
  and a user, created and managed by Miabi like any database you add from the
  **Databases** page

There is no application and there are no inputs.

## After install

Open the database in **Databases** to find its connection details: the
in-network host and port, database name, user and password, plus a ready-made
connection URI. Reference them from other apps' environment instead of copying
the password around.

Backups, restores, resource limits and public exposure are managed from the
same page, as for any Miabi database.

## Upgrades

A PostgreSQL **major** upgrade (for example 17 → 18) changes the on-disk format
and needs a dump and restore. It never happens just by upgrading this template.
