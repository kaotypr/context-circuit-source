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
  `gh pr view`, `gh pr diff`, and `gh pr list`. `gh` runs as the maintainer's
  own login; never read or handle a token. `git fetch --tags origin` is
  allowed: it updates tags and remote-tracking refs, never a tracked file.
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
`git fetch --tags origin`, so the tags and `origin/main` are current. Run the
commands below under `sh`: zsh does not split an unquoted `$(...)` or variable,
and the loop and pathspecs then come back empty without an error.

Name two revisions explicitly and read from them, never from whatever the
checkout holds:

- `HEAD_REV`, what is under review. For a pull request number,
  `git fetch origin pull/<N>/head` and use `FETCH_HEAD`'s commit; for a branch,
  that branch.
- `BASE_REV`, what it merges into, as fetched:
  `origin/<base>`, where `<base>` is the pull request's `baseRefName`
  (`gh pr view <N> --json baseRefName`), or the base the maintainer names.

1. **The previous stable tag.** The highest `cli-v<X.Y.Z>` with no pre-release
   label that is an ancestor of `HEAD_REV`:

   ```sh
   git tag -l 'cli-v*' --sort=-v:refname | while read -r t; do
     case "${t#cli-v}" in *-*) continue ;; esac
     if git merge-base --is-ancestor "$t" "$HEAD_REV"; then echo "$t"; break; fi
   done
   ```

   An empty result is a failure, not "no previous release": stop and say so.

   For the template, the previous version is `git show <tag>:VERSION`. Call the
   tag `PREV` below.
2. **The unreleased change files.** Every `.changelog/*.md` except
   `README.md`, at both revisions:
   `git ls-tree --name-only <rev> .changelog/`, then
   `git show <rev>:.changelog/<file>`. Take the union: a file only at
   `BASE_REV` merged after the branch was cut and still counts; a file only at
   `HEAD_REV` is the branch's own, or one the branch deleted from the base.
   Read each one's frontmatter (`product`, `kind`, `upgrade`, `authors`, `pr`)
   and body, at the revision that holds it, and note where the two differ.
3. **The pull requests behind each file.** The numbers in its `pr:` field, or
   else the pull request whose merge added it:

   ```sh
   git log --first-parent --diff-filter=A --format='%H %s' "$BASE_REV" -- .changelog/<file>
   ```

   The merge subject reads `Merge pull request #N …`; for a commit that is not a
   merge, `gh pr list --state merged --search <sha>` finds it. A file this finds
   nothing for is not on `BASE_REV` yet: it belongs to the branch under
   review. Also read
   `git log --format='%h %s' "$BASE_REV" -- .changelog/<file>`: a file already edited by a
   later pull request has been reconciled once, and that pull request belongs to
   it too.
4. **What each of those pull requests changed.**
   `gh pr diff <N> --name-only`, and `gh pr diff <N>` when a finding needs the
   detail.

The product paths the Change files check uses decide which paths are visible to
users. Read them from the header of `scripts/check-change-files.sh` rather than
restating them here. That classification is the check's, not the user's view:
`assets.go` embeds `product/` and `template/` into the CLI, so `init` from a new
CLI delivers template changes too. For a CLI release, the user-visible paths are
the CLI paths plus `product/` and `template/`. When an entry's `product` leaves
out a release whose users will see the change, report it as an observation for
the maintainer.

## Reviewing a pull request

The input is a pull request number, or a branch to compare with its base:
`main` unless the branch was cut from another one, which the maintainer names.
List its files with `gh pr diff <N> --name-only`, or
`git diff --name-only "$BASE_REV...$HEAD_REV"`, and its change files among
them.

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

The input is the release-prep branch or its pull request, whose base is the
release branch `release/X.Y`, and the product and version being drafted. Review
the release-prep head, not the release branch: the request that will publish is
on that head, and reaches `release/X.Y` only when the pull request merges. Read
the drafted request at `HEAD_REV`,
`release/requests/<cli|template>/<version>.md`, or, before one is written, the
draft from `sh scripts/release-notes.sh <cli|template> <version> --stdout`,
which changes nothing and reads the checkout, so run it with `HEAD_REV` checked
out.

A stable draft consumes its change files: it removes them with `git rm` and
rewrites a `both` file to name the other product. Recover them from `BASE_REV`
with `git diff --name-status "$BASE_REV...$HEAD_REV" -- .changelog/`, then
`git show "$BASE_REV":.changelog/<file>`, and read the originals as step 2
does. Then compare them with the request:

- an entry the request changed, merged, or dropped, against its change file:
  report the edit, and whether it still reads true against `PREV`;
- a change file removed although its product is not this release's, or a
  `both` file removed instead of rewritten, or rewritten to the wrong product;
- a change file for this product with no entry in the request.

Then read the request itself:

- Read `git diff PREV.."$HEAD_REV"` under the release's user-visible paths (for
  a CLI release, including `product/` and `template/`, as above), and for each
  entry ask whether someone using `PREV` will notice it as written.
- **Gone from the branch.** Flag an entry whose change is no longer on the
  branch, with the commit that removed it
  (`git log -S'<name>' PREV.."$HEAD_REV"`).
- **Contradictions.** Flag entries that contradict each other, or that describe
  one feature at two stages of its rework.
- **Fixes that never shipped.** Flag a `fixed` entry whose defect is not in
  `PREV`, as above.
- **Missing entries.** Flag changes in the diff under those paths that
  someone using `PREV` would notice and that no entry describes. Name the
  commits; whether they need an entry is the maintainer's call.
- For a candidate, the change files stay in place and the stable release's
  draft covers them again, so findings on a candidate carry forward.

## Report

Produce one report the maintainer can read in the session or paste into a pull
request as is.

1. Name the branch or pull request reviewed, `PREV`, and the change files read.
2. Group findings by change file, and list edits to the request with the
   change file each came from. Each finding gives:
   - the file;
   - the evidence: the pull requests or commits, and the sentence of the entry
     that no longer holds;
   - the suggested action: keep, edit, merge into another file, or delete,
     with suggested wording where it helps;
   - whether it rests on file overlap alone.
3. At release prep, list the missing entries separately, since they belong to
   no change file: the commits, and why someone using `PREV` would notice.
4. List the files reviewed with nothing to report, so a clean result is
   visible.
5. End with a count: files reviewed, findings by action.

Findings change nothing by themselves: hand the report to the maintainer.
