# Reduce approval-to-plan latency

Plan ID: 0004-approval-plan-latency
Intent: i018-approval-plan-latency
Status: draft

## Original request and coverage

- Approved request: make the normal intent-approval path produce grounded,
  validated plans in about 90 seconds, with a typical two-minute warm ceiling and
  a three-to-four-minute cold target, instead of the current multi-minute path.
- This plan covers revision-bound reusable grounding (LAT-001), plan-ready trace
  evidence and same-repository adaptation (LAT-002), and atomic, measured stack
  materialization (LAT-003).
- It does not execute plans or add another human gate. Independent tracing,
  feasibility, assurance, and authorization remain in the same approval turn.

## Objective and desired behavior

A feasible approval produces either a complete validated plan stack or no stack.
Fresh structural repository maps shorten a trace but never replace first-hand
intent-relevant reading. The coordinator ratifies trace fragments without a
second repository survey. A deterministic runtime materializes the full stack
atomically and records phase timing. Trace anchors remain evidence, while workers
can record necessary same-repository discoveries for independent review.

## Constraints and non-goals

- The runtime stays deterministic and host-neutral; feasibility, cache validity,
  tier, and plan-boundary choices stay with the coordinator.
- Cache entries are revision-bound derived indexes, not source copies or lifecycle
  records; they contain no credentials, ignored files, or provider payloads.
- Cold fallback is a correctness path. Uncertain drift never permits stale reuse.
- Atomic publication covers plan files, dependencies, authorization, index rows,
  and feasible status. Existing one-plan operations remain supported.
- Advisory anchors never become an inferred worker-path allowlist. A different
  repository or a change to the approved decision returns to the coordinator.
- No execution, delivery, publication, commit, or second human approval occurs.

## Product Knowledge grounding

- tracing (`context/domains/tracing/README.md`) — independent post-approval trace,
  freshness reuse, and feasibility must remain before planning.
- assurance (`context/domains/assurance/README.md`) — Standard remains appropriate
  and retains one independent verifier.

The trace manifest identifies `engine.sh` as the deterministic owner of ID,
validation, index, and atomic-write primitives; its schema and role guidance as
the trace-to-plan interface; and the worker/verifier handoffs as the scope-safety
interface. One repository and a single approval/materialization lifecycle justify
one plan with ordered embedded tasks, rather than false independent stacked plans.

## Repositories and source evidence

- context-circuit-source — runtime, contracts, coordinator/role guidance, tests,
  and Product Knowledge. Evidence:
  `intent/i018-approval-plan-latency/trace/context-circuit-source.yaml`, recorded
  2026-09-07T16:25:25Z. Feasible; no out-of-scope reach; Standard remains suitable.

## Tasks

1. LAT-001 (context-circuit-source; trace contract/guidance, tracing context, and
   latency/feasibility tests; depends on: none) — implement immutable,
   revision-bound structural repository-map reuse and exact/delta/fallback evidence.
   Keep the tracer first-hand and independent. Acceptance LAT-AC-005 and
   LAT-AC-002; verification LAT-VT-001 and LAT-VT-002.

2. LAT-002 (context-circuit-source; trace, plan, handoff/verifier contracts and
   role/coordinator guidance/tests; depends on: LAT-001) — emit normalized,
   plan-ready trace fragments with typed advisory anchors and explicit criterion
   coverage, then reconcile worker/verifier expansion handling. Acceptance
   LAT-AC-004, LAT-AC-008, and LAT-AC-009; verification LAT-VT-003 through
   LAT-VT-005.

3. LAT-003 (context-circuit-source; runtime, contracts, coordinator/skills, and
   approval/latency/acceptance tests; depends on: LAT-002) — validate, reserve,
   stage, authorize, index, and publish the full plan stack as one idempotent
   operation, with monotonic timing and benchmark evidence. Acceptance LAT-AC-001,
   LAT-AC-003, LAT-AC-006, and LAT-AC-007; verification LAT-VT-006 through
   LAT-VT-011.

## Acceptance criteria

- LAT-AC-001 — one approval produces all feasible derived plans, fully created and
  validated, without another human action.
- LAT-AC-002 — warm approvals meet the roughly 90-second target and typical
  two-minute ceiling while preserving first-hand grounding.
- LAT-AC-003 — representative cold traces meet the three-to-four-minute target.
- LAT-AC-004 — fast planning retains inspectable grounding, risk, check, tier,
  feasibility, scope, and question evidence.
- LAT-AC-005 — relevant drift refreshes evidence or safely falls back to cold.
- LAT-AC-006 — multi-plan artifacts, IDs, dependencies, authorization, and indexes
  are consistent and never partially visible.
- LAT-AC-007 — phase-level timing exposes trace, feasibility, materialization,
  validation, and fallback costs.
- LAT-AC-008 — necessary same-repository expansion is recorded and independently
  checked without a plan rewrite.
- LAT-AC-009 — second-repository or decision-changing expansion stops for the
  coordinator; ordinary local discoveries do not create a human gate.

## Verification

- `sh test/latency/test-approval-plan.sh --case exact-delta-fallback`
- `sh test/intent/test-feasibility.sh`
- `sh test/intent/test-approval-trace-plan.sh`
- `sh test/latency/test-approval-plan.sh --case same-repository-expansion`
- `sh test/latency/test-approval-plan.sh --case expansion-stop`
- `sh test/latency/test-approval-plan.sh --case warm-standard`
- `sh test/latency/test-approval-plan.sh --case cold-multi-repository`
- `sh test/latency/test-approval-plan.sh --case multi-plan-atomicity`
- `sh test/latency/test-approval-plan.sh --case timing-evidence`
- `sh test/latency/test-approval-plan.sh`
- `sh test/acceptance.sh`

## Assumptions, open questions, risks

Plan-level resolutions during implementation:

- Choose a git-aware, gitignored cache location and revision command that retains
  the approved identity dimensions. Evolve the manifest schema compatibly or use a
  bounded legacy re-trace; never silently accept legacy evidence as plan-ready.
- Choose the smallest deterministic materializer interface. It must accept
  coordinator-ratified input rather than decide feasibility, tier, cache validity,
  or plan boundaries itself.
- Migrate path fields and handoff evidence so anchors are advisory but repository
  boundaries and protected constraints remain enforceable. Preserve old one-plan
  operations and fixtures unless explicitly migrated.

Key risks: stale or sensitive cache content; partially published stacks or duplicate
IDs; moving coordinator policy into the runtime; treating anchors as restrictions;
and timing-only tests that pass without quality evidence. The ordered tasks and
targeted regression cases address these before the complete acceptance suite.

## Expected commits and delivery notes

- One source-repository implementation commit is expected after execution and
  independent verification. No commit, delivery, merge, push, or publish is part of
  this plan derivation.

## Expected Product Knowledge impact

- Reassess `tracing` and `assurance` on completion. Update tracing for the durable
  cache/fragment/materialization interfaces if the implemented ownership or
  lifecycle warrants it; assurance should remain unchanged except for verified
  references to its retained Standard floor.
