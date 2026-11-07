---
name: review-changelog-files
description: Review the unreleased change files in .changelog/ against what will ship, and report which to keep, edit, merge, or delete. Use when reviewing a pull request's change file or Release note, when a pull request may rework or undo an unreleased change, or when preparing a release's notes on a release-prep branch. Maintainer tool; never shipped; edits and posts nothing.
---

# Review change files against what will ship

A change file says what will differ for someone using Context Circuit between
the previous release and the next one, not how the change was built. The rule
is written in `.changelog/README.md` ("Describe the change since the last
release") and in `release/README.md`'s Release notes section: a pull request
that reworks or undoes an unreleased change edits or deletes that file, and the
release-prep edit reads the draft against the last release a user had. This
skill finds where that was missed and reports it. It is maintainer tooling for
the source checkout; it appears in no manifest and ships nowhere.

## Rules

- Report only, when reviewing a pull request and at release prep alike. Edit
  nothing in the checkout: no change file, no request, no draft written to
  disk. The maintainer or the pull request's author makes each change.
- Post nothing to GitHub: no review, comment, label, or status. The report goes
  to the maintainer in the session, who decides what to say on the pull
  request.
- Start no coding host and call no model API.
- Read only this checkout, `git`, and the public pull requests through
  `gh pr view` and `gh pr diff`. `gh` runs as the maintainer's own login; never
  read or handle a token.
- Judge an entry against the previous stable release, never against the merge
  history. The question is always: will someone using that release notice this,
  as written?
- Overlapping files are where to look, not proof. Two pull requests can change
  one file for unrelated reasons, and a revert can live in a file the original
  pull request never touched. Read the diffs before calling an entry stale, and
  say in the report when a finding rests on overlap alone.
- Wording you suggest follows `.changelog/README.md` and the published-text
  rules at the end of `release/README.md`'s Release notes section: written for a
  user, no plan or intent ID, no record path, no private workspace name, no
  AI-assistance credit.

## Set up

Both kinds of review start here. Run from the repository root, after
`git fetch --tags origin`.

1. **The previous stable tag.** The highest `cli-v<X.Y.Z>` with no pre-release
   label that is an ancestor of the branch under review:

   ```sh
   for t in $(git tag -l 'cli-v*' --sort=-v:refname); do
     case "$t" in *-*) continue ;; esac
     git merge-base --is-ancestor "$t" HEAD && { echo "$t"; break; }
   done
   ```

   For the template, the previous version is `git show <tag>:VERSION`. Call the
   tag `PREV` below.
2. **The unreleased change files.** Every `.changelog/*.md` except
   `README.md`. Read each one's frontmatter (`product`, `kind`, `upgrade`,
   `authors`, `pr`) and body.
3. **The pull requests behind each file.** The numbers in its `pr:` field, or
   else the pull request whose merge added it:

   ```sh
   git log --first-parent --diff-filter=A --format='%H %s' origin/main -- .changelog/<file>
   ```

   The merge subject reads `Merge pull request #N …`; for a commit that is not a
   merge, `gh pr list --state merged --search <sha>` finds it. Also read
   `git log --format='%h %s' -- .changelog/<file>`: a file already edited by a
   later pull request has been reconciled once, and that pull request belongs to
   it too.
4. **What each of those pull requests changed.**
   `gh pr diff <N> --name-only`, and `gh pr diff <N>` when a finding needs the
   detail.

The product paths the Change files check uses decide which paths are visible to
users. Read them from the header of `scripts/check-change-files.sh` rather than
restating them here.

## Reviewing a pull request

The input is a pull request number, or a branch to compare with `main`. List
its files with `gh pr diff <N> --name-only`, or
`git diff --name-only main...<branch>`, and its change files among them.

- **Overlap.** For every unreleased change file whose pull requests changed any
  file this one changes, read both diffs. Say whether the entry still describes
  what will ship once this merges. If it does not, say whether to edit it, merge
  it into this pull request's own change file, or delete it, and quote the
  sentence that stops being true.
- **A second file for one feature.** If this pull request adds a change file
  and an unreleased one already describes the same feature, say to merge them
  into one entry that names both pull requests in `pr:` and both authors in
  `authors`, as `.changelog/README.md` asks.
- **A fix for something unreleased.** A `fixed` entry is wrong when the defect
  it fixes is not in `PREV`. Check with `git show PREV:<path>`, or by finding
  the pull request that introduced the code and whether it is an ancestor of
  `PREV`. If it never shipped, the fix needs no file; say which unreleased
  entry, if any, should describe the result instead.
- **Undone changes.** A pull request that removes or reverts code an unreleased
  file describes, without editing or deleting that file, is reported with both
  pull requests. Search the entry's body for the names it uses (commands, flags,
  files, presets) in this pull request's diff as well as the overlapping paths.
- **Release note.** If the pull request touches user-visible paths, it carries a
  change file or says `None` under Release note (`gh pr view <N> --json body`).
  The Change files check already enforces that; only report a `None` the diff
  contradicts.

## Release prep

The input is a release-prep branch, usually `release/X.Y`, and the product and
version being drafted. Read the drafted request,
`release/requests/<cli|template>/<version>.md`, or, before one is written, the
draft from `sh scripts/release-notes.sh <cli|template> <version> --stdout`,
which changes nothing.

- Read `git diff PREV..HEAD` under the product's paths, and for each entry ask
  whether someone using `PREV` will notice it as written.
- **Gone from the branch.** Flag an entry whose change is no longer on the
  branch, with the commit that removed it (`git log -S'<name>' PREV..HEAD`).
- **Contradictions.** Flag entries that contradict each other, or that describe
  one feature at two stages of its rework.
- **Fixes that never shipped.** Flag a `fixed` entry whose defect is not in
  `PREV`, as above.
- **Missing entries.** Flag changes in the diff under the product paths that
  someone using `PREV` would notice and that no entry describes. Name the
  commits; whether they need an entry is the maintainer's call.
- For a candidate, the change files stay in place and the stable release's
  draft covers them again, so findings on a candidate carry forward.

## Report

Produce one report the maintainer can read in the session or paste into a pull
request as is.

1. Name the branch or pull request reviewed, `PREV`, and the change files read.
2. Group findings by change file. Each finding gives:
   - the file;
   - the evidence: the pull requests or commits, and the sentence of the entry
     that no longer holds;
   - the suggested action: keep, edit, merge into another file, or delete,
     with suggested wording where it helps;
   - whether it rests on file overlap alone.
3. List the files reviewed with nothing to report, so a clean result is
   visible.
4. End with a count: files reviewed, findings by action.

Findings change nothing by themselves: hand the report to the maintainer.
