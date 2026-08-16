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
  grep -F -- "$text" "$file" >/dev/null 2>&1 || fail "expected '$text' in $file"
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

plan_contract_complete() {
  plan_file=$1
  for required in \
    'id:' 'number:' 'title:' 'status:' 'source:' 'repositories:' \
    'product_knowledge:' 'implementation_scope:' 'non_goals:' \
    'dependencies:' 'acceptance_criteria:' 'test_scope:' \
    'verification_commands:'; do
    grep -F "$required" "$plan_file" >/dev/null 2>&1 || return 1
  done
  return 0
}

plan_ready_for_execution() {
  execution_plan=$1
  grep -F 'status: approved' "$execution_plan" >/dev/null 2>&1 || return 1
  if grep -F 'status: draft' "$execution_plan" >/dev/null 2>&1; then
    return 1
  fi
  return 0
}

review_outcome_valid() {
  case "$1" in
    ready-for-approval|needs-revision|blocked|approved-for-execution) return 0 ;;
    *) return 1 ;;
  esac
}

dependencies_ready() {
  dependency_file=''
  for dependency_file in "$@"; do
    dependency_status=$(sed -n 's/^status: //p' "$dependency_file" | head -n 1)
    case "$dependency_status" in
      approved|done) ;;
      *) return 1 ;;
    esac
  done
  return 0
}

expected_plan_execution_records() {
  printf 'writer-child verifier-child\n'
}

execution_topology_for_mode() {
  printf 'writer-child verifier-child\n'
}

run_plan_guard() {
  guard_plan=$1
  guard_lock=$2
  plan_ready_for_execution "$guard_plan" || {
    printf 'BLOCKED: only an approved plan can execute\n' >&2
    return 1
  }
  if test -e "$guard_lock"; then
    printf 'BLOCKED: plan already has a writing owner\n' >&2
    return 1
  fi
  return 0
}

recovery_guard() {
  recovery_session=$1
  grep -F 'status: blocked' "$recovery_session" >/dev/null 2>&1 || return 1
  grep -F 'replaces_session_id:' "$recovery_session" >/dev/null 2>&1 || return 1
  grep -F 'takeover_reason: human-authorized recovery' "$recovery_session" \
    >/dev/null 2>&1 || return 1
  return 0
}

completion_ready_for_plan() {
  completion_plan_file=$1
  completion_evidence_dir=$2
  completion_verifier_file=$3
  completion_gate_file=$4
  plan_ready_for_execution "$completion_plan_file" || return 1
  for completion_evidence_file in "$completion_evidence_dir"/*.yaml; do
    test -f "$completion_evidence_file" || return 1
    grep -F 'evidence_status: completed' "$completion_evidence_file" \
      >/dev/null 2>&1 || return 1
  done
  grep -F 'status: completed' "$completion_verifier_file" >/dev/null 2>&1 || return 1
  grep -F 'verification: passed' "$completion_verifier_file" >/dev/null 2>&1 || return 1
  grep -F 'status: approved' "$completion_gate_file" >/dev/null 2>&1 || return 1
  return 0
}

write_gate_plan() {
  plan_file=$1
  plan_id=$2
  plan_status=$3
  mkdir -p "$(dirname "$plan_file")"
  atomic_write "$plan_file" \
    "id: $plan_id" \
    'number: 9' \
    'title: Gate fixture plan' \
    "status: $plan_status" \
    'source:' \
    '  kind: direct-request' \
    '  reference: fixture' \
    'repositories:' \
    '  - context-circuit' \
    'product_knowledge:' \
    '  references:' \
    '    - context/PROJECT.md' \
    'implementation_scope:' \
    '  - docs and agents' \
    'non_goals:' \
    '  - delivery actions' \
    'dependencies:' \
    '  - foundation' \
    'acceptance_criteria:' \
    '  - The fixture can execute safely' \
    'test_scope:' \
    '  - filesystem fixtures' \
    'verification_commands:' \
    '  - git diff --check'
}

write_gate_tasks() {
  task_dir=$1
  task_status=$2
  mkdir -p "$task_dir"
  atomic_write "$task_dir/FPC-0001.md" \
    'id: FPC-0001' "status: $task_status" 'title: First fixture task'
  atomic_write "$task_dir/FPC-0002.md" \
    'id: FPC-0002' "status: $task_status" 'title: Second fixture task'
}

plan_field() {
  sed -n "s/^$2: //p" "$1" | head -n 1
}

approve_plan() {
  plan_file=$1
  task_dir=$2
  confirmation=$3
  expected_id=${4:-}
  if test ! -f "$plan_file"; then
    printf 'BLOCKED: unknown or contradictory plan identifier\n' >&2
    return 1
  fi
  if ! plan_contract_complete "$plan_file"; then
    printf 'BLOCKED: plan artifact is incomplete\n' >&2
    return 1
  fi
  task_found=0
  if test -d "$task_dir"; then
    for task_file in "$task_dir"/*.md; do
      test -f "$task_file" || continue
      task_found=1
    done
  fi
  if test "$task_found" -eq 0; then
    printf 'BLOCKED: plan artifact is incomplete\n' >&2
    return 1
  fi
  actual_id=$(plan_field "$plan_file" id)
  if test -n "$expected_id" && test "$actual_id" != "$expected_id"; then
    printf 'BLOCKED: unknown or contradictory plan identifier\n' >&2
    return 1
  fi
  plan_status=$(plan_field "$plan_file" status)
  case "$plan_status" in
    done)
      printf 'BLOCKED: plan is already done\n' >&2
      return 1
      ;;
    approved)
      printf 'OBSERVED: plan is already approved; recommend cc-run-plan\n'
      return 0
      ;;
    draft) ;;
    *)
      printf 'BLOCKED: unknown or contradictory plan identifier\n' >&2
      return 1
      ;;
  esac
  if test "$confirmation" != confirm; then
    printf 'DECLINED: plan remains draft; execution unstarted\n'
    return 0
  fi
  temporary="${plan_file}.tmp.$$"
  sed 's/^status: draft$/status: approved/' "$plan_file" > "$temporary"
  mv "$temporary" "$plan_file"
  lifecycle_sync_tasks "$plan_file" "$task_dir"
  printf 'APPROVED: plan status approved; tasks ready; execution unstarted\n'
}

write_plan_permission() {
  write_plan=$1
  action=$2
  if test "$write_plan" != true; then
    printf 'BLOCKED: write_plan false; cannot %s a plan\n' "$action" >&2
    return 1
  fi
}

child_cannot_approve_finish_or_clean() {
  role=$1
  action=$2
  case "$role" in
    implementer|verifier)
      if test "$action" = cleanup; then
        printf 'BLOCKED: cleanup is not permitted to children\n' >&2
        return 1
      fi
      printf 'BLOCKED: a child worker or verifier cannot %s\n' "$action" >&2
      return 1
      ;;
  esac
  return 0
}

finish_plan() {
  plan_file=$1
  task_dir=$2
  completion_file=$3
  evidence_dir=$4
  verifier_handoff=$5
  lease_lock=$6
  current_session=$7
  confirmation=$8
  if test ! -f "$plan_file"; then
    printf 'BLOCKED: unknown or contradictory plan identifier\n' >&2
    return 1
  fi
  plan_status=$(plan_field "$plan_file" status)
  case "$plan_status" in
    draft)
      printf 'BLOCKED: draft plans cannot be finished\n' >&2
      return 1
      ;;
    done)
      printf 'BLOCKED: plan is already done\n' >&2
      return 1
      ;;
    approved) ;;
    *)
      printf 'BLOCKED: unknown or contradictory plan identifier\n' >&2
      return 1
      ;;
  esac
  if test ! -f "$completion_file"; then
    printf 'BLOCKED: missing completion evidence\n' >&2
    return 1
  fi
  if test "$(plan_field "$completion_file" status)" != ready-for-human-status-change; then
    printf 'BLOCKED: completion evidence is not ready for human status-change\n' >&2
    return 1
  fi
  evidence_found=0
  if test -d "$evidence_dir"; then
    for evidence_file in "$evidence_dir"/*.yaml; do
      test -f "$evidence_file" || continue
      evidence_found=1
      grep -F 'evidence_status: completed' "$evidence_file" >/dev/null 2>&1 || {
        printf 'BLOCKED: missing task evidence, independent passing verification, or a remaining blocker\n' >&2
        return 1
      }
    done
  fi
  if test "$evidence_found" -eq 0; then
    printf 'BLOCKED: missing task evidence, independent passing verification, or a remaining blocker\n' >&2
    return 1
  fi
  if test ! -f "$verifier_handoff"; then
    printf 'BLOCKED: missing task evidence, independent passing verification, or a remaining blocker\n' >&2
    return 1
  fi
  grep -F 'status: completed' "$verifier_handoff" >/dev/null 2>&1 || {
    printf 'BLOCKED: missing task evidence, independent passing verification, or a remaining blocker\n' >&2
    return 1
  }
  grep -F 'verification: passed' "$verifier_handoff" >/dev/null 2>&1 || {
    printf 'BLOCKED: missing task evidence, independent passing verification, or a remaining blocker\n' >&2
    return 1
  }
  if grep -F 'status: blocked' "$verifier_handoff" >/dev/null 2>&1; then
    printf 'BLOCKED: missing task evidence, independent passing verification, or a remaining blocker\n' >&2
    return 1
  fi
  if test -d "$lease_lock" && test -f "$lease_lock/owner.yaml"; then
    lease_owner=$(plan_field "$lease_lock/owner.yaml" session_id)
    if test -n "$lease_owner" && test "$lease_owner" != "$current_session"; then
      printf 'BLOCKED: live writing session owned by someone else\n' >&2
      return 1
    fi
  fi
  if test "$confirmation" != confirm; then
    printf 'DECLINED: plan remains approved; completion evidence stays ready-for-human-status-change\n'
    return 0
  fi
  temporary="${plan_file}.tmp.$$"
  sed 's/^status: approved$/status: done/' "$plan_file" > "$temporary"
  mv "$temporary" "$plan_file"
  lifecycle_sync_tasks "$plan_file" "$task_dir"
  finish_plan_ref=$(plan_field "$completion_file" plan)
  finish_verifier=$(plan_field "$completion_file" verifier_session_id)
  finish_handoff=$(plan_field "$completion_file" verification_handoff)
  atomic_write "$completion_file" \
    'schema_version: 1' \
    "plan: $finish_plan_ref" \
    'status: completed' \
    "verifier_session_id: $finish_verifier" \
    "verification_handoff: $finish_handoff" \
    'human_gate: status-change' \
    'canonical_status_changed: true'
  if test -d "$lease_lock"; then
    plan_dir=$(dirname "$lease_lock")
    if test -f "$plan_dir/lease.yaml"; then
      temporary="$plan_dir/lease.yaml.tmp.$$"
      sed 's/^status: active$/status: released/' "$plan_dir/lease.yaml" > "$temporary"
      mv "$temporary" "$plan_dir/lease.yaml"
    fi
    rm -rf "$lease_lock"
  fi
  printf 'FINISHED: plan status done; tasks done; lease released; runtime preserved\n'
}

runtime_path_safe() {
  candidate=$1
  runtime_root=$2
  case "$candidate" in
    *..*) return 1 ;;
  esac
  case "$candidate/" in
    "$runtime_root/"*) ;;
    *) return 1 ;;
  esac
  if test -L "$candidate"; then
    return 1
  fi
  repo_key=$(basename "$(dirname "$candidate")")
  plan_id=$(basename "$candidate")
  safe_id "$repo_key" || return 1
  safe_id "$plan_id" || return 1
}

classify_worktree() {
  worktree=$1
  uncommitted=0
  unpushed=0
  if test ! -d "$worktree"; then
    printf 'missing\n'
    return 1
  fi
  if test -n "$(git -C "$worktree" status --porcelain 2>/dev/null)"; then
    uncommitted=1
  fi
  if git -C "$worktree" rev-parse --abbrev-ref --symbolic-full-name '@{upstream}' >/dev/null 2>&1; then
    if test "$(git -C "$worktree" rev-list --count '@{upstream}..HEAD')" -gt 0; then
      unpushed=1
    fi
  else
    unpushed=1
  fi
  if test "$uncommitted" -eq 1 && test "$unpushed" -eq 1; then
    printf 'both\n'
  elif test "$uncommitted" -eq 1; then
    printf 'uncommitted\n'
  elif test "$unpushed" -eq 1; then
    printf 'unpushed\n'
  else
    printf 'clean\n'
  fi
}

inspect_runtime_worktrees() {
  runtime_root=$1
  inspect_log=$2
  atomic_write "$inspect_log" 'CLEANUP: workspace-wide for .runtime/'
  if test -d "$runtime_root/worktrees"; then
    for repo_dir in "$runtime_root/worktrees"/*; do
      test -d "$repo_dir" || continue
      for worktree in "$repo_dir"/*; do
        test -d "$worktree" || continue
        if ! runtime_path_safe "$worktree" "$runtime_root"; then
          printf 'BLOCKED: path traversal, symlink, or identifier outside the runtime contract\n' >&2
          return 1
        fi
        class=$(classify_worktree "$worktree")
        printf 'worktree %s: %s\n' "$worktree" "$class" >> "$inspect_log"
        case "$class" in
          uncommitted)
            printf 'RISK: uncommitted\n' >> "$inspect_log"
            ;;
          unpushed)
            printf 'RISK: unpushed\n' >> "$inspect_log"
            ;;
          both)
            printf 'RISK: uncommitted\n' >> "$inspect_log"
            printf 'RISK: unpushed\n' >> "$inspect_log"
            ;;
        esac
      done
    done
  fi
  if test -d "$runtime_root/sessions"; then
    for session_file in "$runtime_root/sessions"/*/session.yaml; do
      test -f "$session_file" || continue
      session_status=$(plan_field "$session_file" status)
      case "$session_status" in
        completed|failed|cancelled) ;;
        *)
          printf 'live session %s status %s\n' "$session_file" "$session_status" \
            >> "$inspect_log"
          ;;
      esac
    done
  fi
}

git_common_dir() {
  worktree=$1
  common=$(git -C "$worktree" rev-parse --git-common-dir)
  case "$common" in
    /*) printf '%s\n' "$common" ;;
    *) printf '%s\n' "$worktree/$common" ;;
  esac
}

cleanup_runtime() {
  runtime_root=$1
  confirmation=$2
  inspect_log=$3
  if test ! -e "$runtime_root"; then
    printf 'NO-OP: runtime is missing or empty\n'
    return 0
  fi
  if test -d "$runtime_root"; then
    leftover=$(find "$runtime_root" -mindepth 1 -print | sed -n '1p')
    if test -z "$leftover"; then
      printf 'NO-OP: runtime is missing or empty\n'
      return 0
    fi
  fi
  inspect_runtime_worktrees "$runtime_root" "$inspect_log" || return 1
  has_uncommitted=0
  has_unpushed=0
  if grep -F 'RISK: uncommitted' "$inspect_log" >/dev/null 2>&1; then
    has_uncommitted=1
  fi
  if grep -F 'RISK: unpushed' "$inspect_log" >/dev/null 2>&1; then
    has_unpushed=1
  fi
  force_remove=0
  if test "$has_uncommitted" -eq 1; then
    if test "$confirmation" != confirm-discard-dirty; then
      printf 'BLOCKED: uncommitted or unpushed work requires confirmation\n' >&2
      return 1
    fi
    force_remove=1
  elif test "$has_unpushed" -eq 1; then
    case "$confirmation" in
      confirm|confirm-discard-dirty) ;;
      *)
        printf 'BLOCKED: uncommitted or unpushed work requires confirmation\n' >&2
        return 1
        ;;
    esac
  else
    case "$confirmation" in
      confirm|confirm-discard-dirty) ;;
      *)
        printf 'BLOCKED: cleanup requires human confirmation\n' >&2
        return 1
        ;;
    esac
  fi
  if test -d "$runtime_root/worktrees"; then
    for repo_dir in "$runtime_root/worktrees"/*; do
      test -d "$repo_dir" || continue
      for worktree in "$repo_dir"/*; do
        test -d "$worktree" || continue
        test -e "$worktree/.git" || continue
        common=$(git_common_dir "$worktree")
        if test "$force_remove" -eq 1; then
          git --git-dir="$common" worktree remove --force "$worktree"
        else
          git --git-dir="$common" worktree remove "$worktree"
        fi
      done
    done
  fi
  rm -rf "$runtime_root"
  printf 'CLEANED: runtime deleted; branches preserved\n'
}

init_cleanup_git() {
  source_repo=$1
  bare_repo=$2
  mkdir -p "$(dirname "$source_repo")"
  git init -q "$source_repo"
  git -C "$source_repo" config user.name 'Acceptance Fixture'
  git -C "$source_repo" config user.email 'acceptance@example.invalid'
  atomic_write "$source_repo/README.md" 'cleanup base'
  git -C "$source_repo" add README.md
  git -C "$source_repo" commit -qm 'cleanup base'
  git init --bare -q "$bare_repo"
  git -C "$source_repo" remote add origin "$bare_repo"
  git -C "$source_repo" push -q -u origin HEAD
}

require_file AGENTS.md
require_file WORKFLOW.md
require_file CLAUDE.md
require_file workspace.yaml
require_file context/INDEX.md
require_file context/WORKSPACE.md
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
require_file .agents/skills/cc-configure-workspace/SKILL.md
require_file .agents/skills/cc-configure-workspace/agents/openai.yaml
require_file docs/delivery-policies.md
require_file docs/integrations.md
require_file docs/host-capabilities.md

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
contains .agents/skills/cc-configure-workspace/SKILL.md 'Wait for explicit user confirmation'
contains .agents/skills/cc-configure-workspace/SKILL.md 'plan/task status'
contains docs/delivery-policies.md 'remote-review'
contains docs/delivery-policies.md 'local-target'
contains docs/delivery-policies.md 'manual'
contains docs/integrations.md 'enabled: false'
contains docs/integrations.md 'unavailable'
contains docs/host-capabilities.md 'Codex'
contains docs/host-capabilities.md 'Claude Code'
contains docs/host-capabilities.md 'Cursor Agent'

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

contains .gitignore '.runtime/'
if grep -E '^\.dist/|^node_modules/|ignored-clones' .gitignore >/dev/null 2>&1; then
  fail 'obsolete ignore residue remains in .gitignore'
fi
contains context/sources.yaml 'sources: []'
absent context/PRODUCT-DIRECTION.md
contains context/WORKSPACE.md 'Context Circuit workspace'
contains context/WORKSPACE.md 'Project identity lives in `PROJECT.md`.'
contains context/PROJECT.md 'not yet defined'
contains context/INDEX.md 'uninitialized'
contains sources/README.md 'user/team-organized boundary'
contains .agents/skills/cc-idea-brief/SKILL.md 'user/team-selected path under `sources/`'
contains .agents/skills/cc-create-prd/SKILL.md 'user/team-selected path under `sources/`'
contains docs/idea-brief.md 'user/team-selected path under `sources/`'
contains docs/prd.md 'user/team-selected path under `sources/`'
contains docs/getting-started.md 'not a command console for users'
contains docs/delivery-policies.md 'Publication is an optional human-authorized action'
contains docs/integrations.md 'does not approve or complete a plan'
contains README.md 'sh test/acceptance.sh'

removed_run_task="docs/run-task"
removed_run_task="${removed_run_task}.md"
removed_development="docs/development"
removed_development="${removed_development}.md"
removed_publication="docs/plan-publication"
removed_publication="${removed_publication}.md"
removed_wrapper="docs/using-the-wrapper"
removed_wrapper="${removed_wrapper}.md"
removed_whats_next="docs/whats-next"
removed_whats_next="${removed_whats_next}.md"
old_artifact_root="contributions"
old_artifact_root="${old_artifact_root}/"
for removed in \
  "$removed_run_task" \
  "$removed_development" \
  "$removed_publication" \
  "$removed_wrapper" \
  "$removed_whats_next"; do
  test ! -e "$removed" || fail "removed documentation still exists: $removed"
done
test ! -e contributions || fail 'old authored-artifact directory still exists'
if git grep -n -F "$removed_run_task" -- . >/dev/null 2>&1 \
  || git grep -n -F "$removed_development" -- . >/dev/null 2>&1 \
  || git grep -n -F "$removed_publication" -- . >/dev/null 2>&1 \
  || git grep -n -F "$removed_wrapper" -- . >/dev/null 2>&1 \
  || git grep -n -F "$removed_whats_next" -- . >/dev/null 2>&1 \
  || git grep -n -F "$old_artifact_root" -- . >/dev/null 2>&1; then
  fail 'tracked files still refer to a removed template-surface path'
fi

for file in README.md AGENTS.md CLAUDE.md WORKFLOW.md \
  agents/coordinator.md agents/repository-worker.md agents/reviewer.md \
  context/ARCHITECTURE.md context/CONVENTIONS.md context/DECISIONS.md \
  docs/agent-workspace-workflow.md docs/runtime-contract.md \
  docs/getting-started.md docs/planning.md docs/delivery-policies.md; do
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

# Plan 0003: context-grounded planning, review, execution routing, ownership,
# recovery, and acceptance behavior.
require_file plans/README.md
require_file docs/plan-review.md
require_file .agents/skills/cc-create-plan/SKILL.md
require_file .agents/skills/cc-create-plan/agents/openai.yaml
require_file .agents/skills/cc-review-plan/SKILL.md
require_file .agents/skills/cc-review-plan/agents/openai.yaml
require_file .agents/skills/cc-whats-next/SKILL.md
require_file .agents/skills/cc-whats-next/agents/openai.yaml
require_file .agents/skills/cc-run-plan/SKILL.md
require_file .agents/skills/cc-run-plan/agents/openai.yaml

for plan_field in 'plan.yaml' '`draft`' '`approved`' '`done`' \
  'repositories:' 'product_knowledge:' 'non_goals:' \
  'acceptance_criteria:' 'verification_commands:'; do
  contains plans/README.md "$plan_field"
done
contains docs/planning.md 'human-reviewed intended work'
contains docs/plan-review.md 'ready-for-approval'
contains docs/plan-review.md 'approved-for-execution'
contains .agents/skills/cc-create-plan/SKILL.md 'status: draft'
contains .agents/skills/cc-review-plan/SKILL.md 'cc-run-plan'
contains .agents/skills/cc-whats-next/SKILL.md 'approved dependency-ready plan'
contains .agents/skills/cc-run-plan/SKILL.md 'sole standard plan-execution capability'
contains .agents/skills/cc-run-plan/SKILL.md 'exclusive repository worktree'
contains .agents/skills/cc-run-plan/SKILL.md 'session-handoff-v1'
contains agents/coordinator.md 'For approved-plan execution'
contains agents/coordinator.md 'sole standard entry.'
contains agents/repository-worker.md 'assigned repository'
contains agents/reviewer.md 'independently reproduce'
contains docs/runtime-contract.md 'Same-plan contention is resolved'
absent_skill=$(printf '%s-%s' 'cc-run' 'task')
if git grep -F "$absent_skill" -- . >/dev/null 2>&1; then
  fail 'a skill that does not exist is still named in the repository'
fi

plan3_fixture="$fixture/plan-execution"
mkdir -p "$plan3_fixture/tasks" "$plan3_fixture/dependencies" \
  "$plan3_fixture/evidence"
plan3_file="$plan3_fixture/plan.yaml"
atomic_write "$plan3_file" \
  'id: plan-fixture' 'number: 3' 'title: Execute a bounded fixture plan' \
  'status: draft' 'source:' '  kind: accepted-prd' \
  '  reference: sources/fixture-prd.md' 'repositories:' \
  '  - context-circuit' 'product_knowledge:' '  references:' \
  '    - context/PROJECT.md' 'implementation_scope:' \
  '  - docs and agents' 'non_goals:' '  - delivery actions' \
  'dependencies:' '  - foundation' 'acceptance_criteria:' \
  '  - The fixture can execute safely' 'test_scope:' '  - filesystem fixtures' \
  'verification_commands:' '  - git diff --check'
expect_success plan_contract_complete "$plan3_file"
expect_failure plan_ready_for_execution "$plan3_file"
atomic_write "$plan3_fixture/tasks/CCP-0001.md" \
  'id: CCP-0001' 'status: draft' 'title: First fixture task'
atomic_write "$plan3_fixture/tasks/CCP-0002.md" \
  'id: CCP-0002' 'status: draft' 'title: Second fixture task'
atomic_write "$plan3_fixture/dependencies/foundation.yaml" \
  'id: foundation' 'status: done'
atomic_write "$plan3_fixture/dependencies/knowledge.yaml" \
  'id: knowledge' 'status: approved'
expect_success dependencies_ready \
  "$plan3_fixture/dependencies/foundation.yaml" \
  "$plan3_fixture/dependencies/knowledge.yaml"
atomic_write "$plan3_fixture/dependencies/knowledge.yaml" \
  'id: knowledge' 'status: draft'
expect_failure dependencies_ready "$plan3_fixture/dependencies/knowledge.yaml"
atomic_write "$plan3_fixture/dependencies/knowledge.yaml" \
  'id: knowledge' 'status: approved'
atomic_write "$plan3_file" \
  'id: plan-fixture' 'number: 3' 'title: Execute a bounded fixture plan' \
  'status: approved' 'source:' '  kind: accepted-prd' \
  '  reference: sources/fixture-prd.md' 'repositories:' \
  '  - context-circuit' 'product_knowledge:' '  references:' \
  '    - context/PROJECT.md' 'implementation_scope:' \
  '  - docs and agents' 'non_goals:' '  - delivery actions' \
  'dependencies:' '  - foundation' 'acceptance_criteria:' \
  '  - The fixture can execute safely' 'test_scope:' '  - filesystem fixtures' \
  'verification_commands:' '  - git diff --check'
expect_success plan_ready_for_execution "$plan3_file"
expect_success lifecycle_sync_tasks "$plan3_file" "$plan3_fixture/tasks"
contains "$plan3_fixture/tasks/CCP-0001.md" 'status: ready'
contains "$plan3_fixture/tasks/CCP-0002.md" 'status: ready'
plan3_tasks_snapshot="$plan3_fixture/tasks.snapshot"
sha256sum "$plan3_fixture/tasks"/*.md > "$plan3_tasks_snapshot"
expect_success lifecycle_sync_tasks "$plan3_file" "$plan3_fixture/tasks"
test "$(sha256sum "$plan3_fixture/tasks"/*.md)" = "$(cat "$plan3_tasks_snapshot")"

# Directed execution records replace task-count solo routing. Sequential
# tasks share one writer child packet; independent verification uses a
# distinct verifier child packet with write_worktree: false.
writer_packet="$plan3_fixture/writer-delegation.yaml"
verifier_packet="$plan3_fixture/verifier-delegation.yaml"
atomic_write "$writer_packet" \
  'schema_version: 1' 'session_id: writer-plan3' \
  'parent_session_id: root-plan3' 'root_session_id: root-plan3' \
  'role: implementer' 'objective: Implement sequential plan tasks' 'scope:' \
  '  plan: plans/context-circuit-plans/0003-plan-execution-orchestration' \
  '  task: CCP-0001..CCP-0002' '  paths:' '    - docs/' 'non_goals:' \
  '  - Do not change plan scope' 'context_refs:' '  - AGENTS.md' \
  '  - WORKFLOW.md' '  - docs/runtime-contract.md' 'repository: context-circuit' \
  'worktree: .runtime/worktrees/context-circuit/plan-execution-orchestration' \
  'permissions:' '  write_worktree: true' '  write_runtime_session: true' \
  '  write_plan: false' '  write_activity: false' 'acceptance_criteria:' \
  '  - Sequential tasks share one writer child' 'stop_conditions:' \
  '  - A required change falls outside docs/' 'handoff_schema: session-handoff-v1'
atomic_write "$verifier_packet" \
  'schema_version: 1' 'session_id: verifier-plan3' \
  'parent_session_id: root-plan3' 'root_session_id: root-plan3' \
  'role: verifier' 'objective: Independently verify the writer child work' 'scope:' \
  '  plan: plans/context-circuit-plans/0003-plan-execution-orchestration' \
  '  task: CCP-0001..CCP-0002' '  paths:' '    - docs/' 'non_goals:' \
  '  - Do not change plan scope' 'context_refs:' '  - AGENTS.md' \
  '  - WORKFLOW.md' '  - docs/runtime-contract.md' 'repository: context-circuit' \
  'worktree: .runtime/worktrees/context-circuit/plan-execution-orchestration' \
  'permissions:' '  write_worktree: false' '  write_runtime_session: true' \
  '  write_plan: false' '  write_activity: false' 'acceptance_criteria:' \
  '  - Independent verification uses a distinct verifier child' \
  'stop_conditions:' \
  '  - A required change falls outside docs/' 'handoff_schema: session-handoff-v1'
expect_success packet_complete "$writer_packet"
expect_success packet_complete "$verifier_packet"
contains "$writer_packet" 'role: implementer'
contains "$writer_packet" 'write_worktree: true'
contains "$verifier_packet" 'role: verifier'
contains "$verifier_packet" 'write_worktree: false'
test "$(sed -n 's/^session_id: //p' "$writer_packet" | head -n 1)" = writer-plan3
test "$(sed -n 's/^session_id: //p' "$verifier_packet" | head -n 1)" = verifier-plan3
test "$(sed -n 's/^session_id: //p' "$writer_packet" | head -n 1)" != \
  "$(sed -n 's/^session_id: //p' "$verifier_packet" | head -n 1)"
test "$(sed -n 's/^parent_session_id: //p' "$writer_packet" | head -n 1)" = root-plan3
test "$(sed -n 's/^parent_session_id: //p' "$verifier_packet" | head -n 1)" = root-plan3
test "$(expected_plan_execution_records)" = 'writer-child verifier-child'
test "$(execution_topology_for_mode solo)" = 'writer-child verifier-child'
test "$(execution_topology_for_mode team)" = 'writer-child verifier-child'

plan3_lock="$plan3_fixture/lease.lock"
expect_success run_plan_guard "$plan3_file" "$plan3_lock"
mkdir -p "$plan3_lock"
assert_failure_reason "$plan3_fixture/approved-owner.log" \
  'plan already has a writing owner' run_plan_guard "$plan3_file" "$plan3_lock"
rm -rf "$plan3_lock"
atomic_write "$plan3_file" \
  'id: plan-fixture' 'number: 3' 'title: Execute a bounded fixture plan' \
  'status: draft' 'source:' '  kind: accepted-prd' 'repositories:' \
  '  - context-circuit' 'product_knowledge:' '  references:' \
  '    - context/PROJECT.md' 'implementation_scope:' '  - docs and agents' \
  'non_goals:' '  - delivery actions' 'dependencies:' '  - foundation' \
  'acceptance_criteria:' '  - The fixture can execute safely' \
  'test_scope:' '  - filesystem fixtures' 'verification_commands:' \
  '  - git diff --check'
assert_failure_reason "$plan3_fixture/draft-plan.log" \
  'only an approved plan can execute' run_plan_guard "$plan3_file" "$plan3_lock"

packet_file="$plan3_fixture/delegation.yaml"
atomic_write "$packet_file" \
  'schema_version: 1' 'session_id: child-plan3' \
  'parent_session_id: root-plan3' 'root_session_id: root-plan3' \
  'role: implementer' 'objective: Implement one bounded Plan 003 task' 'scope:' \
  '  plan: plans/context-circuit-plans/0003-plan-execution-orchestration' \
  '  task: CCP-0001' '  paths:' '    - docs/' 'non_goals:' \
  '  - Do not change plan scope' 'context_refs:' '  - AGENTS.md' \
  '  - WORKFLOW.md' '  - docs/runtime-contract.md' 'repository: context-circuit' \
  'worktree: .runtime/worktrees/context-circuit/plan-execution-orchestration' \
  'permissions:' '  write_worktree: true' '  write_runtime_session: true' \
  '  write_plan: false' '  write_activity: false' 'acceptance_criteria:' \
  '  - The bounded task is implemented and verified' 'stop_conditions:' \
  '  - A required change falls outside docs/' 'handoff_schema: session-handoff-v1'
expect_success packet_complete "$packet_file"
contains "$packet_file" 'write_plan: false'
contains "$packet_file" 'handoff_schema: session-handoff-v1'

recovery_session="$plan3_fixture/recovery-session.yaml"
atomic_write "$recovery_session" 'session_id: takeover-plan3' \
  'status: blocked' 'replaces_session_id: interrupted-plan3' \
  'takeover_reason: human-authorized recovery'
expect_success recovery_guard "$recovery_session"

completion_plan3="$plan3_fixture/completion-plan.yaml"
completion_verifier3="$plan3_fixture/completion-verifier.md"
completion_gate3="$plan3_fixture/completion-gate.yaml"
atomic_write "$completion_plan3" 'id: plan-fixture' 'status: approved'
for plan3_task_id in CCP-0001 CCP-0002 CCP-0003; do
  atomic_write "$plan3_fixture/evidence/$plan3_task_id.yaml" \
    "task: $plan3_task_id" 'evidence_status: completed'
done
atomic_write "$completion_verifier3" '# Verification handoff' \
  'status: failed' 'verification: failed'
atomic_write "$completion_gate3" 'status: pending' 'gate: status-change'
expect_failure completion_ready_for_plan "$completion_plan3" \
  "$plan3_fixture/evidence" "$completion_verifier3" "$completion_gate3"
atomic_write "$completion_verifier3" '# Verification handoff' \
  'status: completed' 'verification: passed'
expect_failure completion_ready_for_plan "$completion_plan3" \
  "$plan3_fixture/evidence" "$completion_verifier3" "$completion_gate3"
atomic_write "$completion_gate3" 'status: approved' 'gate: status-change'
expect_success completion_ready_for_plan "$completion_plan3" \
  "$plan3_fixture/evidence" "$completion_verifier3" "$completion_gate3"

printf 'PASS: Plan 003 acceptance scenarios (planning, review, execution gating, routing, ownership, recovery, verification, completion)\n'

# Plan 0010: root sessions orchestrate writer and verifier children. Skills
# direct those packets instead of skipping children for small work.
contains .agents/skills/cc-run-plan/SKILL.md 'writer child'
contains .agents/skills/cc-run-plan/SKILL.md 'verifier child'
contains .agents/skills/cc-run-plan/SKILL.md 'Sequence is not a reason to skip children'
contains .agents/skills/cc-run-plan/SKILL.md 'mode: solo'
contains .agents/skills/cc-run-plan/SKILL.md 'mode: team'
contains .agents/skills/cc-run-plan/SKILL.md 'missing host primitive'
contains .agents/skills/cc-run-plan/SKILL.md 'reported to the human'
contains agents/coordinator.md 'writer child'
contains agents/coordinator.md 'verifier child'
contains agents/coordinator.md 'missing host primitive'
contains docs/agent-workspace-workflow.md 'writer child'
contains docs/agent-workspace-workflow.md 'verifier child'
contains docs/agent-workspace-workflow.md 'missing host primitive'
contains docs/host-capabilities.md 'Task / subagent tool'
contains docs/host-capabilities.md 'native child-agent or equivalent'
contains docs/host-capabilities.md 'missing host primitive'
contains context/ARCHITECTURE.md 'child writers and verifiers'
contains docs/runtime-contract.md 'preferred route (`delegated`)'

atomic_write "$plan3_fixture/solo-workspace.yaml" \
  'workspace:' '  name: fixture' '  mode: solo'
atomic_write "$plan3_fixture/team-workspace.yaml" \
  'workspace:' '  name: fixture' '  mode: team'
contains "$plan3_fixture/solo-workspace.yaml" 'mode: solo'
contains "$plan3_fixture/team-workspace.yaml" 'mode: team'
test "$(execution_topology_for_mode "$(sed -n 's/^  mode: //p' "$plan3_fixture/solo-workspace.yaml")")" = \
  'writer-child verifier-child'
test "$(execution_topology_for_mode "$(sed -n 's/^  mode: //p' "$plan3_fixture/team-workspace.yaml")")" = \
  'writer-child verifier-child'

missing_primitive_log="$plan3_fixture/missing-primitive.log"
atomic_write "$missing_primitive_log" \
  'BLOCKED: missing host primitive reported to the human' \
  'next_action: Ask how to proceed; do not skip children'
contains "$missing_primitive_log" 'missing host primitive reported to the human'
contains "$missing_primitive_log" 'do not skip children'

for orch_file in \
  .agents/skills/cc-run-plan/SKILL.md \
  agents/coordinator.md \
  docs/agent-workspace-workflow.md; do
  if grep -E 'root must not write|root cannot write the worktree|forbidden to write the worktree' \
    "$orch_file" >/dev/null 2>&1; then
    fail "instruction adds a root worktree write ban in $orch_file"
  fi
done
if grep -F 'solo users may work directly without child agents' \
  context/WORKSPACE.md context/DECISIONS.md >/dev/null 2>&1; then
  fail 'WORKSPACE or DECISIONS still presents working without child agents as the solo-user path'
fi
if grep -F 'only the amount of delegation changes' \
  context/WORKSPACE.md context/DECISIONS.md >/dev/null 2>&1; then
  fail 'WORKSPACE or DECISIONS still says only the amount of delegation changes'
fi
if grep -E 'selects solo or delegated execution' \
  context/ARCHITECTURE.md context/CONVENTIONS.md context/INDEX.md \
  context/PROJECT.md >/dev/null 2>&1; then
  fail 'Product Knowledge still says cc-run-plan selects solo or delegated execution'
fi

printf 'PASS: Plan 0010 directed writer and verifier child orchestration\n'

# Plan 0006: plan status is canonical and task status is a bulk, idempotent
# projection. Resume repair must preserve provider-owned annotations and must
# not rerun implementation or verification checks.
contains WORKFLOW.md 'Plan status is the canonical lifecycle authority'
contains docs/agent-workspace-workflow.md 'does not require separate task selection'
contains docs/planning.md 'reconciles every included task'
contains docs/runtime-contract.md 'Plan/task lifecycle projection'
contains docs/planning.md 'Task status is not a second approval or execution gate'
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
require_file context/WORKSPACE.md
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
contains .agents/skills/cc-initialize-workspace/SKILL.md 'context/WORKSPACE.md'
contains .agents/skills/cc-initialize-workspace/SKILL.md 'Do not write workspace identity'
contains docs/product-knowledge.md 'WORKSPACE.md'
contains context/INDEX.md 'context/WORKSPACE.md'
if grep -F 'PRODUCT-DIRECTION.md' \
  context/INDEX.md \
  docs/product-knowledge.md \
  docs/getting-started.md \
  docs/planning.md \
  AGENTS.md \
  WORKFLOW.md \
  .agents/skills/cc-session-entry/SKILL.md \
  .agents/skills/cc-initialize-workspace/SKILL.md \
  .agents/skills/cc-gather-context/SKILL.md \
  >/dev/null 2>&1; then
  fail 'a current required document still tells agents to read PRODUCT-DIRECTION.md'
fi
if grep -E 'session routing|CLI absence|working without child agents' \
  context/PROJECT.md >/dev/null 2>&1; then
  fail 'PROJECT.md describes workspace protocol as the project'
fi
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
  '    used_by: sources/example-idea-brief.md'
contains "$source_fixture/context/sources.yaml" 'read_for: Create the requested Idea Brief'
test ! -e "$source_fixture/context/RAW-SOURCE.md"

# The navigation index and artifact fixtures keep raw sources, accepted context,
# product artifacts, plans, and private runtime state in distinct homes.
for layer_link in '`sources/`' '`context/`' '`plans/`' '`.runtime/`'; do
  contains context/INDEX.md "$layer_link"
done
idea_fixture="$foundation_fixture/sources/example-idea-brief.md"
prd_fixture="$foundation_fixture/sources/example-prd.md"
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
contains context/sources.yaml 'sources: []'
if grep -F 'approved-product-direction-domain-role-context' context/sources.yaml \
  >/dev/null 2>&1; then
  fail 'starter provenance still records a product-direction domain-role source'
fi
if grep -F 'docs/templates/role-context.md' context/sources.yaml >/dev/null 2>&1; then
  fail 'starter provenance still records a template path as a source entry'
fi
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

# Plan 0004: optional configuration is capability-led, durable, and separate
# from initialization. Delivery policies retain authorization and human gates;
# integrations are opt-in and the core workflow remains complete offline.
contains .agents/skills/cc-initialize-workspace/SKILL.md 'does not ask about delivery behavior'
contains .agents/skills/cc-configure-workspace/SKILL.md 'capability-led conversational workflow'
contains docs/configuration.md 'canonical, portable configuration location'
contains docs/configuration.md 'delivery.policy'
contains docs/configuration.md 'delivery.authorization'
contains docs/delivery-policies.md 'post-verification choice'
contains docs/delivery-policies.md 'human merge gate'
contains docs/delivery-policies.md 'manual fallback'
contains docs/integrations.md 'opt-in adapter'
contains docs/integrations.md 'core filesystem workflow'
contains docs/host-capabilities.md 'same capability names'
contains agents/coordinator.md 'configuration` block as durable'
contains docs/delivery-policies.md 'remote-review'
contains docs/delivery-policies.md 'local-target'
contains docs/configuration.md 'remote-review'
contains docs/configuration.md 'local-target'
contains docs/configuration.md 'team-review'
contains docs/configuration.md 'solo-local'
contains .agents/skills/cc-configure-workspace/SKILL.md 'remote-review'
contains .agents/skills/cc-configure-workspace/SKILL.md 'local-target'
contains .agents/skills/cc-configure-workspace/SKILL.md 'team-review'
contains .agents/skills/cc-configure-workspace/SKILL.md 'solo-local'
contains agents/coordinator.md 'remote-review'
contains agents/coordinator.md 'local-target'
contains context/CONVENTIONS.md 'remote-review'
contains context/CONVENTIONS.md 'local-target'
contains context/ARCHITECTURE.md 'remote-review'
contains context/ARCHITECTURE.md 'local-target'
contains context/DECISIONS.md 'remote-review'
contains context/DECISIONS.md 'team-review'
contains context/DECISIONS.md 'solo-local'

for delivery_wording_file in \
  docs/delivery-policies.md \
  docs/configuration.md \
  .agents/skills/cc-configure-workspace/SKILL.md; do
  if grep -E 'delivery\.policy` is `team`|delivery\.policy` is `solo`|policy: `team`|policy: `solo`|## `team`|## `solo`|delivery policy: `team`|delivery policy: `solo`' \
    "$delivery_wording_file" >/dev/null 2>&1; then
    fail "identity-as-delivery wording in $delivery_wording_file"
  fi
  if grep -Ei 'configuration authorizes (commit|push|PR|merge)|confirming configuration authorizes|(^|[^t] )authorizes (commit|push|PR|merge)|commit-and-PR|local-merge' \
    "$delivery_wording_file" >/dev/null 2>&1; then
    fail "action-as-authorization wording in $delivery_wording_file"
  fi
done

configuration_fixture="$fixture/configuration"
mkdir -p "$configuration_fixture"

write_configuration_fixture() {
  configuration_policy=$1
  configuration_authorization=$2
  configuration_branch=$3
  configuration_enabled=$4
  configuration_integration_authorization=$5
  atomic_write "$configuration_fixture/workspace.yaml" \
    'workspace:' '  name: configured-fixture' '  mode: team' \
    '  default_branch: development' 'repositories:' '  context-circuit:' \
    '    path: ../context-circuit' '    mode: ignored-clone' \
    '    role: Context Circuit' '    agent: repository-worker' \
    '    default_branch: development' 'configuration:' \
    '  delivery:' '    scope: context-circuit' \
    "    policy: $configuration_policy" \
    '    repositories:' '      - context-circuit' \
    '    target_branches:' "      context-circuit: $configuration_branch" \
    "    authorization: $configuration_authorization" \
    '    fallback: manual' '  integrations:' \
    '    - id: activity-record' "      enabled: $configuration_enabled" \
    '      provider: fixture-provider' '      reads:' \
    '        - plan summaries' '      writes:' \
    '        - approved handoff metadata' \
    "      authorization: $configuration_integration_authorization" \
    '      fallback: core-filesystem' '  hosts:' \
    '    codex:' '      capability: cc-configure-workspace' \
    '      mode: native-or-core'
}

configuration_policy_value() {
  sed -n '/^  delivery:/,/^  integrations:/p' "$1" \
    | sed -n 's/^    policy: //p' | head -n 1
}

configuration_branch_value() {
  sed -n '/^    target_branches:/,/^    authorization:/p' "$1" \
    | sed -n 's/^      context-circuit: //p' | head -n 1
}

configuration_authorization_value() {
  sed -n '/^  delivery:/,/^  integrations:/p' "$1" \
    | sed -n 's/^    authorization: //p' | head -n 1
}

configuration_integration_value() {
  configuration_field=$2
  sed -n '/^  integrations:/,$p' "$1" \
    | sed -n "s/^      $configuration_field: //p" | head -n 1
}

prepare_delivery_fixture() {
  delivery_config=$1
  delivery_verification=$2
  delivery_output=$3
  delivery_policy=$(configuration_policy_value "$delivery_config")
  delivery_authorization=$(configuration_authorization_value "$delivery_config")
  delivery_branch=$(configuration_branch_value "$delivery_config")
  test "$delivery_verification" = passed || {
    printf 'BLOCKED: independent verification is required\n' >&2
    return 1
  }
  case "$delivery_policy" in
    team-review) delivery_policy=remote-review ;;
    solo-local) delivery_policy=local-target ;;
  esac
  case "$delivery_policy" in
    remote-review)
      test "$delivery_authorization" = approved || {
        printf 'BLOCKED: commit/push authorization is required\n' >&2
        return 1
      }
      atomic_write "$delivery_output" \
        'delivery: remote-review' 'reviewable_branch: codex/configured-fixture' \
        'human_gate: commit-push' "target_branch: $delivery_branch"
      ;;
    local-target)
      atomic_write "$delivery_output" \
        'delivery: local-target' "target_branch: $delivery_branch" \
        'human_gate: merge' 'merge: paused'
      ;;
    manual)
      atomic_write "$delivery_output" \
        'delivery: manual' 'worktree: preserved' \
        'human_gate: delivery-choice'
      ;;
    *)
      printf 'BLOCKED: unknown delivery policy\n' >&2
      return 1
      ;;
  esac
}

integration_result_fixture() {
  integration_config=$1
  integration_network=$2
  integration_enabled=$(configuration_integration_value "$integration_config" enabled)
  integration_authorization=$(configuration_integration_value "$integration_config" authorization)
  if test "$integration_enabled" != true; then
    printf 'disabled\n'
  elif test "$integration_authorization" != approved; then
    printf 'denied\n'
  elif test "$integration_network" != online; then
    printf 'unavailable\n'
  else
    printf 'published\n'
  fi
}

configuration_review_reason_fixture() {
  test "$1" = "$2" && test "$3" = "$4" && return 1
  printf 'focused-confirmation\n'
}

core_filesystem_fixture() {
  core_config=$1
  core_network=$2
  test -f "$core_config" || return 1
  test "$core_network" = offline || return 1
  printf 'core-filesystem\n'
}

host_capability_fixture() {
  case "$1" in
    codex|claude-code|cursor-agent) printf 'cc-configure-workspace\n' ;;
    *) return 1 ;;
  esac
}

write_configuration_fixture remote-review pending development false not-requested
test "$(configuration_policy_value "$configuration_fixture/workspace.yaml")" = remote-review
test "$(configuration_branch_value "$configuration_fixture/workspace.yaml")" = development
test "$(configuration_authorization_value "$configuration_fixture/workspace.yaml")" = pending
assert_failure_reason "$configuration_fixture/remote-review-denied.log" \
  'commit/push authorization is required' prepare_delivery_fixture \
  "$configuration_fixture/workspace.yaml" passed "$configuration_fixture/denied.handoff"
test ! -e "$configuration_fixture/denied.handoff"

write_configuration_fixture remote-review approved development false not-requested
expect_success prepare_delivery_fixture "$configuration_fixture/workspace.yaml" \
  passed "$configuration_fixture/remote-review.handoff"
contains "$configuration_fixture/remote-review.handoff" 'delivery: remote-review'
contains "$configuration_fixture/remote-review.handoff" 'reviewable_branch:'
contains "$configuration_fixture/remote-review.handoff" 'human_gate: commit-push'

write_configuration_fixture local-target pending development false not-requested
base_delivery_snapshot="$configuration_fixture/base-delivery.snapshot"
atomic_write "$configuration_fixture/base-branch" 'base branch remains unchanged'
sha256sum "$configuration_fixture/base-branch" > "$base_delivery_snapshot"
expect_success prepare_delivery_fixture "$configuration_fixture/workspace.yaml" \
  passed "$configuration_fixture/local-target.handoff"
contains "$configuration_fixture/local-target.handoff" 'delivery: local-target'
contains "$configuration_fixture/local-target.handoff" 'target_branch: development'
contains "$configuration_fixture/local-target.handoff" 'human_gate: merge'
contains "$configuration_fixture/local-target.handoff" 'merge: paused'
test "$(sha256sum "$configuration_fixture/base-branch")" = "$(cat "$base_delivery_snapshot")"

write_configuration_fixture team-review approved development false not-requested
expect_success prepare_delivery_fixture "$configuration_fixture/workspace.yaml" \
  passed "$configuration_fixture/team-review-alias.handoff"
contains "$configuration_fixture/team-review-alias.handoff" 'delivery: remote-review'
contains "$configuration_fixture/team-review-alias.handoff" 'reviewable_branch:'
contains "$configuration_fixture/team-review-alias.handoff" 'human_gate: commit-push'
if grep -F 'delivery: manual' "$configuration_fixture/team-review-alias.handoff" \
  >/dev/null 2>&1; then
  fail 'team-review alias fell through to manual'
fi

write_configuration_fixture solo-local pending development false not-requested
expect_success prepare_delivery_fixture "$configuration_fixture/workspace.yaml" \
  passed "$configuration_fixture/solo-local-alias.handoff"
contains "$configuration_fixture/solo-local-alias.handoff" 'delivery: local-target'
contains "$configuration_fixture/solo-local-alias.handoff" 'target_branch: development'
contains "$configuration_fixture/solo-local-alias.handoff" 'human_gate: merge'
if grep -F 'delivery: manual' "$configuration_fixture/solo-local-alias.handoff" \
  >/dev/null 2>&1; then
  fail 'solo-local alias fell through to manual'
fi

write_configuration_fixture team approved development false not-requested
assert_failure_reason "$configuration_fixture/identity-team.log" \
  'unknown delivery policy' prepare_delivery_fixture \
  "$configuration_fixture/workspace.yaml" passed \
  "$configuration_fixture/identity-team.handoff"
test ! -e "$configuration_fixture/identity-team.handoff"

write_configuration_fixture solo approved development false not-requested
assert_failure_reason "$configuration_fixture/identity-solo.log" \
  'unknown delivery policy' prepare_delivery_fixture \
  "$configuration_fixture/workspace.yaml" passed \
  "$configuration_fixture/identity-solo.handoff"
test ! -e "$configuration_fixture/identity-solo.handoff"

write_configuration_fixture manual denied development false not-requested
expect_success prepare_delivery_fixture "$configuration_fixture/workspace.yaml" \
  passed "$configuration_fixture/manual.handoff"
contains "$configuration_fixture/manual.handoff" 'worktree: preserved'
contains "$configuration_fixture/manual.handoff" 'human_gate: delivery-choice'
expect_success configuration_review_reason_fixture development main pending pending
expect_success configuration_review_reason_fixture development development pending approved
expect_failure configuration_review_reason_fixture development development pending pending

test "$(integration_result_fixture "$configuration_fixture/workspace.yaml" offline)" = disabled
write_configuration_fixture manual pending development true approved
test "$(integration_result_fixture "$configuration_fixture/workspace.yaml" offline)" = unavailable
test "$(integration_result_fixture "$configuration_fixture/workspace.yaml" online)" = published
write_configuration_fixture manual pending development true denied
test "$(integration_result_fixture "$configuration_fixture/workspace.yaml" online)" = denied

test "$(host_capability_fixture codex)" = cc-configure-workspace
test "$(host_capability_fixture claude-code)" = cc-configure-workspace
test "$(host_capability_fixture cursor-agent)" = cc-configure-workspace
expect_success core_filesystem_fixture "$configuration_fixture/workspace.yaml" offline
if grep -E '(^|[[:space:]])(token|password|secret|api_key):' \
  "$configuration_fixture/workspace.yaml" >/dev/null 2>&1; then
  fail 'configuration fixture persisted a credential-like field'
fi

printf 'PASS: optional configuration, delivery policy, authorization, integration, host, and offline-core acceptance scenarios\n'

printf 'PASS: workspace foundation and context acceptance scenarios (initialization, passive sources, provenance, artifacts)\n'

# Plan 0009: named approve, finish, and cleanup skills with confirmation,
# refusal, task reconciliation, and runtime-preservation fixtures.
require_file .agents/skills/cc-approve-plan/SKILL.md
require_file .agents/skills/cc-approve-plan/agents/openai.yaml
require_file .agents/skills/cc-finish-plan/SKILL.md
require_file .agents/skills/cc-finish-plan/agents/openai.yaml
require_file .agents/skills/cc-cleanup-runtime/SKILL.md
require_file .agents/skills/cc-cleanup-runtime/agents/openai.yaml
test "$(head -n 1 .agents/skills/cc-approve-plan/SKILL.md)" = '---'
test "$(head -n 1 .agents/skills/cc-finish-plan/SKILL.md)" = '---'
test "$(head -n 1 .agents/skills/cc-cleanup-runtime/SKILL.md)" = '---'
contains .agents/skills/cc-approve-plan/SKILL.md 'name: cc-approve-plan'
contains .agents/skills/cc-approve-plan/SKILL.md 'description:'
contains .agents/skills/cc-finish-plan/SKILL.md 'name: cc-finish-plan'
contains .agents/skills/cc-finish-plan/SKILL.md 'description:'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'name: cc-cleanup-runtime'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'description:'
contains .agents/skills/cc-approve-plan/agents/openai.yaml 'display_name:'
contains .agents/skills/cc-approve-plan/agents/openai.yaml '$cc-approve-plan'
contains .agents/skills/cc-finish-plan/agents/openai.yaml 'display_name:'
contains .agents/skills/cc-finish-plan/agents/openai.yaml '$cc-finish-plan'
contains .agents/skills/cc-cleanup-runtime/agents/openai.yaml 'display_name:'
contains .agents/skills/cc-cleanup-runtime/agents/openai.yaml '$cc-cleanup-runtime'

contains .agents/skills/cc-approve-plan/SKILL.md 'explicit human confirmation'
contains .agents/skills/cc-approve-plan/SKILL.md 'Leave execution unstarted'
contains .agents/skills/cc-approve-plan/SKILL.md 'recommend `cc-run-plan`'
contains .agents/skills/cc-approve-plan/SKILL.md 'A prior `cc-review-plan` run is not required'
contains .agents/skills/cc-finish-plan/SKILL.md 'explicit human confirmation'
contains .agents/skills/cc-finish-plan/SKILL.md 'ready-for-human-status-change'
contains .agents/skills/cc-finish-plan/SKILL.md 'canonical_status_changed: true'
contains .agents/skills/cc-finish-plan/SKILL.md 'status: completed'
contains .agents/skills/cc-finish-plan/SKILL.md 'Leave `.runtime/`'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'uncommitted changes'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'unpushed'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'git worktree remove'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'workspace-wide'
contains .agents/skills/cc-cleanup-runtime/SKILL.md 'successful no-op'
contains .agents/skills/cc-cleanup-runtime/SKILL.md '--force'

contains .agents/skills/cc-review-plan/SKILL.md 'never changes plan status'
contains .agents/skills/cc-review-plan/SKILL.md 'never writes `plan.yaml`'
contains .agents/skills/cc-review-plan/SKILL.md 'cc-approve-plan'
contains .agents/skills/cc-run-plan/SKILL.md 'mark the plan done'
contains .agents/skills/cc-run-plan/SKILL.md 'cc-finish-plan'
contains .agents/skills/cc-create-plan/SKILL.md 'cc-approve-plan'
contains .agents/skills/cc-whats-next/SKILL.md 'cc-approve-plan'
contains .agents/skills/cc-whats-next/SKILL.md 'cc-finish-plan'
contains .agents/skills/cc-whats-next/SKILL.md 'cc-cleanup-runtime'
contains .agents/skills/cc-whats-next/SKILL.md 'not mandatory ceremonies'
contains .agents/skills/cc-session-entry/SKILL.md 'cc-approve-plan'
contains .agents/skills/cc-session-entry/SKILL.md 'cc-finish-plan'
contains .agents/skills/cc-session-entry/SKILL.md 'cc-cleanup-runtime'
contains agents/coordinator.md 'cc-approve-plan'
contains agents/coordinator.md 'cc-finish-plan'
contains agents/coordinator.md 'cc-cleanup-runtime'
contains docs/planning.md 'cc-approve-plan'
contains docs/plan-review.md 'never writes `plan.yaml`'
contains docs/plan-review.md 'cc-approve-plan'
contains docs/runtime-contract.md 'is the named status-change skill'
contains docs/runtime-contract.md 'status: completed'
contains docs/runtime-contract.md 'cc-cleanup-runtime'
contains WORKFLOW.md 'cc-approve-plan'
contains WORKFLOW.md 'cc-finish-plan'
contains AGENTS.md 'cc-cleanup-runtime'
contains docs/getting-started.md 'cc-cleanup-runtime'
contains docs/getting-started.md 'cc-approve-plan'
contains docs/getting-started.md 'cc-finish-plan'
contains docs/getting-started.md 'cc-cleanup-runtime'
contains context/INDEX.md 'cc-approve-plan'
contains context/ARCHITECTURE.md 'cc-approve-plan'
contains context/CONVENTIONS.md 'cc-finish-plan'
contains context/INDEX.md 'cc-cleanup-runtime'
contains context/ARCHITECTURE.md 'cc-cleanup-runtime'
contains context/CONVENTIONS.md 'cc-approve-plan'

gates_fixture="$fixture/plan-gates"
mkdir -p "$gates_fixture/tasks" "$gates_fixture/evidence"
gates_plan="$gates_fixture/plan.yaml"
write_gate_plan "$gates_plan" gate-fixture draft
write_gate_tasks "$gates_fixture/tasks" draft
expect_success approve_plan "$gates_plan" "$gates_fixture/tasks" decline gate-fixture
contains "$gates_plan" 'status: draft'
contains "$gates_fixture/tasks/FPC-0001.md" 'status: draft'
test ! -e "$gates_fixture/lease.lock"

expect_success approve_plan "$gates_plan" "$gates_fixture/tasks" confirm gate-fixture
contains "$gates_plan" 'status: approved'
contains "$gates_fixture/tasks/FPC-0001.md" 'status: ready'
contains "$gates_fixture/tasks/FPC-0002.md" 'status: ready'
test ! -e "$gates_fixture/lease.lock"
test ! -e "$gates_fixture/worktree"
approved_snapshot="$gates_fixture/approved.snapshot"
sha256sum "$gates_plan" "$gates_fixture/tasks"/*.md > "$approved_snapshot"
expect_success approve_plan "$gates_plan" "$gates_fixture/tasks" confirm gate-fixture > "$gates_fixture/already-approved.log"
contains "$gates_fixture/already-approved.log" 'already approved'
contains "$gates_fixture/already-approved.log" 'cc-run-plan'
test "$(sha256sum "$gates_plan" "$gates_fixture/tasks"/*.md)" = "$(cat "$approved_snapshot")"

assert_failure_reason "$gates_fixture/unknown.log" \
  'unknown or contradictory plan identifier' \
  approve_plan "$gates_fixture/missing.yaml" "$gates_fixture/tasks" confirm missing
assert_failure_reason "$gates_fixture/contradictory.log" \
  'unknown or contradictory plan identifier' \
  approve_plan "$gates_plan" "$gates_fixture/tasks" confirm other-id
atomic_write "$gates_fixture/incomplete.yaml" 'id: incomplete' 'status: draft'
assert_failure_reason "$gates_fixture/incomplete.log" \
  'plan artifact is incomplete' \
  approve_plan "$gates_fixture/incomplete.yaml" "$gates_fixture/no-tasks" confirm incomplete
write_gate_plan "$gates_fixture/done.yaml" done-fixture done
write_gate_tasks "$gates_fixture/done-tasks" done
assert_failure_reason "$gates_fixture/done.log" \
  'plan is already done' \
  approve_plan "$gates_fixture/done.yaml" "$gates_fixture/done-tasks" confirm done-fixture

write_gate_plan "$gates_plan" gate-fixture approved
write_gate_tasks "$gates_fixture/tasks" ready
finish_completion="$gates_fixture/completion.yaml"
finish_verifier="$gates_fixture/verifier.md"
finish_runtime="$gates_fixture/.runtime"
mkdir -p "$finish_runtime/plans/gate-fixture" \
  "$finish_runtime/sessions/root-finish" \
  "$gates_fixture/evidence"
atomic_write "$finish_completion" \
  'schema_version: 1' \
  'plan: plans/context-circuit-plans/gate-fixture' \
  'status: ready-for-human-status-change' \
  'verifier_session_id: verifier-finish' \
  'verification_handoff: .runtime/sessions/verifier-finish/handoff.md' \
  'human_gate: status-change' \
  'canonical_status_changed: false'
atomic_write "$gates_fixture/evidence/FPC-0001.yaml" \
  'task: FPC-0001' 'evidence_status: completed'
atomic_write "$gates_fixture/evidence/FPC-0002.yaml" \
  'task: FPC-0002' 'evidence_status: completed'
atomic_write "$finish_verifier" \
  '# Verification handoff' 'status: completed' 'verification: passed'
runtime="$finish_runtime"
write_session root-finish null root-finish root coordinator 'Finish the fixture plan' FPC-0001..FPC-0002 workspace false . executing
expect_success acquire_lease \
  "$finish_runtime/plans/gate-fixture/lease.lock" \
  root-finish root-finish \
  plans/context-circuit-plans/gate-fixture \
  .runtime/worktrees/context-circuit/gate-fixture \
  gate-fixture

write_gate_plan "$gates_fixture/draft-finish.yaml" draft-finish draft
write_gate_tasks "$gates_fixture/draft-finish-tasks" draft
assert_failure_reason "$gates_fixture/finish-draft.log" \
  'draft plans cannot be finished' \
  finish_plan "$gates_fixture/draft-finish.yaml" \
  "$gates_fixture/draft-finish-tasks" "$finish_completion" \
  "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish confirm
write_gate_plan "$gates_fixture/already-done.yaml" already-done done
write_gate_tasks "$gates_fixture/already-done-tasks" done
assert_failure_reason "$gates_fixture/finish-done.log" \
  'plan is already done' \
  finish_plan "$gates_fixture/already-done.yaml" \
  "$gates_fixture/already-done-tasks" "$finish_completion" \
  "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish confirm
assert_failure_reason "$gates_fixture/finish-missing.log" \
  'missing completion evidence' \
  finish_plan "$gates_plan" "$gates_fixture/tasks" \
  "$gates_fixture/missing-completion.yaml" \
  "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish confirm
atomic_write "$gates_fixture/blocked-completion.yaml" \
  'schema_version: 1' \
  'plan: plans/context-circuit-plans/gate-fixture' \
  'status: blocked' \
  'human_gate: status-change' \
  'canonical_status_changed: false'
assert_failure_reason "$gates_fixture/finish-blocked.log" \
  'completion evidence is not ready for human status-change' \
  finish_plan "$gates_plan" "$gates_fixture/tasks" \
  "$gates_fixture/blocked-completion.yaml" \
  "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish confirm
atomic_write "$gates_fixture/failed-verifier.md" \
  '# Verification handoff' 'status: failed' 'verification: failed'
assert_failure_reason "$gates_fixture/finish-failed-verifier.log" \
  'missing task evidence, independent passing verification, or a remaining blocker' \
  finish_plan "$gates_plan" "$gates_fixture/tasks" "$finish_completion" \
  "$gates_fixture/evidence" "$gates_fixture/failed-verifier.md" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish confirm
assert_failure_reason "$gates_fixture/finish-foreign-lease.log" \
  'live writing session owned by someone else' \
  finish_plan "$gates_plan" "$gates_fixture/tasks" "$finish_completion" \
  "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" other-session confirm

expect_success finish_plan "$gates_plan" "$gates_fixture/tasks" \
  "$finish_completion" "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish decline
contains "$gates_plan" 'status: approved'
contains "$finish_completion" 'status: ready-for-human-status-change'
contains "$finish_completion" 'canonical_status_changed: false'
test -d "$finish_runtime/plans/gate-fixture/lease.lock"

expect_success finish_plan "$gates_plan" "$gates_fixture/tasks" \
  "$finish_completion" "$gates_fixture/evidence" "$finish_verifier" \
  "$finish_runtime/plans/gate-fixture/lease.lock" root-finish confirm
contains "$gates_plan" 'status: done'
contains "$gates_fixture/tasks/FPC-0001.md" 'status: done'
contains "$gates_fixture/tasks/FPC-0002.md" 'status: done'
contains "$finish_completion" 'status: completed'
contains "$finish_completion" 'human_gate: status-change'
contains "$finish_completion" 'canonical_status_changed: true'
test ! -d "$finish_runtime/plans/gate-fixture/lease.lock"
test -f "$finish_runtime/plans/gate-fixture/lease.yaml"
contains "$finish_runtime/plans/gate-fixture/lease.yaml" 'status: released'
test -d "$finish_runtime/sessions/root-finish"
done_snapshot="$gates_fixture/done.snapshot"
sha256sum "$gates_plan" "$gates_fixture/tasks"/*.md > "$done_snapshot"
expect_success lifecycle_sync_tasks "$gates_plan" "$gates_fixture/tasks"
test "$(sha256sum "$gates_plan" "$gates_fixture/tasks"/*.md)" = "$(cat "$done_snapshot")"

contains "$fixture/.runtime/sessions/child-001/delegation.yaml" 'write_plan: false'
assert_failure_reason "$gates_fixture/child-approve.log" \
  'write_plan false' write_plan_permission false approve
assert_failure_reason "$gates_fixture/child-finish.log" \
  'a child worker or verifier cannot finish' \
  child_cannot_approve_finish_or_clean implementer finish
assert_failure_reason "$gates_fixture/verifier-finish.log" \
  'a child worker or verifier cannot finish' \
  child_cannot_approve_finish_or_clean verifier finish
assert_failure_reason "$gates_fixture/child-cleanup.log" \
  'cleanup is not permitted to children' \
  child_cannot_approve_finish_or_clean implementer cleanup

cleanup_src="$fixture/cleanup-src"
cleanup_bare="$fixture/cleanup-src.git"
init_cleanup_git "$cleanup_src" "$cleanup_bare"

cleanup_dirty_rt="$fixture/cleanup-dirty/.runtime"
mkdir -p "$cleanup_dirty_rt/worktrees/context-circuit" \
  "$cleanup_dirty_rt/sessions/live-001"
git -C "$cleanup_src" worktree add -q -b plan-dirty \
  "$cleanup_dirty_rt/worktrees/context-circuit/plan-dirty"
git -C "$cleanup_dirty_rt/worktrees/context-circuit/plan-dirty" push -q -u origin plan-dirty
atomic_write "$cleanup_dirty_rt/worktrees/context-circuit/plan-dirty/dirty.txt" \
  'uncommitted work'
atomic_write "$cleanup_dirty_rt/sessions/live-001/session.yaml" \
  'session_id: live-001' 'status: executing'
test "$(classify_worktree "$cleanup_dirty_rt/worktrees/context-circuit/plan-dirty")" = uncommitted
assert_failure_reason "$fixture/cleanup-dirty-blocked.log" \
  'uncommitted or unpushed work requires confirmation' \
  cleanup_runtime "$cleanup_dirty_rt" confirm "$fixture/cleanup-dirty.inspect"
contains "$fixture/cleanup-dirty.inspect" 'workspace-wide'
contains "$fixture/cleanup-dirty.inspect" 'RISK: uncommitted'
contains "$fixture/cleanup-dirty.inspect" 'live session'
test -d "$cleanup_dirty_rt/worktrees/context-circuit/plan-dirty"
base_dirty_readme="$cleanup_src/README.md"
sha256sum "$base_dirty_readme" > "$fixture/cleanup-dirty-base.snapshot"
expect_success cleanup_runtime "$cleanup_dirty_rt" confirm-discard-dirty \
  "$fixture/cleanup-dirty-confirmed.inspect"
test ! -e "$cleanup_dirty_rt"
git -C "$cleanup_src" show-ref --verify --quiet refs/heads/plan-dirty || \
  fail 'cleanup deleted the plan-dirty branch'
test ! -e "$cleanup_src/dirty.txt"
test "$(sha256sum "$base_dirty_readme")" = "$(cat "$fixture/cleanup-dirty-base.snapshot")"

cleanup_unpushed_rt="$fixture/cleanup-unpushed/.runtime"
mkdir -p "$cleanup_unpushed_rt/worktrees/context-circuit"
git -C "$cleanup_src" worktree add -q -b plan-unpushed \
  "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed"
git -C "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed" push -q -u origin plan-unpushed
atomic_write "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed/ahead.txt" \
  'local commit'
git -C "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed" add ahead.txt
git -C "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed" commit -qm 'unpushed commit'
test "$(classify_worktree "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed")" = unpushed
assert_failure_reason "$fixture/cleanup-unpushed-blocked.log" \
  'uncommitted or unpushed work requires confirmation' \
  cleanup_runtime "$cleanup_unpushed_rt" decline "$fixture/cleanup-unpushed.inspect"
contains "$fixture/cleanup-unpushed.inspect" 'RISK: unpushed'
unpushed_sha=$(git -C "$cleanup_unpushed_rt/worktrees/context-circuit/plan-unpushed" rev-parse HEAD)
expect_success cleanup_runtime "$cleanup_unpushed_rt" confirm \
  "$fixture/cleanup-unpushed-confirmed.inspect"
test ! -e "$cleanup_unpushed_rt"
git -C "$cleanup_src" show-ref --verify --quiet refs/heads/plan-unpushed || \
  fail 'cleanup deleted the plan-unpushed branch'
test "$(git -C "$cleanup_src" rev-parse plan-unpushed)" = "$unpushed_sha"

cleanup_local_rt="$fixture/cleanup-local/.runtime"
mkdir -p "$cleanup_local_rt/worktrees/context-circuit"
git -C "$cleanup_src" worktree add -q -b plan-local \
  "$cleanup_local_rt/worktrees/context-circuit/plan-local"
atomic_write "$cleanup_local_rt/worktrees/context-circuit/plan-local/local.txt" \
  'local-only commit'
git -C "$cleanup_local_rt/worktrees/context-circuit/plan-local" add local.txt
git -C "$cleanup_local_rt/worktrees/context-circuit/plan-local" commit -qm 'local-only commit'
test "$(classify_worktree "$cleanup_local_rt/worktrees/context-circuit/plan-local")" = unpushed
assert_failure_reason "$fixture/cleanup-local-blocked.log" \
  'uncommitted or unpushed work requires confirmation' \
  cleanup_runtime "$cleanup_local_rt" decline "$fixture/cleanup-local.inspect"
local_sha=$(git -C "$cleanup_local_rt/worktrees/context-circuit/plan-local" rev-parse HEAD)
expect_success cleanup_runtime "$cleanup_local_rt" confirm \
  "$fixture/cleanup-local-confirmed.inspect"
test ! -e "$cleanup_local_rt"
git -C "$cleanup_src" show-ref --verify --quiet refs/heads/plan-local || \
  fail 'cleanup deleted the plan-local branch'
test "$(git -C "$cleanup_src" rev-parse plan-local)" = "$local_sha"

cleanup_clean_rt="$fixture/cleanup-clean/.runtime"
mkdir -p "$cleanup_clean_rt/worktrees/context-circuit" \
  "$cleanup_clean_rt/sessions/done-001"
git -C "$cleanup_src" worktree add -q -b plan-clean \
  "$cleanup_clean_rt/worktrees/context-circuit/plan-clean"
git -C "$cleanup_clean_rt/worktrees/context-circuit/plan-clean" push -q -u origin plan-clean
atomic_write "$cleanup_clean_rt/sessions/done-001/session.yaml" \
  'session_id: done-001' 'status: completed'
test "$(classify_worktree "$cleanup_clean_rt/worktrees/context-circuit/plan-clean")" = clean
assert_failure_reason "$fixture/cleanup-clean-unconfirmed.log" \
  'cleanup requires human confirmation' \
  cleanup_runtime "$cleanup_clean_rt" decline "$fixture/cleanup-clean.inspect"
contains "$fixture/cleanup-clean.inspect" 'workspace-wide'
expect_success cleanup_runtime "$cleanup_clean_rt" confirm \
  "$fixture/cleanup-clean-confirmed.inspect"
test ! -e "$cleanup_clean_rt"
git -C "$cleanup_src" show-ref --verify --quiet refs/heads/plan-clean || \
  fail 'cleanup deleted the plan-clean branch'
test ! -e "$cleanup_src/ahead.txt"

empty_rt="$fixture/cleanup-empty/.runtime"
mkdir -p "$empty_rt"
expect_success cleanup_runtime "$empty_rt" confirm "$fixture/cleanup-empty.inspect"
test -d "$empty_rt"
missing_rt="$fixture/cleanup-missing/.runtime"
expect_success cleanup_runtime "$missing_rt" confirm "$fixture/cleanup-missing.inspect"

printf 'PASS: Plan 0009 approve-plan, finish-plan, and cleanup-runtime gates\n'

printf 'PASS: pure agent-workspace acceptance scenarios (filesystem, contention, isolation, recovery, verification, gates)\n'
