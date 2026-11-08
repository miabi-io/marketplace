# Wiki.js

A modern wiki: write in Markdown or a visual editor, organise pages in a tree,
keep full page history, and control access with groups and fine-grained
permissions. Sign-in can be local or through external providers.

## What gets created

- One application (`wikijs`, image `ghcr.io/requarks/wiki:2.5.317`) on port
  3000
- One PostgreSQL 17 database. The wiki itself (pages, history, users) lives
  there.
- One volume, `content` → `/wiki/data/content`, used only if you enable a
  local or Git storage target

There are no install inputs.

## After install

The first start initialises the database and can take up to a minute. Attach a
domain and open it: Wiki.js shows a **setup wizard** where you create the
administrator account and set the site URL. Use your attached domain as the
site URL.

## Data and backups

Back up the database. For a second copy in plain files, enable the **Git** or
**Local File System** storage target under **Administration → Storage**.

## Upgrades

This template follows Wiki.js 2.x. Wiki.js 3 is a rewrite with its own
migration path; it will ship as a separate, breaking template version.
