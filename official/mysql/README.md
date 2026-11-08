# MySQL

A managed MySQL 8.4 (LTS) database, provisioned as a dedicated instance.

## What gets created

- One dedicated MySQL 8.4 database instance (`mysql`), with no applications

The instance is managed by Miabi like any other database: it appears under
**Databases**, with its credentials, connection details, backups and lifecycle
actions.

## After install

Open the database in Miabi to copy its host, port, user and password.
Applications in the same workspace connect through the instance's network alias
on port 3306. For a web UI, install the **Adminer** or **phpMyAdmin** template
and sign in with the same credentials.

## Notes

- Most app templates that need MySQL create their own database automatically.
  Install this one when you want a standalone instance for your own
  applications.
- Back it up with Miabi's database backups.
