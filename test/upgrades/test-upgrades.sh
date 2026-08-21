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
for fixture in lease.yaml worktree-state.txt handoff.md plan.yaml stack.yaml; do
  require_file "$ROOT/test/upgrades/legacy-live/$fixture"
done
contains "$ROOT/test/upgrades/legacy-live/worktree-state.txt" 'preserve: do-not-discard'
contains "$ROOT/test/upgrades/legacy-live/plan.yaml" 'status: approved'
contains "$ROOT/test/upgrades/legacy-live/stack.yaml" 'status: interrupted'
pass 'compatible, migration-needed, blocked, legacy-unknown, and rollback boundaries'
