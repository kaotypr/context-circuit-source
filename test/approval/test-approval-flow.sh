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

fixture_lines_in() {
  fixture=$1
  file=$2
  require_file "$fixture"
  while IFS= read -r line || [ -n "$line" ]; do
    test -n "$line" || continue
    contains "$file" "$line"
  done < "$fixture"
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
contains "$ROOT/docs/gates.md" 'Action: present-approval-card'
contains "$ROOT/docs/gates.md" 'nothing has changed yet'
contains "$ROOT/docs/gates.md" 'Confirm approval of plan <id>.'
contains "$ROOT/docs/gates.md" 'does not commit Git and does not start Run approved plan'
contains "$ROOT/docs/gates.md" 'cc_confirmation_card'
not_contains "$ROOT/wrapper/runtime/engine.sh" 'cc_approval_card'
not_contains "$ROOT/wrapper/runtime/engine.sh" 'cc_confirm_approval'
not_contains "$ROOT/wrapper/runtime/engine.sh" 'cc_commit_approved_plan_card'
card_route=$(cc_route 'Approve plan card-fixture.')
printf '%s\n' "$card_route" | grep -F 'capability: present-approval-card' >/dev/null || fail 'initial request was not a card'
printf '%s\n' "$card_route" | grep -F 'authorization: confirmed-gate-required' >/dev/null || fail 'initial request was treated as authorized'
assert_eq "$(cc_digest "$runtime/card/plan.yaml")" "$before_plan"
assert_eq "$(cc_digest "$runtime/card/tasks/TASK-01.md")" "$before_task"
test ! -d "$runtime/worktrees" || fail 'card route created a worktree'

expect_failure cc_transition_plan_status "$runtime/card/plan.yaml" "$runtime/card/tasks" approved
expect_failure cc_transition_plan_status "$runtime/card/plan.yaml" "$runtime/card/tasks" approved stale
expect_failure cc_transition_plan_status "$runtime/card/plan.yaml" "$runtime/card/tasks" approved yes
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
cc_transition_plan_status "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks" approved confirmed
assert_eq "$(cc_maintainer_approval_commit_required "$product" "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks")" MAINTAINER_APPROVAL_COMMIT_REQUIRED
product_card="$ROOT/test/approval/fixtures/product-source-follow-on.txt"
wrapped_card="$ROOT/test/approval/fixtures/wrapped-follow-on.txt"
fixture_lines_in "$product_card" "$ROOT/docs/gates.md"
contains "$product_card" 'Action: commit-approved-plan'
contains "$product_card" 'Confirm commit of the approved plan state.'
not_contains "$product_card" 'Next action: Run approved plan'
contains "$ROOT/agents/coordinator.md" 'Do not name `Run approved plan <id>` as the'
fixture_lines_in "$wrapped_card" "$ROOT/agents/coordinator.md"
contains "$wrapped_card" 'Next action: Run approved plan <id>'
not_contains "$wrapped_card" 'Action: commit-approved-plan'
contains "$ROOT/agents/coordinator.md" 'cc_transition_plan_status'
contains "$ROOT/agents/coordinator.md" 'cc_maintainer_approval_commit_required'
contains "$ROOT/docs/planning.md" 'next gate is in'
printf '%s\n' "$(cc_route 'Confirm approval of plan demo.')" | grep -F 'capability: execute-plan' >/dev/null && fail 'confirm approval routed to execution'
contains "$product/plans/app-plans/demo/plan.yaml" 'status: approved'
contains "$product/plans/app-plans/demo/tasks/one.md" 'status: ready'
contains "$product/plans/app-plans/demo/tasks/one.md" 'Keep this body status: draft example.'
assert_eq "$(git -C "$product" rev-parse HEAD)" "$head_before"
test ! -d "$product/.runtime/worktrees" || fail 'status transition created a worktree'
test ! -d "$product/.runtime/plans" || fail 'status transition created a lease'
expect_failure cc_prepare_worktree "$product" "$runtime/worktrees/product-blocked" main
printf '%s\n' unrelated > "$product/unrelated.txt"
expect_failure cc_maintainer_approval_commit_required "$product" "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks"
expect_failure cc_transition_plan_status "$product/plans/app-plans/demo/plan.yaml" "$product/plans/app-plans/demo/tasks" approved confirmed

dirty="$runtime/product-dirty"
init_repo "$dirty" product-source
printf '%s\n' extra > "$dirty/extra.txt"
cc_transition_plan_status "$dirty/plans/app-plans/demo/plan.yaml" "$dirty/plans/app-plans/demo/tasks" approved confirmed
expect_failure cc_maintainer_approval_commit_required "$dirty" "$dirty/plans/app-plans/demo/plan.yaml" "$dirty/plans/app-plans/demo/tasks"
contains "$ROOT/docs/gates.md" 'DIRTY_BASE_BLOCKED'
contains "$dirty/plans/app-plans/demo/plan.yaml" 'status: approved'

wrapped="$runtime/wrapped"
init_repo "$wrapped" instantiated-workspace
wrapped_head=$(git -C "$wrapped" rev-parse HEAD)
cc_transition_plan_status "$wrapped/plans/app-plans/demo/plan.yaml" "$wrapped/plans/app-plans/demo/tasks" approved confirmed
expect_failure cc_maintainer_approval_commit_required "$wrapped" "$wrapped/plans/app-plans/demo/plan.yaml" "$wrapped/plans/app-plans/demo/tasks"
fixture_lines_in "$wrapped_card" "$ROOT/agents/coordinator.md"
contains "$ROOT/docs/gates.md" 'Instantiated or wrapped workspaces do not receive this card'
not_contains "$wrapped_card" 'Action: commit-approved-plan'
assert_eq "$(git -C "$wrapped" rev-parse HEAD)" "$wrapped_head"
pass 'two-turn approval is mutation-free, status-only, and follow-on cards stay outside the engine'
