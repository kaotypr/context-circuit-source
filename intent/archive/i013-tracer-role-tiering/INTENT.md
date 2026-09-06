# Intention — i013

_Status: approved._

## Intention

What you want: **when you ask the coordinator to set up `role-tiering.local.yaml` with the roles you name, including tracer, that tracer config is written into the file and used when the tracer is spawned.**

Today the coordinator creates `role-tiering.local.yaml` but refuses tracer: "Tracer tiering isn’t supported by this workspace contract, so tracers use each adapter’s default." The tracer `(model, effort)` you described is not set in the file.

Worker and verifier stay configurable as they are. Tracer becomes a supported role in the same file, same shape.

## Expectations

- Asking the coordinator to set role-tiering including tracer writes a `tracer` entry into `role-tiering.local.yaml`.
- A spawned tracer uses that `(model, effort)`, not only the adapter default.
- The coordinator no longer reports that tracer tiering is unsupported by the workspace contract.
- Role tiering still changes cost and speed only — not what the tracer does, not independence, not the failure limit.

## The plans

1. **Allow tracer in `role-tiering.local.yaml`.**
   _After this:_ the workspace contract accepts a `tracer` role next to `worker` and `verifier`.
2. **Honor it when spawning the tracer.**
   _After this:_ the coordinator writes the tracer config you asked for and spawns the tracer at that model and effort.
3. **Prove it is no longer refused.**
   _After this:_ checks fail if a requested tracer entry is dropped, or if the coordinator still says tracer tiering is unsupported.

## How carefully this is checked

**`Standard`**

This changes the role-tiering contract, so an independent verifier should confirm a requested tracer config lands in `role-tiering.local.yaml` and is used.

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
