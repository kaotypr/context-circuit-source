#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" contracts development
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# base commits before execution (to prove the anchor checkout is untouched)
api_base=$(git -C "$ws/repositories/api" rev-parse HEAD)

cc_fx_plan "$ws" 0001-checkout "Checkout v2" "contracts api web"
cc_plan_approve "$ws" 0001-checkout >/dev/null
exec=$(cc_execution_begin "$ws" 0001-checkout sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-checkout "$exec")

# --- one isolated branch + worktree per affected repository ---
for r in contracts api web; do
	require_dir "$ws/.runtime/worktrees/0001-checkout/$r"
	contains "$edir/repositories/$r.yaml" "branch: cc/0001-checkout/$r"
	contains "$edir/repositories/$r.yaml" "allowed_paths: [src]"
done

# --- the anchor checkout is not written during execution ---
assert_eq "$api_base" "$(git -C "$ws/repositories/api" rev-parse HEAD)"

# --- commit-before-verification is enforced: no worker commit -> verifier prepare fails ---
expect_failure cc_verifier_prepare "$edir"

# --- one worker commits every repository in dependency order ---
cc_attempt_begin "$edir" >/dev/null
for r in contracts api web; do
	cc_fx_commit "$ws" 0001-checkout "$r" "impl-$r"
	cc_worker_commit_record "$edir" "$r" implementation >/dev/null
done
cc_verifier_prepare "$edir" >/dev/null

# --- independent verifier passes over the latest commits ---
cc_verifier_result_record "$edir" 1 passed >/dev/null
assert_eq "verified" "$(cc_execution_status "$edir")"

# --- worktree isolation: each repo has its own branch tip, distinct from anchor ---
web_tip=$(git -C "$ws/.runtime/worktrees/0001-checkout/web" rev-parse HEAD)
web_base=$(cc_scalar "$edir/repositories/web.yaml" base_commit)
test "$web_tip" != "$web_base" || fail "web worktree tip equals base"

# --- blocked verifier is not a worker failure ---
cc_fx_plan "$ws" 0002-blocked "Blocked" "api"
cc_plan_approve "$ws" 0002-blocked >/dev/null
exec2=$(cc_execution_begin "$ws" 0002-blocked sess2 | sed -n 's/^execution_id: //p')
edir2=$(cc_fx_exec_dir "$ws" 0002-blocked "$exec2")
cc_attempt_begin "$edir2" >/dev/null
cc_fx_commit "$ws" 0002-blocked api impl
cc_worker_commit_record "$edir2" api implementation >/dev/null
cc_verifier_result_record "$edir2" 1 blocked >/dev/null
assert_eq "blocked" "$(cc_execution_status "$edir2")"
assert_eq "0" "$(cc_scalar "$edir2/execution.yaml" worker_failures)"
cc_repair_allowed "$edir2" >/dev/null   # still allowed; block is not a failure

# --- waived is non-passing and is not a worker failure ---
cc_fx_plan "$ws" 0004-waived "Waived" "api"
cc_plan_approve "$ws" 0004-waived >/dev/null
exec4=$(cc_execution_begin "$ws" 0004-waived sess4 | sed -n 's/^execution_id: //p')
edir4=$(cc_fx_exec_dir "$ws" 0004-waived "$exec4")
cc_attempt_begin "$edir4" >/dev/null
cc_fx_commit "$ws" 0004-waived api impl
cc_worker_commit_record "$edir4" api implementation >/dev/null
cc_verifier_result_record "$edir4" 1 waived >/dev/null
assert_eq "0" "$(cc_scalar "$edir4/execution.yaml" worker_failures)"
test "verified" != "$(cc_execution_status "$edir4")" || fail "waived must not verify"

# --- repair reuses the same worktree and creates a new commit ---
cc_fx_plan "$ws" 0003-repair "Repair" "api"
cc_plan_approve "$ws" 0003-repair >/dev/null
exec3=$(cc_execution_begin "$ws" 0003-repair sess3 | sed -n 's/^execution_id: //p')
edir3=$(cc_fx_exec_dir "$ws" 0003-repair "$exec3")
cc_attempt_begin "$edir3" >/dev/null
cc_fx_commit "$ws" 0003-repair api impl
cc_worker_commit_record "$edir3" api implementation >/dev/null
first_commit=$(cc_scalar "$edir3/repositories/api.yaml" latest_commit)
cc_verifier_result_record "$edir3" 1 failed >/dev/null
assert_eq "repairing" "$(cc_execution_status "$edir3")"
cc_attempt_begin "$edir3" >/dev/null
cc_fx_commit "$ws" 0003-repair api fix
cc_worker_commit_record "$edir3" api repair >/dev/null
repair_commit=$(cc_scalar "$edir3/repositories/api.yaml" latest_commit)
test "$first_commit" != "$repair_commit" || fail "repair did not create a new commit"
cc_verifier_result_record "$edir3" 2 passed >/dev/null
assert_eq "verified" "$(cc_execution_status "$edir3")"

# --- concurrent second writer is blocked (one active writer) ---
expect_failure cc_lock_acquire "$ws" 0001-checkout other-session

pass 'execution'
