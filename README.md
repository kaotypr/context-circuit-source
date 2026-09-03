# Context Circuit source

This repository is `context-circuit-source`, the maintainer source that builds
Context Circuit — a universal project workspace for AI-assisted work that
connects and coordinates one or more Git repositories. The source is not bound to
a version; it holds the current design and may lead the version last published as
the released artifact, `context-circuit-template`: a clean, uninitialized
universal project workspace. This checkout is the product source and its
self-hosted maintainer workspace.

## Product experience

In a released workspace you talk to the agent in ordinary language:

> What is this workspace?
> I want to add billing — retry a failed charge before failing the order.
> Approve this intent.
> Execute it, then ship it.

You approve the **intent** — what "correct" means and what scope is in bounds, after
an independent adversary has challenged the criteria. The coordinator derives a plan
within that scope (no separate plan approval); the workspace prepares isolated
repository worktrees; one worker
implements the whole plan and commits each repository; an independent read-only
verifier checks the latest commits; the worker repairs failures with new
commits; and completion follows the tier: Standard after candidate acceptance
plus delivery, Critical after explicit human completion. Explore is planless and
human-supervised. Intent approval and delivery are the two explicit human gates;
archive and restore remain separate organization actions.

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
- `sources/system-design/context-circuit/` — the authoritative maintainer design
  set (core plus scoped increments in versioned subfolders); maintainer material
  only, never shipped.
- `test/` — semantic acceptance suites (the deterministic engine-level laboratory
  is run from here via `agent-harness/test-template-runtime.sh`).
- `agent-harness/` — the built-template behavior laboratory: the deterministic
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

Read the design set under `sources/system-design/context-circuit/` only as
maintainer design material. The released artifact uses the shipped wrapper and
template, not the source repository's maintainer state.
