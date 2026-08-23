#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/wrapper/migrations/upgrade.sh"

assert_eq "$(cc_upgrade_classify legacy-unknown)" legacy-unknown
assert_eq "$(cc_upgrade_classify 1.0.0)" compatible
assert_eq "$(cc_upgrade_classify 0.9.0)" migration-needed
assert_eq "$(cc_upgrade_classify unknown)" blocked
assert_eq "$(cc_migration_classify 0.9.0)" migration-needed
assert_eq "$(cc_migration_classify 1.0.0)" compatible
assert_eq "$(cc_migration_classify '')" legacy-unknown
printf '%s\n' "$(cc_migration_preserves)" | grep -F '.runtime' >/dev/null || fail 'migration preserve boundary missing'
contains "$ROOT/wrapper/migrations/v1-to-v2.yaml" 'never rewrite status without its current human gate'
contains "$ROOT/docs/migration.md" 'Rollback restores the previous wrapper-owned files'
contains "$ROOT/wrapper/manifest.yaml" 'preserve:'
contains "$ROOT/wrapper/manifest.yaml" 'wrapper_only:'
contains "$ROOT/wrapper/manifest.yaml" 'CLAUDE.md'
printf '%s\n' "$(cc_rollback_scope)" | grep -Fx 'CLAUDE.md' >/dev/null || fail 'Claude adapter missing from rollback scope'
printf '%s\n' "$(cc_migration_preserves)" | grep -Fx 'repositories.local.yaml' >/dev/null || fail 'local binding preserve boundary missing'
printf '%s\n' "$(cc_migration_preserves)" | grep -Fx 'repositories' >/dev/null || fail 'repository preserve boundary missing'
for fixture in $(cc_migration_packet_fixtures); do
  require_file "$ROOT/$fixture"
done
upgrade_receipt=$(mktemp "${TMPDIR:-/tmp}/cc-upgrade-receipt.XXXXXX")
printf 'wrapper_version: 1.0.0\n' > "$upgrade_receipt"
assert_eq "$(cc_migration_receipt_status "$upgrade_receipt")" compatible
contains "$ROOT/wrapper/migrations/README.md" 'compatible receipt is retained'
contains "$ROOT/wrapper/migrations/README.md" 'legacy child-start record is retained'
printf '%s\n' "$(cc_migration_packet_fixtures)" | grep -Fx 'wrapper/contracts/schemas/child-start.yaml' >/dev/null || fail 'child-start migration fixture missing'
legacy_child_start=$(mktemp "${TMPDIR:-/tmp}/cc-legacy-child-start.XXXXXX")
printf 'schema_version: 1\nwrapper_version: 1.0.0\ntransaction_state: legacy-readable\ncommit_marker: null\n' > "$legacy_child_start"
assert_eq "$(cc_migration_receipt_status "$legacy_child_start")" compatible
contains "$legacy_child_start" 'transaction_state: legacy-readable'
contains "$legacy_child_start" 'commit_marker: null'
for fixture in lease.yaml worktree-state.txt handoff.md plan.yaml stack.yaml; do
  require_file "$ROOT/test/upgrades/legacy-live/$fixture"
done
contains "$ROOT/test/upgrades/legacy-live/worktree-state.txt" 'preserve: do-not-discard'
contains "$ROOT/test/upgrades/legacy-live/plan.yaml" 'status: approved'
contains "$ROOT/test/upgrades/legacy-live/stack.yaml" 'status: interrupted'
legacy=$(mktemp -d "${TMPDIR:-/tmp}/cc-identity-upgrade.XXXXXX")
trap 'rm -rf "$legacy" "$upgrade_receipt"; rm -f "$legacy_child_start"' EXIT HUP INT TERM
cp -R "$ROOT/test/contracts/fixtures/identity-projection/missing-region/." "$legacy/"
contains "$legacy/workspace.yaml" 'status: accepted'
not_contains "$legacy/context/WORKSPACE.md" 'context-circuit:identity-region:start'
cc_migrate_identity_regions "$legacy"
contains "$legacy/workspace.yaml" 'status: accepted'
contains "$legacy/context/WORKSPACE.md" 'Legacy authored Product Knowledge without an identity region.'
contains "$legacy/context/PROJECT.md" 'Legacy authored Product Knowledge without an identity region.'
contains "$legacy/context/INDEX.md" 'Legacy authored Product Knowledge without an identity region.'
contains "$legacy/context/WORKSPACE.md" 'name: missing-region-workspace'
assert_eq "$(cc_validate_identity_projection "$legacy")" identity-projection-ok
cc_rollback_identity_regions "$legacy"
contains "$legacy/workspace.yaml" 'status: accepted'
not_contains "$legacy/context/WORKSPACE.md" 'context-circuit:identity-region:start'
contains "$legacy/context/WORKSPACE.md" 'Legacy authored Product Knowledge without an identity region.'
contains "$ROOT/docs/configuration.md" 'workspace.yaml` is the workspace identifier'
contains "$ROOT/docs/getting-started.md" 'proposed defaults'
contains "$ROOT/wrapper/migrations/README.md" 'does not change accepted'
contains "$ROOT/wrapper/migrations/README.md" 'Completed historical evidence remains readable without rewrite'
contains "$ROOT/wrapper/migrations/README.md" 'explicit evidence-layer mapping'
printf '%s\n' "$(cc_migration_evidence_rules)" | grep -F 'invent-evidence: forbidden' >/dev/null || fail 'migration evidence invent rule missing'

fx_legacy="$ROOT/test/upgrades/fixtures/evidence-layers"
require_file "$fx_legacy/legacy-completed.yaml"
require_file "$fx_legacy/legacy-unfinished.yaml"
require_file "$fx_legacy/mapped-unfinished.yaml"
assert_eq "$(cc_legacy_evidence_classify "$fx_legacy/legacy-completed.yaml")" legacy-completed-readable
assert_eq "$(cc_require_explicit_evidence_mapping "$fx_legacy/legacy-completed.yaml")" LEGACY_COMPLETED_READABLE
unfinished=$(cc_require_explicit_evidence_mapping "$fx_legacy/legacy-unfinished.yaml" 2>&1 || true)
printf '%s\n' "$unfinished" | grep -F EVIDENCE_MAPPING_REQUIRED >/dev/null || fail 'unfinished legacy work skipped explicit mapping'
assert_eq "$(cc_legacy_evidence_classify "$fx_legacy/legacy-unfinished.yaml")" mapping-required
assert_eq "$(cc_require_explicit_evidence_mapping "$fx_legacy/mapped-unfinished.yaml")" EVIDENCE_MAPPING_PRESENT

round_trip=$(mktemp -d "${TMPDIR:-/tmp}/cc-legacy-evidence.XXXXXX")
trap 'rm -rf "$legacy" "$upgrade_receipt" "$round_trip"; rm -f "$legacy_child_start"' EXIT HUP INT TERM
cc_migrate_legacy_evidence "$fx_legacy/legacy-completed.yaml" "$round_trip/completed.yaml"
cmp -s "$fx_legacy/legacy-completed.yaml" "$round_trip/completed.yaml" || fail 'completed historical evidence was rewritten'
assert_eq "$(cc_legacy_evidence_invented "$fx_legacy/legacy-completed.yaml" "$round_trip/completed.yaml")" EVIDENCE_NOT_INVENTED
cc_migrate_legacy_evidence "$fx_legacy/legacy-unfinished.yaml" "$round_trip/unfinished.yaml"
assert_eq "$(cc_legacy_evidence_invented "$fx_legacy/legacy-unfinished.yaml" "$round_trip/unfinished.yaml")" EVIDENCE_NOT_INVENTED
contains "$round_trip/unfinished.yaml" 'required_layer: omitted'
not_contains "$round_trip/completed.yaml" 'required_layer: process'
pass 'compatible, migration-needed, blocked, legacy-unknown, and rollback boundaries'
