# Invoice Ninja

Invoicing, quotes, recurring billing, online payments, expenses and time tracking
for freelancers and small businesses, with a client portal.

## What gets created

- `app` (image `invoiceninja/invoiceninja-debian:5.13.53`): the PHP-FPM
  application on port 9000, internal only
- `web` (image `nginx:1.30-alpine`) on port 80: the entry point. It serves static
  files and passes PHP requests to `app`.
- A MySQL 8 database (`db`) and a dedicated Redis 8 instance (`redis`) for
  cache, sessions and queues
- Two volumes shared by both apps: `public` → `/var/www/html/public` and
  `storage` → `/var/www/html/storage` (uploads, generated PDFs, logs)
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address clients visit, e.g. `https://invoices.example.com`. It ends up in invoice links, payment redirects and emails. |
| **Admin email / Admin password** | The account created on first start. Leave the password blank to auto-generate. |
| **Application key** | Laravel encryption key. Leave blank to auto-generate, and never change it: stored credentials become unreadable. |

## After install

The first start runs migrations and publishes assets, so `web` waits until the
application files exist. Allow a couple of minutes. Then open the Public URL and
sign in with the admin email and password.

## Notes and limitations

- **Email is logged, not sent.** The template sets `MAIL_MAILER=log`. To email
  invoices, set `MAIL_MAILER=smtp` and the `MAIL_*` SMTP variables on the `app`
  application, then redeploy.
- **TLS terminates at the Miabi gateway.** The app trusts the gateway's forwarded
  headers and does not force HTTPS itself, which would cause a redirect loop.
- The `web` app follows nginx's stable line only.

## Where data lives

Clients, invoices and payments are in the MySQL database. Uploads and generated
PDFs are on the `storage` volume. Back up both.
