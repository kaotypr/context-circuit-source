# Runtime form

Scope landing. The substance and shape of the Context Circuit runtime — today
`wrapper/runtime/engine.sh`.

- [design.md](./design.md) — **normative overview**. The capability, the four
  goals driving the question, the principles, the fixed decisions, the open forks
  (form route, language, distribution, migration, sequencing), and the option
  comparison. Read this end to end to understand the scope without opening the
  detail file.
- [test-strategy.md](./test-strategy.md) — the testing concern at implementation
  depth: moving verification off white-box sourcing and source-text grep onto a
  black-box runtime action contract, and why that migration is contract-first and
  must land before the runtime's form can change.

This scope broadens the concern the v0.7.0 `runtime-opacity` scope opened;
runtime-opacity stays in place and active there, and opacity is treated here as a
property of the chosen form (see `design.md` → Principles and Fixed decisions).
