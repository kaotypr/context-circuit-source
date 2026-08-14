#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

require_file() {
  test -f "$1" || fail "missing required file: $1"
}

require_dir() {
  test -d "$1" || fail "missing required directory: $1"
}

absent() {
  test ! -e "$1" || fail "legacy or forbidden path still exists: $1"
}

contains() {
  file=$1
  text=$2
  grep -F "$text" "$file" >/dev/null 2>&1 || fail "expected '$text' in $file"
}

require_file AGENTS.md
require_file WORKFLOW.md
require_file CLAUDE.md
require_file workspace.yaml
require_file context/INDEX.md
require_file context/PROJECT.md
require_file context/ARCHITECTURE.md
require_file context/CONVENTIONS.md
require_file context/DECISIONS.md
require_file docs/agent-workspace-workflow.md
require_file docs/runtime-contract.md
require_file agents/coordinator.md
require_file agents/repository-worker.md
require_file agents/reviewer.md
require_file .agents/skills/cc-session-entry/SKILL.md

contains AGENTS.md 'Preserve .runtime/'
contains WORKFLOW.md 'There is no single global current session.'
contains docs/agent-workspace-workflow.md 'Start or resume work in this workspace.'
contains docs/agent-workspace-workflow.md 'current-session.yaml'
contains docs/runtime-contract.md 'exclusive-create directory'
contains docs/runtime-contract.md 'handoff.md'
contains .agents/skills/cc-session-entry/SKILL.md 'If the packet or session record is absent,'
contains agents/coordinator.md 'Never use a global current-session or current-plan pointer.'
contains agents/repository-worker.md 'Work only in the assigned isolated worktree.'
contains agents/reviewer.md 'Remain read-only.'

# The rebuild removes the command layer instead of preserving a migration shim.
absent package.json
absent package-lock.json
absent tsconfig.json
absent scripts
absent test/agent-workspace.test.ts
absent test/simplified.test.ts
absent test/helpers.ts
absent .agents/bin
absent .codex/skills
absent .claude/commands
absent .github/workflows/sync-context-circuit-release.yml
absent docs/command-reference.md
absent docs/legacy-cleanup-gate.md

for file in README.md AGENTS.md CLAUDE.md WORKFLOW.md \
  agents/coordinator.md agents/repository-worker.md agents/reviewer.md \
  context/ARCHITECTURE.md context/CONVENTIONS.md context/DECISIONS.md \
  docs/agent-workspace-workflow.md docs/runtime-contract.md docs/development.md \
  docs/getting-started.md docs/planning.md docs/run-task.md; do
  if grep -E 'cc\.mjs|npm run|node --import|legacy-run-task|current command layer|compatibility adapter' "$file" >/dev/null 2>&1; then
    fail "command-layer language remains in normative file: $file"
  fi
done

fixture=$(mktemp -d "${TMPDIR:-/tmp}/context-circuit-acceptance.XXXXXX")
trap 'rm -rf "$fixture"' EXIT HUP INT TERM

runtime="$fixture/.runtime"
mkdir -p "$runtime/sessions/root-001" "$runtime/sessions/child-001"
mkdir -p "$runtime/sessions/verifier-001"
mkdir -p "$runtime/plans/plan-alpha/lease.lock" "$runtime/plans/plan-beta/lease.lock"
mkdir -p "$runtime/plans/plan-alpha/handoffs" "$runtime/plans/plan-beta/handoffs"
mkdir -p "$runtime/worktrees/context-circuit/plan-alpha"
mkdir -p "$runtime/worktrees/context-circuit/plan-beta"

printf '%s\n' \
  'schema_version: 1' \
  'session_id: root-001' \
  'parent_session_id: null' \
  'root_session_id: root-001' \
  'kind: root' \
  'role: coordinator' \
  'objective: Coordinate the workspace request' \
  'scope: workspace' \
  'write_access: false' \
  'status: executing' \
  'next_action: Delegate bounded work' \
  > "$runtime/sessions/root-001/session.yaml"

printf '%s\n' \
  'schema_version: 1' \
  'session_id: child-001' \
  'parent_session_id: root-001' \
  'root_session_id: root-001' \
  'kind: subagent' \
  'role: implementer' \
  'objective: Implement one bounded task' \
  'scope:' \
  '  plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
  '  task: AWF-0001' \
  '  paths:' \
  '    - AGENTS.md' \
  'write_access: true' \
  'worktree: .runtime/worktrees/context-circuit/plan-alpha' \
  'status: executing' \
  'next_action: Verify the assigned files' \
  > "$runtime/sessions/child-001/session.yaml"

printf '%s\n' \
  'schema_version: 1' \
  'session_id: child-001' \
  'parent_session_id: root-001' \
  'root_session_id: root-001' \
  'role: implementer' \
  'objective: Implement one bounded task' \
  'scope:' \
  '  plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
  '  task: AWF-0001' \
  '  paths:' \
  '    - AGENTS.md' \
  'non_goals:' \
  '  - Do not change plan status' \
  'context_refs:' \
  '  - AGENTS.md' \
  'repository: context-circuit' \
  'worktree: .runtime/worktrees/context-circuit/plan-alpha' \
  'permissions:' \
  '  write_worktree: true' \
  '  write_runtime_session: true' \
  '  write_plan: false' \
  '  write_activity: false' \
  'acceptance_criteria:' \
  '  - The assigned behavior is implemented' \
  'stop_conditions:' \
  '  - A required change falls outside scope' \
  'handoff_schema: session-handoff-v1' \
  > "$runtime/sessions/child-001/delegation.yaml"

printf '%s\n' \
  '# Session handoff' \
  '- session_id: root-001' \
  '- status: handoff' \
  '- next_action: Resume from this record' \
  > "$runtime/sessions/root-001/handoff.md"

printf '%s\n' \
  'schema_version: 1' \
  'session_id: verifier-001' \
  'parent_session_id: root-001' \
  'root_session_id: root-001' \
  'kind: subagent' \
  'role: verifier' \
  'objective: Verify the assigned worktree' \
  'write_access: false' \
  'worktree: .runtime/worktrees/context-circuit/plan-alpha' \
  'status: verifying' \
  > "$runtime/sessions/verifier-001/session.yaml"

printf '%s\n' \
  'permissions:' \
  '  write_worktree: false' \
  '  write_runtime_session: true' \
  '  write_plan: false' \
  '  write_activity: false' \
  'acceptance_criteria:' \
  '  - Reproduce the worker evidence' \
  'stop_conditions:' \
  '  - Required evidence is missing' \
  > "$runtime/sessions/verifier-001/delegation.yaml"

contains "$runtime/sessions/root-001/session.yaml" 'kind: root'
contains "$runtime/sessions/child-001/session.yaml" 'parent_session_id: root-001'
contains "$runtime/sessions/child-001/delegation.yaml" 'write_plan: false'
contains "$runtime/sessions/child-001/delegation.yaml" 'stop_conditions:'
contains "$runtime/sessions/root-001/handoff.md" 'Resume from this record'
contains "$runtime/sessions/verifier-001/session.yaml" 'write_access: false'
contains "$runtime/sessions/verifier-001/delegation.yaml" 'write_worktree: false'

printf '%s\n' \
  'schema_version: 1' \
  'plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
  'session_id: child-001' \
  'root_session_id: root-001' \
  'worktree: .runtime/worktrees/context-circuit/plan-alpha' \
  'status: active' \
  'acquired_at: 2026-08-14T12:00:00Z' \
  'heartbeat_at: 2026-08-14T12:10:00Z' \
  'stale_after_seconds: 1800' \
  > "$runtime/plans/plan-alpha/lease.yaml"

printf '%s\n' \
  'session_id: child-001' \
  'root_session_id: root-001' \
  'plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
  'worktree: .runtime/worktrees/context-circuit/plan-alpha' \
  > "$runtime/plans/plan-alpha/lease.lock/owner.yaml"

if mkdir "$runtime/plans/plan-alpha/lease.lock" 2>/dev/null; then
  fail 'same-plan lease contender acquired an existing lock'
fi
contains "$runtime/plans/plan-alpha/lease.yaml" 'session_id: child-001'

printf '%s\n' \
  'schema_version: 1' \
  'plan: plans/context-circuit-plans/0002-independent-plan' \
  'session_id: other-001' \
  'root_session_id: other-001' \
  'worktree: .runtime/worktrees/context-circuit/plan-beta' \
  'status: active' \
  > "$runtime/plans/plan-beta/lease.yaml"
contains "$runtime/plans/plan-beta/lease.yaml" 'worktree: .runtime/worktrees/context-circuit/plan-beta'
test "$runtime/worktrees/context-circuit/plan-alpha" != "$runtime/worktrees/context-circuit/plan-beta"
test ! -e "$runtime/current-session.yaml"
test ! -e "$runtime/current-plan.yaml"

printf '%s\n' \
  '# Session handoff' \
  '- session_id: child-001' \
  '- parent_session_id: root-001' \
  '- status: blocked' \
  '' \
  '## Evidence inspected' \
  '' \
  '- The worker stopped before verification.' \
  '' \
  '## Recommended next action' \
  '' \
  'Request explicit recovery or human takeover.' \
  > "$runtime/sessions/child-001/handoff.md"
contains "$runtime/sessions/child-001/handoff.md" 'status: blocked'
contains "$runtime/sessions/child-001/handoff.md" 'Request explicit recovery'

printf '%s\n' \
  'status: approved' \
  'tasks:' \
  '  - id: AWF-0001' \
  '    status: draft' \
  > "$fixture/plan.yaml"
printf '%s\n' \
  '# Verification handoff' \
  '- status: failed' \
  '- blocker: acceptance check failed' \
  '- next_action: Repair within scope or request a human decision' \
  > "$runtime/sessions/verifier-001/handoff.md"
contains "$runtime/sessions/verifier-001/handoff.md" 'status: failed'
contains "$runtime/sessions/verifier-001/handoff.md" 'human decision'
if grep -F 'status: done' "$fixture/plan.yaml" >/dev/null 2>&1; then
  fail 'failed verification was represented as a completed plan'
fi

printf '%s\n' \
  '# Contradictory source handoff' \
  '- status: blocked' \
  '- blocker: source documents contradict one another' \
  '- next_action: Request human-reviewed context refresh' \
  > "$runtime/sessions/root-001/contradiction-handoff.md"
contains "$runtime/sessions/root-001/contradiction-handoff.md" 'source documents contradict'
contains "$runtime/sessions/root-001/contradiction-handoff.md" 'human-reviewed context refresh'

contains docs/agent-workspace-workflow.md 'Verification fails;'
contains docs/agent-workspace-workflow.md 'Human approval is required'
contains docs/agent-workspace-workflow.md 'source changes during execution'
contains agents/reviewer.md 'Do not repair'
contains agents/reviewer.md 'Do not repair, change'

printf 'PASS: pure agent-workspace acceptance scenarios\n'
