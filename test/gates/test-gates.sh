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
contains "$ROOT/docs/gates.md" 'Action: identity-acceptance'
contains "$ROOT/docs/gates.md" 'Proposed defaults:'
contains "$ROOT/docs/gates.md" 'Immediate effects: workspace.accept_identity'
contains "$ROOT/docs/gates.md" 'Action: repository-registration'
contains "$ROOT/docs/gates.md" 'Immediate effects: workspace.register_repository'
contains "$ROOT/docs/gates.md" 'Later authorized effects: repository-bootstrap; repository-create-empty (reserved, not activated); execute-plan'
contains "$ROOT/docs/gates.md" 'Immediate effects: repository-bootstrap'
contains "$ROOT/docs/gates.md" 'Later authorized effects: git.commit'
contains "$ROOT/docs/gates.md" 'Immediate effects: git.commit'
contains "$ROOT/docs/gates.md" 'Later authorized effects: delivery.push'
contains "$ROOT/docs/gates.md" 'They never grant authorization'

fx=$(mktemp -d "${TMPDIR:-/tmp}/cc-gates-identity.XXXXXX")
trap 'rm -rf "$runtime" "$fx"' EXIT HUP INT TERM
cp -R "$ROOT/test/contracts/fixtures/identity-projection/matching/." "$fx/"
id_card=$(cc_identity_acceptance_card "$fx")
printf '%s\n' "$id_card" | grep -F 'mode: solo' >/dev/null || fail 'identity-acceptance card omitted mode default'
printf '%s\n' "$id_card" | grep -F 'Immediate effects: workspace.accept_identity' >/dev/null || fail 'identity-acceptance card omitted immediate effects'
reg_card=$(cc_repository_registration_card "$fx" app '' '')
printf '%s\n' "$reg_card" | grep -F 'canonical URL: none' >/dev/null || fail 'repository-registration card omitted URL default'
printf '%s\n' "$reg_card" | grep -F 'default branch: main' >/dev/null || fail 'repository-registration card omitted branch default'
printf '%s\n' "$reg_card" | grep -F 'Immediate effects: workspace.register_repository' >/dev/null || fail 'repository-registration card omitted immediate effects'
cc_accept_identity "$fx" confirmed '' '' '' ''
contains "$fx/workspace.yaml" 'status: accepted'
contains "$fx/workspace.yaml" 'mode: solo'
contains "$fx/workspace.yaml" 'roles: none'
contains "$fx/context/WORKSPACE.md" 'status: accepted'
contains "$fx/context/WORKSPACE.md" 'Preserve this authored Product Knowledge in WORKSPACE.md.'
invented=$(cc_accept_identity "$fx" confirmed solo none none main owner=alice 2>&1 || true)
printf '%s\n' "$invented" | grep -F FIELD_INVENTED >/dev/null || fail 'post-confirmation invented field was accepted'
cc_register_repository "$fx" app 'https://github.com/acme/app.git' main confirmed
contains "$fx/workspace.yaml" 'canonical_url: https://github.com/acme/app.git'
contains "$fx/context/PROJECT.md" 'canonical_url: https://github.com/acme/app.git'
contains "$fx/context/INDEX.md" '  app:'
not_contains "$fx/workspace.yaml" 'path:'
test ! -e "$fx/repositories.local.yaml" || fail 'registration created local bindings'
assert_eq "$(cc_validate_identity_projection "$fx")" identity-projection-ok
pass 'confirmation cards, append-only archive, offline fallback, and cleanup gate'
