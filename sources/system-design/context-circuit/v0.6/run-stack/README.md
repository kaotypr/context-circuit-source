# Run-stack — design

The v0.6 scope for executing a set of approved plans in one run (a *plan stack*).
v0.6 is a **delta on v0.5** — see [../README.md](../README.md) for the version
index and [../../v0.5/](../../v0.5/) for the base design.

## Reading order

1. [design.md](./design.md) — the overview: capability, principles, the shape of
   the solution, compatibility, implementation order, and fixed decisions.
2. [plan-dependencies.md](./plan-dependencies.md) — inter-plan dependencies, and
   how same-repo vs cross-repo dependencies differ.
3. [path-leases.md](./path-leases.md) — `(repo, path)` leases (INV-CONCURRENCY-01)
   and what they do and do not guarantee.
4. [execution-bases.md](./execution-bases.md) — anchor / stacked / integration
   base, where and when the merge happens (INV-CONCURRENCY-02), repair invalidation.
5. [scheduling.md](./scheduling.md) — the run-stack loop, readiness, emergent
   waves, and failure/blast-radius containment.
6. [delivery.md](./delivery.md) — dependency-order delivery and the drift guard.
7. [schema-and-versioning.md](./schema-and-versioning.md) — the contract deltas
   and why the plan schema bumps to 2 but execution does not.
8. [examples.md](./examples.md) — two worked walkthroughs and the acceptance criteria.

## Authority

This scope adds no owner and duplicates no rule; it specifies intent and points
to the canonical owners under `wrapper/`. The contract deltas are summarized in
[schema-and-versioning.md](./schema-and-versioning.md).
