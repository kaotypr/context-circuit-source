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

absent() {
  test ! -e "$1" || fail "legacy or forbidden path still exists: $1"
}

contains() {
  file=$1
  text=$2
  grep -F "$text" "$file" >/dev/null 2>&1 || fail "expected '$text' in $file"
}

expect_success() {
  if "$@"; then
    return 0
  fi
  fail "expected command to succeed: $*"
}

expect_failure() {
  if "$@"; then
    fail "expected command to fail: $*"
  fi
}

atomic_write() {
  target=$1
  shift
  temporary="${target}.tmp.$$"
  printf '%s\n' "$@" > "$temporary"
  mv "$temporary" "$target"
}

safe_id() {
  candidate=$1
  case "$candidate" in
    ''|*[!a-z0-9._-]*) return 1 ;;
  esac
  case "$candidate" in
    [a-z0-9]*) test "${#candidate}" -le 64 ;;
    *) return 1 ;;
  esac
}

scope_allows() {
  candidate=$1
  assigned_worktree=$2
  case "$candidate/" in
    "$assigned_worktree/"*) return 0 ;;
    *) return 1 ;;
  esac
}

verifier_write() {
  target=$1
  own_session_dir=$2
  shift 2
  case "$target/" in
    "$own_session_dir/"*) atomic_write "$target" "$@" ;;
    *)
      printf 'BLOCKED: verifier write is outside own session handoff\n' >&2
      return 1
      ;;
  esac
}

assert_failure_reason() {
  failure_log=$1
  expected_reason=$2
  shift 2
  set +e
  "$@" > "$failure_log" 2>&1
  failure_status=$?
  set -e
  test "$failure_status" -ne 0 || fail "fixture unexpectedly succeeded: $*"
  contains "$failure_log" "$expected_reason"
}

lease_guard() {
  guard_lock=$1
  if test -e "$guard_lock"; then
    printf 'BLOCKED: same-plan lease is already owned\n' >&2
    return 1
  fi
  return 0
}

scope_guard() {
  guard_path=$1
  guard_worktree=$2
  if scope_allows "$guard_path" "$guard_worktree"; then
    return 0
  fi
  printf 'BLOCKED: path is outside assigned worktree\n' >&2
  return 1
}

verification_guard() {
  guard_handoff=$1
  if grep -F 'status: failed' "$guard_handoff" >/dev/null 2>&1; then
    printf 'BLOCKED: verification failed; human decision required\n' >&2
    return 1
  fi
  return 0
}

handoff_status_valid() {
  case "$1" in
    completed|blocked|failed|awaiting-human-gate) return 0 ;;
    *) return 1 ;;
  esac
}

write_session() {
  session_id=$1
  parent_id=$2
  root_id=$3
  kind=$4
  role=$5
  objective=$6
  task_id=$7
  scope=$8
  write_access=$9
  worktree=${10}
  session_status=${11}
  session_dir="$runtime/sessions/$session_id"
  mkdir -p "$session_dir"
  atomic_write "$session_dir/session.yaml" \
    'schema_version: 1' \
    "session_id: $session_id" \
    "parent_session_id: $parent_id" \
    "root_session_id: $root_id" \
    "kind: $kind" \
    "role: $role" \
    "objective: $objective" \
    "plan: plans/context-circuit-plans/0001-agent-workspace-workflow" \
    "task: $task_id" \
    "scope: $scope" \
    'context_refs:' \
    '  - AGENTS.md' \
    '  - WORKFLOW.md' \
    '  - docs/runtime-contract.md' \
    'non_goals:' \
    '  - Do not change canonical plan or task status' \
    "write_access: $write_access" \
    "worktree: $worktree" \
    "status: $session_status" \
    'created_at: 2026-08-14T18:00:00Z' \
    'updated_at: 2026-08-14T18:00:00Z' \
    'next_action: Continue only within the recorded scope' \
    'blockers: []'
}

write_delegation() {
  session_id=$1
  parent_id=$2
  root_id=$3
  role=$4
  objective=$5
  task_id=$6
  assigned_paths=$7
  assigned_worktree=$8
  write_worktree=$9
  session_dir="$runtime/sessions/$session_id"
  atomic_write "$session_dir/delegation.yaml" \
    'schema_version: 1' \
    "session_id: $session_id" \
    "parent_session_id: $parent_id" \
    "root_session_id: $root_id" \
    "role: $role" \
    "objective: $objective" \
    'scope:' \
    '  plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
    "  task: $task_id" \
    '  paths:' \
    "    - $assigned_paths" \
    'non_goals:' \
    '  - Do not change canonical plan or task status' \
    '  - Do not modify another session or worktree' \
    'context_refs:' \
    '  - AGENTS.md' \
    '  - WORKFLOW.md' \
    '  - context/ARCHITECTURE.md' \
    '  - docs/runtime-contract.md' \
    'repository: context-circuit' \
    "worktree: $assigned_worktree" \
    'permissions:' \
    "  write_worktree: $write_worktree" \
    '  write_runtime_session: true' \
    '  write_plan: false' \
    '  write_activity: false' \
    'acceptance_criteria:' \
    '  - The assigned behavior is implemented and verified' \
    'stop_conditions:' \
    '  - A required change falls outside the assigned paths' \
    '  - A source contradiction affects the result' \
    '  - A required dependency or verification is missing' \
    'handoff_schema: session-handoff-v1'
}

write_handoff() {
  session_id=$1
  parent_id=$2
  root_id=$3
  handoff_status=$4
  assigned_task=$5
  outcome=$6
  session_dir="$runtime/sessions/$session_id"
  atomic_write "$session_dir/handoff.md" \
    '# Session handoff' \
    "- session_id: $session_id" \
    "- parent_session_id: $parent_id" \
    "- root_session_id: $root_id" \
    "- status: $handoff_status" \
    '- plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
    "- task: $assigned_task" \
    '- worktree: .runtime/worktrees/context-circuit/plan-alpha' \
    '' \
    '## Objective and scope' \
    '' \
    '- Work only within the delegated task and worktree.' \
    '' \
    '## Evidence inspected' \
    '' \
    '- Runtime records, source documents, and Git state.' \
    '' \
    '## Changed files' \
    '' \
    '- None outside the assigned worktree.' \
    '' \
    '## Tests and verification' \
    '' \
    '- The assigned verification command was recorded.' \
    '' \
    '## Decisions and assumptions' \
    '' \
    "- $outcome" \
    '' \
    '## Questions, blockers, and limitations' \
    '' \
    '- Escalate any missing evidence or ownership ambiguity.' \
    '' \
    '## Recommended next action' \
    '' \
    '- Resume within scope or request the required human decision.'
}

session_complete() {
  session_file=$1
  for required in \
    'schema_version: 1' 'session_id:' 'parent_session_id:' \
    'root_session_id:' 'kind:' 'role:' 'objective:' 'plan:' 'task:' \
    'scope:' 'write_access:' 'worktree:' 'status:'; do
    grep -F "$required" "$session_file" >/dev/null 2>&1 || return 1
  done
  return 0
}

packet_complete() {
  packet_file=$1
  for required in \
    'schema_version: 1' 'session_id:' 'parent_session_id:' \
    'root_session_id:' 'role:' 'objective:' 'scope:' 'non_goals:' \
    'context_refs:' 'repository:' 'worktree:' 'permissions:' \
    'acceptance_criteria:' 'stop_conditions:' 'handoff_schema:'; do
    grep -F "$required" "$packet_file" >/dev/null 2>&1 || return 1
  done
  return 0
}

packet_matches_session() {
  session_file=$1
  packet_file=$2
  session_id=$(sed -n 's/^session_id: //p' "$session_file")
  parent_id=$(sed -n 's/^parent_session_id: //p' "$session_file")
  root_id=$(sed -n 's/^root_session_id: //p' "$session_file")
  role=$(sed -n 's/^role: //p' "$session_file")
  objective=$(sed -n 's/^objective: //p' "$session_file")
  plan=$(sed -n 's/^plan: //p' "$session_file")
  task=$(sed -n 's/^task: //p' "$session_file")
  scope=$(sed -n 's/^scope: //p' "$session_file")
  worktree=$(sed -n 's/^worktree: //p' "$session_file")
  grep -F "session_id: $session_id" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "parent_session_id: $parent_id" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "root_session_id: $root_id" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "role: $role" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "objective: $objective" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "  plan: $plan" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "  task: $task" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "$scope" "$packet_file" >/dev/null 2>&1 || return 1
  grep -F "worktree: $worktree" "$packet_file" >/dev/null 2>&1 || return 1
  for context_ref in AGENTS.md WORKFLOW.md docs/runtime-contract.md; do
    grep -F "  - $context_ref" "$packet_file" >/dev/null 2>&1 || return 1
  done
  if [ "$role" = verifier ]; then
    grep -F 'write_worktree: false' "$packet_file" >/dev/null 2>&1 || return 1
    grep -F 'write_plan: false' "$packet_file" >/dev/null 2>&1 || return 1
    grep -F 'write_activity: false' "$packet_file" >/dev/null 2>&1 || return 1
  fi
  return 0
}

route_root() {
  route_runtime=$1
  session_file=''
  any_session_file=''
  for candidate in $(find "$route_runtime/sessions" -mindepth 2 -maxdepth 2 -name session.yaml -print 2>/dev/null); do
    any_session_file=$candidate
    if grep -F 'kind: root' "$candidate" >/dev/null 2>&1; then
      session_file=$candidate
      break
    fi
  done
  if [ -z "$any_session_file" ]; then
    printf 'orienting\n'
    return 0
  fi
  if [ -z "$session_file" ]; then
    printf 'blocked\n'
    return 0
  fi
  if ! session_complete "$session_file"; then
    printf 'blocked\n'
    return 0
  fi
  if grep -R -F 'status: blocked' "$route_runtime/sessions" >/dev/null 2>&1; then
    printf 'blocked\n'
  elif grep -R -F 'status: handoff' "$route_runtime/sessions" >/dev/null 2>&1; then
    printf 'resume\n'
  else
    printf 'execute\n'
  fi
}

acquire_lease() {
  lock_dir=$1
  lease_session=$2
  lease_root=$3
  lease_plan=$4
  lease_worktree=$5
  lease_branch=$6
  plan_dir=$(dirname "$lock_dir")
  if ! mkdir "$lock_dir" 2>/dev/null; then
    return 1
  fi
  atomic_write "$lock_dir/owner.yaml" \
    'schema_version: 1' \
    "plan: $lease_plan" \
    "session_id: $lease_session" \
    "root_session_id: $lease_root" \
    'repository: context-circuit' \
    "worktree: $lease_worktree" \
    "branch: $lease_branch"
  atomic_write "$plan_dir/lease.yaml" \
    'schema_version: 1' \
    "plan: $lease_plan" \
    "session_id: $lease_session" \
    "root_session_id: $lease_root" \
    "worktree: $lease_worktree" \
    "branch: $lease_branch" \
    'status: active' \
    'acquired_at: 2026-08-14T18:00:00Z' \
    'heartbeat_at: 2026-08-14T18:00:00Z' \
    'released_at: null' \
    'stale_after_seconds: 1800'
  return 0
}

valid_lease() {
  lease_lock=$1
  expected_session=$2
  expected_plan=$3
  expected_worktree=$4
  owner_file="$lease_lock/owner.yaml"
  lease_file="$(dirname "$lease_lock")/lease.yaml"
  test -f "$owner_file" || return 1
  test -f "$lease_file" || return 1
  grep -F "session_id: $expected_session" "$owner_file" >/dev/null 2>&1 || return 1
  grep -F "plan: $expected_plan" "$owner_file" >/dev/null 2>&1 || return 1
  grep -F "worktree: $expected_worktree" "$owner_file" >/dev/null 2>&1 || return 1
  grep -F "session_id: $expected_session" "$lease_file" >/dev/null 2>&1 || return 1
  grep -F 'status: active' "$lease_file" >/dev/null 2>&1 || return 1
  return 0
}

completion_ready() {
  completion_plan=$1
  evidence_dir=$2
  verifier_file=$3
  gate_file=$4
  grep -F 'status: approved' "$completion_plan" >/dev/null 2>&1 || return 1
  for task_id in AWF-0001 AWF-0002 AWF-0003 AWF-0004 AWF-0005 AWF-0006 AWF-0007; do
    evidence_file="$evidence_dir/$task_id.yaml"
    grep -F 'evidence_status: completed' "$evidence_file" >/dev/null 2>&1 || return 1
  done
  grep -F 'status: completed' "$verifier_file" >/dev/null 2>&1 || return 1
  grep -F 'verification: passed' "$verifier_file" >/dev/null 2>&1 || return 1
  grep -F 'status: approved' "$gate_file" >/dev/null 2>&1 || return 1
  return 0
}

lifecycle_expected_task_status() {
  case "$1" in
    draft) printf 'draft\n' ;;
    approved) printf 'ready\n' ;;
    done) printf 'done\n' ;;
    *) return 1 ;;
  esac
}

lifecycle_sync_tasks() {
  lifecycle_plan_file=$1
  lifecycle_task_dir=$2
  lifecycle_plan_status=$(sed -n 's/^status: //p' "$lifecycle_plan_file" | head -n 1)
  lifecycle_expected_status=$(lifecycle_expected_task_status "$lifecycle_plan_status") || return 1
  for lifecycle_task_file in "$lifecycle_task_dir"/*.md; do
    test -f "$lifecycle_task_file" || continue
    lifecycle_current_status=$(sed -n 's/^status: //p' "$lifecycle_task_file" | head -n 1)
    if test "$lifecycle_current_status" != "$lifecycle_expected_status"; then
      lifecycle_temporary="${lifecycle_task_file}.tmp.$$"
      sed "s/^status: .*/status: $lifecycle_expected_status/" \
        "$lifecycle_task_file" > "$lifecycle_temporary"
      mv "$lifecycle_temporary" "$lifecycle_task_file"
    fi
  done
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
contains agents/reviewer.md 'Remain read-only with respect'

# The pure workflow has no migration shim or command-specific runtime.
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
  if grep -E 'cc\.mjs|npm run|node --import|legacy-run-task|compatibility adapter' "$file" >/dev/null 2>&1; then
    fail "legacy command-layer reference remains in normative file: $file"
  fi
done

fixture=$(mktemp -d "${TMPDIR:-/tmp}/context-circuit-acceptance.XXXXXX")
trap 'rm -rf "$fixture"' EXIT HUP INT TERM

runtime="$fixture/.runtime"
mkdir -p "$runtime/sessions" "$runtime/plans" "$runtime/worktrees/context-circuit"

# Fresh entry, PRD orientation, resume, malformed state, and blocked recovery.
fresh_runtime="$fixture/fresh/.runtime"
mkdir -p "$fresh_runtime/sessions" "$fixture/fresh/context"
test "$(route_root "$fresh_runtime")" = orienting
atomic_write "$fixture/fresh/prd.md" 'A source document for a new workspace.'
atomic_write "$fixture/fresh/context/PROJECT.md" \
  '# Draft Product Knowledge' 'source: ../prd.md' 'status: draft'
atomic_write "$fixture/fresh/plan.yaml" 'status: draft' 'source: prd.md'
test "$(route_root "$fresh_runtime")" = orienting
contains "$fixture/fresh/context/PROJECT.md" 'source: ../prd.md'
contains "$fixture/fresh/context/PROJECT.md" 'status: draft'
contains "$fixture/fresh/plan.yaml" 'status: draft'
if grep -F 'status: approved' "$fixture/fresh/plan.yaml" >/dev/null 2>&1; then
  fail 'fresh PRD orientation silently approved the proposed plan'
fi
test ! -e "$fresh_runtime/current-session.yaml"
test ! -e "$fresh_runtime/current-plan.yaml"

resume_runtime="$fixture/resume/.runtime"
mkdir -p "$resume_runtime/sessions/root-resume"
runtime="$resume_runtime"
write_session root-resume null root-resume root coordinator 'Resume prior work' AWF-0001..AWF-0007 workspace false . handoff
write_handoff root-resume null root-resume awaiting-human-gate AWF-0001 'Resume from the durable handoff.'
test "$(route_root "$resume_runtime")" = resume
contains "$resume_runtime/sessions/root-resume/handoff.md" 'Resume from the durable handoff.'

blocked_runtime="$fixture/blocked/.runtime"
mkdir -p "$blocked_runtime/sessions/root-blocked"
runtime="$blocked_runtime"
write_session root-blocked null root-blocked root coordinator 'Recover blocked work' AWF-0001..AWF-0007 workspace false . blocked
write_handoff root-blocked null root-blocked blocked AWF-0001 'Request explicit recovery or human takeover.'
test "$(route_root "$blocked_runtime")" = blocked
contains "$blocked_runtime/sessions/root-blocked/handoff.md" 'Request explicit recovery'

runtime="$fixture/.runtime"
mkdir -p "$runtime/sessions/root-001" "$runtime/sessions/child-001" "$runtime/sessions/research-001" "$runtime/sessions/verifier-001"
write_session root-001 null root-001 root coordinator 'Coordinate the workspace request' AWF-0001..AWF-0007 workspace false . executing
write_session child-001 root-001 root-001 subagent implementer 'Implement one bounded task' AWF-0001 'docs/ and test/' true .runtime/worktrees/context-circuit/plan-alpha executing
write_session research-001 root-001 root-001 subagent researcher 'Gather source evidence' AWF-0001 'context/ and source/' false .runtime/worktrees/context-circuit/plan-alpha verifying
write_session verifier-001 root-001 root-001 subagent verifier 'Verify the assigned worktree' AWF-0006 'docs/ and test/' false .runtime/worktrees/context-circuit/plan-alpha verifying
write_delegation child-001 root-001 root-001 implementer 'Implement one bounded task' AWF-0001 'docs/ and test/' .runtime/worktrees/context-circuit/plan-alpha true
write_delegation research-001 root-001 root-001 researcher 'Gather source evidence' AWF-0001 'context/ and source/' .runtime/worktrees/context-circuit/plan-alpha false
write_delegation verifier-001 root-001 root-001 verifier 'Verify the assigned worktree' AWF-0006 'docs/ and test/' .runtime/worktrees/context-circuit/plan-alpha false
write_handoff child-001 root-001 root-001 blocked AWF-0001 'The worker stopped at a scope boundary.'
write_handoff research-001 root-001 root-001 completed AWF-0001 'Research evidence was returned to the parent.'
write_handoff verifier-001 root-001 root-001 failed AWF-0006 'Verification failed and requires a human decision.'

expect_success session_complete "$runtime/sessions/root-001/session.yaml"
expect_success session_complete "$runtime/sessions/child-001/session.yaml"
expect_success packet_complete "$runtime/sessions/child-001/delegation.yaml"
expect_success packet_complete "$runtime/sessions/research-001/delegation.yaml"
expect_success packet_complete "$runtime/sessions/verifier-001/delegation.yaml"
expect_success packet_matches_session "$runtime/sessions/child-001/session.yaml" "$runtime/sessions/child-001/delegation.yaml"
expect_success packet_matches_session "$runtime/sessions/research-001/session.yaml" "$runtime/sessions/research-001/delegation.yaml"
expect_success packet_matches_session "$runtime/sessions/verifier-001/session.yaml" "$runtime/sessions/verifier-001/delegation.yaml"
contains "$runtime/sessions/child-001/delegation.yaml" 'write_plan: false'
contains "$runtime/sessions/research-001/handoff.md" 'status: completed'
contains "$runtime/sessions/verifier-001/delegation.yaml" 'write_worktree: false'
contains "$runtime/sessions/child-001/handoff.md" 'status: blocked'
contains "$runtime/sessions/verifier-001/handoff.md" 'status: failed'
expect_success handoff_status_valid completed
expect_success handoff_status_valid blocked
expect_success handoff_status_valid failed
expect_success handoff_status_valid awaiting-human-gate
expect_failure handoff_status_valid executing

mkdir -p "$runtime/sessions/incomplete"
atomic_write "$runtime/sessions/incomplete/delegation.yaml" \
  'schema_version: 1' 'session_id: incomplete' 'role: implementer' 'objective: incomplete'
expect_failure packet_complete "$runtime/sessions/incomplete/delegation.yaml"

# Safe identifiers reject traversal, absolute paths, and unsafe names.
expect_success safe_id sess-001
expect_success safe_id plan.alpha-001
expect_failure safe_id ../escape
expect_failure safe_id /absolute
expect_failure safe_id 'Unsafe ID'

# Atomic same-plan contention: exactly one contender owns the lock and the
# losing contender cannot alter the plan sentinel or worker session.
same_plan="$runtime/plans/plan-alpha"
mkdir -p "$same_plan"
atomic_write "$same_plan/plan-sentinel" 'must remain unchanged'
mkdir -p "$runtime/sessions/worker-alpha" "$runtime/plans/agent-workspace-workflow"
atomic_write "$runtime/sessions/worker-alpha/session.yaml" 'session_id: worker-alpha' 'status: executing'
sentinel_before=$(sha256sum "$same_plan/plan-sentinel")
worker_before=$(sha256sum "$runtime/sessions/worker-alpha/session.yaml")
set +e
(acquire_lease "$same_plan/lease.lock" worker-alpha root-001 plans/context-circuit-plans/0001-agent-workspace-workflow .runtime/worktrees/context-circuit/plan-alpha plan-alpha) & contender_a=$!
(acquire_lease "$same_plan/lease.lock" worker-beta root-002 plans/context-circuit-plans/0001-agent-workspace-workflow .runtime/worktrees/context-circuit/plan-alpha plan-alpha) & contender_b=$!
wait "$contender_a"
result_a=$?
wait "$contender_b"
result_b=$?
set -e
test "$result_a" -ne "$result_b" || fail 'same-plan contention produced zero or two owners'
if valid_lease "$same_plan/lease.lock" worker-alpha plans/context-circuit-plans/0001-agent-workspace-workflow .runtime/worktrees/context-circuit/plan-alpha; then
  :
else
  expect_success valid_lease "$same_plan/lease.lock" worker-beta plans/context-circuit-plans/0001-agent-workspace-workflow .runtime/worktrees/context-circuit/plan-alpha
fi
test "$(sha256sum "$same_plan/plan-sentinel")" = "$sentinel_before"
test "$(sha256sum "$runtime/sessions/worker-alpha/session.yaml")" = "$worker_before"

# Independent plans obtain independent runtime leases and worktree identities.
beta_plan="$runtime/plans/plan-beta"
mkdir -p "$beta_plan"
expect_success acquire_lease "$beta_plan/lease.lock" worker-gamma root-003 plans/context-circuit-plans/0002-independent-plan .runtime/worktrees/context-circuit/plan-beta plan-beta
expect_success valid_lease "$beta_plan/lease.lock" worker-gamma plans/context-circuit-plans/0002-independent-plan .runtime/worktrees/context-circuit/plan-beta
test "$same_plan/lease.lock" != "$beta_plan/lease.lock"
test ! -e "$runtime/current-session.yaml"
test ! -e "$runtime/current-plan.yaml"

# Missing owner evidence blocks a lease; it cannot be reconstructed.
malformed_plan="$runtime/plans/malformed"
mkdir -p "$malformed_plan/lease.lock"
atomic_write "$malformed_plan/lease.yaml" \
  'schema_version: 1' 'session_id: worker-zeta' 'status: active'
expect_failure valid_lease "$malformed_plan/lease.lock" worker-zeta plans/example .runtime/worktrees/context-circuit/malformed

# Real Git worktrees and branches prove isolation from the base checkout and
# from one another. All writes remain outside the implementation repository.
git_fixture="$fixture/git-fixture"
alpha_worktree="$fixture/worktrees/plan-alpha"
beta_worktree="$fixture/worktrees/plan-beta"
mkdir -p "$fixture/worktrees"
git init -q "$git_fixture"
git -C "$git_fixture" config user.name 'Acceptance Fixture'
git -C "$git_fixture" config user.email 'acceptance@example.invalid'
atomic_write "$git_fixture/README.md" 'base checkout'
git -C "$git_fixture" add README.md
git -C "$git_fixture" commit -qm 'initial fixture commit'
git -C "$git_fixture" worktree add -q -b plan-alpha "$alpha_worktree" HEAD
git -C "$git_fixture" worktree add -q -b plan-beta "$beta_worktree" HEAD
test "$(git -C "$alpha_worktree" branch --show-current)" = plan-alpha
test "$(git -C "$beta_worktree" branch --show-current)" = plan-beta
test "$(git -C "$alpha_worktree" rev-parse --show-toplevel)" = "$alpha_worktree"
test "$(git -C "$beta_worktree" rev-parse --show-toplevel)" = "$beta_worktree"
atomic_write "$alpha_worktree/alpha.txt" 'alpha-only change'
atomic_write "$beta_worktree/beta.txt" 'beta-only change'
test ! -e "$git_fixture/alpha.txt"
test ! -e "$git_fixture/beta.txt"
test ! -e "$alpha_worktree/beta.txt"
test ! -e "$beta_worktree/alpha.txt"
test -z "$(git -C "$git_fixture" status --porcelain)"
expect_success scope_allows "$alpha_worktree/alpha.txt" "$alpha_worktree"
expect_failure scope_allows "$git_fixture/README.md" "$alpha_worktree"
expect_failure scope_allows "$beta_worktree/beta.txt" "$alpha_worktree"

# Interruption, explicit takeover, and contradictory source handling preserve
# the old evidence and record why a new session is allowed to continue.
interrupted_dir="$runtime/sessions/interrupted-001"
mkdir -p "$interrupted_dir"
atomic_write "$interrupted_dir/session.yaml" \
  'session_id: interrupted-001' 'status: blocked' 'worktree: .runtime/worktrees/context-circuit/plan-alpha'
atomic_write "$interrupted_dir/handoff.md" \
  '# Session handoff' '- session_id: interrupted-001' '- status: blocked' \
  '- next_action: Inspect heartbeat and request recovery'
takeover_dir="$runtime/sessions/takeover-001"
mkdir -p "$takeover_dir"
atomic_write "$takeover_dir/session.yaml" \
  'session_id: takeover-001' 'parent_session_id: root-001' 'root_session_id: root-001' \
  'status: executing' 'replaces_session_id: interrupted-001' \
  'takeover_reason: human-authorized recovery'
contains "$takeover_dir/session.yaml" 'replaces_session_id: interrupted-001'
contains "$interrupted_dir/handoff.md" 'Inspect heartbeat and request recovery'

atomic_write "$fixture/source-a.md" 'source version A'
atomic_write "$fixture/source-b.md" 'source version B'
if cmp -s "$fixture/source-a.md" "$fixture/source-b.md"; then
  fail 'contradictory source fixture unexpectedly matched'
fi
atomic_write "$runtime/sessions/root-001/contradiction-handoff.md" \
  '# Contradictory source handoff' '- status: blocked' \
  '- blocker: source documents contradict one another' \
  '- next_action: Request human-reviewed context refresh'
contains "$runtime/sessions/root-001/contradiction-handoff.md" 'source documents contradict'

# A verifier may write only its own handoff. Worker, lease, plan, and activity
# evidence must remain byte-for-byte unchanged.
verifier_plan_snapshot="$fixture/verifier-plan.snapshot"
sha256sum "$same_plan/lease.yaml" "$same_plan/lease.lock/owner.yaml" \
  "$same_plan/plan-sentinel" "$runtime/sessions/worker-alpha/session.yaml" > "$verifier_plan_snapshot"
expect_success verifier_write "$runtime/sessions/verifier-001/allowed-evidence.md" "$runtime/sessions/verifier-001" \
  'evidence: verifier-owned'
assert_failure_reason "$fixture/verifier-lease-failure.log" 'verifier write is outside own session handoff' verifier_write "$same_plan/lease.yaml" "$runtime/sessions/verifier-001" 'tamper: lease'
assert_failure_reason "$fixture/verifier-worker-failure.log" 'verifier write is outside own session handoff' verifier_write "$runtime/sessions/worker-alpha/session.yaml" "$runtime/sessions/verifier-001" 'tamper: worker'
assert_failure_reason "$fixture/verifier-worktree-failure.log" 'verifier write is outside own session handoff' verifier_write "$alpha_worktree/forbidden.txt" "$runtime/sessions/verifier-001" 'tamper: worktree'
assert_failure_reason "$fixture/verifier-activity-failure.log" 'verifier write is outside own session handoff' verifier_write "$fixture/activity.log" "$runtime/sessions/verifier-001" 'tamper: activity'
assert_failure_reason "$fixture/lease-failure.log" 'same-plan lease is already owned' lease_guard "$same_plan/lease.lock"
assert_failure_reason "$fixture/scope-failure.log" 'outside assigned worktree' scope_guard "$git_fixture/README.md" "$alpha_worktree"
assert_failure_reason "$fixture/verification-failure.log" 'verification failed; human decision required' verification_guard "$runtime/sessions/verifier-001/handoff.md"
expect_success verifier_write "$runtime/sessions/verifier-001/handoff.md" "$runtime/sessions/verifier-001" \
  '# Verification handoff' '- status: failed' \
  '- verification: failed' '- next_action: Repair within scope or request a human decision'
verifier_plan_after=$(sha256sum "$same_plan/lease.yaml" "$same_plan/lease.lock/owner.yaml" \
  "$same_plan/plan-sentinel" "$runtime/sessions/worker-alpha/session.yaml")
test "$verifier_plan_after" = "$(sed -n '1,4p' "$verifier_plan_snapshot")"
contains "$runtime/sessions/verifier-001/allowed-evidence.md" 'evidence: verifier-owned'
contains "$runtime/sessions/verifier-001/handoff.md" 'status: failed'
contains "$runtime/sessions/verifier-001/handoff.md" 'human decision'

# Completion is only ready for a human status change after every task has
# evidence, verification passes, and the explicit human gate is present.
completion_plan="$fixture/completion-plan.yaml"
completion_tasks="$fixture/completion-evidence"
completion_verifier="$fixture/completion-verifier.md"
completion_gate="$fixture/completion-human-gate.yaml"
mkdir -p "$completion_tasks"
atomic_write "$completion_plan" 'status: approved' 'id: agent-workspace-workflow'
for task_id in AWF-0001 AWF-0002 AWF-0003 AWF-0004 AWF-0005 AWF-0006 AWF-0007; do
  atomic_write "$completion_tasks/$task_id.yaml" "task: $task_id" 'evidence_status: completed'
done
atomic_write "$completion_verifier" '# Verification handoff' 'status: failed' 'verification: failed'
atomic_write "$completion_gate" 'status: pending'
expect_failure completion_ready "$completion_plan" "$completion_tasks" "$completion_verifier" "$completion_gate"
atomic_write "$completion_verifier" '# Verification handoff' 'status: completed' 'verification: passed'
expect_failure completion_ready "$completion_plan" "$completion_tasks" "$completion_verifier" "$completion_gate"
atomic_write "$completion_gate" 'status: approved' 'gate: status-change'
expect_success completion_ready "$completion_plan" "$completion_tasks" "$completion_verifier" "$completion_gate"
atomic_write "$runtime/plans/agent-workspace-workflow/completion.yaml" \
  'schema_version: 1' 'plan: plans/context-circuit-plans/0001-agent-workspace-workflow' \
  'status: ready-for-human-status-change' 'human_gate: status-change' \
  'canonical_status_changed: false'
contains "$runtime/plans/agent-workspace-workflow/completion.yaml" 'canonical_status_changed: false'
test ! -e "$fixture/status-done-marker"

contains docs/agent-workspace-workflow.md 'Verification fails;'
contains docs/agent-workspace-workflow.md 'Human approval is required'
contains docs/agent-workspace-workflow.md 'source changes during execution'
contains docs/runtime-contract.md 'The canonical plan/task status remains unchanged'
contains agents/reviewer.md 'Do not repair, change'

# Plan 0006: plan status is canonical and task status is a bulk, idempotent
# projection. Resume repair must preserve provider-owned annotations and must
# not rerun implementation or verification checks.
contains WORKFLOW.md 'Plan status is the canonical lifecycle authority'
contains docs/agent-workspace-workflow.md 'does not require separate task selection'
contains docs/planning.md 'reconciles every included task'
contains docs/runtime-contract.md 'Plan/task lifecycle projection'
contains docs/run-task.md 'Task status is not a second approval or execution gate'
contains docs/configuration.md 'external_status'
contains context/CONVENTIONS.md 'synchronized projection'
if grep -F '    - task-selection' workspace.yaml >/dev/null 2>&1; then
  fail 'ordinary plan execution still lists task-selection as a human gate'
fi

lifecycle_fixture="$fixture/plan-lifecycle"
lifecycle_plan="$lifecycle_fixture/plan.yaml"
lifecycle_tasks="$lifecycle_fixture/tasks"
mkdir -p "$lifecycle_tasks"
atomic_write "$lifecycle_plan" 'id: lifecycle-fixture' 'status: draft'
for lifecycle_id in CCL-0001 CCL-0002 CCL-0003; do
  atomic_write "$lifecycle_tasks/$lifecycle_id.md" \
    "id: $lifecycle_id" 'status: draft' 'title: Fixture task'
done
lifecycle_draft_snapshot="$fixture/lifecycle-draft.snapshot"
sha256sum "$lifecycle_tasks"/*.md > "$lifecycle_draft_snapshot"
expect_success lifecycle_sync_tasks "$lifecycle_plan" "$lifecycle_tasks"
test "$(sha256sum "$lifecycle_tasks"/*.md)" = "$(cat "$lifecycle_draft_snapshot")"

atomic_write "$lifecycle_plan" 'id: lifecycle-fixture' 'status: approved'
expect_success lifecycle_sync_tasks "$lifecycle_plan" "$lifecycle_tasks"
for lifecycle_id in CCL-0001 CCL-0002 CCL-0003; do
  contains "$lifecycle_tasks/$lifecycle_id.md" 'status: ready'
done
lifecycle_approved_snapshot="$fixture/lifecycle-approved.snapshot"
sha256sum "$lifecycle_tasks"/*.md > "$lifecycle_approved_snapshot"
expect_success lifecycle_sync_tasks "$lifecycle_plan" "$lifecycle_tasks"
test "$(sha256sum "$lifecycle_tasks"/*.md)" = "$(cat "$lifecycle_approved_snapshot")"

atomic_write "$lifecycle_tasks/CCL-0001.md" \
  'id: CCL-0001' 'status: draft' 'title: Fixture task' \
  'external_status:' '  provider: optional-fixture' '  value: in_progress'
atomic_write "$lifecycle_tasks/CCL-0002.md" \
  'id: CCL-0002' 'status: ready' 'title: Fixture task'
atomic_write "$lifecycle_tasks/CCL-0003.md" \
  'id: CCL-0003' 'status: ready' 'title: Fixture task'
atomic_write "$lifecycle_fixture/verification-runs" '0'
expect_success lifecycle_sync_tasks "$lifecycle_plan" "$lifecycle_tasks"
for lifecycle_id in CCL-0001 CCL-0002 CCL-0003; do
  contains "$lifecycle_tasks/$lifecycle_id.md" 'status: ready'
done
contains "$lifecycle_tasks/CCL-0001.md" 'external_status:'
contains "$lifecycle_tasks/CCL-0001.md" 'value: in_progress'
test "$(cat "$lifecycle_fixture/verification-runs")" = 0

atomic_write "$lifecycle_plan" 'id: lifecycle-fixture' 'status: done'
expect_success lifecycle_sync_tasks "$lifecycle_plan" "$lifecycle_tasks"
for lifecycle_id in CCL-0001 CCL-0002 CCL-0003; do
  contains "$lifecycle_tasks/$lifecycle_id.md" 'status: done'
done
contains "$lifecycle_tasks/CCL-0001.md" 'external_status:'
test "$(cat "$lifecycle_fixture/verification-runs")" = 0

# Plan 0002: workspace initialization is identity-only and supports a valid
# zero-repository state plus repository roles and default active branches.
require_file context/PRODUCT-DIRECTION.md
require_file sources/README.md
require_file .agents/skills/cc-initialize-workspace/SKILL.md
require_file .agents/skills/cc-initialize-workspace/agents/openai.yaml
require_file .agents/skills/cc-idea-brief/SKILL.md
require_file .agents/skills/cc-create-prd/SKILL.md
require_file docs/idea-brief.md
require_file docs/prd.md

contains .agents/skills/cc-initialize-workspace/SKILL.md 'Zero repositories is valid.'
contains .agents/skills/cc-initialize-workspace/SKILL.md 'Recommend `development` only when that branch exists'
contains .agents/skills/cc-initialize-workspace/SKILL.md 'does not ask about delivery behavior'
contains docs/getting-started.md 'no-repository path as an error'
contains workspace.yaml 'repositories: {}'
if grep -E '^[[:space:]]*[0-9]+\..*(delivery|commit|push|merge|publication|deployment|external)' \
  .agents/skills/cc-initialize-workspace/SKILL.md >/dev/null 2>&1; then
  fail 'initialization turned optional delivery behavior into a core question'
fi

foundation_fixture="$fixture/workspace-foundation"
mkdir -p "$foundation_fixture/zero/sources" \
  "$foundation_fixture/existing/sources" "$foundation_fixture/existing/repositories"
atomic_write "$foundation_fixture/zero/workspace.yaml" \
  'workspace:' '  name: new-idea' '  mode: solo' 'repositories: {}'
contains "$foundation_fixture/zero/workspace.yaml" 'repositories: {}'
atomic_write "$foundation_fixture/existing/workspace.yaml" \
  'workspace:' '  name: existing-project' '  mode: team' \
  'repositories:' '  api:' '    path: ../api' '    mode: ignored-clone' \
  '    role: Public API' '    agent: repository-worker' \
  '    default_branch: development' \
  '  docs:' '    path: ../docs' '    mode: ignored-clone' \
  '    role: Product documentation' '    agent: repository-worker' \
  '    default_branch: main'
for registry_field in 'path:' 'mode:' 'role:' 'agent:' 'default_branch:'; do
  contains "$foundation_fixture/existing/workspace.yaml" "$registry_field"
done

recommend_default_active_branch() {
  branch_list=$1
  case " $branch_list " in
    *' development '*) printf 'development\n' ;;
    *) printf '%s\n' "${branch_list%% *}" ;;
  esac
}
test "$(recommend_default_active_branch 'main development feature')" = development
test "$(recommend_default_active_branch 'main feature')" = main

# Source intake is request-scoped. An unrelated entry reads context only, while
# a selected source request leaves an evidence trail and does not copy raw text.
source_fixture="$foundation_fixture/source-boundary"
mkdir -p "$source_fixture/sources" "$source_fixture/context"
atomic_write "$source_fixture/context/INDEX.md" '# Context index'
atomic_write "$source_fixture/sources/selected.md" \
  'selected source evidence' 'do not copy this raw sentence'
atomic_write "$source_fixture/sources/unselected.md" \
  'unselected source evidence' 'must not be read for this request'
source_read_log="$source_fixture/read.log"
normal_session_entry() {
  cat "$source_fixture/context/INDEX.md" >/dev/null
  test ! -e "$source_read_log"
}
read_selected_source() {
  selected_path=$1
  cat "$selected_path" >/dev/null
  printf '%s\n' "${selected_path#"$source_fixture/"}" >> "$source_read_log"
}
expect_success normal_session_entry
expect_success read_selected_source "$source_fixture/sources/selected.md"
contains "$source_read_log" 'sources/selected.md'
if grep -F 'sources/unselected.md' "$source_read_log" >/dev/null 2>&1; then
  fail 'source-based fixture read an unselected source'
fi
atomic_write "$source_fixture/context/sources.yaml" \
  'sources:' '  - id: selected' \
  '    location: sources/selected.md' \
  '    read_for: Create the requested Idea Brief' \
  '    used_by: contributions/idea-briefs/example.md'
contains "$source_fixture/context/sources.yaml" 'read_for: Create the requested Idea Brief'
test ! -e "$source_fixture/context/RAW-SOURCE.md"

# The navigation index and artifact fixtures keep raw sources, accepted context,
# product artifacts, plans, and private runtime state in distinct homes.
for layer_link in '`sources/`' '`context/`' '`contributions/`' '`plans/`' '`.runtime/`'; do
  contains context/INDEX.md "$layer_link"
done
idea_fixture="$foundation_fixture/contributions/idea-briefs/example.md"
prd_fixture="$foundation_fixture/contributions/prds/example.md"
mkdir -p "$(dirname "$idea_fixture")" "$(dirname "$prd_fixture")"
atomic_write "$idea_fixture" \
  'kind: idea-brief' 'status: draft' '# Example idea' \
  '## Intent' '## Assumptions' '## Open questions' '## Provenance' \
  'selected source paths: sources/selected.md'
atomic_write "$prd_fixture" \
  'kind: prd' 'status: draft' '# Example PRD' \
  '## Requirements' '## Non-goals' '## Acceptance criteria' '## Provenance' \
  'selected source paths: sources/selected.md'
contains "$idea_fixture" 'status: draft'
contains "$prd_fixture" 'status: draft'
contains .agents/skills/cc-idea-brief/SKILL.md 'Human gate:'
contains .agents/skills/cc-idea-brief/SKILL.md 'may stop after the Idea Brief'
contains .agents/skills/cc-create-prd/SKILL.md 'human acceptance'
contains .agents/skills/cc-create-prd/SKILL.md 'silently approve a plan'
contains .agents/skills/cc-idea-brief/agents/openai.yaml 'cc-idea-brief'
contains .agents/skills/cc-create-prd/agents/openai.yaml 'cc-create-prd'
contains .agents/skills/cc-initialize-workspace/agents/openai.yaml 'cc-initialize-workspace'

# Plan 0005: Domain Knowledge is canonical for bounded areas and Role Knowledge
# provides a linked cross-domain perspective with request-scoped generation.
require_file docs/templates/domain-context.md
require_file docs/templates/role-context.md
require_file context/domains/README.md
require_file context/roles/README.md
require_file .agents/skills/cc-gather-context/SKILL.md
require_file .agents/skills/cc-gather-context/agents/openai.yaml

for metadata_field in 'kind: domain' 'status: proposed' 'sources:' 'freshness:' \
  'assumptions:' 'unknowns:' 'contradictions:' 'acceptance:'; do
  contains docs/templates/domain-context.md "$metadata_field"
done
for metadata_field in 'kind: role' 'status: proposed' 'domains:' 'sources:' \
  'freshness:' 'assumptions:' 'unknowns:' 'contradictions:' 'acceptance:'; do
  contains docs/templates/role-context.md "$metadata_field"
done
if grep -E '^[[:space:]]*(relevant_domains|contributes_to|acts_on|consumes|approves):' \
  docs/templates/role-context.md >/dev/null 2>&1; then
  fail 'role template contains a relationship matrix instead of a simple domains list'
fi
contains docs/product-knowledge.md 'Domain Knowledge is canonical'
contains docs/product-knowledge.md 'Role Knowledge is a cross-domain view'
contains docs/product-knowledge.md 'Business/project roles are not agent execution roles.'
contains .agents/skills/cc-gather-context/SKILL.md 'Read only those selected files'
contains .agents/skills/cc-gather-context/SKILL.md 'Never silently overwrite accepted context'
contains .agents/skills/cc-gather-context/SKILL.md 'domains:'
contains context/INDEX.md 'smallest relevant'
contains context/INDEX.md 'domain and role'
contains context/domains/README.md 'Do not recursively read every domain'
contains context/roles/README.md 'workflow pages linked'
contains context/sources.yaml 'approved-product-direction-domain-role-context'
contains context/sources.yaml 'docs/templates/role-context.md'
contains .agents/skills/cc-gather-context/agents/openai.yaml 'cc-gather-context'

domain_role_fixture="$fixture/domain-role-context"
mkdir -p "$domain_role_fixture/sources" \
  "$domain_role_fixture/context/domains/payments" \
  "$domain_role_fixture/context/domains/unrelated" \
  "$domain_role_fixture/context/roles"
atomic_write "$domain_role_fixture/sources/selected.md" \
  'selected payment evidence' 'raw-only-secret-must-not-be-copied'
atomic_write "$domain_role_fixture/sources/unselected.md" \
  'unselected evidence must not be read'
domain_read_log="$domain_role_fixture/domain-read.log"
select_domain_sources() {
  selected_source=$1
  test "$selected_source" = "$domain_role_fixture/sources/selected.md"
  cat "$selected_source" >/dev/null
  printf '%s\n' "${selected_source#"$domain_role_fixture/"}" >> "$domain_read_log"
}
expect_success select_domain_sources "$domain_role_fixture/sources/selected.md"
contains "$domain_read_log" 'sources/selected.md'
if grep -F 'sources/unselected.md' "$domain_read_log" >/dev/null 2>&1; then
  fail 'domain generation read an unselected raw source'
fi

domain_page="$domain_role_fixture/context/domains/payments/README.md"
atomic_write "$domain_page" \
  'kind: domain' 'status: proposed' 'title: Payments' \
  'sources:' '  - sources/selected.md' \
  'source_revisions:' '  - sha256:fixture-selected' \
  'freshness: current' 'assumptions:' '  - The payment owner is not confirmed.' \
  'unknowns:' '  - Refund timing is not described.' 'contradictions: []' \
  'acceptance:' '  state: pending' \
  '# Payments' '## Summary' 'Selected payment evidence is summarized here.'
contains "$domain_page" 'kind: domain'
contains "$domain_page" 'status: proposed'
contains "$domain_page" 'sources/selected.md'
contains "$domain_page" 'assumptions:'
contains "$domain_page" 'unknowns:'
contains "$domain_page" 'acceptance:'
if grep -F 'raw-only-secret-must-not-be-copied' "$domain_page" >/dev/null 2>&1; then
  fail 'domain generation copied raw source text into Product Knowledge'
fi
test ! -e "$domain_role_fixture/context/payments.md"

accepted_snapshot="$domain_role_fixture/accepted-domain.snapshot"
atomic_write "$domain_page" \
  'kind: domain' 'status: accepted' 'title: Payments' \
  'sources:' '  - sources/selected.md' 'freshness: current' \
  'acceptance:' '  state: accepted' '  accepted_at: 2026-08-15' \
  '# Payments' '## Accepted behavior' 'The accepted payment behavior.'
sha256sum "$domain_page" > "$accepted_snapshot"
refresh_with_contradiction() {
  test "$(sed -n 's/^status: //p' "$domain_page" | head -n 1)" = accepted
  atomic_write "$domain_page.proposed" \
    'kind: domain' 'status: needs-review' 'refreshes: context/domains/payments/README.md' \
    'contradictions:' '  - New evidence conflicts with accepted payment behavior.' \
    'acceptance:' '  state: needs-review' \
    '# Payments refresh proposal'
}
expect_success refresh_with_contradiction
test "$(sha256sum "$domain_page")" = "$(cat "$accepted_snapshot")"
contains "$domain_page.proposed" 'status: needs-review'
contains "$domain_page.proposed" 'contradictions:'

role_page="$domain_role_fixture/context/roles/planner.md"
atomic_write "$role_page" \
  'kind: role' 'status: proposed' 'title: Product planner' \
  'domains:' '  - ../domains/payments/README.md' \
  '  - ../domains/identity/README.md' \
  'sources:' '  - sources/selected.md' 'freshness: current' \
  'assumptions: []' 'unknowns: []' 'contradictions: []' \
  'acceptance:' '  state: pending' \
  '# Product planner' '## Cross-domain perspective' \
  'Links the payment and identity domains without restating their facts.'
contains "$role_page" 'domains:'
contains "$role_page" '../domains/payments/README.md'
contains "$role_page" '../domains/identity/README.md'
if grep -E '^[[:space:]]*(relevant_domains|contributes_to|acts_on|consumes|approves):' \
  "$role_page" >/dev/null 2>&1; then
  fail 'generated role contains a relationship matrix field'
fi
if grep -F 'The accepted payment behavior.' "$role_page" >/dev/null 2>&1; then
  fail 'generated role duplicated canonical domain facts'
fi
test ! -e "$domain_role_fixture/context/roles/agent-worker.md"

role_selection_log="$domain_role_fixture/role-selection.log"
select_role_context() {
  cat "$role_page" >/dev/null
  cat "$domain_page" >/dev/null
  printf '%s\n' 'context/roles/planner.md' 'context/domains/payments/README.md' \
    >> "$role_selection_log"
}
expect_success select_role_context
contains "$role_selection_log" 'context/roles/planner.md'
contains "$role_selection_log" 'context/domains/payments/README.md'
if grep -F 'unrelated' "$role_selection_log" >/dev/null 2>&1; then
  fail 'role context selection read an unrelated domain'
fi

printf 'PASS: Domain and Role Knowledge acceptance scenarios (templates, selective generation, provenance, links, contradictions, acceptance)\n'

printf 'PASS: workspace foundation and context acceptance scenarios (initialization, passive sources, provenance, artifacts)\n'

printf 'PASS: pure agent-workspace acceptance scenarios (filesystem, contention, isolation, recovery, verification, gates)\n'
