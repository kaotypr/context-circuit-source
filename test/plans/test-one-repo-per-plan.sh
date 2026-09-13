#!/bin/sh
# One plan, one repository, one worker. A two-repository intent cannot collapse.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

inv="$ROOT/.context-circuit/wrapper/contracts/invariants.yaml"
pe="$ROOT/context/domains/plan-execution/README.md"
arch="$ROOT/context/ARCHITECTURE.md"
tr="$ROOT/product/.agents/skills/cc-trace/SKILL.md"
cp="$ROOT/product/.agents/skills/cc-plan/SKILL.md"
tracing="$ROOT/context/domains/tracing/README.md"

contains "$inv" "INV-EXEC-02"
contains "$inv" "that plan's single repository"
contains "$pe" "one repository"
contains "$pe" "one worker"
contains "$arch" "one repository"
contains "$tr" "one repository"
contains "$tr" "plan_dependencies"
contains "$tr" "stacked"
contains "$cp" "one repository"
contains "$cp" "plan_dependencies"
contains "$cp" "Refuse collapsing"
contains "$tracing" "one repository"
contains "$tracing" "plan_dependencies"
not_contains "$inv" "across all mapped repositories"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

# one intent, two stacked one-repo plans remain authorized with no second gate
iid=i001-split
mkdir -p "$ws/intent/$iid"
cat >"$ws/intent/$iid/contract.yaml" <<EOF
schema_version: 2
intent: $iid
title: Split
goal: Split goal.
non_goals:
  - none
constraints:
  - none
acceptance_criteria:
  - id: ac-1
    statement: Split works.
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
printf '# Split\n' >"$ws/intent/$iid/INTENT.md"
cc_intent_index_upsert "$ws" "$iid" >/dev/null
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid" >/dev/null

cc_fx_plan_intent "$ws" 0001-api "API half" api src "$iid"
cc_fx_plan_intent "$ws" 0002-web "Web half" web src "$iid"
# add an ordering edge
awk '
	/^product_knowledge:/ && !added {
		print "plan_dependencies:"
		print "  - id: 0001-api"
		print "    reason: web consumes the API change."
		added=1
	}
	{print}
' "$ws/plans/0002-web/plan.yaml" >"$ws/plans/0002-web/plan.yaml.new"
mv "$ws/plans/0002-web/plan.yaml.new" "$ws/plans/0002-web/plan.yaml"
cc_plan_validate "$ws/plans/0001-api" >/dev/null
cc_plan_validate "$ws/plans/0002-web" >/dev/null
cc_intent_authorized "$ws" 0001-api >/dev/null
cc_intent_authorized "$ws" 0002-web >/dev/null

# collapsing both repositories into one plan fails
mkdir -p "$ws/plans/0003-collapsed/tasks"
cat >"$ws/plans/0003-collapsed/plan.yaml" <<EOF
schema_version: 3
plan: 0003-collapsed
intent: $iid
title: Collapsed
status: draft
objective: invalid collapse
repositories:
  - id: api
  - id: web
product_knowledge:
  - id: project.core
    path: context/PROJECT.md
    reason: Grounds the objective.
context_grounding:
  summary: Invalid.
  constraints: []
  decisions: []
knowledge_impact:
  expected_context_units: []
  review_on_completion: true
tasks:
  - id: T-001
    title: both
    repositories: [api, web]
    paths: [src]
    depends_on: []
    changes: [Both.]
    acceptance:
      - id: T-AC
        statement: both changed
    verification:
      - id: T-VT
        command: true
execution:
  worker: one
  independent_verifier: required
  max_worker_failures: 3
EOF
printf '# Collapsed\n' >"$ws/plans/0003-collapsed/PLAN.md"
expect_failure cc_plan_validate "$ws/plans/0003-collapsed"

pass 'one repo per plan'
