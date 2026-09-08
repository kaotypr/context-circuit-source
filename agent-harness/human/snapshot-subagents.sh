#!/bin/sh
# Snapshot the sub-agent JSONL transcripts for one in-session (cc-test-case) run
# into <run-dir>/subagents/.
#
# The shell --live driver (drivers/claude-code.sh) records the coordinator as one
# headless `claude -p` stream at <run-dir>/coordinator-stream.jsonl. The in-session
# cc-test-case driver instead spawns real nested sub-agents (coordinator, workers,
# verifiers) whose transcripts live in the host session store
# (~/.claude/projects/<project>/<session>/subagents/agent-<id>.jsonl), NOT under
# .out. This copies the ones belonging to a run into the run dir, so both drivers
# leave comparable transcript artifacts alongside transcript.txt.
#
# An agent belongs to the run when its transcript references the run's workspace
# path (from run.yaml) — which scopes the copy to this case and excludes agents
# from other scenarios in the same session.
#
# usage: sh snapshot-subagents.sh <run-dir>
set -eu

[ "$#" -eq 1 ] || { printf 'usage: sh snapshot-subagents.sh <run-dir>\n' >&2; exit 2; }
RUN=$1
[ -f "$RUN/run.yaml" ] || { printf 'FAIL: no run.yaml in %s\n' "$RUN" >&2; exit 2; }
WS=$(awk '/^workspace:[[:space:]]/{sub(/^workspace:[[:space:]]*/,"");print;exit}' "$RUN/run.yaml")
[ -n "$WS" ] || { printf 'FAIL: no workspace path in %s/run.yaml\n' "$RUN" >&2; exit 2; }

OUT="$RUN/subagents"
mkdir -p "$OUT"
n=0
for f in "$HOME"/.claude/projects/*/*/subagents/agent-*.jsonl; do
	[ -f "$f" ] || continue
	# a sub-agent belongs to this run if its transcript references the run workspace
	if grep -Fq "$WS" "$f" 2>/dev/null; then
		cp "$f" "$OUT/$(basename "$f")"
		n=$((n + 1))
	fi
done

printf 'snapshot: copied %s sub-agent transcript(s) into %s\n' "$n" "$OUT"
[ "$n" -gt 0 ] || printf 'note: no sub-agent transcripts matched this run workspace (host store empty, or a different driver was used)\n' >&2
