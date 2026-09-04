#!/bin/sh
# Context Circuit v1.0 — human-simulated scenario fixtures. Plays each scenario in
# scenarios.md / conversations.md as a deterministic engine walkthrough (the human
# turns mapped to coordinator actions) and asserts the state post-conditions, then
# checks the conversation fixture holds the reporting-language discipline (no
# internal vocabulary; the two gates and the honest assurance wording present).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
require_file "$here/scenarios.md"
require_file "$here/conversations.md"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" checkout-service development
cc_fx_repo "$ws" payments-lib development

# ============================ Scenario A ============================
# Intent (Gate 1) -> derived plan WITHIN -> execute+verify -> accept -> deliver ->
# inferred completion -> reconciliation debt pending.
iidA=i0001-checkout-retries
cc_fx_intent "$ws" "$iidA" "Checkout retries" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidA" >/dev/null                       # Gate 1
assert_eq "approved" "$(cc_scalar "$ws/intent/$iidA/contract.yaml" status)"
cc_fx_plan_intent "$ws" 0001-retry "Retry" checkout-service src/checkout "$iidA"
eng intent-authorized "$ws" 0001-retry | grep -q '^authorized: yes' || fail "A: plan should be authorized"
cc_fx_run_ok "$ws" 0001-retry checkout-service src/checkout       # worker + independent verifier pass
edirA=$(cc_fx_exec_dir "$ws" 0001-retry "$(cc_latest_execution "$ws" 0001-retry)")
eng human-acceptance-record "$edirA" maker >/dev/null             # accept the candidate
eng delivery-record "$ws" 0001-retry >/dev/null                  # Gate 2 happened
eng completion-infer "$ws" 0001-retry >/dev/null                 # inferred (Standard)
assert_eq "done" "$(cc_plan_status "$ws" 0001-retry)"
eng knowledge-debt "$ws" | grep -q 'pending_count: 1' || fail "A: completion must leave reconciliation debt"

# ============================ Scenario E ============================
# The next plan in the same knowledge scope is BLOCKED at grounding until reconciled.
iidE=i0002-coupon
cc_fx_intent "$ws" "$iidE" "Coupon field" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidE" >/dev/null
cc_fx_plan_intent "$ws" 0002-coupon "Coupon" checkout-service src/checkout "$iidE"
expect_failure eng knowledge-debt-check "$ws" 0002-coupon         # blocked by A's debt
candA=$(eng candidate-current "$ws" 0001-retry | sed -n 's/^candidate_id: //p')
eng knowledge-reconciled "$ws" "$candA" reconciled >/dev/null     # human reconciles
eng knowledge-debt-check "$ws" 0002-coupon | grep -q 'debt: clear' || fail "E: grounding must clear after reconcile"

# ============================ Scenario D ============================
# Criteria drift self-invalidates: editing an approved intent's criteria breaks the
# frozen digest, so the derived plan's authorization fails (the change re-enters
# Gate 1) and its execution is held; re-approving the changed criteria re-authorizes.
# (There is NO automated scope gate in v1.0 — a reach beyond the coarse scope is a
# feasibility question the coordinator surfaces, and scope-safety is at Gate 2.)
iidD=i0003-drift
cc_fx_intent "$ws" "$iidD" "Drift" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidD" >/dev/null
cc_fx_plan_intent "$ws" 0003-drift "Drift" checkout-service src/checkout "$iidD"
eng intent-authorized "$ws" 0003-drift | grep -q '^authorized: yes' || fail "D: plan should be authorized before the edit"
# edit the approved criteria without re-approval -> the frozen digest no longer matches
sed 's/Drift works\./Drift works within budget./' "$ws/intent/$iidD/contract.yaml" >"$ws/intent/$iidD/c.new"
mv "$ws/intent/$iidD/c.new" "$ws/intent/$iidD/contract.yaml"
eng intent-authorized "$ws" 0003-drift 2>/dev/null | grep -q 'reason: CRITERIA_CHANGED' || fail "D: a criteria change must re-gate"
expect_failure eng execution-begin "$ws" 0003-drift sess-d       # held, not proceeding on the stale approval
eng intent-approve "$ws" "$iidD" >/dev/null                      # re-approve the changed criteria (Gate 1 again)
eng intent-authorized "$ws" 0003-drift | grep -q '^authorized: yes' || fail "D: re-approval should re-authorize"

# ============================ Scenario C ============================
# Two stacked plans delivered as one pull request -> one change-set candidate.
iidC=i0004-ratelimit
cc_fx_intent "$ws" "$iidC" "Rate limit" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidC" >/dev/null
cc_fx_plan_intent "$ws" 0012-api "API" checkout-service src/checkout "$iidC"
cc_fx_plan_intent "$ws" 0013-consumer "Consumer" checkout-service src/checkout "$iidC"
cc_fx_run_ok "$ws" 0012-api checkout-service src/checkout
cc_fx_run_ok "$ws" 0013-consumer checkout-service src/checkout
cs=$(eng change-set-candidate "$ws" 0012-api 0013-consumer | sed -n 's/^change_set_candidate: //p')
case "$cs" in cand-*) : ;; *) fail "C: change-set candidate invalid" ;; esac
# one candidate for the pair, order-independent (one verification, one acceptance)
csr=$(eng change-set-candidate "$ws" 0013-consumer 0012-api | sed -n 's/^change_set_candidate: //p')
assert_eq "$cs" "$csr"

# ============================ Scenario B ============================
# Explore writes no plan; promotion is the first moment an intent and plan exist.
eng pair-begin "$ws" checkout-service explore-b >/dev/null
require_file "$ws/.runtime/pairing/explore-b/pointer.yaml"
test ! -d "$ws/plans/0099-explore" || fail "B: un-promoted Explore must not create a plan file"
# promote: attach an intent (Explore-eligible: single repo, bounded, no signal) and
# a plan of record; the tier can then rise to Standard and gain a verifier.
iidB=i0005-banner
cc_fx_intent "$ws" "$iidB" "Banner copy" checkout-service "src/banner" explore
eng intent-approve "$ws" "$iidB" >/dev/null                      # Explore approves (no risk signal)
assert_eq "explore" "$(cc_scalar "$ws/intent/$iidB/contract.yaml" tier)"
cc_fx_plan_intent "$ws" 0005-banner "Banner" checkout-service src/banner "$iidB"
require_file "$ws/plans/0005-banner/plan.yaml"                   # promotion authored a plan of record

# ==================== conversation language discipline ====================
conv="$here/conversations.md"
# only inspect the dialogue lines (quoted, starting with "> ")
dlg=$(grep '^>' "$conv" || true)
for bad in 'cc_' 'engine.sh' 'worktree' 'plan.yaml' 'contract.yaml' '.runtime' 'cand-' 'contract_digest' 'INV-' 'i0001' '0012-'; do
	printf '%s\n' "$dlg" | grep -F "$bad" >/dev/null 2>&1 && fail "conversation leaks internal term: $bad" || :
done
# the honest assurance wording and both gates are present
contains "$conv" "human-supervised, not verified"
contains "$conv" "independently checked"
contains "$conv" "Approve this and I'll build it"
contains "$conv" "open the pull request"
contains "$conv" "combined result"
contains "$conv" "haven't"

pass 'v1.0 human-simulated scenarios'
