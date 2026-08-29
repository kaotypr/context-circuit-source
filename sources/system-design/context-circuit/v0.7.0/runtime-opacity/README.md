# Runtime opacity — the engine is invoked, never read

The coordinator drives the runtime by **invoking** `wrapper/runtime/engine.sh`
as an opaque tool; it must never **read** the engine's implementation. Today that
discipline is real but soft — worded as *"you never need to read"* rather than
*"you must not read,"* spread across several skills, and enforced only by a
maintainer test gate. A live v0.6 run showed a real coordinator open the engine
source anyway, tripping the access-discipline gate. This scope hardens the
invoke-not-read boundary so a released template does not leak or violate it.

Start at [design.md](./design.md).

## Files

- [design.md](./design.md) — the normative overview: the boundary, the observed
  gap, the principles, the fixed decisions, and the design options (the concrete
  fix is deferred to a plan).

## Authority

No new invariant. This scope **strengthens the coordinator-side corollary of
INV-RUNTIME-01** (the runtime is a thin deterministic tool, not a document to
interpret) and unifies wording that already lives in `wrapper/adapters/AGENTS.md`
and the `cc-execute` / `cc-run-stack` / `cc-publish` skills. Canonical owners stay
under `wrapper/`.
