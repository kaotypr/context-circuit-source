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
require_file "$artifact/WORKFLOW.md"
require_file "$artifact/README.md"
require_file "$artifact/writer-brief.md"
require_file "$artifact/.gitignore"
require_file "$artifact/workspace.yaml"
require_file "$artifact/wrapper/manifest.yaml"
require_file "$artifact/wrapper/runtime/engine.sh"
require_file "$artifact/wrapper/migrations/README.md"
require_file "$artifact/wrapper/contracts/invariants.yaml"
for s in workspace repositories-local plan task execution worker-handoff \
  verifier-result completion context-impact context-proposal context-index \
  lease grounding-manifest publication-config publication-record; do
  require_file "$artifact/wrapper/contracts/schemas/$s.yaml"
done
require_file "$artifact/docs/getting-started.md"
require_file "$artifact/docs/templates/plan.yaml"
require_file "$artifact/context/PROJECT.md"
require_file "$artifact/context/proposals/README.md"
require_file "$artifact/sources/README.md"
require_file "$artifact/plans/README.md"
require_file "$artifact/plans/INDEX.md"

# --- uninitialized identity and thin adapter ---
contains "$artifact/workspace.yaml" 'workspace: uninitialized-workspace'
contains "$artifact/CLAUDE.md" '@AGENTS.md'
not_contains "$artifact/CLAUDE.md" 'transcript:'

# --- porcelain source state recorded for rollback safety ---
if git -C "$ROOT" status --porcelain --untracked-files=all | grep . >/dev/null 2>&1; then
  expected_source_state=dirty
else
  expected_source_state=clean
fi
printf '%s\n' "$result" | grep -F "source_state: $expected_source_state" >/dev/null || fail "release did not record source state: $expected_source_state"

# --- exclusion boundary: no maintainer, source, test, or runtime state ---
for leaked in .runtime test .github scripts template repositories repositories.local.yaml \
  wrapper/adapters plans/context-circuit-plans sources/system-design sources/reports \
  docs/release.md wrapper/contracts/routes.yaml wrapper/contracts/context-sets.yaml \
  wrapper/contracts/schemas/delegation.yaml .agents/skills/cc-entry; do
  test ! -e "$artifact/$leaked" || fail "leaked into artifact: $leaked"
done

# --- only the shipped skills ship (v0.5 seven + v0.6 cc-run-stack, cc-system-design) ---
count=0
for skill_dir in "$artifact"/.agents/skills/cc-*; do
  [ -d "$skill_dir" ] || continue
  name=${skill_dir##*/}
  case "$name" in
    cc-workspace|cc-plan|cc-execute|cc-run-stack|cc-system-design|cc-verify|cc-complete|cc-archive|cc-deliver|cc-publish) count=$((count + 1)) ;;
    *) fail "unexpected skill leaked into artifact: $name" ;;
  esac
done
test "$count" -eq 10 || fail "shipped skill count: $count"

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
