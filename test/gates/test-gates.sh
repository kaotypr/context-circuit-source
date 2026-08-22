#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-gates.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
card=$(cc_confirmation_card approve-plan plans/app-plans/checkout draft 'plan status draft→approved; tasks draft→ready' 'execution, Git, delivery, publication' none)
printf '%s\n' "$card" | grep -F 'Confirmation requested:' >/dev/null || fail 'confirmation card missing gate'
printf '%s\n' "$card" | grep -F 'Will not change:' >/dev/null || fail 'confirmation card missing non-effects'
approval_card=$(cc_approval_card checkout-validation)
printf '%s\n' "$approval_card" | grep -F 'Action: present-approval-card' >/dev/null || fail 'approval card used the mutation action'
printf '%s\n' "$approval_card" | grep -F 'nothing has changed yet' >/dev/null || fail 'approval card omitted pre-confirmation state'
printf '%s\n' "$approval_card" | grep -F 'Confirm approval of plan checkout-validation.' >/dev/null || fail 'approval card omitted exact confirmation'
printf '%s\n' "$approval_card" | grep -F 'does not start Run approved plan' >/dev/null || fail 'approval card implied execution'
commit_card=$(cc_commit_approved_plan_card checkout-validation)
printf '%s\n' "$commit_card" | grep -F 'Action: commit-approved-plan' >/dev/null || fail 'commit card missing action'
printf '%s\n' "$commit_card" | grep -F 'Confirm commit of the approved plan state.' >/dev/null || fail 'commit card omitted exact confirmation'
contains "$ROOT/docs/gates.md" 'Confirm approval of plan <id>.'
contains "$ROOT/docs/gates.md" 'Confirm commit of the approved plan state.'

cat > "$runtime/plan.yaml" <<'EOF'
schema_version: 2
id: archive-fixture
status: approved
EOF
cc_archive_event "$runtime/archive.yaml" archive-fixture archived root replaced
contains "$runtime/archive.yaml" 'event: archived'
contains "$runtime/plan.yaml" 'status: approved'
cc_archive_event "$runtime/archive.yaml" archive-fixture restored root restored
contains "$runtime/archive.yaml" 'event: restored'
test "$(grep -c '^  - event:' "$runtime/archive.yaml")" -eq 2 || fail 'archive history was overwritten'

assert_cleanup=$(cc_cleanup_disposition dirty unpushed)
assert_eq "$assert_cleanup" CLEANUP_REQUIRES_DISCARD_CONFIRMATION
assert_clean=$(cc_cleanup_disposition clean clean)
assert_eq "$assert_clean" CLEANUP_REQUIRES_HUMAN_CONFIRMATION
cc_select_sources sources/brief.md sources/constraints.md >/dev/null
expect_failure cc_select_sources sources
assert_eq "$(cc_provider_fallback unavailable)" offline-fallback
assert_eq "$(cc_provider_fallback denied)" provider-denied-manual-fallback
contains "$ROOT/docs/gates.md" 'separate explicit gates'
contains "$ROOT/docs/integrations.md" 'offline'
pass 'confirmation cards, append-only archive, offline fallback, and cleanup gate'
