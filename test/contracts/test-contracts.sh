#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

for path in \
  wrapper/manifest.yaml wrapper/contracts/invariants.yaml wrapper/contracts/routes.yaml \
  wrapper/contracts/context-sets.yaml wrapper/contracts/tier0.yaml \
  wrapper/contracts/schemas/workspace.yaml wrapper/contracts/schemas/session.yaml \
  wrapper/contracts/schemas/delegation.yaml wrapper/contracts/schemas/context-receipt.yaml \
  wrapper/contracts/schemas/lease.yaml wrapper/contracts/schemas/handoff.yaml \
  wrapper/contracts/schemas/plan.yaml wrapper/contracts/schemas/task.yaml \
  wrapper/contracts/schemas/archive.yaml wrapper/contracts/schemas/completion.yaml \
  wrapper/contracts/schemas/stack.yaml wrapper/runtime/engine.sh; do
  require_file "$ROOT/$path"
done

contains "$ROOT/wrapper/manifest.yaml" 'wrapper_version: 1.0.0'
contains "$ROOT/wrapper/manifest.yaml" 'artifact_kind: blank-workspace'
contains "$ROOT/wrapper/contracts/routes.yaml" 'owner: wrapper/contracts/routes.yaml'
contains "$ROOT/wrapper/contracts/context-sets.yaml" 'overrun: Report'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'status: [draft, approved, done]'
contains "$ROOT/wrapper/contracts/schemas/task.yaml" 'status: [draft, ready, done]'
contains "$ROOT/wrapper/contracts/schemas/task.yaml" 'frontmatter:'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'repositories.<key>.canonical_url'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'repositories.local.<key>.path'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-REPO-04'
contains "$ROOT/wrapper/manifest.yaml" 'repositories.local.yaml'
contains "$ROOT/wrapper/contracts/schemas/session.yaml" 'host_evidence:'
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'host_id'
contains "$ROOT/wrapper/contracts/schemas/handoff.yaml" 'offline_fallback'
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'provider-payloads'
contains "$ROOT/wrapper/contracts/routes.yaml" 'host_binding:'
contains "$ROOT/wrapper/contracts/routes.yaml" 'unavailable_child: block-missing-child-primitive'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-HOST-01'
contains "$ROOT/.gitignore" 'repositories/'

assert_task_frontmatter() {
  task_file=$1
  test "$(sed -n '1p' "$task_file")" = '---' || fail "missing task frontmatter opening: $task_file"
  awk 'NR > 1 && /^---$/ { found=1; exit } END { exit found ? 0 : 1 }' "$task_file" \
    || fail "missing task frontmatter closing: $task_file"
}

assert_task_frontmatter "$ROOT/docs/templates/task.md"
for task_file in "$ROOT"/plans/*-plans/*/tasks/*.md; do
  [ -f "$task_file" ] || continue
  assert_task_frontmatter "$task_file"
done

ids=$(awk '/^  - id: INV-/{print $3}' "$ROOT/wrapper/contracts/invariants.yaml")
test "$(printf '%s\n' "$ids" | sort | uniq | wc -l)" -eq "$(printf '%s\n' "$ids" | wc -l)" || fail 'duplicate invariant IDs'
not_contains "$ROOT/.agents/skills/cc-entry/SKILL.md" 'cc-session-entry'
not_contains "$ROOT/.agents/skills/cc-execute/SKILL.md" 'cc-run-plan'
pass 'contract inventory, schema ownership, and unique invariant IDs'
