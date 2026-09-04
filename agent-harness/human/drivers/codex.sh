#!/bin/sh
# Codex live driver for the human-simulated harness.
#
# Invoked by run-scenario.sh with the CC_* environment. It:
#   1. drives the real product coordinator as a resumable `codex exec` session
#      rooted at the instantiated workspace;
#   2. plays the case's scripted human turns (say / on_* conditionals);
#   3. records the transcript, raw Codex JSONL, and best-effort telemetry;
#   4. runs an ephemeral, read-only Codex session as the human-simulator to emit
#      the conversational verdict over the finished transcript.
#
# Codex's JSONL stream does not expose dependable file-read events comparable to
# Claude Code's Read/Grep/Glob tool calls. This driver intentionally emits no
# CC_TRACE rather than claiming an inaccurate access audit; dimension C degrades
# to warning-only while the raw stream remains available for inspection.
set -eu

: "${CC_WORKSPACE:?}" "${CC_CASE_FILE:?}" "${CC_TRANSCRIPT:?}" "${CC_TRACE:?}"
: "${CC_TELEMETRY:?}" "${CC_HUMAN_SIM:?}" "${CC_RUN_DIR:?}"
CODEX=${CC_CODEX_BIN:-codex}
TURN_TIMEOUT=${CC_TURN_TIMEOUT:-300}

command -v "$CODEX" >/dev/null 2>&1 || { printf 'FAIL: no codex CLI (%s)\n' "$CODEX" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { printf 'FAIL: jq required for Codex JSONL processing\n' >&2; exit 1; }

# A per-turn timeout is best-effort: GNU `timeout` (Linux) or `gtimeout`
# (coreutils on macOS). When neither is present, Codex runs without a wall clock.
TIMEOUT_CMD=
if command -v timeout >/dev/null 2>&1; then TIMEOUT_CMD=timeout
elif command -v gtimeout >/dev/null 2>&1; then TIMEOUT_CMD=gtimeout
else printf '[driver] no timeout command found; running Codex without a per-turn wall clock\n' >&2
fi

# run_codex WORKDIR ARGS... (prompt is supplied on stdin).
run_codex() {
	rc_wd=$1; shift
	if [ -n "$TIMEOUT_CMD" ]; then
		( cd "$rc_wd" && "$TIMEOUT_CMD" "$TURN_TIMEOUT" "$CODEX" "$@" )
	else
		( cd "$rc_wd" && "$CODEX" "$@" )
	fi
}

RAW="$CC_RUN_DIR/coordinator-stream.jsonl"
TELRAW="$CC_RUN_DIR/.telemetry-raw.tsv"
THREAD_FILE="$CC_RUN_DIR/.codex-thread-id"
ROLE_EVIDENCE=${CC_ROLE_EVIDENCE:-$CC_RUN_DIR/role-evidence.tsv}
CODEX_STATE_DIR=${CC_CODEX_STATE_DIR:-${CODEX_HOME:-$HOME/.codex}}
: > "$CC_TRANSCRIPT"
: > "$CC_TRACE"
: > "$RAW"
: > "$TELRAW"
: > "$ROLE_EVIDENCE"
rm -f "$THREAD_FILE"

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

# One coordinator turn: send $1, capture the final reply, and record telemetry
# tagged by $2. The thread id lives in a file because command substitution runs
# this function in a subshell.
coordinator_turn() {
	msg=$1; action=$2
	turn_raw="$CC_RUN_DIR/.codex-turn.$$.jsonl"
	turn_reply="$CC_RUN_DIR/.codex-reply.$$"
	rm -f "$turn_raw" "$turn_reply"
	effective_msg=$msg
	case "${CC_FAULT:-}" in
		verifier-unavailable)
			effective_msg="[Harness condition: the independent verifier child is unavailable for this run. Do not attempt to create or imitate that child; report the result as host-blocked and preserve the existing work.]\n\n$msg"
			;;
	esac

	if [ ! -f "$THREAD_FILE" ]; then
		if ! printf '%b\n' "$effective_msg" | run_codex "$CC_WORKSPACE" exec \
			--cd "$CC_WORKSPACE" --approve-for-me \
			--json --output-last-message "$turn_reply" - \
			>"$turn_raw" 2>>"$CC_RUN_DIR/driver.err"; then
			[ -f "$turn_raw" ] && cat "$turn_raw" >> "$RAW"
			printf 'FAIL: Codex coordinator turn failed; see %s/driver.err\n' "$CC_RUN_DIR" >&2
			return 1
		fi
		thread_id=$(jq -r 'select(.type == "thread.started") | .thread_id // empty' "$turn_raw" | sed -n '1p')
		[ -n "$thread_id" ] || { cat "$turn_raw" >> "$RAW"; printf 'FAIL: Codex emitted no thread.started id\n' >&2; return 1; }
		printf '%s\n' "$thread_id" > "$THREAD_FILE"
	else
		thread_id=$(sed -n '1p' "$THREAD_FILE")
		if ! printf '%b\n' "$effective_msg" | run_codex "$CC_WORKSPACE" exec resume \
			-c 'sandbox_mode="workspace-write"' --json \
			--output-last-message "$turn_reply" "$thread_id" - \
			>"$turn_raw" 2>>"$CC_RUN_DIR/driver.err"; then
			[ -f "$turn_raw" ] && cat "$turn_raw" >> "$RAW"
			printf 'FAIL: Codex coordinator resume failed; see %s/driver.err\n' "$CC_RUN_DIR" >&2
			return 1
		fi
	fi

	cat "$turn_raw" >> "$RAW"
	reply=
	[ -f "$turn_reply" ] && reply=$(cat "$turn_reply")
	if [ -z "$reply" ]; then
		reply=$(jq -rs 'map(select(.type == "item.completed" and .item.type == "agent_message") | .item.text) | last // ""' "$turn_raw" 2>/dev/null || printf '')
	fi
	[ -n "$reply" ] || { printf 'FAIL: Codex coordinator emitted no final message\n' >&2; return 1; }

	# Codex reports per-turn token usage but not price. Record the harness ledger as:
	# action, agent-loop turns, output tokens, context peak, cost (always unknown/0).
	usage=$(jq -rs '
		map(select(.type == "turn.completed")) as $turns
		| ($turns | last | .usage // {}) as $usage
		| [($turns | length), ($usage.output_tokens // 0), ($usage.input_tokens // 0), 0]
		| @tsv
	' "$turn_raw" 2>/dev/null || printf '')
	[ -n "$usage" ] && printf '%s\t%s\n' "$action" "$usage" >> "$TELRAW"

	rm -f "$turn_raw" "$turn_reply"
	printf '%s' "$reply"
}

offered_plan() { printf '%s' "$1" | grep -Eiq 'plan|outline|sketch|write .*up|get (you )?(set up|started)'; }
asked_question() { printf '%s' "$1" | grep -q '?'; }

ACTION=orient
LAST_REPLY=
printf '# transcript: %s (host=codex)\n' "$CC_CASE_ID" >> "$CC_TRANSCRIPT"

turns_tsv | while IFS='	' read -r kind payload; do
	[ -n "$kind" ] || continue
	send=
	case "$kind" in
		say) send=$payload ;;
		on_offer_plan)
			if offered_plan "$LAST_REPLY"; then send=$payload; ACTION=create-plan
			else printf '(skipped on_offer_plan: coordinator did not offer a plan)\n' >> "$CC_TRANSCRIPT"; continue; fi ;;
		on_open_questions)
			if asked_question "$LAST_REPLY"; then
				send="The simplest option sounds fine. For now, please only show me the plan and stop there."
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
done

# Aggregate per-turn usage into the six-column efficiency ledger:
# action, conversational turns, agent-loop turns, generated output tokens,
# context peak, cost. Codex JSONL has no cost field, so cost remains 0.
if [ -s "$TELRAW" ]; then
	awk -F'\t' '
		{ conv[$1]++; at[$1]+=$2; ot[$1]+=$3; if(($4+0)>cp[$1]) cp[$1]=$4+0; cost[$1]+=$5 }
		END { for (a in conv) printf "%s\t%d\t%d\t%d\t%d\t%.4f\n", a, conv[a], at[a], ot[a], cp[a], cost[a] }
	' "$TELRAW" > "$CC_TELEMETRY"
	printf '[driver] wrote efficiency telemetry: %s\n' "$CC_TELEMETRY" >&2

fi
rm -f "$TELRAW"

# Codex persists bounded child metadata separately from `exec --json`. Record
# only role/model/effort for immediate native sub-agent children whose task
# names follow the adapter's *_worker / *_verifier convention. The coordinator
# itself may use a resumable thread, but a generic peer thread is not a child
# role and must never be counted as one. Never copy prompts or payloads.
record_role_evidence() {
	[ -f "$THREAD_FILE" ] || return 0
	parent_id=$(sed -n '1p' "$THREAD_FILE")
	sessions="$CODEX_STATE_DIR/sessions"
	[ -d "$sessions" ] || return 0
	find "$sessions" -type f -name '*.jsonl' 2>/dev/null \
	| while IFS= read -r session_file; do
		agent_path=$(jq -r --arg parent "$parent_id" 'select(.type == "session_meta")
			| if ((.payload.source | type) == "object")
			     and (.payload.source.subagent.thread_spawn.parent_thread_id == $parent)
			  then (.payload.source.subagent.thread_spawn.agent_path // "")
			  else ""
			  end' \
			"$session_file" 2>/dev/null | sed -n '1p')
		case "$agent_path" in
			*_worker) role=worker ;;
			*_verifier) role=verifier ;;
			*) continue ;;
		esac
		model_effort=$(jq -r 'select(.type == "turn_context")
			| [(.payload.model // ""), (.payload.effort // "")] | @tsv' \
			"$session_file" 2>/dev/null | sed -n '1p')
		[ -n "$model_effort" ] && printf '%s\t%s\n' "$role" "$model_effort" >> "$ROLE_EVIDENCE"
	done
	if [ -s "$ROLE_EVIDENCE" ]; then
		sort -u "$ROLE_EVIDENCE" > "$ROLE_EVIDENCE.tmp" && mv "$ROLE_EVIDENCE.tmp" "$ROLE_EVIDENCE"
		printf '[driver] wrote bounded role evidence: %s\n' "$ROLE_EVIDENCE" >&2
	else
		rm -f "$ROLE_EVIDENCE"
	fi
}
record_role_evidence

# Conversational verdict: an isolated, ephemeral Codex session judges only the
# simulator definition, case human block, and visible transcript supplied here.
HUMAN_BLOCK=$(awk '/^human:/{f=1} /^grader:/{f=0} f{print}' "$CC_CASE_FILE")
SIM_SYS=$(cat "$CC_HUMAN_SIM")
VERDICT="$CC_RUN_DIR/conversational-verdict.txt"
SCRATCH="$CC_RUN_DIR/.sim"
mkdir -p "$SCRATCH"
{
	printf '%s\n\n' "$SIM_SYS"
	printf 'You are judging a finished conversation as the human-simulator.\n'
	printf 'Here is the case human: block you were given:\n\n%s\n\n' "$HUMAN_BLOCK"
	printf 'Here is the full transcript of what actually happened:\n\n'
	cat "$CC_TRANSCRIPT"
	printf '\n\nFor EACH visible_expectation, output one line: "PASS: <evidence>" or "FAIL: <evidence>", in order. Then a final line "verdict: pass" or "verdict: fail". Do not use tools or inspect files.\n'
} > "$SCRATCH/prompt.txt"
if ! run_codex "$SCRATCH" exec --cd "$SCRATCH" --skip-git-repo-check \
	--sandbox read-only --ephemeral --output-last-message "$VERDICT" - \
	< "$SCRATCH/prompt.txt" >"$SCRATCH/stdout.txt" 2>>"$CC_RUN_DIR/driver.err"; then
	printf 'FAIL: Codex human-simulator failed; see %s/driver.err\n' "$CC_RUN_DIR" >&2
	exit 1
fi
rm -rf "$SCRATCH"

# Isolation guard: surface any repository binding whose path escapes the
# disposable workspace. The Codex sandbox prevents writes outside its root, but
# recording an invalid binding still makes the run diagnosable.
LOCAL="$CC_WORKSPACE/repositories.local.yaml"
if [ -f "$LOCAL" ]; then
	awk '/^[[:space:]]+path:[[:space:]]/{sub(/^[[:space:]]+path:[[:space:]]*/,"");sub(/[[:space:]]+$/ ,"");print}' "$LOCAL" \
	| while IFS= read -r bp; do
		[ -n "$bp" ] || continue
		case "$bp" in
			"$CC_WORKSPACE"/*|repositories/*|./repositories/*|.|./) : ;;
			/*) printf '[driver] ISOLATION WARNING: binding path OUTSIDE the workspace: %s\n' "$bp" >&2
			    printf 'isolation_violation: %s\n' "$bp" >> "$CC_RUN_DIR/run.yaml" ;;
			*) : ;;
		esac
	done
fi

printf '\n[driver] coordinator turns complete; transcript=%s\n' "$CC_TRANSCRIPT"
printf '[driver] conversational verdict=%s\n' "$VERDICT"
rm -f "$CC_TRACE"
printf '[driver] Codex exposes no dependable file-read trace; dimension C will degrade\n'
exit 0
