#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

# Deterministic maintainer identity; the pipeline must never author as a bot.
export MAINTAINER_NAME='Release Tester'
export MAINTAINER_EMAIL='release-tester@example.invalid'

work=$(mktemp -d "${TMPDIR:-/tmp}/cc-publish.XXXXXX")
tpl="$work/template"
fix="$work/fixture"
gho="$work/gh_output"
testreq="$ROOT/release/requests/0.0.1-test.1.md"
cleanup() {
  rm -rf "$work"
  rm -f "$testreq"
  rm -f "$ROOT/dist/context-circuit-v0.0.1-alpha.1.tar.gz" \
        "$ROOT/dist/context-circuit-v0.0.1-test.1.tar.gz"
}
trap cleanup EXIT HUP INT TERM

git_as() { GIT_AUTHOR_NAME=$MAINTAINER_NAME GIT_AUTHOR_EMAIL=$MAINTAINER_EMAIL \
           GIT_COMMITTER_NAME=$MAINTAINER_NAME GIT_COMMITTER_EMAIL=$MAINTAINER_EMAIL git "$@"; }

new_template_repo() { d=$1; rm -rf "$d"; mkdir -p "$d"; git -C "$d" init -q; \
  git -C "$d" symbolic-ref HEAD refs/heads/main; }

# ---------------------------------------------------------------------------
# publish-template.sh
# ---------------------------------------------------------------------------
new_template_repo "$tpl"

# A) First release into an empty repo: diff-gate is skipped, one commit + tag.
out=$(sh "$ROOT/scripts/publish-template.sh" 0.0.1-alpha.1 "$tpl")
printf '%s\n' "$out" | grep -F 'first-release' >/dev/null || fail 'first release did not skip the diff-gate'
git -C "$tpl" rev-parse -q --verify refs/tags/v0.0.1-alpha.1 >/dev/null || fail 'first release did not tag v0.0.1-alpha.1'
commits=$(git -C "$tpl" rev-list --count HEAD); assert_eq "$commits" 1
subject=$(git -C "$tpl" log -1 --format='%s'); assert_eq "$subject" 'chore(release): v0.0.1-alpha.1'
require_file "$tpl/.context-circuit/wrapper/manifest.yaml"
contains "$tpl/.context-circuit/wrapper/manifest.yaml" 'template_version: 0.0.1-alpha.1'
require_file "$tpl/CHANGELOG.md"
contains "$tpl/CHANGELOG.md" '## v0.0.1-alpha.1'
test ! -e "$tpl/release" || fail 'release ledger leaked into published tree'

# E) No attribution: author + committer are the maintainer, no bot, no trailer.
assert_eq "$(git -C "$tpl" log -1 --format='%ae')" "$MAINTAINER_EMAIL"
assert_eq "$(git -C "$tpl" log -1 --format='%ce')" "$MAINTAINER_EMAIL"
body=$(git -C "$tpl" log -1 --format='%B')
printf '%s' "$body" | grep -iE 'co-authored-by|generated with|noreply@anthropic' >/dev/null \
  && fail 'release commit carried an attribution trailer' || :

# B) Guard: an already-published version is never re-published.
expect_failure sh "$ROOT/scripts/publish-template.sh" 0.0.1-alpha.1 "$tpl"

# C) Diff-gate: a new version with no source change has nothing to publish.
printf -- '---\nlevel: patch\nreason: test\n---\n\nTest entry.\n' > "$testreq"
expect_failure sh "$ROOT/scripts/publish-template.sh" 0.0.1-test.1 "$tpl"

# D) Real difference: when the latest published tree differs, it publishes.
printf '\ntest drift\n' >> "$tpl/README.md"
git_as -C "$tpl" commit -qam 'chore: simulate drifted prior release'
out=$(sh "$ROOT/scripts/publish-template.sh" 0.0.1-test.1 "$tpl")
git -C "$tpl" rev-parse -q --verify refs/tags/v0.0.1-test.1 >/dev/null || fail 'real-change publish did not tag v0.0.1-test.1'
contains "$tpl/.context-circuit/wrapper/manifest.yaml" 'template_version: 0.0.1-test.1'
# newest-first changelog keeps every prior entry
contains "$tpl/CHANGELOG.md" '## v0.0.1-test.1'
contains "$tpl/CHANGELOG.md" '## v0.0.1-alpha.1'

# ---------------------------------------------------------------------------
# resolve-release-version.sh (hermetic fixture repo)
# ---------------------------------------------------------------------------
mkdir -p "$fix/release/requests"
( cd "$fix" && git init -q && mkdir -p release/requests )
printf 'a\n' > "$fix/release/requests/1.0.0.md"
git_as -C "$fix" add -A && git_as -C "$fix" commit -q -m 'init'
before=$(git -C "$fix" rev-parse HEAD)
printf 'b\n' > "$fix/release/requests/1.1.0.md"
git_as -C "$fix" add -A && git_as -C "$fix" commit -q -m 'add 1.1.0'
sha=$(git -C "$fix" rev-parse HEAD)

# explicit input wins and is written to GITHUB_OUTPUT
: > "$gho"
( cd "$fix" && INPUT_VERSION=1.0.0 GITHUB_OUTPUT="$gho" \
    sh "$ROOT/scripts/resolve-release-version.sh" "$tpl" >/dev/null )
contains "$gho" 'version=1.0.0'

# push resolves the changed request file in the range
: > "$gho"
( cd "$fix" && EVENT_NAME=push BEFORE="$before" SHA="$sha" GITHUB_OUTPUT="$gho" \
    sh "$ROOT/scripts/resolve-release-version.sh" "$tpl" >/dev/null )
contains "$gho" 'version=1.1.0'

# unpublished detection: tag 1.0.0 in the template, leaving exactly one candidate
new_template_repo "$fix/tpl"
git_as -C "$fix/tpl" commit -q --allow-empty -m 'seed'
git_as -C "$fix/tpl" tag -a v1.0.0 -m v1.0.0
: > "$gho"
( cd "$fix" && GITHUB_OUTPUT="$gho" \
    sh "$ROOT/scripts/resolve-release-version.sh" "$fix/tpl" >/dev/null )
contains "$gho" 'version=1.1.0'

# ambiguity (two unpublished, no input) is a hard failure, never a guess
( cd "$fix" && expect_failure sh "$ROOT/scripts/resolve-release-version.sh" "$tpl" )

pass 'template publication and version resolution'
