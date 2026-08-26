# Context Circuit source

This repository is `context-circuit-source`, the maintainer source that builds
Context Circuit v0.5 — a universal project workspace for AI-assisted work that
connects and coordinates one or more Git repositories. The released artifact is
`context-circuit-template`: a clean, uninitialized universal project workspace.
This checkout is the product source and its self-hosted maintainer workspace.

## Product experience

In a released workspace you talk to the agent in ordinary language:

> What is this workspace?
> Create a plan for adding billing.
> Approve plan 0001-billing-v2 and execute it.
> Mark 0001-billing-v2 complete.

The coordinator drafts or reviews a readable plan; you approve it in
conversation; the workspace prepares isolated repository worktrees; one worker
implements the whole plan and commits each repository; an independent read-only
verifier checks the latest commits; the worker repairs failures with new
commits; and you decide when the plan is complete. Approval, execution,
completion, archive, restore, and delivery are separate explicit human actions.

## Source layout

- `wrapper/` — shipped runtime, contracts, schemas, adapters, and migration
  boundary. `wrapper/runtime/engine.sh` is the small host-neutral deterministic
  runtime; `wrapper/contracts/invariants.yaml` is the one-rule-one-owner map.
- `template/` — the blank mutable seed for a new workspace.
- `.agents/skills/` and `agents/` — thin host skills and worker/verifier/
  coordinator role deltas.
- `docs/` — shipped guides and plan/task templates.
- `context/`, `plans/` — source-only maintainer Product Knowledge and plans; not
  released.
- `sources/system-design/context-circuit/v0.5/` — the authoritative v0.5 design;
  maintainer material only, never shipped.
- `test/` — semantic acceptance suites (the deterministic engine-level laboratory
  is run from here via `template-harness/test-template-runtime.sh`).
- `template-harness/` — the built-template behavior laboratory: the deterministic
  engine-level suite plus the human-simulated harness (`human/`, `scenarios/`).
  Source-only; never shipped in `context-circuit-template`.
- `scripts/` — maintainer-only release assembly.
- `repositories/`, `repositories.local.yaml`, `.runtime/` — host-local, ignored,
  never released.

## Tests

Run the complete semantic acceptance suite:

```sh
sh test/acceptance.sh
```

Read `sources/system-design/context-circuit/v0.5/core/design.md` only as
maintainer design material. The released artifact uses the shipped wrapper and
template, not the source repository's maintainer state.
