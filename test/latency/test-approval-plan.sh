#!/bin/sh
# Approval-to-plan latency semantics are deterministic contract and transaction
# checks; no live model is timed by this harness.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

case_name=all
if [ "${1:-}" = --case ]; then case_name=${2:-}; fi

make_request() { # ws req invocation intent tier trace-mode
	mr_ws=$1; mr_req=$2; mr_inv=$3; mr_iid=$4; mr_tier=$5; mr_mode=$6
	mkdir -p "$mr_req/fragments"
	cat >"$mr_req/request.yaml" <<Y
schema_version: 1
invocation_id: $mr_inv
intent: $mr_iid
ratified: true
trace_current: true
criterion_coverage: complete
question_dispositions: complete
tier: $mr_tier
trace_mode: $mr_mode
approval_freeze_ms: 10
trace_dispatch_ms: 20
trace_complete_ms: 30
feasibility_disposition_ms: 40
ratification_ms: 50
fallback_ms: 50
plans:
Y
}

add_fragment() { # req key slug title repo intent [dep-key]
	af_req=$1; af_key=$2; af_slug=$3; af_title=$4; af_repo=$5; af_iid=$6; af_dep=${7:-}
	cat >>"$af_req/request.yaml" <<Y
  - key: $af_key
    slug: $af_slug
Y
	mkdir -p "$af_req/fragments/$af_key"
	{
		printf 'schema_version: 3\nplan: @plan:%s\ntitle: %s\nstatus: draft\nobjective: %s objective\nintent: %s\n' "$af_key" "$af_title" "$af_title" "$af_iid"
		printf 'repositories:\n  - id: %s\n    purpose: %s scope\n' "$af_repo" "$af_repo"
		if [ -n "$af_dep" ]; then printf 'plan_dependencies:\n  - id: @plan:%s\n    reason: builds on %s\n' "$af_dep" "$af_dep"; fi
		printf 'product_knowledge:\n  - id: project.core\n    path: context/PROJECT.md\n    reason: Grounds the objective.\n'
		printf 'context_grounding:\n  summary: Ratified plan-ready trace.\n  constraints: []\n  decisions:\n    - One repository per plan.\n'
		printf 'knowledge_impact:\n  expected_context_units: []\n  review_on_completion: true\n'
		printf 'tasks:\n  - id: %s-001\n    title: Work in %s\n    repositories: [%s]\n    paths: [src]\n    evidence_anchors: [src/seed.txt]\n    depends_on: []\n' "$af_key" "$af_repo" "$af_repo"
		printf '    changes:\n      - Implement %s.\n    acceptance:\n      - id: %s-AC\n        statement: %s works.\n' "$af_title" "$af_key" "$af_title"
		printf '    verification:\n      - id: %s-VT\n        command: test -d src\n' "$af_key"
		printf 'execution:\n  worker: one\n  independent_verifier: required\n  max_worker_failures: 3\n'
	} >"$af_req/fragments/$af_key/plan.yaml"
	printf '# %s\n\nPlan @plan:%s.\n' "$af_title" "$af_key" >"$af_req/fragments/$af_key/PLAN.md"
}

new_ws() {
	nw_ws=$(cc_fx_ws)
	nw_iid=i001-latency
	cc_fx_intent "$nw_ws" "$nw_iid" "Approval latency" api src standard
	sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$nw_ws" "$nw_iid" >/dev/null
}

case_exact_delta_fallback() {
	manifest="$ROOT/.context-circuit/wrapper/contracts/schemas/trace-manifest.yaml"
	trace="$ROOT/.agents/skills/cc-trace/SKILL.md"
	for word in plan_ready_version repository_identity observed current cold exact delta fallback reread_paths fallback_reason discovery_rules_digest; do contains "$manifest" "$word"; done
	contains "$trace" '`.runtime/trace-cache/maps/'
	contains "$trace" "never secrets, ignored files, provider payloads, or source copies"
	contains "$trace" "reread the current intent sites"
	contains "$trace" "changed sites and"
	contains "$trace" "broad drift"
}

case_same_repository_expansion() {
	contains "$ROOT/.context-circuit/agents/worker.md" "scope_expansions"
	contains "$ROOT/.context-circuit/agents/worker.md" "not an exhaustive write allowlist"
	contains "$ROOT/.context-circuit/agents/verifier.md" "complete committed diff"
	contains "$ROOT/.context-circuit/wrapper/contracts/schemas/worker-handoff.yaml" "consistent with the approved intent"
}

case_expansion_stop() {
	contains "$ROOT/.context-circuit/agents/worker.md" "required second repository"
	contains "$ROOT/.context-circuit/agents/worker.md" "approved-decision change stops"
	contains "$ROOT/.context-circuit/agents/verifier.md" "coordinator finding"
}

case_materialize() { # mode and number of plans
	cm_mode=$1; cm_count=$2
	new_ws; cm_req="$nw_ws/request-$cm_mode"; make_request "$nw_ws" "$cm_req" "invoke-$cm_mode-$cm_count" "$nw_iid" standard "$cm_mode"
	add_fragment "$cm_req" core core "Core" api "$nw_iid"
	if [ "$cm_count" -gt 1 ]; then add_fragment "$cm_req" ui ui "UI" web "$nw_iid" core; fi
	out=$(sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" plan-stack-materialize "$nw_ws" "$cm_req")
	printf '%s\n' "$out" | grep -q '^materialized: created$' || fail "stack not created"
	core_id=$(printf '%s\n' "$out" | sed -n 's/^plan: core=//p')
	[ -d "$nw_ws/plans/$core_id" ] || fail "core plan absent"
	cc_plan_validate "$nw_ws/plans/$core_id" >/dev/null
	if [ "$cm_count" -gt 1 ]; then
		ui_id=$(printf '%s\n' "$out" | sed -n 's/^plan: ui=//p')
		[ -d "$nw_ws/plans/$ui_id" ] || fail "ui plan absent"
		contains "$nw_ws/plans/$ui_id/plan.yaml" "id: $core_id"
		cc_plan_validate "$nw_ws/plans/$ui_id" >/dev/null
	fi
	out2=$(sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" plan-stack-materialize "$nw_ws" "$cm_req")
	printf '%s\n' "$out2" | grep -q '^materialized: existing$' || fail "retry did not reuse allocation"
	assert_eq "$core_id" "$(printf '%s\n' "$out2" | sed -n 's/^plan: core=//p')"
	cm_rec="$nw_ws/.runtime/materialization/invoke-$cm_mode-$cm_count"
	contains "$cm_rec/timing.yaml" "trace_mode: $cm_mode"
	prev=0
	for phase in approval_freeze_ms trace_dispatch_ms trace_complete_ms feasibility_disposition_ms ratification_ms fallback_ms allocation_render_ms validation_ms authorization_publication_ms total_ms; do
		val=$(cc_scalar "$cm_rec/timing.yaml" "$phase"); [ "$val" -ge "$prev" ] || fail "timing not monotonic at $phase"; prev=$val
	done
	rm -rf "$nw_ws"
}

case_atomicity() {
	new_ws; req="$nw_ws/request-bad"; make_request "$nw_ws" "$req" invoke-bad "$nw_iid" standard exact
	add_fragment "$req" one one One api "$nw_iid"
	add_fragment "$req" two two Two api "$nw_iid" one
	sed 's/schema_version: 3/schema_version: 99/' "$req/fragments/two/plan.yaml" >"$req/fragments/two/plan.yaml.tmp"
	mv "$req/fragments/two/plan.yaml.tmp" "$req/fragments/two/plan.yaml"
	set +e
	out=$(sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" plan-stack-materialize "$nw_ws" "$req" 2>/dev/null); rc=$?
	set -e
	[ "$rc" -ne 0 ] || fail "invalid stack succeeded"
	printf '%s\n' "$out" | grep -Fq "stage: validation" || fail "missing validation stage diagnostic"
	printf '%s\n' "$out" | grep -Fq "key: two" || fail "missing failing key diagnostic"
	printf '%s\n' "$out" | grep -Fq "error: PLAN_INVALID" || fail "missing stable error diagnostic"
	not_contains "$nw_ws/plans/INDEX.md" "one"
	not_contains "$nw_ws/plans/INDEX.md" "two"
	find "$nw_ws/plans" -mindepth 1 -maxdepth 1 -type d ! -name archive | grep -q . && fail "partial plan published"
	rm -rf "$nw_ws"
}

run_case() {
	case "$1" in
		exact-delta-fallback) case_exact_delta_fallback ;;
		same-repository-expansion) case_same_repository_expansion ;;
		expansion-stop) case_expansion_stop ;;
		warm-standard) case_materialize exact 1 ;;
		cold-multi-repository) case_materialize cold 2 ;;
		multi-plan-atomicity) case_atomicity ;;
		timing-evidence) case_materialize delta 1 ;;
		*) fail "unknown case $1" ;;
	esac
}

if [ "$case_name" = all ]; then
	for c in exact-delta-fallback same-repository-expansion expansion-stop warm-standard cold-multi-repository multi-plan-atomicity timing-evidence; do run_case "$c"; done
else
	run_case "$case_name"
fi

pass 'approval-to-plan latency'
