#!/bin/sh
# Context Circuit v1.0 — human-simulated scenario fixtures. Plays each scenario in
# scenarios.md / conversations.md as a deterministic engine walkthrough (the human
# turns mapped to coordinator actions) and asserts the state post-conditions, then
# checks the conversation fixture holds the reporting-language discipline (no
# internal vocabulary; the two gates and the honest assurance wording present).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" "$@"; }
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
require_file "$here/scenarios.md"
require_file "$here/conversations.md"

add_embedded_task() {
	file=$1
	awk '
    /^execution:/ && !added {
        print "  - id: CHECKOUT-SERVICE-002"
        print "    title: Bound the checkout latency task"
        print "    repositories: [checkout-service]"
        print "    paths: [src/checkout/latency]"
        print "    depends_on: [CHECKOUT-SERVICE-001]"
        print "    changes:"
        print "      - Keep the latency check after the retry behavior."
        print "    acceptance:"
        print "      - id: CHECKOUT-SERVICE-002-AC"
        print "        statement: The latency task is represented."
        print "    verification:"
        print "      - id: CHECKOUT-SERVICE-002-VT"
        print "        command: test -f src/checkout/latency/mod.txt"
        added=1
    }
    {print}
    ' "$file" >"$file.new"
	mv "$file.new" "$file"
}

add_plan_dependency() {
	file=$1
	dependency=$2
	awk -v dependency="$dependency" '
    /^product_knowledge:/ && !added {
        print "plan_dependencies:"
        print "  - id: " dependency
        print "    reason: Consumer depends on the API plan"
        added=1
    }
    {print}
    ' "$file" >"$file.new"
	mv "$file.new" "$file"
}

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" checkout-service development
cc_fx_repo "$ws" payments-lib development

# ============================ Scenario A ============================
# Intent (Gate 1) -> derived plan WITHIN -> execute+verify -> accept -> deliver.
# Delivery does not mark done. Explicit mark-done does. No reconcile-starting
# debt marker.
iidA=i001-checkout-retries
cc_fx_intent "$ws" "$iidA" "Checkout retries" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidA" >/dev/null                       # Gate 1
assert_eq "approved" "$(cc_scalar "$ws/intent/$iidA/contract.yaml" status)"
cc_fx_plan_intent "$ws" 0001-retry "Retry" checkout-service src/checkout "$iidA"
add_embedded_task "$ws/plans/0001-retry/plan.yaml"
cc_plan_validate "$ws/plans/0001-retry" >/dev/null
contains "$ws/plans/0001-retry/plan.yaml" "depends_on: [CHECKOUT-SERVICE-001]"
not_contains "$ws/plans/0001-retry/plan.yaml" "plan_dependencies:"
eng intent-authorized "$ws" 0001-retry | grep -q '^authorized: yes' || fail "A: plan should be authorized"
cc_fx_run_ok "$ws" 0001-retry checkout-service src/checkout       # worker + independent verifier pass
edirA=$(cc_fx_exec_dir "$ws" 0001-retry "$(cc_latest_execution "$ws" 0001-retry)")
eng human-acceptance-record "$edirA" maker >/dev/null             # accept the candidate
eng delivery-record "$ws" 0001-retry >/dev/null                  # Gate 2 happened
assert_eq "draft" "$(cc_plan_status "$ws" 0001-retry)"            # delivery does not complete
expect_failure eng completion-infer "$ws" 0001-retry
eng plan-complete "$ws" 0001-retry >/dev/null                    # explicit mark-done
assert_eq "done" "$(cc_plan_status "$ws" 0001-retry)"
eng knowledge-debt "$ws" | grep -q 'pending_count: 0' || fail "A: delivery/mark-done must not leave reconcile-starting debt"

# ============================ Scenario E ============================
# The next plan in the same knowledge scope is NOT blocked.
iidE=i002-coupon
cc_fx_intent "$ws" "$iidE" "Coupon field" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidE" >/dev/null
cc_fx_plan_intent "$ws" 0002-coupon "Coupon" checkout-service src/checkout "$iidE"
eng knowledge-debt-check "$ws" 0002-coupon | grep -q 'debt: clear' || fail "E: overlapping later plan must not be blocked"

# ============================ Scenario D ============================
# Criteria drift self-invalidates: editing an approved intent's criteria breaks the
# frozen digest, so the derived plan's authorization fails (the change re-enters
# Gate 1) and its execution is held; re-approving the changed criteria re-authorizes.
# (There is NO automated scope gate in v1.0 — a reach beyond the coarse scope is a
# feasibility question the coordinator surfaces, and scope-safety is at Gate 2.)
iidD=i003-drift
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
iidC=i004-ratelimit
cc_fx_intent "$ws" "$iidC" "Rate limit" checkout-service "src/checkout"
eng intent-approve "$ws" "$iidC" >/dev/null
cc_fx_plan_intent "$ws" 0012-api "API" checkout-service src/checkout/api "$iidC"
cc_fx_plan_intent "$ws" 0013-consumer "Consumer" checkout-service src/checkout/consumer "$iidC"
add_plan_dependency "$ws/plans/0013-consumer/plan.yaml" 0012-api
cc_plan_validate "$ws/plans/0012-api" >/dev/null
cc_plan_validate "$ws/plans/0013-consumer" >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" 0012-api)"
assert_eq "draft" "$(cc_plan_status "$ws" 0013-consumer)"
beforeC=$(cc_plan_ready "$ws" 0013-consumer || :)
printf '%s' "$beforeC" | grep -q '^readiness: waiting' || fail "C: consumer should wait for API"
cc_fx_run_ok "$ws" 0012-api checkout-service src/checkout/api
afterC=$(cc_plan_ready "$ws" 0013-consumer || :)
printf '%s' "$afterC" | grep -q '^readiness: ready' || fail "C: consumer should become ready"
cc_fx_run_ok "$ws" 0013-consumer checkout-service src/checkout/consumer
assert_eq "verified" "$(cc_execution_status "$(cc_fx_exec_dir "$ws" 0012-api "$(cc_latest_execution "$ws" 0012-api)")")"
assert_eq "verified" "$(cc_execution_status "$(cc_fx_exec_dir "$ws" 0013-consumer "$(cc_latest_execution "$ws" 0013-consumer)")")"
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
iidB=i005-banner
cc_fx_intent "$ws" "$iidB" "Banner copy" checkout-service "src/banner" explore
eng intent-approve "$ws" "$iidB" >/dev/null                      # Explore approves (no risk signal)
assert_eq "explore" "$(cc_scalar "$ws/intent/$iidB/contract.yaml" tier)"
cc_fx_plan_intent "$ws" 0005-banner "Banner" checkout-service src/banner "$iidB"
require_file "$ws/plans/0005-banner/plan.yaml"                   # promotion authored a plan of record

# ==================== conversation language discipline ====================
conv="$here/conversations.md"
# only inspect the dialogue lines (quoted, starting with "> ")
dlg=$(grep '^>' "$conv" || true)
for bad in 'cc_' 'engine.sh' 'worktree' 'plan.yaml' 'contract.yaml' '.runtime' 'cand-' 'contract_digest' 'INV-' 'i001' '0012-'; do
	printf '%s\n' "$dlg" | grep -F "$bad" >/dev/null 2>&1 && fail "conversation leaks internal term: $bad" || :
done
# the honest assurance wording and both gates are present
contains "$conv" "human-supervised, not verified"
contains "$conv" "independently checked"
contains "$conv" "Approve this and I'll build it"
contains "$conv" "open the pull request"
contains "$conv" "combined result"
contains "$conv" "coupon field"
contains "$here/scenarios.md" "split solely to match the assurance tier"
contains "$here/scenarios.md" "partitions are combined into one lifecycle"

pass 'v1.0 human-simulated scenarios'
