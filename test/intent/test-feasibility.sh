#!/bin/sh
# Context Circuit v1.0 — the feasibility check + derives-from-an-approved-intent
# authorization (INV-INTENT-02, reworked). v1.0 REPLACES the automated scope-envelope
# check with a feasibility check (coordinator judgment, not an engine verb) and keeps
# only a scope-FREE runtime authorization: a plan derives from an APPROVED intent whose
# criteria are unchanged. Scope-safety is settled at delivery (Gate 2), so a plan that
# reaches beyond the intent's coarse scope is still AUTHORIZED here — this suite pins
# that deliberate absence of a scope gate as much as the authorization itself.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# the trace manifest is a first-class schema (the tracer's recorded report)
require_file "$ROOT/.context-circuit/wrapper/contracts/schemas/trace-manifest.yaml"
contains "$ROOT/.context-circuit/wrapper/contracts/schemas/trace-manifest.yaml" "done_checks"
contains "$ROOT/.context-circuit/wrapper/contracts/schemas/trace-manifest.yaml" "out_of_scope_reach"

iid=i001-checkout
cc_fx_intent "$ws" "$iid" "Checkout" checkout-service "src/checkout test/checkout"
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid" >/dev/null

auth() { sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-authorized "$ws" "$1" 2>/dev/null; }
auth_reason() { sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-authorized "$ws" "$1" 2>/dev/null | sed -n 's/^reason: //p'; }

# --- AUTHORIZED: an in-scope plan derived from the approved intent ---
cc_fx_plan_intent "$ws" 0001-child "Child" checkout-service src/checkout/retry "$iid"
auth 0001-child | grep -q '^authorized: yes' || fail "in-scope plan should be authorized"

# One approved intent can authorize more than one stacked plan; the second plan
# does not introduce another approval state or gate.
cc_fx_plan_intent "$ws" 0002-sibling "Sibling" checkout-service src/checkout/summary "$iid"
auth 0002-sibling | grep -q '^authorized: yes' || fail "stacked sibling should be authorized"
assert_eq "draft" "$(cc_plan_status "$ws" 0002-sibling)"

# --- NO SCOPE GATE (the deliberate v1.0 change): a plan that reaches a repository or
# path OUTSIDE the intent's coarse scope is STILL authorized — scope-safety is at
# delivery (Gate 2), not an automated gate here. Under the old envelope check these
# would have been held; under v1.0 they proceed. ---
cc_fx_plan_intent "$ws" 0003-newrepo "New repo" payments-lib src/pay "$iid"
auth 0003-newrepo | grep -q '^authorized: yes' || fail "an out-of-scope repository is NOT gated here (Gate 2 owns scope)"
cc_fx_plan_intent "$ws" 0004-broader "Broader" checkout-service src "$iid"
auth 0004-broader | grep -q '^authorized: yes' || fail "a broader-than-scope path is NOT gated here (Gate 2 owns scope)"
cc_fx_plan_intent "$ws" 0005-wide "Repo-wide" checkout-service . "$iid"
auth 0005-wide | grep -q '^authorized: yes' || fail "a repo-wide region is NOT gated here (Gate 2 owns scope)"

# --- UNAUTHORIZED (fail upward): a criteria edit after approval, until re-approved ---
cc_fx_plan_intent "$ws" 0007-ok "Ok" checkout-service src/checkout "$iid"
auth 0007-ok | grep -q '^authorized: yes' || fail "in-scope plan should be authorized before edit"
cp "$ws/intent/$iid/contract.yaml" "$ws/intent/$iid/contract.yaml.bak"
sed 's/works\./works within budget./' "$ws/intent/$iid/contract.yaml.bak" >"$ws/intent/$iid/contract.yaml"
expect_failure auth 0007-ok
assert_eq "CRITERIA_CHANGED" "$(auth_reason 0007-ok)"
# re-approving on the changed criteria restores authorization (Gate 1 re-entered)
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid" >/dev/null
auth 0007-ok | grep -q '^authorized: yes' || fail "re-approval on changed criteria should re-authorize"
rm -f "$ws/intent/$iid/contract.yaml.bak"

# --- UNAUTHORIZED: an intent that is not approved yet ---
iid2=i002-draft
cc_fx_intent "$ws" "$iid2" "Draft" api src
cc_fx_plan_intent "$ws" 0008-draft "On draft" api src "$iid2"
expect_failure auth 0008-draft
assert_eq "INTENT_NOT_APPROVED" "$(auth_reason 0008-draft)"

# --- UNAUTHORIZED: a plan naming a missing intent ---
cc_fx_plan_intent "$ws" 0009-missing "Missing" api src i099-nope
expect_failure auth 0009-missing
assert_eq "INTENT_MISSING" "$(auth_reason 0009-missing)"

# --- UNAUTHORIZED: a plan that names no parent intent at all ---
mkdir -p "$ws/plans/0010-noparent/tasks"
cat >"$ws/plans/0010-noparent/plan.yaml" <<Y
schema_version: 3
plan: 0010-noparent
title: No parent
status: draft
objective: a plan with no intent
repositories:
  - id: checkout-service
    purpose: none
tasks:
  - id: CS-001
    title: Work
    repositories: [checkout-service]
    paths: [src/checkout]
    depends_on: []
Y
printf '# No parent\n' >"$ws/plans/0010-noparent/PLAN.md"
expect_failure auth 0010-noparent

# --- authorization gates execution as a preflight ---
cc_fx_repo "$ws" checkout-service development
# a plan on an unapproved intent cannot begin execution (re-gate, not proceed)
expect_failure sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" execution-begin "$ws" 0008-draft sess-x
# an authorized plan begins execution even though the plan itself was never explicitly
# approved (INV-EXEC-01 reworked); it stays draft (no intermediate "approved")
assert_eq "draft" "$(cc_plan_status "$ws" 0007-ok)"
exec=$(sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" execution-begin "$ws" 0007-ok sess-ok | sed -n 's/^execution_id: //p')
test -n "$exec" || fail "authorized plan should begin execution without plan approval"
assert_eq "draft" "$(cc_plan_status "$ws" 0007-ok)"

pass 'intent feasibility + authorization'
