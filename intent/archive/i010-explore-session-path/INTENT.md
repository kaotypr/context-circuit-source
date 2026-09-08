# Intention — i010

_Status: approved._

## Intention

What you want: **when you ask to start Explore, you choose the name, and that Explore worktree lives under `.runtime/explore/<name>` — not mixed with plan-execution worktrees under `.runtime/worktrees/cc-pair/`.**

Today the Explore worktree sits with plan-execution worktrees, under a name the agent picks. Instead, the coordinator asks what this Explore session is about, and that name is used for a worktree in a dedicated Explore location.

## Expectations

- Starting Explore asks you for a short name or what it is about, and uses that name.
- The Explore worktree lives at `.runtime/explore/<human-name>`, separate from `.runtime/worktrees/`.
- The agent does not invent the folder name.
- Explore still uses one repository, you still judge the result live, and promote still works.

## The plans

1. **Ask for the name and keep Explore worktrees in their own place.**
   _After this:_ an Explore session is named by you and stored under `.runtime/explore/`.
2. **Prove the split.**
   _After this:_ checks fail if an Explore worktree is still created under `.runtime/worktrees/cc-pair/` or under an agent-invented name.

## How carefully this is checked

**`Standard`**

This changes where Explore worktrees live, so an independent verifier should confirm they are named by you and kept apart from plan-execution worktrees.

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
