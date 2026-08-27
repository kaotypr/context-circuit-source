# System-design authoring — design

The v0.6 scope that ships a **skill for authoring and structuring a system
design** as a piece of source material. It is **not** a lifecycle stage: a system
design is one kind of source, and this scope adds only the guidance for producing
a well-structured one. v0.6 is a **delta on v0.5** — see
[../README.md](../README.md) for the version index.

## Reading order

1. [design.md](./design.md) — the overview: what this is (a source, not a stage),
   why it needs a skill, principles, and fixed decisions.
2. [authoring-rubric.md](./authoring-rubric.md) — the crux: where a system design
   lives, the three-tier layout, **how much detail per file**, and **how to
   separate scopes by concern**.
3. [product-knowledge-relationship.md](./product-knowledge-relationship.md) — how a
   system design (a source) feeds Product Knowledge and plans through the
   **existing** v0.5 flow, with no new mechanism.
4. [contracts.md](./contracts.md) — the surface: one shipped skill and its
   allowlist wiring. No engine, no runtime records, no schemas, no invariants.
5. [examples.md](./examples.md) — a worked example and the acceptance criteria.

## Authority

This scope adds no runtime owner and duplicates no rule. Its one owner is the
shipped `cc-system-design` skill; everything else is the existing v0.5 behavior.
