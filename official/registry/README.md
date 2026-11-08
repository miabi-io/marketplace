# Docker Registry

A private OCI image registry built on CNCF Distribution (the reference
`registry` image). Push and pull images with any Docker- or OCI-compatible
client.

## What gets created

- One application (`registry`, image `registry:3.1`) on port 5000
- One volume, `data` → `/var/lib/registry`, holding every pushed layer and
  manifest

## Inputs

| Input | Notes |
|---|---|
| **Allow deletes** | Permits `DELETE` of manifests and blobs (default `true`). Garbage collection needs it. |

## Security: there is no authentication

This registry ships **without authentication**: anyone who can reach it can
push and pull. Either keep it internal (apps in the workspace reach it over the
workspace network), or attach a domain only behind access control.

## After install

With a domain attached, tag and push as usual:

```sh
docker tag myapp:1.0 registry.example.com/myapp:1.0
docker push registry.example.com/myapp:1.0
```

## Freeing space

Deleting a tag does not free disk by itself. With deletes allowed, run
`registry garbage-collect /etc/distribution/config.yml` in the container
(ideally while nothing is pushing) to reclaim unreferenced layers.
