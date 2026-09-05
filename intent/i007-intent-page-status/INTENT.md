# Intention — i007

_Status: draft, waiting for your approval._

## Intention

What you want: **`INTENT.md` must not show a status that is wrong after the tracer and feasibility check have finished.**

After you approve an intent, the tracer runs. When it finds no follow-up questions and the feasibility check marks the intent feasible, the status line in `INTENT.md` still does not move. Either that status line is kept in sync with what actually happened, or it is removed from `INTENT.md` so it cannot go stale.

## Expectations

- After a feasible tracer result with no follow-up questions, `INTENT.md` does not keep a misleading status.
- You pick one rule: keep the `INTENT.md` status line and make it match, or remove the status line from `INTENT.md`.
- `contract.yaml` remains the machine authority for approval; the `INTENT.md` status line does not replace it.
- Plans still follow approval plus a feasible tracer; this does not add a second approval.

## The plans

1. **Apply the status rule you pick.**
   _After this:_ `INTENT.md` either shows a status that matches a completed feasible tracer, or it no longer has a status line.
2. **Prove `INTENT.md` cannot lie.**
   _After this:_ checks fail if `INTENT.md` still shows a stale status after a feasible tracer with no follow-up questions.

## How carefully this is checked

**`Standard`**

This is the human-facing intent file, so an independent verifier should confirm `INTENT.md` cannot go stale.

Explanations:
- **Explore:** you check it yourself as you work alongside the agent — no separate
  verification. Meant for trying things out.
- **Standard:** a second, independent agent verifies the finished work against your
  definition of done before it's called done.
- **Critical:** the strictest — independent verification plus a repair-and-recheck
  loop. For anything risky or impossible to undo (money, security, production, data
  you can't get back).

## Open questions

**Should `INTENT.md` keep a status line that updates after the tracer finishes with a feasible finding and no follow-up questions, or should the status line be removed from `INTENT.md`?**
