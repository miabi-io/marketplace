# Ghost

A professional publishing platform for blogs, newsletters and paid memberships,
with a modern editor and built-in SEO.

## What gets created

- One application (`ghost`, image `ghost:6.75-alpine`) on port 2368, limited to
  512 MiB of memory and half a CPU core
- A MySQL 8.4 database (`db`), shared with a compatible instance when one exists
  but always with its own database and user
- One volume, `content` → `/var/lib/ghost/content`, for images, themes and files
- A route for the **Public URL** when its domain belongs to this workspace

## Inputs

| Input | Notes |
|---|---|
| **Public URL** | The full address the blog is served on, e.g. `https://blog.example.com`. Ghost builds every link, feed and email from it, so it must match the domain readers visit. |

## After install

Open `<Public URL>/ghost` to reach the admin panel. The first visit asks you to
create the site and its owner account; do this straight after installing, before
anyone else finds the page.

## Notes and limitations

- **Email needs configuring.** Ghost sends staff invitations and sign-in links
  over SMTP (`mail__*` environment variables), and sends newsletters in bulk
  through Mailgun. Add the relevant settings to the app's environment before
  inviting staff or enabling members.
- **The memory cap is modest.** Raise it in the app's resources if you run a
  large site or many members.

## Where data lives

Posts, members and settings are in the MySQL database. Images, themes and
uploaded files are on the `content` volume. Back up both.
