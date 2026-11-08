# pgAdmin 4

The standard web interface for administering PostgreSQL: browse schemas, run
queries, inspect plans, manage roles and take backups from the browser.

## What gets created

- One application (`pgadmin`, image `dpage/pgadmin4:8.14`) on port 80, capped
  at 512 MiB of memory and half a CPU
- One volume, `data` → `/var/lib/pgadmin`, holding pgAdmin's own configuration
  database (registered servers, saved passwords, preferences)

No PostgreSQL is installed: pgAdmin connects to databases you already run.

## Inputs

| Input | Notes |
|---|---|
| **Admin email** | Login email for the pgAdmin web interface. |
| **Admin password** | Leave blank to auto-generate. It is stored as a secret env var. |

## After install

Attach a domain, sign in with the admin email and password, then register a
server with **Add New Server**. For a Miabi-managed PostgreSQL in the same
workspace, use its in-network host and port and the credentials from the
database's **Connection** panel.

The admin credentials only seed the first start. Change the password later from
inside pgAdmin, not by editing the env var.

## Notes

pgAdmin exposes every database you register to anyone who can sign in. Keep the
admin password strong, or restrict access to the route.
