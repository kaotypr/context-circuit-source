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
#   INPUT_VERSION  explicit version from workflow_dispatch (optional)
#   GH_TOKEN       credential for the "is it already published" lookup
# Arg:
#   $1  product: template | cli
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
    repository=kaotypr/context-circuit
    tag="v$version"
    ;;
  cli)
    version=${INPUT_VERSION:-$(cat "$source_root/CLI_VERSION")}
    repository=kaotypr/context-circuit-source
    tag="cli-v$version"
    ;;
  *) fail 'usage: resolve-release-version.sh <template|cli>' ;;
esac

case "$version" in ''|v*|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $version" ;; esac

if gh release view "$tag" --repo "$repository" >/dev/null 2>&1; then
  # Already out. Say so and stop looking for notes: a release that has shipped
  # does not need a request, and demanding one would make the product that did
  # not change this time fail the run for the product that did.
  printf 'resolved: %s %s already published as %s\n' "$product" "$version" "$tag"
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    printf 'version=%s\npublish=false\n' "$version" >> "$GITHUB_OUTPUT"
  fi
  exit 0
fi

# Each product carries its own request, so a release of one describes only what
# that one changed, and a reader of a CLI release is not handed template notes.
request="$source_root/release/requests/$product/$version.md"
[ -f "$request" ] || fail "no $product release request: release/requests/$product/$version.md"
publish=true

printf 'resolved: %s %s (tag %s, publish %s)\n' "$product" "$version" "$tag" "$publish"
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  printf 'version=%s\npublish=%s\nrequest=%s\n' "$version" "$publish" "release/requests/$product/$version.md" >> "$GITHUB_OUTPUT"
fi
