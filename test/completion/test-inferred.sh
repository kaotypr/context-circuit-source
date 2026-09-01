#!/bin/sh
# Context Circuit v1.0 — inferred completion + change-set delivery (Phase 5).
# Completion is inferred from candidate acceptance + delivery at Explore/Standard,
# and explicit at Critical (INV-COMPLETE-01). A change set — plans delivered as one
# pull request — has one candidate, so it is verified and accepted once (pain 6).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# --- Standard: completion inferred from candidate acceptance + delivery ---
iid=i0001-std
cc_fx_intent "$ws" "$iid" "Std" api "src"
eng intent-approve "$ws" "$iid" >/dev/null
cc_fx_plan_intent "$ws" 0001-std "Std" api src "$iid"
cc_fx_run_ok "$ws" 0001-std api src
edir=$(cc_fx_exec_dir "$ws" 0001-std "$(cc_latest_execution "$ws" 0001-std)")

# acceptance alone does not complete — delivery has not happened
eng human-acceptance-record "$edir" alice >/dev/null
expect_failure eng completion-infer "$ws" 0001-std
assert_eq "approved" "$(cc_plan_status "$ws" 0001-std)"

# record delivery (Gate 2 happened); completion is then inferred, not a manual flip
eng delivery-record "$ws" 0001-std >/dev/null
out=$(eng completion-infer "$ws" 0001-std)
printf '%s\n' "$out" | grep -q 'human_completion: inferred' || fail "Standard completion must be inferred"
assert_eq "done" "$(cc_plan_status "$ws" 0001-std)"
contains "$edir/completion.yaml" "human_completion: inferred"
# inferred completion still emits the reconciliation-debt marker (M4)
printf '%s\n' "$out" | grep -q 'reconciliation_debt: pending' || fail "inferred completion must emit debt"
# idempotent: inferring again reports done, does not error
eng completion-infer "$ws" 0001-std >/dev/null

# a post-delivery change (new candidate) makes the delivery signal stale
# (guards against completing a candidate that moved after delivery was recorded)
cc_fx_intent "$ws" i0009-drift "Drift" web "src"
eng intent-approve "$ws" i0009-drift >/dev/null
cc_fx_plan_intent "$ws" 0009-drift "Drift" web src i0009-drift
cc_fx_run_ok "$ws" 0009-drift web src
edir9=$(cc_fx_exec_dir "$ws" 0009-drift "$(cc_latest_execution "$ws" 0009-drift)")
eng human-acceptance-record "$edir9" alice >/dev/null
eng delivery-record "$ws" 0009-drift >/dev/null
wt9="$ws/.runtime/worktrees/0009-drift/web"
printf 'z\n' >>"$wt9/src/mod.txt"; git -C "$wt9" add -A; git -C "$wt9" commit -q -m "feat(web): more"
cc_worker_commit_record "$edir9" web repair >/dev/null
expect_failure eng completion-infer "$ws" 0009-drift    # delivery signal is now stale

# --- Critical: inferred completion is refused; explicit is required ---
iid2=i0002-crit
cc_fx_intent "$ws" "$iid2" "Crit" api "src/crit"
awk '/^tier:/{print "tier: critical"; next}{print}' "$ws/intent/$iid2/contract.yaml" >"$ws/intent/$iid2/c.new"
mv "$ws/intent/$iid2/c.new" "$ws/intent/$iid2/contract.yaml"
eng intent-approve "$ws" "$iid2" >/dev/null
cc_fx_plan_intent "$ws" 0002-crit "Crit" api src/crit "$iid2"
cc_fx_run_ok "$ws" 0002-crit api src/crit
edir2=$(cc_fx_exec_dir "$ws" 0002-crit "$(cc_latest_execution "$ws" 0002-crit)")
contains "$edir2/execution.yaml" "tier: critical"
eng human-acceptance-record "$edir2" alice >/dev/null
eng delivery-record "$ws" 0002-crit >/dev/null
expect_failure eng completion-infer "$ws" 0002-crit           # Critical must be explicit
assert_eq "approved" "$(cc_plan_status "$ws" 0002-crit)"
eng plan-complete "$ws" 0002-crit >/dev/null                  # explicit human completion
assert_eq "done" "$(cc_plan_status "$ws" 0002-crit)"
contains "$edir2/completion.yaml" "human_completion: accepted"

# --- change set: one candidate over several plans, deterministic ---
cs1=$(eng change-set-candidate "$ws" 0001-std 0002-crit | sed -n 's/^change_set_candidate: //p')
case "$cs1" in cand-*) : ;; *) fail "change-set candidate not cand-<hash>: $cs1" ;; esac
cs2=$(eng change-set-candidate "$ws" 0001-std 0002-crit | sed -n 's/^change_set_candidate: //p')
assert_eq "$cs1" "$cs2"                                        # deterministic
# a change set of one equals that plan's own candidate identity space (cand-*)
cs_one=$(eng change-set-candidate "$ws" 0001-std | sed -n 's/^change_set_candidate: //p')
case "$cs_one" in cand-*) : ;; *) fail "single-plan change set invalid: $cs_one" ;; esac
# order independence: the set is a set, not a sequence
cs_rev=$(eng change-set-candidate "$ws" 0002-crit 0001-std | sed -n 's/^change_set_candidate: //p')
assert_eq "$cs1" "$cs_rev"

pass 'inferred completion and change-set delivery'
