#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-runtime.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
mkdir -p "$runtime/sess"
cat > "$runtime/sess/receipt.yaml" <<'EOF'
schema_version: 1
wrapper_version: 1.0.0
session_id: sess-a
route_decision_digest: sha256:test
context_set: run-plan
references: []
invariants: [INV-CTX-03]
created_at: 2026-08-21T00:00:00Z
EOF
cc_validate_receipt "$runtime/sess/receipt.yaml"
cat > "$runtime/sess/delta-receipt.yaml" <<EOF
schema_version: 1
wrapper_version: 1.0.0
session_id: sess-a
route_decision_digest: sha256:test
context_set: resume
references:
  - path: $ROOT/wrapper/manifest.yaml
    revision: sha256:old
invariants: [INV-CTX-03]
created_at: 2026-08-21T00:00:00Z
EOF
cc_receipt_delta "$runtime/sess/delta-receipt.yaml" | grep -F 'RELOAD:' >/dev/null || fail 'changed receipt evidence was not reloaded'
cat > "$runtime/sess/bad-receipt.yaml" <<'EOF'
schema_version: 1
session_id: sess-a
EOF
expect_failure cc_validate_receipt "$runtime/sess/bad-receipt.yaml"

cc_atomic_write "$runtime/sess/handoff.md" \
  'Observed state: writer complete' \
  'Selected route and reason: verification' \
  'Evidence read: plan and worktree' \
  'Action performed or proposed: verify' \
  'Files/state changed: handoff only' \
  'Verification: pending' \
  'Blockers and human decisions: none' \
  'One next safe action: independent verification'
for section in 'Observed state' 'Selected route and reason' 'Evidence read' 'One next safe action'; do
  contains "$runtime/sess/handoff.md" "$section"
done
pass 'receipt validation, atomic handoff, and resumable evidence shape'
