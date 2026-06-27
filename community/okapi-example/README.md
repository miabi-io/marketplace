# Okapi Example

A reference application for the [Okapi](https://github.com/jkaninda/okapi) Go API
framework, showcasing middleware, routing, real-time communication (SSE), and
automatic OpenAPI documentation.

## What gets created

- One application (`jkaninda/okapi-example`) listening on port **8080** (HTTP)

## Notes

- `JWT_SECRET` is auto-generated and stored encrypted; you never need to set it.
- API docs are enabled (`ENABLE_DOCS=true`) — browse them at `/docs` once the app
  is running and a domain is attached.
