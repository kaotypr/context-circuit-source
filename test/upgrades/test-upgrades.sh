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
for fixture in lease.yaml worktree-state.txt handoff.md plan.yaml stack.yaml; do
  require_file "$ROOT/test/upgrades/legacy-live/$fixture"
done
contains "$ROOT/test/upgrades/legacy-live/worktree-state.txt" 'preserve: do-not-discard'
contains "$ROOT/test/upgrades/legacy-live/plan.yaml" 'status: approved'
contains "$ROOT/test/upgrades/legacy-live/stack.yaml" 'status: interrupted'
legacy=$(mktemp -d "${TMPDIR:-/tmp}/cc-identity-upgrade.XXXXXX")
trap 'rm -rf "$legacy" "$upgrade_receipt"' EXIT HUP INT TERM
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
pass 'compatible, migration-needed, blocked, legacy-unknown, and rollback boundaries'
