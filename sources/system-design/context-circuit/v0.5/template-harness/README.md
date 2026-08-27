# Template harness — design

The maintainer **test-harness** scope of Context Circuit v0.5: how the built
`context-circuit-template` is exercised as a real workspace by a human-simulated
conversation. Normative design: [design.md](./design.md).

This is the *design* of the harness; the harness *implementation* lives at the
repository root under `template-harness/`. For the version index and the other
v0.5 scope, see [../README.md](../README.md).

## Reading order

1. [design.md](./design.md) — what the harness tests and why it is separate from
   the deterministic runtime laboratory; the maintainer-only boundary.
2. [roles/](./roles/) — the actors and their interaction contract:
   human-simulator, coordinator-under-test, and grader.
3. [scenario-library.md](./scenario-library.md) — the natural-language scenario
   prompt libraries under `test/`, the case-file format, and the enumerated cases.
4. [harness-and-evaluation.md](./harness-and-evaluation.md) — run lifecycle,
   workspace isolation, the turn protocol, grading, result artifacts, and boundaries.
5. [host-matrix.md](./host-matrix.md) — running the same scenarios under each host
   adapter (Codex CLI, Claude Code, Cursor Agent).

## Authority

Maintainer test design. It does not change the product contract in `../core/` or
`wrapper/contracts/`. Where it observes product behavior, the core design and
`wrapper/contracts/invariants.yaml` remain the authority for what "correct" means.
