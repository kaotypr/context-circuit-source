#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
repo="$ws/repositories/api"
cc_fx_ignored_content "$repo"

# A newly created execution worktree receives ignored files, not ordinary
# untracked files. The direct prepare surface is also used by execution setup.
cc_worktree_prepare "$ws" 0000-overlay api >/dev/null
preflight_wt="$ws/.runtime/worktrees/0000-overlay/api"
assert_eq 'fixture-overlay-value' "$(cat "$preflight_wt/.env")"
assert_eq 'module fixture' "$(cat "$preflight_wt/node_modules/nested/file.txt")"
assert_eq 'spaced fixture' "$(cat "$preflight_wt/ignored directory with spaces/file.txt")"
test ! -e "$preflight_wt/untracked.txt" || fail 'non-ignored untracked file was overlaid'
test ! -L "$preflight_wt/.env" || fail 'overlaid file is a symlink'
test ! -L "$preflight_wt/node_modules" || fail 'overlaid directory is a symlink'
test -z "$(git -C "$preflight_wt" status --porcelain)" || fail 'overlay dirtied worktree status'

# execution-begin intentionally requires a clean bound checkout, so remove only
# the non-ignored fixture file before proving its full setup path.
rm -f "$repo/untracked.txt"
cc_fx_plan_ex "$ws" 0001-overlay "Overlay" api src ""
exec=$(cc_execution_begin "$ws" 0001-overlay overlay-worker | sed -n 's/^execution_id: //p')
edir="$ws/.runtime/executions/0001-overlay/$exec"
wt="$ws/.runtime/worktrees/0001-overlay/api"
assert_eq 'fixture-overlay-value' "$(cat "$wt/.env")"
assert_eq 'module fixture' "$(cat "$wt/node_modules/nested/file.txt")"
assert_eq 'spaced fixture' "$(cat "$wt/ignored directory with spaces/file.txt")"
test ! -L "$wt/.env" || fail 'overlaid file is a symlink'
test ! -L "$wt/node_modules" || fail 'overlaid directory is a symlink'
test -z "$(git -C "$wt" status --porcelain)" || fail 'overlay dirtied worktree status'
not_contains "$edir/repositories/api.yaml" '.env'
not_contains "$edir/execution.yaml" 'fixture-overlay-value'
not_contains "$edir/grounding/api.yaml" 'fixture-overlay-value'

# A reused worktree is worker-owned and is never overlaid again.
printf 'worker override\n' >"$wt/.env"
cc_worktree_prepare "$ws" 0001-overlay api >/dev/null
assert_eq 'worker override' "$(cat "$wt/.env")"

# Explore worktrees receive the same overlay before their pointer is written.
pair=$(cc_pair_begin "$ws" api overlay-explore)
printf '%s\n' "$pair" | grep -q '^copy_mode: ' || fail 'pair output omitted copy mode'
pwt="$ws/.runtime/explore/overlay-explore/api"
assert_eq 'fixture-overlay-value' "$(cat "$pwt/.env")"
assert_eq 'module fixture' "$(cat "$pwt/node_modules/nested/file.txt")"
assert_eq 'spaced fixture' "$(cat "$pwt/ignored directory with spaces/file.txt")"
test ! -e "$pwt/untracked.txt" || fail 'Explore copied non-ignored untracked file'
test ! -L "$pwt/.env" || fail 'Explore overlay is a symlink'
test -z "$(git -C "$pwt" status --porcelain)" || fail 'Explore overlay dirtied worktree status'

# A workspace-root binding must not recursively overlay its own state.
root_ws=$(cc_fx_ws)
printf '.runtime/\nrepositories/\nrole-tiering.local.yaml\nmember.local.yaml\nrepositories.local.yaml\n' >"$root_ws/.gitignore"
cc_fx_repo "$root_ws" rootrepo development .
cc_fx_ignored_content "$root_ws"
mkdir -p "$root_ws/repositories" "$root_ws/.runtime/overlay-source"
printf 'runtime state\n' >"$root_ws/.runtime/overlay-source/state.txt"
printf 'registered checkout\n' >"$root_ws/repositories/state.txt"
printf 'host state\n' >"$root_ws/role-tiering.local.yaml"
root_wt="$root_ws/.runtime/worktrees/0002-root/rootrepo"
cc_worktree_prepare "$root_ws" 0002-root rootrepo >/dev/null
assert_eq 'fixture-overlay-value' "$(cat "$root_wt/.env")"
test ! -e "$root_wt/.runtime" || fail 'runtime state was overlaid'
test ! -e "$root_wt/repositories" || fail 'registered repositories were overlaid'
test ! -e "$root_wt/role-tiering.local.yaml" || fail 'role tiering was overlaid'
test ! -e "$root_wt/member.local.yaml" || fail 'member identity was overlaid'
test ! -e "$root_wt/repositories.local.yaml" || fail 'repository bindings were overlaid'
rm -rf "$root_ws"

# A forced overlay failure leaves no resumable Explore state and blocks execution
# before grounding is written.
cc_overlay_ignored() { cc_fail OVERLAY_FAILED forced; }
expect_failure cc_pair_begin "$ws" api blocked-explore
test ! -e "$ws/.runtime/pairing/blocked-explore/pointer.yaml" || fail 'failed Explore setup wrote pointer'
test ! -e "$ws/.runtime/explore/blocked-explore/api" || fail 'failed Explore setup left worktree'
git -C "$repo" show-ref --verify --quiet refs/heads/cc-pair/blocked-explore \
	&& fail 'failed Explore setup left branch'
cc_fx_plan_ex "$ws" 0003-blocked "Blocked overlay" api src ""
expect_failure cc_execution_begin "$ws" 0003-blocked blocked-worker
for blocked_edir in "$ws/.runtime/executions/0003-blocked"/exec-*; do break; done
contains "$blocked_edir/execution.yaml" 'status: blocked'
contains "$blocked_edir/execution.yaml" 'blocked_reason: WORKTREE_OVERLAY_FAILED'
test ! -e "$blocked_edir/grounding/api.yaml" || fail 'blocked overlay wrote grounding evidence'

pass 'worktree ignored overlay'
