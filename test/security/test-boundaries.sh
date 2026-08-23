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
contains "$ROOT/wrapper/contracts/schemas/child-start.yaml" 'evidence_only: child-start.yaml records validated launch evidence and never authorizes execution'
contains "$ROOT/wrapper/contracts/schemas/child-start.yaml" 'forbidden: [prompt, credentials, tokens, provider-payloads, transcripts, auth-state, permission-grants]'
for fixture in "$ROOT"/test/hosts/fixtures/*.yaml; do
  require_file "$fixture"
  not_contains "$fixture" 'password:'
  not_contains "$fixture" 'api_key:'
  not_contains "$fixture" 'provider_payload:'
  not_contains "$fixture" 'transcript:'
done
if rg -n 'force_is_authorization: true|host_local: required' "$ROOT/test/hosts" >/dev/null 2>&1; then
  fail 'host-local permission evidence can authorize Context Circuit work'
fi
contains "$ROOT/docs/configuration.md" 'host-local'
contains "$ROOT/docs/host-capabilities.md" 'host-blocked'

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

# The live packet boundary rejects an overread before bytes are exposed. The
# verifier packet is small enough for an offline fixture and has no mutable
# workspace or provider evidence.
packet_receipt=$(mktemp "${TMPDIR:-/tmp}/cc-context-receipt.XXXXXX")
packet_session=security-packet-fixture
packet_paths="wrapper/contracts/schemas/delegation.yaml wrapper/contracts/schemas/completion.yaml agents/verifier.md"
cc_context_packet_load "$ROOT" verifier "$packet_receipt" sha256:security "$packet_session" $packet_paths >/dev/null
contains "$packet_receipt" 'packet_id: verifier'
contains "$packet_receipt" 'actual_bytes:'
assert_eq "$(cc_validate_context_receipt "$ROOT" "$packet_receipt" verifier)" context-receipt-ok
expect_failure cc_context_packet_measure "$ROOT" verifier wrapper/runtime/engine.sh
expect_failure cc_context_packet_read "$ROOT" verifier wrapper/runtime/engine.sh

sed -i 's/^packet_id: verifier$/packet_id: wrong-packet/' "$packet_receipt"
expect_failure cc_validate_context_receipt "$ROOT" "$packet_receipt" verifier

packet_fixture=$(mktemp -d "${TMPDIR:-/tmp}/cc-context-overrun.XXXXXX")
trap 'rm -rf "$binding_root" "$packet_fixture" "$packet_receipt"' EXIT HUP INT TERM
mkdir -p "$packet_fixture/wrapper/contracts/schemas" "$packet_fixture/agents"
cp "$ROOT/wrapper/contracts/context-sets.yaml" "$packet_fixture/wrapper/contracts/context-sets.yaml"
cp "$ROOT/wrapper/contracts/schemas/delegation.yaml" "$packet_fixture/wrapper/contracts/schemas/delegation.yaml"
cp "$ROOT/wrapper/contracts/schemas/completion.yaml" "$packet_fixture/wrapper/contracts/schemas/completion.yaml"
cp "$ROOT/agents/verifier.md" "$packet_fixture/agents/verifier.md"
sed -i '/^  - id: verifier$/,/^  - id: resume$/ s/^    budget_bytes: 11264$/    budget_bytes: 1/' "$packet_fixture/wrapper/contracts/context-sets.yaml"
expect_failure cc_context_packet_measure "$packet_fixture" verifier $packet_paths

neg="$ROOT/test/security/fixtures/negatives"
require_file "$neg/provider-payload-record.yaml"
require_file "$neg/forbidden-fields.yaml"
contains "$neg/provider-payload-record.yaml" 'provider_payload:'
contains "$neg/provider-payload-record.yaml" 'transcript:'
contains "$neg/forbidden-fields.yaml" 'ship: never'
for tree in \
  "$ROOT/wrapper/contracts" \
  "$ROOT/docs/host-capabilities.md" \
  "$ROOT/wrapper/migrations/README.md" \
  "$ROOT/wrapper/migrations/upgrade.sh" \
  "$ROOT/test/hosts/fixtures" \
  "$ROOT/test/contracts/fixtures/evidence-layers" \
  "$ROOT/test/runtime/fixtures/evidence-layers"; do
  if rg -n '^[[:space:]]*(password|api_key|access_token|client_secret|provider_payload|transcript):' "$tree" >/dev/null 2>&1; then
    fail "credential or provider payload leaked into $tree"
  fi
done
not_contains "$ROOT/docs/host-capabilities.md" 'provider_payload:'
not_contains "$ROOT/wrapper/contracts/schemas/plan.yaml" 'provider_payload:'
not_contains "$ROOT/wrapper/contracts/schemas/handoff.yaml" 'transcript:'
pass 'path traversal, credential boundary, request-scoped source, and live packet overread boundaries'
