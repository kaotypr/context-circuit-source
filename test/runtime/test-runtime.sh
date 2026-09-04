#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

cc_fx_repo "$ws" api development
cc_fx_plan "$ws" 0001-alpha "Alpha" "api"

# --- plan-validate rejects an undeclared task repository ---
mkdir -p "$ws/plans/9999-bad/tasks"
cat >"$ws/plans/9999-bad/plan.yaml" <<'EOF'
schema_version: 3
plan: 9999-bad
title: Bad
status: draft
objective: bad
intent: i9999-bad
repositories:
  - id: api
tasks:
  - id: X-001
    title: t
    repositories: [ghost]
    paths: [src]
    depends_on: []
EOF
printf '# Bad\n' >"$ws/plans/9999-bad/PLAN.md"
expect_failure cc_plan_validate "$ws/plans/9999-bad"
rm -rf "$ws/plans/9999-bad"

# --- execute on the approved intent authorization: one branch + worktree per repo ---
exec=$(cc_execution_begin "$ws" 0001-alpha sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-alpha "$exec")
require_dir "$ws/.runtime/worktrees/0001-alpha/api"
require_file "$edir/repositories/api.yaml"
contains "$edir/repositories/api.yaml" "branch: cc/0001-alpha/api"
require_file "$edir/snapshot/plan.yaml"

# --- dirty base checkout blocks a fresh execution ---
printf 'dirty\n' >"$ws/repositories/api/src/dirty.txt"
git -C "$ws/repositories/api" add -A
# leave uncommitted -> dirty
cc_fx_plan "$ws" 0002-beta "Beta" "api"
expect_failure cc_execution_begin "$ws" 0002-beta sess2
git -C "$ws/repositories/api" reset -q --hard HEAD
git -C "$ws/repositories/api" clean -fdq

# --- worker commit recorded before verification ---
cc_attempt_begin "$edir" >/dev/null
cc_fx_commit "$ws" 0001-alpha api impl
cc_worker_commit_record "$edir" api implementation >/dev/null
assert_eq "verifying" "$(cc_execution_status "$edir")"
contains "$edir/repositories/api.yaml" "latest_commit: "

# --- verifier write is rejected; verification does not change plan status ---
cc_verifier_prepare "$edir" >/dev/null
expect_failure cc_verifier_result_record "$edir" 1 passed --wrote-products
cc_verifier_result_record "$edir" 1 failed >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" 0001-alpha)"

# --- repair creates a new commit record ---
cc_attempt_begin "$edir" >/dev/null
cc_fx_commit "$ws" 0001-alpha api repair
cc_worker_commit_record "$edir" api repair >/dev/null
require_file "$edir/attempts/002/worker.yaml"
cc_verifier_prepare "$edir" >/dev/null

# --- third worker failure stops execution ---
cc_verifier_result_record "$edir" 2 failed >/dev/null
cc_attempt_begin "$edir" >/dev/null
cc_fx_commit "$ws" 0001-alpha api repair2
cc_worker_commit_record "$edir" api repair >/dev/null
cc_verifier_prepare "$edir" >/dev/null
out=$(cc_verifier_result_record "$edir" 3 failed)
printf '%s\n' "$out" | grep -Fq "stop: FAILURE_LIMIT_REACHED" || fail "expected FAILURE_LIMIT_REACHED"
assert_eq "failed" "$(cc_execution_status "$edir")"
expect_failure cc_repair_allowed "$edir"

# --- completion requires passed evidence (failed execution cannot complete) ---
expect_failure cc_plan_complete "$ws" 0001-alpha
assert_eq "draft" "$(cc_plan_status "$ws" 0001-alpha)"

# --- successful path: completion records commits + preserves PK boundary ---
cc_fx_repo "$ws" web development
cc_fx_plan "$ws" 0003-gamma "Gamma" "web"
exec3=$(cc_execution_begin "$ws" 0003-gamma sess3 | sed -n 's/^execution_id: //p')
edir3=$(cc_fx_exec_dir "$ws" 0003-gamma "$exec3")
cc_attempt_begin "$edir3" >/dev/null
cc_fx_commit "$ws" 0003-gamma web impl
cc_worker_commit_record "$edir3" web implementation >/dev/null
cc_verifier_prepare "$edir3" >/dev/null
cc_verifier_result_record "$edir3" 1 passed >/dev/null
assert_eq "verified" "$(cc_execution_status "$edir3")"
assert_eq "draft" "$(cc_plan_status "$ws" 0003-gamma)"   # verified execution != done plan
expect_failure cc_completion_ready "$ws" 0003-gamma
cc_human_acceptance_record "$edir3" alice >/dev/null
cc_completion_ready "$ws" 0003-gamma >/dev/null
expect_failure cc_plan_complete "$ws" 0003-gamma
cc_delivery_record "$ws" 0003-gamma >/dev/null
cc_completion_infer "$ws" 0003-gamma >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0003-gamma)"
require_file "$edir3/completion.yaml"
contains "$edir3/completion.yaml" "human_completion: inferred"
# runtime writes no Product Knowledge
not_contains "$ws/context/PROJECT.md" "web"

# --- archive preserves files/status, removes index; no status validation ---
prevstatus=$(cc_plan_status "$ws" 0003-gamma)
cc_plan_archive "$ws" 0003-gamma >/dev/null
require_dir "$ws/plans/archive/0003-gamma"
require_file "$ws/plans/archive/0003-gamma/plan.yaml"
not_contains "$ws/plans/INDEX.md" "| 0003-gamma |"
assert_eq "$prevstatus" "$(cc_scalar "$ws/plans/archive/0003-gamma/plan.yaml" status)"

# --- archive collision leaves state unchanged ---
cc_fx_plan "$ws" 0004-delta "Delta" "web"
mkdir -p "$ws/plans/archive/0004-delta"
expect_failure cc_plan_archive "$ws" 0004-delta
require_dir "$ws/plans/0004-delta"

# --- restore returns the plan and index row ---
cc_plan_restore "$ws" 0003-gamma >/dev/null
require_dir "$ws/plans/0003-gamma"
contains "$ws/plans/INDEX.md" "| 0003-gamma |"

# --- interruption preserves execution record and evidence ---
cc_recovery_inspect "$edir" >/dev/null
require_file "$edir/execution.yaml"
require_dir "$ws/.runtime/worktrees/0001-alpha/api"

# --- resume eligibility: mid-flight + matching owner is eligible; a terminal or
#     foreign-owner request is not ---
cc_fx_plan "$ws" 0011-resume "Resume" "api"
rexec=$(cc_execution_begin "$ws" 0011-resume owner-a | sed -n 's/^execution_id: //p')
redir=$(cc_fx_exec_dir "$ws" 0011-resume "$rexec")
cc_attempt_begin "$redir" >/dev/null
ri=$(cc_recovery_inspect "$redir" owner-a)
printf '%s\n' "$ri" | grep -Fq "resume_eligible: true" || fail "owner should resume"
rj=$(cc_recovery_inspect "$redir" owner-b)
printf '%s\n' "$rj" | grep -Fq "resume_eligible: false" || fail "foreign owner must not resume"
printf '%s\n' "$rj" | grep -Fq "resume_reason: OWNER_MISMATCH" || fail "expected OWNER_MISMATCH"
# the failed execution above is terminal and not resume-eligible
rt=$(cc_recovery_inspect "$edir" sess1)
printf '%s\n' "$rt" | grep -Fq "resume_eligible: false" || fail "failed execution not resumable"

# --- provider-specific launch and old-design machinery absent from the runtime ---
not_contains "$ROOT/wrapper/runtime/engine.sh" "cc_probe"
not_contains "$ROOT/wrapper/runtime/engine.sh" "cc_route"
not_contains "$ROOT/wrapper/runtime/engine.sh" "cc_confirmation_card"
not_contains "$ROOT/wrapper/runtime/engine.sh" "cc_context_packet"
not_contains "$ROOT/wrapper/runtime/engine.sh" "codex exec"
not_contains "$ROOT/wrapper/runtime/engine.sh" "claude -p"

pass 'runtime boundary'
