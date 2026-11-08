# Matomo

Privacy-respecting, self-hosted web analytics: a Google Analytics alternative
where you own the data, with goals, funnels, heatmaps and GDPR tooling.

## What gets created

- One application (`matomo`, image `matomo:5-apache`) on port 80
- A MariaDB 11.4 database (`db`), wired in through environment variables
- One volume, `data` → `/var/www/html`, holding Matomo's code, `config.ini.php`,
  plugins and generated files

## After install

Attach a domain to port **80** and open it. Matomo runs its installation wizard.
The database settings are prefilled from the environment, so confirm them, then
create the superuser and the first website to track. The wizard ends with the
JavaScript tracking code to add to your site.

For sites with real traffic, set up report archiving (a scheduled job running
`core:archive`) and turn off browser-triggered archiving under
**System → General settings**.

## Where data lives

Visits and reports are in the MariaDB database. The `data` volume holds
`config/config.ini.php` and installed plugins. Back up both.

## Notes

The `5-apache` tag follows Matomo 5.x, so a redeploy can pick up a newer 5.x
release. Matomo then asks a superuser to run its database upgrade on the next
visit.
