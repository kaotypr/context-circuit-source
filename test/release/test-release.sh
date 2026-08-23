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
require_file "$artifact/CLAUDE.md"
require_file "$artifact/WORKFLOW.md"
require_file "$artifact/README.md"
require_file "$artifact/workspace.yaml"
require_file "$artifact/wrapper/manifest.yaml"
require_file "$artifact/wrapper/runtime/engine.sh"
require_file "$artifact/wrapper/migrations/upgrade.sh"
require_file "$artifact/wrapper/contracts/schemas/context-receipt.yaml"
require_file "$artifact/wrapper/contracts/schemas/delegation.yaml"
require_file "$artifact/wrapper/contracts/schemas/plan.yaml"
require_file "$artifact/wrapper/contracts/schemas/task.yaml"
require_file "$artifact/docs/gates.md"
require_file "$artifact/docs/templates/plan.md"
require_file "$artifact/docs/templates/plan.yaml"
require_file "$artifact/docs/templates/task.md"
require_file "$artifact/docs/templates/prd.md"
require_file "$artifact/context/PROJECT.md"
require_file "$artifact/sources/README.md"
require_file "$artifact/plans/README.md"
contains "$artifact/workspace.yaml" 'name: uninitialized-workspace'
contains "$artifact/AGENTS.md" 'small safety spine'
contains "$artifact/CLAUDE.md" '@AGENTS.md'
not_contains "$artifact/CLAUDE.md" '.mcp.json'
not_contains "$artifact/CLAUDE.md" 'transcript:'
if git -C "$ROOT" status --porcelain --untracked-files=all | grep . >/dev/null 2>&1; then
  expected_source_state=dirty
else
  expected_source_state=clean
fi
printf '%s\n' "$result" | grep -F "source_state: $expected_source_state" >/dev/null || fail "release did not record porcelain source state: $expected_source_state"
test ! -e "$artifact/.runtime" || fail 'runtime leaked into artifact'
test ! -e "$artifact/test" || fail 'semantic tests leaked into artifact'
test ! -e "$artifact/.cursor" || fail 'host-local Cursor state leaked into artifact'
test ! -e "$artifact/PLAN.md" || fail 'maintainer plan leaked into artifact'
test ! -e "$artifact/sources/context-circuit-design" || fail 'source design material leaked into artifact'
test ! -e "$artifact/template" || fail 'source template directory leaked into artifact'
test ! -e "$artifact/plans/context-circuit-plans" || fail 'maintainer plan stack leaked into artifact'
test ! -e "$artifact/repositories.local.yaml" || fail 'local binding leaked into artifact'
test ! -e "$artifact/repositories" || fail 'repository checkout leaked into artifact'
test ! -e "$artifact/sources/secret.txt" || fail 'source inbox leaked into artifact'
shipped_skill_count=0
for skill_dir in "$artifact"/.agents/skills/cc-*; do
  [ -d "$skill_dir" ] || continue
  skill_name=${skill_dir##*/}
  case "$skill_name" in
    cc-entry|cc-next|cc-plan|cc-execute|cc-verify|cc-gates|cc-upgrade) shipped_skill_count=$((shipped_skill_count + 1)) ;;
    *) fail "unexpected legacy skill leaked into artifact: $skill_name" ;;
  esac
done
test "$shipped_skill_count" -eq 7 || fail "shipped skill allowlist count: $shipped_skill_count"
test ! -e "$artifact/.agents/skills/cc-review-plan" || fail 'cc-review-plan leaked into artifact as a separate skill'
test ! -e "$ROOT/.agents/skills/cc-review-plan" || fail 'cc-review-plan exists in the source checkout'
test ! -e "$artifact/Opus-4.8-plan.md" || fail 'unlisted root plan leaked into artifact'
test ! -e "$artifact/Sol-5.6-plan.md" || fail 'unlisted root plan leaked into artifact'
not_contains "$artifact/workspace.yaml" 'credential'
contains "$artifact/context/WORKSPACE.md" 'context-circuit:identity-region:start'
contains "$artifact/context/PROJECT.md" 'context-circuit:identity-region:end'
contains "$artifact/context/INDEX.md" 'name: uninitialized-workspace'
contains "$artifact/wrapper/contracts/schemas/workspace.yaml" 'missing_or_disagreeing_result: projection-mismatch'
contains "$artifact/wrapper/contracts/routes.yaml" 'authorization: never'
contains "$artifact/wrapper/manifest.yaml" 'identity_projection:'
not_contains "$artifact/workspace.yaml" 'path: /'
not_contains "$artifact/context/WORKSPACE.md" 'repositories.local.yaml'

build_result=$(sh "$ROOT/scripts/build-dist.sh" v0.5.0 "$build_out")
build_artifact="$build_out/context-circuit-v0.5.0"
require_file "$build_artifact/README.md"
require_file "$build_out/context-circuit-v0.5.0.tar.gz"
test ! -e "$build_artifact/template" || fail 'source template directory leaked into dist build'
test ! -e "$build_artifact/sources/context-circuit-design" || fail 'source design material leaked into dist build'
printf '%s\n' "$build_result" | grep -F "dist_dir: $build_out" >/dev/null || fail 'dist build did not report output directory'
pass 'staged artifact identity, blank seed, exclusion boundary, and rollback-safe source state'
