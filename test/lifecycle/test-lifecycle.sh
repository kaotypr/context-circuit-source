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
contains "$ROOT/docs/plan-review.md" 'Human decisions'
contains "$ROOT/docs/plan-review.md" 'Risks/assumptions'
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

partial="$runtime/partial"
mkdir -p "$partial/tasks"
cat > "$partial/plan.yaml" <<'EOF'
schema_version: 2
id: partial-fixture
status: draft
tasks: [TASK-01, TASK-02]
EOF
cat > "$partial/tasks/TASK-01.md" <<'EOF'
---
id: TASK-01
status: draft
---
# Task
EOF
cat > "$partial/tasks/TASK-02.md" <<'EOF'
---
id: TASK-02
Status: draft
---
# Task
EOF
expect_failure cc_transition_plan_status "$partial/plan.yaml" "$partial/tasks" approved confirmed
contains "$partial/plan.yaml" 'status: draft'
contains "$partial/tasks/TASK-01.md" 'status: draft'
contains "$partial/tasks/TASK-02.md" 'Status: draft'
not_contains "$partial/plan.yaml" 'status: approved'
not_contains "$partial/tasks/TASK-01.md" 'status: ready'

named_review=$(cc_route 'Review plan checkout-validation.')
printf '%s\n' "$named_review" | grep -F 'capability: review-plan' >/dev/null ||
  fail 'named-plan review did not route to review-plan'
printf '%s\n' "$named_review" | grep -F 'authorization: read-only' >/dev/null ||
  fail 'named-plan review was not read-only'
printf '%s\n' "$named_review" | grep -F 'context_set: plan-review' >/dev/null ||
  fail 'named-plan review did not load the plan-review context set'

named_review_variant=$(cc_route 'Walk me through plan checkout-validation.')
printf '%s\n' "$named_review_variant" | grep -F 'capability: review-plan' >/dev/null ||
  fail 'equivalent named-plan review phrasing did not route to review-plan'

unnamed_review=$(cc_route 'Is this plan approved-for-execution?')
printf '%s\n' "$unnamed_review" | grep -F 'capability: clarify-target' >/dev/null ||
  fail 'unnamed review request was not routed to clarify-target'
printf '%s\n' "$unnamed_review" | grep -F 'capability: review-plan' >/dev/null &&
  fail 'unnamed review request incorrectly inspected an unrelated plan' || :

pass 'separate approval/run/finish lifecycle, read-only named-plan review, and task projection'
