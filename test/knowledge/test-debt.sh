#!/bin/sh
# Context Circuit v1.0 — Mechanism 4: the closed knowledge loop. Completion emits a
# reconciliation-debt marker keyed to the accepted candidate; the next plan's
# grounding preflight BLOCKS (Standard/Critical) or loudly WARNS (Explore) while
# delivered work in its knowledge scope remains unreconciled. The gate never
# auto-accepts knowledge; it only refuses to let the debt be forgotten
# (INV-COMPLETE-02, INV-KNOWLEDGE-02).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# --- a completed Standard plan emits a pending reconciliation-debt marker ---
iid=i0001-retry
cc_fx_intent "$ws" "$iid" "Retry" api "src"
eng intent-approve "$ws" "$iid" >/dev/null
cc_fx_plan_intent "$ws" 0001-retry "Retry" api src "$iid"
cc_fx_run_ok "$ws" 0001-retry api src
edir=$(cc_fx_exec_dir "$ws" 0001-retry "$(cc_latest_execution "$ws" 0001-retry)")
eng human-acceptance-record "$edir" alice >/dev/null
out=$(eng plan-complete "$ws" 0001-retry)
printf '%s\n' "$out" | grep -q 'reconciliation_debt: pending' || fail "completion must emit a pending debt marker"
cand=$(eng candidate-current "$ws" 0001-retry | sed -n 's/^candidate_id: //p')
require_file "$ws/.runtime/knowledge-debt/$cand.yaml"
contains "$ws/.runtime/knowledge-debt/$cand.yaml" "resolved: pending"

# --- knowledge-debt lists the pending candidate ---
list=$(eng knowledge-debt "$ws")
printf '%s\n' "$list" | grep -q "debt: $cand" || fail "knowledge-debt must list the pending candidate"
printf '%s\n' "$list" | grep -q 'pending_count: 1' || fail "expected exactly one pending debt"

# --- a new Standard plan on the overlapping repo is BLOCKED at grounding ---
iid2=i0002-coupon
cc_fx_intent "$ws" "$iid2" "Coupon" api "src"
eng intent-approve "$ws" "$iid2" >/dev/null
cc_fx_plan_intent "$ws" 0002-coupon "Coupon" api src "$iid2"
expect_failure eng knowledge-debt-check "$ws" 0002-coupon
eng knowledge-debt-check "$ws" 0002-coupon 2>/dev/null | grep -q 'debt: blocking' \
	|| eng knowledge-debt-check "$ws" 0002-coupon 2>&1 | grep -q 'debt: blocking' || fail "expected blocking"

# --- a new Standard plan on a DIFFERENT repo (no overlap) is CLEAR ---
iid3=i0003-web
cc_fx_intent "$ws" "$iid3" "Web" web "src"
eng intent-approve "$ws" "$iid3" >/dev/null
cc_fx_plan_intent "$ws" 0003-web "Web" web src "$iid3"
eng knowledge-debt-check "$ws" 0003-web | grep -q 'debt: clear' || fail "non-overlapping plan must be clear"

# --- an Explore plan on the overlapping repo WARNS (does not block; rc 0) ---
iid4=i0004-explore
cc_fx_intent "$ws" "$iid4" "Explore" api "src/widget" explore
eng intent-approve "$ws" "$iid4" >/dev/null
cc_fx_plan_intent "$ws" 0004-explore "Explore" api src/widget "$iid4"
ewarn=$(eng knowledge-debt-check "$ws" 0004-explore)   # rc 0 (would exit set -e otherwise)
printf '%s\n' "$ewarn" | grep -q 'debt: warn' || fail "Explore plan must warn, not block"

# --- reconcile clears the debt; the blocked plan then grounds cleanly ---
eng knowledge-reconciled "$ws" "$cand" reconciled >/dev/null
contains "$ws/.runtime/knowledge-debt/$cand.yaml" "resolved: reconciled"
eng knowledge-debt-check "$ws" 0002-coupon | grep -q 'debt: clear' || fail "reconciled debt must clear the block"
eng knowledge-debt "$ws" | grep -q 'pending_count: 0' || fail "no debt should remain pending"

# --- "deferred" is a first-class no-update-needed resolution ---
cc_fx_run_ok "$ws" 0003-web web src
edir3=$(cc_fx_exec_dir "$ws" 0003-web "$(cc_latest_execution "$ws" 0003-web)")
eng human-acceptance-record "$edir3" bob >/dev/null
eng plan-complete "$ws" 0003-web >/dev/null
cand3=$(eng candidate-current "$ws" 0003-web | sed -n 's/^candidate_id: //p')
eng knowledge-reconciled "$ws" "$cand3" deferred >/dev/null
contains "$ws/.runtime/knowledge-debt/$cand3.yaml" "resolved: deferred"
eng knowledge-debt "$ws" | grep -q 'pending_count: 0' || fail "deferred also clears the pending debt"

# --- reconciliation never accepts Product Knowledge on its own (INV-KNOWLEDGE-02) ---
project_before=$(cc_digest "$ws/context/PROJECT.md")
assert_eq "$project_before" "$(cc_digest "$ws/context/PROJECT.md")"

pass 'closed knowledge loop'
