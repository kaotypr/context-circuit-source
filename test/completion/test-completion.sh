#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

cc_fx_plan "$ws" 0001-billing "Billing" "api web"
exec=$(cc_execution_begin "$ws" 0001-billing sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-billing "$exec")
cc_attempt_begin "$edir" >/dev/null
cc_fx_commit "$ws" 0001-billing api impl
cc_worker_commit_record "$edir" api implementation >/dev/null
cc_fx_commit "$ws" 0001-billing web impl
cc_worker_commit_record "$edir" web implementation >/dev/null
cc_verifier_prepare "$edir" >/dev/null

# --- completion is blocked before a verifier pass ---
expect_failure cc_completion_ready "$ws" 0001-billing
expect_failure cc_plan_complete "$ws" 0001-billing
assert_eq "draft" "$(cc_plan_status "$ws" 0001-billing)"

# --- a verifier pass produces verified evidence but not completion ---
cc_verifier_result_record "$edir" 1 passed >/dev/null
expect_failure cc_completion_ready "$ws" 0001-billing
cc_human_acceptance_record "$edir" alice >/dev/null
cc_completion_ready "$ws" 0001-billing >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" 0001-billing)"   # verified execution, not yet done

# --- Standard completion is inferred from the accepted candidate plus delivery ---
expect_failure cc_plan_complete "$ws" 0001-billing
cc_delivery_record "$ws" 0001-billing >/dev/null
cc_completion_infer "$ws" 0001-billing >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" 0001-billing)"
require_file "$edir/completion.yaml"
contains "$edir/completion.yaml" "human_completion: inferred"
api_latest=$(cc_scalar "$edir/repositories/api.yaml" latest_commit)
contains "$edir/completion.yaml" "api: $api_latest"

# --- knowledge reconciliation references are recorded without touching PK ---
project_before=$(cc_digest "$ws/context/PROJECT.md")
cat >"$ws/.tmp-impact.yaml" <<EOF
schema_version: 1
execution_id: $exec
plan: 0001-billing
status: review-needed
proposals: [context-impact-0001-billing-001]
EOF
cc_context_impact_record "$edir" "$ws/.tmp-impact.yaml" >/dev/null
require_file "$edir/context-impact.yaml"
contains "$edir/context-impact.yaml" "status: review-needed"
# Product Knowledge is unchanged by completion or impact recording
assert_eq "$project_before" "$(cc_digest "$ws/context/PROJECT.md")"

# --- a failed execution can never complete ---
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
expect_failure cc_plan_complete "$ws" 0002-fail

pass 'completion'
