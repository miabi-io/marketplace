# Wallabag

A self-hosted read-it-later app. Save web articles from the browser extension
or the mobile apps, and Wallabag stores a clean, readable copy you can tag,
annotate and read offline, even after the original page disappears.

## What gets created

- One application (`wallabag`, image `wallabag/wallabag:2.6.15`) on port 80
- One PostgreSQL 17 database for your articles, tags and annotations
- One volume, `images` → `/var/www/wallabag/web/assets/images`, for images
  downloaded with saved articles
- A route for the Public URL, created when its domain belongs to the workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address you use. Links and asset URLs are built from it, so a wrong value breaks images and redirects. |

## After install

The image creates a default account, **`wallabag` / `wallabag`**. Sign in and
**change its password immediately** under **Config → User information**, or
create your own user and remove the default one.

Install the browser extension or mobile app and point it at the Public URL.
The apps need an API client, which you create under **API clients management**.

## Data and backups

Back up the database (articles and metadata) and the `images` volume.
