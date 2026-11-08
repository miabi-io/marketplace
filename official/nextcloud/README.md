# Nextcloud

A self-hosted content collaboration platform: file sync and sharing, calendar,
contacts and a large app store, from the desktop and mobile clients you already
know.

## What gets created

- One application (`nextcloud`, image `nextcloud:34-apache`) on port 80
- One PostgreSQL 17 database
- One volume, `data` → `/var/www/html`, which holds the Nextcloud installation,
  its `config.php`, installed apps and every user's files

## Inputs

| Input | Notes |
|---|---|
| **Admin username** | The first administrator, created on first launch (default `admin`). |
| **Admin password** | Leave blank to auto-generate. It is shown in the install summary and stored as a secret env var. |

## After install

The first start installs Nextcloud into the database, which takes a minute or
two. Attach your domain in Miabi, then sign in with the admin account.

Nextcloud only answers on hostnames it trusts. If the page complains about an
untrusted domain, add yours to `trusted_domains` in `config/config.php` (in the
`data` volume) and set `overwriteprotocol` to `https`, since TLS terminates at
the Miabi gateway. Then work through **Administration settings → Overview**,
which lists any remaining setup warnings.

## Data and backups

Back up both the `data` volume (files, config, apps) and the database. They
must be restored together: the database indexes the files in the volume.

## Upgrades

Nextcloud can only upgrade **one major version at a time** (34 → 35, never 34 →
36). Upgrade through each intermediate template version in turn, and let each
start finish its migration before moving on.
