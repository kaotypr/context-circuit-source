#!/bin/sh
set -eu
. "$(git -C "$(dirname -- "$0")" rev-parse --show-toplevel)/test/lib/assert.sh"

fixture="$ROOT/agent-harness/fixtures/role-tiering.yaml"
require_file "$fixture"
for host in codex claude-code cursor-agent; do
	contains "$fixture" "  $host:"
done
assert_eq 6 "$(grep -c '^      effort: medium$' "$fixture")"
not_contains "$ROOT/agent-harness/scenarios/13-execution-latency/case.yaml" 'role_tiering:'
not_contains "$ROOT/agent-harness/scenarios/14-codex-direct-collaboration/case.yaml" 'role_tiering:'

matrix_output=$(sh "$ROOT/agent-harness/human/run-matrix.sh" --prepare-only)
matrix_dir=$(printf '%s\n' "$matrix_output" | sed -n 's/^matrix: //p' | sed -n '1p')
[ -n "$matrix_dir" ] || fail 'matrix runner did not report its output directory'
summary="$matrix_dir/summary.tsv"
require_file "$summary"
trap '
	while IFS="	" read -r _ _ _ run_dir _; do
		[ -n "$run_dir" ] && rm -rf "$run_dir"
	done < "$summary"
	rm -rf "$matrix_dir"
' EXIT HUP INT TERM

assert_eq 66 "$(awk 'NR > 1 && NF {n++} END {print n+0}' "$summary")"
for host in codex claude-code cursor-agent; do
	assert_eq 22 "$(awk -F '\t' -v h="$host" 'NR > 1 && $2 == h {n++} END {print n+0}' "$summary")"
	if awk -F '\t' 'NR > 1 && $3 != "prepared" {bad=1} END {exit bad ? 0 : 1}' "$summary"; then
		fail "matrix runner did not prepare every selected case for $host"
	fi
done

first_run=$(awk -F '\t' 'NR == 2 {print $4}' "$summary")
require_file "$first_run/workspace/role-tiering.local.yaml"
cmp "$fixture" "$first_run/workspace/role-tiering.local.yaml" >/dev/null || fail 'generated role-tiering config differs from shared fixture'

pass 'shared role-tiering fixture and three-host matrix'
