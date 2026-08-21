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
