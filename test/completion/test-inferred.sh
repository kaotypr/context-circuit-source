#!/bin/sh
# Context Circuit v1.0 — inferred completion + change-set delivery (Phase 5).
# Completion is inferred from candidate acceptance + delivery at Standard, and
# explicit at Critical (INV-COMPLETE-01). Explore is planless. A change set — plans
# delivered as one
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
assert_eq "draft" "$(cc_plan_status "$ws" 0001-std)"

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
assert_eq "draft" "$(cc_plan_status "$ws" 0002-crit)"
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
# a change set of one IS that plan's own candidate identity (candidate-current), exactly
cc_cur=$(eng candidate-current "$ws" 0001-std | sed -n 's/^candidate_id: //p')
assert_eq "$cc_cur" "$cs_one"
# order independence: the set is a set, not a sequence
cs_rev=$(eng change-set-candidate "$ws" 0002-crit 0001-std | sed -n 's/^change_set_candidate: //p')
assert_eq "$cs1" "$cs_rev"

# --- change-set integration verification: one PR, one integration tip, verified and
#     accepted ONCE against the change-set candidate (pain 6). Member executions
#     supply worker commits but do not require separate verifier results: the
#     integration verifier is the one assurance record for the set. ---
cc_fx_plan_ex "$ws" 0020-csa "CS A" api src/a ""
cc_fx_plan_ex "$ws" 0021-csb "CS B" api src/b ""
run_worker_only() {
	ws_arg=$1; pid_arg=$2; repo_arg=$3; path_arg=$4
	exec_arg=$(eng execution-begin "$ws_arg" "$pid_arg" "$pid_arg-w" | sed -n 's/^execution_id: //p')
	edir_arg=$(cc_fx_exec_dir "$ws_arg" "$pid_arg" "$exec_arg")
	cc_attempt_begin "$edir_arg" >/dev/null
	wt_arg="$ws_arg/.runtime/worktrees/$pid_arg/$repo_arg"
	mkdir -p "$wt_arg/$path_arg"
	printf 'member\n' >"$wt_arg/$path_arg/mod.txt"
	git -C "$wt_arg" add -A
	git -C "$wt_arg" commit -q -m "feat($repo_arg): add $path_arg/mod.txt"
	cc_worker_commit_record "$edir_arg" "$repo_arg" implementation >/dev/null
}
run_worker_only "$ws" 0020-csa api src/a
run_worker_only "$ws" 0021-csb api src/b
csp=$(eng change-set-prepare "$ws" 0020-csa 0021-csb)
csid=$(printf '%s\n' "$csp" | sed -n 's/^change_set: //p')
printf '%s\n' "$csp" | grep -q 'status: prepared' || fail "change-set prepare should succeed"
printf '%s\n' "$csp" | grep -q 'tier: standard' || fail "change-set tier should be max(members)=standard"
require_dir "$ws/.runtime/change-sets/$csid/integration/api"        # the integration tip worktree
# the integration tip contains BOTH members' files (verified once, together)
require_file "$ws/.runtime/change-sets/$csid/integration/api/src/a/mod.txt"
require_file "$ws/.runtime/change-sets/$csid/integration/api/src/b/mod.txt"
# not ready until the ONE verifier pass AND the ONE acceptance bind to the candidate
expect_failure eng change-set-ready "$ws" "$csid"
expect_failure eng change-set-verifier-record "$ws" "$csid" passed --wrote-products   # read-only
eng change-set-verifier-prepare "$ws" "$csid" >/dev/null
eng change-set-verifier-record "$ws" "$csid" passed >/dev/null
eng change-set-accept "$ws" "$csid" alice >/dev/null
csr=$(eng change-set-ready "$ws" "$csid")
printf '%s\n' "$csr" | grep -q 'change_set_ready: eligible' || fail "change set should be eligible"
printf '%s\n' "$csr" | grep -q 'assurance: independent' || fail "standard change set needs an independent pass"
# complete the whole set from the ONE change-set acceptance — every member is marked
# done (pain 6: accept once, complete the set), with a reconciliation-debt marker each
csc=$(eng change-set-complete "$ws" "$csid")
printf '%s\n' "$csc" | grep -q 'completed: 2' || fail "change-set-complete must complete both members"
assert_eq "done" "$(cc_plan_status "$ws" 0020-csa)"
assert_eq "done" "$(cc_plan_status "$ws" 0021-csb)"
require_file "$(cc_fx_exec_dir "$ws" 0021-csb "$(cc_latest_execution "$ws" 0021-csb)")/completion.yaml"
# a member commit after prepare moves the candidate and voids the prepared evidence
wta="$ws/.runtime/worktrees/0020-csa/api"
printf 'x\n' >>"$wta/src/a/mod.txt"; git -C "$wta" add -A; git -C "$wta" commit -q -m "feat(api): more a"
cc_worker_commit_record "$(cc_fx_exec_dir "$ws" 0020-csa "$(cc_latest_execution "$ws" 0020-csa)")" api repair >/dev/null
expect_failure eng change-set-ready "$ws" "$csid"                   # stale: a member moved

# base drift after preparation also voids the integrated candidate, even when
# no member branch moved (the pull request target changed underneath it).
cc_fx_plan_ex "$ws" 0024-base "Base" api src/base ""
run_worker_only "$ws" 0024-base api src/base
base_set=$(eng change-set-prepare "$ws" 0024-base)
base_csid=$(printf '%s\n' "$base_set" | sed -n 's/^change_set: //p')
api_repo="$ws/repositories/api"
printf 'base moved\n' >"$api_repo/base.txt"
git -C "$api_repo" add base.txt
git -C "$api_repo" commit -q -m "chore(api): advance base"
expect_failure eng change-set-ready "$ws" "$base_csid"

# --- change-set integration that will not build is BASE_UNBUILDABLE, not a failure ---
cc_fx_plan_ex "$ws" 0022-cx "CX" api src/shared ""
cc_fx_plan_ex "$ws" 0023-cy "CY" api src/shared ""
cc_fx_run_ok "$ws" 0022-cx api src/shared
excy=$(eng execution-begin "$ws" 0023-cy 0023-cy-w | sed -n 's/^execution_id: //p')
cc_attempt_begin "$(cc_fx_exec_dir "$ws" 0023-cy "$excy")" >/dev/null
wtcy="$ws/.runtime/worktrees/0023-cy/api"
mkdir -p "$wtcy/src/shared"; printf 'DIFFERENT\n' >"$wtcy/src/shared/mod.txt"
git -C "$wtcy" add -A; git -C "$wtcy" commit -q -m "feat(api): cy"
cc_worker_commit_record "$(cc_fx_exec_dir "$ws" 0023-cy "$excy")" api implementation >/dev/null
uo=$(eng change-set-prepare "$ws" 0022-cx 0023-cy 2>&1 || :)
printf '%s\n' "$uo" | grep -q 'BASE_UNBUILDABLE' || fail "conflicting change set must be BASE_UNBUILDABLE"

pass 'inferred completion and change-set delivery'
