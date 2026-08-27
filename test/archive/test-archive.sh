#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development

# --- archive without any status validation: works for draft, approved, done ---
cc_fx_plan "$ws" 0001-draft "Draft plan" "api"
cc_plan_archive "$ws" 0001-draft >/dev/null
require_dir "$ws/plans/archive/0001-draft"
assert_eq "draft" "$(cc_scalar "$ws/plans/archive/0001-draft/plan.yaml" status)"
not_contains "$ws/plans/INDEX.md" "| 0001-draft |"

cc_fx_plan "$ws" 0002-appr "Approved plan" "api"
cc_plan_approve "$ws" 0002-appr >/dev/null
cc_plan_archive "$ws" 0002-appr >/dev/null
assert_eq "approved" "$(cc_scalar "$ws/plans/archive/0002-appr/plan.yaml" status)"

# --- archiving does not stop or rename an active execution ---
cc_fx_plan "$ws" 0003-live "Live" "api"
cc_plan_approve "$ws" 0003-live >/dev/null
exec=$(cc_execution_begin "$ws" 0003-live sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0003-live "$exec")
cc_attempt_begin "$edir" >/dev/null
cc_fx_commit "$ws" 0003-live api impl
cc_worker_commit_record "$edir" api implementation >/dev/null
cc_plan_archive "$ws" 0003-live >/dev/null
# execution continues from its immutable snapshot; runtime evidence preserved
require_file "$edir/execution.yaml"
require_file "$edir/snapshot/plan.yaml"
require_dir "$ws/.runtime/worktrees/0003-live/api"
cc_recovery_inspect "$edir" >/dev/null
# the archived plan directory is not touched by the ongoing execution
require_dir "$ws/plans/archive/0003-live"

# --- restore returns files, status, and the index row ---
cc_plan_restore "$ws" 0003-live >/dev/null
require_dir "$ws/plans/0003-live"
contains "$ws/plans/INDEX.md" "| 0003-live |"

# --- archive/restore collision safety leaves everything unchanged ---
mkdir -p "$ws/plans/archive/0001-draft-collide"
cc_fx_plan "$ws" 0001-draft-collide "Collide" "api"
# pre-place a collision at the archive target
mv "$ws/plans/archive/0001-draft-collide" "$ws/plans/archive/0001-draft-collide.keep" 2>/dev/null || true
mkdir -p "$ws/plans/archive/0001-draft-collide"
printf 'x\n' >"$ws/plans/archive/0001-draft-collide/marker"
expect_failure cc_plan_archive "$ws" 0001-draft-collide
require_dir "$ws/plans/0001-draft-collide"
contains "$ws/plans/INDEX.md" "| 0001-draft-collide |"

# --- restore of a missing archived plan fails cleanly ---
expect_failure cc_plan_restore "$ws" 9999-nope

pass 'archive'
