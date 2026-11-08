# Redis

A standalone, managed Redis 8 instance for caching, queues, sessions and
pub/sub.

## What gets created

- One **dedicated** Redis instance (`redis:8-alpine`), password-protected and
  managed by Miabi like any database you add from the **Databases** page

There is no application and there are no inputs.

## After install

Open the instance in **Databases** to find its in-network host, port and
password, plus a ready-made `redis://` URI. Reference them from your apps'
environment.

Persistence, backups and resource limits are managed from the same page. Treat
Redis as a cache unless you have confirmed its persistence settings suit your
data.
