#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

cc_fx_plan "$ws" 0001-deliver "Deliver" "api web"
cc_plan_approve "$ws" 0001-deliver >/dev/null
exec=$(cc_execution_begin "$ws" 0001-deliver sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-deliver "$exec")

# --- the runtime records the pull-request source (execution branch) and
#     default target (recorded anchor branch) for each repository ---
contains "$edir/repositories/api.yaml" "branch: cc/0001-deliver/api"
contains "$edir/repositories/api.yaml" "anchor_branch: development"

# --- the runtime performs no delivery action itself (report-only) ---
not_contains "$ROOT/wrapper/runtime/engine.sh" "git push"
not_contains "$ROOT/wrapper/runtime/engine.sh" "gh pr"
not_contains "$ROOT/wrapper/runtime/engine.sh" "git -C \"\$cc_dt_abs\" push"
# v0.6: the runtime authors an integration BASE merge (INV-CONCURRENCY-02) on the
# dependent plan's OWN worktree, before the worker starts — never a delivery merge.
# Every `git ... merge ` in the engine must be that base merge ($cc_bp_tree), so no
# merge ever runs against a delivery target or the anchor checkout. (merge-base and
# rebase are not merges into a branch.)
if grep -nE 'git[^\n]*merge( |$)' "$ROOT/wrapper/runtime/engine.sh" \
	| grep -vE 'merge --abort' | grep -vE '\$cc_bp_tree' >/dev/null 2>&1; then
	grep -nE 'git[^\n]*merge( |$)' "$ROOT/wrapper/runtime/engine.sh" \
		| grep -vE 'merge --abort' | grep -vE '\$cc_bp_tree'
	fail 'engine performs a git merge outside the integration-base builder'
fi

# --- delivery is not a runtime command ---
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" open-pull-request "$ws" 0001-deliver

# --- the read-only delivery-targets helper reports source and target per repo ---
dt=$(cc_delivery_targets "$ws" 0001-deliver)
printf '%s\n' "$dt" | grep -Fq "source_branch: cc/0001-deliver/api" || fail "missing api source branch"
printf '%s\n' "$dt" | grep -Fq "target_branch: development" || fail "missing target branch"
printf '%s\n' "$dt" | grep -Fq "source_present: true" || fail "execution branch should exist"
# it never performs delivery and needs an execution to report
expect_failure cc_delivery_targets "$ws" 9999-none

# --- the delivery skill fixes source, target, and the blocked boundary ---
skill="$ROOT/.agents/skills/cc-deliver/SKILL.md"
require_file "$skill"
contains "$skill" "execution branch"
contains "$skill" "anchor_branch"
contains "$skill" "default_branch"
contains "$skill" "blocked"

pass 'delivery'
