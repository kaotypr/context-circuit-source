#!/bin/sh
set -eu

usage() {
  printf 'usage: sh scripts/build-dist.sh [version] [output-dir]\n' >&2
  exit 2
}

[ "$#" -le 2 ] || usage

source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

# The dist artifact is the context-circuit-template; its version is the template
# release identity. Default it from the single source of truth
# (wrapper/manifest.yaml template_version), matching the published archive name
# in publish-template.sh. An explicit [version] argument still overrides.
version=${1:-}
if [ -z "$version" ]; then
  manifest="$source_root/wrapper/manifest.yaml"
  [ -f "$manifest" ] || { printf 'FAIL: missing manifest: %s\n' "$manifest" >&2; exit 1; }
  template_version=$(sed -n 's/^template_version:[[:space:]]*//p' "$manifest" | head -n1)
  [ -n "$template_version" ] || { printf 'FAIL: no template_version in %s\n' "$manifest" >&2; exit 1; }
  version=v$template_version
fi
output_dir=${2:-$source_root/dist}

mkdir -p "$output_dir"
staging_dir=$(mktemp -d "$output_dir/.staging.XXXXXX")
cleanup() { rm -rf "$staging_dir"; }
trap cleanup EXIT HUP INT TERM

sh "$source_root/scripts/release-artifact.sh" "$staging_dir" "$output_dir" "$version"

printf 'dist_dir: %s\n' "$output_dir"
