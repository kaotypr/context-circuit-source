#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# --- completion-ready stays a query; unverified mark-done still flips status ---
cc_fx_plan "$ws" 0001-unverified "Unverified" "api"
exec_u=$(cc_execution_begin "$ws" 0001-unverified sess-u | sed -n 's/^execution_id: //p')
edir_u=$(cc_fx_exec_dir "$ws" 0001-unverified "$exec_u")
cc_attempt_begin "$edir_u" >/dev/null
cc_fx_commit "$ws" 0001-unverified api impl
cc_worker_commit_record "$edir_u" api implementation >/dev/null
cc_verifier_prepare "$edir_u" >/dev/null
expect_failure cc_completion_ready "$ws" 0001-unverified
cc_plan_complete "$ws" 0001-unverified >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0001-unverified)"

# --- a verifier pass produces verified evidence but not completion ---
cc_fx_plan "$ws" 0001-billing "Billing" "api"
exec=$(cc_execution_begin "$ws" 0001-billing sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-billing "$exec")
cc_attempt_begin "$edir" >/dev/null
cc_fx_commit "$ws" 0001-billing api impl
cc_worker_commit_record "$edir" api implementation >/dev/null
cc_verifier_prepare "$edir" >/dev/null
cc_verifier_result_record "$edir" 1 passed >/dev/null
expect_failure cc_completion_ready "$ws" 0001-billing
cc_human_acceptance_record "$edir" alice >/dev/null
cc_completion_ready "$ws" 0001-billing >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" 0001-billing)"   # verified execution, not yet done

# --- Standard completion is an explicit mark-done (delivery is not required) ---
cc_plan_complete "$ws" 0001-billing >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0001-billing)"
require_file "$edir/completion.yaml"
contains "$edir/completion.yaml" "human_completion: accepted"
api_latest=$(cc_scalar "$edir/repositories/api.yaml" latest_commit)
contains "$edir/completion.yaml" "api: $api_latest"

# --- knowledge impact refs are recorded without touching PK ---
project_before=$(cc_digest "$ws/context/PROJECT.md")
cat >"$ws/.tmp-impact.yaml" <<EOF
schema_version: 1
execution_id: $exec
plan: 0001-billing
status: no-update-needed
units: []
EOF
cc_context_impact_record "$edir" "$ws/.tmp-impact.yaml" >/dev/null
require_file "$edir/context-impact.yaml"
contains "$edir/context-impact.yaml" "status: no-update-needed"
# Product Knowledge is unchanged by completion or impact recording
assert_eq "$project_before" "$(cc_digest "$ws/context/PROJECT.md")"

# --- a failed execution still becomes done when asked ---
cc_fx_plan "$ws" 0002-fail "Fail" "api"
exec2=$(cc_execution_begin "$ws" 0002-fail sess2 | sed -n 's/^execution_id: //p')
edir2=$(cc_fx_exec_dir "$ws" 0002-fail "$exec2")
i=1
while [ "$i" -le 3 ]; do
	cc_attempt_begin "$edir2" >/dev/null
	cc_fx_commit "$ws" 0002-fail api "try$i"
	cc_worker_commit_record "$edir2" api implementation >/dev/null
	cc_verifier_prepare "$edir2" >/dev/null
	cc_verifier_result_record "$edir2" "$i" failed >/dev/null || true
	i=$((i + 1))
done
assert_eq "failed" "$(cc_execution_status "$edir2")"
cc_plan_complete "$ws" 0002-fail >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0002-fail)"

# --- a never-built draft (no execution) still becomes done when asked ---
cc_fx_plan "$ws" 0004-never "Never" "api"
cc_plan_complete "$ws" 0004-never >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0004-never)"

# --- verified with no human-acceptance still becomes done when asked ---
cc_fx_plan "$ws" 0005-noacc "NoAcc" "api"
exec5=$(cc_execution_begin "$ws" 0005-noacc sess5 | sed -n 's/^execution_id: //p')
edir5=$(cc_fx_exec_dir "$ws" 0005-noacc "$exec5")
cc_attempt_begin "$edir5" >/dev/null
cc_fx_commit "$ws" 0005-noacc api impl
cc_worker_commit_record "$edir5" api implementation >/dev/null
cc_verifier_prepare "$edir5" >/dev/null
cc_verifier_result_record "$edir5" 1 passed >/dev/null
expect_failure cc_completion_ready "$ws" 0005-noacc
cc_plan_complete "$ws" 0005-noacc >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0005-noacc)"

# --- several named plans: the same plan-complete path, once per id ---
cc_fx_plan "$ws" 0006-one "One" "api"
cc_fx_plan "$ws" 0007-two "Two" "web"
cc_plan_complete "$ws" 0006-one >/dev/null
cc_plan_complete "$ws" 0007-two >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0006-one)"
assert_eq "done" "$(cc_plan_status "$ws" 0007-two)"

pass 'completion'
