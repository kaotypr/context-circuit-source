#!/bin/sh
# Hidden ok easter-egg: engine stdout, coordinator trigger, punchline hiding.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

ENG="$ROOT/.context-circuit/wrapper/runtime/engine.sh"
COORD="$ROOT/.context-circuit/agents/coordinator.md"
PUNCHLINE='Haha, meat proxy!'
PLOT="$ROOT/agent-harness/conversations/plots/lone-ok-easter-egg.md"
CASE="$ROOT/agent-harness/conversations/generated/lone-ok-easter-egg.case.yaml"
SEAM="$ROOT/sources/system-design/context-circuit/v1.0.0/core/engine-and-seam.md"

case_name=all
if [ "${1:-}" = --case ]; then case_name=${2:-}; fi

eng() { sh "$ENG" "$@"; }

assert_exact_stdout() {
	got=$(
		eng "$@" || { printf 'engine-exit:%s\n' "$?" >&2; exit 1; }
		printf x
	) || fail "engine $* failed"
	want=$(printf '%s\nx' "$PUNCHLINE")
	assert_eq "$want" "$got"
}

case_engine() {
	require_file "$ENG"
	contains "$ENG" "cc_ok_easter_egg"
	contains "$ENG" "ok-easter-egg"

	assert_exact_stdout ok-easter-egg
	# extra argv and stdin are ignored — the engine is not a conversation router
	assert_exact_stdout ok-easter-egg "user said ok"
	got=$(
		printf 'ok\nplease route this\n' | eng ok-easter-egg \
			|| { printf 'engine-exit:%s\n' "$?" >&2; exit 1; }
		printf x
	) || fail "engine ok-easter-egg (stdin) failed"
	want=$(printf '%s\nx' "$PUNCHLINE")
	assert_eq "$want" "$got"

	eng ok-easter-egg >/dev/null
	assert_eq "0" "$?"

	for help_cmd in help -h --help ''; do
		if [ -n "$help_cmd" ]; then
			help_out=$(eng "$help_cmd") || fail "engine help ($help_cmd) failed"
		else
			help_out=$(eng) || fail "engine (no args) failed"
		fi
		printf '%s\n' "$help_out" | grep -F "$PUNCHLINE" >/dev/null 2>&1 \
			&& fail "punchline leaked into engine help ($help_cmd)" || :
		printf '%s\n' "$help_out" | grep -F "ok-easter-egg" >/dev/null 2>&1 \
			&& fail "ok-easter-egg listed in engine help ($help_cmd)" || :
	done
}

case_coordinator() {
	require_file "$COORD"
	contains "$COORD" "trimmed and case-folded"
	contains "$COORD" "equals exactly \`ok\`"
	contains "$COORD" "sh .context-circuit/wrapper/runtime/engine.sh ok-easter-egg"
	contains "$COORD" "Prefix the user-visible reply with that command's stdout"
	contains "$COORD" "immediately pending coordinator ask"
	contains "$COORD" "If nothing is pending"
	contains "$COORD" "wait"
	contains "$COORD" "not a new approval, execute, complete, or deliver token"
	contains "$COORD" "okay"
	contains "$COORD" "ok."
	contains "$COORD" "ok!"
	not_contains "$COORD" "$PUNCHLINE"

	for role in worker verifier planner; do
		require_file "$ROOT/.context-circuit/agents/$role.md"
		not_contains "$ROOT/.context-circuit/agents/$role.md" "ok-easter-egg"
		not_contains "$ROOT/.context-circuit/agents/$role.md" "$PUNCHLINE"
	done
}

assert_no_punchline_tree() {
	tree=$1
	label=$2
	[ -e "$tree" ] || fail "missing surface: $label ($tree)"
	hits=$(grep -r -F -l -- "$PUNCHLINE" "$tree" 2>/dev/null || true)
	[ -z "$hits" ] || fail "punchline leaked into $label: $hits"
}

case_hidden() {
	require_file "$SEAM"
	not_contains "$SEAM" "$PUNCHLINE"
	not_contains "$SEAM" "ok-easter-egg"

	assert_no_punchline_tree "$ROOT/.context-circuit/agents" "role packets"
	assert_no_punchline_tree "$ROOT/product/.agents/skills" "skills"
	assert_no_punchline_tree "$ROOT/.context-circuit/wrapper/adapters" "host adapters"
	assert_no_punchline_tree "$ROOT/.context-circuit/docs" "product docs"
	assert_no_punchline_tree "$ROOT/context" "Product Knowledge"
	assert_no_punchline_tree "$ROOT/product/.claude" "Claude host routes"
	assert_no_punchline_tree "$ROOT/product/.codex" "Codex host routes"
	assert_no_punchline_tree "$ROOT/product/.cursor" "Cursor host routes"

	if [ -d "$ROOT/template/.agents/skills" ]; then
		assert_no_punchline_tree "$ROOT/template/.agents/skills" "template skills"
	fi
	if [ -d "$ROOT/template/context" ]; then
		assert_no_punchline_tree "$ROOT/template/context" "template Product Knowledge"
	fi

	# engine implementation is the allowed home of the punchline
	contains "$ENG" "$PUNCHLINE"
}

case_harness() {
	require_file "$PLOT"
	require_file "$CASE"
	contains "$PLOT" "id: lone-ok-easter-egg"
	contains "$PLOT" "AC-LONE-OK"
	contains "$PLOT" "AC-NO-FALSE-TRIGGER"
	contains "$PLOT" "AC-IDLE-OK"
	contains "$PLOT" "okay"
	contains "$PLOT" "ok."
	contains "$PLOT" "ok!"
	contains "$PLOT" "pending"
	contains "$CASE" "generated from plots/lone-ok-easter-egg.md"
	contains "$CASE" "id: lone-ok-easter-egg"
	contains "$CASE" "AC-LONE-OK"
	contains "$CASE" "AC-NO-FALSE-TRIGGER"
	contains "$CASE" "AC-IDLE-OK"
}

run_case() {
	case "$1" in
		engine) case_engine ;;
		coordinator) case_coordinator ;;
		hidden) case_hidden ;;
		harness) case_harness ;;
		*) fail "unknown case $1" ;;
	esac
}

if [ "$case_name" = all ]; then
	run_case engine
	run_case coordinator
	run_case hidden
	run_case harness
else
	run_case "$case_name"
fi

pass 'ok easter-egg'
