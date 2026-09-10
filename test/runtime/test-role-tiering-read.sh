#!/bin/sh
# role-tiering-read verb: local override, committed fallback, and missing-both.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

W="$ROOT/.context-circuit/wrapper"
ENG="$W/runtime/engine.sh"

eng() { sh "$ENG" "$@"; }

expect_code() {
	ec_want=$1
	shift
	ec_err=$( { "$@" >/tmp/cc-rtr-out.$$; } 2>&1 || true )
	printf '%s\n' "$ec_err" | grep -Fq "$ec_want" || fail "expected '$ec_want' in stderr, got: $ec_err"
	rm -f /tmp/cc-rtr-out.$$
}

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- neither file exists: ROLE_TIERING_MISSING, no invented output ---
rm -f "$ws/role-tiering.local.yaml"
rm -f "$ws/.context-circuit/role-tiering.fallback.yaml"
expect_code ROLE_TIERING_MISSING eng role-tiering-read "$ws"

# --- only the committed fallback exists: source: fallback, body echoed ---
mkdir -p "$ws/.context-circuit"
cp "$ROOT/.context-circuit/role-tiering.fallback.yaml" "$ws/.context-circuit/role-tiering.fallback.yaml"
out=$(eng role-tiering-read "$ws")
printf '%s\n' "$out" | head -1 | grep -Fxq 'source: fallback' || fail "fallback: wrong source line"
printf '%s\n' "$out" | grep -Fq "path: $ws/.context-circuit/role-tiering.fallback.yaml" || fail "fallback: wrong path line"
printf '%s\n' "$out" | grep -Fxq -- '---' || fail "fallback: missing separator"
printf '%s\n' "$out" | grep -q '^hosts:' || fail "fallback: body not echoed"

# --- local override present: source: local, wins over fallback, body echoed ---
printf 'hosts:\n  claude-code:\n    worker:\n      model: test-model\n      effort: high\n' \
	>"$ws/role-tiering.local.yaml"
out=$(eng role-tiering-read "$ws")
printf '%s\n' "$out" | head -1 | grep -Fxq 'source: local' || fail "local: wrong source line"
printf '%s\n' "$out" | grep -Fq "path: $ws/role-tiering.local.yaml" || fail "local: wrong path line"
printf '%s\n' "$out" | grep -Fq 'model: test-model' || fail "local: body not echoed"

rm -f "$ws/role-tiering.local.yaml"
rm -f "$ws/.context-circuit/role-tiering.fallback.yaml"

trap - EXIT HUP INT TERM
rm -rf "$ws"

# --- print-only boundary: no model id, no host-group selection in engine.sh ---
not_contains "$W/runtime/engine.sh" "claude-opus"
not_contains "$W/runtime/engine.sh" "gpt-5"
not_contains "$W/runtime/engine.sh" "composer-"
not_contains "$W/runtime/engine.sh" "cursor-grok"
not_contains "$W/runtime/engine.sh" "hosts.claude-code"
not_contains "$W/runtime/engine.sh" "hosts.codex"
not_contains "$W/runtime/engine.sh" "hosts.cursor-agent"

pass 'role-tiering-read verb: local, fallback, missing'
