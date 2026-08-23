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

binding_root="$runtime/binding-workspace"
mkdir -p "$binding_root"
cat > "$binding_root/workspace.yaml" <<'EOF'
version: 1
template_version: 1.0.0
workspace:
  name: binding-fixture
repositories:
  app:
    canonical_url: https://github.com/acme/app.git
    default_branch: main
EOF
bound_source="$runtime/external-app"
git init -q -b main "$bound_source"
git -C "$bound_source" config user.email test@example.invalid
git -C "$bound_source" config user.name 'Context Circuit Test'
printf '%s\n' bound > "$bound_source/file.txt"
git -C "$bound_source" add file.txt
git -C "$bound_source" commit -qm initial
git -C "$bound_source" remote add origin https://github.com/acme/app.git
cat > "$binding_root/repositories.local.yaml" <<EOF
repositories:
  app:
    path: $bound_source
    remote: git@github.com:acme/app.git
EOF
assert_eq "$(cc_validate_repository_binding "$binding_root" app)" "$bound_source"
mkdir -p "$binding_root/projects"
mv "$bound_source" "$binding_root/projects/app"
sed -i 's#path: .*#path: projects/app#' "$binding_root/repositories.local.yaml"
assert_eq "$(cc_validate_repository_binding "$binding_root" app)" "$binding_root/projects/app"
mkdir -p "$binding_root/repositories"
mv "$binding_root/projects/app" "$binding_root/repositories/app"
sed -i 's#path: .*#path: repositories/app#' "$binding_root/repositories.local.yaml"
assert_eq "$(cc_validate_repository_binding "$binding_root" app)" "$binding_root/repositories/app"
printf '%s\n' dirty >> "$binding_root/repositories/app/file.txt"
expect_failure cc_validate_repository_binding "$binding_root" app
git -C "$binding_root/repositories/app" checkout -- file.txt
prepared=$(cc_prepare_bound_worktree "$binding_root" app repository-bootstrap main)
assert_eq "$prepared" "$binding_root/.runtime/worktrees/app/repository-bootstrap"
test -e "$prepared/.git" || fail 'bound repository worktree was not prepared'
test -z "$(git -C "$binding_root/repositories/app" status --porcelain --untracked-files=all)" || fail 'bound source repository was modified'

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
rm -f "$maintainer/unrelated.txt"
git -C "$maintainer" checkout -- plans
cc_transition_plan_status "$maintainer/plans/context-circuit-plans/demo/plan.yaml" "$maintainer/plans/context-circuit-plans/demo/tasks" approved confirmed
assert_eq "$(cc_maintainer_approval_commit_required "$maintainer" "$maintainer/plans/context-circuit-plans/demo/plan.yaml" "$maintainer/plans/context-circuit-plans/demo/tasks")" MAINTAINER_APPROVAL_COMMIT_REQUIRED
contains "$ROOT/docs/gates.md" 'present this existing commit card in the same session immediately'
contains "$ROOT/agents/coordinator.md" 'Do not add engine helpers for card text'
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
cat > "$runtime/expanded-launch.yaml" <<'EOF'
role: verifier
assigned_root: .runtime/worktrees/app/checkout
delegation_locator: delegation.yaml
handoff_locator: handoff.md
prompt: expand scope
EOF
expect_failure cc_validate_launch_projection "$runtime/expanded-launch.yaml" "$runtime/verifier.yaml"
for host in codex claude-code cursor-agent; do
  fixture="$ROOT/test/hosts/fixtures/$host.yaml"
  require_file "$fixture"
  contains "$fixture" 'child_mapping:'
  contains "$fixture" 'worktree: exclusive'
  contains "$fixture" 'write_worktree: true'
  contains "$fixture" 'worktree: independent'
  contains "$fixture" 'write_worktree: false'
  contains "$fixture" 'write_plan: false'
  contains "$fixture" 'write_activity: false'
done
pass 'atomic single-plan lease, exclusive scope, and read-only verifier permissions'
