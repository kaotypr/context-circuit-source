#!/bin/sh
# Run-stack semantics (v0.6): base selection (base branch / stack / integration),
# clean integration merge, BASE_UNBUILDABLE, stale-base rebuild, readiness
# AND-join, failure containment (held descendant), the run-stack partition, and
# the delivery drift guard (drift / rebase / rebase-conflict). These edge cases
# are proven deterministically here because a live model cannot be made to
# fail/conflict/contend on command.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
apidir="$ws/repositories/api"

edir_of() { printf '%s/.runtime/executions/%s/%s' "$ws" "$1" "$(cc_latest_execution "$ws" "$1")"; }

# run a plan writing CONTENT to <path>/mod.txt (a clean, verifying execution)
run_content() { # pid repo path content
	rc_ex=$(cc_execution_begin "$ws" "$1" "$1-w" | sed -n 's/^execution_id: //p')
	rc_ed="$ws/.runtime/executions/$1/$rc_ex"; rc_wt="$ws/.runtime/worktrees/$1/$2"
	cc_attempt_begin "$rc_ed" >/dev/null
	mkdir -p "$rc_wt/$3"; printf '%s\n' "$4" >"$rc_wt/$3/mod.txt"
	git -C "$rc_wt" add -A; git -C "$rc_wt" commit -q -m "feat($2): $3"
	cc_worker_commit_record "$rc_ed" "$2" implementation >/dev/null
	cc_verifier_prepare "$rc_ed" >/dev/null
	cc_verifier_result_record "$rc_ed" 1 passed >/dev/null
}

# ============================================================================
# 1. Three base kinds + clean integration merge + the based_on / built_on proofs
# ============================================================================
cc_fx_plan_ex "$ws" 0001-root  "Root"  api src/root  ""
cc_fx_plan_ex "$ws" 0002-stack "Stack" api src/sp    "0001-root"
cc_fx_plan_ex "$ws" 0003-side  "Side"  api src/op    "0001-root"
cc_fx_plan_ex "$ws" 0004-integ "Integ" api src/ip    "0002-stack 0003-side"

cc_fx_run_ok "$ws" 0001-root api src/root
cc_fx_run_ok "$ws" 0002-stack api src/sp
cc_fx_run_ok "$ws" 0003-side api src/op
cc_fx_run_ok "$ws" 0004-integ api src/ip

# base branch: no based_on recorded for the root plan
r1="$(edir_of 0001-root)/repositories/api.yaml"
not_contains "$r1" "based_on:"
basetip=$(git -C "$apidir" rev-parse --verify refs/heads/development)

# stack base: single predecessor, based_on names it, base == predecessor tip
r2="$(edir_of 0002-stack)/repositories/api.yaml"
contains "$r2" "based_on: [0001-root]"
p1tip=$(git -C "$apidir" rev-parse --verify refs/heads/cc/0001-root/api)
base2=$(cc_scalar "$r2" base_commit)
assert_eq "$p1tip" "$base2"

# integration base: two predecessors, runtime-authored merge, both are ancestors
r4="$(edir_of 0004-integ)/repositories/api.yaml"
contains "$r4" "based_on: [0002-stack, 0003-side]"
base4=$(cc_scalar "$r4" base_commit)
nparents=$(git -C "$apidir" rev-list --parents -n1 "$base4" | wc -w | tr -d ' ')
[ "$nparents" -ge 3 ] || fail "integration base is not a merge (parents token count=$nparents)"
# built_on ORDER proof: each predecessor's verified commit is an ancestor of the base
p2commit=$(cc_scalar "$(edir_of 0002-stack)/repositories/api.yaml" latest_commit)
p3commit=$(cc_scalar "$(edir_of 0003-side)/repositories/api.yaml" latest_commit)
git -C "$apidir" merge-base --is-ancestor "$p2commit" "$base4" || fail "0002 not an ancestor of the integration base"
git -C "$apidir" merge-base --is-ancestor "$p3commit" "$base4" || fail "0003 not an ancestor of the integration base"
# the kept base ref lives in the reserved namespace, never nested under a branch ref
git -C "$apidir" rev-parse --verify "refs/cc-base/0004-integ/api" >/dev/null 2>&1 \
	|| fail "integration base ref missing at refs/cc-base/0004-integ/api"

# ============================================================================
# 2. BASE_UNBUILDABLE: two predecessors edit the SAME path with different content
#    -> the integration merge conflicts -> blocked, not a worker failure.
# ============================================================================
cc_fx_plan_ex "$ws" 0005-ca "ConflA" api src/shared ""
cc_fx_plan_ex "$ws" 0006-cb "ConflB" api src/shared ""
cc_fx_plan_ex "$ws" 0007-cc "ConflC" api src/cc "0005-ca 0006-cb"
run_content 0005-ca api src/shared "AAA"
run_content 0006-cb api src/shared "BBB"
expect_failure cc_execution_begin "$ws" 0007-cc 0007-cc-w
# a blocked execution record is preserved; it is not a worker failure
ce="$(edir_of 0007-cc)"
assert_eq "blocked" "$(cc_execution_status "$ce")"
assert_eq "0" "$(cc_scalar "$ce/execution.yaml" worker_failures)"
contains "$ce/execution.yaml" "blocked_reason: BASE_UNBUILDABLE"

# ============================================================================
# 3. Stale base rebuild: a predecessor is repaired after the dependent's base was
#    built; a fresh cc_base_prepare invalidates and rebuilds from the new tip.
# ============================================================================
cc_fx_plan_ex "$ws" 0008-s1 "S1" api src/s1 ""
cc_fx_plan_ex "$ws" 0009-s2 "S2" api src/s2 ""
cc_fx_plan_ex "$ws" 0010-sd "SD" api src/sd "0008-s1 0009-s2"
cc_fx_run_ok "$ws" 0008-s1 api src/s1
cc_fx_run_ok "$ws" 0009-s2 api src/s2
out1=$(cc_base_prepare "$ws" 0010-sd api)
base_first=$(printf '%s' "$out1" | sed -n 's/^base_commit: //p')
# repair predecessor 0008: advance cc/0008-s1/api with a new commit
s1wt="$ws/.runtime/worktrees/0008-s1/api"
printf 'repair\n' >>"$s1wt/src/s1/mod.txt"
git -C "$s1wt" add -A; git -C "$s1wt" commit -q -m "fix(api): repair s1"
s1new=$(git -C "$apidir" rev-parse --verify refs/heads/cc/0008-s1/api)
out2=$(cc_base_prepare "$ws" 0010-sd api)
base_second=$(printf '%s' "$out2" | sed -n 's/^base_commit: //p')
test "$base_first" != "$base_second" || fail "stale base was not rebuilt"
git -C "$apidir" merge-base --is-ancestor "$s1new" "$base_second" || fail "rebuilt base excludes the repaired predecessor tip"
printf '%s' "$out2" | grep -q '^worktree_reused: false' || fail "stale rebuild should not report a reuse"

# ============================================================================
# 4. Readiness AND-join + run-stack partition + lease gate
# ============================================================================
cc_fx_plan_ex "$ws" 0011-ra "RA" api src/ra ""
cc_fx_plan_ex "$ws" 0012-rb "RB" api src/rb "0011-ra"
# before 0011 runs, 0012 waits on its dependency
ready11=$(cc_plan_ready "$ws" 0011-ra || :); ready12=$(cc_plan_ready "$ws" 0012-rb || :)
printf '%s' "$ready11" | grep -q '^readiness: ready'   || fail "0011 should be ready"
printf '%s' "$ready12" | grep -q '^readiness: waiting' || fail "0012 should wait on its dep"
part=$(cc_run_stack_ready "$ws" 0011-ra 0012-rb)
printf '%s\n' "$part" | grep -q '^0011-ra: ready'   || fail "partition: 0011 not ready"
printf '%s\n' "$part" | grep -q '^0012-rb: waiting' || fail "partition: 0012 not waiting"
# run 0011; now 0012 becomes ready and 0011 partitions as verified (not runnable)
cc_fx_run_ok "$ws" 0011-ra api src/ra
ready12b=$(cc_plan_ready "$ws" 0012-rb || :)
printf '%s' "$ready12b" | grep -q '^readiness: ready' || fail "0012 should be ready after its dep verified"
part2=$(cc_run_stack_ready "$ws" 0011-ra 0012-rb)
printf '%s\n' "$part2" | grep -q '^0011-ra: verified' || fail "partition: verified plan is not runnable"
printf '%s\n' "$part2" | grep -q '^0012-rb: ready'    || fail "partition: 0012 not ready after dep"

# One approved intent can yield a real stack: the dependency is inter-plan order,
# while each plan keeps its own bounded task and verifier lifecycle.
shared_iid=i025-shared-decomposition
cc_fx_intent "$ws" "$shared_iid" "Shared decomposition" api "src/shared-stack"
cc_intent_approve "$ws" "$shared_iid" >/dev/null
cc_fx_plan_intent "$ws" 0024-api "Shared API" api src/shared-stack/api "$shared_iid"
cc_fx_plan_intent "$ws" 0025-consumer "Shared consumer" api src/shared-stack/consumer "$shared_iid"
shared_consumer="$ws/plans/0025-consumer/plan.yaml"
awk '
/^product_knowledge:/ && !added {
	print "plan_dependencies:"
	print "  - id: 0024-api"
	print "    reason: Consumer depends on the API plan"
	added=1
}
{print}
' "$shared_consumer" >"$shared_consumer.new"
mv "$shared_consumer.new" "$shared_consumer"
cc_plan_validate "$ws/plans/0024-api" >/dev/null
cc_plan_validate "$ws/plans/0025-consumer" >/dev/null
contains "$shared_consumer" "plan_dependencies:"
contains "$shared_consumer" "reason: Consumer depends on the API plan"
contains "$shared_consumer" "worker: one"
contains "$shared_consumer" "independent_verifier: required"
shared_before=$(cc_run_stack_ready "$ws" 0024-api 0025-consumer)
printf '%s\n' "$shared_before" | grep -q '^0024-api: ready' || fail "shared API should be ready"
printf '%s\n' "$shared_before" | grep -q '^0025-consumer: waiting' || fail "shared consumer should wait"
cc_fx_run_ok "$ws" 0024-api api src/shared-stack/api
shared_after=$(cc_run_stack_ready "$ws" 0024-api 0025-consumer)
printf '%s\n' "$shared_after" | grep -q '^0024-api: verified' || fail "shared API should be verified"
printf '%s\n' "$shared_after" | grep -q '^0025-consumer: ready' || fail "shared consumer should become ready"
cc_fx_run_ok "$ws" 0025-consumer api src/shared-stack/consumer
assert_eq "verified" "$(cc_execution_status "$(edir_of 0024-api)")"
assert_eq "verified" "$(cc_execution_status "$(edir_of 0025-consumer)")"

# lease gate: an unrelated held lease on the same region makes a plan wait
cc_fx_plan_ex "$ws" 0013-l1 "L1" api src/lease ""
cc_fx_plan_ex "$ws" 0014-l2 "L2" api src/lease ""
cc_lease_acquire "$ws" api 0013-l1 "src/lease" >/dev/null
readyL=$(cc_plan_ready "$ws" 0014-l2 || :)
printf '%s' "$readyL" | grep -q '^readiness: waiting' || fail "0014 should wait on a held lease"
cc_lease_release "$ws" api 0013-l1 >/dev/null
readyLb=$(cc_plan_ready "$ws" 0014-l2 || :)
printf '%s' "$readyLb" | grep -q '^readiness: ready' || fail "0014 should be ready once the lease frees"

# ============================================================================
# 5. Failure containment: a failed plan holds only its descendant; unrelated
#    verified/ready plans are unaffected. Terminal/in-progress plans are NOT ready.
# ============================================================================
cc_fx_plan_ex "$ws" 0015-fa "FA" api src/fa ""
cc_fx_plan_ex "$ws" 0016-fb "FB" api src/fb "0015-fa"
cc_fx_plan_ex "$ws" 0017-fc "FC" api src/fc ""
# fail 0015 with three rejected attempts
fx_ex=$(cc_execution_begin "$ws" 0015-fa 0015-fa-w | sed -n 's/^execution_id: //p')
fx_ed="$ws/.runtime/executions/0015-fa/$fx_ex"; fx_wt="$ws/.runtime/worktrees/0015-fa/api"
i=1
while [ "$i" -le 3 ]; do
	cc_attempt_begin "$fx_ed" >/dev/null
	mkdir -p "$fx_wt/src/fa"; printf 'attempt %s\n' "$i" >>"$fx_wt/src/fa/wip.txt"
	git -C "$fx_wt" add -A; git -C "$fx_wt" commit -q -m "wip $i"
	cc_worker_commit_record "$fx_ed" api implementation >/dev/null
	cc_verifier_prepare "$fx_ed" >/dev/null
	cc_verifier_result_record "$fx_ed" "$i" failed >/dev/null
	i=$((i + 1))
done
assert_eq "failed" "$(cc_execution_status "$fx_ed")"
fpart=$(cc_run_stack_ready "$ws" 0015-fa 0016-fb 0017-fc)
printf '%s\n' "$fpart" | grep -q '^0015-fa: failed'  || fail "containment: 0015 not failed"
printf '%s\n' "$fpart" | grep -q '^0016-fb: blocked' || fail "containment: 0016 (descendant) not held/blocked"
printf '%s\n' "$fpart" | grep -q '^0017-fc: ready'   || fail "containment: unrelated 0017 not ready"

# refused: a plan whose parent intent is NOT approved cannot be authorized, so it is
# refused without blocking the rest (v1.0: authorization is the approved intent with
# unchanged criteria, not a separate plan-approval status and not a scope gate)
cc_fx_plan_ex "$ws" 0018-draft "Draft" api src/dr ""
awk '/^status:/{print "status: draft"; next}{print}' "$ws/intent/i018-draft/contract.yaml" >"$ws/intent/i018-draft/c.new"
mv "$ws/intent/i018-draft/c.new" "$ws/intent/i018-draft/contract.yaml"
dpart=$(cc_run_stack_ready "$ws" 0017-fc 0018-draft)
printf '%s\n' "$dpart" | grep -q '^0018-draft: refused' || fail "partition: unauthorized plan not refused"
printf '%s\n' "$dpart" | grep -q '^0017-fc: ready'      || fail "partition: refusal blocked the rest"

# ============================================================================
# 6. Delivery drift guard: drift detect, rebase + re-verify flag, rebase conflict.
# ============================================================================
cc_fx_plan_ex "$ws" 0019-dv "DV" api src/dv ""
cc_fx_run_ok "$ws" 0019-dv api src/dv
# no drift yet
cc_delivery_drift "$ws" 0019-dv | grep -q '^drift_detected: false' || fail "no drift expected before a sibling merges"
# a sibling merged: advance the base branch on a DIFFERENT path (clean rebase)
git -C "$apidir" checkout -q development
mkdir -p "$apidir/src/sibling"; printf 'sib\n' >"$apidir/src/sibling/f.txt"
git -C "$apidir" add -A; git -C "$apidir" commit -q -m "feat(api): sibling merged"
newtip=$(git -C "$apidir" rev-parse --verify refs/heads/development)
cc_delivery_drift "$ws" 0019-dv | grep -q '^drift_detected: true' || fail "drift should be detected after the base branch advanced"
reb=$(cc_delivery_rebase "$ws" 0019-dv)
printf '%s\n' "$reb" | grep -q '^reverify_required: true' || fail "rebase must flag re-verification"
dvbase=$(cc_scalar "$(edir_of 0019-dv)/repositories/api.yaml" base_commit)
assert_eq "$newtip" "$dvbase"
git -C "$apidir" merge-base --is-ancestor "$newtip" "$(cc_scalar "$(edir_of 0019-dv)/repositories/api.yaml" latest_commit)" \
	|| fail "rebased branch does not contain the new base tip"

# rebase conflict: the base branch advances on the SAME file the plan added -> conflict
cc_fx_plan_ex "$ws" 0020-cf "CF" api src/cf ""
cc_fx_run_ok "$ws" 0020-cf api src/cf
git -C "$apidir" checkout -q development
mkdir -p "$apidir/src/cf"; printf 'base-side\n' >"$apidir/src/cf/mod.txt"
git -C "$apidir" add -A; git -C "$apidir" commit -q -m "feat(api): base touches src/cf/mod.txt"
expect_failure cc_delivery_rebase "$ws" 0020-cf
# the plan's work is preserved after the aborted rebase
require_dir "$ws/.runtime/worktrees/0020-cf/api"

# ============================================================================
# 7. Multi-repo: a CROSS-repo dependency is an ordering gate only — it gates
#    readiness but never produces a git base (INV-CONCURRENCY-02). A same-repo
#    dependency in the other repo still stacks.
# ============================================================================
cc_fx_repo "$ws" web development
webdir="$ws/repositories/web"
cc_fx_plan_ex "$ws" 0021-xrapi "XR api"  api src/xr ""
cc_fx_plan_ex "$ws" 0022-xrweb "XR web"  web src/xr "0021-xrapi"   # cross-repo dep (web -> api)
cc_fx_plan_ex "$ws" 0023-xrweb2 "XR web2" web src/xr2 "0022-xrweb" # same-repo dep (web -> web)
# the cross-repo dependent has NO same-repo predecessor
expect_failure cc_plan_has_same_repo_pred "$ws" 0022-xrweb web
# readiness gates on the cross-repo dependency until it verifies
xr_before=$(cc_plan_ready "$ws" 0022-xrweb || :)
printf '%s' "$xr_before" | grep -q '^readiness: waiting' || fail "0022 must wait on its cross-repo dep"
cc_fx_run_ok "$ws" 0021-xrapi api src/xr
xr_after=$(cc_plan_ready "$ws" 0022-xrweb || :)
printf '%s' "$xr_after" | grep -q '^readiness: ready' || fail "0022 must be ready once the cross-repo dep verified"
# execute 0022: base is web's base tip and NO based_on is recorded (gate, not base)
cc_fx_run_ok "$ws" 0022-xrweb web src/xr
xr_rf="$(edir_of 0022-xrweb)/repositories/web.yaml"
not_contains "$xr_rf" "based_on:"
web_tip=$(git -C "$webdir" rev-parse --verify refs/heads/development)
git -C "$webdir" merge-base --is-ancestor "$web_tip" "$(cc_scalar "$xr_rf" base_commit)" \
	|| fail "cross-repo dependent must be based on its own base tip"
# a SAME-repo dependent in web still stacks on its predecessor's branch
cc_fx_run_ok "$ws" 0023-xrweb2 web src/xr2
xr2_rf="$(edir_of 0023-xrweb2)/repositories/web.yaml"
contains "$xr2_rf" "based_on: [0022-xrweb]"
xrweb_tip=$(git -C "$webdir" rev-parse --verify refs/heads/cc/0022-xrweb/web)
assert_eq "$xrweb_tip" "$(cc_scalar "$xr2_rf" base_commit)"

pass 'run-stack'
