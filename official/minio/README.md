# MinIO

High-performance, S3-compatible object storage with a web console.

> **Deprecated.** MinIO no longer publishes public container images, so new
> installs of this template can't pull the image. Existing installs keep running;
> the template stays in the catalog only so they can still be managed.

## What it deployed

- One application (`minio`, image `minio/minio:RELEASE.2025-04-22T22-12-26Z`)
  with the S3 API on port 9000 and the web console on port 9001
- One volume, `data` → `/data`, holding every bucket and object
- A root user and password taken from the install inputs

## For existing installs

Your data stays on the `data` volume. Keep backing it up. To move off MinIO,
copy buckets to another S3-compatible store with an S3 client such as
`mc mirror` or `rclone sync`, then point your applications at the new endpoint.
