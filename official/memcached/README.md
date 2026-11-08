# Memcached

A distributed in-memory key/value cache for speeding up dynamic web applications
by caching database results, rendered fragments and sessions.

## What gets created

- One application (`memcached`, image `memcached:1.6-alpine`) listening on port
  11211 inside the workspace network
- No volumes: the cache lives in memory and starts empty on every restart

## Inputs

| Input | Notes |
|---|---|
| **Cache size (MB)** | Memory the cache may use before it evicts items (default 64). |
| **Max connections** | Simultaneous client connections (default 1024). |

## After install

Point your applications at the Memcached network alias on port **11211**. There
is nothing to attach a domain to: Memcached speaks its own protocol, not HTTP,
and has no authentication, so keep it internal to the workspace.

## Notes

- Size the app's memory limit above the cache size to leave room for connection
  buffers.
- The `1.6-alpine` tag follows Memcached 1.6.x.
