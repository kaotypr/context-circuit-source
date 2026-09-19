#!/bin/sh
# Build into a directory. The default output under dist/ is clean-rebuilt each
# run. An explicit output directory must be new and is never cleaned or replaced.
#
# What this builds is the workspace a user receives, which is not the tree the
# published template repository carries. That repository's own landing page —
# its README, license, changelog, conduct, contributing and security pages, and
# the artwork the guide renders — is restored by publish-template.sh over this
# artifact and reaches no workspace, so it is absent here by design.
set -eu
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
[ "$#" -le 2 ] || { printf 'usage: sh scripts/build-dist.sh [version] [new-output-dir]\n' >&2; exit 2; }
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
version=${1:-v$(cat "$source_root/VERSION")}
case "$version" in ''|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $version" ;; esac
if [ "$#" -ge 2 ]; then
  output_dir=$2
  [ ! -e "$output_dir" ] && [ ! -L "$output_dir" ] || fail "output already exists: $output_dir"
else
  # Own subdirectory per build so a workspace rebuild never removes CLI assets.
  output_dir=$source_root/dist/workspace-${version#v}
  case "$output_dir" in ''|/|"$source_root"|"$source_root"/) fail "refusing to clean: $output_dir" ;; esac
  rm -rf "$output_dir"
fi
staging_dir=$(mktemp -d)
trap 'rm -rf "$staging_dir"' EXIT HUP INT TERM
sh "$source_root/scripts/release-artifact.sh" "$staging_dir" "$output_dir" "$version"
printf 'dist_dir: %s\n' "$output_dir"
printf 'contents: the workspace a user receives, without the landing page the\n'
printf '          published repository owns and publication restores\n'
