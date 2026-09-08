#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
repo="$ws/repositories/api"
base_before=$(git -C "$repo" rev-parse HEAD)

# A dirty connected checkout is untouched: pairing branches from its committed
# base tip into a separate worktree and never needs to clean the checkout.
printf 'user work\n' >"$repo/src/user-work.txt"
base_dirty_before=$(git -C "$repo" status --porcelain)
begin=$(cc_pair_begin "$ws" api tighten-checkout)
printf '%s\n' "$begin" | grep -Fq 'branch: cc-pair/tighten-checkout' || fail 'pair branch missing'
printf '%s\n' "$begin" | grep -Fq 'supervision: human-supervised' || fail 'supervision label missing'
assert_eq "$base_dirty_before" "$(git -C "$repo" status --porcelain)"
assert_eq "$base_before" "$(git -C "$repo" rev-parse HEAD)"

pdir="$ws/.runtime/pairing/tighten-checkout"
pointer="$pdir/pointer.yaml"
wt="$ws/.runtime/explore/tighten-checkout/api"
require_file "$pointer"
require_dir "$wt"
test ! -d "$ws/.runtime/worktrees/cc-pair/tighten-checkout" || fail 'Explore mixed into plan-execution worktrees'
contains "$pointer" 'schema_version: 1'
contains "$pointer" 'repo: api'
contains "$pointer" "worktree: $wt"
contains "$pointer" 'branch: cc-pair/tighten-checkout'
contains "$pointer" "base: $base_before"
not_contains "$pointer" 'verifier'
not_contains "$pointer" 'failure'
not_contains "$pointer" 'plan:'
test ! -d "$ws/.runtime/executions/tighten-checkout" || fail 'pairing created an execution'
test ! -e "$ws/.runtime/locks/tighten-checkout" || fail 'pairing created a lease or lock'

inspect=$(cc_pair_inspect "$ws" tighten-checkout)
printf '%s\n' "$inspect" | grep -Fq 'status: active' || fail 'active state missing'
printf '%s\n' "$inspect" | grep -Fq 'resumable: true' || fail 'active session not resumable'
expect_failure cc_pair_begin "$ws" api tighten-checkout
expect_failure cc_pair_begin "$ws" api Bad_Session

# Delivery is separate and an active or dirty session cannot be closed/delivered.
expect_failure cc_pair_delivery_targets "$ws" tighten-checkout
printf 'pair change\n' >"$wt/src/pair-change.txt"
expect_failure cc_pair_close "$ws" tighten-checkout
require_file "$pointer"

# The user explicitly asks for the commit; the runtime never creates it.
git -C "$wt" add -A
git -C "$wt" commit -q -m 'feat(api): tighten checkout'
pair_tip=$(git -C "$wt" rev-parse HEAD)
cc_pair_close "$ws" tighten-checkout >/dev/null
test ! -e "$pointer" || fail 'active pointer remained after close'
require_file "$pdir/closed.yaml"
require_dir "$wt"
git -C "$repo" show-ref --verify --quiet refs/heads/cc-pair/tighten-checkout \
  || fail 'pair branch was not preserved'
assert_eq "$pair_tip" "$(git -C "$wt" rev-parse HEAD)"

closed=$(cc_pair_inspect "$ws" tighten-checkout)
printf '%s\n' "$closed" | grep -Fq 'status: closed' || fail 'closed state missing'
printf '%s\n' "$closed" | grep -Fq 'resumable: false' || fail 'closed session reported resumable'

targets=$(cc_pair_delivery_targets "$ws" tighten-checkout)
printf '%s\n' "$targets" | grep -Fq 'source_branch: cc-pair/tighten-checkout' || fail 'pair source missing'
printf '%s\n' "$targets" | grep -Fq 'target_branch: development' || fail 'pair target missing'
printf '%s\n' "$targets" | grep -Fq 'drift_detected: false' || fail 'unexpected pair drift'
printf '%s\n' "$targets" | grep -Fq 'result_label: human-supervised' || fail 'delivery label missing'

# Explicit post-session worktree cleanup does not remove or make the branch
# undeliverable; the closed pointer retains the small amount of routing state.
git -C "$repo" worktree remove "$wt"
test ! -d "$wt" || fail 'explicit worktree cleanup failed'
cc_pair_delivery_targets "$ws" tighten-checkout >/dev/null

# If the base branch advances outside the pairing branch, delivery blocks rather than
# silently rebasing or inventing independent verification.
rm -f "$repo/src/user-work.txt"
printf 'base branch advance\n' >"$repo/src/base-change.txt"
git -C "$repo" add -A
git -C "$repo" commit -q -m 'feat(api): advance base branch'
drift_out=$(mktemp "${TMPDIR:-/tmp}/cc-pair-drift.XXXXXX")
drift_err=$(mktemp "${TMPDIR:-/tmp}/cc-pair-drift.XXXXXX")
if cc_pair_delivery_targets "$ws" tighten-checkout >"$drift_out" 2>"$drift_err"; then
  fail 'drifted pair delivery succeeded'
fi
contains "$drift_out" 'drift_detected: true'
contains "$drift_err" 'PAIR_BASE_DRIFT tighten-checkout'
rm -f "$drift_out" "$drift_err"

# An explicitly named base is resolved to a commit but still receives a fresh
# pairing branch/worktree; no plan branch is edited in place.
cc_pair_begin "$ws" api from-explicit-base "$base_before" >/dev/null
explicit="$ws/.runtime/pairing/from-explicit-base/pointer.yaml"
contains "$explicit" "base: $base_before"
assert_eq "$base_before" "$(git -C "$ws/.runtime/explore/from-explicit-base/api" rev-parse HEAD)"

# Host adapters invoke the engine from the workspace with `.`. Pairing state
# and Git must agree on one absolute worktree path in that mode.
(cd "$ws" && sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" pair-begin . api relative-root >/dev/null)
relative_pointer="$ws/.runtime/pairing/relative-root/pointer.yaml"
relative_wt="$ws/.runtime/explore/relative-root/api"
require_dir "$relative_wt"
contains "$relative_pointer" "worktree: $relative_wt"
(cd "$ws" && sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" pair-inspect . relative-root) \
	| grep -Fq 'resumable: true' || fail 'relative-root pair is not resumable'

contains "$ROOT/.agents/skills/cc-pair/SKILL.md" "Model & effort per role"
contains "$ROOT/.agents/skills/cc-pair/SKILL.md" "role-tiering"
contains "$ROOT/.agents/skills/cc-pair/SKILL.md" "workspace root"
contains "$ROOT/.agents/skills/cc-pair/SKILL.md" "isolated working copy"
contains "$ROOT/.agents/skills/cc-pair/SKILL.md" "runtime/explore"
contains "$ROOT/.agents/skills/cc-pair/SKILL.md" "never invent"
contains "$ROOT/.agents/skills/cc-execute/SKILL.md" "Model & effort per role"
contains "$ROOT/.agents/skills/cc-execute/SKILL.md" "workspace root"

# Live Explore sessions survive cleanup; closed ones are removed.
cc_runtime_cleanup "$ws" >/dev/null
require_dir "$relative_wt"
require_dir "$ws/.runtime/explore/from-explicit-base/api"
cc_pair_close "$ws" relative-root >/dev/null
require_dir "$relative_wt"
mkdir -p "$ws/.code-review-graph"
printf 'junk\n' >"$ws/.code-review-graph/graph.db"
cc_runtime_cleanup "$ws" >/dev/null
test ! -d "$relative_wt" || fail 'cleanup left closed Explore worktree'
require_dir "$ws/.runtime/explore/from-explicit-base/api"
test ! -e "$ws/.code-review-graph" || fail 'cleanup kept leftover .code-review-graph'

pass 'direct collaboration'
