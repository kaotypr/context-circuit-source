---
name: ff
description: Context Circuit's /ff — fast-forward this machine's base branches to their remotes, and after a person says their merge or pull request merged, offer that together with completing the delivered plans.
---

# Bring base branches up to date

A merged request leaves this machine's base branch behind its remote, and the
next worktree prepared from it starts stale. This skill closes that gap with a
fast-forward, and only a fast-forward. It is named `ff` rather than `cc-ff`
because people type it.

Moving a base branch changes the person's own working copy, so it happens only
on their word: invoking `/ff`, asking for it in words where the host has no
slash commands, or saying yes to the offer below. Never on your own.

## `/ff`

| What the person typed | What you run |
| --- | --- |
| `/ff` | `repo sync`, then `workspace sync` |
| `/ff api web` | `repo sync --id api --id web` |
| `/ff workspace` | `workspace sync` |
| `/ff api workspace` | `repo sync --id api`, then `workspace sync` |

```sh
sh <root>/.context-circuit/cli --json repo sync
sh <root>/.context-circuit/cli --json workspace sync
```

The commands fetch, then move each base branch to its remote counterpart as a
fast-forward. They never switch a checked-out branch, merge, rebase, reset,
stash, or discard, and neither do you: do not pull, rebase, or run raw Git to get
past a status. Every status other than `updated` and `current` left that
repository exactly as it was, and its `detail` names what the person can do.

The workspace repository is held stricter than a code checkout. Any uncommitted
path in it, untracked included, reports `dirty` with the paths and fetches
nothing, while the code checkouts are still synced. Tell the person which paths
kept it back.

## Reporting

One table, one row per repository: its ID, base branch, status, and either the
commits taken (`from`, `to`, `commits`) or why nothing moved. Say when something
was `ahead` — those are commits nobody has delivered yet.

When `workspace sync` moved the workspace it returns `check_recommended`: run
`check` and report what came in, including any finding the pull brought, which
is reported and does not block. When it also returns `cli_version_changed`, the
workspace now pins another CLI release and the launcher will refuse until it is
installed: say so, and point at `.agents/skills/cc-cli/SKILL.md`.

## After a merge

When a person says their delivered work merged — "PR merged", "the MR is in",
"sudah di-merge" — they have not asked for anything yet. Offer, in one reply,
the two things that usually follow, and do neither until they answer.

1. Work out which delivery merged: the `cc/<plan>/<repo>` branches delivered in
   this conversation or named by the person, and the plans those branches carry.
   Where that is unclear, ask which request merged rather than guessing.
2. Offer to sync those repositories' base branches — not the workspace, which a
   merged code change does not move — and to mark those plans completed. Name
   the repositories with their base branch, and the plans by ID and title.
3. One yes covers both. Run `repo sync --id …` for those repositories, then
   complete the plans through `.agents/skills/cc-complete/SKILL.md`, which writes
   the note and reconciles the knowledge they changed. That yes is the explicit
   completion request; saying "merged" alone completes nothing.
4. A no, or a yes to only one of them, does only what was agreed.

Leave worktrees and `cc/*` branches alone. Removing them after a merge is a
separate request.
