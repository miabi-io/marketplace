# Mealie

A self-hosted recipe manager and meal planner. Import recipes from any website by
URL, plan the week, and generate shopping lists. Households can share it.

## What gets created

- One application (`mealie`, image `ghcr.io/mealie-recipes/mealie:v2.8.0`) on port 9000
- A PostgreSQL 17 database (`db`)
- One volume, `data` → `/app/data`, for recipe images, backups and uploads
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address people visit, e.g. `https://recipes.example.com`. Links and invitations are built from it. |

## After install

Open the Public URL and sign in with Mealie's default administrator account
(documented upstream), then change its email and password straight away under
your profile. Invite other users from **Settings → Users**.

## Where data lives

Recipes, plans and users are in PostgreSQL. Recipe images and Mealie's own
backups are on the `data` volume. Back up both.

## Upgrades

Mealie v3 is available upstream. It is a major release, so read Mealie's release
notes and take a backup before moving an existing install to it.
