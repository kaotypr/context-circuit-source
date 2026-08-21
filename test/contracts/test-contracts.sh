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

ids=$(awk '/^  - id: INV-/{print $3}' "$ROOT/wrapper/contracts/invariants.yaml")
test "$(printf '%s\n' "$ids" | sort | uniq | wc -l)" -eq "$(printf '%s\n' "$ids" | wc -l)" || fail 'duplicate invariant IDs'
not_contains "$ROOT/.agents/skills/cc-entry/SKILL.md" 'cc-session-entry'
not_contains "$ROOT/.agents/skills/cc-execute/SKILL.md" 'cc-run-plan'
pass 'contract inventory, schema ownership, and unique invariant IDs'
