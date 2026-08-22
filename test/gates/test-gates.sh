#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-gates.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
card=$(cc_confirmation_card archive-plan plans/app-plans/checkout approved 'archive sidecar append' 'plan status, Git, delivery' none)
printf '%s\n' "$card" | grep -F 'Confirmation requested:' >/dev/null || fail 'confirmation card missing gate'
printf '%s\n' "$card" | grep -F 'Will not change:' >/dev/null || fail 'confirmation card missing non-effects'
printf '%s\n' "$card" | grep -F 'Action: archive-plan' >/dev/null || fail 'generic card used an approval action'
contains "$ROOT/docs/gates.md" 'Action: present-approval-card'
contains "$ROOT/docs/gates.md" 'nothing has changed yet'
contains "$ROOT/docs/gates.md" 'Confirm approval of plan <id>.'
contains "$ROOT/docs/gates.md" 'does not start Run approved plan'
contains "$ROOT/docs/gates.md" 'Action: commit-approved-plan'
contains "$ROOT/docs/gates.md" 'Confirm commit of the approved plan state.'
not_contains "$ROOT/wrapper/runtime/engine.sh" 'cc_approval_card'
not_contains "$ROOT/wrapper/runtime/engine.sh" 'cc_commit_approved_plan_card'

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
