#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

# Product host-native set (source extras are checked separately).
product_files="
.claude/agents/worker.md
.claude/agents/verifier.md
.claude/agents/planner.md
.claude/rules/role-tiering-spawn.md
.claude/rules/commit-convention.md
.codex/agents/worker.toml
.codex/agents/verifier.toml
.codex/agents/planner.toml
.cursor/agents/worker.md
.cursor/agents/verifier.md
.cursor/agents/planner.md
.cursor/rules/role-tiering-spawn.mdc
.cursor/rules/commit-convention.mdc
.cursor/rules/gh-unsandboxed.mdc
"

# --- product stubs exist in source and template, and stay aligned ---
for relpath in $product_files; do
	require_file "$ROOT/$relpath"
	require_file "$ROOT/template/$relpath"
	cmp -s "$ROOT/$relpath" "$ROOT/template/$relpath" || fail "source/template mismatch: $relpath"
done

test ! -e "$ROOT/template/.claude/agents/cc-human-simulator.md" || fail 'template shipped source-only cc-human-simulator'
test ! -e "$ROOT/template/.claude/skills/cc-test-case" || fail 'template shipped source-only cc-test-case'
test ! -e "$ROOT/.codex/rules" || fail 'source invented .codex/rules'
test ! -e "$ROOT/template/.codex/rules" || fail 'template invented .codex/rules'
require_file "$ROOT/.claude/agents/cc-human-simulator.md"

for role in worker verifier planner; do
	contains "$ROOT/.claude/agents/$role.md" ".context-circuit/agents/$role.md"
	contains "$ROOT/.cursor/agents/$role.md" ".context-circuit/agents/$role.md"
	contains "$ROOT/.codex/agents/$role.toml" ".context-circuit/agents/$role.md"
done

if grep -E '^# Worker role|^# Verifier role|^# Planner role' \
	"$ROOT/.claude/agents/worker.md" \
	"$ROOT/.claude/agents/verifier.md" \
	"$ROOT/.claude/agents/planner.md" \
	"$ROOT/.cursor/agents/worker.md" \
	"$ROOT/.cursor/agents/verifier.md" \
	"$ROOT/.cursor/agents/planner.md" >/dev/null 2>&1; then
	fail 'host agent stub copied a role-body heading'
fi

# --- Claude skill links are pointers to .agents/skills/cc-*, not copies ---
for skill_dir in "$ROOT"/.agents/skills/cc-*; do
	[ -d "$skill_dir" ] || continue
	name=${skill_dir##*/}
	test -e "$ROOT/.claude/skills/$name" || fail "missing Claude skill route: $name"
	test -L "$ROOT/template/.claude/skills/$name" || fail "missing template Claude skill route: $name"
	src_link=$(readlink "$ROOT/.claude/skills/$name") || fail "Claude skill $name is not a symlink"
	tpl_link=$(readlink "$ROOT/template/.claude/skills/$name") || fail "template Claude skill $name is not a symlink"
	[ "$src_link" = "$tpl_link" ] || fail "Claude skill link mismatch: $name"
	case "$src_link" in
		../../.agents/skills/"$name"|../../.agents/skills/"$name"/) ;;
		*) fail "Claude skill $name does not point at .agents/skills/$name (got $src_link)" ;;
	esac
done

# --- standing rules restate always-on clauses and cite the owner ---
contains "$ROOT/.claude/rules/commit-convention.md" "INV-COMMIT-01"
contains "$ROOT/.cursor/rules/commit-convention.mdc" "INV-COMMIT-01"
contains "$ROOT/.claude/rules/commit-convention.md" "Co-authored-by"
contains "$ROOT/.cursor/rules/commit-convention.mdc" "Co-authored-by"
contains "$ROOT/.claude/rules/commit-convention.md" "inspect the recorded message"
contains "$ROOT/.cursor/rules/commit-convention.mdc" "inspect the recorded message"
not_contains "$ROOT/.claude/rules/commit-convention.md" "Do not restate the rule here"
not_contains "$ROOT/.cursor/rules/commit-convention.mdc" "Do not restate the rule here"
contains "$ROOT/.claude/rules/role-tiering-spawn.md" "CLAUDE.md"
contains "$ROOT/.cursor/rules/role-tiering-spawn.mdc" "CURSOR.md"
contains "$ROOT/.cursor/rules/role-tiering-spawn.mdc" "Do not import"

# --- host-blocked fail-closed language is unchanged in adapters and coordinator ---
for f in \
	"$ROOT/.context-circuit/wrapper/adapters/AGENTS.md" \
	"$ROOT/.context-circuit/wrapper/adapters/CLAUDE.md" \
	"$ROOT/.context-circuit/wrapper/adapters/CURSOR.md" \
	"$ROOT/.context-circuit/agents/coordinator.md"; do
	contains "$f" "host-blocked"
done

# --- no credentials, transcripts, or auth state in shipped host folders ---
host_trees="$ROOT/.claude $ROOT/.codex $ROOT/.cursor $ROOT/template/.claude $ROOT/template/.codex $ROOT/template/.cursor"
if grep -REn '^[[:space:]]*(password|api_key|access_token|provider_payload|transcript):' $host_trees >/dev/null 2>&1; then
	fail 'host-native folders contained credential or transcript keys'
fi
if grep -RIlF 'BEGIN RSA PRIVATE KEY' $host_trees >/dev/null 2>&1; then
	fail 'host-native folders contained a private key block'
fi
test ! -e "$ROOT/.claude/settings.local.json"
test ! -e "$ROOT/template/.claude/settings.local.json"
test ! -e "$ROOT/AGENTS.override.md"

# --- shipped CLAUDE.md no longer treats .claude/ as unshipped host-local ---
not_contains "$ROOT/.context-circuit/wrapper/adapters/CLAUDE.md" 'never part of workspace or shipped state'

pass 'host-native routes stay pointers'
