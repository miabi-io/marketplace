# Adminer

A single-file web UI for managing databases: browse tables, run SQL, edit rows,
and import or export data. It speaks PostgreSQL, MySQL, MariaDB, SQLite and
MongoDB.

Adminer holds no data of its own. It is a client that you point at a database,
which makes it a handy companion to the managed databases in your workspace.

## What gets created

- One application (`adminer`, image `adminer:5.5.1`) on port 8080
- No volumes and no database: nothing is stored on the node

## Inputs

| Input | Notes |
|---|---|
| **Default server** | Host prefilled on the login form, for example the network alias of a managed database in this workspace. Leave blank to type it each time. |
| **Theme** | Adminer skin (`pepa-linha` by default). |

## After install

Attach a domain to port **8080** and open it. Pick the database system, then
sign in with that database's host, user and password. You can copy them from the
database's **Connection** details in Miabi.

## Notes and limitations

- **Anyone who can reach the page can try to log in.** Adminer has no accounts
  of its own; it relies entirely on the database credentials. Don't attach a
  public domain unless you need to, or put the route behind authentication.
- Databases in the workspace are reachable by their network alias, so a database
  doesn't need to be exposed outside Miabi for Adminer to manage it.
