#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

# The template-runtime laboratory assembles the distributable
# context-circuit-template, initializes an isolated project workspace from it,
# connects fixture repositories, and exercises the product experience through
# the artifact's own runtime boundary. It must not import the source
# repository's Product Knowledge, plans, .runtime, or implementation state.

stage=$(mktemp -d "${TMPDIR:-/tmp}/cc-tr-stage.XXXXXX")
out=$(mktemp -d "${TMPDIR:-/tmp}/cc-tr-out.XXXXXX")
lab=$(mktemp -d "${TMPDIR:-/tmp}/cc-tr-lab.XXXXXX")
trap 'rm -rf "$stage" "$out" "$lab"' EXIT HUP INT TERM

# 1. Assemble the template and instantiate an isolated workspace from it.
sh "$ROOT/scripts/release-artifact.sh" "$stage" "$out" v0.5.0 >/dev/null
ws="$lab/project"
cp -R "$out/context-circuit-v0.5.0" "$ws"

# 2. Use the ARTIFACT's shipped engine, not the source engine under test.
ENGINE="$ws/wrapper/runtime/engine.sh"
require_file "$ENGINE"
# shellcheck disable=SC1090
. "$ENGINE"

# 3. The instantiated workspace carries only the blank seed (no source state).
not_contains "$ws/plans/INDEX.md" "0001-"
test ! -e "$ws/plans/context-circuit-plans" || fail 'source maintainer plans leaked into workspace'
test ! -e "$ws/.runtime/executions" || fail 'source runtime state leaked into workspace'
contains "$ws/workspace.yaml" 'uninitialized-workspace'

# 4. Initialize and confirm identity.
printf 'schema_version: 1\nworkspace: lab-commerce\ntitle: Lab Commerce\nrepositories: []\n' >"$ws/workspace.yaml"
cc_workspace_init "$ws" >/dev/null
cc_workspace_validate "$ws" >/dev/null

# 5. Connect two fixture repositories with team anchor branches.
for r in api web; do
	mkdir -p "$ws/repositories/$r/src"
	git -C "$ws/repositories/$r" init -q -b development
	git -C "$ws/repositories/$r" config user.email t@t.t
	git -C "$ws/repositories/$r" config user.name t
	printf 'seed\n' >"$ws/repositories/$r/src/seed.txt"
	git -C "$ws/repositories/$r" add -A
	git -C "$ws/repositories/$r" commit -q -m seed
done
printf 'schema_version: 1\nbindings:\n  api:\n    path: repositories/api\n    anchor_branch: development\n  web:\n    path: repositories/web\n    anchor_branch: development\n' >"$ws/repositories.local.yaml"

# 6. Create, approve, and execute a multi-repository plan.
pid=$(cc_plan_allocate_id "$ws" checkout-v2)
assert_eq "0001-checkout-v2" "$pid"
mkdir -p "$ws/plans/$pid/tasks"
cat >"$ws/plans/$pid/plan.yaml" <<EOF
schema_version: 1
plan: $pid
title: Checkout v2
status: draft
objective: Add saved checkout sessions.
repositories:
  - id: api
  - id: web
product_knowledge:
  - id: project.core
    path: context/PROJECT.md
    reason: Grounds the objective.
context_grounding:
  summary: API owns state; web renders it.
  constraints: []
  decisions: []
knowledge_impact:
  expected_context_units: []
  review_on_completion: true
tasks:
  - id: API-001
    title: API
    repositories: [api]
    paths: [src]
    depends_on: []
    changes: [Add API.]
    acceptance:
      - id: API-AC
        statement: API works.
    verification:
      - id: API-VT
        command: test/api.sh
  - id: WEB-001
    title: Web
    repositories: [web]
    paths: [src]
    depends_on: [API-001]
    changes: [Add UI.]
    acceptance:
      - id: WEB-AC
        statement: UI works.
    verification:
      - id: WEB-VT
        command: test/web.sh
execution:
  worker: one
  independent_verifier: required
  max_worker_failures: 3
EOF
printf '# Checkout v2\n' >"$ws/plans/$pid/PLAN.md"
cc_plan_validate "$ws/plans/$pid" >/dev/null
cc_plan_index_upsert "$ws" "$pid" >/dev/null
cc_plan_approve "$ws" "$pid" >/dev/null
exec=$(cc_execution_begin "$ws" "$pid" lab-session | sed -n 's/^execution_id: //p')
edir="$ws/.runtime/executions/$pid/$exec"

# 7. One worker commits both repositories; independent verifier passes.
cc_attempt_begin "$edir" >/dev/null
for r in api web; do
	printf 'work\n' >>"$ws/.runtime/worktrees/$pid/$r/src/change.txt"
	git -C "$ws/.runtime/worktrees/$pid/$r" add -A
	git -C "$ws/.runtime/worktrees/$pid/$r" commit -q -m impl
	cc_worker_commit_record "$edir" "$r" implementation >/dev/null
done
cc_verifier_prepare "$edir" >/dev/null
cc_verifier_result_record "$edir" 1 passed >/dev/null
assert_eq "verified" "$(cc_execution_status "$edir")"

# 8. Human completion, then archive and restore.
cc_plan_complete "$ws" "$pid" >/dev/null
assert_eq "done" "$(cc_plan_status "$ws" "$pid")"
cc_plan_archive "$ws" "$pid" >/dev/null
require_dir "$ws/plans/.archived/$pid"
cc_plan_restore "$ws" "$pid" >/dev/null
require_dir "$ws/plans/$pid"

# 9. The lab workspace's runtime state is separate from the source repository's.
test "$ws/.runtime" != "$ROOT/.runtime" || fail 'lab runtime collided with source runtime'
require_dir "$ws/.runtime/executions/$pid"

pass 'template-runtime laboratory'
