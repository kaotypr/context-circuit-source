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

fx_evid="$ROOT/test/runtime/fixtures/evidence-layers"
require_file "$fx_evid/schema-versus-browser.yaml"
require_file "$fx_evid/build-versus-interaction.yaml"
require_file "$fx_evid/sufficient-process.yaml"
require_file "$fx_evid/unavailable-host.yaml"
require_file "$fx_evid/waiver.yaml"
require_file "$fx_evid/writer-claim.yaml"
require_file "$fx_evid/missing-observed.yaml"
require_file "$fx_evid/verifier-read-only.yaml"
require_file "$fx_evid/verifier-repair.yaml"

assert_eq "$(cc_compare_evidence_layer process process passed independent)" passed
assert_eq "$(cc_compare_evidence_layer browser schema passed independent)" failed
assert_eq "$(cc_compare_evidence_layer human process passed independent)" failed
assert_eq "$(cc_compare_evidence_layer store '' passed independent)" failed
assert_eq "$(cc_compare_evidence_layer browser omitted blocked independent unavailable)" blocked
assert_eq "$(cc_compare_evidence_layer human human waived independent)" waived
assert_eq "$(cc_compare_evidence_layer process process passed writer-claim)" failed

assert_eq "$(cc_validate_evidence_mapping "$fx_evid/sufficient-process.yaml")" EVIDENCE_LAYER_MATCH
schema_vs_browser=$(cc_validate_evidence_mapping "$fx_evid/schema-versus-browser.yaml" 2>&1 || true)
printf '%s\n' "$schema_vs_browser" | grep -F EVIDENCE_LAYER_MISMATCH >/dev/null || fail 'schema evidence was accepted for a browser criterion'
build_vs_human=$(cc_validate_evidence_mapping "$fx_evid/build-versus-interaction.yaml" 2>&1 || true)
printf '%s\n' "$build_vs_human" | grep -F EVIDENCE_LAYER_MISMATCH >/dev/null || fail 'process smoke was accepted for a human criterion'
missing=$(cc_validate_evidence_mapping "$fx_evid/missing-observed.yaml" 2>&1 || true)
printf '%s\n' "$missing" | grep -F EVIDENCE_LAYER_MISSING >/dev/null || fail 'missing observed layer was not rejected'
blocked=$(cc_validate_evidence_mapping "$fx_evid/unavailable-host.yaml" 2>&1 || true)
printf '%s\n' "$blocked" | grep -F EVIDENCE_LAYER_BLOCKED >/dev/null || fail 'unavailable host capability was not blocked'
printf '%s\n' "$blocked" | grep -F EVIDENCE_LAYER_MATCH >/dev/null && fail 'unavailable host capability became a pass'
waived=$(cc_validate_evidence_mapping "$fx_evid/waiver.yaml" 2>&1 || true)
printf '%s\n' "$waived" | grep -F EVIDENCE_LAYER_WAIVED >/dev/null || fail 'waiver was not reported as waived'
printf '%s\n' "$waived" | grep -F EVIDENCE_LAYER_MATCH >/dev/null && fail 'waiver counted as verifier success'
writer_claim=$(cc_validate_evidence_mapping "$fx_evid/writer-claim.yaml" 2>&1 || true)
printf '%s\n' "$writer_claim" | grep -F EVIDENCE_WRITER_CLAIM >/dev/null || fail 'writer claim substituted for independent evidence'
assert_eq "$(cc_validate_verifier_readonly "$fx_evid/verifier-read-only.yaml")" EVIDENCE_VERIFIER_READONLY
repair=$(cc_validate_verifier_readonly "$fx_evid/verifier-repair.yaml" 2>&1 || true)
printf '%s\n' "$repair" | grep -E 'EVIDENCE_VERIFIER_WRITE|EVIDENCE_VERIFIER_REPAIR' >/dev/null || fail 'verifier repair packet was accepted'
expect_failure cc_validate_delegation "$fx_evid/verifier-repair.yaml"

layer_request="$runtime/layer-request.yaml"
cat > "$layer_request" <<EOF
transaction_id: layer-1
session_id: writer-1
parent_session_id: root
root_session_id: root
kind: child
role: writer
objective: carry evidence layers
scope: VEL-002
non_goals: [provider]
context_set: writer
repository: app
worktree: .runtime/worktrees/app/checkout
acceptance_criteria: [VEL-AC-02]
verification: [VEL-VT-02]
expected_evidence: [layers]
stop_conditions: [mismatch]
plan: graph-plan
route_decision_digest: sha256:graph-route
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
required_layer: process
produced_layer: process
observed_layer: process
outcome: passed
evidence_ref: sufficient-process
evidence_source: independent
EOF
cc_construct_delegation "$runtime/layer-delegation.yaml" "$layer_request" layer-1
contains "$runtime/layer-delegation.yaml" 'required_layer: process'
contains "$runtime/layer-delegation.yaml" 'observed_layer: process'
cc_validate_delegation "$runtime/layer-delegation.yaml" >/dev/null
cc_construct_handoff "$runtime/layer-handoff.md" "$layer_request" layer-1
contains "$runtime/layer-handoff.md" 'required_layer: process'
contains "$runtime/layer-handoff.md" 'observed_layer: process'
contains "$runtime/layer-handoff.md" 'Writer-recorded mappings are claims'
cc_validate_handoff_evidence "$runtime/layer-handoff.md" >/dev/null
cc_construct_completion "$runtime/layer-completion.yaml" "$layer_request" layer-1
contains "$runtime/layer-completion.yaml" 'required_layer: process'
contains "$runtime/layer-completion.yaml" 'status: not-ready'
cc_validate_completion_evidence "$runtime/layer-completion.yaml" >/dev/null

mismatch_request="$runtime/mismatch-request.yaml"
sed 's/observed_layer: process/observed_layer: schema/; s/produced_layer: process/produced_layer: schema/' "$layer_request" > "$mismatch_request"
cc_construct_delegation "$runtime/mismatch-delegation.yaml" "$mismatch_request" layer-mismatch
mismatch_del=$(cc_validate_delegation "$runtime/mismatch-delegation.yaml" 2>&1 || true)
printf '%s\n' "$mismatch_del" | grep -F EVIDENCE_LAYER_MISMATCH >/dev/null || fail 'delegation accepted schema evidence for a process criterion'
cc_construct_handoff "$runtime/mismatch-handoff.md" "$mismatch_request" layer-mismatch
mismatch_hand=$(cc_validate_handoff_evidence "$runtime/mismatch-handoff.md" 2>&1 || true)
printf '%s\n' "$mismatch_hand" | grep -F EVIDENCE_LAYER_MISMATCH >/dev/null || fail 'handoff accepted non-matching observed evidence'
cc_construct_completion "$runtime/mismatch-completion.yaml" "$mismatch_request" layer-mismatch
mismatch_comp=$(cc_validate_completion_evidence "$runtime/mismatch-completion.yaml" 2>&1 || true)
printf '%s\n' "$mismatch_comp" | grep -F EVIDENCE_LAYER_MISMATCH >/dev/null || fail 'completion accepted non-matching observed evidence'

ready_pass="$runtime/ready-pass.yaml"
cc_atomic_write "$ready_pass" \
  'status: ready-for-human-status-change' \
  'evidence_layers:' \
  '  required_layer: process' \
  '  observed_layer: process' \
  '  outcome: passed' \
  '  evidence_source: independent'
cc_validate_completion_evidence "$ready_pass" >/dev/null
ready_waived="$runtime/ready-waived.yaml"
cc_atomic_write "$ready_waived" \
  'status: ready-for-human-status-change' \
  'evidence_layers:' \
  '  required_layer: human' \
  '  observed_layer: human' \
  '  outcome: waived' \
  '  evidence_source: independent'
ready_waived_result=$(cc_validate_completion_evidence "$ready_waived" 2>&1 || true)
printf '%s\n' "$ready_waived_result" | grep -E 'EVIDENCE_LAYER_WAIVED|EVIDENCE_COMPLETION_NOT_READY' >/dev/null || fail 'ready completion accepted a waiver as pass'

contains "$ROOT/.agents/skills/cc-verify/SKILL.md" 'wrapper/contracts/schemas/plan.yaml'
contains "$ROOT/agents/verifier.md" 'wrapper/contracts/schemas/plan.yaml'
contains "$ROOT/docs/runtime-contract.md" 'wrapper/contracts/schemas/plan.yaml'
not_contains "$ROOT/.agents/skills/cc-verify/SKILL.md" 'vocabulary: [schema, store, api, process, browser, human]'
not_contains "$ROOT/agents/verifier.md" 'comparison: exact-match-only'

pass 'receipt validation, atomic handoff, resumable evidence shape, and evidence-layer enforcement'
