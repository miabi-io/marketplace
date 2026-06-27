# Nginx

A lightweight, high-performance web server for serving static sites and assets.

This official template deploys a single Nginx container exposing port 80 over
HTTP. Attach a domain in Miabi to get automatic TLS via Goma Gateway.

## What gets created

- One application (`nginx`, image `nginx:1.30-alpine`)

## After install

Mount a volume at `/usr/share/nginx/html` (or bake your own image) to serve your
own content.
