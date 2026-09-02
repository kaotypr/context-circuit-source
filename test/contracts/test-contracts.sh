#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

W="$ROOT/wrapper"

# --- invariants owner map: v0.5 rules and owners present ---
inv="$W/contracts/invariants.yaml"
require_file "$inv"
for id in INV-INTENT-01 INV-INTENT-02 INV-CANDIDATE-01 INV-ASSURE-01 INV-PLAN-01 INV-PLAN-05 INV-APPROVE-01 INV-EXEC-01 \
	INV-VERIFY-01 INV-VERIFY-02 \
	INV-REPAIR-01 INV-COMPLETE-01 INV-ARCHIVE-01 INV-REPO-02 INV-DELIVER-01 \
	INV-RUNTIME-01 INV-KNOWLEDGE-02 INV-OWN-01 INV-CONCURRENCY-01 INV-CONCURRENCY-02 \
	INV-GROUND-01 INV-GROUND-02 INV-GROUND-03 INV-EXTERNAL-01 INV-EXTERNAL-02 \
	INV-EXTERNAL-03 INV-PAIR-01; do
	contains "$inv" "$id"
done
for concern in intent_contract intent_gate scope_envelope spec_adversary_role \
	candidate_identity human_acceptance assurance_tiering reconciliation_debt \
	plan_lifecycle runtime repository_identity local_binding \
	execution_records verifier_result completion_record context_proposals \
	worker_role verifier_role coordinator_role path_leases path_lease_records \
	base_selection run_stack_action repository_grounding grounding_manifest \
	worker_brief external_surface publication_config publication_field_intent \
	publication_record publication_thread_record pairing_mode pairing_session; do
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
contains "$man" "runtime_version: 1.0.0"
contains "$man" "plan: [3]"
contains "$man" "repositories-local: [1, 2]"
contains "$man" "execution: [1, 2]"
contains "$man" "intent-contract: [1]"
contains "$man" "candidate: [1]"
contains "$man" "human-acceptance: [1]"
contains "$man" "host-neutral-deterministic-library"
contains "$man" "automatic plan completion"
contains "$man" "provider-specific child-agent launch"
contains "$man" "pairing-session: [1]"
contains "$W/contracts/schemas/repositories-local.yaml" "base_branch"
contains "$W/contracts/schemas/repositories-local.yaml" "anchor_branch"
contains "$W/contracts/schemas/execution.yaml" "base_branch"

# --- shipped schemas present; old-design schemas absent ---
for s in workspace repositories-local intent-contract plan task execution worker-handoff \
	verifier-result candidate human-acceptance completion context-impact context-proposal context-index lease \
	grounding-manifest pairing-session publication-config publication-field-intent \
	publication-record \
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

# --- shipped skills present; old-design skills absent ---
for sk in cc-workspace cc-intent cc-plan cc-execute cc-run-stack cc-system-design cc-verify cc-complete cc-archive cc-deliver cc-pair cc-publish; do
	require_file "$ROOT/.agents/skills/$sk/SKILL.md"
done
# --- v1.0 intent front door: skill, schema, role, and gate wording present ---
ci="$ROOT/.agents/skills/cc-intent/SKILL.md"
contains "$ci" "spec adversary"
contains "$ci" "intent-approve"
contains "$ci" "contract_digest"
contains "$ci" "envelope"
require_file "$ROOT/agents/spec-adversary.md"
contains "$ROOT/agents/spec-adversary.md" "before any code"
contains "$W/contracts/schemas/intent-contract.yaml" "i<NNNN>-<kebab-slug>"
contains "$ROOT/.agents/skills/cc-plan/SKILL.md" "intent-envelope-check"

# --- direct collaboration is outside the plan lifecycle ---
pair="$ROOT/.agents/skills/cc-pair/SKILL.md"
contains "$pair" "outside the plan lifecycle"
contains "$pair" "Do not launch a verifier"
contains "$pair" "human-supervised"
contains "$pair" "pair-begin"
contains "$pair" "pair-close"
contains "$pair" "repository-relative path"
# v1.0: pairing is the Explore tier with an explicit promote step (INV-ASSURE-01)
contains "$pair" "Explore tier"
contains "$pair" "Promote"
contains "$inv" "INV-ASSURE-01"
contains "$W/adapters/WORKFLOW.md" "/cc-pair"
contains "$ROOT/agents/coordinator.md" "INV-PAIR-01"
contains "$ROOT/agents/worker.md" "Direct-collaboration mode"
contains "$W/adapters/AGENTS.md" "reads and follows \`agents/coordinator.md\`"
contains "$ROOT/agents/coordinator.md" "Start at the user's vocabulary level"
contains "$ROOT/docs/terminology.md" '| Workspace | "workspace" when the user has not introduced that term |'
contains "$ROOT/docs/terminology.md" '| Runtime, skill, tool, or command failure |'
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
# execution-latency: additive-only, no new invariant id. Per-role tiering
# guidance is coordinator-read, so it ships in docs/ (a wholesale-shipped tree),
# not in the source-only wrapper/adapters/ staging directory.
require_file "$ROOT/docs/role-tiering.md"
contains "$inv" "role_tiering: docs/role-tiering.md"
not_contains "$inv" "INV-LATENCY"
not_contains "$inv" "INV-TIER"
contains "$man" "runtime_version: 1.0.0"
for r in coordinator worker verifier; do
	require_file "$ROOT/agents/$r.md"
done
for alias in repository-worker reviewer; do
	test ! -e "$ROOT/agents/$alias.md" || fail "role alias remains: $alias"
done

pass 'contracts'
