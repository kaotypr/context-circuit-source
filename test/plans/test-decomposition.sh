#!/bin/sh
# Decomposition acceptance: one bounded plan with ordered tasks versus a stack of
# independently executable plans under one approved intent.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development

iid=i001-decomposition
cc_fx_intent "$ws" "$iid" "Decomposition" api "src/decomposition"
cc_intent_approve "$ws" "$iid" >/dev/null

single_plan_shape() {
	file=$1
	grep -Fq "depends_on: [API-001]" "$file" &&
	! grep -Fq "plan_dependencies:" "$file" &&
	grep -Fq "worker: one" "$file" &&
	grep -Fq "independent_verifier: required" "$file"
}

stacked_plan_shape() {
	file=$1
	grep -Fq "plan_dependencies:" "$file" &&
	grep -Fq "reason: Consumer depends on the API plan" "$file" &&
	grep -Fq "worker: one" "$file" &&
	grep -Fq "independent_verifier: required" "$file"
}

# A single bounded change can carry several ordered tasks without an inter-plan edge.
cc_fx_plan_intent "$ws" 0001-one "One plan" api src/decomposition/one "$iid"
single="$ws/plans/0001-one/plan.yaml"
awk '
/^execution:/ && !added {
	print "  - id: API-002"
	print "    title: Work in src/decomposition/one/detail"
	print "    repositories: [api]"
	print "    paths: [src/decomposition/one/detail]"
	print "    depends_on: [API-001]"
	print "    changes:"
	print "      - Add the detail task after the primary task."
	print "    acceptance:"
	print "      - id: API-002-AC"
	print "        statement: The detail task is represented."
	print "    verification:"
	print "      - id: API-002-VT"
	print "        command: test -f src/decomposition/one/detail/mod.txt"
	added=1
}
{print}
' "$single" >"$single.new"
mv "$single.new" "$single"
cc_plan_validate "$ws/plans/0001-one" >/dev/null
single_plan_shape "$single" || fail "one-plan fixture lost its embedded task boundary"

# Independent partitions are separate plans under the same intent, with an explicit
# edge reason and a separate worker/verifier lifecycle for each plan.
cc_fx_plan_intent "$ws" 0002-api "Stack API" api src/decomposition/api "$iid"
cc_fx_plan_intent "$ws" 0003-consumer "Stack consumer" api src/decomposition/consumer "$iid"
consumer="$ws/plans/0003-consumer/plan.yaml"
awk '
/^product_knowledge:/ && !added {
	print "plan_dependencies:"
	print "  - id: 0002-api"
	print "    reason: Consumer depends on the API plan"
	added=1
}
{print}
' "$consumer" >"$consumer.new"
mv "$consumer.new" "$consumer"
cc_plan_validate "$ws/plans/0002-api" >/dev/null
cc_plan_validate "$ws/plans/0003-consumer" >/dev/null
stacked_plan_shape "$consumer" || fail "stacked fixture lost its dependency or lifecycle"

# The two wrong interpretations fail the corresponding shape checks: an embedded
# task is not a stack, and a dependent partition is not a combined one-plan record.
expect_failure stacked_plan_shape "$single"
expect_failure single_plan_shape "$consumer"

wrong_combined="$ws/plans/0003-consumer/plan-without-edge.yaml"
awk '
/^plan_dependencies:/ {skip=1; next}
skip && /^product_knowledge:/ {skip=0}
!skip {print}
' "$consumer" >"$wrong_combined"
expect_failure stacked_plan_shape "$wrong_combined"

# The coordinator guidance rejects tier-only splitting and combining independent
# boundaries; these are review failures rather than new runtime gates.
contains "$ROOT/.context-circuit/docs/planning.md" "Do not split solely to match an assurance tier"
contains "$ROOT/.context-circuit/docs/planning.md" "do not combine genuinely independent execution"
contains "$ROOT/product/.agents/skills/cc-plan/SKILL.md" "One approved intent may"

# The stack is independently ready and verified one plan at a time, while retaining
# the same upstream authorization and allowing one later change-set candidate.
cc_intent_authorized "$ws" 0002-api | grep -q '^authorized: yes' || fail "API plan not authorized"
cc_intent_authorized "$ws" 0003-consumer | grep -q '^authorized: yes' || fail "consumer plan not authorized"
before=$(cc_plan_ready "$ws" 0003-consumer || :)
printf '%s' "$before" | grep -q '^readiness: waiting' || fail "dependent plan should wait"
cc_fx_run_ok "$ws" 0002-api api src/decomposition/api
after=$(cc_plan_ready "$ws" 0003-consumer || :)
printf '%s' "$after" | grep -q '^readiness: ready' || fail "dependent plan should become ready"
cc_fx_run_ok "$ws" 0003-consumer api src/decomposition/consumer
assert_eq "verified" "$(cc_execution_status "$ws/.runtime/executions/0002-api/$(cc_latest_execution "$ws" 0002-api)")"
assert_eq "verified" "$(cc_execution_status "$ws/.runtime/executions/0003-consumer/$(cc_latest_execution "$ws" 0003-consumer)")"
assert_eq "draft" "$(cc_plan_status "$ws" 0002-api)"
assert_eq "draft" "$(cc_plan_status "$ws" 0003-consumer)"

cs=$(sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" change-set-candidate "$ws" 0002-api 0003-consumer | sed -n 's/^change_set_candidate: //p')
cs_reversed=$(sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" change-set-candidate "$ws" 0003-consumer 0002-api | sed -n 's/^change_set_candidate: //p')
case "$cs" in cand-*) : ;; *) fail "stacked change-set candidate is invalid" ;; esac
assert_eq "$cs" "$cs_reversed"

pass 'decomposition'
