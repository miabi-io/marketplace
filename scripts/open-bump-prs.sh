#!/usr/bin/env bash
# Opens one pull request per template that has a newer upstream image; run from the repository root on main.
set -euo pipefail

base="${BASE_BRANCH:-main}"
mode="${1:-minor}"

flags=()
[[ "$mode" == "major" ]] && flags=(--major)

mapfile -t names < <(go run ./cmd/marketplace bump --all --dry-run --json "${flags[@]}" | jq -r '.[].template')
for name in "${names[@]}"; do
  git checkout -q "$base"
  git reset -q --hard "origin/$base"

  out=$(go run ./cmd/marketplace bump "$name" --json "${flags[@]}")
  to=$(jq -r '.[0].to // empty' <<<"$out")
  [[ -z "$to" ]] && continue
  branch="bump/${name}-${to}"
  if git ls-remote --exit-code --heads origin "$branch" >/dev/null; then
    echo "$branch already open, skipping"
    git reset -q --hard && git clean -qfd
    continue
  fi

  go run ./cmd/marketplace generate
  git switch -q -c "$branch"
  git add -A official community export.json registry/index.json
  changes=$(jq -r '.[0].changes[]' <<<"$out" | sed 's/^/- /')
  title="chore(${name}): bump to template ${to}"
  git commit -q -m "$title" -m "$changes"
  git push -q origin "$branch"

  body="Automated image update for \`${name}\` (template ${to}).

${changes}

Release notes are in \`$(jq -r '.[0].path' <<<"$out")/release.yaml\`; edit them before merging if the upgrade needs manual steps."
  draft=()
  if [[ "$(jq -r '.[0].breaking' <<<"$out")" == "true" ]]; then
    draft=(--draft)
    body+="

**Major upstream release.** Read the upstream upgrade notes, describe any manual steps in \`notes\`, and set \`requiresFrom\` if installs can't skip versions."
  fi
  gh pr create --base "$base" --head "$branch" --title "$title" --body "$body" "${draft[@]}"
done
git checkout -q "$base"
