# Keycloak

Open-source identity and access management: single sign-on over OpenID Connect,
OAuth 2.0 and SAML, with user federation (LDAP, Active Directory), social login
and fine-grained authorization.

## What gets created

- One application (`keycloak`, image `quay.io/keycloak/keycloak:26.8.2`) on port
  8080, running `start` (production mode)
- A PostgreSQL 17 database (`db`) holding realms, clients, users and sessions
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The address Keycloak is served on, e.g. `https://auth.example.com`. It becomes the issuer in every token and the base of every redirect, so changing it later breaks existing clients. |
| **Admin username / password** | A temporary bootstrap admin created on first start (`admin` by default). Leave the password blank to auto-generate. |

## After install

The first boot builds Keycloak's optimized server and runs migrations, so it can
take a minute or two before it reports healthy. Then open the Public URL and sign
in to the admin console with the bootstrap credentials.

That bootstrap account is meant to be temporary. In the `master` realm, create a
permanent admin user, then delete the bootstrap one.

## Notes and limitations

- **TLS terminates at the Miabi gateway.** Keycloak serves plain HTTP and trusts
  the gateway's `X-Forwarded-*` headers to build `https://` URLs.
- **Health endpoints are on the management port (9000)**, which stays inside the
  container.

## Where data lives

Everything Keycloak stores is in the PostgreSQL database; there is no volume. Back
up the database.
