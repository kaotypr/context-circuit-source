#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
cc_fx_repo "$ws" web development

cc_fx_plan "$ws" 0001-deliver "Deliver" "api"
exec=$(cc_execution_begin "$ws" 0001-deliver sess1 | sed -n 's/^execution_id: //p')
edir=$(cc_fx_exec_dir "$ws" 0001-deliver "$exec")

# --- the runtime records the pull-request source (execution branch) and
#     default target (recorded base branch) for each repository ---
contains "$edir/repositories/api.yaml" "branch: cc/0001-deliver/api"
contains "$edir/repositories/api.yaml" "base_branch: development"

# --- the runtime performs no delivery action itself (report-only) ---
not_contains "$ROOT/wrapper/runtime/engine.sh" "git push"
not_contains "$ROOT/wrapper/runtime/engine.sh" "gh pr"
not_contains "$ROOT/wrapper/runtime/engine.sh" "git -C \"\$cc_dt_abs\" push"
# The runtime authors integration merges only into a dedicated integration worktree,
# before the worker/verifier — never a delivery merge, a delivery target, or the base branch
# checkout (INV-CONCURRENCY-02, INV-DELIVER-01). The two integration builders are the
# per-plan base merge ($cc_bp_tree) and the change-set integration tip ($cc_csp_wt).
# Every `git ... merge ` in the engine must be one of those. (merge-base and rebase are
# not merges into a branch.)
if grep -nE 'git[^\n]*merge( |$)' "$ROOT/wrapper/runtime/engine.sh" \
	| grep -vE 'merge --abort' | grep -vE '\$cc_bp_tree|\$cc_csp_wt' >/dev/null 2>&1; then
	grep -nE 'git[^\n]*merge( |$)' "$ROOT/wrapper/runtime/engine.sh" \
		| grep -vE 'merge --abort' | grep -vE '\$cc_bp_tree|\$cc_csp_wt'
	fail 'engine performs a git merge outside an integration builder'
fi

# --- delivery is not a runtime command ---
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" open-pull-request "$ws" 0001-deliver

# --- the read-only delivery-targets helper reports source and target per repo ---
dt=$(cc_delivery_targets "$ws" 0001-deliver)
printf '%s\n' "$dt" | grep -Fq "source_branch: cc/0001-deliver/api" || fail "missing api source branch"
printf '%s\n' "$dt" | grep -Fq "target_branch: development" || fail "missing target branch"
printf '%s\n' "$dt" | grep -Fq "source_present: true" || fail "execution branch should exist"

# Legacy schema-1 execution repository records remain readable after the field
# rename; delivery still resolves their recorded base branch.
awk '{sub(/^base_branch:/, "anchor_branch:"); print}' "$edir/repositories/api.yaml" >"$edir/repositories/api.yaml.new"
mv "$edir/repositories/api.yaml.new" "$edir/repositories/api.yaml"
legacy_dt=$(cc_delivery_targets "$ws" 0001-deliver)
printf '%s\n' "$legacy_dt" | grep -Fq "target_branch: development" || fail "legacy base branch record not readable"

# it never performs delivery and needs an execution to report
expect_failure cc_delivery_targets "$ws" 9999-none

# --- the delivery skill fixes source, target, and the blocked boundary ---
skill="$ROOT/.agents/skills/cc-deliver/SKILL.md"
require_file "$skill"
contains "$skill" "execution branch"
contains "$skill" "base_branch"
contains "$skill" "default_branch"
contains "$skill" "blocked"

pass 'delivery'
