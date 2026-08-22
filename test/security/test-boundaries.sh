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
expect_failure cc_repository_canonical_safe 'https://user:password@example.invalid/app.git'
expect_failure cc_repository_remote_safe 'https://token@example.invalid/app.git'
cc_repository_canonical_safe 'https://github.com/acme/app.git'
not_contains "$ROOT/workspace.yaml" 'password:'
not_contains "$ROOT/workspace.yaml" 'api_key:'
not_contains "$ROOT/template/workspace.yaml" 'secret:'
if rg -n '^[[:space:]]*(password|api_key|access_token|client_secret):' "$ROOT/wrapper" "$ROOT/template" "$ROOT/workspace.yaml" >/dev/null 2>&1; then
  fail 'credential-shaped value stored in workspace or wrapper data'
fi
contains "$ROOT/wrapper/contracts/invariants.yaml" 'sources/ is passive'
contains "$ROOT/wrapper/contracts/context-sets.yaml" 'broad sources scan'

binding_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-binding-security.XXXXXX")
trap 'rm -rf "$binding_root"' EXIT HUP INT TERM
mkdir -p "$binding_root/repositories"
cat > "$binding_root/workspace.yaml" <<'EOF'
version: 1
repositories:
  app:
    canonical_url: https://github.com/acme/app.git
EOF
cat > "$binding_root/repositories.local.yaml" <<'EOF'
repositories:
  app:
    path: ../outside
EOF
expect_failure cc_resolve_repository_binding "$binding_root" app
mkdir -p "$binding_root/real-app"
ln -s "$binding_root/real-app" "$binding_root/repositories/app"
sed -i 's#path: ../outside#path: repositories/app#' "$binding_root/repositories.local.yaml"
expect_failure cc_resolve_repository_binding "$binding_root" app
pass 'path traversal, credential boundary, and request-scoped source boundary'
