#!/bin/sh
# Maintainer-only helper: stage a clean cloneable template from a tracked revision.
# Usage: sh scripts/release-artifact.sh <staging-dir> <output-dir> <version>
set -eu

DESTINATION_REPO='kaotypr/context-circuit-release'
DESTINATION_REF='main'

usage() {
  printf 'usage: sh scripts/release-artifact.sh <staging-dir> <output-dir> <version>\n' >&2
  exit 2
}

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

path_unsafe() {
  candidate=$1
  case "$candidate" in
    ''|/*|..|../*|*/..|*/../*) return 0 ;;
  esac
  return 1
}

resolve_source_root() {
  if git rev-parse --show-toplevel >/dev/null 2>&1; then
    git rev-parse --show-toplevel
    return 0
  fi
  script_parent=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
  git -C "$script_parent" rev-parse --show-toplevel 2>/dev/null || \
    fail 'not a git repository; run from the Context Circuit source checkout'
}

if [ "$#" -ne 3 ]; then
  usage
fi

staging_dir=$1
output_dir=$2
version=$3

case "$staging_dir" in
  ''|.) fail 'staging-dir must be an explicit directory' ;;
esac
case "$output_dir" in
  ''|.) fail 'output-dir must be an explicit directory' ;;
esac
case "$version" in
  ''|*/*|.*|-*|*' '*|*..*) fail "invalid version: $version" ;;
esac

source_root=$(resolve_source_root)
manifest="$source_root/scripts/release-manifest.txt"
test -f "$manifest" || fail "missing release manifest: $manifest"

source_sha=$(git -C "$source_root" rev-parse HEAD) || fail 'unable to resolve source revision'
git -C "$source_root" rev-parse --verify --quiet HEAD >/dev/null || \
  fail 'source repository has no tracked revision'

mkdir -p "$staging_dir" "$output_dir"
stage_tree="$staging_dir/tree"
if [ -e "$stage_tree" ]; then
  fail "staging tree already exists: $stage_tree"
fi
mkdir -p "$stage_tree"

# Stage tracked files only. Dirty and ignored working-tree state cannot enter.
git -C "$source_root" archive --format=tar HEAD | tar -xf - -C "$stage_tree"

required_files=
exclude_paths=
while IFS= read -r line || [ -n "$line" ]; do
  case "$line" in
    ''|\#*) continue ;;
  esac
  set -- $line
  kind=${1:-}
  relpath=${2:-}
  [ -n "$kind" ] && [ -n "$relpath" ] || fail "invalid manifest line: $line"
  path_unsafe "$relpath" && fail "unsafe manifest path: $relpath"
  case "$kind" in
    required)
      required_files="$required_files $relpath"
      ;;
    exclude)
      exclude_paths="$exclude_paths $relpath"
      ;;
    *)
      fail "unknown manifest directive: $kind"
      ;;
  esac
done < "$manifest"

for relpath in $exclude_paths; do
  target="$stage_tree/$relpath"
  if [ -e "$target" ] || [ -L "$target" ]; then
    rm -rf "$target"
  fi
done

for relpath in $required_files; do
  test -e "$stage_tree/$relpath" || fail "missing required file: $relpath"
done

for relpath in $exclude_paths; do
  if [ -e "$stage_tree/$relpath" ] || [ -L "$stage_tree/$relpath" ]; then
    fail "forbidden path remains: $relpath"
  fi
done

# Package-manager files are never silently stripped; their presence fails packaging.
pm_names='package.json package-lock.json npm-shrinkwrap.json yarn.lock pnpm-lock.yaml bun.lock bun.lockb tsconfig.json'
for name in $pm_names; do
  found=$(find "$stage_tree" -name "$name" -print)
  if [ -n "$found" ]; then
    fail "forbidden package-manager file in artifact: $name"
  fi
done
found_modules=$(find "$stage_tree" -type d -name node_modules -print)
if [ -n "$found_modules" ]; then
  fail 'forbidden package-manager file in artifact: node_modules'
fi

artifact_name="context-circuit-$version"
artifact_dir="$output_dir/$artifact_name"
if [ -e "$artifact_dir" ]; then
  fail "output artifact already exists: $artifact_dir"
fi
mv "$stage_tree" "$artifact_dir"

archive_path=
if command -v gzip >/dev/null 2>&1; then
  archive_path="$output_dir/$artifact_name.tar.gz"
  if [ -e "$archive_path" ]; then
    fail "output archive already exists: $archive_path"
  fi
  (CDPATH= cd -- "$artifact_dir" && tar -cf - .) | gzip -n > "$archive_path"
fi

printf 'version: %s\n' "$version"
printf 'source_revision: %s\n' "$source_sha"
printf 'destination: %s %s\n' "$DESTINATION_REPO" "$DESTINATION_REF"
printf 'artifact_dir: %s\n' "$artifact_dir"
if [ -n "$archive_path" ]; then
  printf 'artifact_archive: %s\n' "$archive_path"
fi
printf 'remaining_gate: publication\n'
printf 'inventory:\n'
(CDPATH= cd -- "$artifact_dir" && find . -type f -print | sort | sed 's|^\./||')
