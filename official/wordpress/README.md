# WordPress

The world's most popular CMS, for blogs, company sites and stores, with
thousands of themes and plugins.

## What gets created

- One application (`wordpress`, image `wordpress:7.1-php8.3-apache`) on port 80
- One MySQL 8.4 database
- One volume, `content` → `/var/www/html`, holding the WordPress install,
  `wp-config.php`, themes, plugins and uploads

There are no install inputs: the database connection is wired automatically.

## After install

Attach a domain and open it. WordPress runs its **five-minute install**: choose
the site title and create the administrator account. Do that straight away,
since whoever reaches the page first claims the site.

TLS terminates at the Miabi gateway. If WordPress builds `http://` links or
loops on redirects, make sure **Settings → General** uses your `https://`
address.

## Data and backups

Back up both the database (posts, pages, settings) and the `content` volume
(themes, plugins, uploads). They belong together.

## Upgrades

The PHP version is part of the image tag (`php8.3`). A PHP minor change can
break plugins, so it ships as a deliberate template version, separate from
WordPress updates.
