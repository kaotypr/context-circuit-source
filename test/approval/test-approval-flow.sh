#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

init_repo() {
  repo=$1
  kind=$2
  mkdir -p "$repo/plans/app-plans/demo/tasks"
  cat > "$repo/workspace.yaml" <<EOF
version: 1
workspace:
  name: approval-flow-fixture
identity:
  kind: $kind
  status: accepted
EOF
  cat > "$repo/plans/app-plans/demo/plan.yaml" <<'EOF'
schema_version: 2
id: demo
status: draft
tasks: [one, two]
EOF
  for task_id in one two; do
    cat > "$repo/plans/app-plans/demo/tasks/$task_id.md" <<EOF
---
schema_version: 2
id: $task_id
status: draft
---
# $task_id

Keep this body status: draft example.
EOF
  done
  git init -q -b main "$repo"
  git -C "$repo" config user.email test@example.invalid
  git -C "$repo" config user.name 'Context Circuit Test'
  git -C "$repo" add workspace.yaml plans
  git -C "$repo" commit -qm initial
}

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-approval.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM

mkdir -p "$runtime/card/tasks"
cat > "$runtime/card/plan.yaml" <<'EOF'
schema_version: 2
id: card-fixture
status: draft
EOF
cat > "$runtime/card/tasks/TASK-01.md" <<'EOF'
---
id: TASK-01
status: draft
---
# Task
EOF
before_plan=$(cc_digest "$runtime/card/plan.yaml")
before_task=$(cc_digest "$runtime/card/tasks/TASK-01.md")
card=$(cc_approval_card card-fixture)
printf '%s\n' "$card" | grep -F 'Action: present-approval-card' >/dev/null || fail 'initial card used approve-plan'
printf '%s\n' "$card" | grep -F 'nothing has changed yet' >/dev/null || fail 'initial card omitted unchanged-state wording'
printf '%s\n' "$card" | grep -F 'Confirm approval of plan card-fixture.' >/dev/null || fail 'initial card omitted exact confirmation'
printf '%s\n' "$card" | grep -F 'does not commit Git' >/dev/null || fail 'initial card omitted commit non-effect'
printf '%s\n' "$card" | grep -F 'does not start Run approved plan' >/dev/null || fail 'initial card implied execution'
assert_eq "$(cc_digest "$runtime/card/plan.yaml")" "$before_plan"
assert_eq "$(cc_digest "$runtime/card/tasks/TASK-01.md")" "$before_task"
test ! -d "$runtime/worktrees" || fail 'card helper created a worktree'

expect_failure cc_confirm_approval "$runtime/card" "$runtime/card/plan.yaml" "$runtime/card/tasks" card-fixture
expect_failure cc_confirm_approval "$runtime/card" "$runtime/card/plan.yaml" "$runtime/card/tasks" card-fixture stale
expect_failure cc_confirm_approval "$runtime/card" "$runtime/card/plan.yaml" "$runtime/card/tasks" card-fixture yes
assert_eq "$(cc_digest "$runtime/card/plan.yaml")" "$before_plan"
assert_eq "$(cc_digest "$runtime/card/tasks/TASK-01.md")" "$before_task"

printf '%s' '---
id: TASK-01
status: draft
---
Body must keep this status: draft line.

# Task' > "$runtime/card/tasks/TASK-01.md"
cc_transition_plan_status "$runtime/card/plan.yaml" "$runtime/card/tasks" approved confirmed
contains "$runtime/card/plan.yaml" 'status: approved'
contains "$runtime/card/tasks/TASK-01.md" 'status: ready'
contains "$runtime/card/tasks/TASK-01.md" 'Body must keep this status: draft line.'
test "$(grep -c '^status: ready$' "$runtime/card/tasks/TASK-01.md")" -eq 1 || fail 'transition rewrote a body status line'
test "$(tail -c 1 "$runtime/card/tasks/TASK-01.md" | od -An -tx1 | tr -d ' \n')" != 0a || fail 'transition added a trailing newline'
expect_failure cc_transition_plan_status "$runtime/card/plan.yaml" "$runtime/card/tasks" approved confirmed

product="$runtime/product-source"
init_repo "$product" product-source
head_before=$(git -C "$product" rev-parse HEAD)
product_out=$(cc_confirm_approval "$product" "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks" demo confirmed)
printf '%s\n' "$product_out" | grep -F 'Approval is complete. Execution has not started.' >/dev/null || fail 'product-source omitted approval-complete wording'
printf '%s\n' "$product_out" | grep -F 'Action: commit-approved-plan' >/dev/null || fail 'product-source omitted commit card'
printf '%s\n' "$product_out" | grep -F 'Confirm commit of the approved plan state.' >/dev/null || fail 'product-source omitted commit confirmation'
printf '%s\n' "$product_out" | grep -F MAINTAINER_APPROVAL_COMMIT_REQUIRED >/dev/null || fail 'product-source omitted maintainer-commit class'
printf '%s\n' "$product_out" | grep -F 'Next action: Run approved plan' >/dev/null && fail 'product-source named Run as the immediate next action'
contains "$product/plans/app-plans/demo/plan.yaml" 'status: approved'
contains "$product/plans/app-plans/demo/tasks/one.md" 'status: ready'
contains "$product/plans/app-plans/demo/tasks/one.md" 'Keep this body status: draft example.'
assert_eq "$(git -C "$product" rev-parse HEAD)" "$head_before"
test ! -d "$product/.runtime/worktrees" || fail 'confirm approval created a worktree'
test ! -d "$product/.runtime/plans" || fail 'confirm approval created a lease'
assert_eq "$(cc_maintainer_approval_commit_required "$product" "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks")" MAINTAINER_APPROVAL_COMMIT_REQUIRED
expect_failure cc_prepare_worktree "$product" "$runtime/worktrees/product-blocked" main
printf '%s\n' unrelated > "$product/unrelated.txt"
expect_failure cc_maintainer_approval_commit_required "$product" "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks"
expect_failure cc_confirm_approval "$product" "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks" demo confirmed

dirty="$runtime/product-dirty"
init_repo "$dirty" product-source
printf '%s\n' extra > "$dirty/extra.txt"
dirty_out=$(cc_confirm_approval "$dirty" "$dirty/plans/app-plans/demo/plan.yaml" "$dirty/plans/app-plans/demo/tasks" demo confirmed)
printf '%s\n' "$dirty_out" | grep -F 'Approval is complete. Execution has not started.' >/dev/null || fail 'dirty product-source omitted approval-complete wording'
printf '%s\n' "$dirty_out" | grep -F DIRTY_BASE_BLOCKED >/dev/null || fail 'unrelated dirty files were not classified as DIRTY_BASE_BLOCKED'
printf '%s\n' "$dirty_out" | grep -F 'Action: commit-approved-plan' >/dev/null && fail 'unrelated dirt presented a commit card'
contains "$dirty/plans/app-plans/demo/plan.yaml" 'status: approved'
expect_failure cc_maintainer_approval_commit_required "$dirty" "$dirty/plans/app-plans/demo/plan.yaml" "$dirty/plans/app-plans/demo/tasks"

wrapped="$runtime/wrapped"
init_repo "$wrapped" instantiated-workspace
wrapped_head=$(git -C "$wrapped" rev-parse HEAD)
wrapped_out=$(cc_confirm_approval "$wrapped" "$wrapped/plans/app-plans/demo/plan.yaml" "$wrapped/plans/app-plans/demo/tasks" demo confirmed)
printf '%s\n' "$wrapped_out" | grep -F 'Approval is complete. Execution has not started.' >/dev/null || fail 'wrapped workspace omitted approval-complete wording'
printf '%s\n' "$wrapped_out" | grep -F 'Next action: Run approved plan demo' >/dev/null || fail 'wrapped workspace omitted Run as next action'
printf '%s\n' "$wrapped_out" | grep -F 'Action: commit-approved-plan' >/dev/null && fail 'wrapped workspace presented a commit card'
printf '%s\n' "$wrapped_out" | grep -F 'Confirm commit of the approved plan state.' >/dev/null && fail 'wrapped workspace asked for a maintainer commit'
printf '%s\n' "$wrapped_out" | grep -F MAINTAINER_APPROVAL_COMMIT_REQUIRED >/dev/null && fail 'wrapped workspace used maintainer-commit class'
assert_eq "$(git -C "$wrapped" rev-parse HEAD)" "$wrapped_head"
expect_failure cc_maintainer_approval_commit_required "$wrapped" "$wrapped/plans/app-plans/demo/plan.yaml" "$wrapped/plans/app-plans/demo/tasks"
contains "$ROOT/docs/gates.md" 'present this existing commit card in the same session immediately'
contains "$ROOT/agents/coordinator.md" 'cc_transition_plan_status'
pass 'two-turn approval is mutation-free, status-only, and sequenced before commit or run'
