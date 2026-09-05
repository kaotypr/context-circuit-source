# Intention — i005

_Status: draft, waiting for your approval._

## Intention

What you want: **a workspace started from the Context Circuit template should not keep a leftover `.code-review-graph` folder.**

That folder is not part of the product. It should not ship with the template, and a workspace that already has it should delete it.

## Expectations

- A new workspace from the template does not include `.code-review-graph`.
- A workspace that already has `.code-review-graph` no longer keeps it.
- Intent approval, tracing, planning, execution, and delivery are unchanged.

## The plans

1. **Stop keeping `.code-review-graph`.**
   _After this:_ the template and a real workspace no longer carry that folder.
2. **Prove a fresh workspace is clean.**
   _After this:_ checks show a new workspace does not include it.

## How carefully this is checked

**`Standard`**

This changes what a workspace ships with, so an independent verifier should confirm the leftover is gone and that the rest of the workflow is untouched.

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
