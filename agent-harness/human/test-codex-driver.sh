#!/bin/sh
set -eu
. "$(git -C "$(dirname -- "$0")" rev-parse --show-toplevel)/test/lib/assert.sh"

lab=$(mktemp -d "${TMPDIR:-/tmp}/cc-codex-driver.XXXXXX")
trap 'rm -rf "$lab"' EXIT HUP INT TERM
mkdir -p "$lab/run" "$lab/workspace"

case_file="$lab/case.yaml"
simulator="$lab/human-simulator.md"
fake_codex="$lab/codex"
fake_log="$lab/codex.log"
codex_state="$lab/codex-state"

cat > "$case_file" <<'EOF'
id: codex-driver-contract
mode: conversation-only
human:
  persona: A person asking for help.
  turns:
    - say: "Please help me shape this idea."
    - on_offer_plan: "Yes, make that plan."
  visible_expectations:
    - The assistant responds naturally.
grader:
  post_conditions: []
EOF
printf 'Act as the tool-less human-simulator and judge only supplied text.\n' > "$simulator"

cat > "$fake_codex" <<'EOF'
#!/bin/sh
set -eu
out=
mode=coordinator
previous=
args=
saw_approve=0
saw_workspace_sandbox=0
for arg in "$@"; do
	args="$args [$arg]"
	if [ "$previous" = output ]; then out=$arg; previous=; continue; fi
	if [ "$previous" = sandbox ]; then
		[ "$arg" = workspace-write ] && saw_workspace_sandbox=1
		previous=
	fi
	[ "$arg" = "--output-last-message" ] && { previous=output; continue; }
	[ "$arg" = --sandbox ] && previous=sandbox
	[ "$arg" = --approve-for-me ] && saw_approve=1
	[ "$arg" = resume ] && mode=resume
	[ "$arg" = --ephemeral ] && mode=simulator
done
printf '%s\n' "$args" >> "$FAKE_CODEX_LOG"
[ "$saw_approve" -eq 0 ] || [ "$saw_workspace_sandbox" -eq 0 ] || {
	printf 'incompatible approval and sandbox flags\n' >&2
	exit 64
}
: "${out:?missing --output-last-message}"
case "$mode" in
	coordinator)
		mkdir -p "$FAKE_CODEX_STATE_DIR/sessions/test"
		cat > "$FAKE_CODEX_STATE_DIR/sessions/test/worker.jsonl" <<'JSON'
{"type":"session_meta","payload":{"source":{"subagent":{"thread_spawn":{"parent_thread_id":"11111111-1111-4111-8111-111111111111","agent_path":"/root/contract_worker"}}}}}
{"type":"turn_context","payload":{"model":"gpt-5.6-luna","effort":"medium"}}
JSON
		mkdir -p "$FAKE_CODEX_STATE_DIR/sessions/peer"
		cat > "$FAKE_CODEX_STATE_DIR/sessions/peer/thread.jsonl" <<'JSON'
{"type":"session_meta","payload":{"source":{"thread_spawn":{"parent_thread_id":"11111111-1111-4111-8111-111111111111","agent_path":"/root/peer_worker"}}}}
{"type":"turn_context","payload":{"model":"peer-model","effort":"low"}}
JSON
		printf '%s\n' 'I can outline a plan and get you started.' > "$out"
		printf '%s\n' '{"type":"thread.started","thread_id":"11111111-1111-4111-8111-111111111111"}'
		printf '%s\n' '{"type":"turn.completed","usage":{"input_tokens":120,"cached_input_tokens":20,"output_tokens":30}}' ;;
	resume)
		printf '%s\n' 'The plan is ready for your review.' > "$out"
		printf '%s\n' '{"type":"turn.completed","usage":{"input_tokens":180,"cached_input_tokens":40,"output_tokens":45}}' ;;
	simulator)
		printf '%s\n' 'PASS: The assistant responded naturally.' 'verdict: pass' > "$out" ;;
esac
EOF
chmod +x "$fake_codex"

FAKE_CODEX_LOG="$fake_log" \
FAKE_CODEX_STATE_DIR="$codex_state" \
CC_CODEX_BIN="$fake_codex" \
CC_CODEX_STATE_DIR="$codex_state" \
CC_RUN_DIR="$lab/run" \
CC_WORKSPACE="$lab/workspace" \
CC_BASELINE="$lab/baseline" \
CC_CASE_FILE="$case_file" \
CC_CASE_ID=codex-driver-contract \
CC_HOST=codex \
CC_MODE=conversation-only \
CC_HUMAN_SIM="$simulator" \
CC_TRANSCRIPT="$lab/run/transcript.txt" \
CC_TRACE="$lab/run/file-access-trace.tsv" \
CC_TELEMETRY="$lab/run/telemetry.tsv" \
CC_ROLE_EVIDENCE="$lab/run/role-evidence.tsv" \
	sh "$ROOT/agent-harness/human/drivers/codex.sh" >/dev/null

require_file "$lab/run/transcript.txt"
contains "$lab/run/transcript.txt" '# transcript: codex-driver-contract (host=codex)'
contains "$lab/run/transcript.txt" 'human: Please help me shape this idea.'
contains "$lab/run/transcript.txt" 'coordinator: I can outline a plan and get you started.'
contains "$lab/run/transcript.txt" 'human: Yes, make that plan.'
contains "$lab/run/transcript.txt" 'coordinator: The plan is ready for your review.'

require_file "$lab/run/coordinator-stream.jsonl"
contains "$lab/run/coordinator-stream.jsonl" '"type":"thread.started"'
require_file "$lab/run/telemetry.tsv"
contains "$lab/run/telemetry.tsv" 'orient'
contains "$lab/run/telemetry.tsv" 'create-plan'
require_file "$lab/run/conversational-verdict.txt"
contains "$lab/run/conversational-verdict.txt" 'verdict: pass'
require_file "$lab/run/role-evidence.tsv"
contains "$lab/run/role-evidence.tsv" 'worker'
contains "$lab/run/role-evidence.tsv" 'gpt-5.6-luna'
contains "$lab/run/role-evidence.tsv" 'medium'
not_contains "$lab/run/role-evidence.tsv" 'peer-model'
not_contains "$lab/run/role-evidence.tsv" 'peer_worker'
test ! -e "$lab/run/file-access-trace.tsv" || fail 'Codex driver must not claim a file-access trace'

contains "$fake_log" '[--approve-for-me]'
contains "$fake_log" '[resume]'
contains "$fake_log" '[11111111-1111-4111-8111-111111111111]'
contains "$fake_log" '[--sandbox] [read-only]'
contains "$fake_log" '[--ephemeral]'
not_contains "$fake_log" '[--dangerously-bypass-approvals-and-sandbox]'
assert_eq 3 "$(wc -l < "$fake_log" | tr -d ' ')"

pass 'Codex human-harness driver contract'
