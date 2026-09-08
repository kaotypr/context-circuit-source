#!/bin/sh
# Member roster schema bootstrap and runtime validate/resolve.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

W="$ROOT/.context-circuit/wrapper"
ENG="$W/runtime/engine.sh"
case_name=all
if [ "${1:-}" = --case ]; then case_name=${2:-}; fi

eng() { sh "$ENG" "$@"; }

expect_code() {
	ec_want=$1
	shift
	ec_err=$( { "$@" >/tmp/cc-mr-out.$$; } 2>&1 || true )
	printf '%s\n' "$ec_err" | grep -Fq "$ec_want" || fail "expected '$ec_want' in stderr, got: $ec_err"
	if grep -E '^(intent_start|intent_end|plan_start|plan_end):' /tmp/cc-mr-out.$$ >/dev/null 2>&1; then
		fail "failure invented a band: $(cat /tmp/cc-mr-out.$$)"
	fi
	rm -f /tmp/cc-mr-out.$$
}

case_schema_bootstrap() {
	require_file "$ROOT/members.yaml"
	require_file "$ROOT/template/members.yaml"
	require_file "$W/contracts/schemas/members.yaml"
	require_file "$W/contracts/schemas/member-local.yaml"
	contains "$ROOT/members.yaml" "kao:"
	contains "$ROOT/members.yaml" "intent_band_start: 1"
	contains "$ROOT/members.yaml" "intent_band_end: 99"
	contains "$ROOT/members.yaml" "plan_band_start: 1"
	contains "$ROOT/members.yaml" "plan_band_end: 999"
	contains "$ROOT/.gitignore" "/member.local.yaml"
	contains "$ROOT/.gitignore" "/template/member.local.yaml"
	contains "$ROOT/.gitignore" "Gitignore is not a read block"
	contains "$ROOT/template/.gitignore" "/member.local.yaml"
	contains "$W/manifest.yaml" "member.local.yaml"
	contains "$W/manifest.yaml" "members: [1]"
	contains "$W/manifest.yaml" "member-local: [1]"
	contains "$W/contracts/invariants.yaml" "member_roster:"
	contains "$W/contracts/invariants.yaml" "member_identity:"
	contains "$W/contracts/schemas/members.yaml" "MEMBER_BAND_OVERLAP"
	contains "$W/runtime/engine.sh" "member-roster-validate"
	contains "$W/runtime/engine.sh" "member-identity-read"
	contains "$W/runtime/engine.sh" "member-band-resolve"
}

case_runtime_validate() {
	ws=$(cc_fx_ws)
	trap 'rm -rf "$ws"' EXIT HUP INT TERM

	# init seeds an empty roster and must not write member.local.yaml
	require_file "$ws/members.yaml"
	test ! -e "$ws/member.local.yaml" || fail "workspace-init wrote member.local.yaml"
	eng member-roster-validate "$ws" >/dev/null

	# overlapping intent bands fail closed
	printf '%s\n' \
		'schema_version: 1' \
		'members:' \
		'  a:' \
		'    display_name: a' \
		'    intent_band_start: 1' \
		'    intent_band_end: 50' \
		'    plan_band_start: 1' \
		'    plan_band_end: 100' \
		'  b:' \
		'    display_name: b' \
		'    intent_band_start: 50' \
		'    intent_band_end: 99' \
		'    plan_band_start: 101' \
		'    plan_band_end: 200' \
		>"$ws/members.yaml"
	expect_code MEMBER_BAND_OVERLAP eng member-roster-validate "$ws"

	# inverted bounds
	cc_fx_roster "$ws" kao 99 1 1 999
	expect_code MEMBER_BAND_INVALID eng member-roster-validate "$ws"

	# malformed roster (missing fields)
	printf 'schema_version: 1\nmembers:\n  kao:\n    display_name: kao\n' >"$ws/members.yaml"
	expect_code MEMBER_ROSTER_INVALID eng member-roster-validate "$ws"

	# missing roster
	rm -f "$ws/members.yaml"
	expect_code MEMBER_ROSTER_MISSING eng member-roster-validate "$ws"

	# valid roster + missing identity: stable error, no invented band
	cc_fx_roster "$ws" kao 1 99 1 999
	rm -f "$ws/member.local.yaml"
	expect_code MEMBER_IDENTITY_MISSING eng member-identity-read "$ws"
	expect_code MEMBER_IDENTITY_MISSING eng member-band-resolve "$ws"
	# extra numeric args must not invent a band
	expect_code MEMBER_IDENTITY_MISSING eng member-identity-read "$ws" 10 20
	expect_code MEMBER_IDENTITY_MISSING eng member-band-resolve "$ws" 10 20

	# unknown member
	cc_fx_member_identity "$ws" ghost
	expect_code MEMBER_IDENTITY_UNKNOWN eng member-identity-read "$ws"
	expect_code MEMBER_IDENTITY_UNKNOWN eng member-band-resolve "$ws"

	# valid resolve emits the roster bands, not CLI numbers
	cc_fx_member_identity "$ws" kao
	out=$(eng member-band-resolve "$ws" 50 60)
	printf '%s\n' "$out" | grep -Fq 'member: kao' || fail "resolve missing member"
	printf '%s\n' "$out" | grep -Fq 'intent_start: 1' || fail "resolve missing intent_start"
	printf '%s\n' "$out" | grep -Fq 'intent_end: 99' || fail "resolve missing intent_end"
	printf '%s\n' "$out" | grep -Fq 'plan_start: 1' || fail "resolve missing plan_start"
	printf '%s\n' "$out" | grep -Fq 'plan_end: 999' || fail "resolve missing plan_end"
	printf '%s\n' "$out" | grep -Fq 'intent_start: 50' && fail "CLI range was accepted as a band" || :

	idout=$(eng member-identity-read "$ws")
	printf '%s\n' "$idout" | grep -Fq 'member: kao' || fail "identity-read missing member"

	eng member-roster-validate "$ws" >/dev/null

	trap - EXIT HUP INT TERM
	rm -rf "$ws"
}

case_guidance() {
	contains "$ROOT/.agents/skills/cc-workspace/SKILL.md" "member.local.yaml"
	contains "$ROOT/.agents/skills/cc-workspace/SKILL.md" "members.yaml"
	contains "$ROOT/.agents/skills/cc-workspace/SKILL.md" "numeric range"
	contains "$ROOT/.agents/skills/cc-intent/SKILL.md" "member-band-resolve"
	contains "$ROOT/.agents/skills/cc-intent/SKILL.md" "block number"
	contains "$ROOT/.context-circuit/agents/coordinator.md" "member identity"
	contains "$ROOT/.context-circuit/agents/coordinator.md" "block number"
	contains "$ROOT/.context-circuit/docs/getting-started.md" "member"
	contains "$ROOT/.context-circuit/docs/getting-started.md" "host-local"
}

run_case() {
	case "$1" in
		schema-bootstrap) case_schema_bootstrap ;;
		runtime-validate) case_runtime_validate ;;
		guidance) case_guidance ;;
		*) fail "unknown case $1" ;;
	esac
}

if [ "$case_name" = all ]; then
	run_case schema-bootstrap
	run_case runtime-validate
	run_case guidance
else
	run_case "$case_name"
fi

pass 'member roster'
