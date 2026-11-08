# Roundcube

A browser-based webmail client for any IMAP server: folders, search, an HTML
editor, address books, filters and a plugin system.

Roundcube is a **client**, not a mail server. It reads mail from the IMAP server
and sends through the SMTP server you point it at.

## What gets created

- One application (`roundcube`, image `roundcube/roundcubemail:1.7.5-apache`)
  on port 80
- One PostgreSQL 17 database for user preferences, address books and cached
  data. The mail itself stays on your mail server.

## Inputs

| Input | Notes |
|---|---|
| **IMAP server** | Your mail server's IMAP host. Prefix with `ssl://` for implicit TLS (993) or `tls://` for STARTTLS (143). |
| **IMAP port** | `993` (default) for SSL/TLS, `143` for STARTTLS. |
| **SMTP server** | Your outgoing server. Prefix with `tls://` for STARTTLS or `ssl://` for implicit TLS. |
| **SMTP port** | `587` (default) for STARTTLS, `465` for SSL/TLS. |

## After install

Attach a domain and sign in with a mailbox's own IMAP username and password.
There is no separate Roundcube account, and the first login creates the user's
preferences.

If login fails, check the IMAP host prefix and port pair first; a TLS mismatch
is the usual cause.

## Data and backups

Back up the database for preferences and contacts. Mail is backed up wherever
your mail server keeps it.
