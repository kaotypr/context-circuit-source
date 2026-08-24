#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- bind an existing checkout with a personal anchor branch ---
cc_fx_repo "$ws" api kao/development/v0.5
out=$(cc_repo_resolve "$ws" api)
printf '%s\n' "$out" | grep -Fq "anchor_branch: kao/development/v0.5" || fail "anchor not resolved"

# --- default clone destination is under the gitignored repositories/ ---
printf '%s\n' "$out" | grep -Fq "repositories/api" || fail "default destination not repositories/"

# --- missing binding is rejected without filesystem scanning ---
expect_failure cc_repo_resolve "$ws" ghost

# --- a non-Git path is rejected ---
mkdir -p "$ws/repositories/plain"
cc_fx_bindings_header "$ws"
printf '  plain:\n    path: repositories/plain\n    anchor_branch: main\n' >>"$ws/repositories.local.yaml"
expect_failure cc_repo_resolve "$ws" plain

# --- git init + initial commit: a repo with no commit is not executable ---
mkdir -p "$ws/repositories/fresh"
git -C "$ws/repositories/fresh" init -q -b development
printf '  fresh:\n    path: repositories/fresh\n    anchor_branch: development\n' >>"$ws/repositories.local.yaml"
cc_repo_resolve "$ws" fresh >/dev/null   # binding resolves
expect_failure cc_repo_anchor_commit "$ws" fresh   # but no anchor tip yet
git -C "$ws/repositories/fresh" config user.email t@t.t
git -C "$ws/repositories/fresh" config user.name t
git -C "$ws/repositories/fresh" commit -q --allow-empty -m bootstrap
cc_repo_anchor_commit "$ws" fresh >/dev/null   # now executable

# --- anchor_branch is validated, not inferred from default_branch ---
mkdir -p "$ws/repositories/nomatch"
git -C "$ws/repositories/nomatch" init -q -b main
git -C "$ws/repositories/nomatch" config user.email t@t.t
git -C "$ws/repositories/nomatch" config user.name t
git -C "$ws/repositories/nomatch" commit -q --allow-empty -m seed
printf '  nomatch:\n    path: repositories/nomatch\n    anchor_branch: development\n' >>"$ws/repositories.local.yaml"
expect_failure cc_repo_anchor_commit "$ws" nomatch   # anchor 'development' does not exist

# --- workspace root as the reserved logical repository `workspace` ---
wsroot=$(mktemp -d "${TMPDIR:-/tmp}/cc-wsroot.XXXXXX")
printf 'schema_version: 1\nworkspace: root-ws\nrepositories: []\n' >"$wsroot/workspace.yaml"
cc_workspace_init "$wsroot" >/dev/null
git -C "$wsroot" init -q -b kao/development/v0.5
git -C "$wsroot" config user.email t@t.t
git -C "$wsroot" config user.name t
git -C "$wsroot" add -A
git -C "$wsroot" commit -q -m seed
printf 'schema_version: 1\nbindings:\n  workspace:\n    path: .\n    anchor_branch: kao/development/v0.5\n' >"$wsroot/repositories.local.yaml"
wout=$(cc_repo_resolve "$wsroot" workspace)
printf '%s\n' "$wout" | grep -Fq "anchor_branch: kao/development/v0.5" || fail "workspace repo anchor"
rm -rf "$wsroot"

# --- multi-repository preflight over three repositories ---
cc_fx_repo "$ws" web development
cc_fx_repo "$ws" contracts development
cc_fx_plan "$ws" 0001-multi "Multi" "contracts api web"
# fix api anchor for the plan (api was bound to a personal branch already)
pf=$(cc_repository_preflight "$ws" "$ws/plans/0001-multi")
printf '%s\n' "$pf" | grep -Fq "repositories_ready: 3" || fail "expected 3 repos ready"

# --- dirty anchor checkout is rejected by preflight ---
printf 'dirty\n' >"$ws/repositories/web/src/x.txt"
expect_failure cc_repository_preflight "$ws" "$ws/plans/0001-multi"
git -C "$ws/repositories/web" checkout -q -- . 2>/dev/null || rm -f "$ws/repositories/web/src/x.txt"

# --- deterministic registration records identity + binding, then resolves ---
ws2=$(mktemp -d "${TMPDIR:-/tmp}/cc-reg.XXXXXX")
printf 'schema_version: 1\nworkspace: reg-ws\nrepositories: []\n' >"$ws2/workspace.yaml"
cc_workspace_init "$ws2" >/dev/null
mkdir -p "$ws2/repositories/svc/src"
git -C "$ws2/repositories/svc" init -q -b development
git -C "$ws2/repositories/svc" config user.email t@t.t
git -C "$ws2/repositories/svc" config user.name t
git -C "$ws2/repositories/svc" commit -q --allow-empty -m seed
cc_repository_register "$ws2" svc repositories/svc development https://example.com/svc.git main >/dev/null
contains "$ws2/workspace.yaml" "- id: svc"
contains "$ws2/workspace.yaml" "canonical_url: https://example.com/svc.git"
contains "$ws2/repositories.local.yaml" "anchor_branch: development"
cc_repo_resolve "$ws2" svc >/dev/null
# credential-bearing URL is rejected; duplicate binding is rejected
expect_failure cc_repository_register "$ws2" bad repositories/svc main https://user:pw@example.com/x.git
expect_failure cc_repository_register "$ws2" svc repositories/svc development
rm -rf "$ws2"

# --- regression: register when workspace.yaml has content AFTER `repositories:` ---
# The shipped seed places an `identity:` block after `repositories:`. An EOF append
# would drop the new id into the wrong section, where cc_plan_repositories cannot
# see it (breaking plan validation). Registration must insert under `repositories:`
# and preserve trailing blocks.
ws3=$(mktemp -d "${TMPDIR:-/tmp}/cc-reg2.XXXXXX")
printf 'schema_version: 1\nworkspace: seed-ws\nrepositories: []\n\nidentity:\n  kind: instantiated-workspace\n  status: uninitialized\n' >"$ws3/workspace.yaml"
cc_workspace_init "$ws3" >/dev/null
cc_repository_register "$ws3" alpha repositories/alpha main >/dev/null
cc_repository_register "$ws3" beta repositories/beta develop >/dev/null
regids=$(cc_plan_repositories "$ws3/workspace.yaml" | sort | tr '\n' ' ')
[ "$regids" = "alpha beta " ] || fail "identity not recorded under repositories: (got: '$regids')"
contains "$ws3/workspace.yaml" "status: uninitialized"   # trailing identity block preserved
rm -rf "$ws3"

pass 'repositories'
