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

draft=$(cc_route 'Run draft plan checkout-validation.')
printf '%s\n' "$draft" | grep -F 'eligibility: blocked' >/dev/null || fail 'draft execution was not blocked'
printf '%s\n' "$draft" | grep -F 'authorization: absent' >/dev/null || fail 'blocked execution retained authorization'
approved=$(cc_route 'Run approved plan checkout-validation.')
printf '%s\n' "$approved" | grep -F 'capability: execute-plan' >/dev/null || fail 'approved execution route missing'
pass "two-stage router and $fixture_count normalized fixtures"
