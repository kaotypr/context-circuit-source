#!/bin/sh
set -eu

DESTINATION_REPO='kaotypr/context-circuit-release'
DESTINATION_REF='main'
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
usage() { printf 'usage: sh scripts/release-artifact.sh <staging-dir> <output-dir> <version>\n' >&2; exit 2; }
path_unsafe() { case "$1" in ''|/*|..|../*|*/..|*/../*) return 0 ;; esac; return 1; }

[ "$#" -eq 3 ] || usage
staging_dir=$1
output_dir=$2
version=$3
case "$version" in ''|*/*|.*|-*|*' '*|*..*) fail "invalid version: $version" ;; esac
source_root=$(git rev-parse --show-toplevel 2>/dev/null) || fail 'not a git source checkout'
manifest="$source_root/scripts/release-manifest.txt"
[ -f "$manifest" ] || fail "missing release manifest: $manifest"
source_sha=$(git -C "$source_root" rev-parse HEAD 2>/dev/null || printf uncommitted)
if git -C "$source_root" status --porcelain --untracked-files=all | grep . >/dev/null 2>&1; then source_state=dirty; else source_state=clean; fi

mkdir -p "$staging_dir" "$output_dir"
stage_tree="$staging_dir/tree"
[ ! -e "$stage_tree" ] || fail "staging tree already exists: $stage_tree"
mkdir -p "$stage_tree"

(CDPATH= cd "$source_root" && tar -cf - \
  .gitignore .agents agents docs wrapper) | tar -xf - -C "$stage_tree"

cp "$source_root/wrapper/adapters/AGENTS.md" "$stage_tree/AGENTS.md"
cp "$source_root/wrapper/adapters/CLAUDE.md" "$stage_tree/CLAUDE.md"
cp "$source_root/wrapper/adapters/WORKFLOW.md" "$stage_tree/WORKFLOW.md"
cp "$source_root/wrapper/adapters/README.md" "$stage_tree/README.md"
cp "$source_root/template/workspace.yaml" "$stage_tree/workspace.yaml"
mkdir -p "$stage_tree/context" "$stage_tree/sources" "$stage_tree/plans"
cp -R "$source_root/template/context/." "$stage_tree/context/"
cp "$source_root/template/sources/README.md" "$stage_tree/sources/README.md"
cp "$source_root/template/plans/README.md" "$stage_tree/plans/README.md"

required_files=''
exclude_paths=''
while IFS= read -r line || [ -n "$line" ]; do
  case "$line" in ''|'#'*) continue ;; esac
  set -- $line
  kind=${1:-}; relpath=${2:-}
  [ -n "$kind" ] && [ -n "$relpath" ] || fail "invalid manifest line: $line"
  path_unsafe "$relpath" && fail "unsafe manifest path: $relpath"
  case "$kind" in
    required) required_files="$required_files $relpath" ;;
    exclude) exclude_paths="$exclude_paths $relpath" ;;
    *) fail "unknown manifest directive: $kind" ;;
  esac
done < "$manifest"

for relpath in $exclude_paths; do
  case "$relpath" in
    *) target="$stage_tree/$relpath"; if [ -e "$target" ] || [ -L "$target" ]; then rm -rf "$target"; fi ;;
  esac
done
for relpath in $required_files; do
  [ -e "$stage_tree/$relpath" ] || fail "missing required file: $relpath"
done
for fixture in \
  wrapper/contracts/schemas/context-receipt.yaml \
  wrapper/contracts/schemas/delegation.yaml \
  wrapper/contracts/schemas/plan.yaml \
  wrapper/contracts/schemas/task.yaml \
  docs/gates.md docs/templates/plan.md docs/templates/plan.yaml \
  docs/templates/task.md docs/templates/prd.md; do
  [ -f "$stage_tree/$fixture" ] || fail "missing canonical packet fixture: $fixture"
done
for forbidden_path in repositories.local.yaml repositories; do
  [ ! -e "$stage_tree/$forbidden_path" ] || fail "forbidden repository state in artifact: $forbidden_path"
done
[ ! -e "$stage_tree/plans/context-circuit-plans" ] || fail 'maintainer plan stack leaked into artifact'
find "$stage_tree" -type f \( -name repositories.local.yaml -o -name '*.credentials' \) -print -quit | grep . && fail 'forbidden repository or credential file' || :
for skill_dir in "$stage_tree"/.agents/skills/cc-*; do
  [ -d "$skill_dir" ] || continue
  skill_name=${skill_dir##*/}
  case "$skill_name" in
    cc-entry|cc-next|cc-plan|cc-execute|cc-verify|cc-gates|cc-upgrade) ;;
    *) fail "unexpected legacy skill remains: $skill_name" ;;
  esac
done
for name in package.json package-lock.json npm-shrinkwrap.json yarn.lock pnpm-lock.yaml bun.lock bun.lockb tsconfig.json; do
  find "$stage_tree" -name "$name" -print -quit | grep . && fail "forbidden package-manager file: $name" || :
done
find "$stage_tree" -type d -name node_modules -print -quit | grep . && fail 'forbidden node_modules' || :

artifact_name="context-circuit-$version"
artifact_dir="$output_dir/$artifact_name"
[ ! -e "$artifact_dir" ] || fail "output artifact already exists: $artifact_dir"
mv "$stage_tree" "$artifact_dir"
archive_path="$output_dir/$artifact_name.tar.gz"
(CDPATH= cd "$artifact_dir" && tar -cf - .) | gzip -n > "$archive_path"

printf 'version: %s\n' "$version"
printf 'wrapper_version: 1.0.0\n'
printf 'source_revision: %s\n' "$source_sha"
printf 'source_state: %s\n' "$source_state"
printf 'destination: %s %s\n' "$DESTINATION_REPO" "$DESTINATION_REF"
printf 'artifact_dir: %s\n' "$artifact_dir"
printf 'artifact_archive: %s\n' "$archive_path"
printf 'remaining_gate: publication\n'
printf 'inventory:\n'
(CDPATH= cd "$artifact_dir" && find . -type f -print | sort | sed 's|^\./||')
