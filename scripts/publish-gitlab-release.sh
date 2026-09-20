#!/bin/sh
# Upload a release's assets and create the release that links them.
#
# A GitLab release links to files, it does not store them, so the assets go to
# the generic package registry first and the release records where they landed.
# Those links are exactly what `install.sh --source gitlab` reads: it asks the
# release for an asset by name and fetches that asset's own url with the same
# PRIVATE-TOKEN, so a file uploaded here is reachable by the installer without
# any second credential or public path.
#
# Inputs (env):
#   CI_API_V4_URL  API base
#   GITLAB_TOKEN   token that can write packages and releases on --project
set -eu
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

project= package= version= tag= name= notes= ref= assets=
while [ "$#" -gt 0 ]; do
  case "$1" in
    --project|--package|--version|--tag|--name|--notes|--ref|--assets) ;;
    *) fail "unknown option: $1" ;;
  esac
  [ "$#" -ge 2 ] || fail "$1 requires a value"
  case "$1" in
    --project) project=$2 ;; --package) package=$2 ;; --version) version=$2 ;;
    --tag) tag=$2 ;; --name) name=$2 ;; --notes) notes=$2 ;;
    --ref) ref=$2 ;; --assets) assets=$2 ;;
  esac
  shift 2
done
for required in project package version tag name notes assets; do
  eval "value=\$$required"
  [ -n "$value" ] || fail "--$required is required"
done
[ -f "$notes" ] || fail "missing notes file: $notes"
[ -d "$assets" ] || fail "missing assets directory: $assets"
[ -n "${CI_API_V4_URL:-}" ] || fail 'set CI_API_V4_URL'
[ -n "${GITLAB_TOKEN:-}" ] || fail 'set GITLAB_TOKEN'
command -v jq >/dev/null 2>&1 || fail 'jq is required to build the release payload'

# A project may be given as a numeric id or as a group/project path, which the
# API expects percent-encoded rather than literal in the path segment.
case "$project" in
  *[!0-9]*) encoded=$(printf '%s' "$project" | sed 's,/,%2F,g') ;;
  *) encoded=$project ;;
esac
api="$CI_API_V4_URL/projects/$encoded"

request() { curl --fail --silent --show-error --header "PRIVATE-TOKEN: $GITLAB_TOKEN" "$@"; }

links=$(jq -n '[]')
for path in "$assets"/*; do
  [ -f "$path" ] || continue
  file=${path##*/}
  # The registry rejects a name it cannot put in a path, and a file this
  # release cannot name is a file the installer could never ask for.
  case "$file" in *[!A-Za-z0-9._-]*) fail "unexpected asset name: $file" ;; esac
  url="$api/packages/generic/$package/$version/$file"
  request --request PUT --upload-file "$path" "$url" >/dev/null ||
    fail "upload failed: $file"
  printf 'uploaded: %s\n' "$file"
  links=$(printf '%s' "$links" | jq --arg name "$file" --arg url "$url" \
    '. + [{name: $name, url: $url, link_type: "package"}]')
done
[ "$(printf '%s' "$links" | jq 'length')" -gt 0 ] || fail "no assets found in $assets"

# Built with jq rather than by hand: release notes are arbitrary prose, and a
# quote or backslash in them would otherwise produce a malformed request that
# fails after the assets are already uploaded.
payload=$(jq -n \
  --arg name "$name" \
  --arg tag "$tag" \
  --rawfile description "$notes" \
  --argjson links "$links" \
  '{name: $name, tag_name: $tag, description: $description, assets: {links: $links}}')
if [ -n "$ref" ]; then
  payload=$(printf '%s' "$payload" | jq --arg ref "$ref" '. + {ref: $ref}')
fi

request --request POST --header 'Content-Type: application/json' \
  --data "$payload" "$api/releases" >/dev/null ||
  fail "could not create release $tag"
printf 'released: %s on %s\n' "$tag" "$project"
