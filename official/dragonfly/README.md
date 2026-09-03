# Dragonfly

A drop-in replacement for Redis and Memcached that scales across cores instead of
running single-threaded. Existing Redis clients connect unchanged.

## Connecting

Dragonfly speaks the Redis protocol, which the gateway does not route, so it is
reachable only on the workspace network:

```
redis://:<password>@mb-app-<handle>:6379
```

The password is generated at install and stored as a workspace secret. Do not
attach a domain — there is nothing here for the gateway to serve.

## Memory and eviction

**Max memory** bounds what the datastore uses. Left at `0`, Dragonfly sizes
itself against the container's memory limit, so set a limit on the app if you
leave it at the default.

**Cache mode** decides what happens when that bound is reached: on, the least
valuable keys are evicted; off, writes are refused. Turn it on for a cache, and
leave it off for anything you would be unhappy to lose.

## Snapshots

Snapshots are written to the `data` volume on the schedule you set — hourly by
default — and on a clean shutdown. Clear the schedule to snapshot only on
shutdown, which means an unclean stop loses everything since the last one.
