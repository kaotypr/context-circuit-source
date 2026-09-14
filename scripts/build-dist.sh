#!/bin/sh
# Build into a new directory. Existing output is never cleaned or replaced.
set -eu
[ "$#" -le 2 ] || { printf 'usage: sh scripts/build-dist.sh [version] [new-output-dir]\n' >&2; exit 2; }
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
version=${1:-v$(cat "$source_root/VERSION")}
output_dir=${2:-$source_root/dist}
[ ! -e "$output_dir" ] && [ ! -L "$output_dir" ] || { printf 'FAIL: output already exists: %s\n' "$output_dir" >&2; exit 1; }
staging_dir=$(mktemp -d)
trap 'rm -rf "$staging_dir"' EXIT HUP INT TERM
sh "$source_root/scripts/release-artifact.sh" "$staging_dir" "$output_dir" "$version"
printf 'dist_dir: %s\n' "$output_dir"
