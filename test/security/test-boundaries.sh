#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

for path in ../x /absolute ./bad 'path with space' ''; do
  expect_failure cc_safe_relative "$path"
done
cc_safe_relative context/PROJECT.md
cc_safe_id safe-plan-01
expect_failure cc_safe_id '../unsafe'
expect_failure cc_safe_id 'unsafe id'
not_contains "$ROOT/workspace.yaml" 'password:'
not_contains "$ROOT/workspace.yaml" 'api_key:'
not_contains "$ROOT/template/workspace.yaml" 'secret:'
if rg -n '^[[:space:]]*(password|api_key|access_token|client_secret):' "$ROOT/wrapper" "$ROOT/template" "$ROOT/workspace.yaml" >/dev/null 2>&1; then
  fail 'credential-shaped value stored in workspace or wrapper data'
fi
contains "$ROOT/wrapper/contracts/invariants.yaml" 'sources/ is passive'
contains "$ROOT/wrapper/contracts/context-sets.yaml" 'broad sources scan'
pass 'path traversal, credential boundary, and request-scoped source boundary'
