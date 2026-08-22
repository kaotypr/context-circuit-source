#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-lifecycle.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
mkdir -p "$runtime/tasks"
cat > "$runtime/plan.yaml" <<'EOF'
schema_version: 2
id: lifecycle-fixture
status: draft
tasks: [TASK-01]
EOF
cat > "$runtime/tasks/TASK-01.md" <<'EOF'
---
id: TASK-01
status: draft
---
# Task
EOF
expect_failure cc_transition_plan_status "$runtime/plan.yaml" "$runtime/tasks" approved
contains "$runtime/plan.yaml" 'status: draft'
cc_transition_plan_status "$runtime/plan.yaml" "$runtime/tasks" approved confirmed
contains "$runtime/plan.yaml" 'status: approved'
contains "$runtime/tasks/TASK-01.md" 'status: ready'
cc_transition_plan_status "$runtime/plan.yaml" "$runtime/tasks" done confirmed
contains "$runtime/plan.yaml" 'status: done'
contains "$runtime/tasks/TASK-01.md" 'status: done'
expect_failure cc_transition_plan_status "$runtime/plan.yaml" "$runtime/tasks" approved confirmed
contains "$ROOT/template/plans/README.md" 'Approval never starts execution'
contains "$ROOT/docs/plan-review.md" 'read-only'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'schema v1 remains readable'

preserve="$runtime/preserve"
mkdir -p "$preserve/tasks"
cat > "$preserve/plan.yaml" <<'EOF'
schema_version: 2
id: preserve-fixture
status: draft
tasks: [TASK-01]
EOF
printf '%s' '---
id: TASK-01
status: draft
---
Body must keep this status: draft line.

# Task' > "$preserve/tasks/TASK-01.md"
cp "$preserve/tasks/TASK-01.md" "$preserve/tasks/TASK-01.md.before"
expect_failure cc_transition_plan_status "$preserve/plan.yaml" "$preserve/tasks" approved stale
expect_failure cc_transition_plan_status "$preserve/plan.yaml" "$preserve/tasks" approved yes
cmp -s "$preserve/tasks/TASK-01.md" "$preserve/tasks/TASK-01.md.before" || fail 'invalid confirmation mutated a task body'
cc_transition_plan_status "$preserve/plan.yaml" "$preserve/tasks" approved confirmed
contains "$preserve/plan.yaml" 'status: approved'
contains "$preserve/tasks/TASK-01.md" 'status: ready'
contains "$preserve/tasks/TASK-01.md" 'Body must keep this status: draft line.'
test "$(tail -c 1 "$preserve/tasks/TASK-01.md" | od -An -tx1 | tr -d ' \n')" != 0a || fail 'task file gained a trailing newline'
test "$(grep -c '^status: ready$' "$preserve/tasks/TASK-01.md")" -eq 1 || fail 'status projection replaced more than the first status line'

pass 'separate approval/run/finish lifecycle and task projection'
