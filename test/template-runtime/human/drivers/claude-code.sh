#!/bin/sh
# Claude Code live driver for the human-simulated harness.
#
# Invoked by run-scenario.sh with the CC_* environment. It:
#   1. drives the real product coordinator as a headless `claude -p` session
#      whose working directory is the instantiated workspace (so it loads that
#      workspace's own CLAUDE.md/AGENTS.md/skills — the real product);
#   2. plays the case's scripted human turns (say / on_* conditionals);
#   3. records the transcript, and a best-effort file-access trace from the
#      coordinator's tool calls (stream-json), tagged by product action;
#   4. runs the human-simulator once, tool-less, to emit the conversational
#      verdict over the finished transcript.
#
# The coordinator runs with --permission-mode bypassPermissions because it is
# headless; this is safe here only because its cwd is a DISPOSABLE isolated
# workspace under .out/ and no extra directories are granted.
set -eu

: "${CC_WORKSPACE:?}" "${CC_CASE_FILE:?}" "${CC_TRANSCRIPT:?}" "${CC_TRACE:?}"
: "${CC_HUMAN_SIM:?}" "${CC_RUN_DIR:?}"
CLAUDE=${CC_CLAUDE_BIN:-claude}
TURN_TIMEOUT=${CC_TURN_TIMEOUT:-300}

command -v "$CLAUDE" >/dev/null 2>&1 || { printf 'FAIL: no claude CLI (%s)\n' "$CLAUDE" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { printf 'FAIL: jq required for trace capture\n' >&2; exit 1; }

new_uuid() { if [ -r /proc/sys/kernel/random/uuid ]; then cat /proc/sys/kernel/random/uuid; else "$CLAUDE" --version >/dev/null; printf '00000000-0000-4000-8000-%012d' "$$"; fi; }
CID=$(new_uuid)
RAW="$CC_RUN_DIR/coordinator-stream.jsonl"    # accumulated stream-json (all turns)
: > "$CC_TRANSCRIPT"; : > "$CC_TRACE"; : > "$RAW"

# Extract the ordered human turns as TYPE<TAB>PAYLOAD.
turns_tsv() {
	awk '
		/^  turns:/{t=1;next}
		t && /^  [A-Za-z]/ && $0 !~ /^    /{t=0}
		t && /^    -[[:space:]]/{
			l=$0; sub(/^    -[[:space:]]*/,"",l);
			k=l; sub(/:.*/,"",k);
			v=l; sub(/^[^:]*:[[:space:]]*/,"",v); sub(/[[:space:]]*#.*$/,"",v); sub(/[[:space:]]+$/,"",v);
			sub(/^"/,"",v); sub(/"$/,"",v);
			print k "\t" v;
		}
	' "$CC_CASE_FILE"
}

# One coordinator turn: send $1, capture reply text + reads, tag reads with $2.
# A flag FILE (not a shell var) tracks session start, because each turn runs in a
# command-substitution subshell that cannot write a variable back to the caller.
STARTED_FLAG="$CC_RUN_DIR/.coord-started"
rm -f "$STARTED_FLAG"
coordinator_turn() {
	msg=$1; action=$2; turn_raw="$CC_RUN_DIR/.turn.$$.jsonl"
	if [ ! -f "$STARTED_FLAG" ]; then
		( cd "$CC_WORKSPACE" && timeout "$TURN_TIMEOUT" "$CLAUDE" -p "$msg" \
			--session-id "$CID" --output-format stream-json --verbose \
			--permission-mode bypassPermissions </dev/null ) > "$turn_raw" 2>>"$CC_RUN_DIR/driver.err" || true
		: > "$STARTED_FLAG"
	else
		( cd "$CC_WORKSPACE" && timeout "$TURN_TIMEOUT" "$CLAUDE" -p "$msg" \
			--resume "$CID" --output-format stream-json --verbose \
			--permission-mode bypassPermissions </dev/null ) > "$turn_raw" 2>>"$CC_RUN_DIR/driver.err" || true
	fi
	cat "$turn_raw" >> "$RAW"
	reply=$(jq -rs 'map(select(.type=="result")) | last | .result // ""' "$turn_raw" 2>/dev/null || printf '')
	# best-effort file-access trace: coordinator read/grep/glob tool calls
	jq -r 'select(.type=="assistant") | .message.content[]?
		| select(.type=="tool_use") | select(.name=="Read" or .name=="Grep" or .name=="Glob")
		| (.input.file_path // .input.path // "") | select(. != "")' "$turn_raw" 2>/dev/null \
	| while IFS= read -r p; do
		rel=$(printf '%s' "$p" | sed "s|^$CC_WORKSPACE/||; s|^\./||")
		printf '%s\tRead\t%s\n' "$action" "$rel" >> "$CC_TRACE"
	done
	rm -f "$turn_raw"
	printf '%s' "$reply"
}

offered_plan() { printf '%s' "$1" | grep -Eiq 'plan|outline|sketch|write .*up|get (you )?(set up|started)'; }
asked_question() { printf '%s' "$1" | grep -q '?'; }

ACTION=orient
LAST_REPLY=""
printf '# transcript: %s (host=claude-code)\n' "$CC_CASE_ID" >> "$CC_TRANSCRIPT"

turns_tsv | while IFS='	' read -r kind payload; do
	[ -n "$kind" ] || continue
	send="";
	case "$kind" in
		say) send=$payload ;;
		on_offer_plan)
			if offered_plan "$LAST_REPLY"; then send=$payload; ACTION=create-plan
			else printf '(skipped on_offer_plan: coordinator did not offer a plan)\n' >> "$CC_TRANSCRIPT"; continue; fi ;;
		on_open_questions)
			if asked_question "$LAST_REPLY"; then
				send="Whatever's simplest is fine — I don't have strong preferences, just make it work for me."
			else printf '(skipped on_open_questions: coordinator asked nothing)\n' >> "$CC_TRANSCRIPT"; continue; fi ;;
		*) printf '(unknown turn kind: %s)\n' "$kind" >> "$CC_TRANSCRIPT"; continue ;;
	esac
	# a plan-related say also advances the action label
	case "$send" in *plan*|*"show me"*) ACTION=create-plan ;; esac
	printf '\nhuman: %s\n' "$send" >> "$CC_TRANSCRIPT"
	reply=$(coordinator_turn "$send" "$ACTION")
	printf 'coordinator: %s\n' "$reply" >> "$CC_TRANSCRIPT"
	LAST_REPLY=$reply
done

# Conversational verdict: the human-simulator judges visible_expectations from the
# transcript alone, tool-less, in a scratch cwd it cannot use to inspect state.
HUMAN_BLOCK=$(awk '/^human:/{f=1} /^grader:/{f=0} f{print}' "$CC_CASE_FILE")
SIM_SYS=$(cat "$CC_HUMAN_SIM")
VERDICT="$CC_RUN_DIR/conversational-verdict.txt"
SCRATCH=$(mktemp -d "${TMPDIR:-/tmp}/cc-sim.XXXXXX")
{
	printf 'You are judging a finished conversation as the human-simulator.\n'
	printf 'Here is the case human: block you were given:\n\n%s\n\n' "$HUMAN_BLOCK"
	printf 'Here is the full transcript of what actually happened:\n\n'
	cat "$CC_TRANSCRIPT"
	printf '\n\nFor EACH visible_expectation, output one line: "PASS: <evidence>" or "FAIL: <evidence>", in order. Then a final line "verdict: pass" or "verdict: fail".\n'
} > "$SCRATCH/prompt.txt"
( cd "$SCRATCH" && timeout "$TURN_TIMEOUT" "$CLAUDE" -p "$(cat "$SCRATCH/prompt.txt")" \
	--append-system-prompt "$SIM_SYS" --output-format text \
	--permission-mode bypassPermissions </dev/null ) > "$VERDICT" 2>>"$CC_RUN_DIR/driver.err" || true
rm -rf "$SCRATCH"

printf '\n[driver] coordinator turns complete; transcript=%s\n' "$CC_TRANSCRIPT"
printf '[driver] conversational verdict=%s\n' "$VERDICT"
[ -s "$CC_TRACE" ] || { rm -f "$CC_TRACE"; printf '[driver] no file reads captured; dimension C will degrade\n'; }
exit 0
