#!/bin/sh
set -eu
. "$(git -C "$(dirname -- "$0")" rev-parse --show-toplevel)/test/lib/assert.sh"

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
# v1.0: author + approve a parent intent scoped to the plan's repos before the plan.
mkdir -p "$ws/intent/i$pid"
cat >"$ws/intent/i$pid/contract.yaml" <<EOF
schema_version: 1
intent: i$pid
title: Checkout v2
goal: Add saved checkout sessions.
non_goals:
  - none
constraints:
  - none
acceptance_criteria:
  - id: ac-1
    statement: Saved checkout works.
    method: test
done_when: ac-1 passes and a human accepts the candidate.
scope:
  repositories:
    - id: api
      paths: [src]
    - id: web
      paths: [src]
tier: standard
status: draft
contract_digest:
EOF
printf '# Checkout v2\n' >"$ws/intent/i$pid/INTENT.md"
cc_intent_index_upsert "$ws" "i$pid" >/dev/null
cc_intent_approve "$ws" "i$pid" >/dev/null
mkdir -p "$ws/plans/$pid/tasks"
cat >"$ws/plans/$pid/plan.yaml" <<EOF
schema_version: 3
plan: $pid
intent: i$pid
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
require_dir "$ws/plans/archive/$pid"
cc_plan_restore "$ws" "$pid" >/dev/null
require_dir "$ws/plans/$pid"

# 9. The lab workspace's runtime state is separate from the source repository's.
test "$ws/.runtime" != "$ROOT/.runtime" || fail 'lab runtime collided with source runtime'
require_dir "$ws/.runtime/executions/$pid"

# 10. Run-stack (v0.6): the SHIPPED template carries base selection + path leases.
#     Two roots and one integration dependent, all in api; plus one lease check.
mkplan() { # pid title repo path "deps"
	mp_dir="$ws/plans/$1"; mkdir -p "$mp_dir/tasks"
	mp_iid="i$1"; mkdir -p "$ws/intent/$mp_iid"
	{
		printf 'schema_version: 1\nintent: %s\ntitle: %s\ngoal: %s goal.\n' "$mp_iid" "$2" "$2"
		printf 'non_goals:\n  - none\nconstraints:\n  - none\n'
		printf 'acceptance_criteria:\n  - id: ac-1\n    statement: %s works.\n    method: test\n' "$2"
		printf 'done_when: ac-1 passes and a human accepts the candidate.\n'
		printf 'scope:\n  repositories:\n    - id: %s\n      paths: [%s]\n' "$3" "$4"
		printf 'tier: standard\nstatus: draft\ncontract_digest:\n'
	} >"$ws/intent/$mp_iid/contract.yaml"
	printf '# %s\n' "$2" >"$ws/intent/$mp_iid/INTENT.md"
	cc_intent_index_upsert "$ws" "$mp_iid" >/dev/null
	cc_intent_approve "$ws" "$mp_iid" >/dev/null
	{
		printf 'schema_version: 3\nplan: %s\ntitle: %s\nstatus: draft\nobjective: %s objective.\nintent: %s\n' "$1" "$2" "$2" "$mp_iid"
		printf 'repositories:\n  - id: %s\n' "$3"
		if [ -n "$5" ]; then printf 'plan_dependencies:\n'; for d in $5; do printf '  - id: %s\n    reason: builds on %s\n' "$d" "$d"; done; fi
		printf 'product_knowledge:\n  - id: project.core\n    path: context/PROJECT.md\n    reason: Grounds it.\n'
		printf 'context_grounding:\n  summary: s\n  constraints: []\n  decisions: []\n'
		printf 'knowledge_impact:\n  expected_context_units: []\n  review_on_completion: true\n'
		printf 'tasks:\n  - id: T-001\n    title: Work\n    repositories: [%s]\n    paths: [%s]\n    depends_on: []\n' "$3" "$4"
		printf '    changes: [Create %s/mod.txt.]\n    acceptance:\n      - id: T-AC\n        statement: %s/mod.txt exists.\n    verification:\n      - id: T-VT\n        command: test -f %s/mod.txt\n' "$4" "$4" "$4"
		printf 'execution:\n  worker: one\n  independent_verifier: required\n  max_worker_failures: 3\n'
	} >"$mp_dir/plan.yaml"
	printf '# %s\n' "$2" >"$mp_dir/PLAN.md"
	cc_plan_index_upsert "$ws" "$1" >/dev/null
}
runok() { # pid repo path
	cc_plan_approve "$ws" "$1" >/dev/null
	ro_ex=$(cc_execution_begin "$ws" "$1" "$1-w" | sed -n 's/^execution_id: //p')
	ro_ed="$ws/.runtime/executions/$1/$ro_ex"; ro_wt="$ws/.runtime/worktrees/$1/$2"
	cc_attempt_begin "$ro_ed" >/dev/null
	mkdir -p "$ro_wt/$3"; printf 'm\n' >"$ro_wt/$3/mod.txt"
	git -C "$ro_wt" add -A; git -C "$ro_wt" commit -q -m "feat($2): $3"
	cc_worker_commit_record "$ro_ed" "$2" implementation >/dev/null
	cc_verifier_prepare "$ro_ed" >/dev/null
	cc_verifier_result_record "$ro_ed" 1 passed >/dev/null
}
mkplan 0002-sa  "Stack A"   api src/a ""
mkplan 0003-sb  "Stack B"   api src/b ""
mkplan 0004-int "Integrate" api src/c "0002-sa 0003-sb"
runok 0002-sa api src/a
runok 0003-sb api src/b
cc_plan_approve "$ws" 0004-int >/dev/null
bp=$(cc_base_prepare "$ws" 0004-int api)
printf '%s' "$bp" | grep -q 'based_on: \[0002-sa, 0003-sb\]' || fail 'shipped integration base is missing based_on'
intbase=$(printf '%s' "$bp" | sed -n 's/^base_commit: //p')
a_tip=$(git -C "$ws/repositories/api" rev-parse --verify refs/heads/cc/0002-sa/api)
b_tip=$(git -C "$ws/repositories/api" rev-parse --verify refs/heads/cc/0003-sb/api)
git -C "$ws/repositories/api" merge-base --is-ancestor "$a_tip" "$intbase" || fail 'shipped integration base not built on predecessor A'
git -C "$ws/repositories/api" merge-base --is-ancestor "$b_tip" "$intbase" || fail 'shipped integration base not built on predecessor B'
# one lease: a held region blocks an unrelated plan but exempts a declared descendant
cc_lease_acquire "$ws" api 0002-sa "src/a" >/dev/null
expect_failure cc_lease_check "$ws" api 0003-sb "src/a"
cc_lease_check "$ws" api 0004-int "src/a" >/dev/null || fail 'shipped lease must exempt a descendant'

# 11. Repository grounding (v0.6): the SHIPPED template discovers the repo's own
#     agent guidance and assembles a worker brief with the required grounding slot.
require_file "$ws/wrapper/runtime/worker-brief.md"    # runtime-only brief, promoted beside the runtime
git -C "$ws/repositories/api" checkout -q development
printf '# API agent guide\n\nStart every new source file with `// @grounded`.\n' >"$ws/repositories/api/AGENTS.md"
git -C "$ws/repositories/api" add -A; git -C "$ws/repositories/api" commit -q -m 'chore(api): add agent guidance'
mkplan 0005-grounded "Grounded" api src/g ""
cc_plan_approve "$ws" 0005-grounded >/dev/null
gex=$(cc_execution_begin "$ws" 0005-grounded sess-g | sed -n 's/^execution_id: //p')
gedir="$ws/.runtime/executions/0005-grounded/$gex"
require_file "$gedir/grounding/api.yaml"
contains "$gedir/grounding/api.yaml" "- AGENTS.md"
cc_worker_brief_assemble "$ws" "$gedir" api "Add the grounded module." >/dev/null
require_file "$gedir/brief-api.md"
contains "$gedir/brief-api.md" "## Repository grounding"
contains "$gedir/brief-api.md" "AGENTS.md — read and honor it"
cc_brief_preflight "$gedir/brief-api.md" >/dev/null

# 12. Execution latency: the SHIPPED template records bounded per-attempt
#     host evidence (per-role model/effort), refuses an unknown evidence key,
#     validates the optional plan complexity hint, and ships the per-role tiering
#     guidance and the 0.7.0 runtime version.
contains "$ws/wrapper/manifest.yaml" "runtime_version: 1.0.0"
require_file "$ws/docs/role-tiering.md"                  # per-role tiering guidance ships (docs/)
mkplan 0006-latency "Latency" api src/lat ""
cc_plan_approve "$ws" 0006-latency >/dev/null
lex=$(cc_execution_begin "$ws" 0006-latency sess-l | sed -n 's/^execution_id: //p')
ledir="$ws/.runtime/executions/0006-latency/$lex"
cc_attempt_begin "$ledir" >/dev/null
mkdir -p "$ws/.runtime/worktrees/0006-latency/api/src/lat"
printf 'm\n' >"$ws/.runtime/worktrees/0006-latency/api/src/lat/mod.txt"
git -C "$ws/.runtime/worktrees/0006-latency/api" add -A
git -C "$ws/.runtime/worktrees/0006-latency/api" commit -q -m "feat(api): lat"
cc_worker_commit_record "$ledir" api implementation >/dev/null
cc_verifier_prepare "$ledir" >/dev/null
# the coordinator records bounded, credential-free per-attempt host evidence:
# the (model, effort) each role ran at (no wall-clock — timing is engine-stamped)
cc_attempt_evidence_record "$ledir" 1 worker_model=model-hi worker_effort=high verifier_model=model-lo verifier_effort=medium >/dev/null
contains "$ledir/attempts/001/host-evidence.yaml" "worker_model: model-hi"
contains "$ledir/attempts/001/host-evidence.yaml" "verifier_model: model-lo"
expect_failure cc_attempt_evidence_record "$ledir" 1 bogus_key=1   # bounded surface: unknown key refused
cc_verifier_result_record "$ledir" 1 passed >/dev/null
# the optional per-plan complexity hint validates in the shipped engine
printf 'complexity: high\n' >>"$ws/plans/0006-latency/plan.yaml"
cc_plan_validate "$ws/plans/0006-latency" >/dev/null || fail 'shipped engine rejects a valid complexity hint'

pass 'template-runtime laboratory'
