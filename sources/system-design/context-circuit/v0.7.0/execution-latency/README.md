# execution-latency

The v0.7.0 scope for making Context Circuit **actions finish sooner without
changing what they mean**. Current execution is already acceptable; this is a
"faster where we safely can" pass. It measures first, then applies the two
inference-layer levers whose safety the contract already guarantees — overlapping
independent plans and tiering the model — while the deterministic runtime is left
thin and untouched.

## Reading order

1. [design.md](./design.md) — the normative overview: the one capability, the
   two-layer latency model, principles, the central decision (latency is an
   inference-layer problem and the safe levers already exist), what changes
   relative to v0.6, and the implementation order. Stop here for a review-level
   understanding.
2. [measurement.md](./measurement.md) — the two-observer model: deterministic
   phase durations the engine stamps vs inference wall-clock only the coordinator
   can see, the additive record fields, and the acceptance assertions that make
   "faster" provable. Ships first because it aims the other two.
3. [concurrent-run-stack.md](./concurrent-run-stack.md) — overlapping
   provably-independent plans: why the safety is already built
   (INV-CONCURRENCY-01/02), the coordinator fan-out over the ready bucket, the
   lease as race arbiter, fan-out width as policy, and the coordination/reporting
   cost that is the real price.
4. [model-tiering.md](./model-tiering.md) — concrete per-role `(model, effort)` in
   host-local config with adapter defaults: model/effort as a coordinator/host
   concern (INV-RUNTIME-01, INV-HOST-01), per-role `escalate_on_repair` composing
   with the failure counter (INV-REPAIR-01) without changing accounting, the
   optional per-plan complexity hint, and how verifier independence is preserved.

## Authority

This scope drafts source only. It proposes **no new invariant**: its concurrency
safety is existing INV-CONCURRENCY-01/02, and its tiering safety is existing
INV-HOST-01 / INV-RUNTIME-01. Its only contract surface is **additive schema
fields** (execution timing fields; an optional `complexity` plan field) and
**coordinator policy** (fan-out width, tier ladder), which settle in
`wrapper/contracts/` on acceptance. It grants no route, role, or authority by
itself (INV-SKILL-01).
