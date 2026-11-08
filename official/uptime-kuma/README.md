# Uptime Kuma

A self-hosted monitoring tool: watch HTTP(S) endpoints, TCP ports, DNS records,
containers and more, get alerts through dozens of notification channels, and
publish public status pages.

## What gets created

- One application (`uptime-kuma`, image `louislam/uptime-kuma:2`) on port 3001
- One volume, `data` → `/app/data`, holding monitors, history, notification
  settings and status pages

## After install

Attach a domain and open it. On first visit Uptime Kuma asks you to choose a
database and create the admin account. Keep the built-in embedded option unless
you have a reason not to, and do it right away, since whoever reaches the page
first claims the instance.

Then add monitors and notification channels from the dashboard.

## Notes

- Monitor the services you care about from **somewhere other than the server
  they run on**. If this node goes down, so does the monitor.
- The `2` tag follows Uptime Kuma 2.x releases. Each template version pins the
  image digest it was published with.

## Data and backups

Everything lives in the `data` volume. Back it up.
