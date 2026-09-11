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

# --- archive releases the numeric prefix; restore refuses an active collision ---
cc_plan_archive "$ws" "$id2" >/dev/null
id3=$(cc_plan_allocate_id "$ws" checkout-v2)
assert_eq "0002-checkout-v2" "$id3"
cc_fx_plan "$ws" "$id3" "Checkout v2" "api"
restore_collision=$(cc_plan_restore "$ws" "$id2" 2>&1 || true)
printf '%s\n' "$restore_collision" | grep -Fq RESTORE_PREFIX_COLLISION \
	|| fail "restore did not report active prefix collision"
cc_plan_archive "$ws" "$id3" >/dev/null
cc_plan_restore "$ws" "$id2" >/dev/null

# --- bad slugs and ids are rejected ---
expect_failure cc_plan_allocate_id "$ws" "Bad_Slug"
expect_failure cc_plan_allocate_id "$ws" "with spaces"
cc_plan_id_valid "0001-billing-v2" || fail "valid id rejected"
expect_failure cc_plan_id_valid "1-x"
expect_failure cc_plan_id_valid "0001-Bad"

# --- validation catches missing tasks ---
mkdir -p "$ws/plans/0009-empty/tasks"
printf 'schema_version: 3\nplan: 0009-empty\ntitle: E\nstatus: draft\nobjective: e\nintent: i009-empty\nrepositories:\n  - id: api\ntasks:\n' >"$ws/plans/0009-empty/plan.yaml"
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

# --- authorization: a plan derives from its approved intent (INV-INTENT-02). There is
#     no separate plan-approval status and no automated scope gate — the plan stays
#     `draft`, and the derives-from-an-approved-intent authorization is the preflight
#     that execution re-runs. ---
cc_intent_authorized "$ws" "$id1" >/dev/null
assert_eq "draft" "$(cc_plan_status "$ws" "$id1")"   # no intermediate "approved" status

# --- a plan authored with block-list task fields also validates ---
mkdir -p "$ws/plans/0011-block/tasks"
cat >"$ws/plans/0011-block/plan.yaml" <<'EOF'
schema_version: 3
plan: 0011-block
title: Block form
status: draft
objective: block
intent: i011-block
repositories:
  - id: api
tasks:
  - id: T-001
    title: t1
    repositories:
      - api
    paths:
      - src
    depends_on: []
EOF
printf '# Block\n' >"$ws/plans/0011-block/PLAN.md"
cc_plan_validate "$ws/plans/0011-block" >/dev/null
assert_eq "api" "$(cc_plan_affected_repositories "$ws/plans/0011-block/plan.yaml")"

# --- a plan that lists two repositories is invalid ---
mkdir -p "$ws/plans/0013-tworepo/tasks"
cat >"$ws/plans/0013-tworepo/plan.yaml" <<'EOF'
schema_version: 3
plan: 0013-tworepo
title: Two repos
status: draft
objective: invalid
intent: i013-tworepo
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
printf '# Two\n' >"$ws/plans/0013-tworepo/PLAN.md"
expect_failure cc_plan_validate "$ws/plans/0013-tworepo"

# --- execute on the approved intent's authorization (no separate plan-approval step) ---
cc_fx_plan "$ws" 0010-compound "Compound" "api"
cc_execution_begin "$ws" 0010-compound s2 >/dev/null
require_dir "$ws/.runtime/executions/0010-compound"
assert_eq "draft" "$(cc_plan_status "$ws" 0010-compound)"   # stays draft through execution

# --- one bounded Standard change remains one plan with ordered embedded tasks ---
cc_fx_plan "$ws" 0012-bounded "Bounded" "api"
multi="$ws/plans/0012-bounded/plan.yaml"
awk '
/^execution:/ && !added {
    print "  - id: API-002"
    print "    title: Work in src/billing/receipt"
    print "    repositories: [api]"
    print "    paths: [src/billing/receipt]"
    print "    depends_on: [API-001]"
    print "    changes:" 
    print "      - Change the receipt path after the bounded billing task."
    print "    acceptance:"
    print "      - id: API-002-AC"
    print "        statement: The receipt path is changed."
    print "    verification:"
    print "      - id: API-002-VT"
    print "        command: test -f src/billing/receipt/mod.txt"
    added=1
}
{print}
' "$multi" >"$multi.new"
mv "$multi.new" "$multi"
cc_plan_validate "$ws/plans/0012-bounded" >/dev/null
contains "$multi" "depends_on: [API-001]"
not_contains "$multi" "plan_dependencies:"
contains "$multi" "worker: one"
contains "$multi" "independent_verifier: required"

pass 'plans'
