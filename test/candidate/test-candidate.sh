#!/bin/sh
# Context Circuit v1.0 — Mechanism 2: candidate identity + mechanical evidence
# staleness (INV-CANDIDATE-01). Evidence and acceptance bind to a candidate; any
# new commit OR a criteria change yields a new candidate and voids prior evidence
# and acceptance. "It passed earlier" is impossible by construction.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
iid=i0001-retries
cc_fx_intent "$ws" "$iid" "Retries" api "src"
eng intent-approve "$ws" "$iid" >/dev/null
cc_fx_plan_intent "$ws" 0001-retries "Retries" api src "$iid"

exec=$(eng execution-begin "$ws" 0001-retries sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-retries "$exec")

# --- the execution records the intent's frozen contract digest (candidate input) ---
frozen=$(cc_scalar "$ws/intent/$iid/contract.yaml" contract_digest)
contains "$edir/execution.yaml" "contract_digest: $frozen"

cc_attempt_begin "$edir" >/dev/null
wt="$ws/.runtime/worktrees/0001-retries/api"
mkdir -p "$wt/src"; printf 'a\n' >"$wt/src/mod.txt"
git -C "$wt" add -A; git -C "$wt" commit -q -m "feat(api): add mod.txt"
cc_worker_commit_record "$edir" api implementation >/dev/null
cc_verifier_prepare "$edir" >/dev/null

# --- candidate-digest records candidate.yaml and is deterministic ---
cand=$(eng candidate-digest "$ws" 0001-retries "$exec" | sed -n 's/^candidate_id: //p')
require_file "$edir/candidate.yaml"
contains "$edir/candidate.yaml" "candidate_id: $cand"
case "$cand" in cand-*) : ;; *) fail "candidate id not cand-<hash>: $cand" ;; esac
again=$(eng candidate-digest "$ws" 0001-retries "$exec" | sed -n 's/^candidate_id: //p')
assert_eq "$cand" "$again"                                    # deterministic
cur=$(eng candidate-current "$ws" 0001-retries | sed -n 's/^candidate_id: //p')
assert_eq "$cand" "$cur"

# --- verifier pass binds its result to the candidate ---
cc_verifier_result_record "$edir" 1 passed >/dev/null
contains "$edir/attempts/001/verifier.yaml" "candidate_id: $cand"
contains "$edir/execution.yaml" "verified_candidate: $cand"
expect_failure eng completion-ready "$ws" 0001-retries

# --- human acceptance is first-class and candidate-bound ---
eng human-acceptance-record "$edir" alice >/dev/null
require_file "$edir/human-acceptance.yaml"
contains "$edir/human-acceptance.yaml" "candidate_id: $cand"
contains "$edir/human-acceptance.yaml" "accepted_by: alice"
eng human-acceptance-current "$edir" >/dev/null

# --- a NEW commit yields a new candidate and VOIDS prior evidence + acceptance ---
cc_attempt_begin "$edir" >/dev/null
printf 'b\n' >>"$wt/src/mod.txt"; git -C "$wt" add -A; git -C "$wt" commit -q -m "feat(api): more"
cc_worker_commit_record "$edir" api repair >/dev/null
cand2=$(eng candidate-current "$ws" 0001-retries | sed -n 's/^candidate_id: //p')
test "$cand2" != "$cand" || fail "a new commit must change the candidate"
expect_failure eng completion-ready "$ws" 0001-retries      # verified evidence is stale/void
expect_failure eng human-acceptance-current "$edir"          # acceptance no longer matches

# re-verify + re-accept against the new candidate restores eligibility
cc_verifier_prepare "$edir" >/dev/null
cc_verifier_result_record "$edir" 2 passed >/dev/null
contains "$edir/execution.yaml" "verified_candidate: $cand2"
eng human-acceptance-record "$edir" alice >/dev/null
eng completion-ready "$ws" 0001-retries >/dev/null
eng human-acceptance-current "$edir" >/dev/null

# --- a CRITERIA change (re-approved on the intent) also voids the candidate ---
sed 's/Retries works/Retries works within budget/' "$ws/intent/$iid/contract.yaml" >"$ws/intent/$iid/contract.yaml.new"
mv "$ws/intent/$iid/contract.yaml.new" "$ws/intent/$iid/contract.yaml"
eng intent-approve "$ws" "$iid" >/dev/null                   # re-freeze the changed criteria
cand3=$(eng candidate-current "$ws" 0001-retries | sed -n 's/^candidate_id: //p')
test "$cand3" != "$cand2" || fail "a criteria change must change the candidate"
expect_failure eng completion-ready "$ws" 0001-retries       # old green cannot survive the goalpost move
expect_failure eng human-acceptance-current "$edir"

pass 'candidate identity'
