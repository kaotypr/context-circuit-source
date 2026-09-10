#!/bin/sh
# Repository grounding (INV-GROUND-01/02/03): deterministic discovery + manifest
# (rich / doc-less / greenfield), worktree hardening's environment field, the
# grounding directive, brief assembly + the preflight refusal, and the
# friction->proposal contract. Proven here because a live model cannot be made to
# omit a slot or ship a doc-less repo on command.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# the shipped brief template is promoted beside the runtime on release; mirror
# that here so cc_worker_brief_assemble can find it in the fixture.
mkdir -p "$ws/.context-circuit/wrapper/runtime"
cp "$ROOT/.context-circuit/wrapper/adapters/worker-brief.md" "$ws/.context-circuit/wrapper/runtime/worker-brief.md"

# --- a rich repo: agent guidance + a skill + cursor rules + a lockfile ---
cc_fx_repo "$ws" widgets development
wr="$ws/repositories/widgets"
printf '# Widgets agent guide\n\nEvery new source file must begin with `// @grounded`.\n' >"$wr/AGENTS.md"
printf '# Claude\n\nSee AGENTS.md.\n' >"$wr/CLAUDE.md"
mkdir -p "$wr/.agents/skills/widget-style" "$wr/.cursor/rules"
printf -- '---\nname: widget-style\ndescription: DLS components via the private registry\n---\n\n# style\n' >"$wr/.agents/skills/widget-style/SKILL.md"
printf 'always use tabs\n' >"$wr/.cursor/rules/base.md"
printf '{ "name": "widgets" }\n' >"$wr/package-lock.json"
printf 'node_modules/\n' >"$wr/.gitignore"
mkdir -p "$wr/node_modules/widget"
printf 'prepared\n' >"$wr/node_modules/widget/index.js"
cc_toolchain_stamp_write "$wr/node_modules" "$(cc_digest "$wr/package-lock.json")"
git -C "$wr" add -A && git -C "$wr" commit -q -m 'chore: add agent guidance + toolchain'

cc_fx_plan_ex "$ws" 0001-widget "Widget" widgets src/widget ""
exec=$(cc_execution_begin "$ws" 0001-widget sess1 | sed -n 's/^execution_id: //p')
edir="$ws/.runtime/executions/0001-widget/$exec"

# 1. discovery recorded the manifest as execution evidence
mf="$edir/grounding/widgets.yaml"
require_file "$mf"
contains "$mf" "schema_version: 1"
contains "$mf" "- AGENTS.md"
contains "$mf" "- CLAUDE.md"
contains "$mf" "- .cursor/rules/"
contains "$mf" "name: widget-style"
contains "$mf" "description: DLS components via the private registry"
contains "$mf" "environment: ready"          # the overlaid install tree was prepared

# 2. the grounding directive renders the discovered guidance (non-empty variant)
dir=$(cc_grounding_directive "$mf")
printf '%s' "$dir" | grep -q "read and apply this repository" || fail "directive missing read-and-apply"
printf '%s' "$dir" | grep -q "AGENTS.md — read and honor it" || fail "directive missing AGENTS.md"
printf '%s' "$dir" | grep -q "widget-style — DLS components" || fail "directive missing the skill menu entry"
printf '%s' "$dir" | grep -q "STOP and report" || fail "directive missing the precedence/conflict stop"

# 3. brief assembly fills every slot and passes the preflight
cc_worker_brief_assemble "$ws" "$edir" widgets "Add the widget module." >/dev/null
brief="$edir/brief-widgets.md"
require_file "$brief"
contains "$brief" "Worker brief — plan 0001-widget"
contains "$brief" "## Repository grounding"
contains "$brief" "AGENTS.md — read and honor it"
not_contains "$brief" "@@GROUNDING@@"          # the slot is filled, not left raw
not_contains "$brief" "@@ENVIRONMENT@@"
contains "$brief" "Add the widget module."     # the coordinator's task focus
cc_brief_preflight "$brief" >/dev/null

# --- a doc-less repo: code, but no agent guidance ---
cc_fx_repo "$ws" plain development
cc_fx_plan_ex "$ws" 0002-plain "Plain" plain src/plain ""
exec2=$(cc_execution_begin "$ws" 0002-plain sess2 | sed -n 's/^execution_id: //p')
edir2="$ws/.runtime/executions/0002-plain/$exec2"
mf2="$edir2/grounding/plain.yaml"
require_file "$mf2"
contains "$mf2" "files: []"
contains "$mf2" "skills: []"
contains "$mf2" "environment: no-toolchain"    # no lockfile in the seed repo
dir2=$(cc_grounding_directive "$mf2")
printf '%s' "$dir2" | grep -q "No repository agent guidance was discovered" || fail "doc-less directive wrong"
# the brief still carries the (empty-variant) grounding section and passes preflight
cc_worker_brief_assemble "$ws" "$edir2" plain "Add the plain module." >/dev/null
contains "$edir2/brief-plain.md" "## Repository grounding"
cc_brief_preflight "$edir2/brief-plain.md" >/dev/null

# --- hardening provisions a toolchain vs preserving a greenfield directly ---
cc_harden_worktree "$wr" "$wr" | grep -q "environment: ready" || fail "prepared lockfile repo should harden ready"
cc_harden_worktree "$ws/repositories/plain" "$ws/repositories/plain" | grep -q "environment: no-toolchain" || fail "toolchain-less repo should be no-toolchain"

# --- preflight refuses a brief missing / with an unfilled grounding slot ---
printf '# Worker brief\n\nno grounding here\n' >"$ws/bad-brief.md"
expect_failure cc_brief_preflight "$ws/bad-brief.md"
printf '# Worker brief\n\n## Repository grounding\n\n@@GROUNDING@@\n' >"$ws/unfilled-brief.md"
expect_failure cc_brief_preflight "$ws/unfilled-brief.md"

# --- contracts: precedence + friction->proposal are owned, not duplicated ---
contains "$ROOT/.context-circuit/wrapper/contracts/invariants.yaml" "INV-GROUND-02"
contains "$ROOT/.context-circuit/wrapper/contracts/schemas/worker-handoff.yaml" "repository_friction"

# --- the runtime discovery emits DATA, not a model prompt (INV-RUNTIME-01) ---
not_contains "$ROOT/.context-circuit/wrapper/runtime/engine.sh" "claude -p"

pass 'repository grounding'
