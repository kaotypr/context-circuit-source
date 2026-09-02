#!/bin/sh
# Context Circuit v1.0 — crown jewel 1: the intent-envelope check (INV-INTENT-02).
# This is one of the two safety-critical checks; it carries the deepest fixtures
# and must FAIL UPWARD (any ambiguity re-gates, never passes).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

iid=i0001-checkout
cc_fx_intent "$ws" "$iid" "Checkout" checkout-service "src/checkout test/checkout"
sh "$ROOT/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid" >/dev/null

env_check() { sh "$ROOT/wrapper/runtime/engine.sh" intent-envelope-check "$ws" "$1" 2>/dev/null; }
env_reason() { sh "$ROOT/wrapper/runtime/engine.sh" intent-envelope-check "$ws" "$1" 2>/dev/null | sed -n 's/^reason: //p'; }

# --- WITHIN: a plan region that is a strict child of a scope path ---
cc_fx_plan_intent "$ws" 0001-child "Child" checkout-service src/checkout/retry "$iid"
env_check 0001-child | grep -q '^envelope: within' || fail "strict child should be WITHIN"

# --- WITHIN: a plan region equal to a scope path ---
cc_fx_plan_intent "$ws" 0002-equal "Equal" checkout-service test/checkout "$iid"
env_check 0002-equal | grep -q '^envelope: within' || fail "equal region should be WITHIN"

# --- EXCEEDS: a repository not named in the intent scope ---
cc_fx_plan_intent "$ws" 0003-newrepo "New repo" payments-lib src/pay "$iid"
expect_failure env_check 0003-newrepo
assert_eq "NEW_REPOSITORY" "$(env_reason 0003-newrepo)"

# --- EXCEEDS: a plan region broader than any scope path (parent of scope) ---
cc_fx_plan_intent "$ws" 0004-broader "Broader" checkout-service src "$iid"
expect_failure env_check 0004-broader
assert_eq "PATH_OUTSIDE_SCOPE" "$(env_reason 0004-broader)"

# --- EXCEEDS: a repository-wide "." plan region against a bounded scope ---
cc_fx_plan_intent "$ws" 0005-wide "Repo-wide" checkout-service . "$iid"
expect_failure env_check 0005-wide
assert_eq "PATH_OUTSIDE_SCOPE" "$(env_reason 0005-wide)"

# --- EXCEEDS: a task that NAMES the repo but OMITS paths (repo-wide work) must be
# checked for containment, never pass on repository membership alone (the bypass fix) ---
mkdir -p "$ws/plans/0010-nopaths/tasks"
cat >"$ws/plans/0010-nopaths/plan.yaml" <<Y
schema_version: 3
plan: 0010-nopaths
title: No paths
status: draft
objective: a task naming the repo with no bounded paths
intent: $iid
repositories:
  - id: checkout-service
    purpose: repo-wide
tasks:
  - id: CS-001
    title: Repo-wide work
    repositories: [checkout-service]
    depends_on: []
Y
printf '# No paths\n' >"$ws/plans/0010-nopaths/PLAN.md"
expect_failure env_check 0010-nopaths
assert_eq "PATH_OUTSIDE_SCOPE" "$(env_reason 0010-nopaths)"

# --- RE-GATE (fail upward): an unresolvable / relative plan region ---
cc_fx_plan_intent "$ws" 0006-rel "Relative" checkout-service ../secrets "$iid"
expect_failure env_check 0006-rel
assert_eq "INDETERMINATE" "$(env_reason 0006-rel)"

# --- EXCEEDS: a criteria edit after approval, until re-approved ---
cc_fx_plan_intent "$ws" 0007-ok "Ok" checkout-service src/checkout "$iid"
env_check 0007-ok | grep -q '^envelope: within' || fail "in-scope plan should be WITHIN before edit"
cp "$ws/intent/$iid/contract.yaml" "$ws/intent/$iid/contract.yaml.bak"
sed 's/works\./works within budget./' "$ws/intent/$iid/contract.yaml.bak" >"$ws/intent/$iid/contract.yaml"
expect_failure env_check 0007-ok
assert_eq "CRITERIA_CHANGED" "$(env_reason 0007-ok)"

# --- RE-GATE: an intent that is not approved yet ---
mv "$ws/intent/$iid/contract.yaml.bak" "$ws/intent/$iid/contract.yaml"   # restore criteria
iid2=i0002-draft
cc_fx_intent "$ws" "$iid2" "Draft" api src
cc_fx_plan_intent "$ws" 0008-draft "On draft" api src "$iid2"
expect_failure env_check 0008-draft
assert_eq "INTENT_NOT_APPROVED" "$(env_reason 0008-draft)"

# --- RE-GATE: a plan naming a missing intent ---
cc_fx_plan_intent "$ws" 0009-missing "Missing" api src i0099-nope
expect_failure env_check 0009-missing
assert_eq "INTENT_MISSING" "$(env_reason 0009-missing)"

# --- the envelope gates execution as a preflight ---
# an EXCEEDS plan cannot begin execution (re-gate, not proceed)
cc_fx_repo "$ws" checkout-service development
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" execution-begin "$ws" 0004-broader sess-x

# a WITHIN plan is authorized by the approved intent's envelope: execution begins
# even though the plan itself was never explicitly approved (INV-EXEC-01 reworked)
assert_eq "draft" "$(cc_plan_status "$ws" 0007-ok)"
exec=$(sh "$ROOT/wrapper/runtime/engine.sh" execution-begin "$ws" 0007-ok sess-ok | sed -n 's/^execution_id: //p')
test -n "$exec" || fail "WITHIN intent-authorized plan should begin execution without plan approval"
assert_eq "approved" "$(cc_plan_status "$ws" 0007-ok)"   # derived, not a separate human gate

pass 'intent envelope (crown jewel 1)'
