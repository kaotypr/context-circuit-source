#!/bin/sh
# Context Circuit v1.0 — crown jewel 2: consequence tier classification and the
# per-tier verifier floor (INV-ASSURE-01). The second safety-critical check; it
# carries the deepest fixtures and must FAIL UPWARD (when unsure, tier higher).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

eng() { sh "$ROOT/wrapper/runtime/engine.sh" "$@"; }
tier_of() { eng tier-classify "$1" "$2" | sed -n 's/^classified_tier: //p'; }
explore_ok() { eng tier-classify "$1" "$2" | sed -n 's/^explore_ok: //p'; }

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- a scope path under an auth/secrets surface -> Critical, verifier required ---
cc_fx_intent "$ws" i0001-auth "Auth" api "src/auth"
assert_eq "critical" "$(tier_of "$ws" i0001-auth)"
assert_eq "no" "$(explore_ok "$ws" i0001-auth)"

# --- a money / payments surface -> Critical ---
cc_fx_intent "$ws" i0002-pay "Billing" api "src/billing"
assert_eq "critical" "$(tier_of "$ws" i0002-pay)"

# --- a data migration surface -> Critical ---
cc_fx_intent "$ws" i0003-mig "Migrate" api "db/migrations"
assert_eq "critical" "$(tier_of "$ws" i0003-mig)"

# --- a single-repo, bounded, no-signal change -> Standard by default (fail upward);
# the engine cannot see reversibility/coverage/novelty, so no-signal is NOT auto-
# Explore. Explore stays available only as an explicit human lowering (explore_ok=yes) ---
cc_fx_intent "$ws" i0004-widget "Widget" api "src/widget"
assert_eq "standard" "$(tier_of "$ws" i0004-widget)"
assert_eq "yes" "$(explore_ok "$ws" i0004-widget)"

# --- an unrecognized / uncontained shape (repo-wide scope) -> Standard, not Explore ---
cc_fx_intent "$ws" i0005-wide "Wide" api "."
assert_eq "standard" "$(tier_of "$ws" i0005-wide)"
assert_eq "no" "$(explore_ok "$ws" i0005-wide)"

# --- more than one repository -> Standard, not Explore ---
mkdir -p "$ws/intent/i0006-multi"
cat >"$ws/intent/i0006-multi/contract.yaml" <<Y
schema_version: 1
intent: i0006-multi
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
printf '# m\n' >"$ws/intent/i0006-multi/INTENT.md"
assert_eq "standard" "$(tier_of "$ws" i0006-multi)"
assert_eq "no" "$(explore_ok "$ws" i0006-multi)"

# --- lowering guard: Explore on a risk surface is refused; Standard is allowed ---
expect_failure eng tier-lower-check "$ws" i0001-auth explore     # Critical->Explore for security: refused
eng tier-lower-check "$ws" i0001-auth standard >/dev/null        # Critical->Standard: allowed (verifier kept)

# --- the floor is enforced at the gate: approving Explore on a risk surface is refused ---
awk '/^tier:/{print "tier: explore"; next}{print}' "$ws/intent/i0001-auth/contract.yaml" >"$ws/intent/i0001-auth/c.new"
mv "$ws/intent/i0001-auth/c.new" "$ws/intent/i0001-auth/contract.yaml"
expect_failure eng intent-approve "$ws" i0001-auth
# an Explore-eligible intent approves at Explore
awk '/^tier:/{print "tier: explore"; next}{print}' "$ws/intent/i0004-widget/contract.yaml" >"$ws/intent/i0004-widget/c.new"
mv "$ws/intent/i0004-widget/c.new" "$ws/intent/i0004-widget/contract.yaml"
eng intent-approve "$ws" i0004-widget >/dev/null

# --- per-tier verifier floor at completion ---
cc_fx_repo "$ws" api development

# Explore: human-supervised, no verifier — eligible on a candidate-bound acceptance,
# and completion-ready reports assurance human-supervised, never "verified".
cc_fx_plan_intent "$ws" 0001-widget "Widget" api src/widget i0004-widget
exec=$(eng execution-begin "$ws" 0001-widget sess-e | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-widget "$exec")
contains "$edir/execution.yaml" "tier: explore"
cc_attempt_begin "$edir" >/dev/null
wt="$ws/.runtime/worktrees/0001-widget/api"
mkdir -p "$wt/src/widget"; printf 'x\n' >"$wt/src/widget/mod.txt"
git -C "$wt" add -A; git -C "$wt" commit -q -m "feat(api): widget"
cc_worker_commit_record "$edir" api implementation >/dev/null
# no verifier runs at Explore; completion is blocked until a human accepts
expect_failure eng completion-ready "$ws" 0001-widget
eng human-acceptance-record "$edir" alice >/dev/null
out=$(eng completion-ready "$ws" 0001-widget)
printf '%s\n' "$out" | grep -q 'assurance: human-supervised' || fail "Explore completion must report human-supervised"
printf '%s\n' "$out" | grep -q 'verified' && fail "Explore output must never be labeled verified" || :

# Standard: the independent verifier floor is required (acceptance alone is not enough).
cc_fx_intent "$ws" i0007-std "Std" api "src/std"
eng intent-approve "$ws" i0007-std >/dev/null
cc_fx_plan_intent "$ws" 0002-std "Std" api src/std i0007-std
exec2=$(eng execution-begin "$ws" 0002-std sess-s | sed -n 's/^execution_id: //p')
edir2=$(cc_fx_exec_dir "$ws" 0002-std "$exec2")
contains "$edir2/execution.yaml" "tier: standard"
cc_attempt_begin "$edir2" >/dev/null
wt2="$ws/.runtime/worktrees/0002-std/api"
mkdir -p "$wt2/src/std"; printf 'y\n' >"$wt2/src/std/mod.txt"
git -C "$wt2" add -A; git -C "$wt2" commit -q -m "feat(api): std"
cc_worker_commit_record "$edir2" api implementation >/dev/null
# a human acceptance without the verifier does NOT satisfy a Standard completion
eng human-acceptance-record "$edir2" alice >/dev/null
expect_failure eng completion-ready "$ws" 0002-std
# the independent verifier pass satisfies the floor
cc_verifier_prepare "$edir2" >/dev/null
cc_verifier_result_record "$edir2" 1 passed >/dev/null
out2=$(eng completion-ready "$ws" 0002-std)
printf '%s\n' "$out2" | grep -q 'assurance: independent' || fail "Standard completion must report independent assurance"

pass 'consequence tiering (crown jewel 2)'
