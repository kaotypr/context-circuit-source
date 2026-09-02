#!/bin/sh
# Run the template-harness scenario library across the configured host matrix.
#
# By default this starts one lane per host. Each lane runs every case in its own
# disposable workspace, so Codex, Claude Code, and Cursor Agent can progress at
# the same time without sharing runtime state. Use --prepare-only for a
# credential-free deterministic matrix smoke test.
set -eu

HERE=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT=$(git -C "$HERE" rev-parse --show-toplevel 2>/dev/null) || { printf 'FAIL: not a git source checkout\n' >&2; exit 1; }
SCENARIOS="$HERE/../scenarios"
OUT_ROOT="$HERE/.out"
HOSTS="codex claude-code cursor-agent"
CASES=
PREPARE_ONLY=0

usage() {
	cat <<'EOF'
usage: sh run-matrix.sh [--prepare-only] [--case CASE]...

Runs every selected scenario on codex, claude-code, and cursor-agent in three
parallel host lanes. Live runs use each host's built-in driver and are graded.
--prepare-only creates and records isolated runs without invoking a host CLI.
EOF
}

while [ "$#" -gt 0 ]; do
	case "$1" in
		--prepare-only) PREPARE_ONLY=1; shift ;;
		--case) CASES="$CASES ${2:?--case needs a value}"; shift 2 ;;
		--case=*) CASES="$CASES ${1#--case=}"; shift ;;
		-h|--help) usage; exit 0 ;;
		*) printf 'FAIL: unknown argument: %s\n' "$1" >&2; usage >&2; exit 2 ;;
	esac
done

case_selected() {
	[ -z "$CASES" ] && return 0
	for wanted in $CASES; do [ "$wanted" = "$1" ] && return 0; done
	return 1
}

MATRIX_ID="$(date -u +%Y%m%dT%H%M%SZ)-$$"
MATRIX_DIR="$OUT_ROOT/matrix/$MATRIX_ID"
mkdir -p "$MATRIX_DIR"
printf 'case_id\thost\tstatus\trun_dir\tlog\n' > "$MATRIX_DIR/summary.tsv"

run_lane() {
	lane_host=$1
	lane_rows="$MATRIX_DIR/$lane_host.tsv"
	: > "$lane_rows"
	lane_rc=0
	for case_dir in "$SCENARIOS"/*; do
		[ -f "$case_dir/case.yaml" ] || continue
		case_id=$(awk '$0 ~ /^id:[[:space:]]/ {sub(/^id:[[:space:]]*/,""); print; exit}' "$case_dir/case.yaml")
		[ -n "$case_id" ] || continue
		case_selected "$case_id" || continue
		case_log="$MATRIX_DIR/$lane_host-$case_id.log"
		if [ "$PREPARE_ONLY" -eq 1 ]; then
			if sh "$HERE/run-scenario.sh" --host "$lane_host" --no-grade "$case_id" > "$case_log" 2>&1; then
				status=prepared
			else
				status=fail; lane_rc=1
			fi
		else
			if sh "$HERE/run-scenario.sh" --host "$lane_host" --live "$case_id" > "$case_log" 2>&1; then
				status=pass
			else
				status=fail; lane_rc=1
			fi
		fi
		run_dir=$(sed -n 's/^prepared run: //p' "$case_log" | sed -n '1p')
		printf '%s\t%s\t%s\t%s\t%s\n' "$case_id" "$lane_host" "$status" "$run_dir" "$case_log" >> "$lane_rows"
	done
	return "$lane_rc"
}

pids=
for host in $HOSTS; do
	run_lane "$host" &
	pids="$pids $!"
done

MATRIX_RC=0
for pid in $pids; do
	if wait "$pid"; then :; else MATRIX_RC=1; fi
done

for host in $HOSTS; do
	[ -f "$MATRIX_DIR/$host.tsv" ] && cat "$MATRIX_DIR/$host.tsv" >> "$MATRIX_DIR/summary.tsv"
done

rows=$(awk 'NR>1 && NF{n++} END{print n+0}' "$MATRIX_DIR/summary.tsv")
printf 'matrix: %s\n' "$MATRIX_DIR"
printf 'summary: %s\n' "$MATRIX_DIR/summary.tsv"
printf 'runs: %s (hosts=%s, mode=%s)\n' "$rows" "$HOSTS" "$([ "$PREPARE_ONLY" -eq 1 ] && printf prepare || printf live)"
cat "$MATRIX_DIR/summary.tsv"
exit "$MATRIX_RC"
