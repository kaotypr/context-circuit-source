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

# --- the runtime performs no delivery itself ---
not_contains "$ROOT/wrapper/runtime/engine.sh" "git push"
not_contains "$ROOT/wrapper/runtime/engine.sh" "git merge"
not_contains "$ROOT/wrapper/runtime/engine.sh" "pull request"
not_contains "$ROOT/wrapper/runtime/engine.sh" "gh pr"

# --- delivery is not a runtime command ---
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" open-pull-request "$ws" 0001-deliver

# --- the delivery skill fixes source, target, and the blocked boundary ---
skill="$ROOT/.agents/skills/cc-deliver/SKILL.md"
require_file "$skill"
contains "$skill" "execution branch"
contains "$skill" "anchor_branch"
contains "$skill" "default_branch"
contains "$skill" "blocked"

pass 'delivery'
