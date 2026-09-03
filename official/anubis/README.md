# Anubis

Anubis sits in front of an application and makes every visitor's browser solve a
small proof-of-work puzzle before the request reaches the backend. Scrapers that
do not run JavaScript never get through; a real browser pays a barely noticeable
cost once.

## Before you install

Deploy the app you want to protect first. Anubis needs its address on the
workspace network — `http://mb-app-<handle>:<port>` — as the **Upstream URL**.

Attach your domain to **Anubis**, not to the app behind it. If the app keeps its
own route, visitors can reach it directly and bypass the challenge entirely.

## The bot policy

The rules live in a Miabi config (**Configs → anubis-policy**), mounted at
`/data/cfg/botPolicy.yaml`. Editing it redeploys Anubis with the new rules.

The shipped policy allows `/.well-known`, `/favicon.ico` and `/robots.txt`,
denies Cloudflare Workers, and challenges anything claiming to be a browser.

**A request matching no rule is allowed through.** That is Anubis's default, and
it is why feed readers and uptime monitors keep working — but it also means a
scraper with an unusual User-Agent passes untouched. Add explicit rules as you
find them.

Actions are `ALLOW`, `DENY`, `CHALLENGE` and `WEIGH`. See the
[policy reference](https://anubis.techaro.lol/docs/admin/policies).

## Difficulty

The default of 4 is right for almost everyone. Every visitor pays the cost of
whatever you set, so raise it only after you have seen scrapers get through, and
expect slower first loads on low-end phones when you do.
