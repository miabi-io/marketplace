# Mongo Express

A web interface to administer MongoDB: browse databases and collections, query
and edit documents, and manage indexes.

Mongo Express doesn't provision a database. It connects to an existing MongoDB,
typically one created from the MongoDB template in this workspace.

## What gets created

- One application (`mongo-express`, image `mongo-express:1.0.2-20`) on port 8081,
  limited to 256 MiB of memory and half a CPU core
- No volumes and no database

## Inputs

| Input | Notes |
|---|---|
| **MongoDB connection string** | Copy it from the MongoDB instance's **Credentials** page, e.g. `mongodb://admin:password@<alias>:27017/?authSource=admin`. Stored as a secret. |
| **UI username / password** | The login that protects the Mongo Express page. Leave the password blank to auto-generate. |

## After install

Attach a domain to port **8081**, open it, and sign in with the UI username and
password. Admin mode is on, so every database the connection string's user can
see is listed.

## Notes

The UI login is the only thing between the internet and full write access to your
data. Use a strong password, and avoid a public domain if you can reach the app
another way.
