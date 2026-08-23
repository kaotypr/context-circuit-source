#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

for path in \
  wrapper/manifest.yaml wrapper/contracts/invariants.yaml wrapper/contracts/routes.yaml \
  wrapper/contracts/context-sets.yaml wrapper/contracts/tier0.yaml \
  wrapper/contracts/schemas/workspace.yaml wrapper/contracts/schemas/session.yaml \
  wrapper/contracts/schemas/delegation.yaml wrapper/contracts/schemas/context-receipt.yaml \
  wrapper/contracts/schemas/child-start.yaml \
  wrapper/contracts/schemas/lease.yaml wrapper/contracts/schemas/handoff.yaml \
  wrapper/contracts/schemas/plan.yaml wrapper/contracts/schemas/task.yaml \
  wrapper/contracts/schemas/archive.yaml wrapper/contracts/schemas/completion.yaml \
  wrapper/contracts/schemas/stack.yaml \
  wrapper/runtime/engine.sh docs/templates/plan.yaml; do
  require_file "$ROOT/$path"
done

contains "$ROOT/wrapper/manifest.yaml" 'wrapper_version: 1.0.0'
contains "$ROOT/wrapper/manifest.yaml" 'artifact_kind: blank-workspace'
contains "$ROOT/wrapper/contracts/routes.yaml" 'owner: wrapper/contracts/routes.yaml'
contains "$ROOT/wrapper/contracts/context-sets.yaml" 'overrun: Report'
contains "$ROOT/wrapper/contracts/context-sets.yaml" 'id: initialization'
contains "$ROOT/wrapper/contracts/context-sets.yaml" 'id: plan-draft'
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
contains "$ROOT/wrapper/contracts/schemas/child-start.yaml" 'evidence_only: child-start.yaml records validated launch evidence and never authorizes execution'
contains "$ROOT/wrapper/contracts/schemas/child-start.yaml" 'commit_marker: engine-generated marker is required for authoritative graph publication'
contains "$ROOT/wrapper/contracts/schemas/child-start.yaml" 'legacy-readable: preserve readable legacy evidence without rewriting or upgrading authority'
contains "$ROOT/wrapper/contracts/schemas/handoff.yaml" 'offline_fallback'
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'provider-payloads'
contains "$ROOT/wrapper/contracts/routes.yaml" 'host_binding:'
contains "$ROOT/wrapper/contracts/routes.yaml" 'unavailable_child: block-missing-child-primitive'
contains "$ROOT/wrapper/contracts/routes.yaml" 'context_set_map:'
contains "$ROOT/wrapper/contracts/schemas/context-receipt.yaml" 'packet_digest'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-HOST-01'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-REC-03'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-REC-04'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-REC-05'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'child_scope_authority: wrapper/contracts/schemas/delegation.yaml'
contains "$ROOT/wrapper/manifest.yaml" 'child_start: [1]'

for record_schema in session context-receipt delegation child-start handoff completion; do
  contains "$ROOT/wrapper/contracts/schemas/$record_schema.yaml" 'constructor_inputs:'
  contains "$ROOT/wrapper/contracts/schemas/$record_schema.yaml" 'engine_generated:'
  contains "$ROOT/wrapper/contracts/schemas/$record_schema.yaml" 'transaction_states: [staged, valid, committed, interrupted, legacy-readable]'
  contains "$ROOT/wrapper/contracts/schemas/$record_schema.yaml" 'compatibility:'
done
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'engine_generated: [schema_version, wrapper_version, delegation_id'
contains "$ROOT/wrapper/contracts/schemas/handoff.yaml" 'engine_generated: [schema_version, wrapper_version, handoff_id'
contains "$ROOT/wrapper/contracts/schemas/completion.yaml" 'engine_generated: [schema_version, wrapper_version, completion_id'
contains "$ROOT/wrapper/contracts/schemas/child-start.yaml" 'engine_generated: [schema_version, wrapper_version, child_start_id'

child_start_fixture=$(mktemp "${TMPDIR:-/tmp}/cc-child-start.XXXXXX")
sed -n '/^required:/,/^record_contract:/p' "$ROOT/wrapper/contracts/schemas/child-start.yaml" >"$child_start_fixture"
contains "$child_start_fixture" 'ownership_graph_digest'
contains "$child_start_fixture" 'transaction_state'
contains "$child_start_fixture" 'commit_marker'
rm -f "$child_start_fixture"
contains "$ROOT/.gitignore" 'repositories/'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'workspace.roles'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'canonical_projection: none-or-comma-separated'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'product-source: exempt'
contains "$ROOT/wrapper/contracts/schemas/workspace.yaml" 'missing_files: projection-mismatch'
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
missing_files=$(cc_validate_identity_projection "$fx/missing-files" || true)
printf '%s\n' "$missing_files" | grep -Fx projection-mismatch >/dev/null || fail 'missing summaries did not yield projection-mismatch'
entry_missing=$(cc_entry_preflight "$fx/missing-files" || true)
printf '%s\n' "$entry_missing" | grep -Fx projection-mismatch >/dev/null || fail 'entry preflight fail-opened when summaries were missing'
printf '%s\n' "$entry_missing" | grep -F entry-preflight-ok >/dev/null && fail 'entry preflight reported ok without summaries'
write_missing=$(cc_write_preflight "$fx/missing-files" || true)
printf '%s\n' "$write_missing" | grep -Fx projection-mismatch >/dev/null || fail 'write preflight fail-opened when summaries were missing'
assert_eq "$(cc_validate_identity_projection "$fx/roles-sequence")" identity-projection-ok
roles_disagree=$(mktemp -d "${TMPDIR:-/tmp}/cc-roles-disagree.XXXXXX")
trap 'rm -rf "$roles_disagree"' EXIT HUP INT TERM
cp -R "$fx/roles-sequence/." "$roles_disagree/"
sed -i 's/roles: owner/roles: none/' "$roles_disagree/context/WORKSPACE.md"
roles_mismatch=$(cc_validate_identity_projection "$roles_disagree" || true)
printf '%s\n' "$roles_mismatch" | grep -Fx projection-mismatch >/dev/null || fail 'sequence roles collapsed to none instead of disagreeing'
undeclared=$(cc_validate_delegated_effects "$fx/undeclared-effect/approved-effects.txt" "$fx/undeclared-effect/delegated-effects.txt" || true)
assert_eq "$undeclared" UNDECLARED_EFFECT
assert_eq "$(cc_validate_delegated_effects "$fx/undeclared-effect/approved-effects.txt" "$fx/undeclared-effect/valid-delegated-effects.txt")" delegated-effects-ok
card=$(cc_identity_acceptance_card "$fx/displayed-default")
printf '%s\n' "$card" | grep -F 'mode: solo' >/dev/null || fail 'displayed-default card omitted mode default'
printf '%s\n' "$card" | grep -F 'roles: none' >/dev/null || fail 'displayed-default card omitted roles default'
printf '%s\n' "$card" | grep -F 'default branches: main' >/dev/null || fail 'displayed-default card omitted branch default'
incomplete=$(cc_register_repository "$fx/incomplete" '' '' '' confirmed 2>&1 || true)
printf '%s\n' "$incomplete" | grep -F INCOMPLETE_FIELD >/dev/null || fail 'missing logical key was not incomplete'

contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'vocabulary: [schema, store, api, process, browser, human]'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'comparison: exact-match-only'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'strength_hierarchy: none'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'required_layer: owned by each acceptance entry'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'produced_layer: owned by each verification mapping'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'values: [passed, failed, blocked, waived]'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'waiver_gate: none'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'completed historical evidence remains readable and is never rewritten or invented'
contains "$ROOT/wrapper/contracts/schemas/task.yaml" 'task_fields: [acceptance, verification]'
contains "$ROOT/wrapper/contracts/schemas/task.yaml" 'tasks retain acceptance and verification ids only'
contains "$ROOT/wrapper/contracts/schemas/delegation.yaml" 'runtime_fields: [observed_layer, outcome, evidence_ref]'
contains "$ROOT/wrapper/contracts/schemas/handoff.yaml" 'runtime_fields: [observed_layer, outcome, evidence_ref]'
contains "$ROOT/wrapper/contracts/schemas/completion.yaml" 'completion_rule: only passed satisfies completion'
contains "$ROOT/wrapper/contracts/schemas/completion.yaml" 'waiver_gate: none'
contains "$ROOT/wrapper/contracts/schemas/completion.yaml" 'migration: must not invent required_layer, produced_layer, observed_layer, or evidence_ref'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-EVID-01'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-EVID-02'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-EVID-03'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-EVID-04'
contains "$ROOT/wrapper/contracts/invariants.yaml" 'evidence_layers: wrapper/contracts/schemas/plan.yaml'
contains "$ROOT/wrapper/manifest.yaml" 'vocabulary: [schema, store, api, process, browser, human]'
contains "$ROOT/wrapper/manifest.yaml" 'comparison: exact-match-only'
contains "$ROOT/wrapper/manifest.yaml" 'completion_rule: only-passed'
contains "$ROOT/wrapper/manifest.yaml" 'strength_hierarchy: none'
contains "$ROOT/wrapper/manifest.yaml" 'waiver_gate: none'
not_contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'weaker_layer_passes'
not_contains "$ROOT/wrapper/contracts/schemas/completion.yaml" 'waived_counts_as_success: true'

fx_evid="$ROOT/test/contracts/fixtures/evidence-layers"
require_file "$fx_evid/matching.yaml"
require_file "$fx_evid/non-matching.yaml"
require_file "$fx_evid/missing.yaml"
require_file "$fx_evid/blocked.yaml"
require_file "$fx_evid/waived.yaml"
require_file "$fx_evid/legacy-completed.yaml"
require_file "$fx_evid/legacy-unfinished.yaml"

yaml_field() {
  awk -F': ' -v field="$2" '
    {
      key=$1
      sub(/^[ \t]+/, "", key)
      if (key == field) { print $2; exit }
    }
  ' "$1"
}

match_req=$(yaml_field "$fx_evid/matching.yaml" required_layer)
match_obs=$(yaml_field "$fx_evid/matching.yaml" observed_layer)
test "$match_req" = "$match_obs" || fail 'matching fixture layers are not exact-match'
test "$(yaml_field "$fx_evid/matching.yaml" outcome)" = passed || fail 'matching fixture must be passed'
test "$(yaml_field "$fx_evid/matching.yaml" satisfies_completion)" = true || fail 'matching passed evidence must satisfy completion'

weak_req=$(yaml_field "$fx_evid/non-matching.yaml" required_layer)
weak_obs=$(yaml_field "$fx_evid/non-matching.yaml" observed_layer)
test "$weak_req" = browser || fail 'non-matching fixture must require browser'
test "$weak_obs" = process || fail 'non-matching fixture must observe process'
test "$weak_req" != "$weak_obs" || fail 'non-matching fixture layers must differ'
test "$(yaml_field "$fx_evid/non-matching.yaml" outcome)" = failed || fail 'non-matching fixture must be failed'
test "$(yaml_field "$fx_evid/non-matching.yaml" satisfies_completion)" = false || fail 'non-matching evidence must not satisfy completion'
test "$(yaml_field "$fx_evid/non-matching.yaml" inferred_equivalence)" = false || fail 'non-matching fixture inferred a substitute layer'

test "$(yaml_field "$fx_evid/missing.yaml" observed_present)" = false || fail 'missing fixture still has observed evidence'
test "$(yaml_field "$fx_evid/missing.yaml" comparison)" = missing || fail 'missing fixture comparison is not missing'
test "$(yaml_field "$fx_evid/missing.yaml" satisfies_completion)" = false || fail 'missing evidence must not satisfy completion'

test "$(yaml_field "$fx_evid/blocked.yaml" outcome)" = blocked || fail 'blocked fixture missing blocked outcome'
test "$(yaml_field "$fx_evid/blocked.yaml" host_capability)" = unavailable || fail 'blocked fixture must record unavailable capability'
test "$(yaml_field "$fx_evid/blocked.yaml" converted_to_pass)" = false || fail 'blocked fixture converted limitation into pass'
test "$(yaml_field "$fx_evid/blocked.yaml" satisfies_completion)" = false || fail 'blocked evidence must not satisfy completion'

test "$(yaml_field "$fx_evid/waived.yaml" outcome)" = waived || fail 'waived fixture missing waived outcome'
test "$(yaml_field "$fx_evid/waived.yaml" waiver_gate)" = none || fail 'waived fixture introduced a waiver gate'
test "$(yaml_field "$fx_evid/waived.yaml" counts_as_success)" = false || fail 'waived fixture counted as success'
test "$(yaml_field "$fx_evid/waived.yaml" satisfies_completion)" = false || fail 'waived evidence must not satisfy completion'

test "$(yaml_field "$fx_evid/legacy-completed.yaml" rewrite)" = false || fail 'legacy-completed fixture rewrote historical evidence'
test "$(yaml_field "$fx_evid/legacy-completed.yaml" invent_evidence)" = false || fail 'legacy-completed fixture invented evidence'
test "$(yaml_field "$fx_evid/legacy-completed.yaml" readable)" = true || fail 'legacy-completed fixture is not readable'
test "$(yaml_field "$fx_evid/legacy-unfinished.yaml" requires_explicit_mapping)" = true || fail 'legacy-unfinished fixture skipped mapping'
test "$(yaml_field "$fx_evid/legacy-unfinished.yaml" new_verification)" = blocked || fail 'legacy-unfinished fixture allowed new verification'
test "$(yaml_field "$fx_evid/legacy-unfinished.yaml" invent_evidence)" = false || fail 'legacy-unfinished fixture invented evidence'

pass 'contract inventory, schema ownership, unique invariant IDs, identity-projection fixtures, and evidence-layer contracts'
