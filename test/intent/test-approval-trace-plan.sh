#!/bin/sh
# Approval → planner → feasibility → same-turn plans; INTENT.md status after a
# feasible look. contract.yaml remains the approval authority.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ci="$ROOT/product/.agents/skills/cc-intent/SKILL.md"
tr="$ROOT/product/.agents/skills/cc-trace/SKILL.md"
cp="$ROOT/product/.agents/skills/cc-plan/SKILL.md"
coord="$ROOT/.context-circuit/agents/coordinator.md"
tmpl="$ROOT/.context-circuit/docs/templates/intent.md"

# --- 0008: planner and feasibility run after approval, before any plan ---
contains "$ci" "before any plan is written"
contains "$ci" "cc-trace"
contains "$ci" "feasibility"
contains "$ci" "intent-approve"
contains "$tr" "before any plan"
contains "$tr" "feasibility"
contains "$coord" "before writing any plan"
contains "$coord" "cc-trace"
contains "$ci" "cannot skip or stand in for the planner"
contains "$ci" "in that same turn"
contains "$cp" "in the same"
contains "$cp" "write no plan"
contains "$tr" "in the same turn"
contains "$coord" "in that same turn"
contains "$cp" "does not start execution"
contains "$cp" "plan-stack-materialize"
contains "$cp" "publishes all plans and index rows or none"
contains "$tr" "planner"

# --- 0009: INTENT.md status stays truthful; contract.yaml is still approval ---
contains "$ci" "intent-human-status"
contains "$ci" "approved, look complete, feasible"
contains "$tmpl" "approved, look complete, feasible"
contains "$tmpl" "Keep the line; do not remove it"
contains "$ROOT/.context-circuit/wrapper/runtime/engine.sh" "intent-human-status"
contains "$ROOT/.context-circuit/wrapper/runtime/engine.sh" "cc_intent_human_status"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
iid=i001-status-sync
cc_fx_intent "$ws" "$iid" "Status sync" api "src"
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid" >/dev/null
assert_eq "approved" "$(cc_scalar "$ws/intent/$iid/contract.yaml" status)"
assert_eq "approved" "$(sed -n 's/^_Status:[[:space:]]*\([^,._]*\).*$/\1/p' "$ws/intent/$iid/INTENT.md")"
grep -q '^_Status:' "$ws/intent/$iid/INTENT.md" || fail "INTENT.md status line was removed"

sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-human-status "$ws" "$iid" "approved, look complete, feasible" >/dev/null
assert_eq "approved" "$(cc_scalar "$ws/intent/$iid/contract.yaml" status)"
contains "$ws/intent/$iid/INTENT.md" "_Status: approved, look complete, feasible."
not_contains "$ws/intent/$iid/INTENT.md" "waiting for your approval"
grep -q '^_Status:' "$ws/intent/$iid/INTENT.md" || fail "INTENT.md status line missing after feasible sync"

pass 'approval then trace then plan'
