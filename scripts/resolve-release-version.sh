#!/bin/sh
# Resolve what one product should publish, and whether it still needs to.
#
# Publication keys off the product's own version file rather than off which
# request file a push happened to touch. That is what makes the two version
# lines independent in practice: bump CLI_VERSION alone and only the CLI
# publishes, bump VERSION alone and only the template does. It also removes a
# whole class of failure — a push that carried an unrelated request file, or
# deleted a superseded one, used to resolve to two versions and refuse.
#
# Inputs (env):
#   INPUT_VERSION       explicit version from workflow_dispatch (optional)
#   FORGE               github (default) | gitlab
#   GH_TOKEN            credential for the lookup (github)
#   GITLAB_TOKEN        credential for the lookup (gitlab)
#   CI_API_V4_URL       API base for the lookup (gitlab)
#   TEMPLATE_REPOSITORY where the template publishes (product: template)
#   CLI_REPOSITORY      where the CLI publishes (product: cli)
#   RESOLVE_OUTPUT      file to append key=value lines to; defaults to
#                       GITHUB_OUTPUT so the GitHub workflow needs no change
# Arg:
#   $1  product: template | cli
#
# Only the "has this already shipped" lookup differs by forge. Which version to
# publish, and which request describes it, stay written once: two pipelines
# disagreeing about the version is the failure this file exists to prevent,
# and it is not a failure either pipeline would report — both would succeed.
#
# Neither repository is written here. A publish target is a property of the
# pipeline that runs this, not of the product: the same source publishes to a
# different pair of repositories depending on where it was pushed, and a
# hardcoded name would silently make one of those two the real one.
#
# Prints "version=<v>" and "publish=true|false" for $GITHUB_OUTPUT, and a
# human line for the log. Publishing an already-published version is not an
# error; it is the ordinary case for the product that did not change.
set -eu

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

product=${1:-}
source_root=$(git rev-parse --show-toplevel 2>/dev/null) || fail 'not a git source checkout'

case "$product" in
  template)
    version=${INPUT_VERSION:-$(cat "$source_root/VERSION")}
    repository=${TEMPLATE_REPOSITORY:-}
    [ -n "$repository" ] || fail 'set TEMPLATE_REPOSITORY to the repository the template publishes to'
    tag="v$version"
    ;;
  cli)
    version=${INPUT_VERSION:-$(cat "$source_root/CLI_VERSION")}
    repository=${CLI_REPOSITORY:-}
    [ -n "$repository" ] || fail 'set CLI_REPOSITORY to the repository the CLI publishes to'
    tag="cli-v$version"
    ;;
  *) fail 'usage: resolve-release-version.sh <template|cli>' ;;
esac

case "$version" in ''|v*|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $version" ;; esac

output_file=${RESOLVE_OUTPUT:-${GITHUB_OUTPUT:-}}

published() {
  case "${FORGE:-github}" in
    github) gh release view "$tag" --repo "$repository" >/dev/null 2>&1 ;;
    gitlab)
      [ -n "${CI_API_V4_URL:-}" ] || fail 'set CI_API_V4_URL for a gitlab lookup'
      [ -n "${GITLAB_TOKEN:-}" ] || fail 'set GITLAB_TOKEN for a gitlab lookup'
      encoded=$(printf '%s' "$repository" | sed 's,/,%2F,g')
      # A release that exists answers 200 and one that does not answers 404.
      # Any other status is a lookup that did not happen — treating it as "not
      # published" would republish over a release that is already out.
      status=$(curl --silent --output /dev/null --write-out '%{http_code}' \
        --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
        "$CI_API_V4_URL/projects/$encoded/releases/$tag")
      case "$status" in
        200) return 0 ;;
        404) return 1 ;;
        *) fail "release lookup for $tag answered $status" ;;
      esac
      ;;
    *) fail 'unknown FORGE: expected github or gitlab' ;;
  esac
}

if published; then
  # Already out. Say so and stop looking for notes: a release that has shipped
  # does not need a request, and demanding one would make the product that did
  # not change this time fail the run for the product that did.
  printf 'resolved: %s %s already published as %s\n' "$product" "$version" "$tag"
  if [ -n "$output_file" ]; then
    printf 'version=%s\npublish=false\n' "$version" >> "$output_file"
  fi
  exit 0
fi

# Each product carries its own request, so a release of one describes only what
# that one changed, and a reader of a CLI release is not handed template notes.
request="$source_root/release/requests/$product/$version.md"
[ -f "$request" ] || fail "no $product release request: release/requests/$product/$version.md"
publish=true

printf 'resolved: %s %s (tag %s, publish %s)\n' "$product" "$version" "$tag" "$publish"
if [ -n "$output_file" ]; then
  printf 'version=%s\npublish=%s\nrequest=%s\n' "$version" "$publish" "release/requests/$product/$version.md" >> "$output_file"
fi
