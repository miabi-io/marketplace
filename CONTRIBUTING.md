# Contributing templates

The Miabi marketplace hosts **official** (curated) and **community**
(contributed) templates in one repo, separated by folder. Trust = which folder a
template came from, surfaced as a badge in the storefront and the Miabi console.

## Add a community template

You can develop and test a template entirely in the Miabi console first —
**Marketplace → Import**, paste your `template.yaml`, and install it into a
workspace — before adding it here (see the tip under step 2).

1. Create the directory layout under `community/<name>/`:

   ```
   community/<name>/
     metadata.yaml          # optional storefront enrichment (featured, screenshots, sourceRepo)
     README.md              # optional long description, shown on the detail page
     <version>/template.yaml  # the install manifest (apiVersion: miabi.io/v1)
   ```

   - `<name>` is the template **handle**: lowercase `^[a-z0-9][a-z0-9-]*$`,
     **unique across both `official/` and `community/`** (CI enforces this), and it
     must equal the manifest's `metadata.name`.
   - One `template.yaml` per version directory; a version is immutable once merged.
   - The manifest schema is documented in [`schema/template.schema.json`](schema/template.schema.json).
     Point your editor's YAML language server at it for inline validation.

   The manifest's `metadata` block identifies the template:

   ```yaml
   apiVersion: miabi.io/v1
   kind: Template
   metadata:
     name: my-app          # the handle — lowercase [a-z0-9-], equals the directory name
     displayName: My App   # free-text label shown in the storefront and Miabi console
     version: 1.0.0        # immutable per version directory
   ```

2. Validate locally before opening a PR:

   ```sh
   go run ./cmd/marketplace lint            # parse + validate every template, verify digests
   go run ./cmd/marketplace generate-index  # regenerate registry/index.json
   git diff --exit-code registry/index.json # must be clean (CI runs this)
   ```

   > **Tip — try it in Miabi first.** Beyond the linter, you can test the
   > template end to end: in the Miabi console open **Marketplace → Import**,
   > paste your `template.yaml`, then install it into a workspace. Installing it
   > for real is the surest way to confirm it deploys and runs as intended before
   > you open a PR.

3. Open a PR. CI re-runs the validator and the index drift check. A maintainer
   reviews and merges; merging to `main` redeploys the service, so your template
   is live immediately and Miabi installs pick it up on the next sync.

## Official templates

Changes under `official/` require core-maintainer review (enforced by
`CODEOWNERS` + branch protection). Official is the source of truth that Miabi
vendors into its binary as the offline floor. Promoting a community template to
official is a maintainer-reviewed folder move within this repo.

## Safety

Community manifests are untrusted input. The validator rejects host binds,
privileged flags, unknown fields, and malformed values. Keep templates minimal
and pin image tags.
