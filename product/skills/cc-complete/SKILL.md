---
name: cc-complete
description: Mark Context Circuit plans completed and reconcile the durable product knowledge those plans changed, which together finish the request.
---

# Mark a plan completed

Delivery does not mark a plan done. Completion happens when a person explicitly
asks for it, and it is two acts, not one: append a short completion note to each
record, **and** reconcile the durable product knowledge those plans changed. The
request is finished only when both are done.

```sh
sh <root>/.context-circuit/cli --json record complete --id p0001 --text 'User requested completion.'
```

Write the note in the words the person used, including their language. It
records that someone asked for this, so a translation of it is a paraphrase of
evidence rather than a tidier form of it. The knowledge reconciled afterwards is
English regardless.

The command stamps `completed_at` in the plan's frontmatter. Its result carries
`completed`, the ID of the plan it marked, the catalog entries scoped to that plan's
repositories, and `reconcile_required` naming the judging whenever there are any.
In a Git workspace it also carries `commit_recommended`, the plan, intent, and ID
files left uncommitted: ask the person once whether to commit and push them, or
commit only. Their answer authorizes that and nothing more; a refused push is
handled as `AGENTS.md` says under delivery. Completion reporting success is the record written, not the request finished.

## Judge, do not rewrite

The returned entries are a candidate set to judge, never a list to rewrite.

Most completions change no durable knowledge. Recording *that* in the completion
note is the normal outcome, not a skipped step — and it is an assertion about the
code, so make it after reading the entries rather than before.

Where meaning did change, edit the note and its catalog entry in one pass and
move its reviewed date. A renamed or retired identifier is a glossary row. See
`.agents/skills/cc-knowledge/SKILL.md` for what a note may say and how it anchors.

Completing several plans at once reconciles once across the set, not once per
plan.

Keep implementation-specific evidence in the plan, and keep those plans' own
identifiers out of every note reconciled from them. Run the diagnostic afterwards:

```sh
sh <root>/.context-circuit/cli check
```

Other work need not wait for unrelated knowledge updates.

## Offer to archive

When the result carries `archive_offer`, the person has that many finished
records of their own not archived. Once the request is finished, offer to
archive them, and run `record archive` only on a yes; a no changes nothing.
It moves only their own finished records and returns `commit_recommended`,
which is asked about like any other: commit and push, or commit only, on their
authorization.

```sh
sh <root>/.context-circuit/cli --json record archive
```

## What completion is not

It does not remove a worktree, delete a branch, or deliver anything. Worktrees
and branches remain after completion, and each of those is a separate act needing
its own request.

Completing a plan does not mean its work was merged. When another approach
replaces a plan, complete it anyway and say so in the note: the record should
describe what happened, and the alternative is a plan left open forever.

One consequence needs care. `record order` treats any plan carrying
`completed_at` as a satisfied dependency, so a plan that merged nothing still
satisfies whatever depends on it. Move those dependencies to the plan that did
the work.
