#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-ownership.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
cc_acquire_lease "$runtime" checkout sess-a sess-a app .runtime/worktrees/app/checkout main
contains "$runtime/plans/checkout/lease.lock/owner.yaml" 'session_id: sess-a'
expect_failure cc_acquire_lease "$runtime" checkout sess-b sess-b app .runtime/worktrees/app/checkout-b main
cc_scope_allows .runtime/worktrees/app/checkout/src/file .runtime/worktrees/app/checkout
expect_failure cc_scope_allows .runtime/worktrees/other/file .runtime/worktrees/app/checkout

repo="$runtime/repository"
git init -q -b main "$repo"
git -C "$repo" config user.email test@example.invalid
git -C "$repo" config user.name 'Context Circuit Test'
printf '%s\n' clean > "$repo/file.txt"
git -C "$repo" add file.txt
git -C "$repo" commit -qm initial
cc_prepare_worktree "$repo" "$runtime/worktrees/app/checkout" main
test -e "$runtime/worktrees/app/checkout/.git" || fail 'exclusive Git worktree was not created'
printf '%s\n' dirty > "$repo/dirty.txt"
expect_failure cc_prepare_worktree "$repo" "$runtime/worktrees/app/dirty" main

maintainer="$runtime/maintainer"
mkdir -p "$maintainer/plans/context-circuit-plans/demo/tasks"
cat > "$maintainer/workspace.yaml" <<'EOF'
version: 1
workspace:
  name: context-circuit-source
identity:
  kind: product-source
EOF
cat > "$maintainer/plans/context-circuit-plans/demo/plan.yaml" <<'EOF'
schema_version: 2
id: demo
status: draft
EOF
for task_id in one two; do
  cat > "$maintainer/plans/context-circuit-plans/demo/tasks/$task_id.md" <<'EOF'
---
schema_version: 2
status: draft
---
# task
EOF
done
git init -q -b main "$maintainer"
git -C "$maintainer" config user.email test@example.invalid
git -C "$maintainer" config user.name 'Context Circuit Test'
git -C "$maintainer" add workspace.yaml plans
git -C "$maintainer" commit -qm initial
sed -i 's/^status: draft$/status: approved/' "$maintainer/plans/context-circuit-plans/demo/plan.yaml"
sed -i 's/^status: draft$/status: ready/' "$maintainer/plans/context-circuit-plans/demo/tasks/one.md"
sed -i 's/^status: draft$/status: ready/' "$maintainer/plans/context-circuit-plans/demo/tasks/two.md"
assert_eq "$(cc_maintainer_approval_commit_required "$maintainer" "$maintainer/plans/context-circuit-plans/demo/plan.yaml" "$maintainer/plans/context-circuit-plans/demo/tasks")" MAINTAINER_APPROVAL_COMMIT_REQUIRED
expect_failure cc_prepare_worktree "$maintainer" "$runtime/worktrees/maintainer/blocked" main
printf '%s\n' unrelated > "$maintainer/unrelated.txt"
expect_failure cc_maintainer_approval_commit_required "$maintainer" "$maintainer/plans/context-circuit-plans/demo/plan.yaml" "$maintainer/plans/context-circuit-plans/demo/tasks"

cat > "$runtime/verifier.yaml" <<'EOF'
schema_version: 1
wrapper_version: 1.0.0
session_id: verifier
parent_session_id: root
root_session_id: root
role: verifier
objective: Verify
scope: bounded
non_goals: []
context_set: verifier
context_receipt: receipt
repository: app
worktree: .runtime/worktrees/app/checkout
permissions:
  write_worktree: false
  write_plan: false
  write_activity: false
acceptance_criteria: [AC-01]
verification: [VT-01]
expected_evidence: [handoff]
stop_conditions: [scope]
handoff_schema: v1
EOF
cc_validate_delegation "$runtime/verifier.yaml"
pass 'atomic single-plan lease, exclusive scope, and read-only verifier permissions'
