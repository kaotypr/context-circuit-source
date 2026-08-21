#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

rows=$(awk '/^  - / {count++} END {print count+0}' "$ROOT/test/behavior-matrix/matrix.yaml")
test "$rows" -eq 40 || fail "behavior matrix rows: $rows"
test "$(grep -c 'positive_and_negative_required: true' "$ROOT/test/behavior-matrix/matrix.yaml")" -eq 1 || fail 'negative matrix requirement missing'
test "$(grep -c '^  - id: R' "$ROOT/test/routing/fixtures.yaml")" -ge 50 || fail 'routing matrix coverage below 50'
pair_rows=$(grep -v '^#' "$ROOT/test/behavior-matrix/fixture-pairs.tsv" | sed '/^$/d' | wc -l)
test "$pair_rows" -eq 40 || fail "behavior fixture pairs: $pair_rows"
while IFS='|' read -r case_name positive negative; do
  case "$case_name" in ''|'#'*) continue ;; esac
  grep -F "  - id: $positive" "$ROOT/test/routing/fixtures.yaml" >/dev/null || fail "missing positive fixture $positive"
  grep -F "  - id: $negative" "$ROOT/test/routing/fixtures.yaml" >/dev/null || fail "missing negative fixture $negative"
done < "$ROOT/test/behavior-matrix/fixture-pairs.tsv"
for request in 'What is this workspace?' 'Run approved plan checkout-validation.' 'Clean runtime after showing me dirty work.'; do
  cc_route "$request" >/dev/null || fail "router unavailable for $request"
done
contains "$ROOT/wrapper/contracts/invariants.yaml" 'INV-GATE-05'
pass '40 behavior-matrix routes with positive/negative obligation and 60 fixtures'
