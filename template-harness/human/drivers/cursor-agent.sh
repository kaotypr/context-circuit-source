#!/bin/sh
# Cursor Agent live driver for the human-simulated template harness.
#
# Cursor is a thin transport here. It runs the coordinator in the disposable
# workspace and records the same transcript shape as the other host drivers.
# Cursor does not expose a dependable file-read trace, so dimension C degrades.
set -eu

: "${CC_WORKSPACE:?}" "${CC_CASE_FILE:?}" "${CC_TRANSCRIPT:?}" "${CC_TRACE:?}"
: "${CC_HUMAN_SIM:?}" "${CC_RUN_DIR:?}" "${CC_TELEMETRY:?}"
CURSOR=${CC_CURSOR_BIN:-cursor-agent}
TURN_TIMEOUT=${CC_TURN_TIMEOUT:-300}

command -v "$CURSOR" >/dev/null 2>&1 || { printf 'FAIL: no Cursor Agent CLI (%s)\n' "$CURSOR" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { printf 'FAIL: jq required for Cursor JSON processing\n' >&2; exit 1; }

TIMEOUT_CMD=
if command -v timeout >/dev/null 2>&1; then TIMEOUT_CMD=timeout
elif command -v gtimeout >/dev/null 2>&1; then TIMEOUT_CMD=gtimeout
else printf '[driver] no timeout command found; running Cursor without a wall clock\n' >&2
fi

run_cursor() {
	rc_wd=$1; shift
	if [ -n "$TIMEOUT_CMD" ]; then
		( cd "$rc_wd" && "$TIMEOUT_CMD" "$TURN_TIMEOUT" "$CURSOR" "$@" )
	else
		( cd "$rc_wd" && "$CURSOR" "$@" )
	fi
}

RAW="$CC_RUN_DIR/coordinator-stream.jsonl"
CHAT_FILE="$CC_RUN_DIR/.cursor-chat-id"
: > "$CC_TRANSCRIPT"
: > "$CC_TRACE"
: > "$RAW"
: > "$CC_TELEMETRY"
rm -f "$CHAT_FILE"

turns_tsv() {
	awk '
		/^  turns:/{t=1;next}
		t && /^  [A-Za-z]/ && $0 !~ /^    /{t=0}
		t && /^    -[[:space:]]/{
			l=$0; sub(/^    -[[:space:]]*/,"",l)
			k=l; sub(/:.*/,"",k)
			v=l; sub(/^[^:]*:[[:space:]]*/,"",v); sub(/[[:space:]]*#.*$/,"",v); sub(/[[:space:]]+$/,"",v)
			sub(/^"/,"",v); sub(/"$/,"",v)
			print k "\t" v
		}
	' "$CC_CASE_FILE"
}

reply_from_stream() {
	jq -rs '
		[ .[] | (.result? // .message?.content? // .text? // empty)
		  | if type == "string" then .
		    elif type == "array" then map(if type == "object" then (.text? // .content? // "") else . end) | join("")
		    else "" end
		] | map(select(length > 0)) | last // ""
	' "$1" 2>/dev/null || printf ''
}

chat_from_stream() {
	jq -rs '[.[] | (.chat_id? // .conversation_id? // .thread_id? // empty)
		| select(type == "string" and length > 0)] | last // ""' "$1" 2>/dev/null || printf ''
}

coordinator_turn() {
	msg=$1; action=$2
	effective_msg=$msg
	case "${CC_FAULT:-}" in
		verifier-unavailable)
			effective_msg="[Harness condition: the independent verifier child is unavailable for this run. Do not attempt to create or imitate that child; report the result as host-blocked and preserve the existing work.]\n\n$msg"
			;;
	esac
	turn_raw="$CC_RUN_DIR/.cursor-turn.$$.jsonl"
	rm -f "$turn_raw"
	if [ ! -f "$CHAT_FILE" ]; then
		if ! run_cursor "$CC_WORKSPACE" --workspace "$CC_WORKSPACE" --print \
			--output-format stream-json --sandbox enabled --force --trust "$effective_msg" \
			> "$turn_raw" 2>>"$CC_RUN_DIR/driver.err"; then
			cat "$turn_raw" >> "$RAW" 2>/dev/null || :
			printf 'FAIL: Cursor coordinator turn failed; see %s/driver.err\n' "$CC_RUN_DIR" >&2
			return 1
		fi
		chat=$(chat_from_stream "$turn_raw")
		[ -n "$chat" ] && printf '%s\n' "$chat" > "$CHAT_FILE" || :
	else
		chat=$(sed -n '1p' "$CHAT_FILE")
		if ! run_cursor "$CC_WORKSPACE" --workspace "$CC_WORKSPACE" --resume "$chat" \
			--print --output-format stream-json --sandbox enabled --force --trust "$effective_msg" \
			> "$turn_raw" 2>>"$CC_RUN_DIR/driver.err"; then
			cat "$turn_raw" >> "$RAW" 2>/dev/null || :
			printf 'FAIL: Cursor coordinator resume failed; see %s/driver.err\n' "$CC_RUN_DIR" >&2
			return 1
		fi
	fi

	cat "$turn_raw" >> "$RAW"
	reply=$(reply_from_stream "$turn_raw")
	[ -n "$reply" ] || { printf 'FAIL: Cursor coordinator emitted no final message\n' >&2; return 1; }
	# Cursor usage fields vary by CLI release; preserve a zero-valued ledger row.
	printf '%s\t0\t0\t0\t0\t0.0000\n' "$action" >> "$CC_TELEMETRY"
	rm -f "$turn_raw"
	printf '%s' "$reply"
}

offered_plan() { printf '%s' "$1" | grep -Eiq 'plan|outline|sketch|write .*up|get (you )?(set up|started)'; }
asked_question() { printf '%s' "$1" | grep -q '?'; }

printf '# transcript: %s (host=cursor-agent)\n' "$CC_CASE_ID" >> "$CC_TRANSCRIPT"
turns_file="$CC_RUN_DIR/.cursor-turns.tsv"
turns_tsv > "$turns_file"
ACTION=orient
LAST_REPLY=
while IFS='	' read -r kind payload; do
	[ -n "$kind" ] || continue
	send=
	case "$kind" in
		say) send=$payload ;;
		on_offer_plan)
			if offered_plan "$LAST_REPLY"; then send=$payload; ACTION=create-plan
			else printf '(skipped on_offer_plan: coordinator did not offer a plan)\n' >> "$CC_TRANSCRIPT"; continue; fi ;;
		on_open_questions)
			if asked_question "$LAST_REPLY"; then send="The simplest option sounds fine. For now, please only show me the plan and stop there."
			else printf '(skipped on_open_questions: coordinator asked nothing)\n' >> "$CC_TRANSCRIPT"; continue; fi ;;
		*) printf '(unknown turn kind: %s)\n' "$kind" >> "$CC_TRANSCRIPT"; continue ;;
	esac
	send_lc=$(printf '%s' "$send" | tr '[:upper:]' '[:lower:]')
	case "$send_lc" in
		*"work with me directly"*|*collaborate*directly*) ACTION=direct-collaboration ;;
		*approve*build*|*approve*execute*|*approve*run*) ACTION=execute-plan ;;
		*"build the"*|*"build all"*|*execute*|*"run them"*|*"run all"*|*"run the"*|*"go ahead and build"*) ACTION=execute-plan ;;
		*approve*) ACTION=approve ;;
		*"show me the plan"*|*review*) ACTION=review ;;
		*connect*|*"hook it up"*|*"hook up"*) ACTION=connect-repo ;;
		*plan*|*"show me"*) ACTION=create-plan ;;
	esac
	printf '\nhuman: %s\n' "$send" >> "$CC_TRANSCRIPT"
	reply=$(coordinator_turn "$send" "$ACTION")
	printf 'coordinator: %s\n' "$reply" >> "$CC_TRANSCRIPT"
	LAST_REPLY=$reply
done < "$turns_file"
rm -f "$turns_file"

HUMAN_BLOCK=$(awk '/^human:/{f=1} /^grader:/{f=0} f{print}' "$CC_CASE_FILE")
VERDICT="$CC_RUN_DIR/conversational-verdict.txt"
SCRATCH="$CC_RUN_DIR/.cursor-sim"
mkdir -p "$SCRATCH"
{
	printf 'Act as the tool-less human-simulator and judge only supplied text.\n\n'
	printf 'Here is the case human: block:\n\n%s\n\n' "$HUMAN_BLOCK"
	printf 'Here is the full transcript:\n\n'
	cat "$CC_TRANSCRIPT"
	printf '\n\nFor EACH visible_expectation, output one line: "PASS: <evidence>" or "FAIL: <evidence>", in order. Then a final line "verdict: pass" or "verdict: fail". Do not use tools or inspect files.\n'
} > "$SCRATCH/prompt.txt"
run_cursor "$SCRATCH" --workspace "$SCRATCH" --print --output-format text --mode ask \
	--sandbox enabled --trust "$(cat "$SCRATCH/prompt.txt")" > "$VERDICT" 2>>"$CC_RUN_DIR/driver.err" || {
	printf 'FAIL: Cursor human-simulator failed; see %s/driver.err\n' "$CC_RUN_DIR" >&2
	exit 1
}
rm -rf "$SCRATCH"
rm -f "$CC_TRACE"
printf '[driver] Cursor exposes no dependable file-read trace; dimension C will degrade\n'
exit 0
