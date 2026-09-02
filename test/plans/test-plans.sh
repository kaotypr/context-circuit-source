#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development

# --- stable, sequential, human-friendly plan ids ---
id1=$(cc_plan_allocate_id "$ws" billing-v2)
assert_eq "0001-billing-v2" "$id1"
cc_fx_plan "$ws" "$id1" "Billing v2" "api"
id2=$(cc_plan_allocate_id "$ws" saved-checkout)
assert_eq "0002-saved-checkout" "$id2"
cc_fx_plan "$ws" "$id2" "Saved checkout" "api"

# --- sequence never reused after archive ---
cc_plan_archive "$ws" "$id2" >/dev/null
id3=$(cc_plan_allocate_id "$ws" checkout-v2)
assert_eq "0003-checkout-v2" "$id3"
cc_plan_restore "$ws" "$id2" >/dev/null

# --- bad slugs and ids are rejected ---
expect_failure cc_plan_allocate_id "$ws" "Bad_Slug"
expect_failure cc_plan_allocate_id "$ws" "with spaces"
cc_plan_id_valid "0001-billing-v2" || fail "valid id rejected"
expect_failure cc_plan_id_valid "1-x"
expect_failure cc_plan_id_valid "0001-Bad"

# --- validation catches missing tasks ---
mkdir -p "$ws/plans/0009-empty/tasks"
printf 'schema_version: 3\nplan: 0009-empty\ntitle: E\nstatus: draft\nobjective: e\nintent: i0009-empty\nrepositories:\n  - id: api\ntasks:\n' >"$ws/plans/0009-empty/plan.yaml"
printf '# E\n' >"$ws/plans/0009-empty/PLAN.md"
expect_failure cc_plan_validate "$ws/plans/0009-empty"
rm -rf "$ws/plans/0009-empty"

# --- validation (review) never changes plan status ---
cc_plan_validate "$ws/plans/$id1" >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" "$id1")"

# --- active index maintenance ---
contains "$ws/plans/INDEX.md" "| $id1 |"
cc_plan_index_remove "$ws" "$id1"
not_contains "$ws/plans/INDEX.md" "| $id1 |"
cc_plan_index_upsert "$ws" "$id1" >/dev/null
contains "$ws/plans/INDEX.md" "| $id1 |"

# --- authorization: a plan is authorized by its approved intent within the scope
#     envelope (INV-INTENT-02). There is no separate plan-approval status — the plan
#     stays `draft`, and the envelope check is the preflight that execution re-runs. ---
cc_intent_envelope_check "$ws" "$id1" >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" "$id1")"   # no intermediate "approved" status

# --- a plan authored with block-list task fields also validates ---
mkdir -p "$ws/plans/0011-block/tasks"
cat >"$ws/plans/0011-block/plan.yaml" <<'EOF'
schema_version: 3
plan: 0011-block
title: Block form
status: draft
objective: block
intent: i0011-block
repositories:
  - id: api
  - id: web
tasks:
  - id: T-001
    title: t1
    repositories:
      - api
    paths:
      - src
    depends_on: []
  - id: T-002
    title: t2
    repositories:
      - web
    paths:
      - app
    depends_on:
      - T-001
EOF
printf '# Block\n' >"$ws/plans/0011-block/PLAN.md"
cc_plan_validate "$ws/plans/0011-block" >/dev/null
assert_eq "api
web" "$(cc_plan_affected_repositories "$ws/plans/0011-block/plan.yaml")"

# --- execute within the intent envelope (no separate plan-approval step) ---
cc_fx_plan "$ws" 0010-compound "Compound" "api"
cc_execution_begin "$ws" 0010-compound s2 >/dev/null
require_dir "$ws/.runtime/executions/0010-compound"
assert_eq "draft" "$(cc_plan_status "$ws" 0010-compound)"   # stays draft through execution

pass 'plans'
