#!/bin/sh
# Context Circuit v1.0 — explicit mark-done + change-set delivery.
# A plan becomes done only on an explicit mark-done at Standard and Critical
# (INV-COMPLETE-01). Delivery does not mark done and does not start reconcile
# (INV-COMPLETE-02). Explore is planless. A change set — same-repository plans
# delivered as one pull request — has one tip-map candidate, accepted once, with
# no delivery-time verifier. Change-set complete records delivery only and does
# not mark members done. Plans in different repositories never form a change
# set (INV-DELIVER-01). Named plans partition by repository covering tip.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# --- Standard: delivery and inferred completion must not mark done ---
iid=i001-std
cc_fx_intent "$ws" "$iid" "Std" api "src"
eng intent-approve "$ws" "$iid" >/dev/null
cc_fx_plan_intent "$ws" 0001-std "Std" api src "$iid"
cc_fx_run_ok "$ws" 0001-std api src
edir=$(cc_fx_exec_dir "$ws" 0001-std "$(cc_latest_execution "$ws" 0001-std)")

# acceptance alone does not complete
eng human-acceptance-record "$edir" alice >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" 0001-std)"

# record delivery (Gate 2); it must not mark done or start reconcile
eng delivery-record "$ws" 0001-std >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" 0001-std)"
expect_failure eng completion-infer "$ws" 0001-std
assert_eq "draft" "$(cc_plan_status "$ws" 0001-std)"
test ! -f "$edir/completion.yaml" || fail "delivery must not write a completion record"
test ! -d "$ws/.runtime/knowledge-debt" || [ "$(find "$ws/.runtime/knowledge-debt" -name '*.yaml' 2>/dev/null | wc -l | tr -d ' ')" = "0" ] \
	|| fail "delivery must not emit a reconcile-starting debt marker"

# explicit mark-done is the Standard done path
out=$(eng plan-complete "$ws" 0001-std)
printf '%s\n' "$out" | grep -q 'human_completion: accepted' || fail "Standard completion must be explicit mark-done"
assert_eq "done" "$(cc_plan_status "$ws" 0001-std)"
contains "$edir/completion.yaml" "human_completion: accepted"
printf '%s\n' "$out" | grep -q 'reconciliation_debt:' && fail "mark-done must not emit a reconcile-starting debt marker" || :

# a post-verify change (new candidate) refuses mark-done
cc_fx_intent "$ws" i009-drift "Drift" web "src"
eng intent-approve "$ws" i009-drift >/dev/null
cc_fx_plan_intent "$ws" 0009-drift "Drift" web src i009-drift
cc_fx_run_ok "$ws" 0009-drift web src
edir9=$(cc_fx_exec_dir "$ws" 0009-drift "$(cc_latest_execution "$ws" 0009-drift)")
eng human-acceptance-record "$edir9" alice >/dev/null
wt9="$ws/.runtime/worktrees/0009-drift/web"
printf 'z\n' >>"$wt9/src/mod.txt"; git -C "$wt9" add -A; git -C "$wt9" commit -q -m "feat(web): more"
cc_worker_commit_record "$edir9" web repair >/dev/null
expect_failure eng plan-complete "$ws" 0009-drift
assert_eq "draft" "$(cc_plan_status "$ws" 0009-drift)"

# --- Critical: inferred completion is refused; explicit is required ---
iid2=i002-crit
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
expect_failure eng completion-infer "$ws" 0002-crit
assert_eq "draft" "$(cc_plan_status "$ws" 0002-crit)"
eng plan-complete "$ws" 0002-crit >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0002-crit)"
contains "$edir2/completion.yaml" "human_completion: accepted"

# --- change set: one candidate over several plans, deterministic ---
cs1=$(eng change-set-candidate "$ws" 0001-std 0002-crit | sed -n 's/^change_set_candidate: //p')
case "$cs1" in cand-*) : ;; *) fail "change-set candidate not cand-<hash>: $cs1" ;; esac
cs2=$(eng change-set-candidate "$ws" 0001-std 0002-crit | sed -n 's/^change_set_candidate: //p')
assert_eq "$cs1" "$cs2"
cs_one=$(eng change-set-candidate "$ws" 0001-std | sed -n 's/^change_set_candidate: //p')
case "$cs_one" in cand-*) : ;; *) fail "single-plan change set invalid: $cs_one" ;; esac
cc_cur=$(eng candidate-current "$ws" 0001-std | sed -n 's/^candidate_id: //p')
assert_eq "$cc_cur" "$cs_one"
cs_rev=$(eng change-set-candidate "$ws" 0002-crit 0001-std | sed -n 's/^change_set_candidate: //p')
assert_eq "$cs1" "$cs_rev"

# --- change-set delivery: one PR from the covering execution tip; no delivery
#     verifier. Members keep their own candidate-bound passes (INV-DELIVER-01).
#     change-set-complete records Gate 2 only — it must not mark members done. ---
cc_fx_plan_ex "$ws" 0020-csa "CS A" api src/a ""
cc_fx_plan_ex "$ws" 0021-csb "CS B" api src/b "0020-csa"
cc_fx_run_ok "$ws" 0020-csa api src/a
cc_fx_run_ok "$ws" 0021-csb api src/b
csp=$(eng change-set-prepare "$ws" 0020-csa 0021-csb)
csid=$(printf '%s\n' "$csp" | sed -n 's/^change_set: //p')
printf '%s\n' "$csp" | grep -q 'status: prepared' || fail "change-set prepare should succeed"
printf '%s\n' "$csp" | grep -q 'tier: standard' || fail "change-set tier should be max(members)=standard"
printf '%s\n' "$csp" | grep -q 'source_plan: 0021-csb' || fail "covering tip should be the stacked dependent"
printf '%s\n' "$csp" | grep -q 'source_branch: cc/0021-csb/api' || fail "PR source should be the covering execution branch"
test ! -d "$ws/.runtime/change-sets/$csid/integration" || fail "delivery must not author an integration worktree"
vprep=$(eng change-set-verifier-prepare "$ws" "$csid" 2>&1 || :)
printf '%s\n' "$vprep" | grep -q 'CHANGE_SET_DELIVERY_HAS_NO_VERIFIER' || fail "change-set-verifier-prepare must fail: $vprep"
vrec=$(eng change-set-verifier-record "$ws" "$csid" passed 2>&1 || :)
printf '%s\n' "$vrec" | grep -q 'CHANGE_SET_DELIVERY_HAS_NO_VERIFIER' || fail "change-set-verifier-record must fail: $vrec"
expect_failure eng change-set-ready "$ws" "$csid"
eng change-set-accept "$ws" "$csid" alice >/dev/null
csr=$(eng change-set-ready "$ws" "$csid")
printf '%s\n' "$csr" | grep -q 'change_set_ready: eligible' || fail "change set should be eligible"
printf '%s\n' "$csr" | grep -q 'assurance: independent' || fail "standard change set reuses member independent passes"
csc=$(eng change-set-complete "$ws" "$csid")
printf '%s\n' "$csc" | grep -q 'status: delivered' || fail "change-set-complete must record delivery: $csc"
printf '%s\n' "$csc" | grep -q 'delivered: 2' || fail "change-set-complete must count both members: $csc"
printf '%s\n' "$csc" | grep -q 'delivered-and-completed' && fail "change-set-complete must not fuse delivered-and-completed" || :
assert_eq "draft" "$(cc_plan_status "$ws" 0020-csa)"
assert_eq "draft" "$(cc_plan_status "$ws" 0021-csb)"
test ! -f "$(cc_fx_exec_dir "$ws" 0021-csb "$(cc_latest_execution "$ws" 0021-csb)")/completion.yaml" \
	|| fail "change-set-complete must not write member completion records"
# a member commit after prepare moves the candidate and voids the recorded set
wta="$ws/.runtime/worktrees/0020-csa/api"
printf 'x\n' >>"$wta/src/a/mod.txt"; git -C "$wta" add -A; git -C "$wta" commit -q -m "feat(api): more a"
cc_worker_commit_record "$(cc_fx_exec_dir "$ws" 0020-csa "$(cc_latest_execution "$ws" 0020-csa)")" api repair >/dev/null
expect_failure eng change-set-ready "$ws" "$csid"

# --- parallel tips cannot ship as one pull request (no delivery merge) ---
cc_fx_plan_ex "$ws" 0022-cx "CX" api src/cx ""
cc_fx_plan_ex "$ws" 0023-cy "CY" api src/cy ""
cc_fx_run_ok "$ws" 0022-cx api src/cx
cc_fx_run_ok "$ws" 0023-cy api src/cy
uo=$(eng change-set-prepare "$ws" 0022-cx 0023-cy 2>&1 || :)
printf '%s\n' "$uo" | grep -q 'CHANGE_SET_NO_SINGLE_TIP' || fail "parallel tips must fail CHANGE_SET_NO_SINGLE_TIP: $uo"
po=$(eng change-set-partition "$ws" 0022-cx 0023-cy)
printf '%s\n' "$po" | grep -q 'groups: 2' || fail "parallel same-repo tips are two covering groups: $po"
printf '%s\n' "$po" | grep -q 'pull_requests: 2' || fail "parallel tips are two pull requests: $po"
printf '%s\n' "$po" | grep -q 'covering_plan: 0022-cx' || fail "each parallel tip is a covering plan: $po"
printf '%s\n' "$po" | grep -q 'covering_plan: 0023-cy' || fail "each parallel tip is a covering plan: $po"
if printf '%s\n' "$po" | grep -q 'blocked_reason:'; then fail "partition must not block sibling tips: $po"; fi

# --- cross-repo members never form a change set (INV-DELIVER-01) ---
cc_fx_plan_ex "$ws" 0025-xapi "XAPI" api src/xapi ""
cc_fx_plan_ex "$ws" 0026-xweb "XWEB" web src/xweb ""
cc_fx_run_ok "$ws" 0025-xapi api src/xapi
cc_fx_run_ok "$ws" 0026-xweb web src/xweb
xo=$(eng change-set-candidate "$ws" 0025-xapi 0026-xweb 2>&1 || :)
printf '%s\n' "$xo" | grep -q 'CHANGE_SET_CROSS_REPO' || fail "cross-repo change-set-candidate must fail CHANGE_SET_CROSS_REPO: $xo"
printf '%s\n' "$xo" | grep -q 'change_set_candidate:' && fail "cross-repo change-set-candidate must not emit a candidate"
xo2=$(eng change-set-prepare "$ws" 0025-xapi 0026-xweb 2>&1 || :)
printf '%s\n' "$xo2" | grep -q 'CHANGE_SET_CROSS_REPO' || fail "cross-repo change-set-prepare must fail CHANGE_SET_CROSS_REPO: $xo2"
printf '%s\n' "$xo2" | grep -q 'change_set:' && fail "cross-repo prepare must not emit a change set"

# --- four plans, two repositories: partition to two covering-tip PRs ---
cc_fx_plan_ex "$ws" 0030-a1 "A1" api src/a1 ""
cc_fx_plan_ex "$ws" 0031-a2 "A2" api src/a2 "0030-a1"
cc_fx_plan_ex "$ws" 0032-a3 "A3" api src/a3 "0031-a2"
cc_fx_plan_ex "$ws" 0033-b1 "B1" web src/b1 ""
cc_fx_run_ok "$ws" 0030-a1 api src/a1
cc_fx_run_ok "$ws" 0031-a2 api src/a2
cc_fx_run_ok "$ws" 0032-a3 api src/a3
cc_fx_run_ok "$ws" 0033-b1 web src/b1
part=$(eng change-set-partition "$ws" 0030-a1 0031-a2 0032-a3 0033-b1)
printf '%s\n' "$part" | grep -q 'groups: 2' || fail "3 stacked in api + 1 in web must be two groups: $part"
printf '%s\n' "$part" | grep -q 'pull_requests: 2' || fail "two covering tips are two pull requests: $part"
printf '%s\n' "$part" | grep -q 'covering_plan: 0032-a3' || fail "api covering tip is the last stacked member: $part"
printf '%s\n' "$part" | grep -q 'source_branch: cc/0032-a3/api' || fail "api PR source is the covering branch: $part"
printf '%s\n' "$part" | grep -q 'covering_plan: 0033-b1' || fail "web covering tip is the lone plan: $part"
printf '%s\n' "$part" | grep -q 'source_branch: cc/0033-b1/web' || fail "web PR source is the lone branch: $part"
ncand=$(printf '%s\n' "$part" | grep -c '^candidate_id: ')
[ "$ncand" -eq 2 ] || fail "two groups must emit two candidates: $part"
xo4=$(eng change-set-prepare "$ws" 0030-a1 0031-a2 0032-a3 0033-b1 2>&1 || :)
printf '%s\n' "$xo4" | grep -q 'CHANGE_SET_CROSS_REPO' || fail "mixed-repo prepare must still fail: $xo4"
csp3=$(eng change-set-prepare "$ws" 0030-a1 0031-a2 0032-a3)
printf '%s\n' "$csp3" | grep -q 'source_plan: 0032-a3' || fail "api stack prepare covering tip: $csp3"
printf '%s\n' "$csp3" | grep -q 'source_branch: cc/0032-a3/api' || fail "api stack PR source: $csp3"

# --- sibling stacks in one repository are two covering-tip PRs, not zero ---
cc_fx_plan_ex "$ws" 0040-root "Root" api src/root ""
cc_fx_plan_ex "$ws" 0041-left "Left" api src/left "0040-root"
cc_fx_plan_ex "$ws" 0042-right "Right" api src/right "0040-root"
cc_fx_run_ok "$ws" 0040-root api src/root
cc_fx_run_ok "$ws" 0041-left api src/left
cc_fx_run_ok "$ws" 0042-right api src/right
fork=$(eng change-set-partition "$ws" 0040-root 0041-left 0042-right)
printf '%s\n' "$fork" | grep -q 'groups: 2' || fail "1→2 and 1→3 must be two covering groups: $fork"
printf '%s\n' "$fork" | grep -q 'pull_requests: 2' || fail "sibling stacks are two pull requests: $fork"
printf '%s\n' "$fork" | grep -q 'covering_plan: 0041-left' || fail "left stack covering tip: $fork"
printf '%s\n' "$fork" | grep -q 'covering_plan: 0042-right' || fail "right stack covering tip: $fork"
printf '%s\n' "$fork" | grep -q 'source_branch: cc/0041-left/api' || fail "left PR source: $fork"
printf '%s\n' "$fork" | grep -q 'source_branch: cc/0042-right/api' || fail "right PR source: $fork"
printf '%s\n' "$fork" | grep -q 'members: 0040-root, 0041-left' || fail "left PR contains the shared root: $fork"
printf '%s\n' "$fork" | grep -q 'members: 0040-root, 0042-right' || fail "right PR contains the shared root: $fork"
fo=$(eng change-set-prepare "$ws" 0040-root 0041-left 0042-right 2>&1 || :)
printf '%s\n' "$fo" | grep -q 'CHANGE_SET_NO_SINGLE_TIP' || fail "forcing sibling tips into one prepare must fail: $fo"

pass 'explicit mark-done and change-set delivery'
