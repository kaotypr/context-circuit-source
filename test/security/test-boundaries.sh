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

# --- credentials never appear in shipped or workspace-owned files ---
for term in password api_key access_token client_secret provider_payload BEGIN\ RSA\ PRIVATE\ KEY; do
	if grep -RIlF "$term" "$ROOT/wrapper" "$ROOT/template" "$ROOT/.agents" "$ROOT/agents" 2>/dev/null | grep -v '/test/' | grep . ; then
		fail "credential-like term '$term' found in shipped files"
	fi
done

# --- the template ignores local bindings, checkouts, and runtime state ---
contains "$ROOT/template/.gitignore" "repositories.local.yaml"
contains "$ROOT/template/.gitignore" "/repositories/"
contains "$ROOT/template/.gitignore" "/.runtime/"

# --- negative fixtures document forbidden record fields ---
require_file "$ROOT/test/security/fixtures/negatives/forbidden-fields.yaml"
require_file "$ROOT/test/security/fixtures/negatives/provider-payload-record.yaml"

pass 'security boundaries'
