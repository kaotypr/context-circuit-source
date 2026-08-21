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

contains "$ROOT/plans/README.md" 'Approval never starts execution'
contains "$ROOT/docs/plan-review.md" 'read-only'
contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'schema v1 remains readable'
pass 'separate approval/run/finish lifecycle and task projection'
