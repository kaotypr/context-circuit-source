#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-recovery.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
cc_acquire_lease "$runtime" live sess-live sess-live app .runtime/worktrees/app/live main
expect_failure cc_acquire_lease "$runtime" live sess-other sess-other app .runtime/worktrees/app/takeover main
test -d "$runtime/plans/live/lease.lock" || fail 'live lease disappeared during contention'
expect_failure cc_takeover_lease "$runtime" live sess-live sess-other takeover 'not-confirmed'
cc_takeover_lease "$runtime" live sess-live sess-other takeover confirmed
contains "$runtime/plans/live/takeover.yaml" 'previous_owner: sess-live'
assert_eq "$(cc_upgrade_classify '')" legacy-unknown
assert_eq "$(cc_upgrade_classify 1.0.0)" compatible
assert_eq "$(cc_upgrade_classify 0.4.0)" migration-needed
assert_eq "$(cc_upgrade_classify 9.0.0)" blocked
digest=$(cc_digest "$ROOT/wrapper/manifest.yaml")
printf '%s\n' "$digest" | grep -E '^(sha256|cksum):' >/dev/null || fail 'receipt digest missing'
contains "$ROOT/docs/migration.md" 'never rewrites accepted context'
pass 'interruption preservation, legacy receipt classification, and recovery safety'
