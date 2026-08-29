#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

W="$ROOT/wrapper"

# --- invariants owner map: v0.5 rules and owners present ---
inv="$W/contracts/invariants.yaml"
require_file "$inv"
for id in INV-PLAN-01 INV-PLAN-05 INV-APPROVE-01 INV-EXEC-01 INV-VERIFY-01 INV-VERIFY-02 \
	INV-REPAIR-01 INV-COMPLETE-01 INV-ARCHIVE-01 INV-REPO-02 INV-DELIVER-01 \
	INV-RUNTIME-01 INV-KNOWLEDGE-02 INV-OWN-01 INV-CONCURRENCY-01 INV-CONCURRENCY-02 \
	INV-GROUND-01 INV-GROUND-02 INV-GROUND-03 INV-EXTERNAL-01 INV-EXTERNAL-02 \
	INV-EXTERNAL-03; do
	contains "$inv" "$id"
done
for concern in plan_lifecycle runtime repository_identity local_binding \
	execution_records verifier_result completion_record context_proposals \
	worker_role verifier_role coordinator_role path_leases path_lease_records \
	base_selection run_stack_action repository_grounding grounding_manifest \
	worker_brief external_surface publication_config publication_record \
	publication_thread_record; do
	contains "$inv" "$concern:"
done
# old-design owners are gone
not_contains "$inv" "context_sets:"
not_contains "$inv" "routes:"
# Path A: git delivery no longer says "publish/publication"
not_contains "$inv" "push, publication, deployment"
not_contains "$inv" "push/publish/deploy"

# --- manifest declares runtime exclusions and the release boundary ---
man="$W/manifest.yaml"
require_file "$man"
contains "$man" "runtime_version: 0.6.1"
contains "$man" "plan: [1, 2]"
contains "$man" "host-neutral-deterministic-library"
contains "$man" "automatic plan completion"
contains "$man" "provider-specific child-agent launch"

# --- v0.5 schemas present; old-design schemas absent ---
for s in workspace repositories-local plan task execution worker-handoff \
	verifier-result completion context-impact context-proposal context-index lease \
	grounding-manifest publication-config publication-record \
	publication-thread-record; do
	require_file "$W/contracts/schemas/$s.yaml"
done
for old in delegation session stack child-start context-receipt handoff archive; do
	test ! -e "$W/contracts/schemas/$old.yaml" || fail "old schema remains: $old"
done
for oldc in routes context-sets tier0; do
	test ! -e "$W/contracts/$oldc.yaml" || fail "old contract remains: $oldc"
done

# --- runtime engine present and host-neutral ---
require_file "$W/runtime/engine.sh"
not_contains "$W/runtime/engine.sh" "cc_probe"
not_contains "$W/runtime/engine.sh" "cc_route"
not_contains "$W/runtime/engine.sh" "cc_confirmation_card"

# --- v0.5 skills present; old-design skills absent ---
for sk in cc-workspace cc-plan cc-execute cc-run-stack cc-system-design cc-verify cc-complete cc-archive cc-deliver cc-publish; do
	require_file "$ROOT/.agents/skills/$sk/SKILL.md"
done
for old in cc-entry cc-gates cc-next cc-upgrade; do
	test ! -e "$ROOT/.agents/skills/$old" || fail "old skill remains: $old"
done

# --- cc-system-design carries the authoring rubric and adds no runtime surface ---
sd="$ROOT/.agents/skills/cc-system-design/SKILL.md"
contains "$sd" "sources/system-design/"
contains "$sd" "three-tier"
contains "$sd" "design.md"
contains "$sd" "by concern"
contains "$sd" "never by repository"
contains "$sd" "Never approve"
# system-design is authoring-only: no runtime schema, no invariant, no WORKFLOW action
test ! -e "$W/contracts/schemas/system-design.yaml" || fail "system-design must add no runtime schema"
test ! -e "$W/contracts/schemas/design-acceptance.yaml" || fail "system-design must add no acceptance schema"
not_contains "$inv" "INV-DESIGN"
not_contains "$W/adapters/WORKFLOW.md" "design the system"

# --- adapters and roles present; role aliases removed ---
for a in AGENTS.md WORKFLOW.md CLAUDE.md README.md worker-brief.md; do
	require_file "$W/adapters/$a"
done
for r in coordinator worker verifier; do
	require_file "$ROOT/agents/$r.md"
done
for alias in repository-worker reviewer; do
	test ! -e "$ROOT/agents/$alias.md" || fail "role alias remains: $alias"
done

pass 'contracts'
