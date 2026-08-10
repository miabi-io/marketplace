# GitLab

A complete DevOps platform: Git hosting, merge requests, issues, CI/CD pipelines
and a package registry, in one self-hosted application.

This template deploys **GitLab CE** from the official omnibus image
(`gitlab/gitlab-ce`), which bundles PostgreSQL, Redis, Puma, Sidekiq and Nginx in
a single container. That is deliberate: GitLab supports an external database only
when a superuser has created the `pg_trgm` and `btree_gist` extensions in it, so
a managed Miabi database would not work out of the box. Everything stateful lives
in the three volumes below.

## Requirements

GitLab is the heaviest application in the catalog. Give the node at least:

- **4 CPU cores** and **8 GB RAM** (4 GB is the documented floor and leaves no
  headroom for CI or a second app on the node)
- **20 GB** of free disk for repositories, uploads, artifacts and the bundled
  database

The template sets no resource caps, so GitLab uses what the node offers. It
already trims the default install — two Puma workers, Sidekiq concurrency 10, and
the built-in Prometheus suite disabled — which saves roughly 2 GB against the
stock configuration.

## What gets created

- One application (`gitlab`, image `gitlab/gitlab-ce:19.2.1-ce.0`) on port 80
- Three volumes: `config` → `/etc/gitlab`, `logs` → `/var/log/gitlab`,
  `data` → `/var/opt/gitlab`

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The domain you attach in Miabi, e.g. `https://gitlab.example.com`. It ends up in clone URLs, notification emails and every redirect — set it correctly before installing, since changing it later means editing the app's environment and redeploying. |
| **Root password** | Initial password for the built-in `root` account. Leave it blank to auto-generate one. |
| **Time zone** | IANA zone for timestamps and scheduled pipelines (default `UTC`). |

## After install

**The first boot takes 5–10 minutes.** Omnibus runs `gitlab-ctl reconfigure` and
the database migrations before Nginx answers, so the app reports unhealthy for a
while — the healthcheck allows 10 minutes before it counts a failure. Watch
progress in the application's logs.

Then attach your domain in Miabi (it must match the Public URL you entered), and
sign in as `root` with the password from the install summary. Change it, and
consider turning off open registration under **Admin → Settings → General →
Sign-up restrictions**.

## Notes and limitations

- **TLS terminates at the Miabi gateway.** The container serves plain HTTP on
  port 80 and is told to forward `X-Forwarded-Proto: https`, so GitLab still
  generates `https://` URLs. Omnibus' own Let's Encrypt integration is disabled —
  it cannot work in a container that owns no domain.
- **SSH clone is not exposed.** A template cannot publish a host port, so port 22
  stays inside the container: clone over HTTPS instead. To enable SSH, publish
  the app's port 22 to a host port after install and set
  `gitlab_rails['gitlab_shell_ssh_port']` in `GITLAB_OMNIBUS_CONFIG` to match.
- **The container registry and Pages are off.** Both need their own hostname and
  certificate; add `registry_external_url` / `pages_external_url` to
  `GITLAB_OMNIBUS_CONFIG` and route the extra hostnames if you want them.
- **CI runners are separate.** Install a GitLab Runner elsewhere and register it
  against this instance — GitLab CE ships no runner of its own.
- **Configuration lives in `GITLAB_OMNIBUS_CONFIG`.** Edit that environment
  variable and redeploy rather than editing `/etc/gitlab/gitlab.rb` inside the
  container: omnibus rewrites the file from the variable on every start.

## Upgrades

GitLab must be upgraded along its supported [upgrade
path](https://docs.gitlab.com/ee/update/#upgrade-paths) — you cannot jump several
major versions in one step. Check the path for your current version before
changing the image tag.
