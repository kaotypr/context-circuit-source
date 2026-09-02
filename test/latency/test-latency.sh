#!/bin/sh
# Execution-latency semantics: the coordinator-recorded per-role
# (model, effort) as bounded host evidence (the attempt's own started_at /
# checked_at give its duration, so no wall-clock is recorded), per-role tiering
# that never touches the failure counter, and concurrent-run-stack overlap
# arbitrated by the lease.
# These prove the additive fields exist and are monotonic, that escalation shows
# only in the recorded (model, effort), and that independent plans may overlap
# while conflicting ones
# serialize — deterministically, because a live model cannot be timed on command.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development

edir_of() { printf '%s/.runtime/executions/%s/%s' "$ws" "$1" "$(cc_latest_execution "$ws" "$1")"; }

# ============================================================================
# 1. Evidence + timing: the coordinator records only the (model, effort) each role
#    ran at; the attempt's own started_at / checked_at timestamps give its duration
#    with no coordinator step.
# ============================================================================
cc_fx_plan "$ws" 0001-meas "Measure" "api"
cc_plan_approve "$ws" 0001-meas >/dev/null
exec=$(cc_execution_begin "$ws" 0001-meas sess1 | sed -n 's/^execution_id: //p')
edir="$ws/.runtime/executions/0001-meas/$exec"
ey="$edir/execution.yaml"

# per-attempt start boundary already exists (started_at), no field needed
cc_attempt_begin "$edir" >/dev/null
wy="$edir/attempts/001/worker.yaml"
astart=$(cc_scalar "$wy" started_at)
[ -n "$astart" ] || fail "started_at empty"

cc_fx_commit "$ws" 0001-meas api impl
cc_worker_commit_record "$edir" api implementation >/dev/null
cc_verifier_prepare "$edir" >/dev/null

# coordinator records only the (model, effort) each role ran at — the one thing
# the engine cannot know. No wall-clock: the attempt boundaries give the duration.
cc_attempt_evidence_record "$edir" 1 \
	worker_model=model-hi worker_effort=high \
	verifier_model=model-lo verifier_effort=medium >/dev/null
hy="$edir/attempts/001/host-evidence.yaml"
for kv in "worker_model: model-hi" "worker_effort: high" "verifier_model: model-lo"; do
	contains "$hy" "$kv"
done
assert_eq "1" "$(cc_scalar "$ey" current_attempt)"

# the attempt's end boundary (checked_at) bounds its duration; monotonic with start
cc_verifier_result_record "$edir" 1 passed >/dev/null
vy="$edir/attempts/001/verifier.yaml"
aend=$(cc_scalar "$vy" checked_at)
[ -n "$aend" ] || fail "checked_at empty"
# ISO-8601 UTC sorts lexicographically == chronologically
{ [ "$aend" = "$astart" ] || [ "$aend" \> "$astart" ]; } || fail "attempt end precedes start (not monotonic)"

# ============================================================================
# 2. Bounded evidence surface: unknown key and unsafe value are refused.
# ============================================================================
expect_failure cc_attempt_evidence_record "$edir" 1 bogus_key=x
expect_failure cc_attempt_evidence_record "$edir" 1 "worker_model=model-hi;rm"
expect_failure cc_attempt_evidence_record "$edir" 1 malformed_no_equals
# a refused call writes nothing new
not_contains "$hy" "bogus_key:"

# ============================================================================
# 3. Host-blocked verifier: worker_model recorded, verifier_model absent (no
#    verifier ran) — the missing value, with status: blocked, is the evidence of
#    the block.
# ============================================================================
cc_fx_plan "$ws" 0002-blk "Blocked" "api"
cc_plan_approve "$ws" 0002-blk >/dev/null
bexec=$(cc_execution_begin "$ws" 0002-blk sess2 | sed -n 's/^execution_id: //p')
bedir="$ws/.runtime/executions/0002-blk/$bexec"
cc_attempt_begin "$bedir" >/dev/null
cc_fx_commit "$ws" 0002-blk api impl
cc_worker_commit_record "$bedir" api implementation >/dev/null
cc_attempt_evidence_record "$bedir" 1 worker_model=model-hi worker_effort=high >/dev/null
cc_verifier_result_record "$bedir" 1 blocked >/dev/null
assert_eq "blocked" "$(cc_execution_status "$bedir")"
bhy="$bedir/attempts/001/host-evidence.yaml"
contains "$bhy" "worker_model: model-hi"
not_contains "$bhy" "verifier_model:"                 # the block's evidence is the absence

# a BASE_UNBUILDABLE blocked trace also records the deterministic phases that ran.
# two predecessors write DIFFERENT content to the same path -> the runtime-authored
# integration merge conflicts -> blocked (not a worker failure).
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
cc_fx_plan_ex "$ws" 0003-ca "CA" api src/shared ""
cc_fx_plan_ex "$ws" 0004-cb "CB" api src/shared ""
cc_fx_plan_ex "$ws" 0005-cc "CC" api src/cc "0003-ca 0004-cb"
cc_plan_approve "$ws" 0003-ca >/dev/null
cc_plan_approve "$ws" 0004-cb >/dev/null
run_content 0003-ca api src/shared "AAA"
run_content 0004-cb api src/shared "BBB"
cc_plan_approve "$ws" 0005-cc >/dev/null
expect_failure cc_execution_begin "$ws" 0005-cc 0005-cc-w
ce="$(edir_of 0005-cc)"
assert_eq "blocked" "$(cc_execution_status "$ce")"
contains "$ce/execution.yaml" "blocked_reason: BASE_UNBUILDABLE"

# ============================================================================
# 4. Tiering is evidence-only: escalation across a repair changes the recorded
#    (model, effort) but never the worker-failure counter (INV-REPAIR-01).
# ============================================================================
cc_fx_plan "$ws" 0006-esc "Escalate" "api"
cc_plan_approve "$ws" 0006-esc >/dev/null
eex=$(cc_execution_begin "$ws" 0006-esc sess6 | sed -n 's/^execution_id: //p')
eed="$ws/.runtime/executions/0006-esc/$eex"
# attempt 1 at the configured start -> rejected
cc_attempt_begin "$eed" >/dev/null
cc_fx_commit "$ws" 0006-esc api a1
cc_worker_commit_record "$eed" api implementation >/dev/null
cc_attempt_evidence_record "$eed" 1 worker_model=model-lo worker_effort=high >/dev/null
cc_verifier_result_record "$eed" 1 failed >/dev/null
assert_eq "1" "$(cc_scalar "$eed/execution.yaml" worker_failures)"
# attempt 2 escalated above the start -> passes; escalation did NOT buy an attempt
cc_attempt_begin "$eed" >/dev/null
cc_fx_commit "$ws" 0006-esc api a2
cc_worker_commit_record "$eed" api repair >/dev/null
cc_attempt_evidence_record "$eed" 2 worker_model=model-hi worker_effort=max >/dev/null
cc_verifier_result_record "$eed" 2 passed >/dev/null
assert_eq "verified" "$(cc_execution_status "$eed")"
# the counter is exactly the one rejection — escalation changed the model, not the accounting
assert_eq "1" "$(cc_scalar "$eed/execution.yaml" worker_failures)"
assert_eq "2" "$(cc_scalar "$eed/execution.yaml" current_attempt)"
contains "$eed/attempts/001/host-evidence.yaml" "worker_model: model-lo"
contains "$eed/attempts/002/host-evidence.yaml" "worker_model: model-hi"
assert_eq "2" "$(cc_scalar "$eed/execution.yaml" current_attempt)"   # escalation via worker_model, not a boolean flag

# ============================================================================
# 6. Concurrent run-stack: the safety is already built — two independent plans are
#    both ready (overlap permitted); a held lease serializes a conflicting pair
#    (the lease is the arbiter). No new invariant; only a coordinator fan-out.
# ============================================================================
cc_fx_plan_ex "$ws" 0009-ia "IndepA" api src/ia ""
cc_fx_plan_ex "$ws" 0010-ib "IndepB" api src/ib ""
cc_plan_approve "$ws" 0009-ia >/dev/null
cc_plan_approve "$ws" 0010-ib >/dev/null
part=$(cc_run_stack_ready "$ws" 0009-ia 0010-ib)
printf '%s\n' "$part" | grep -q '^0009-ia: ready' || fail "independent 0009 not ready (overlap must be permitted)"
printf '%s\n' "$part" | grep -q '^0010-ib: ready' || fail "independent 0010 not ready (overlap must be permitted)"
# a held lease on a shared region makes a conflicting peer wait — the arbiter
cc_fx_plan_ex "$ws" 0011-la "LeaseA" api src/lz ""
cc_fx_plan_ex "$ws" 0012-lb "LeaseB" api src/lz ""
cc_plan_approve "$ws" 0011-la >/dev/null
cc_plan_approve "$ws" 0012-lb >/dev/null
cc_lease_acquire "$ws" api 0011-la "src/lz" >/dev/null
# the loser of the race gets LEASE_CONFLICT and stays waiting
expect_failure cc_lease_acquire "$ws" api 0012-lb "src/lz"
lpart=$(cc_run_stack_ready "$ws" 0011-la 0012-lb)
printf '%s\n' "$lpart" | grep -q '^0012-lb: waiting' || fail "conflicting peer must serialize (waiting)"
cc_lease_release "$ws" api 0011-la >/dev/null

pass 'execution latency'
