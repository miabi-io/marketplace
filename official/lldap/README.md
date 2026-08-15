# LLDAP

A light LDAP server for authentication. It gives you a real LDAP directory —
users, groups, bind DNs — behind a small web interface, so the many apps that
only speak LDAP get a user database without anyone having to operate OpenLDAP.

It is deliberately a subset: read-mostly, no schema editing, no replication. If
you need a full directory server, this is the wrong tool. If you need "one place
to manage the accounts my self-hosted apps authenticate against", it is the right
size.

Source: **[github.com/lldap/lldap](https://github.com/lldap/lldap)**.

## What gets created

- One application (`lldap`, image `lldap/lldap:v0.6.3`) serving the web interface
  on 17170 and LDAP on 3890
- One PostgreSQL 17 database — users, groups and password hashes
- One small volume at `/data`, holding only the config file the image writes on
  first start (see *Why a volume* below)

## Before you install

Decide your **Base DN** first. It is the namespace every user and group lives
under, conventionally your domain written as `dc=example,dc=com`. You do not need
to own the domain, but **you cannot change it later** without recreating the
directory — and every application you connect will be configured with it.

## Inputs

| Input | Notes |
|---|---|
| **Base DN** | See above. Decide once. |
| **Public URL** | The URL the web interface is served on; used for password-reset links. |
| **Admin username / email / password** | The first administrator, created once. Leave the password blank to auto-generate. |
| **Key seed** | Derives the key that protects stored passwords. Leave blank to auto-generate. **Back it up** — see below. |
| **JWT secret** | Signs web sessions. Rotating it only signs everyone out. |

## After install: save the key seed

The key seed derives the private key LLDAP protects every stored password with.
**Change it and every existing password stops working** — the hashes were
computed against the key it derives, so users cannot log in and the admin
password reverts to nothing usable. A database backup does not help: the seed is
not in the database.

Read it from the app's environment (**Application → Environment →
`LLDAP_KEY_SEED` → reveal**) and store it somewhere that is *not* this server,
alongside the admin password. Then:

- Attach your domain, matching the Public URL you entered.
- Sign in as the admin and change the password.
- Create the groups your applications will map to roles, before you connect them.

This template generates a seed per install. That matters: the LLDAP image ships a
config file containing a **published** default seed, so an install that does not
override it protects its passwords with a value anyone can read on GitHub. The
generated seed is passed as an environment variable, which takes precedence over
that file.

## Connecting an application

Apps on the same workspace reach LDAP over the internal network, on port
**3890**, using this app's network alias as the host — nothing needs to be
published to the internet for an app in Miabi to authenticate against it.

With the defaults and a base DN of `dc=example,dc=com`:

| Setting | Value |
|---|---|
| Host / port | the `lldap` app's alias, `3890` |
| Bind DN | `cn=admin,ou=people,dc=example,dc=com` |
| Bind password | the admin password |
| User search base | `ou=people,dc=example,dc=com` |
| Group search base | `ou=groups,dc=example,dc=com` |
| Username attribute | `uid` |

The upstream repository keeps worked examples for Nextcloud, Authelia, Grafana,
Jellyfin, Portainer and others under
[`example_configs/`](https://github.com/lldap/lldap/tree/main/example_configs).

**Prefer a dedicated bind account** over the admin: create a regular user, add it
to the `lldap_strict_readonly` group, and bind with that. An app holding your
admin bind credentials can change every account in the directory.

## Notes and limitations

- **Why a volume, when the data is in Postgres.** The image's entrypoint exits
  immediately if `/data` is missing or unwritable — it copies its default config
  there on first start. Nothing durable lives in it: the directory is in Postgres
  and the private key is derived from the seed at boot.
- **LDAPS is off.** LDAP runs unencrypted on 3890, which is fine over the
  workspace's internal network and is *not* fine over the internet. If you must
  expose it, terminate TLS in front of it or enable LDAPS
  (`LLDAP_LDAPS_OPTIONS__ENABLED`) with a certificate mounted into the container.
- **Password reset by email** needs SMTP, which this template does not configure.
  Add `LLDAP_SMTP_OPTIONS__*` to the app's environment to enable it.
- **The web UI cannot change the Base DN**, by design. Reinstalling is the way
  out, which is why it is worth deciding before you install.
- **Version `v0.6.3`** is what is published today. Templates are immutable per
  version, so a newer LLDAP ships as a new template version.
