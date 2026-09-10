#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

# --- safe relative paths reject traversal, absolute, and space tricks ---
cc_safe_relative "src/billing" || fail "safe path rejected"
expect_failure cc_safe_relative "/etc/passwd"
expect_failure cc_safe_relative "../escape"
expect_failure cc_safe_relative "a/../b"
expect_failure cc_safe_relative "with space"

# --- safe identifiers ---
cc_safe_id "0001-billing-v2" || fail "safe id rejected"
expect_failure cc_safe_id "../evil"
expect_failure cc_safe_id "-leading"
expect_failure cc_safe_id "with space"

# --- plan id form ---
cc_plan_id_valid "0001-billing-v2" || fail "valid plan id rejected"
expect_failure cc_plan_id_valid "billing"
expect_failure cc_plan_id_valid "0001-Billing"

# --- unsafe symlink repository path is rejected ---
ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
cc_fx_repo "$ws" api development
ln -s "$ws/repositories/api" "$ws/repositories/linky"
cc_fx_bindings_header "$ws"
printf '  linky:\n    path: repositories/linky\n    base_branch: development\n' >>"$ws/repositories.local.yaml"
expect_failure cc_repo_resolve "$ws" linky

# --- credentials never appear in runtime or other shipped files ---
credential_scan_paths="$ROOT/.context-circuit/wrapper $ROOT/template $ROOT/.agents $ROOT/.context-circuit/agents"
for term in password api_key access_token client_secret provider_payload BEGIN\ RSA\ PRIVATE\ KEY; do
	if grep -RIlF "$term" $credential_scan_paths 2>/dev/null | grep -v '/test/' | grep . ; then
		fail "credential-like term '$term' found in shipped files"
	fi
done

# Overlay content belongs only in the worktree; runtime records carry metadata,
# never a copied path or its content.
ws_overlay=$(cc_fx_ws)
cc_fx_repo "$ws_overlay" api development
cc_fx_ignored_content "$ws_overlay/repositories/api"
rm -f "$ws_overlay/repositories/api/untracked.txt"
cc_fx_plan_ex "$ws_overlay" 0001-overlay-security "Overlay security" api src ""
exec_overlay=$(cc_execution_begin "$ws_overlay" 0001-overlay-security security-worker | sed -n 's/^execution_id: //p')
overlay_records="$ws_overlay/.runtime/executions/0001-overlay-security/$exec_overlay"
if grep -RIlF 'fixture-overlay-value' "$overlay_records" | grep . >/dev/null 2>&1; then
	fail 'overlay content was written to a runtime record'
fi
if grep -RIlF '.env' "$overlay_records" | grep . >/dev/null 2>&1; then
	fail 'overlay path was written to a runtime record'
fi
rm -rf "$ws_overlay"

# --- the template ignores local bindings, checkouts, and runtime state ---
contains "$ROOT/template/.gitignore" "repositories.local.yaml"
contains "$ROOT/template/.gitignore" "/member.local.yaml"
contains "$ROOT/template/.gitignore" "/repositories/"
contains "$ROOT/template/.gitignore" "/.runtime/"
contains "$ROOT/template/.gitignore" ".code-review-graph"

# --- leftover .code-review-graph is dropped on workspace-init ---
ws_crg=$(cc_fx_ws)
mkdir -p "$ws_crg/.code-review-graph"
printf 'junk\n' >"$ws_crg/.code-review-graph/graph.db"
cc_workspace_init "$ws_crg" >/dev/null
test ! -e "$ws_crg/.code-review-graph" || fail "workspace-init kept leftover .code-review-graph"
rm -rf "$ws_crg"

# --- negative fixtures document forbidden record fields ---
require_file "$ROOT/test/security/fixtures/negatives/forbidden-fields.yaml"
require_file "$ROOT/test/security/fixtures/negatives/provider-payload-record.yaml"

pass 'security boundaries'
