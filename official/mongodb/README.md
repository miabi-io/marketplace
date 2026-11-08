# MongoDB

A managed MongoDB 7.0 document database, provisioned as a dedicated instance.

## What gets created

- One dedicated MongoDB 7.0 database instance (`mongodb`), with no applications

The instance is managed by Miabi like any other database: it appears under
**Databases**, with its credentials, connection string, backups and lifecycle
actions.

## After install

Open the database in Miabi and copy the connection string from its
**Credentials** page. Applications in the same workspace connect through the
instance's network alias on port 27017. To browse it in a web UI, install the
**Mongo Express** template and paste the same connection string.

## Notes

- Back it up with Miabi's database backups. Data lives in the instance's managed
  storage, not in a template volume.
- Moving to a newer MongoDB major is a migration (MongoDB upgrades one major
  version at a time), not a simple tag change.
