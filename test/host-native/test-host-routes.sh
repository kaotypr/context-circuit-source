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

# --- product discovery exists only in packaging, not source root or seed ---
for relpath in $product_files; do
	require_file "$ROOT/product/$relpath"
	test ! -e "$ROOT/$relpath" || fail "product route active in source: $relpath"
	test ! -e "$ROOT/template/$relpath" || fail "product route duplicated in seed: $relpath"
done
require_file "$ROOT/.claude/agents/cc-human-simulator.md"
require_file "$ROOT/.claude/skills/cc-test-case/SKILL.md"
require_file "$ROOT/.agents/skills/cc-source-develop/SKILL.md"
for skill_dir in "$ROOT"/.agents/skills/*; do
	[ -d "$skill_dir" ] || continue
	[ "${skill_dir##*/}" = cc-source-develop ] || fail "unexpected source skill: $skill_dir"
done
for host in .claude .codex .cursor; do
	test ! -e "$ROOT/template/$host" || fail "host integration remains in seed: $host"
done
test ! -e "$ROOT/product/.codex/rules" || fail 'invented .codex/rules'
test ! -e "$ROOT/product/.claude/agents/cc-human-simulator.md" || fail 'source agent packaged'
test ! -e "$ROOT/product/.claude/skills/cc-test-case" || fail 'source skill packaged'

for role in worker verifier planner; do
	contains "$ROOT/product/.claude/agents/$role.md" ".context-circuit/agents/$role.md"
	contains "$ROOT/product/.cursor/agents/$role.md" ".context-circuit/agents/$role.md"
	contains "$ROOT/product/.codex/agents/$role.toml" ".context-circuit/agents/$role.md"
done

if grep -E '^# Worker role|^# Verifier role|^# Planner role' \
	"$ROOT/product/.claude/agents/worker.md" \
	"$ROOT/product/.claude/agents/verifier.md" \
	"$ROOT/product/.claude/agents/planner.md" \
	"$ROOT/product/.cursor/agents/worker.md" \
	"$ROOT/product/.cursor/agents/verifier.md" \
	"$ROOT/product/.cursor/agents/planner.md" >/dev/null 2>&1; then
	fail 'host agent stub copied a role-body heading'
fi

# --- Claude skill stubs route to .agents/skills/cc-*, not copies ---
for skill_dir in "$ROOT/product"/.agents/skills/cc-*; do
	[ -d "$skill_dir" ] || continue
	name=${skill_dir##*/}
	test ! -e "$ROOT/.agents/skills/$name" || fail "product skill active in source: $name"
	test ! -e "$ROOT/.claude/skills/$name" || fail "product slash skill active in source: $name"
	src="$ROOT/product/.claude/skills/$name/SKILL.md"
	require_file "$src"
	test ! -L "$ROOT/product/.claude/skills/$name" || fail "Claude skill $name is still a symlink"
	contains "$src" ".agents/skills/$name/SKILL.md"
	contains "$src" "Do not elaborate a second skill policy here"
	if grep -E '^## ' "$src" >/dev/null 2>&1; then
		fail "Claude skill stub $name copied a skill-body heading"
	fi
done

# --- standing rules restate always-on clauses and cite the owner ---
contains "$ROOT/product/.claude/rules/commit-convention.md" "INV-COMMIT-01"
contains "$ROOT/product/.cursor/rules/commit-convention.mdc" "INV-COMMIT-01"
contains "$ROOT/product/.claude/rules/commit-convention.md" "Co-authored-by"
contains "$ROOT/product/.cursor/rules/commit-convention.mdc" "Co-authored-by"
contains "$ROOT/product/.claude/rules/commit-convention.md" "inspect the recorded message"
contains "$ROOT/product/.cursor/rules/commit-convention.mdc" "inspect the recorded message"
not_contains "$ROOT/product/.claude/rules/commit-convention.md" "Do not restate the rule here"
not_contains "$ROOT/product/.cursor/rules/commit-convention.mdc" "Do not restate the rule here"
contains "$ROOT/product/.claude/rules/role-tiering-spawn.md" "CLAUDE.md"
contains "$ROOT/product/.cursor/rules/role-tiering-spawn.mdc" "CURSOR.md"
contains "$ROOT/product/.cursor/rules/role-tiering-spawn.mdc" "Do not import"

# --- host-blocked fail-closed language is unchanged in adapters and coordinator ---
for f in \
	"$ROOT/.context-circuit/wrapper/adapters/AGENTS.md" \
	"$ROOT/.context-circuit/wrapper/adapters/CLAUDE.md" \
	"$ROOT/.context-circuit/wrapper/adapters/CURSOR.md" \
	"$ROOT/.context-circuit/agents/coordinator.md"; do
	contains "$f" "host-blocked"
done

# --- no credentials, transcripts, or auth state in shipped host folders ---
host_trees="$ROOT/product/.claude $ROOT/product/.codex $ROOT/product/.cursor"
if grep -REn '^[[:space:]]*(password|api_key|access_token|provider_payload|transcript):' $host_trees >/dev/null 2>&1; then
	fail 'host-native folders contained credential or transcript keys'
fi
if grep -RIlF 'BEGIN RSA PRIVATE KEY' $host_trees >/dev/null 2>&1; then
	fail 'host-native folders contained a private key block'
fi
test ! -e "$ROOT/product/.claude/settings.local.json"
test ! -e "$ROOT/template/.claude/settings.local.json"
test ! -e "$ROOT/AGENTS.override.md"

# --- shipped CLAUDE.md no longer treats .claude/ as unshipped host-local ---
not_contains "$ROOT/.context-circuit/wrapper/adapters/CLAUDE.md" 'never part of workspace or shipped state'

pass 'host-native routes stay pointers'
