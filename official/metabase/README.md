# Metabase

Open-source business intelligence: connect your databases, ask questions with a
point-and-click builder or SQL, and share dashboards and alerts with your team.

## What gets created

- One application (`metabase`, image `metabase/metabase:v0.52.17`) on port 3000
- A PostgreSQL 17 database (`db`) for Metabase's own application data:
  questions, dashboards, users and settings

The template has no install inputs. Metabase's application database is wired in
for you, so nothing is kept in the container's default embedded H2 file.

## After install

The first start runs migrations and can take a minute. Attach a domain to port
**3000** and open it: Metabase's setup wizard creates the admin account, then asks
you to connect the first database you want to analyse. Apps' databases in the
workspace are reachable by their network alias.

Set **Site URL** under **Admin → Settings → General** to the domain you attached,
so links in emails and alerts point to the right place.

## Where data lives

Everything Metabase stores is in the PostgreSQL `db` database; there is no
volume. Your analysed data stays in your own databases.

## Upgrades

Newer Metabase releases (v0.6x) are available upstream. They migrate the
application database on first start and can't be rolled back, so back up `db`
before upgrading an existing install.
