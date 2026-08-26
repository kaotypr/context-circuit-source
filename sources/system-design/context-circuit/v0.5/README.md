# Context Circuit v0.5

Source design for Context Circuit v0.5, organized by scope. This README is the
version index and reading-order authority; each scope owns its own design.

## Scopes

- [core/](./core/) — the core product design: workspace model, readable plans,
  conversational workflow, multi-repository execution, worker/verifier behavior,
  runtime, and acceptance criteria. Start at [core/design.md](./core/design.md).
- [template-harness/](./template-harness/) — the maintainer test-harness design:
  harness and evaluation, host matrix, scenario library, and roles. Start at
  [template-harness/design.md](./template-harness/design.md). (This is the
  *design* of the harness; the harness *implementation* lives at the repository
  root under `template-harness/`.)

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. A **scope** is a bounded
area of the design — an enduring aspect (like `core` or `template-harness`) or a
bounded feature. The version folder holds one scope folder per area it covers;
this version README indexes them and carries the reading order, so scope grouping
never fights the linear read.

Every folder has a **`README.md`** as its landing and index — it auto-renders
when the folder is browsed. Each scope's normative design is **`design.md`**; the
scope's `README.md` states the scope's purpose and gives the reading order over
its files. The version `README.md` (this file) indexes the scopes. So both files
are predictable everywhere: `README.md` to land and navigate, `design.md` for the
normative design.

## Source boundary

These documents are maintainer design material. They are not copied into an
instantiated workspace as Product Knowledge and are not normal agent context. An
implementation session should receive only the contracts and detail its bounded
change needs.
