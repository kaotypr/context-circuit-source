#!/bin/sh
# Context Circuit v1.0 — unupdated Product Knowledge must not block the next
# plan (INV-KNOWLEDGE-02). Delivery and inferred completion must not emit a
# reconcile-starting marker. Explore is planless. The runtime never writes
# Product Knowledge (INV-COMPLETE-02).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# --- a completed Standard plan does not emit a reconcile-starting debt marker ---
iid=i001-retry
cc_fx_intent "$ws" "$iid" "Retry" api "src"
eng intent-approve "$ws" "$iid" >/dev/null
cc_fx_plan_intent "$ws" 0001-retry "Retry" api src "$iid"
cc_fx_run_ok "$ws" 0001-retry api src
edir=$(cc_fx_exec_dir "$ws" 0001-retry "$(cc_latest_execution "$ws" 0001-retry)")
eng human-acceptance-record "$edir" alice >/dev/null
eng delivery-record "$ws" 0001-retry >/dev/null
expect_failure eng completion-infer "$ws" 0001-retry
assert_eq "draft" "$(cc_plan_status "$ws" 0001-retry)"
out=$(eng plan-complete "$ws" 0001-retry)
assert_eq "done" "$(cc_plan_status "$ws" 0001-retry)"
printf '%s\n' "$out" | grep -q 'reconciliation_debt:' && fail "mark-done must not emit a reconcile-starting debt marker" || :
cand=$(eng candidate-current "$ws" 0001-retry | sed -n 's/^candidate_id: //p')
test ! -f "$ws/.runtime/knowledge-debt/$cand.yaml" || fail "mark-done must not write a pending debt marker"

# --- a new Standard plan on the overlapping repo is NOT blocked ---
iid2=i002-coupon
cc_fx_intent "$ws" "$iid2" "Coupon" api "src"
eng intent-approve "$ws" "$iid2" >/dev/null
cc_fx_plan_intent "$ws" 0002-coupon "Coupon" api src "$iid2"
eng knowledge-debt-check "$ws" 0002-coupon | grep -q 'debt: clear' \
	|| fail "overlapping later plan must not be blocked solely because Product Knowledge is unupdated"

# --- a new Standard plan on a DIFFERENT repo is also CLEAR ---
iid3=i003-web
cc_fx_intent "$ws" "$iid3" "Web" web "src"
eng intent-approve "$ws" "$iid3" >/dev/null
cc_fx_plan_intent "$ws" 0003-web "Web" web src "$iid3"
sed 's/project\.core/web.core/g' "$ws/plans/0003-web/plan.yaml" >"$ws/plans/0003-web/plan.yaml.new"
mv "$ws/plans/0003-web/plan.yaml.new" "$ws/plans/0003-web/plan.yaml"
eng knowledge-debt-check "$ws" 0003-web | grep -q 'debt: clear' || fail "disjoint later plan must be clear"

# --- Explore is planless; a plan-shaped Explore record is invalid ---
iid4=i004-explore
cc_fx_intent "$ws" "$iid4" "Explore" api "src/widget" explore
eng intent-approve "$ws" "$iid4" >/dev/null
cc_fx_plan_intent "$ws" 0004-explore "Explore" api src/widget "$iid4"
expect_failure eng plan-validate "$ws/plans/0004-explore"

# --- leftover bookkeeping markers, if present, still do not block ---
mkdir -p "$ws/.runtime/knowledge-debt"
printf 'schema_version: 1\ncandidate_id: leftover\nplan: 0001-retry\nexecution_id: exec-x\nrepositories: [api]\nknowledge_units: [project.core]\nresolved: pending\ncreated_at: 0\n' \
	>"$ws/.runtime/knowledge-debt/leftover.yaml"
eng knowledge-debt-check "$ws" 0002-coupon | grep -q 'debt: clear' \
	|| fail "a leftover pending marker must not block the next plan"
eng knowledge-reconciled "$ws" leftover reconciled >/dev/null
contains "$ws/.runtime/knowledge-debt/leftover.yaml" "resolved: reconciled"

# --- deferred is still a first-class bookkeeping resolution ---
printf 'schema_version: 1\ncandidate_id: leftover2\nplan: 0003-web\nexecution_id: exec-y\nrepositories: [web]\nknowledge_units: [web.core]\nresolved: pending\ncreated_at: 0\n' \
	>"$ws/.runtime/knowledge-debt/leftover2.yaml"
eng knowledge-reconciled "$ws" leftover2 deferred >/dev/null
contains "$ws/.runtime/knowledge-debt/leftover2.yaml" "resolved: deferred"

# --- reconciliation never accepts Product Knowledge on its own (INV-KNOWLEDGE-02) ---
project_before=$(cc_digest "$ws/context/PROJECT.md")
assert_eq "$project_before" "$(cc_digest "$ws/context/PROJECT.md")"

pass 'next plan not blocked by unupdated Product Knowledge'
