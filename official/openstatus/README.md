# OpenStatus

Status pages for publishing incidents and scheduled maintenance, managed from a
dashboard. This template installs OpenStatus's status-page setup: the dashboard,
the public status page and a libSQL database.

It does **not** include automated monitoring. OpenStatus's monitoring needs
Tinybird for its data and probes deployed elsewhere, set up by hand after
install. Use it alongside the monitoring you already have, and post incidents
here yourself.

## Before you install

**Dashboard URL** is the address you will open the dashboard on. Sign-in links
are built from it, so attach that domain to the **dashboard** app after install.

Put the status page somewhere that stays up when your product doesn't — another
provider, or at least another server. A status page that goes down with the
outage it should report is worse than none.

## Signing in

There is no email setup. Enter your email on the sign-in page, then open the
**dashboard** app's logs in Miabi: the sign-in link is printed there as
`>>> Magic Link: …`. Open it to finish signing in. The first sign-in creates
your account and workspace.

To sign in without the logs, fill in the GitHub or OIDC fields at install
(authentik and Keycloak both work as the OIDC provider).

## Set the workspace limits

A self-hosted workspace starts on OpenStatus's free-plan limits: one status
page and three components. Raise them once, after your first sign-in, from the
**dashboard** app's terminal:

```sh
q() { curl -s -X POST "$DATABASE_URL" -H "Authorization: Bearer $DATABASE_AUTH_TOKEN" -H 'Content-Type: application/json' -d @-; echo; }

# Your workspace id (normally 1)
echo '{"statements":["SELECT id, slug FROM workspace"]}' | q

q <<'JSON'
{"statements":["UPDATE workspace SET limits = '{\"status-pages\":20,\"page-components\":500,\"maintenance\":true,\"status-subscribers\":true,\"custom-domain\":true,\"password-protection\":true,\"email-domain-protection\":true,\"white-label\":true,\"no-index\":true,\"custom-theme\":true,\"i18n\":true,\"members\":\"Unlimited\",\"audit-log\":true,\"notifications\":true,\"notification-channels\":50}' WHERE id = 1"]}
JSON
```

## Publish a status page

1. In the dashboard, create a status page. Its **slug** is `acme` below.
2. Attach a domain to the **status-page** app, for example `status.example.com`.
3. Point the page at that domain. The dashboard's domain setting needs a Vercel
   account, so set it from the **dashboard** app's terminal instead (with `q`
   defined as above):

   ```sh
   q <<'JSON'
   {"statements":["UPDATE page SET custom_domain = 'status.example.com' WHERE slug = 'acme'"]}
   JSON
   ```

Each page is served on its own domain: repeat steps 2 and 3 for another page.

## Upgrades and data

The **migrate** app applies database migrations and stops; a stopped migrate app
is expected. Redeploying it after an upgrade applies the new migrations.

Everything lives in the libSQL database: back it up from **Databases**.
