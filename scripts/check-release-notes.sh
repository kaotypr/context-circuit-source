#!/bin/sh
# Exercise check-change-files.sh and release-notes.sh in a throwaway repository.
# Merge commits carry GitHub's subject, so the pull request links and fallback
# contributors are read the way they are on a real release branch, offline.
set -eu
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

repo="$work/repo"
mkdir -p "$repo/scripts" "$repo/.changelog" "$repo/internal/cli" "$repo/product/skills" "$repo/release/requests"
cp "$source_root/scripts/release-notes.sh" "$source_root/scripts/check-change-files.sh" "$repo/scripts/"
cd "$repo"
git init -q -b main
git config user.name Fixture
git config user.email fixture@example.invalid
git config commit.gpgsign false
printf 'format\n' > .changelog/README.md
printf '2.4.1\n' > VERSION
printf 'package cli\n' > internal/cli/cli.go
printf 'skill\n' > product/skills/skill.md
git add -A
git commit -q -m 'chore: fixture'

# change <file> <product> <kind> [upgrade] [authors]: writes a change file.
change() {
  {
    printf '%s\n' --- "product: $2" "kind: $3"
    [ -z "${4:-}" ] || printf 'upgrade: %s\n' "$4"
    [ -z "${5:-}" ] || printf 'authors: "%s"\n' "$5"
    printf '%s\n' --- '' "Entry for $1."
  } > ".changelog/$1.md"
}
# merge <number> <owner> <branch>: merges the current topic branch into main.
merge() {
  git checkout -q main
  git merge -q --no-ff -m "Merge pull request #$1 from $2/$3" "$3"
}
topic() { git checkout -q -b "$1" main; }
commit() { git add -A; git commit -q -m "$1"; }

# check <case> <expect pass|fail> <base> <head> <body>
check() {
  if PR_BODY=$5 sh scripts/check-change-files.sh "$3" "$4" > "$work/out" 2>&1; then
    [ "$2" = pass ] || fail "change-file check accepted $1: $(cat "$work/out")"
  else
    [ "$2" = fail ] || fail "change-file check refused $1: $(cat "$work/out")"
  fi
}

# The check, on topic branches against main.
topic cli-with-note
printf 'package cli // changed\n' > internal/cli/cli.go
change cli-fix cli fixed
commit 'fix(cli): a fix'
check 'a CLI change with a CLI change file' pass main cli-with-note ''

topic cli-wrong-note
printf 'package cli // other\n' > internal/cli/cli.go
change wrong-product template changed
commit 'fix(cli): a fix with the wrong note'
check 'a CLI change with only a template change file' fail main cli-wrong-note ''
grep -q 'changes cli files' "$work/out" || fail "the refusal did not name the product: $(cat "$work/out")"

topic cli-none
printf 'package cli // none\n' > internal/cli/cli.go
commit 'refactor(cli): nothing users see'
check 'a CLI change without a change file' fail main cli-none '## Release note

Something.'
crlf=$(printf '## What changed\r\n\r\nx\r\n\r\n## Release note\r\n\r\n<!-- One paragraph,\r\nor None. -->\r\nnone.\r\n\r\n## Linked issue\r\n\r\nnone\r\n')
check 'None with CRLF and a comment' pass main cli-none "$crlf"
check 'None inside a comment only' fail main cli-none '## Release note

<!-- None -->'
check 'None under another heading' fail main cli-none '## Linked issue

None'

topic tests-only
printf 'package cli\n' > internal/cli/cli_test.go
commit 'test(cli): a test'
check 'a Go test only' pass main tests-only ''

topic bad-note
printf 'x\n' > product/skills/skill.md
printf '%s\n' --- 'product: template' 'kind: tweak' --- '' 'Body.' > .changelog/bad-kind.md
commit 'docs(template): a change with a bad note'
check 'an unparseable change file' fail main bad-note 'None'
grep -q 'bad-kind.md: unknown kind: tweak' "$work/out" || fail "the refusal did not name the file: $(cat "$work/out")"

topic prep-only
printf '2.5.0\n' > VERSION
mkdir -p release/requests/cli
printf 'request\n' > release/requests/cli/2.5.0.md
git rm -q .changelog/README.md
commit 'chore(release): prepare'
check 'a release-prep change' pass main prep-only ''
git checkout -q main
git branch -q -D cli-wrong-note cli-none tests-only bad-note prep-only

# Drafting, on main after real merges.
merge 11 kaotypr cli-with-note
topic both-break
change both-break both breaking 'Run `agent setup` once.'
commit 'feat!: both'
merge 7 kaotypr both-break
topic template-change
change template-change template changed '' '@hallode @dependabot[bot]'
commit 'docs(template): wording'
merge 12 kaotypr template-change
topic fork-add
change fork-add cli added
commit 'feat(cli): from a fork'
merge 13 contributor fork-add
change direct-add cli added
commit 'feat(cli): committed without a pull request'

sh scripts/release-notes.sh cli 2.5.0-rc.1 > "$work/out" 2> "$work/err" || fail "rc.1 draft failed: $(cat "$work/err")"
grep -q 'direct-add.md reached this branch without a pull request merge' "$work/err" || fail "an unlinked file was not reported: $(cat "$work/err")"
cat > "$work/want-cli" <<'EOF'
---
level: prerelease
reason: One sentence for the person cutting the release; no script reads it.
---

# Context Circuit CLI 2.5.0-rc.1

## Upgrade notes

Run `agent setup` once. (#7)

## Breaking

Entry for both-break. (#7)

## Added

Entry for fork-add. (#13)

Entry for direct-add.

## Fixed

Entry for cli-fix. (#11)

## Contributors

@kaotypr, @contributor
EOF
diff -u "$work/want-cli" release/requests/cli/2.5.0-rc.1.md > "$work/diff" || fail "CLI candidate draft differs:
$(cat "$work/diff")"
[ -f .changelog/cli-fix.md ] || fail 'a candidate removed a change file'

sh scripts/release-notes.sh template 2.5.0-rc.1 > /dev/null 2>&1 || fail 'template candidate draft failed'
grep -qx 'changelog: skip' release/requests/template/2.5.0-rc.1.md || fail 'a template candidate does not skip the changelog'
grep -qx 'Entry for template-change. (kaotypr/context-circuit-source#12)' release/requests/template/2.5.0-rc.1.md ||
  fail "a template entry does not link with the repository: $(cat release/requests/template/2.5.0-rc.1.md)"
grep -qx '@kaotypr, @hallode' release/requests/template/2.5.0-rc.1.md ||
  fail "template contributors are wrong: $(grep -A2 Contributors release/requests/template/2.5.0-rc.1.md)"
if grep -q 'Entry for cli-fix' release/requests/template/2.5.0-rc.1.md; then fail 'a CLI entry reached the template draft'; fi

if sh scripts/release-notes.sh cli 2.5.0-rc.1 > /dev/null 2>&1; then fail 'an existing request was overwritten without --force'; fi
sh scripts/release-notes.sh cli 2.5.0-rc.1 --force > /dev/null 2>&1 || fail '--force did not draft over an existing request'
sh scripts/release-notes.sh cli 2.5.0 --stdout > "$work/stdout" 2>/dev/null || fail '--stdout failed'
[ ! -e release/requests/cli/2.5.0.md ] || fail '--stdout wrote the request'
[ -f .changelog/cli-fix.md ] || fail '--stdout removed a change file'
grep -qx 'level: stable' "$work/stdout" || fail 'a stable draft is not level: stable'

sh scripts/release-notes.sh cli 2.5.0 > /dev/null 2>&1 || fail 'stable CLI draft failed'
for f in cli-fix fork-add direct-add; do
  [ ! -e ".changelog/$f.md" ] || fail "a stable release kept the shipped change file $f.md"
done
grep -qx 'product: template' .changelog/both-break.md || fail "a both file was not rewritten to the other product: $(cat .changelog/both-break.md)"
grep -qx 'kind: breaking' .changelog/both-break.md || fail 'rewriting a both file lost its other fields'
[ -f .changelog/template-change.md ] || fail 'a CLI release removed a template change file'

sh scripts/release-notes.sh template 2.5.0 > /dev/null 2>&1 || fail 'stable template draft failed'
grep -qx 'changelog: keep' release/requests/template/2.5.0.md || fail 'a stable template release does not keep the changelog'
grep -qx 'Entry for both-break. (kaotypr/context-circuit-source#7)' release/requests/template/2.5.0.md ||
  fail 'the rewritten both file did not reach the template draft'
left=$(find .changelog -name '*.md' ! -name README.md)
[ -z "$left" ] || fail "change files remain after both stable releases: $left"

sh scripts/release-notes.sh cli 2.5.1 --stdout > "$work/stdout" 2>/dev/null || fail 'an empty draft failed'
grep -qx 'None.' "$work/stdout" || fail 'an empty draft does not say None. under Upgrade notes'
if grep -q '^## Added' "$work/stdout"; then fail 'an empty draft printed an empty heading'; fi

# A file can name the pull requests it describes: entries added after the fact,
# a change several pull requests made, or a cherry-pick of an earlier fix.
topic named-prs
printf '%s\n' --- 'product: both' 'kind: added' 'pr: 180 177' 'authors: "@hallode"' --- '' 'Entry for named-prs.' > .changelog/named-prs.md
commit 'docs(changelog): describe earlier pull requests'
merge 30 kaotypr named-prs
sh scripts/release-notes.sh cli 2.5.1 --stdout > "$work/stdout" 2> "$work/err" || fail "a draft with pr: failed: $(cat "$work/err")"
grep -qx 'Entry for named-prs. (#180, #177)' "$work/stdout" || fail "pr: did not set the links: $(cat "$work/stdout")"
[ ! -s "$work/err" ] || fail "a file with pr: was reported as unlinked: $(cat "$work/err")"
grep -qx '@hallode' "$work/stdout" || fail 'pr: changed who is credited'
sh scripts/release-notes.sh template 2.5.1 --stdout > "$work/stdout" 2>/dev/null
grep -qx 'Entry for named-prs. (kaotypr/context-circuit-source#180, kaotypr/context-circuit-source#177)' "$work/stdout" ||
  fail "pr: did not set the template links: $(cat "$work/stdout")"
printf '%s\n' --- 'product: cli' 'kind: fixed' 'pr: " "' --- '' 'Entry for blank-pr.' > .changelog/blank-pr.md
git add -A
sh scripts/release-notes.sh cli 2.5.1 --stdout > "$work/stdout" 2> "$work/err" || fail "a blank pr: failed: $(cat "$work/err")"
grep -q 'blank-pr.md reached this branch without a pull request merge' "$work/err" ||
  fail "a blank pr: hid the missing link: $(cat "$work/err")"
git rm -q -f .changelog/blank-pr.md
for bad in abc '12 x' 012; do
  printf '%s\n' --- 'product: cli' 'kind: fixed' "pr: $bad" --- '' 'Body.' > "$work/bad-pr.md"
  if sh scripts/release-notes.sh --validate "$work/bad-pr.md" 2> "$work/err"; then fail "pr: $bad was accepted"; fi
  grep -q 'pr is not a list of pull request numbers' "$work/err" || fail "pr: $bad failed for another reason: $(cat "$work/err")"
done
git rm -q .changelog/named-prs.md
commit 'chore: clear'

# A variant's own label is a release of that variant, not a candidate.
sh scripts/release-notes.sh template 2.5.0-acme.1 --stdout > "$work/stdout" 2>/dev/null || fail 'a variant draft failed'
grep -qx 'changelog: keep' "$work/stdout" || fail 'a variant release was treated as a candidate'

printf 'PASS: change files are required where a product changes, and drafted into requests\n'
