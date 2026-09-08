#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

stage=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-stage.XXXXXX")
out=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-out.XXXXXX")
build_out=$(mktemp -d "${TMPDIR:-/tmp}/cc-release-build.XXXXXX")
trap 'rm -rf "$stage" "$out" "$build_out"' EXIT HUP INT TERM

result=$(sh "$ROOT/scripts/release-artifact.sh" "$stage" "$out" v0.5.0)
artifact="$out/context-circuit-v0.5.0"

# --- artifact identity and required files ---
require_file "$artifact/AGENTS.md"
require_file "$artifact/CLAUDE.md"
require_file "$artifact/CURSOR.md"
require_file "$artifact/WORKFLOW.md"
require_file "$artifact/README.md"
require_file "$artifact/.context-circuit/wrapper/runtime/worker-brief.md"
require_file "$artifact/.gitignore"
require_file "$artifact/workspace.yaml"
require_file "$artifact/members.yaml"
require_file "$artifact/.context-circuit/wrapper/manifest.yaml"
require_file "$artifact/.context-circuit/wrapper/runtime/engine.sh"
require_file "$artifact/.context-circuit/wrapper/migrations/README.md"
require_file "$artifact/.context-circuit/wrapper/contracts/invariants.yaml"
for s in workspace repositories-local intent-contract trace-manifest plan task execution worker-handoff \
  verifier-result candidate human-acceptance completion context-impact context-index \
  lease grounding-manifest pairing-session publication-config publication-field-intent \
  publication-record \
  publication-thread-record members member-local; do
  require_file "$artifact/.context-circuit/wrapper/contracts/schemas/$s.yaml"
done
require_file "$artifact/.context-circuit/docs/getting-started.md"
require_file "$artifact/.context-circuit/docs/templates/plan.yaml"
require_file "$artifact/context/PROJECT.md"
require_file "$artifact/sources/README.md"
require_file "$artifact/intent/README.md"
require_file "$artifact/intent/INDEX.md"
require_file "$artifact/intent/archive/README.md"
require_file "$artifact/plans/README.md"
require_file "$artifact/plans/INDEX.md"
require_file "$artifact/.context-circuit/agents/planner.md"
require_file "$artifact/.context-circuit/agents/planner-brief.md"
test ! -e "$artifact/.context-circuit/agents/tracer.md" || fail "tracer role leaked into artifact"
test ! -e "$artifact/.context-circuit/agents/spec-adversary.md" || fail "spec-adversary role leaked into artifact"
test ! -e "$artifact/wrapper" || fail "root wrapper leaked into artifact"
test ! -e "$artifact/agents" || fail "root agents leaked into artifact"
test ! -d "$artifact/docs" || fail "root docs leaked into artifact"
require_dir "$artifact/.agents/skills"
require_file "$artifact/.agents/skills/cc-workspace/SKILL.md"
require_file "$artifact/.claude/agents/worker.md"
require_file "$artifact/.claude/agents/verifier.md"
require_file "$artifact/.claude/agents/planner.md"
require_file "$artifact/.claude/rules/role-tiering-spawn.md"
require_file "$artifact/.claude/rules/commit-convention.md"
require_file "$artifact/.codex/agents/worker.toml"
require_file "$artifact/.codex/agents/verifier.toml"
require_file "$artifact/.codex/agents/planner.toml"
require_file "$artifact/.cursor/agents/worker.md"
require_file "$artifact/.cursor/agents/verifier.md"
require_file "$artifact/.cursor/agents/planner.md"
require_file "$artifact/.cursor/rules/role-tiering-spawn.mdc"
require_file "$artifact/.cursor/rules/commit-convention.mdc"
test ! -e "$artifact/.claude/agents/cc-human-simulator.md" || fail "source-only Claude agent leaked into artifact"
test ! -e "$artifact/.claude/skills/cc-test-case" || fail "source-only Claude skill leaked into artifact"
test ! -e "$artifact/.codex/rules" || fail "non-native .codex/rules leaked into artifact"
for s in cc-workspace cc-intent cc-trace cc-plan cc-execute cc-run-stack cc-system-design cc-verify cc-complete cc-archive cc-deliver cc-pair cc-publish; do
  require_file "$artifact/.claude/skills/$s/SKILL.md"
done

# --- uninitialized identity and thin adapter ---
contains "$artifact/workspace.yaml" 'workspace: uninitialized-workspace'
contains "$artifact/CLAUDE.md" '@AGENTS.md'
contains "$artifact/CURSOR.md" '@AGENTS.md'
not_contains "$artifact/CLAUDE.md" 'transcript:'
not_contains "$artifact/CURSOR.md" 'transcript:'

# --- porcelain source state recorded for rollback safety ---
if git -C "$ROOT" status --porcelain --untracked-files=all | grep . >/dev/null 2>&1; then
  expected_source_state=dirty
else
  expected_source_state=clean
fi
printf '%s\n' "$result" | grep -F "source_state: $expected_source_state" >/dev/null || fail "release did not record source state: $expected_source_state"
printf '%s\n' "$result" | grep -F "runtime_version: 1.0.0" >/dev/null || fail 'release did not report manifest runtime version'

# --- exclusion boundary: no maintainer, source, test, or runtime state ---
for leaked in .runtime test .github scripts template repositories repositories.local.yaml \
  member.local.yaml \
  .context-circuit/wrapper/adapters plans/context-circuit-plans plans/0002-mark-done-no-precheck \
  plans/archive/0001-contracts-runtime-foundation sources/system-design sources/reports \
  .context-circuit/docs/release.md .context-circuit/wrapper/contracts/routes.yaml .context-circuit/wrapper/contracts/context-sets.yaml \
  .context-circuit/wrapper/contracts/schemas/delegation.yaml .agents/skills/cc-entry .code-review-graph; do
  test ! -e "$artifact/$leaked" || fail "leaked into artifact: $leaked"
done
contains "$artifact/.gitignore" ".code-review-graph"

# --- only the shipped skills ship ---
count=0
for skill_dir in "$artifact"/.agents/skills/cc-*; do
  [ -d "$skill_dir" ] || continue
  name=${skill_dir##*/}
  case "$name" in
    cc-workspace|cc-intent|cc-trace|cc-plan|cc-execute|cc-run-stack|cc-system-design|cc-verify|cc-complete|cc-archive|cc-deliver|cc-pair|cc-publish) count=$((count + 1)) ;;
    *) fail "unexpected skill leaked into artifact: $name" ;;
  esac
done
test "$count" -eq 13 || fail "shipped skill count: $count"

# --- no credentials or provider payloads anywhere in the artifact ---
if grep -REn '^[[:space:]]*(password|api_key|access_token|provider_payload|transcript):' "$artifact" >/dev/null 2>&1; then
  fail 'release artifact contained credentials or provider payloads'
fi

# --- build-dist wrapper produces a versioned artifact + archive ---
build_result=$(sh "$ROOT/scripts/build-dist.sh" v0.5.0 "$build_out")
require_file "$build_out/context-circuit-v0.5.0/README.md"
require_file "$build_out/context-circuit-v0.5.0.tar.gz"
printf '%s\n' "$build_result" | grep -F "dist_dir: $build_out" >/dev/null || fail 'dist build did not report output directory'

pass 'release artifact assembly and exclusion boundary'
