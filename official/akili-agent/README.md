# Akili Agent

An [Akili](https://github.com/goakili/akili) agent, run as a Miabi app. The agent dials out to your Akili
control plane over an encrypted tunnel and runs chat and task sessions within the signed policy it
receives. It holds no model or forge credentials: every model call and every git push goes through the
control plane, where each tool call is checked again, approved when needed and audited.

You need a running control plane first, for example from the [Akili](../akili) template.

Source: **[github.com/goakili/akili](https://github.com/goakili/akili)**.

## What gets created

- One application (`akili-agent`, image `jkaninda/akili-agent:0.0.2`). It exposes **no ports**: the
  agent only opens an outbound connection to the control plane.
- One volume, `state`, mounted at `/var/lib/akili-agent`. It holds the agent's Ed25519 identity and its
  working directory, so restarts and redeploys reconnect without enrolling again.

## Before you install

In Akili, create the agent (**Agents → Add agent**): give it a name, pick its policy and autonomy level
(L0 approves every call, L3 runs up to high risk alone; critical actions always need a human), and copy
the **join token**. It is shown once.

## Inputs

| Input | Notes |
|---|---|
| **Control plane URL** | The public URL of your Akili control plane, e.g. `https://akili.example.com`. |
| **Join token** | The one-time token from **Agents → Add agent**. It is used once to enroll; after that the agent authenticates with its own key from the volume. |
| **Control plane CA** | Only for a control plane with a self-signed or private-CA certificate: the CA certificate as PEM, or base64-encoded PEM (`base64 < ca.pem \| tr -d '\n'`). |

## After install

The agent shows up **online** under **Agents** in Akili within a few seconds. Give it work from a chat
or a task.

- **Keep the `state` volume.** Deleting it deletes the agent's identity; enrolling again needs a new
  join token from Akili.
- **To operate Miabi apps,** add a Miabi integration in Akili and enable the workspaces agents may use.
  The Miabi API key stays on the control plane, never in this container. See the
  [Akili](../akili) template for the steps.

## What this agent can do

The container runs as an unprivileged user (UID 10001) on Alpine, with `bash`, `curl`, `git` and `jq`.
Inside its own container, it can:

- **Code:** clone projects through the control plane's git proxy, work on a branch per task, and open
  pull requests. Pushes are allowed only to the session's own branch, and forge credentials stay on the
  control plane.
- **Run commands** within its policy: shell, file reads and writes in its working directory, and HTTP
  checks.
- **Drive Miabi** through the Miabi integration: deploys, rollbacks, logs, scaling and databases, each
  checked and audited by the control plane.

It **cannot** reach the Miabi host or other containers' internals:

- **Host tools** (`service_status`, `journal_logs`, `docker_ps`, `service_restart`, …) have no host to
  act on here. To manage a server, install the agent on that server with the command Akili shows.
- **`sandbox_exec`** needs a Docker daemon, which this container doesn't have, so it reports that it is
  unavailable. Tests can still run through `shell`, which is high risk and needs approval below L3. For
  sandboxed coding tasks, use an agent on a dedicated build host with rootless Docker.

## Notes and limitations

- **Match the control plane's version.** This template pins the agent image to `0.0.2`. When the
  control plane is upgraded, move agents to the template version with the matching tag. An agent
  installed with the control plane's own install command always matches it.
- **One agent per install.** Each install is one agent with one identity. For several agents, create
  each one in Akili and install the template once per token.
- **Version `0.0.2`** is what this template runs. Templates are immutable per version, so a newer agent
  ships as a new template version.
