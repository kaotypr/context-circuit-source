#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

stage=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-stage.XXXXXX")
out=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-out.XXXXXX")
build_out=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-build.XXXXXX")
trap 'rm -rf "$stage" "$out" "$build_out"' EXIT HUP INT TERM
result=$(sh "$ROOT/scripts/release-artifact.sh" "$stage" "$out" v1.0.0)
artifact="$out/context-circuit-v1.0.0"
require_file "$artifact/AGENTS.md"
require_file "$artifact/WORKFLOW.md"
require_file "$artifact/README.md"
require_file "$artifact/workspace.yaml"
require_file "$artifact/wrapper/manifest.yaml"
require_file "$artifact/wrapper/runtime/engine.sh"
require_file "$artifact/wrapper/migrations/upgrade.sh"
require_file "$artifact/context/PROJECT.md"
require_file "$artifact/sources/README.md"
require_file "$artifact/plans/README.md"
contains "$artifact/workspace.yaml" 'name: uninitialized-workspace'
contains "$artifact/AGENTS.md" 'small safety spine'
printf '%s\n' "$result" | grep -F 'source_state: dirty' >/dev/null || fail 'release did not record dirty source state'
test ! -e "$artifact/.runtime" || fail 'runtime leaked into artifact'
test ! -e "$artifact/test" || fail 'semantic tests leaked into artifact'
test ! -e "$artifact/PLAN.md" || fail 'maintainer plan leaked into artifact'
test ! -e "$artifact/sources/context-circuit-design" || fail 'source design material leaked into artifact'
test ! -e "$artifact/template" || fail 'source template directory leaked into artifact'
test ! -e "$artifact/sources/secret.txt" || fail 'source inbox leaked into artifact'
for skill_dir in "$artifact"/.agents/skills/cc-*; do
  [ -d "$skill_dir" ] || continue
  skill_name=${skill_dir##*/}
  case "$skill_name" in
    cc-entry|cc-next|cc-plan|cc-execute|cc-verify|cc-gates|cc-upgrade) ;;
    *) fail "unexpected legacy skill leaked into artifact: $skill_name" ;;
  esac
done
test ! -e "$artifact/Opus-4.8-plan.md" || fail 'unlisted root plan leaked into artifact'
test ! -e "$artifact/Sol-5.6-plan.md" || fail 'unlisted root plan leaked into artifact'
not_contains "$artifact/workspace.yaml" 'credential'

build_result=$(sh "$ROOT/scripts/build-dist.sh" preview "$build_out")
build_artifact="$build_out/context-circuit-preview"
require_file "$build_artifact/README.md"
require_file "$build_out/context-circuit-preview.tar.gz"
test ! -e "$build_artifact/template" || fail 'source template directory leaked into dist build'
test ! -e "$build_artifact/sources/context-circuit-design" || fail 'source design material leaked into dist build'
printf '%s\n' "$build_result" | grep -F "dist_dir: $build_out" >/dev/null || fail 'dist build did not report output directory'
pass 'staged artifact identity, blank seed, exclusion boundary, and rollback-safe source state'
