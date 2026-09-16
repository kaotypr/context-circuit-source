---
name: cc-direct
description: Make a change directly in a bound checkout when a person chooses to bypass intent and planning — what their request authorizes, where the work happens, when to stop and offer the intent path, and how the knowledge it changed is reconciled.
---

# Change it directly

A person has chosen to bypass the intent and plan path for this change. Work in
the bound checkout, on the branch it is already on, and keep the change to what
they asked for.

## What the request authorizes

The request is the approval. The intent gate exists because an agent's reading of
a request can drift from what the person meant, so the outcome is written down
and approved before code is investigated. Here the outcome is what they said: no
reading was derived, so there is nothing derived to approve. That is the whole of
the bypass, and it holds only because a person asked for it.

It authorizes changing what the request describes and nothing beyond it. No
record is written, no intent or plan ID is allocated, and nothing is reserved or
archived. Every other gate stands exactly as it does elsewhere: committing in a
bound checkout, pushing, opening a pull request, merging, deploying, publishing,
and deleting each still need explicit authorization.

## Where the work happens

The bound checkout this machine records, on its current branch. No worktree is
prepared and no branch is created, because isolation is what the person declined.

- Never checkout, reset, stash, or force anything to make the change fit. Work
  with the tree as it is.
- Preserve everything already there — staged, dirty, untracked, and ignored.
  Uncommitted work belonging to someone else is not yours to tidy.
- Leave the change uncommitted unless the person authorizes a commit. The diff
  in front of them is what a plan record and a branch would otherwise have been.
- The repository's own `AGENTS.md` or `CLAUDE.md` applies here exactly as it does
  in a prepared worktree; follow it where it is more specific.
- Run the checks the change deserves, and report what was actually run and what
  it actually returned.

## When to stop and offer the intent path

The lane fits a change whose outcome is already settled. Stop as soon as it is
not, say what was found, and offer to write an intent instead:

- the change reaches a second repository, or has to be coordinated with work
  already planned;
- settling it needs a decision the person has not made;
- the fix reveals the real problem is somewhere else;
- the reading needed to make the change safely is large enough to have been a
  plan.

Leave what is already correct in place, and say plainly what is done and what is
not. A change that grew past this lane and was carried through anyway is an
unrecorded feature nobody approved.

## Reconcile the knowledge it changed

A completion names this work for a plan. A change made without a plan has no
record to name it, so ask the repository:

```sh
context-circuit-cli --workspace <root> context find --repo REPO
```

It returns the catalog entries claiming that repository and names the judging as
`reconcile_required`. Judge each against what the change actually altered: edit
the note and its catalog entry together and move its reviewed date, or say in
your report that it changed nothing. Most direct changes change nothing durable,
and saying so is the normal outcome rather than a skipped step.

Use `.agents/skills/cc-knowledge/SKILL.md` for the writing itself. This is the
one trace a direct change leaves besides the diff, which is why it is not
optional.
