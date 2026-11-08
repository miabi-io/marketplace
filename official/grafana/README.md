# Grafana

An open-source observability platform: dashboards over metrics, logs and traces,
with alerting and a large catalog of data source plugins.

## What gets created

- One application (`grafana`, image `grafana/grafana:11.6.16`) on port 3000
- One volume, `data` → `/var/lib/grafana`, holding Grafana's SQLite database
  (dashboards, users, data sources), plugins and alert state

## Inputs

| Input | Notes |
|---|---|
| **Admin password** | Password for the built-in `admin` account. Leave blank to auto-generate; the install summary shows it. |

## After install

Attach a domain to port **3000** and sign in as `admin` with the password from
the install summary. Add a data source (Prometheus, Loki, PostgreSQL and so on)
under **Connections**, then build or import dashboards.

Set `GF_SERVER_ROOT_URL` to the domain you attached in the app's environment.
Without it, links in alert notifications and shared dashboards point at the
wrong address.

## Where data lives

Everything Grafana stores is on the `data` volume. Back it up to keep dashboards
and users.

## Upgrades

Grafana 12 and 13 are available upstream. Major versions remove deprecated
features and plugins (for example AngularJS-based panels), so check the
upstream "What's new" and breaking-changes notes before moving to them.
