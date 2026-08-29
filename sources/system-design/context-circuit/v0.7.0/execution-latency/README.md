# execution-latency

The v0.7.0 scope for making Context Circuit **actions finish sooner without
changing what they mean**. Current execution is already acceptable; this is a
"faster where we safely can" pass. It applies the two inference-layer levers whose
safety the contract already guarantees — overlapping independent plans and tiering
the model per role — while the deterministic runtime is left thin and untouched.
(This scope reflects what shipped: an early "measure the layers" cut was built and
then trimmed to the honest minimum — see [measurement.md](./measurement.md).)

## Reading order

1. [design.md](./design.md) — the normative overview: the one capability, the
   two-layer latency model, principles, the central decision (latency is an
   inference-layer problem and the safe levers already exist), what changes
   relative to v0.6, and the implementation order. Stop here for a review-level
   understanding.
2. [concurrent-run-stack.md](./concurrent-run-stack.md) — overlapping
   provably-independent plans: why the safety is already built
   (INV-CONCURRENCY-01/02), the coordinator fan-out over the ready bucket, the
   lease as race arbiter, fan-out width as policy, and the coordination/reporting
   cost that is the real price.
3. [model-tiering.md](./model-tiering.md) — per-role `(model, effort)` in a
   host-local, host-grouped config with adapter defaults: model/effort as a
   coordinator/host concern (INV-RUNTIME-01, INV-HOST-01), that **recording a tier
   is not running it** (the adapter applies it on the child spawn; per-child effort
   may be unavailable), per-role `escalate_on_repair` composing with the failure
   counter (INV-REPAIR-01) without changing accounting, the optional per-plan
   complexity hint, and how verifier independence is preserved.
4. [measurement.md](./measurement.md) — what an execution actually records (the
   per-role `(model, effort)`), timing from the pre-existing `started_at`/
   `checked_at` span, and the honest account of the phase_ms / coordinator-
   wall-clock machinery that was built and then cut.

## Authority

This scope drafts source only. It proposes **no new invariant**: its concurrency
safety is existing INV-CONCURRENCY-01/02, and its tiering safety is existing
INV-HOST-01 / INV-RUNTIME-01. Its only contract surface is **additive schema
fields** (the per-attempt `(model, effort)` host evidence; an optional
`complexity` plan field) and **coordinator/host policy** (fan-out width, tier
ladder, the spawn-application mechanism), which settle in `wrapper/contracts/` and
`wrapper/adapters/` on acceptance. It grants no route, role, or authority by
itself (INV-SKILL-01).
