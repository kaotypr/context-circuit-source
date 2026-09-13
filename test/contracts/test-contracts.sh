#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

W="$ROOT/.context-circuit/wrapper"

# --- invariants owner map: v0.5 rules and owners present ---
inv="${W}/contracts/invariants.yaml"
require_file "$inv"
for id in INV-INTENT-01 INV-INTENT-02 INV-MEMBER-01 INV-CANDIDATE-01 INV-ASSURE-01 INV-PLAN-01 INV-PLAN-03 INV-PLAN-05 INV-APPROVE-01 INV-EXEC-01 \
	INV-VERIFY-01 INV-VERIFY-02 \
	INV-REPAIR-01 INV-COMPLETE-01 INV-ARCHIVE-01 INV-REPO-02 INV-DELIVER-01 \
	INV-RUNTIME-01 INV-KNOWLEDGE-02 INV-KNOWLEDGE-03 INV-OWN-01 INV-CONCURRENCY-01 INV-CONCURRENCY-02 \
	INV-GROUND-01 INV-GROUND-02 INV-GROUND-03 INV-EXTERNAL-01 INV-EXTERNAL-02 \
	INV-EXTERNAL-03 INV-PAIR-01; do
	contains "$inv" "$id"
done
for concern in intent_contract intent_gate feasibility_check planner_role trace_manifest \
	candidate_identity human_acceptance assurance_tiering \
	plan_lifecycle runtime repository_identity local_binding \
	member_roster member_identity \
	execution_records verifier_result completion_record \
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
contains "$inv" "plans in different repositories are delivered"
contains "${W}/runtime/engine.sh" "cc_change_set_require_same_repo"
contains "${W}/runtime/engine.sh" "cc_change_set_partition"
contains "${W}/runtime/engine.sh" "CHANGE_SET_CROSS_REPO"
contains "${W}/runtime/engine.sh" "CHANGE_SET_NO_SINGLE_TIP"
contains "${W}/runtime/engine.sh" "CHANGE_SET_DELIVERY_HAS_NO_VERIFIER"
contains "$inv" "does not spawn a verifier"
contains "$ROOT/product/.agents/skills/cc-deliver/SKILL.md" "CHANGE_SET_CROSS_REPO"
contains "$ROOT/product/.agents/skills/cc-deliver/SKILL.md" "change-set-partition"
contains "$ROOT/product/.agents/skills/cc-deliver/SKILL.md" "CHANGE_SET_DELIVERY_HAS_NO_VERIFIER"
contains "$ROOT/product/.agents/skills/cc-run-stack/SKILL.md" "CHANGE_SET_CROSS_REPO"
contains "$ROOT/product/.agents/skills/cc-run-stack/SKILL.md" "change-set-partition"
contains "$ROOT/.context-circuit/agents/coordinator.md" "one pull request per covering tip"
contains "$ROOT/.context-circuit/agents/coordinator.md" "Delivery does not spawn a verifier"

# --- manifest declares runtime exclusions and the release boundary ---
man="${W}/manifest.yaml"
require_file "$man"
contains "$man" "runtime_version: 1.0.0"
contains "$man" "plan: [3]"
contains "$man" "repositories-local: [1, 2]"
contains "$man" "execution: [1, 2]"
contains "$man" "intent-contract: [1, 2]"
contains "$man" "trace-manifest: [1]"
contains "$man" "candidate: [1]"
contains "$man" "human-acceptance: [1]"
contains "$man" "lease: [1]"
contains "$man" "grounding-manifest: [1]"
contains "$man" "host-neutral-deterministic-library"
contains "$man" "automatic plan completion"
contains "$man" "provider-specific child-agent launch"
contains "$man" "pairing-session: [1]"
contains "$man" "members: [1]"
contains "$man" "member-local: [1]"
contains "${W}/contracts/schemas/repositories-local.yaml" "base_branch"
contains "${W}/contracts/schemas/repositories-local.yaml" "anchor_branch"
contains "${W}/contracts/schemas/execution.yaml" "base_branch"
contains "${W}/contracts/schemas/workspace.yaml" "values: [1]"
contains "${W}/contracts/schemas/verifier-result.yaml" "schema_version: { type: integer, required: true, values: [1] }"
contains "${W}/contracts/schemas/context-impact.yaml" "schema_version: { type: integer, required: true, values: [1] }"
assert_eq "schema_version: 2" "$(sed -n '1p' "${W}/contracts/schemas/repositories-local.yaml")"
assert_eq "schema_version: 2" "$(sed -n '1p' "${W}/contracts/schemas/execution.yaml")"
assert_eq "schema_version: 3" "$(sed -n '1p' "${W}/contracts/schemas/plan.yaml")"

# --- shipped schemas present; old-design schemas absent ---
for s in workspace repositories-local intent-contract plan task execution worker-handoff \
	verifier-result candidate human-acceptance completion context-impact context-index lease \
	grounding-manifest pairing-session publication-config publication-field-intent \
	publication-record \
	publication-thread-record members member-local; do
	require_file "${W}/contracts/schemas/$s.yaml"
done
for old in delegation session stack child-start context-receipt handoff archive; do
	test ! -e "${W}/contracts/schemas/$old.yaml" || fail "old schema remains: $old"
done
for oldc in routes context-sets tier0; do
	test ! -e "${W}/contracts/$oldc.yaml" || fail "old contract remains: $oldc"
done

# --- grounding environment is preparation, never lockfile detection alone ---
grounding_schema="${W}/contracts/schemas/grounding-manifest.yaml"
contains "$grounding_schema" "values: [ready, no-toolchain]"
contains "$grounding_schema" "gitignored overlay"
contains "$grounding_schema" "match its own lockfile"
contains "$grounding_schema" "not ready"
contains "$grounding_schema" "blocks setup"
contains "${W}/runtime/engine.sh" "Ready: gitignored overlay complete"
for f in "${W}/adapters/worker-brief.md" "${W}/runtime/engine.sh"; do
	not_contains "$f" "no-verify"
	not_contains "$f" "verify-deps"
	not_contains "$f" "shared node_modules"
done

# --- runtime engine present and host-neutral ---
require_file "${W}/runtime/engine.sh"
not_contains "${W}/runtime/engine.sh" "cc_probe"
not_contains "${W}/runtime/engine.sh" "cc_route"
not_contains "${W}/runtime/engine.sh" "cc_confirmation_card"
# persisted record formats have independent schema-version owners; there is no
# generic runtime schema version that can be mistaken for a record schema.
not_contains "${W}/runtime/engine.sh" "CC_SCHEMA_VERSION"
for record_schema in \
	CC_REPOSITORIES_LOCAL_SCHEMA_VERSION CC_PAIRING_SESSION_SCHEMA_VERSION \
	CC_MEMBERS_SCHEMA_VERSION CC_MEMBER_LOCAL_SCHEMA_VERSION \
	CC_GROUNDING_MANIFEST_SCHEMA_VERSION CC_LEASE_SCHEMA_VERSION \
	CC_EXECUTION_SCHEMA_VERSION CC_CANDIDATE_SCHEMA_VERSION \
	CC_CHANGE_SET_SCHEMA_VERSION CC_VERIFIER_RESULT_SCHEMA_VERSION \
	CC_HUMAN_ACCEPTANCE_SCHEMA_VERSION CC_COMPLETION_SCHEMA_VERSION \
	CC_DELIVERY_SCHEMA_VERSION CC_KNOWLEDGE_DEBT_SCHEMA_VERSION; do
	contains "${W}/runtime/engine.sh" "$record_schema"
done

# --- shipped skills present; old-design skills absent ---
for sk in cc-workspace cc-intent cc-trace cc-plan cc-execute cc-run-stack cc-system-design cc-verify cc-complete cc-archive cc-deliver cc-pair cc-publish; do
	require_file "$ROOT/product/.agents/skills/$sk/SKILL.md"
done
# --- v1.0 intent front door: skill, schema, role, and gate wording present ---
ci="$ROOT/product/.agents/skills/cc-intent/SKILL.md"
contains "$ci" "planner"
contains "$ci" "intent-approve"
contains "$ci" "contract_digest"
contains "$ci" "feasibility"
contains "$ci" "phase-aware"
contains "$ci" "intent-level"
contains "$ci" "plan-level"
# --- tracing/feasibility replaced the spec adversary + envelope check ---
not_contains "$ci" "spec adversary"
not_contains "$ci" "scope envelope"
test ! -e "$ROOT/.context-circuit/agents/spec-adversary.md" || fail "spec-adversary role must be removed"
require_file "$ROOT/.context-circuit/agents/planner.md"
contains "$ROOT/.context-circuit/agents/planner.md" "real code"
contains "$ROOT/.context-circuit/agents/planner.md" "child per repository"
contains "$ROOT/.context-circuit/agents/planner.md" "Do not record already-answered"
contains "$ROOT/.context-circuit/agents/planner.md" "more than"
require_file "$ROOT/.context-circuit/agents/planner-brief.md"
contains "$ROOT/.context-circuit/agents/planner-brief.md" "one or more plans"
contains "$ROOT/.context-circuit/agents/planner-brief.md" "Do not record already-answered"
test ! -e "$ROOT/.context-circuit/agents/tracer.md" || fail "tracer role must be removed"
require_file "${W}/contracts/schemas/trace-manifest.yaml"
contains "${W}/contracts/schemas/intent-contract.yaml" "i<NNN>-<kebab-slug>"
tr="$ROOT/product/.agents/skills/cc-trace/SKILL.md"
contains "$tr" "one planner child per repository"
contains "$tr" "feasibility check"
contains "$tr" "Classify every question before feasibility"
contains "$tr" "already-answered"
contains "$tr" "done"
contains "$tr" "Do not add extra classification homework"
ce="$ROOT/product/.agents/skills/cc-plan/SKILL.md"
contains "$ce" "intent-authorized"
contains "$ce" "unresolved intent-level question"
not_contains "$ce" "intent-envelope-check"
# the envelope verb is gone from the runtime; the authorization verb replaces it
not_contains "${W}/runtime/engine.sh" "intent-envelope-check"
not_contains "${W}/runtime/engine.sh" "cc_intent_envelope_check"
contains "${W}/runtime/engine.sh" "cc_intent_authorized"
contains "${W}/runtime/engine.sh" "cc_attempt_norm"
cv="$ROOT/product/.agents/skills/cc-verify/SKILL.md"
contains "$cv" "current_attempt"
contains "$cv" "003"
contains "$ROOT/product/.agents/skills/cc-execute/SKILL.md" "003"

# --- direct collaboration is outside the plan lifecycle ---
pair="$ROOT/product/.agents/skills/cc-pair/SKILL.md"
contains "$pair" "outside the plan lifecycle"
contains "$pair" "Do not launch a verifier"
contains "$pair" "human-supervised"
contains "$pair" "pair-begin"
contains "$pair" "pair-close"
contains "$pair" "repository-relative path"
contains "$pair" "Model & effort per role"
contains "$pair" "role-tiering"
contains "$pair" "workspace root"
contains "$pair" "isolated working copy"
contains "$pair" "runtime/explore"
contains "$pair" "runtime-cleanup"
# v1.0: pairing is the Explore tier with an explicit promote step (INV-ASSURE-01)
contains "$pair" "Explore tier"
contains "$pair" "Promote"
contains "$inv" "INV-ASSURE-01"
contains "${W}/adapters/WORKFLOW.md" "/cc-pair"
contains "$ROOT/.context-circuit/agents/coordinator.md" "INV-PAIR-01"
contains "$ROOT/.context-circuit/agents/worker.md" "Direct-collaboration mode"
contains "${W}/adapters/AGENTS.md" "reads and follows \`.context-circuit/agents/coordinator.md\`"
# Source-checkout host files import workspace-root AGENTS.md. Pointing
# CLAUDE.md/CURSOR.md at nested adapters would replace maintainer AGENTS.md
# with the shipped product copy.
contains "$ROOT/AGENTS.md" "context-circuit-source repository safety"
contains "$ROOT/CLAUDE.md" "@AGENTS.md"
contains "$ROOT/CURSOR.md" "@AGENTS.md"
not_contains "$ROOT/CLAUDE.md" "@.context-circuit/wrapper/adapters/CLAUDE.md"
not_contains "$ROOT/CURSOR.md" "@.context-circuit/wrapper/adapters/CURSOR.md"
contains "$ROOT/.context-circuit/agents/coordinator.md" "Start at the user's vocabulary level"
contains "$ROOT/.context-circuit/docs/terminology.md" '| Workspace | "workspace" when the user has not introduced that term |'
contains "$ROOT/.context-circuit/docs/terminology.md" '| Runtime, skill, tool, or command failure |'
contains "$ROOT/.context-circuit/agents/coordinator.md" "classify every trace question"
contains "$ROOT/.context-circuit/docs/templates/intent.md" "At draft time"
for old in cc-entry cc-gates cc-next cc-upgrade; do
	test ! -e "$ROOT/product/.agents/skills/$old" || fail "old skill remains: $old"
done

# --- cc-system-design carries the authoring rubric and adds no runtime surface ---
sd="$ROOT/product/.agents/skills/cc-system-design/SKILL.md"
contains "$sd" "sources/system-design/"
contains "$sd" "intent/<id>/detail/"
contains "$sd" "three-tier"
contains "$sd" "design.md"
contains "$sd" "by concern"
contains "$sd" "never by repository"
contains "$sd" "Never approve"
# dual-root: product-level home is passive source material; intent detail is a second home, not a new skill
contains "$ROOT/context/domains/system-design-authoring/README.md" "intent/<id>/detail/"
contains "$ROOT/context/domains/system-design-authoring/README.md" "passive source material"
contains "$ROOT/context/domains/system-design-authoring/README.md" "<product-or-project>/<grouping>/<scope>/"
test ! -e "$ROOT/product/.agents/skills/cc-intent-detail" || fail "cc-intent-detail must not exist"
test ! -e "$ROOT/sources/system-design/context-circuit/v1.1.0" || fail "unsolicited v1.1.0 design tree must be absent"
# system-design is authoring-only: no runtime schema, no invariant, no WORKFLOW action
test ! -e "${W}/contracts/schemas/system-design.yaml" || fail "system-design must add no runtime schema"
test ! -e "${W}/contracts/schemas/design-acceptance.yaml" || fail "system-design must add no acceptance schema"
not_contains "$inv" "INV-DESIGN"
not_contains "${W}/adapters/WORKFLOW.md" "design the system"

# --- adapters and roles present; role aliases removed ---
for a in AGENTS.md WORKFLOW.md CLAUDE.md CURSOR.md README.md worker-brief.md; do
	require_file "${W}/adapters/$a"
done
# execution-latency: additive-only, no new invariant id. Per-role tiering
# guidance is coordinator-read, so it ships in .context-circuit/docs/ (a wholesale-shipped tree),
# not in the source-only .context-circuit/wrapper/adapters/ staging directory.
require_file "$ROOT/.context-circuit/docs/role-tiering.md"
contains "$inv" "role_tiering: .context-circuit/docs/role-tiering.md"
contains "$ROOT/.context-circuit/docs/role-tiering.md" "planner"
contains "$ROOT/.context-circuit/docs/role-tiering.md" "workspace root"
contains "$ROOT/.context-circuit/docs/role-tiering.md" "A missing file in an isolated working copy is not an absent config"
contains "$ROOT/product/.agents/skills/cc-execute/SKILL.md" "workspace root"
contains "$ROOT/product/.agents/skills/cc-trace/SKILL.md" "workspace root"
contains "$ROOT/product/.agents/skills/cc-trace/SKILL.md" "role-tiering"
contains "${W}/adapters/AGENTS.md" "workspace root"
contains "${W}/adapters/CLAUDE.md" "workspace root"
contains "${W}/adapters/CURSOR.md" "workspace root"
contains "${W}/adapters/AGENTS.md" "Applying a configured tier on Cursor"
contains "${W}/adapters/CURSOR.md" "Do not default to inherit"
not_contains "$ROOT/.context-circuit/docs/role-tiering.md" "tracer"
not_contains "$ROOT/product/.agents/skills/cc-intent/SKILL.md" "tracer"
not_contains "$ROOT/product/.agents/skills/cc-trace/SKILL.md" "tracer"
not_contains "${W}/adapters/AGENTS.md" "tracer"
not_contains "${W}/adapters/CLAUDE.md" "tracer"
not_contains "${W}/adapters/CURSOR.md" "tracer"
not_contains "$ROOT/CLAUDE.md" "tracer"
not_contains "$ROOT/CURSOR.md" "tracer"
not_contains "$ROOT/.context-circuit/agents/coordinator.md" "tracer"
not_contains "$ROOT/.context-circuit/agents/planner.md" "tracer"
not_contains "$inv" "tracer"
not_contains "${W}/contracts/schemas/trace-manifest.yaml" "tracer"

# --- engine-verb access to role-tiering config, with a committed fallback ---
contains "${W}/runtime/engine.sh" "role-tiering-read"
contains "${W}/runtime/engine.sh" "ROLE_TIERING_MISSING"
require_file "$ROOT/.context-circuit/role-tiering.fallback.yaml"
contains "$ROOT/.context-circuit/docs/role-tiering.md" "role-tiering-read"
contains "$ROOT/.context-circuit/docs/role-tiering.md" "fallback"
for f in "${W}/adapters/AGENTS.md" "${W}/adapters/CLAUDE.md" "${W}/adapters/CURSOR.md" \
	"$ROOT/.context-circuit/agents/coordinator.md" "$ROOT/.context-circuit/agents/worker.md" \
	"$ROOT/.context-circuit/agents/verifier.md" "$ROOT/.context-circuit/agents/planner.md"; do
	contains "$f" ".context-circuit/docs/role-tiering.md"
	not_contains "$f" "role-tiering.fallback.yaml"
done
contains "${W}/manifest.yaml" ".context-circuit/role-tiering.fallback.yaml"
contains "${W}/manifest.yaml" "role-tiering.local.yaml"
if git -C "$ROOT" status --porcelain -- .gitignore .rtk '*/config.toml' 2>/dev/null | grep . >/dev/null 2>&1; then
	fail "this plan must not modify .gitignore or rtk config"
fi
not_contains "$inv" "INV-LATENCY"
not_contains "$inv" "INV-TIER"
contains "$man" "runtime_version: 1.0.0"
for r in coordinator worker verifier; do
	require_file "$ROOT/.context-circuit/agents/$r.md"
done
for alias in repository-worker reviewer tracer; do
	test ! -e "$ROOT/.context-circuit/agents/$alias.md" || fail "role alias remains: $alias"
done

# --- band-scoped allocation: roster, local identity, in-band never-reuse ---
not_contains "$inv" "until band-scoped allocation"
contains "$inv" "INV-MEMBER-01"
contains "$inv" "current member's intent band"
contains "$inv" "current member's plan band"
contains "$inv" "never asks the user to type a block"
contains "$inv" "cc/<plan-id>/<repository-id>"
contains "${W}/contracts/schemas/intent-contract.yaml" "current member's band"
contains "${W}/contracts/schemas/intent-contract.yaml" "member.local.yaml"
contains "${W}/contracts/schemas/plan.yaml" "current member's band"
contains "${W}/contracts/schemas/plan.yaml" "consecutive ids inside that band"
contains "${W}/contracts/schemas/members.yaml" "member.local.yaml"
contains "${W}/contracts/schemas/members.yaml" "INV-MEMBER-01"
contains "${W}/contracts/schemas/member-local.yaml" "members.yaml"
contains "${W}/contracts/schemas/member-local.yaml" "intent-allocate-id"
contains "${W}/contracts/schemas/execution.yaml" "cc/<plan-id>/<repository-id>"

pass 'contracts'
