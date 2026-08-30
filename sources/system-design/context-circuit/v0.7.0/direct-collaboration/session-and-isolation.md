# session-and-isolation

Continues [design.md](./design.md). This file specifies **where a `cc-pair` session
writes and how it is bracketed**: the own-branch-and-worktree rule, base selection,
the light resumable pointer, `host-blocked`, convergence, and how a pairing branch
reaches delivery. It is deliberately light — a pairing session is not an execution.

## Own branch and worktree, from a base commit

A session always works in a **fresh `cc-pair/<session>` branch and its own
worktree**, created from one chosen **base commit**:

```
git worktree add <session-worktree> -b cc-pair/<session> <base-commit>
```

The base is *just a commit*, so worktree creation is one uniform operation
regardless of where the base came from. This is why a fresh worktree is **easier
and safer** than reusing an existing branch/worktree: no need to check whether some
other worktree still exists, is clean, or is mid-repair, and the pairing commits
stay cleanly separate from any plan's commits.

Two hard rules:

- **Never the user's active branch.** The user's checkout is untouched.
- **Never a plan's execution branch in place.** If the base is a completed plan or
  stack, branch a *new* `cc-pair` branch **from its tip** — do not commit onto
  `cc/<plan>/<repo>`. Editing a plan's branch in place would collide with a possible
  repair and entangle human-supervised commits with the plan's verified commits.

## Base selection — a commit, chosen mostly by context

| Entry | Base commit |
| --- | --- |
| Fresh (repo bound, no plan) | the **anchor tip** (the branch the user works from) |
| Offered after a plan/stack execution completes | that execution's **tip** (a completed dependency-chained stack's leaf tip already contains the integrated result — see below) |
| Explicit | any base commit the user names |

There is no separate "integration" base. For a dependency-chained stack, run-stack
already integrates along the dependency edges, so the **stack's leaf tip already is
the integrated result** — pairing from that tip needs no extra merge. Force-combining
*independent* plans that run-stack kept as separate deliverables is an
integration/delivery concern, not a pairing one, and is out of scope here.

## The loop's working state

Within the session, the worker's edits are ordinary working-tree changes in the
session worktree. The worker commits when the work warrants it (the user's call),
not on a fixed cadence; every commit follows Conventional Commits and carries no AI
attribution (INV-COMMIT-01). Nothing is marked "done" or "verified" — there is
nothing in this mode that can (see [design.md](./design.md)).

## Light resumable pointer — not an execution record

The only state kept is a small pointer so the user can leave and resume:

```
.runtime/pairing/<session>/pointer.yaml
  repo:     <repository id>
  worktree: <session worktree path>
  branch:   cc-pair/<session>
  base:     <base commit>
```

That is all. There is **no** attempt log, verifier result, evidence set, failure
counter, or handoff — this is not an execution. Resuming a session is reading the
pointer and reattaching to the worktree.

## `host-blocked`

`cc-pair` needs to create the worker child and the session worktree. If the host
cannot create the worker child, or cannot create a worktree, the outcome is
**`host-blocked`, read-only** — the coordinator does not do the work itself (it
never writes), and it does not pretend a change was made. This mirrors the shared
`host-blocked` contract.

## Convergence and after

- **Convergence** is the user's explicit "that's it." The session's branch now
  holds the human-supervised work; the session pointer can be closed.
- **The branch is the user's to deliver** — opening a pull request from a
  `cc-pair/<session>` branch is a normal, separate, explicit delivery action that
  targets the recorded `anchor_branch`, exactly like any other branch (INV-DELIVER-01).
  Because the work is human-supervised and not "verified," the coordinator says so
  when offering delivery; it never implies a verified result.
- **Preservation.** A pairing branch and its commits are preserved, never silently
  cleaned up; the user decides whether to deliver, keep, or discard it. The session
  worktree may be removed once the session ends, but the branch stays until the user
  acts on it.

## Not a lease

No lease is taken, and none is needed. Plan/stack executions never write the anchor
checkout (INV-EXEC-03) and each pairing session has its own branch and worktree, so
a pairing session cannot race an execution or another pairing session over the same
files. Isolation comes from the separate worktree, not from a lock.
