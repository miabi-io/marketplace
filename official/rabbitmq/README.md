# RabbitMQ

A widely used message broker speaking AMQP 0-9-1, here with the management
plugin enabled for queues, exchanges, users and live message rates in the
browser.

## What gets created

- One application (`rabbitmq`, image `rabbitmq:4.4-management-alpine`). Only
  the management UI (port 15672) is published; AMQP stays on 5672 inside the
  workspace network.
- One volume, `data` → `/var/lib/rabbitmq`, holding durable queues, messages
  and definitions

## Inputs

| Input | Notes |
|---|---|
| **Username** | Initial user, administrator on the default vhost (default `admin`). |
| **Password** | Leave blank to auto-generate. It is stored as a secret env var. |

## After install

Attach a domain to open the management UI and sign in with that user. Apps in
the same workspace connect to
`amqp://<user>:<password>@<rabbitmq alias>:5672/`. Create a separate user per
application from **Admin → Users**, rather than sharing the admin account.

The initial user and password only seed an empty data volume. Change them later
in the UI, not by editing the env vars.

## Upgrades

RabbitMQ upgrades one minor series at a time and requires all stable feature
flags to be enabled first. Check **Admin → Feature Flags** before moving to a
newer template version.
