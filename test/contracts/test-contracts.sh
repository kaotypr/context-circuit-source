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
contains "$ROOT/docs/getting-started.md" 'create the ignored root'
contains "$ROOT/wrapper/contracts/schemas/session.yaml" 'host_evidence:'
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'host_id'
contains "$ROOT/wrapper/contracts/schemas/handoff.yaml" 'offline_fallback'
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'provider-payloads'
contains "$ROOT/wrapper/contracts/routes.yaml" 'host_binding:'
contains "$ROOT/wrapper/contracts/routes.yaml" 'unavailable_child: block-missing-child-primitive'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-HOST-01'
contains "$ROOT/.gitignore" 'repositories/'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'workspace.roles'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'identity_region:'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'missing_or_disagreeing_result: projection-mismatch'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'proposed_default: solo'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'deriving_product_knowledge_from_workspace.yaml'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-PROJ-01'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-PROJ-02'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-GATE-07'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-GATE-09'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'identity_projection: wrapper/contracts/schemas/workspace.yaml'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'gate_effects: wrapper/contracts/routes.yaml'
contains "$ROOT/wrapper/contracts/routes.yaml" 'workspace.accept_identity'
contains "$ROOT/wrapper/contracts/routes.yaml" 'workspace.register_repository'
contains "$ROOT/wrapper/contracts/routes.yaml" 'repository-create-empty'
contains "$ROOT/wrapper/contracts/routes.yaml" 'git.commit'
contains "$ROOT/wrapper/contracts/routes.yaml" 'delivery.push'
contains "$ROOT/wrapper/contracts/routes.yaml" 'authorization: never'
contains "$ROOT/wrapper/contracts/routes.yaml" 'not_execute_plan_effects: [git.init, git.clone]'
contains "$ROOT/wrapper/manifest.yaml" 'mismatch_result: projection-mismatch'
contains "$ROOT/wrapper/manifest.yaml" 'gate_effects:'
not_contains "$ROOT/.agents/skills/cc-entry/SKILL.md" 'workspace.accept_identity'
not_contains "$ROOT/.agents/skills/cc-execute/SKILL.md" 'projection-mismatch'
not_contains "$ROOT/.agents/skills/cc-gates/SKILL.md" 'workspace.register_repository'
not_contains "$ROOT/agents/writer.md" 'workspace.accept_identity'
not_contains "$ROOT/wrapper/adapters/AGENTS.md" 'Immediate effects: workspace.accept_identity'

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

fx="$ROOT/test/contracts/fixtures/identity-projection"
require_file "$fx/matching/workspace.yaml"
require_file "$fx/mismatched/context/WORKSPACE.md"
require_file "$fx/missing-region/context/INDEX.md"
require_file "$fx/displayed-default/workspace.yaml"
require_file "$fx/incomplete/workspace.yaml"
require_file "$fx/undeclared-effect/delegated-effects.txt"
assert_eq "$(cc_validate_identity_projection "$fx/matching")" identity-projection-ok
mismatch=$(cc_validate_identity_projection "$fx/mismatched" || true)
printf '%s\n' "$mismatch" | grep -Fx projection-mismatch >/dev/null || fail 'mismatched region did not yield projection-mismatch'
missing=$(cc_validate_identity_projection "$fx/missing-region" || true)
printf '%s\n' "$missing" | grep -Fx projection-mismatch >/dev/null || fail 'missing region did not yield projection-mismatch'
undeclared=$(cc_validate_delegated_effects "$fx/undeclared-effect/approved-effects.txt" "$fx/undeclared-effect/delegated-effects.txt" || true)
assert_eq "$undeclared" UNDECLARED_EFFECT
assert_eq "$(cc_validate_delegated_effects "$fx/undeclared-effect/approved-effects.txt" "$fx/undeclared-effect/valid-delegated-effects.txt")" delegated-effects-ok
card=$(cc_identity_acceptance_card "$fx/displayed-default")
printf '%s\n' "$card" | grep -F 'mode: solo' >/dev/null || fail 'displayed-default card omitted mode default'
printf '%s\n' "$card" | grep -F 'roles: none' >/dev/null || fail 'displayed-default card omitted roles default'
printf '%s\n' "$card" | grep -F 'default branches: main' >/dev/null || fail 'displayed-default card omitted branch default'
incomplete=$(cc_register_repository "$fx/incomplete" '' '' '' confirmed 2>&1 || true)
printf '%s\n' "$incomplete" | grep -F INCOMPLETE_FIELD >/dev/null || fail 'missing logical key was not incomplete'
pass 'contract inventory, schema ownership, unique invariant IDs, and identity-projection fixtures'
