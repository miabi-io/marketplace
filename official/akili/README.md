# Akili

A security-first control plane for autonomous AI operator agents: **default deny, explicit allow,
always auditable.**

Akili lets AI agents write code, operate servers and drive deployments without handing them the keys.
Agents run on your servers and dial out to the control plane. Every model call goes through its LLM
gateway, and every tool call is checked against a signed policy twice, once on the agent and once on
the control plane, before anything runs. Risky actions wait for a human, and every decision lands in a
hash-chained audit log.

This template runs the control plane: the web UI, the REST API (OpenAPI docs at `/docs`), the LLM
gateway and the agent tunnels, in one container. Add agents afterwards, on your servers or with the
[Akili Agent](../akili-agent) template.

Source: **[github.com/goakili/akili](https://github.com/goakili/akili)**.

## What gets created

- One application (`akili`, image `jkaninda/akili:0.0.2`) on port 8080, health-checked on `/healthz`
- One PostgreSQL 17 database: agents, policies, tasks, encrypted secrets and the audit log
- One dedicated Redis 8: events, presence, leases and leader election
- A route from the public URL to the app, when its domain belongs to the workspace

No volume: Akili keeps its state in PostgreSQL, so the database and the encryption key (see below) are
what you back up.

## Before you install

- **Attach the domain** of the public URL to the workspace first. The route is created only when the
  domain belongs to it; otherwise the install continues with a warning and you add the route later.
- **Pick the final public URL.** Agents connect to it, and it appears in agent install commands and in
  links. Changing it later means re-pointing every agent.

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The HTTPS URL Akili is served on, e.g. `https://akili.example.com`. |
| **Admin email / password** | The first owner account, created once on an empty database. Leave the password blank to generate one. |
| **JWT secret** | Signs session and API tokens. Rotating it signs everyone out. Generated when blank. |
| **Encryption key** | Encrypts provider keys, integration tokens and the policy-signing key. Generated when blank. **Back it up**, see below. |
| **Anthropic API key** | Optional. Seeds Claude as the default model provider on the first start. |
| **Enterprise license** | Optional. An Akili Enterprise token, bound to this deployment's Install ID. It can also be installed later. |
| **Trusted proxies** | CIDRs whose `X-Forwarded-For` is trusted, so rate limits and the audit log see real client IPs. The default covers Miabi's gateway on a private network. |

## After install

1. **Save the encryption key.** Read it from the app's environment (**Application → Environment →
   `AKILI_ENCRYPTION_KEY` → reveal**) and store it somewhere that is not this server. Without it, the
   provider keys, integration tokens and the policy-signing key in the database cannot be read, and a
   database backup does not help.
2. **Sign in** with the admin email. If the password was generated, reveal `AKILI_ADMIN_PASSWORD` the
   same way, then change it.
3. **Add a model provider.** Without an Anthropic key, Akili starts with a scripted test provider that
   understands `run: <cmd>`, `read: <path>` and a few other commands, so you can try the whole flow. Add
   a real provider under **Settings → Model providers**.
4. **Add agents** (**Agents → Add agent**): choose a policy and an autonomy level (L0 approves every
   call, L3 runs up to high risk alone; critical actions always need a human), then copy the one-time
   join token.
   - **On a server:** run the install command Akili shows. It installs a systemd service running as an
     unprivileged `akili` user. The control plane serves the agent binary, so agents always match it.
   - **On Miabi:** install the [Akili Agent](../akili-agent) template with the public URL and the token.

## Let agents operate this Miabi

Akili has native Miabi tools: deploys, rollbacks, restarts, scaling, logs, databases and backups, cron
jobs and pipelines, each with a risk level the policy checks.

1. In Miabi, create an API key for a dedicated bot user that belongs **only to the workspaces Akili
   should manage**, with the scopes `read`, `write` and `deploy`.
2. In Akili, add a Miabi integration with this Miabi's URL and the key. The key stays on the control
   plane; agents never see it.
3. **Enable** the workspaces agents may use, and scope policies by name, e.g. allow `staging/*` and
   deny `prod/*`.
4. **Watch** an app, a pattern or a whole workspace (`*`). Every deploy is then verified on real
   traffic, and failed deploys, crashes and out-of-memory kills open triage tasks. Events arrive over
   Miabi's live event stream, so no webhook or public URL is needed.

Restoring a database is critical risk and always needs a human approval.

## Akili Enterprise

The Community edition is complete and AGPL-licensed. Akili Enterprise adds teams, SAML/SCIM, approval
governance and compliance features. After install, the owner finds this deployment's **Install ID**
under **Settings → License**; paste a license issued for it there, or set it as the template's
Enterprise license input.

## Notes and limitations

- **Gateway timeouts.** Agent tunnels, the browser terminal and the live event streams are long-lived
  connections. Miabi's gateway uses Goma's default timeouts (read 30 s, write 60 s, idle 90 s), which
  close them, so agents reconnect and terminals drop. Raise them under `timeouts` in
  `/etc/miabi/goma.yml`, ideally to an hour (`3600`) and at least 15 minutes (`900`), or set
  `write: 0` for no limit. See [Reverse proxy and load balancer](https://github.com/goakili/akili#reverse-proxy-and-load-balancer).
- **Do not change the encryption key** on an existing install. It is not a rotation setting; it is the
  key to every secret already stored.
- **Install agents from the same version.** An agent installed from the control plane's command always
  matches it. With the Akili Agent template, use the template version whose image tag matches this one.
- **Version `0.0.2`** is what this template runs. Templates are immutable per version, so a newer Akili
  ships as a new template version.
