#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

stage=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-stage.XXXXXX")
out=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-out.XXXXXX")
trap 'rm -rf "$stage" "$out"' EXIT HUP INT TERM
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
test ! -e "$artifact/sources/secret.txt" || fail 'source inbox leaked into artifact'
if find "$artifact/.agents/skills" -mindepth 1 -maxdepth 1 -type d -name 'cc-*' -print -quit | grep .; then
  fail 'legacy cc skill leaked into artifact'
fi
not_contains "$artifact/workspace.yaml" 'credential'
pass 'staged artifact identity, blank seed, exclusion boundary, and rollback-safe source state'
