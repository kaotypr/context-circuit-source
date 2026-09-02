# 0013 — Run-stack substrate: plan dependencies, path leases, execution bases

- **Plan ID:** `0013-run-stack-substrate`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0004-planning`, `0005-approval-and-execution`
- **Owns (invariants):** INV-PLAN-05, INV-CONCURRENCY-01, INV-CONCURRENCY-02

## Original request

Retroactive plan for the run-stack substrate, built as if from an empty repo.
Source design: `sources/system-design/context-circuit/v0.6/run-stack/`
(`schema-and-versioning`, `plan-dependencies`, `path-leases`, `execution-bases`).

## Objective and desired behavior

Let a set of already-approved plans run together so any conflict is a scheduling
decision made **before** a worker runs, not a merge collision found after.

- `plan.yaml` gains optional `plan_dependencies` (`{id, reason}`); using the
  field requires `schema_version 2`; the graph is acyclic with existing,
  non-self references. This plan owns the coordinated bump — plan schema `[1,2]`
  and `runtime_version 0.6.0` — that repository-grounding (0015) also rides.
- Path leases generalize the one-writer lock to `(repository, path-region)`
  scope under `.runtime/locks/paths/`; a non-descendant plan cannot acquire an
  overlapping region; a lease is held from execution start until delivery.
- The runtime selects a base per repository — base tip, a single predecessor
  branch, or a runtime-authored integration merge — recorded as `base_commit`,
  with base refs at `refs/cc-base/<plan>/<repo>`; a stale base is rebuilt.

## Constraints and non-goals

- Non-goal: the scheduling loop, run action, and delivery drift guard (0014).

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `design-deltas`
  (`context/DESIGN-DELTAS.md`).

## Tasks

1. **RSS-001** — plan-dependency contract + coordinated version bump.
2. **RSS-002** — path leases.
3. **RSS-003** — execution base selection.

## Acceptance & verification

- schema_version 2 enforced; cyclic/unknown deps rejected; overlapping region
  blocked for non-descendants; integration merge is runtime-authored; unbuildable
  base is blocked.
- `sh test/contracts/test-contracts.sh`, `sh test/concurrency/test-leases.sh`,
  `sh test/run-stack/test-run-stack.sh`.

## Assumptions, open questions, risks

- Risk: a v0.5 engine mis-scheduling dependency-bearing plans — prevented by the
  asymmetric schema bump (such a plan is `schema_version 2`, refused by v0.5).

## Expected commits and delivery notes

Coordinated with 0015 in one release; the manifest bump is authored here.
