#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

fixture_count=0
while IFS='	' read -r request expected; do
  test -n "$request" || continue
  fixture_count=$((fixture_count + 1))
  route=$(cc_route "$request")
  actual=$(printf '%s\n' "$route" | sed -n 's/^capability: //p')
  assert_eq "$actual" "$expected"
  context_set=$(printf '%s\n' "$route" | sed -n 's/^context_set: //p')
  test -n "$context_set" || fail "route emitted no context set: $request"
  cc_context_set_registered "$ROOT" "$context_set" || fail "route emitted unregistered context set: $context_set"
done <<EOF
$(awk '
  /^    request: / { request=substr($0,14) }
  /^    action: / { print request "\t" $2 }
' "$ROOT/test/routing/fixtures.yaml")
EOF
test "$fixture_count" -ge 50 || fail "routing fixture count: $fixture_count"

while IFS= read -r action; do
  test -n "$action" || continue
  mapped=$(cc_route_context_set "$action" "$ROOT") || fail "action was not mapped: $action"
  cc_context_set_registered "$ROOT" "$mapped" || fail "action mapped to unregistered set: $action -> $mapped"
done <<EOF
$(awk '
  /^  actions:$/ { in_actions=1; next }
  in_actions && /^[^[:space:]]/ { exit }
  in_actions && /^  - / { print $2 }
' "$ROOT/wrapper/contracts/routes.yaml")
EOF

workspace_fixture=$(mktemp -d "${TMPDIR:-/tmp}/cc-uninitialized-route.XXXXXX")
trap 'rm -rf "$workspace_fixture"' EXIT HUP INT TERM
mkdir -p "$workspace_fixture/context" "$workspace_fixture/wrapper/contracts"
cp "$ROOT/template/workspace.yaml" "$workspace_fixture/workspace.yaml"
cp "$ROOT/template/context/WORKSPACE.md" "$workspace_fixture/context/WORKSPACE.md"
cp "$ROOT/template/context/PROJECT.md" "$workspace_fixture/context/PROJECT.md"
cp "$ROOT/template/context/INDEX.md" "$workspace_fixture/context/INDEX.md"
cp "$ROOT/wrapper/contracts/routes.yaml" "$workspace_fixture/wrapper/contracts/routes.yaml"
cp "$ROOT/wrapper/contracts/context-sets.yaml" "$workspace_fixture/wrapper/contracts/context-sets.yaml"
uninitialized_build=$(cc_route 'Can you help me build this?' "$workspace_fixture")
printf '%s\n' "$uninitialized_build" | grep -F 'capability: initialize' >/dev/null || fail 'uninitialized build request was not routed to initialization'
printf '%s\n' "$uninitialized_build" | grep -F 'human_gate: identity-acceptance' >/dev/null || fail 'uninitialized build request missed identity gate'
printf '%s\n' "$uninitialized_build" | grep -F 'authorization: confirmed-gate-required' >/dev/null || fail 'uninitialized build request was treated as authorized'
uninitialized_read=$(cc_route 'What is this workspace?' "$workspace_fixture")
printf '%s\n' "$uninitialized_read" | grep -F 'capability: orient' >/dev/null || fail 'uninitialized read-only orientation was blocked'
printf '%s\n' "$uninitialized_read" | grep -F 'authorization: read-only' >/dev/null || fail 'uninitialized read-only orientation was not read-only'
cc_accept_identity "$workspace_fixture" confirmed
accepted_build=$(cc_route 'Create a plan from this request.' "$workspace_fixture")
printf '%s\n' "$accepted_build" | grep -F 'capability: draft-plan' >/dev/null || fail 'accepted workspace did not retain plan route'
generic_build=$(cc_route 'I want a command-line to-do app. Can you help me build this?' "$workspace_fixture")
printf '%s\n' "$generic_build" | grep -F 'phase: plan-draft' >/dev/null || fail 'generic accepted build did not route to plan drafting'
printf '%s\n' "$generic_build" | grep -F 'capability: draft-plan' >/dev/null || fail 'generic accepted build was not constrained to draft-plan'

bootstrap_card=$(cc_route 'Bootstrap repository app.')
printf '%s\n' "$bootstrap_card" | grep -F 'capability: present-repository-bootstrap-card' >/dev/null || fail 'bootstrap request did not present a confirmation card'
printf '%s\n' "$bootstrap_card" | grep -F 'human_gate: repository-bootstrap' >/dev/null || fail 'bootstrap request missed repository gate'
bootstrap_confirmed=$(cc_route 'Confirm repository bootstrap for app.')
printf '%s\n' "$bootstrap_confirmed" | grep -F 'capability: repository-bootstrap' >/dev/null || fail 'confirmed bootstrap request was not authorized'

approval_card=$(cc_route 'Approve plan checkout-validation.')
printf '%s\n' "$approval_card" | grep -F 'capability: present-approval-card' >/dev/null || fail 'initial approval request did not present a card'
printf '%s\n' "$approval_card" | grep -F 'authorization: confirmed-gate-required' >/dev/null || fail 'initial approval request was treated as authorized'
printf '%s\n' "$approval_card" | grep -F 'capability: approve-plan' >/dev/null && fail 'initial approval request mutated into approve-plan'
approval_confirmed=$(cc_route 'Confirm approval of plan checkout-validation.')
printf '%s\n' "$approval_confirmed" | grep -F 'capability: approve-plan' >/dev/null || fail 'exact approval confirmation was not authorized'
printf '%s\n' "$approval_confirmed" | grep -F 'authorization: confirmed-gate' >/dev/null || fail 'exact approval confirmation missed confirmed-gate'
printf '%s\n' "$approval_confirmed" | grep -F 'capability: execute-plan' >/dev/null && fail 'approval confirmation was routed to execution'
printf '%s\n' "$approval_confirmed" | grep -F 'capability: commit-approved-plan' >/dev/null && fail 'approval confirmation was folded into commit'

draft=$(cc_route 'Run draft plan checkout-validation.')
printf '%s\n' "$draft" | grep -F 'eligibility: blocked' >/dev/null || fail 'draft execution was not blocked'
printf '%s\n' "$draft" | grep -F 'authorization: absent' >/dev/null || fail 'blocked execution retained authorization'
approved=$(cc_route 'Run approved plan checkout-validation.')
printf '%s\n' "$approved" | grep -F 'capability: execute-plan' >/dev/null || fail 'approved execution route missing'

mismatch_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-projection-mismatch.XXXXXX")
trap 'rm -rf "$workspace_fixture" "$mismatch_root"' EXIT HUP INT TERM
cp -R "$ROOT/test/contracts/fixtures/identity-projection/mismatched/." "$mismatch_root/"
mkdir -p "$mismatch_root/wrapper/contracts"
cp "$ROOT/wrapper/contracts/routes.yaml" "$mismatch_root/wrapper/contracts/routes.yaml"
cp "$ROOT/wrapper/contracts/context-sets.yaml" "$mismatch_root/wrapper/contracts/context-sets.yaml"
mismatch_route=$(cc_route 'What is this workspace?' "$mismatch_root")
printf '%s\n' "$mismatch_route" | grep -F 'reason_codes:' >/dev/null || fail 'projection-mismatch route omitted reason codes'
printf '%s\n' "$mismatch_route" | grep -F 'projection-mismatch' >/dev/null || fail 'mismatched identity region did not block entry'
printf '%s\n' "$mismatch_route" | grep -F 'authorization: absent' >/dev/null || fail 'projection-mismatch retained authorization'
printf '%s\n' "$mismatch_route" | grep -F 'eligibility: blocked' >/dev/null || fail 'projection-mismatch was eligible'
missing_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-projection-missing.XXXXXX")
trap 'rm -rf "$workspace_fixture" "$mismatch_root" "$missing_root"' EXIT HUP INT TERM
cp -R "$ROOT/test/contracts/fixtures/identity-projection/missing-region/." "$missing_root/"
mkdir -p "$missing_root/wrapper/contracts"
cp "$ROOT/wrapper/contracts/routes.yaml" "$missing_root/wrapper/contracts/routes.yaml"
cp "$ROOT/wrapper/contracts/context-sets.yaml" "$missing_root/wrapper/contracts/context-sets.yaml"
missing_route=$(cc_route 'Help me build this.' "$missing_root")
printf '%s\n' "$missing_route" | grep -F 'projection-mismatch' >/dev/null || fail 'missing identity region did not block write-like entry'
missing_files_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-projection-missing-files.XXXXXX")
trap 'rm -rf "$workspace_fixture" "$mismatch_root" "$missing_root" "$missing_files_root"' EXIT HUP INT TERM
cp -R "$ROOT/test/contracts/fixtures/identity-projection/missing-files/." "$missing_files_root/"
mkdir -p "$missing_files_root/wrapper/contracts"
cp "$ROOT/wrapper/contracts/routes.yaml" "$missing_files_root/wrapper/contracts/routes.yaml"
cp "$ROOT/wrapper/contracts/context-sets.yaml" "$missing_files_root/wrapper/contracts/context-sets.yaml"
missing_files_route=$(cc_route 'Help me build this.' "$missing_files_root")
printf '%s\n' "$missing_files_route" | grep -F 'projection-mismatch' >/dev/null || fail 'missing identity summaries did not fail closed'
printf '%s\n' "$missing_files_route" | grep -F 'capability: initialize' >/dev/null && fail 'missing identity summaries still routed to initialize'
printf '%s\n' "$missing_files_route" | grep -F 'authorization: absent' >/dev/null || fail 'missing identity summaries retained authorization'
fx_effects="$ROOT/test/contracts/fixtures/identity-projection/undeclared-effect"
undeclared=$(cc_validate_delegated_effects "$fx_effects/approved-effects.txt" "$fx_effects/delegated-effects.txt" || true)
assert_eq "$undeclared" UNDECLARED_EFFECT
matching_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-projection-ok.XXXXXX")
trap 'rm -rf "$workspace_fixture" "$mismatch_root" "$missing_root" "$missing_files_root" "$matching_root"' EXIT HUP INT TERM
cp -R "$ROOT/test/contracts/fixtures/identity-projection/matching/." "$matching_root/"
mkdir -p "$matching_root/wrapper/contracts"
cp "$ROOT/wrapper/contracts/routes.yaml" "$matching_root/wrapper/contracts/routes.yaml"
cp "$ROOT/wrapper/contracts/context-sets.yaml" "$matching_root/wrapper/contracts/context-sets.yaml"
ok_route=$(cc_route 'What is this workspace?' "$matching_root")
printf '%s\n' "$ok_route" | grep -F 'capability: orient' >/dev/null || fail 'matching identity region blocked orientation'
printf '%s\n' "$ok_route" | grep -F 'projection-mismatch' >/dev/null && fail 'matching identity region emitted projection-mismatch'
pass "two-stage router and $fixture_count normalized fixtures"
