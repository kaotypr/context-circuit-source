#!/bin/sh
set -eu

usage() {
  printf 'usage: sh scripts/build-dist.sh [version] [output-dir]\n' >&2
  exit 2
}

[ "$#" -le 2 ] || usage

source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
version=${1:-v0.5.0}
output_dir=${2:-$source_root/dist}

mkdir -p "$output_dir"
staging_dir=$(mktemp -d "$output_dir/.staging.XXXXXX")
cleanup() { rm -rf "$staging_dir"; }
trap cleanup EXIT HUP INT TERM

sh "$source_root/scripts/release-artifact.sh" "$staging_dir" "$output_dir" "$version"

printf 'dist_dir: %s\n' "$output_dir"
