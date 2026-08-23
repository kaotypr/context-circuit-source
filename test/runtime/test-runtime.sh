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

interrupted_graph="$runtime/.transactions/interrupted"
mkdir -p "$interrupted_graph"
cc_atomic_write "$interrupted_graph/transaction.yaml" \
  'schema_version: 1' 'transaction_id: interrupted' 'transaction_state: valid' 'commit_marker: null'
expect_failure cc_runtime_graph_authoritative "$interrupted_graph"
test ! -e "$interrupted_graph/commit.marker" || fail 'interrupted graph published a commit marker'

graph_request="$runtime/graph-request.yaml"
graph_lease="$runtime/graph-lease.yaml"
cat > "$graph_lease" <<'EOF'
schema_version: 1
wrapper_version: 1.0.0
plan: graph-plan
session_id: writer-1
root_session_id: root
repository: app
worktree: .runtime/worktrees/app/checkout
branch: main
status: active
EOF
cat > "$graph_request" <<EOF
transaction_id: graph-1
session_id: writer-1
parent_session_id: root
root_session_id: root
kind: child
role: writer
objective: bounded graph fixture
scope: ERR-004
non_goals: [provider]
context_set: writer
repository: app
worktree: .runtime/worktrees/app/checkout
acceptance_criteria: [ERR-AC-02]
verification: [ERR-VT-02]
expected_evidence: [graph]
stop_conditions: [mismatch]
plan: graph-plan
route_decision_digest: sha256:graph-route
lease_source: $graph_lease
host_id: codex
observed_version: unavailable
instruction_surface: AGENTS.md
session_role: writer
root_capability: available
child_capability: available
verifier_capability: available
resume_capability: available
permission_mode: bounded-write
provider_status: enabled
offline_fallback: filesystem-only
live_smoke_status: unavailable
EOF
graph=$(cc_construct_runtime_graph "$runtime" "$graph_request")
cc_runtime_graph_authoritative "$graph" >/dev/null
contains "$graph/records/delegation.yaml" 'delegation_id:'
contains "$graph/records/handoff.md" 'handoff_id:'
contains "$graph/records/completion.yaml" 'completion_id:'
contains "$graph/records/child-start.yaml" 'child_start_id:'

tampered_digest="$runtime/tampered-digest"
cp -R "$graph" "$tampered_digest"
sed -i '0,/^record_digest: /s//record_digest: sha256:tampered/' "$tampered_digest/records/session.yaml"
expect_failure cc_runtime_graph_authoritative "$tampered_digest"

marker_mismatch="$runtime/marker-mismatch"
cp -R "$graph" "$marker_mismatch"
sed -i 's/^transaction_state: committed$/transaction_state: valid/' "$marker_mismatch/transaction.yaml"
expect_failure cc_runtime_graph_authoritative "$marker_mismatch"

stale_non_child="$runtime/stale-non-child"
cp -R "$graph" "$stale_non_child"
sed -i 's/^plan: graph-plan$/plan: changed-plan/' "$stale_non_child/records/completion.yaml"
expect_failure cc_runtime_graph_authoritative "$stale_non_child"

lease_mismatch="$runtime/lease-mismatch"
cp -R "$graph" "$lease_mismatch"
sed -i 's/^repository: app$/repository: foreign/' "$lease_mismatch/records/lease.yaml"
expect_failure cc_runtime_graph_authoritative "$lease_mismatch"

receipt_mismatch="$runtime/receipt-mismatch"
cp -R "$graph" "$receipt_mismatch"
sed -i 's/^context_set: writer$/context_set: verifier/' "$receipt_mismatch/records/context-receipt.yaml"
expect_failure cc_runtime_graph_authoritative "$receipt_mismatch"

session_worktree_backup="$runtime/session-worktree.yaml"
cp "$graph/records/session.yaml" "$session_worktree_backup"
sed -i 's#^worktree: .runtime/worktrees/app/checkout$#worktree: .runtime/worktrees/foreign/checkout#' "$graph/records/session.yaml"
session_worktree_result=$(cc_runtime_graph_authoritative "$graph" 2>&1 || true)
mv "$session_worktree_backup" "$graph/records/session.yaml"
printf '%s\n' "$session_worktree_result" | grep -F OWNERSHIP_WORKTREE_MISMATCH >/dev/null || fail 'session worktree mismatch did not fail with stable ownership result'

receipt_packet_backup="$runtime/receipt-packet.yaml"
cp "$graph/records/context-receipt.yaml" "$receipt_packet_backup"
sed -i 's/^packet_id: writer$/packet_id: verifier/' "$graph/records/context-receipt.yaml"
receipt_packet_result=$(cc_runtime_graph_authoritative "$graph" 2>&1 || true)
mv "$receipt_packet_backup" "$graph/records/context-receipt.yaml"
printf '%s\n' "$receipt_packet_result" | grep -F RECEIPT_OWNERSHIP_MISMATCH >/dev/null || fail 'receipt packet identity mismatch did not fail closed'
pass 'receipt validation, atomic handoff, and resumable evidence shape'
