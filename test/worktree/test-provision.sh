#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# A matching ignored install tree is retained without reaching the install seam.
cc_fx_repo "$ws" api development
repo="$ws/repositories/api"
cc_fx_node_toolchain "$repo" 'lock-v1'
cc_worktree_prepare "$ws" 0040-keep api >/dev/null
wt="$ws/.runtime/worktrees/0040-keep/api"
cc_toolchain_install() { fail 'matching tree unexpectedly installed'; }
keep=$(cc_provision_worktree "$repo" "$wt")
printf '%s\n' "$keep" | grep -Fqx 'provision: kept' || fail 'matching tree was not kept'
printf '%s\n' "$keep" | grep -Fqx 'environment: ready' || fail 'kept tree was not ready'
assert_eq "$(cc_digest "$wt/package-lock.json")" "$(cc_toolchain_stamp_read "$wt/node_modules")"

# A changed worktree lockfile replaces the copied tree using the one install seam.
printf 'lock-v2\n' >"$wt/package-lock.json"
rm -f "$wt/node_modules/.context-circuit-toolchain-lock-digest"
install_mark="$ws/install-ran"
cc_toolchain_install() {
	mkdir -p "$1/node_modules"
	printf '%s\n' "$2" >"$install_mark"
}
installed=$(cc_provision_worktree "$repo" "$wt")
printf '%s\n' "$installed" | grep -Fqx 'provision: installed' || fail 'mismatching tree was not installed'
assert_eq 'package-lock.json' "$(cat "$install_mark")"
test ! -e "$wt/node_modules/fixture/index.js" || fail 'stale install tree was kept'
assert_eq "$(cc_digest "$wt/package-lock.json")" "$(cc_toolchain_stamp_read "$wt/node_modules")"

# Explicit global-cache rows and greenfield repositories never invoke installs.
go_src="$ws/go-source"; go_wt="$ws/go-worktree"
mkdir -p "$go_src" "$go_wt"
printf 'module example.test/go\n' >"$go_src/go.mod"
cp "$go_src/go.mod" "$go_wt/go.mod"
global=$(cc_provision_worktree "$go_src" "$go_wt")
printf '%s\n' "$global" | grep -Fqx 'provision: none' || fail 'go unexpectedly provisioned a tree'
printf '%s\n' "$global" | grep -Fqx 'environment: ready' || fail 'go was not ready'
test ! -e "$go_wt/node_modules" || fail 'global-cache row created an install tree'

empty_src="$ws/empty-source"; empty_wt="$ws/empty-worktree"
mkdir -p "$empty_src" "$empty_wt"
greenfield=$(cc_provision_worktree "$empty_src" "$empty_wt")
printf '%s\n' "$greenfield" | grep -Fqx 'provision: none' || fail 'greenfield provisioned'
printf '%s\n' "$greenfield" | grep -Fqx 'environment: no-toolchain' || fail 'greenfield was not no-toolchain'

# A missing tool fails closed for both execution and Explore setup. Override only
# the runtime table row; no package manager or network is involved in this suite.
. "$ROOT/.context-circuit/wrapper/runtime/engine.sh"
cc_toolchain_row() {
	case "$1" in
		bun.lockb) printf 'tree: node_modules\ninstall: cc-fixture-missing-tool\n' ;;
		*) cc_fail TOOLCHAIN_KEY_UNKNOWN "$1"; return 1 ;;
	esac
}
cc_fx_repo "$ws" broken development
broken="$ws/repositories/broken"
printf 'lock\n' >"$broken/bun.lockb"
git -C "$broken" add bun.lockb
git -C "$broken" commit -q -m 'test: add broken toolchain fixture'
cc_fx_plan_ex "$ws" 0041-provision-block "Provision block" broken src ""
expect_failure cc_execution_begin "$ws" 0041-provision-block provision-worker
for edir in "$ws/.runtime/executions/0041-provision-block"/exec-*; do break; done
contains "$edir/execution.yaml" 'status: blocked'
contains "$edir/execution.yaml" 'blocked_reason: TOOLCHAIN_PROVISION_FAILED'
test ! -e "$edir/brief-broken.md" || fail 'blocked provisioning assembled a brief'
test ! -e "$edir/grounding/broken.yaml" || fail 'blocked provisioning wrote grounding evidence'

expect_failure cc_pair_begin "$ws" broken provision-block
test ! -e "$ws/.runtime/pairing/provision-block/pointer.yaml" || fail 'failed Explore setup wrote a pointer'
test ! -e "$ws/.runtime/explore/provision-block/broken" || fail 'failed Explore setup left a worktree'
git -C "$broken" show-ref --verify --quiet refs/heads/cc-pair/provision-block \
	&& fail 'failed Explore setup left a branch'

pass 'worktree toolchain provisioning'
