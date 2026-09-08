# Intention — i012

_Status: approved._

## Intention

What you want: **`role-tiering.local.yaml` must apply in Explore the same way it applies in `cc-execute`.**

You already set per-role model and effort in `role-tiering.local.yaml`. In Explore, the worker did not follow that file. The `cc-pair` skill should include the same "Model & effort per role" step that `cc-execute` already has, so the Explore worker runs at the configured model and effort.

## Expectations

- `role-tiering.local.yaml` is honored during Explore, including the worker.
- The `cc-pair` skill includes the same Model & effort per role step that `cc-execute` uses.
- This changes cost and speed only — not what Explore means, who judges it, or whether a verifier is required.
- `cc-execute` role-tiering behavior is not weakened.

## The plans

1. **Honor `role-tiering.local.yaml` in Explore.**
   _After this:_ the Explore worker runs at the model and effort set for that role.
2. **Add the Model & effort step to `cc-pair`.**
   _After this:_ `cc-pair` includes the same Model & effort per role step as `cc-execute`.
3. **Prove the worker follows it.**
   _After this:_ checks fail if an Explore worker ignores `role-tiering.local.yaml`.

## How carefully this is checked

**`Standard`**

This is a real-life miss of a setting you already rely on, so an independent verifier should confirm an Explore worker follows `role-tiering.local.yaml`.

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
