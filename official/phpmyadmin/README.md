# phpMyAdmin

The classic web interface for administering MySQL and MariaDB: browse tables,
run SQL, edit rows, and import or export dumps.

## What gets created

- One application (`phpmyadmin`, image `phpmyadmin:5.2-apache`) on port 80,
  capped at 256 MiB of memory and half a CPU
- No volume and no database: phpMyAdmin keeps no state of its own

## After install

Attach a domain and open it. The login screen has a **Server** field, because
the template enables arbitrary servers (`PMA_ARBITRARY=1`). Enter the in-network
host (and port, if not 3306) of any MySQL or MariaDB instance in the workspace,
plus that database's user and password from its **Connection** panel.

Imports up to **256 MB** are accepted (`UPLOAD_LIMIT`). Use the CLI for larger
dumps.

## Notes

Anyone who can reach the page can try credentials against every database the
container can reach. Restrict access to the route, or stop the app when you're
not using it.
