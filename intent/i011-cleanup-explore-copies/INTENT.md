# Intention — i011

_Status: draft, waiting for your approval._

## Intention

What you want: **when you ask for runtime cleanup, Explore worktrees are included — not only plan-execution worktrees.**

If Explore worktrees live under `.runtime/explore/`, cleanup must know that path too. An explicit runtime-cleanup request removes those Explore worktrees as well.

## Expectations

- An explicit runtime-cleanup request also removes Explore worktrees.
- Plan-execution worktrees continue to be cleaned up as they are today.
- Cleanup stays something you ask for; it does not run by itself when an Explore session closes.
- A still-live Explore session is not silently deleted.

## The plans

1. **Teach cleanup about Explore worktrees.**
   _After this:_ an explicit runtime-cleanup request also removes worktrees under `.runtime/explore/`.
2. **Prove they are included.**
   _After this:_ checks fail if a cleanup request leaves Explore worktrees behind, and fail if a live Explore session is deleted without you asking.

## How carefully this is checked

**`Standard`**

This is cleanup of real worktrees, so an independent verifier should confirm Explore worktrees are included and that a live Explore session is not dropped by accident.

Explanations:
- **Explore:** you check it yourself as you work alongside the agent — no separate
  verification. Meant for trying things out.
- **Standard:** a second, independent agent verifies the finished work against your
  definition of done before it's called done.
- **Critical:** the strictest — independent verification plus a repair-and-recheck
  loop. For anything risky or impossible to undo (money, security, production, data
  you can't get back).

## Open questions

No known unresolved human decisions at draft time.
