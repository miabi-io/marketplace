# Forgejo

A lightweight, self-hosted software forge: Git repositories, pull requests,
issues, wikis, packages and Actions CI. Forgejo is a community-governed soft
fork of Gitea.

## What gets created

- One application (`forgejo`, image `codeberg.org/forgejo/forgejo:15`) on port 3000
- A PostgreSQL 17 database (`db`), wired in through environment variables
- One volume, `data` → `/data`, holding repositories, LFS objects, attachments
  and `app.ini`

The template has no install inputs: the database connection is configured for you.

## After install

Attach a domain to port **3000** and open it. Because the database is already
configured, Forgejo skips its web installer. **The first account you register
becomes the site administrator**, so register yourself straight after install.
Afterwards, consider disabling open registration under **Site administration**
or in `app.ini`.

Set the instance's public URL (`ROOT_URL`) to the domain you attached, either
in `app.ini` on the `data` volume or as `FORGEJO__server__ROOT_URL` in the
app's environment. Otherwise clone URLs and emails point at the wrong address.

## Notes and limitations

- **SSH clone is not exposed.** A template cannot publish a host port, so clone
  over HTTPS, or publish the container's SSH port yourself after install.
- **Actions runners are separate.** Run a Forgejo Runner elsewhere and register
  it against this instance.
- **The tag `15` follows the 15.x line**, so a redeploy can pick up a newer
  15.x patch release.
