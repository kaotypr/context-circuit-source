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

You approve the **intent** — what "correct" means — from the plain ask. On approval a
planner reads the real code and writes the plan; the coordinator runs a feasibility check
and publishes that plan (no separate plan approval and no automated scope
gate; scope-safety is settled at delivery). The workspace prepares isolated
repository worktrees; one worker
implements the whole plan in that plan's single repository and commits it; an independent read-only
verifier checks the latest commits; the worker repairs failures with new
commits; and completion follows the tier: Standard after candidate acceptance
plus delivery, Critical after explicit human completion. Explore is planless and
human-supervised. Intent approval and delivery are the two explicit human gates;
archive and restore remain separate organization actions.

## Source layout

- `.context-circuit/wrapper/` — shipped runtime, contracts, schemas, adapters, and migration
  boundary. `.context-circuit/wrapper/runtime/engine.sh` is the small host-neutral deterministic
  runtime; `.context-circuit/wrapper/contracts/invariants.yaml` is the one-rule-one-owner map.
- `template/` — the blank mutable seed for a new workspace.
- `product/` — packaging-only shipped skills and native host integrations; its
  nested discovery trees become workspace-root paths in the output.
- `.agents/skills/` — source-only direct-development skills, never released.
- `.context-circuit/agents/` — shipped planner/worker/verifier/coordinator role deltas.
- `.context-circuit/docs/` — shipped guides and plan/task templates.
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

## Source development

Read and edit files directly on the active branch under root `AGENTS.md` and
`WORKFLOW.md`. Use `cc-source-develop` for maintainer work; product lifecycle
skills are packaged under `product/` and do not govern this checkout. Generated
workspaces retain their own adapters and lifecycle unchanged.

## Tests

Run the complete semantic acceptance suite:

```sh
sh test/acceptance.sh
```

Read the design set under `sources/system-design/context-circuit/` only as
maintainer design material. The released artifact uses the shipped wrapper and
template, not the source repository's maintainer state.
