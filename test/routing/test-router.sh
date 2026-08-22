#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

fixture_count=0
while IFS='	' read -r request expected; do
  test -n "$request" || continue
  fixture_count=$((fixture_count + 1))
  actual=$(cc_route "$request" | sed -n 's/^capability: //p')
  assert_eq "$actual" "$expected"
done <<EOF
$(awk '
  /^    request: / { request=substr($0,14) }
  /^    action: / { print request "\t" $2 }
' "$ROOT/test/routing/fixtures.yaml")
EOF
test "$fixture_count" -ge 50 || fail "routing fixture count: $fixture_count"

workspace_fixture=$(mktemp -d "${TMPDIR:-/tmp}/cc-uninitialized-route.XXXXXX")
trap 'rm -rf "$workspace_fixture"' EXIT HUP INT TERM
cp "$ROOT/template/workspace.yaml" "$workspace_fixture/workspace.yaml"
uninitialized_build=$(cc_route 'Can you help me build this?' "$workspace_fixture")
printf '%s\n' "$uninitialized_build" | grep -F 'capability: initialize' >/dev/null || fail 'uninitialized build request was not routed to initialization'
printf '%s\n' "$uninitialized_build" | grep -F 'human_gate: identity-acceptance' >/dev/null || fail 'uninitialized build request missed identity gate'
printf '%s\n' "$uninitialized_build" | grep -F 'authorization: confirmed-gate-required' >/dev/null || fail 'uninitialized build request was treated as authorized'
uninitialized_read=$(cc_route 'What is this workspace?' "$workspace_fixture")
printf '%s\n' "$uninitialized_read" | grep -F 'capability: orient' >/dev/null || fail 'uninitialized read-only orientation was blocked'
printf '%s\n' "$uninitialized_read" | grep -F 'authorization: read-only' >/dev/null || fail 'uninitialized read-only orientation was not read-only'
sed -i 's/status: uninitialized/status: accepted/' "$workspace_fixture/workspace.yaml"
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

draft=$(cc_route 'Run draft plan checkout-validation.')
printf '%s\n' "$draft" | grep -F 'eligibility: blocked' >/dev/null || fail 'draft execution was not blocked'
printf '%s\n' "$draft" | grep -F 'authorization: absent' >/dev/null || fail 'blocked execution retained authorization'
approved=$(cc_route 'Run approved plan checkout-validation.')
printf '%s\n' "$approved" | grep -F 'capability: execute-plan' >/dev/null || fail 'approved execution route missing'
pass "two-stage router and $fixture_count normalized fixtures"
