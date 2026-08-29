#!/bin/sh
# Execution-latency semantics (v0.7.0): the two-observer measurement model
# (engine-stamped deterministic phase_ms + per-attempt timing boundaries;
# coordinator-recorded inference wall-clock and (model, effort) as bounded host
# evidence), per-role tiering that never touches the failure counter, the optional
# plan complexity hint, and concurrent-run-stack overlap arbitrated by the lease.
# These prove the additive fields exist and are monotonic, that escalation is
# evidence-only, and that independent plans may overlap while conflicting ones
# serialize — deterministically, because a live model cannot be timed on command.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development

edir_of() { printf '%s/.runtime/executions/%s/%s' "$ws" "$1" "$(cc_latest_execution "$ws" "$1")"; }

# ============================================================================
# 1. Measurement: engine-stamped deterministic phases on a normal execution.
# ============================================================================
cc_fx_plan "$ws" 0001-meas "Measure" "api"
cc_plan_approve "$ws" 0001-meas >/dev/null
exec=$(cc_execution_begin "$ws" 0001-meas sess1 | sed -n 's/^execution_id: //p')
edir="$ws/.runtime/executions/0001-meas/$exec"
ey="$edir/execution.yaml"
# the phase_ms map exists with all four deterministic phases, each a non-negative integer
contains "$ey" "phase_ms:"
for ph in base_prepare grounding_discovery integration_merge verifier_prepare; do
	grep -Eq "^  $ph: [0-9]+$" "$ey" || fail "phase_ms.$ph missing or non-integer"
done

# per-attempt start boundary is stamped by attempt-begin
cc_attempt_begin "$edir" >/dev/null
wy="$edir/attempts/001/worker.yaml"
contains "$wy" "attempt_started_at:"
astart=$(cc_scalar "$wy" attempt_started_at)
[ -n "$astart" ] || fail "attempt_started_at empty"

cc_fx_commit "$ws" 0001-meas api impl
cc_worker_commit_record "$edir" api implementation >/dev/null
# verifier-prepare folds its (recurring) duration into phase_ms.verifier_prepare
cc_verifier_prepare "$edir" >/dev/null
grep -Eq "^  verifier_prepare: [0-9]+$" "$ey" || fail "verifier_prepare not stamped after verifier-prepare"

# coordinator records inference wall-clock + (model, effort) as bounded evidence
cc_attempt_evidence_record "$edir" 1 \
	worker_wall_s=12.4 worker_model=opus worker_effort=high \
	verifier_wall_s=3.1 verifier_model=sonnet verifier_effort=medium >/dev/null
hy="$edir/attempts/001/host-evidence.yaml"
for kv in "worker_wall_s: 12.4" "worker_model: opus" "verifier_wall_s: 3.1" "verifier_model: sonnet"; do
	contains "$hy" "$kv"
done
# attempt_count is current_attempt, surfaced as a timing dimension
assert_eq "1" "$(cc_scalar "$ey" current_attempt)"

# per-attempt end boundary is stamped by verifier-result-record, monotonic with start
cc_verifier_result_record "$edir" 1 passed >/dev/null
vy="$edir/attempts/001/verifier.yaml"
contains "$vy" "attempt_ended_at:"
aend=$(cc_scalar "$vy" attempt_ended_at)
[ -n "$aend" ] || fail "attempt_ended_at empty"
# ISO-8601 UTC sorts lexicographically == chronologically
{ [ "$aend" = "$astart" ] || [ "$aend" \> "$astart" ]; } || fail "attempt end precedes start (not monotonic)"

# the engine phases are a rounding error against inference — proven by shape, not
# magnitude: both are recorded, so the split is real and attributable.

# ============================================================================
# 2. Bounded evidence surface: unknown key and unsafe value are refused.
# ============================================================================
expect_failure cc_attempt_evidence_record "$edir" 1 bogus_key=x
expect_failure cc_attempt_evidence_record "$edir" 1 "worker_model=opus;rm"
expect_failure cc_attempt_evidence_record "$edir" 1 malformed_no_equals
# a refused call writes nothing new
not_contains "$hy" "bogus_key:"

# ============================================================================
# 3. Host-blocked verifier: worker_wall_s recorded, verifier_wall_s absent — the
#    missing value is the evidence of the block. Deterministic phases still present.
# ============================================================================
cc_fx_plan "$ws" 0002-blk "Blocked" "api"
cc_plan_approve "$ws" 0002-blk >/dev/null
bexec=$(cc_execution_begin "$ws" 0002-blk sess2 | sed -n 's/^execution_id: //p')
bedir="$ws/.runtime/executions/0002-blk/$bexec"
contains "$bedir/execution.yaml" "phase_ms:"          # phases present before any verify
cc_attempt_begin "$bedir" >/dev/null
cc_fx_commit "$ws" 0002-blk api impl
cc_worker_commit_record "$bedir" api implementation >/dev/null
cc_attempt_evidence_record "$bedir" 1 worker_wall_s=9.0 worker_model=opus >/dev/null
cc_verifier_result_record "$bedir" 1 blocked >/dev/null
assert_eq "blocked" "$(cc_execution_status "$bedir")"
bhy="$bedir/attempts/001/host-evidence.yaml"
contains "$bhy" "worker_wall_s: 9.0"
not_contains "$bhy" "verifier_wall_s:"                # the block's evidence is the absence

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
contains "$ce/execution.yaml" "phase_ms:"             # a blocked trace still records what ran

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
cc_attempt_evidence_record "$eed" 1 worker_model=sonnet worker_effort=high escalated=false >/dev/null
cc_verifier_result_record "$eed" 1 failed >/dev/null
assert_eq "1" "$(cc_scalar "$eed/execution.yaml" worker_failures)"
# attempt 2 escalated above the start -> passes; escalation did NOT buy an attempt
cc_attempt_begin "$eed" >/dev/null
cc_fx_commit "$ws" 0006-esc api a2
cc_worker_commit_record "$eed" api repair >/dev/null
cc_attempt_evidence_record "$eed" 2 worker_model=opus worker_effort=max escalated=true >/dev/null
cc_verifier_result_record "$eed" 2 passed >/dev/null
assert_eq "verified" "$(cc_execution_status "$eed")"
# the counter is exactly the one rejection — escalation changed the model, not the accounting
assert_eq "1" "$(cc_scalar "$eed/execution.yaml" worker_failures)"
assert_eq "2" "$(cc_scalar "$eed/execution.yaml" current_attempt)"
contains "$eed/attempts/001/host-evidence.yaml" "worker_model: sonnet"
contains "$eed/attempts/002/host-evidence.yaml" "worker_model: opus"
contains "$eed/attempts/002/host-evidence.yaml" "escalated: true"

# ============================================================================
# 5. Optional per-plan complexity hint: additive, validated, absent by default.
# ============================================================================
cc_fx_plan "$ws" 0007-cx "Complex" "api"
cc_plan_validate "$ws/plans/0007-cx" >/dev/null   # no complexity -> valid (default absent)
# inject a valid hint -> still valid
printf 'complexity: high\n' >>"$ws/plans/0007-cx/plan.yaml"
cc_plan_validate "$ws/plans/0007-cx" >/dev/null || fail "complexity: high must validate"
# an invalid value is refused
cc_fx_plan "$ws" 0008-cxbad "ComplexBad" "api"
printf 'complexity: enormous\n' >>"$ws/plans/0008-cxbad/plan.yaml"
expect_failure cc_plan_validate "$ws/plans/0008-cxbad"

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
