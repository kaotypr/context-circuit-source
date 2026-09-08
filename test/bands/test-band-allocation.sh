#!/bin/sh
# Band-scoped intent and plan id allocation.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

W="$ROOT/.context-circuit/wrapper"
ENG="$W/runtime/engine.sh"
case_name=all
if [ "${1:-}" = --case ]; then case_name=${2:-}; fi

eng() { sh "$ENG" "$@"; }

expect_code() {
	ec_want=$1
	shift
	ec_err=$( { "$@" >/tmp/cc-ba-out.$$; } 2>&1 || true )
	printf '%s\n' "$ec_err" | grep -Fq "$ec_want" || fail "expected '$ec_want' in stderr, got: $ec_err"
	if grep -E '^(i[0-9]{3}-|[0-9]{4}-)' /tmp/cc-ba-out.$$ >/dev/null 2>&1; then
		fail "failure emitted an id: $(cat /tmp/cc-ba-out.$$)"
	fi
	if grep -E '^(intent_start|intent_end|plan_start|plan_end):' /tmp/cc-ba-out.$$ >/dev/null 2>&1; then
		fail "failure invented a band: $(cat /tmp/cc-ba-out.$$)"
	fi
	rm -f /tmp/cc-ba-out.$$
}

intent_prefix() {
	ip=${1#i}; ip=${ip%%-*}
	ip=$(printf '%s' "$ip" | sed 's/^0*//'); [ -n "$ip" ] || ip=0
	printf '%s' "$ip"
}

plan_prefix() {
	pp=${1%%-*}
	pp=$(printf '%s' "$pp" | sed 's/^0*//'); [ -n "$pp" ] || pp=0
	printf '%s' "$pp"
}

plant_intent() { mkdir -p "$1/intent/$2"; }
plant_plan() { mkdir -p "$1/plans/$2"; }

ba_make_request() { # ws req invocation intent
	bm_ws=$1; bm_req=$2; bm_inv=$3; bm_iid=$4
	mkdir -p "$bm_req/fragments"
	cat >"$bm_req/request.yaml" <<Y
schema_version: 1
invocation_id: $bm_inv
intent: $bm_iid
ratified: true
trace_current: true
criterion_coverage: complete
question_dispositions: complete
tier: standard
trace_mode: exact
approval_freeze_ms: 10
trace_dispatch_ms: 20
trace_complete_ms: 30
feasibility_disposition_ms: 40
ratification_ms: 50
fallback_ms: 50
plans:
Y
}

ba_add_fragment() { # req key slug title repo intent [dep-key]
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

case_intent_parallel() {
	ws=$(cc_fx_ws)
	trap 'rm -rf "$ws"' EXIT HUP INT TERM
	cc_fx_roster_pair "$ws"
	cc_fx_member_identity "$ws" kao

	kao1=$(eng intent-allocate-id "$ws" kao-first)
	assert_eq "i001-kao-first" "$kao1"
	plant_intent "$ws" "$kao1"

	cc_fx_member_identity "$ws" bobby
	bob1=$(eng intent-allocate-id "$ws" bobby-first)
	assert_eq "i100-bobby-first" "$bob1"
	plant_intent "$ws" "$bob1"

	cc_fx_member_identity "$ws" kao
	kao2=$(eng intent-allocate-id "$ws" kao-second)
	assert_eq "i002-kao-second" "$kao2"
	plant_intent "$ws" "$kao2"

	cc_fx_member_identity "$ws" bobby
	bob2=$(eng intent-allocate-id "$ws" bobby-second)
	assert_eq "i101-bobby-second" "$bob2"

	k1=$(intent_prefix "$kao1"); k2=$(intent_prefix "$kao2")
	b1=$(intent_prefix "$bob1"); b2=$(intent_prefix "$bob2")
	[ "$k1" -ne "$b1" ] && [ "$k1" -ne "$b2" ] && [ "$k2" -ne "$b1" ] && [ "$k2" -ne "$b2" ] \
		|| fail "parallel intent prefixes overlapped: $kao1 $kao2 $bob1 $bob2"
	printf '%s\n' "$kao1" "$kao2" "$bob1" "$bob2" | grep -E 'kao|bobby' >/dev/null \
		|| fail "expected slugs only, missing fixture slugs"

	trap - EXIT HUP INT TERM
	rm -rf "$ws"
}

case_plan_stack_consecutive() {
	ws=$(cc_fx_ws)
	trap 'rm -rf "$ws"' EXIT HUP INT TERM
	cc_fx_roster_pair "$ws"
	cc_fx_intent "$ws" i001-stack "Stack intent" api src standard
	eng intent-approve "$ws" i001-stack >/dev/null

	# Bobby allocates first so kao's 0001 range must not wait on a higher global max.
	cc_fx_member_identity "$ws" bobby
	bob_req="$ws/request-bobby"
	ba_make_request "$ws" "$bob_req" invoke-bobby i001-stack
	ba_add_fragment "$bob_req" core core "Core" api i001-stack
	ba_add_fragment "$bob_req" ui ui "UI" web i001-stack core
	bob_out=$(eng plan-stack-materialize "$ws" "$bob_req") \
		|| fail "bobby stack materialize failed"
	printf '%s\n' "$bob_out" | grep -q '^materialized: created$' || fail "bobby stack not created"
	bob_core=$(printf '%s\n' "$bob_out" | sed -n 's/^plan: core=//p')
	bob_ui=$(printf '%s\n' "$bob_out" | sed -n 's/^plan: ui=//p')
	assert_eq "1000-core" "$bob_core"
	assert_eq "1001-ui" "$bob_ui"
	[ "$(( $(plan_prefix "$bob_ui") - $(plan_prefix "$bob_core") ))" -eq 1 ] \
		|| fail "bobby stack ids were not consecutive"
	[ -d "$ws/plans/$bob_core" ] && [ -d "$ws/plans/$bob_ui" ] || fail "bobby plans absent"

	cc_fx_member_identity "$ws" kao
	kao_req="$ws/request-kao"
	ba_make_request "$ws" "$kao_req" invoke-kao i001-stack
	ba_add_fragment "$kao_req" core core "Core" api i001-stack
	ba_add_fragment "$kao_req" ui ui "UI" web i001-stack core
	kao_out=$(eng plan-stack-materialize "$ws" "$kao_req") \
		|| fail "kao stack materialize failed"
	printf '%s\n' "$kao_out" | grep -q '^materialized: created$' || fail "kao stack not created"
	kao_core=$(printf '%s\n' "$kao_out" | sed -n 's/^plan: core=//p')
	kao_ui=$(printf '%s\n' "$kao_out" | sed -n 's/^plan: ui=//p')
	assert_eq "0001-core" "$kao_core"
	assert_eq "0002-ui" "$kao_ui"
	[ "$(( $(plan_prefix "$kao_ui") - $(plan_prefix "$kao_core") ))" -eq 1 ] \
		|| fail "kao stack ids were not consecutive"

	printf '%s\n' "$kao_core" "$kao_ui" "$bob_core" "$bob_ui" | grep -E '^[0-9]{4}-[a-z0-9-]+$' >/dev/null \
		|| fail "plan ids left NNNN-slug form"
	printf '%s\n' "$kao_core" "$kao_ui" "$bob_core" "$bob_ui" | grep -E 'kao|bobby|member' \
		&& fail "plan id contains a member token" || :

	retry=$(eng plan-stack-materialize "$ws" "$kao_req")
	printf '%s\n' "$retry" | grep -q '^materialized: existing$' || fail "retry did not reuse allocation"
	assert_eq "$kao_core" "$(printf '%s\n' "$retry" | sed -n 's/^plan: core=//p')"

	contains "$W/runtime/engine.sh" 'cc_wp_branch="cc/$cc_wp_plan/$cc_wp_id"'
	grep -E 'cc_wp_branch="cc/\$' "$W/runtime/engine.sh" | grep -E 'member|identity' \
		&& fail "execution branch form gained a member token" || :

	trap - EXIT HUP INT TERM
	rm -rf "$ws"
}

case_never_reuse_archived() {
	ws=$(cc_fx_ws)
	trap 'rm -rf "$ws"' EXIT HUP INT TERM
	cc_fx_roster "$ws" kao 1 99 1 999
	cc_fx_member_identity "$ws" kao

	iid=$(eng intent-allocate-id "$ws" keep)
	assert_eq "i001-keep" "$iid"
	cc_fx_intent "$ws" "$iid" "Keep" api src
	eng intent-archive "$ws" "$iid" >/dev/null
	iid2=$(eng intent-allocate-id "$ws" after-archive)
	assert_eq "i002-after-archive" "$iid2"
	require_dir "$ws/intent/archive/$iid"

	pid=$(eng plan-allocate-id "$ws" keep-plan)
	assert_eq "0001-keep-plan" "$pid"
	plant_plan "$ws" "$pid"
	cc_plan_archive "$ws" "$pid" >/dev/null
	pid2=$(eng plan-allocate-id "$ws" after-plan-archive)
	assert_eq "0002-after-plan-archive" "$pid2"
	require_dir "$ws/plans/archive/$pid"

	trap - EXIT HUP INT TERM
	rm -rf "$ws"
}

case_exhaustion_overlap_collision() {
	ws=$(cc_fx_ws)
	trap 'rm -rf "$ws"' EXIT HUP INT TERM

	# missing identity fails closed without inventing a band or emitting an id
	rm -f "$ws/member.local.yaml"
	expect_code MEMBER_IDENTITY_MISSING eng intent-allocate-id "$ws" no-id
	expect_code MEMBER_IDENTITY_MISSING eng plan-allocate-id "$ws" no-id

	# overlapping roster bands
	printf '%s\n' \
		'schema_version: 1' \
		'members:' \
		'  a:' \
		'    display_name: a' \
		'    intent_band_start: 1' \
		'    intent_band_end: 50' \
		'    plan_band_start: 1' \
		'    plan_band_end: 100' \
		'  b:' \
		'    display_name: b' \
		'    intent_band_start: 50' \
		'    intent_band_end: 99' \
		'    plan_band_start: 101' \
		'    plan_band_end: 200' \
		>"$ws/members.yaml"
	cc_fx_member_identity "$ws" a
	expect_code MEMBER_BAND_OVERLAP eng intent-allocate-id "$ws" overlap
	expect_code MEMBER_BAND_OVERLAP eng plan-allocate-id "$ws" overlap

	# duplicate numeric prefixes fail closed
	cc_fx_roster "$ws" kao 1 99 1 999
	cc_fx_member_identity "$ws" kao
	mkdir -p "$ws/intent/i010-alpha" "$ws/intent/i010-beta"
	expect_code INTENT_PREFIX_COLLISION eng intent-allocate-id "$ws" dup-intent
	rm -rf "$ws/intent/i010-alpha" "$ws/intent/i010-beta"
	mkdir -p "$ws/plans/0010-alpha" "$ws/plans/0010-beta"
	expect_code PLAN_PREFIX_COLLISION eng plan-allocate-id "$ws" dup-plan
	rm -rf "$ws/plans/0010-alpha" "$ws/plans/0010-beta"

	# full intent band does not wrap
	cc_fx_roster "$ws" kao 1 1 1 1
	mkdir -p "$ws/intent/i001-full"
	overflow=$(eng intent-allocate-id "$ws" wrap 2>&1 || true)
	printf '%s\n' "$overflow" | grep -Fq 'INTENT_ID_EXHAUSTED' || fail "full intent band must fail clearly"
	printf '%s\n' "$overflow" | grep -E '^i002-|^i001-' && fail "intent allocation wrapped or recycled" || :
	mkdir -p "$ws/plans/0001-full"
	poverflow=$(eng plan-allocate-id "$ws" wrap 2>&1 || true)
	printf '%s\n' "$poverflow" | grep -Fq 'PLAN_ID_EXHAUSTED' || fail "full plan band must fail clearly"
	printf '%s\n' "$poverflow" | grep -E '^0002-|^0001-' && fail "plan allocation wrapped or recycled" || :

	# stack that would exceed plan_end fails without wrapping
	cc_fx_roster "$ws" kao 1 99 1 1
	rm -rf "$ws/plans/0001-full"
	cc_fx_intent "$ws" i002-stack "Tiny band stack" api src standard
	eng intent-approve "$ws" i002-stack >/dev/null
	req="$ws/request-tiny"
	ba_make_request "$ws" "$req" invoke-tiny i002-stack
	ba_add_fragment "$req" core core "Core" api i002-stack
	ba_add_fragment "$req" ui ui "UI" web i002-stack core
	set +e
	sout=$(eng plan-stack-materialize "$ws" "$req" 2>/dev/null); src=$?
	set -e
	[ "$src" -ne 0 ] || fail "oversize stack succeeded"
	printf '%s\n' "$sout" | grep -Fq "PLAN_ID_EXHAUSTED" || fail "oversize stack missing PLAN_ID_EXHAUSTED"
	find "$ws/plans" -mindepth 1 -maxdepth 1 -type d ! -name archive | grep -q . && fail "exhausted stack published a plan" || :

	# existing ids keep their numbers; only new allocation uses the band
	cc_fx_roster "$ws" kao 1 99 1 999
	mkdir -p "$ws/intent/i050-legacy"
	nid=$(eng intent-allocate-id "$ws" next-after-legacy)
	assert_eq "i051-next-after-legacy" "$nid"
	require_dir "$ws/intent/i050-legacy"
	mkdir -p "$ws/plans/0050-legacy"
	npid=$(eng plan-allocate-id "$ws" next-plan-legacy)
	assert_eq "0051-next-plan-legacy" "$npid"
	require_dir "$ws/plans/0050-legacy"

	trap - EXIT HUP INT TERM
	rm -rf "$ws"
}

run_case() {
	case "$1" in
		intent-parallel) case_intent_parallel ;;
		plan-stack-consecutive) case_plan_stack_consecutive ;;
		never-reuse-archived) case_never_reuse_archived ;;
		exhaustion-overlap-collision) case_exhaustion_overlap_collision ;;
		*) fail "unknown case $1" ;;
	esac
}

if [ "$case_name" = all ]; then
	run_case intent-parallel
	run_case plan-stack-consecutive
	run_case never-reuse-archived
	run_case exhaustion-overlap-collision
else
	run_case "$case_name"
fi

pass 'band allocation'
