# n8n

Workflow automation with a visual editor: connect hundreds of services, APIs and
databases, transform data, and run workflows on schedules, webhooks or events.

## What gets created

- One application (`n8n`, image `n8nio/n8n:1.123.92`) on port 5678, set to UTC
- One volume, `data` → `/home/node/.n8n`, holding n8n's SQLite database
  (workflows, credentials, executions) and its encryption key

This is the standalone setup: a single container with SQLite and no queue
workers. It suits personal and small-team use.

## After install

Attach a domain to port **5678** and open it. On first visit, n8n asks you to
create the owner account.

For webhooks to show the right URL, set `WEBHOOK_URL` (and `N8N_HOST` /
`N8N_PROTOCOL`) on the app to the domain you attached, then redeploy. Set
`GENERIC_TIMEZONE` if your schedules should run in a timezone other than UTC.

## Where data lives

Everything is on the `data` volume, including the key n8n uses to encrypt stored
credentials. Back up the whole volume: losing the key makes saved credentials
unreadable.

## Upgrades

n8n 2.x and 3.x are available upstream. They are major releases with breaking
changes, so read n8n's migration notes and back up the `data` volume before
moving an existing install.
