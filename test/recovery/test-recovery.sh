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

bootstrap_root="$runtime/bootstrap-workspace"
origin="$runtime/bootstrap-origin"
mkdir -p "$bootstrap_root"
cat > "$bootstrap_root/workspace.yaml" <<'EOF'
version: 1
template_version: 1.0.0
workspace:
  name: bootstrap-fixture
repositories:
  app:
    canonical_url: https://github.com/acme/app.git
    default_branch: main
EOF
git init -q -b main "$origin"
git -C "$origin" config user.email test@example.invalid
git -C "$origin" config user.name 'Context Circuit Test'
printf '%s\n' bootstrap > "$origin/README.md"
git -C "$origin" add README.md
git -C "$origin" commit -qm initial
card=$(cc_bootstrap_repository "$bootstrap_root" app "$origin" main repositories/app declined 2>&1 || true)
printf '%s\n' "$card" | grep -F 'Canonical URL:' >/dev/null || fail 'bootstrap card omitted canonical URL'
printf '%s\n' "$card" | grep -F 'Existing-path check: absent' >/dev/null || fail 'bootstrap card omitted target check'
test ! -e "$bootstrap_root/repositories/app" || fail 'unconfirmed bootstrap created a destination'
if cc_bootstrap_repository "$bootstrap_root" app "$origin" main repositories/app confirmed >/dev/null 2>&1; then :; else fail 'confirmed local bootstrap failed'; fi
test -e "$bootstrap_root/repositories/app/.git" || fail 'confirmed bootstrap did not clone'
contains "$bootstrap_root/.runtime/bootstrap/app.yaml" 'status: cloned'
mkdir -p "$bootstrap_root/repositories/existing"
printf '%s\n' preserve > "$bootstrap_root/repositories/existing/marker.txt"
expect_failure cc_bootstrap_repository "$bootstrap_root" app "$origin" main repositories/existing confirmed
contains "$bootstrap_root/repositories/existing/marker.txt" preserve
if CC_OFFLINE=1 cc_bootstrap_repository "$bootstrap_root" offline "$origin" main repositories/offline confirmed >/dev/null 2>&1; then fail 'offline bootstrap unexpectedly cloned'; fi
contains "$bootstrap_root/.runtime/bootstrap/offline.yaml" 'status: offline'
if rg -n 'password:|api_key:|access_token:|client_secret:' "$bootstrap_root" >/dev/null 2>&1; then fail 'bootstrap fixture persisted credential-shaped data'; fi

identity_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-identity-recovery.XXXXXX")
cp -R "$ROOT/test/contracts/fixtures/identity-projection/matching/." "$identity_root/"
before_ws=$(cat "$identity_root/context/WORKSPACE.md")
before_proj=$(cat "$identity_root/context/PROJECT.md")
before_idx=$(cat "$identity_root/context/INDEX.md")
before_yaml=$(cat "$identity_root/workspace.yaml")
if CC_IDENTITY_PUBLISH_FAIL_AFTER=1 cc_register_repository "$identity_root" app 'https://github.com/acme/app.git' main confirmed >/dev/null 2>&1; then
  fail 'interrupted identity publish unexpectedly succeeded'
fi
after_ws=$(cat "$identity_root/context/WORKSPACE.md")
after_proj=$(cat "$identity_root/context/PROJECT.md")
after_idx=$(cat "$identity_root/context/INDEX.md")
after_yaml=$(cat "$identity_root/workspace.yaml")
test "$before_ws" = "$after_ws" || fail 'interrupted publish mutated WORKSPACE.md'
test "$before_proj" = "$after_proj" || fail 'interrupted publish mutated PROJECT.md'
test "$before_idx" = "$after_idx" || fail 'interrupted publish mutated INDEX.md'
test "$before_yaml" = "$after_yaml" || fail 'interrupted publish mutated workspace.yaml'
test ! -f "$identity_root/.runtime/identity-publish.journal" || fail 'interrupted publish left a journal that could authorize a partial region'
printf 'status: committing\n' > "$identity_root/.runtime/identity-publish.journal"
journal_mismatch=$(cc_validate_identity_projection "$identity_root" || true)
printf '%s\n' "$journal_mismatch" | grep -Fx projection-mismatch >/dev/null || fail 'leftover publish journal did not fail closed'
rm -f "$identity_root/.runtime/identity-publish.journal"
trap 'rm -rf "$runtime" "$identity_root"' EXIT HUP INT TERM
pass 'interruption preservation, legacy receipt classification, and recovery safety'
