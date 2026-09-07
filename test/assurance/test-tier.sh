#!/bin/sh
# Context Circuit v1.0 — the one safety-critical automated check: consequence tier classification and the
# per-tier verifier floor (INV-ASSURE-01). This is now the ONLY safety-critical
# automated check (v1.0 removed the scope-envelope check; scope-safety is at Gate 2);
# it carries the deepest fixtures and must FAIL UPWARD (when unsure, tier higher).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" "$@"; }
tier_of() { eng tier-classify "$1" "$2" | sed -n 's/^classified_tier: //p'; }
explore_ok() { eng tier-classify "$1" "$2" | sed -n 's/^explore_ok: //p'; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- a scope path under an auth/secrets surface -> Critical, verifier required ---
cc_fx_intent "$ws" i001-auth "Auth" api "src/auth"
assert_eq "critical" "$(tier_of "$ws" i001-auth)"
assert_eq "no" "$(explore_ok "$ws" i001-auth)"

# --- a money / payments surface -> Critical ---
cc_fx_intent "$ws" i002-pay "Billing" api "src/billing"
assert_eq "critical" "$(tier_of "$ws" i002-pay)"

# --- a data migration surface -> Critical ---
cc_fx_intent "$ws" i003-mig "Migrate" api "db/migrations"
assert_eq "critical" "$(tier_of "$ws" i003-mig)"

# --- a single-repo, bounded, no-signal change -> Standard by default (fail upward);
# the engine cannot see reversibility/coverage/novelty, so no-signal is NOT auto-
# Explore. Explore stays available only as an explicit human lowering (explore_ok=yes) ---
cc_fx_intent "$ws" i004-widget "Widget" api "src/widget"
assert_eq "standard" "$(tier_of "$ws" i004-widget)"
assert_eq "yes" "$(explore_ok "$ws" i004-widget)"

# --- an unrecognized / uncontained shape (repo-wide scope) -> Standard, not Explore ---
cc_fx_intent "$ws" i005-wide "Wide" api "."
assert_eq "standard" "$(tier_of "$ws" i005-wide)"
assert_eq "no" "$(explore_ok "$ws" i005-wide)"

# --- more than one repository -> Standard, not Explore ---
mkdir -p "$ws/intent/i006-multi"
cat >"$ws/intent/i006-multi/contract.yaml" <<Y
schema_version: 1
intent: i006-multi
title: Multi
goal: g
acceptance_criteria:
  - id: ac-1
    statement: s
    method: test
scope:
  repositories:
    - id: api
      paths: [src]
    - id: web
      paths: [src]
tier: standard
status: draft
contract_digest:
Y
printf '# m\n' >"$ws/intent/i006-multi/INTENT.md"
assert_eq "standard" "$(tier_of "$ws" i006-multi)"
assert_eq "no" "$(explore_ok "$ws" i006-multi)"

# --- DEEP BOUNDARY FIXTURES (fail-upward corners; each pins one under-tiering trap). ---

# Hidden critical among benign paths: a Critical surface sharing scope with an
# otherwise Explore-eligible path must still tier Critical and forbid Explore. The
# highest signal across all scope paths wins; a benign-first classifier that stopped
# at src/widget would ship an auth change unverified. This is the headline tier guard.
cc_fx_intent "$ws" i008-hidden "Hidden critical" api "src/widget src/auth"
assert_eq "critical" "$(tier_of "$ws" i008-hidden)"
assert_eq "no" "$(explore_ok "$ws" i008-hidden)"

# Keyword-coverage guards: the Critical surface set is more than "auth". Pin a few
# distinct high-consequence surfaces so narrowing the signal set is caught. Each is a
# single scope path; a regression that drops one silently under-tiers that whole class.
for cc_ti_case in secrets:i009-sec src/payments:i010-pay deploy:i011-dep src/credentials:i012-cred config/production:i013-prod; do
	cc_ti_path=${cc_ti_case%%:*}; cc_ti_id=${cc_ti_case#*:}
	cc_fx_intent "$ws" "$cc_ti_id" "Critical surface" api "$cc_ti_path"
	assert_eq "critical" "$(tier_of "$ws" "$cc_ti_id")"
	assert_eq "no" "$(explore_ok "$ws" "$cc_ti_id")"
done

# --- lowering guard: Explore on a risk surface is refused; Standard is allowed ---
expect_failure eng tier-lower-check "$ws" i001-auth explore     # Critical->Explore for security: refused
eng tier-lower-check "$ws" i001-auth standard >/dev/null        # Critical->Standard: allowed (verifier kept)

# --- the floor is enforced at the gate: approving Explore on a risk surface is refused ---
awk '/^tier:/{print "tier: explore"; next}{print}' "$ws/intent/i001-auth/contract.yaml" >"$ws/intent/i001-auth/c.new"
mv "$ws/intent/i001-auth/c.new" "$ws/intent/i001-auth/contract.yaml"
expect_failure eng intent-approve "$ws" i001-auth
# an Explore-eligible intent approves at Explore
awk '/^tier:/{print "tier: explore"; next}{print}' "$ws/intent/i004-widget/contract.yaml" >"$ws/intent/i004-widget/c.new"
mv "$ws/intent/i004-widget/c.new" "$ws/intent/i004-widget/contract.yaml"
eng intent-approve "$ws" i004-widget >/dev/null

# --- per-tier verifier floor as an eligibility query ---
cc_fx_repo "$ws" api development

# Explore is human-supervised and planless in v1.0. Promotion to Standard/Critical
# is the first point at which execution creates a plan of record.
cc_fx_plan_intent "$ws" 0001-widget "Widget" api src/widget i004-widget
expect_failure eng plan-validate "$ws/plans/0001-widget"
expect_failure eng execution-begin "$ws" 0001-widget sess-e
# Standard: the independent verifier floor is required (acceptance alone is not enough).
cc_fx_intent "$ws" i007-std "Std" api "src/std"
eng intent-approve "$ws" i007-std >/dev/null
cc_fx_plan_intent "$ws" 0002-std "Std" api src/std i007-std
exec2=$(eng execution-begin "$ws" 0002-std sess-s | sed -n 's/^execution_id: //p')
edir2=$(cc_fx_exec_dir "$ws" 0002-std "$exec2")
contains "$edir2/execution.yaml" "tier: standard"
cc_attempt_begin "$edir2" >/dev/null
wt2="$ws/.runtime/worktrees/0002-std/api"
mkdir -p "$wt2/src/std"; printf 'y\n' >"$wt2/src/std/mod.txt"
git -C "$wt2" add -A; git -C "$wt2" commit -q -m "feat(api): std"
cc_worker_commit_record "$edir2" api implementation >/dev/null
# a human acceptance without the verifier does NOT satisfy the Standard
# eligibility query (completion-ready); it is not a mark-done gate.
eng human-acceptance-record "$edir2" alice >/dev/null
expect_failure eng completion-ready "$ws" 0002-std
# the independent verifier pass satisfies the floor
cc_verifier_prepare "$edir2" >/dev/null
cc_verifier_result_record "$edir2" 1 passed >/dev/null
out2=$(eng completion-ready "$ws" 0002-std)
printf '%s\n' "$out2" | grep -q 'assurance: independent' || fail "Standard completion must report independent assurance"

pass 'consequence tiering (the one safety-critical automated check)'
