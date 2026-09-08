---
kind: domain
status: accepted
title: Run-stack
slug: run-stack
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: cb84870
    basis: current-wrapper
generated_at: 2026-08-27T00:00:00Z
review_date: 2026-11-27
freshness: accepted-from-current-wrapper
assumptions:
  - Concurrency changes only order and overlap of already-authorized plans; it adds no authority.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-27
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/getting-started.md
---

# Run-stack

## Summary

Executing a *set* of intent-authorized plans in one request — a "plan stack" — so any
conflict between them becomes a scheduling decision made before a worker runs,
not a merge collision discovered afterward. Route "execute plans `<X>`…`<Z>`" and
"run the ready stack" here. Owned by the `cc-run-stack` skill; the run-stack
action is owned by `.context-circuit/wrapper/adapters/WORKFLOW.md`. Adds no new authority: every
plan is still individually authorized by its intent, verified, completed, and delivered.

## Scope

Inside: inter-plan dependencies, path leases, execution-base selection, the
readiness/scheduling loop, failure containment, and the delivery drift guard.

Outside: single-plan execution ([plan-execution](../plan-execution/README.md)),
verification ([verification](../verification/README.md)), completion
([completion](../completion/README.md)), and delivery
([delivery](../delivery/README.md)) — each unchanged and reused per plan.

## Behavior

- **Inter-plan dependencies (INV-PLAN-05).** `plan.yaml` may declare optional
  `plan_dependencies` (`{id, reason}`), distinct from a task's intra-plan
  `depends_on`. A v1.0 plan is `schema_version: 3` (the only supported version);
  entries reference an existing, non-self plan id and the graph is acyclic.
- **Path leases (INV-CONCURRENCY-01).** The one-worker lock generalizes to
  `(repository, path-region)` scope, recorded under `.runtime/locks/paths/`.
  Regions overlap when equal, when one is a path-prefix ancestor of the other, or
  when either is the repository-wide `.`. A non-descendant holder blocks an
  overlapping acquire (`LEASE_CONFLICT`); a declared descendant is exempt because
  it builds on the holder. A lease is held from execution start until delivery;
  release preserves the record. Composes with INV-OWN-01.
- **Execution bases (INV-CONCURRENCY-02).** A dependent plan's base in each
  repository is the base tip (no same-repo predecessor), the single predecessor
  branch (stack), or a runtime-authored integration merge of the predecessor
  branches (two or more), recorded as `base_commit` with `based_on`. The kept ref
  lives at `refs/cc-base/<plan>/<repo>` (never nested under a branch ref). The
  integration merge is authored by the runtime before the worker starts and is
  never a worker attempt; a base that cannot be built cleanly is *blocked*, not a
  worker failure; a predecessor repair invalidates and rebuilds a stale base.
- **Readiness and scheduling.** A plan is runnable when every dependency is
  `verified` (an AND-join) and every needed region is free. `cc_run_stack_ready`
  partitions a set into verified / ready / waiting / failed / blocked / refused; a
  plan whose latest execution is terminal or in progress is not runnable. The
  runtime detects readiness deterministically; the coordinator decides how many
  ready plans to launch — there is no scheduler heuristic in the runtime
  (INV-RUNTIME-01).
- **Failure containment.** A failed or blocked plan holds only its descendants;
  unrelated verified and ready plans are unaffected.
- **Delivery drift guard (INV-DELIVER-01).** See [delivery](../delivery/README.md).

## Workflows

- Execute a batch of intent-authorized plans in one go: `.context-circuit/docs/getting-started.md`

## Interfaces

- Human request: "Execute plans `0001` through `0010`" / "run the ready stack"
- Records: `plan.yaml` `plan_dependencies`; `.runtime/locks/paths/<repo>/<holder>.yaml`;
  `execution.yaml` `based_on`; base ref `refs/cc-base/<plan>/<repo>`

## Data

Path-lease records (`plan`, `repository`, `regions`, `acquired_at`, `released_at`)
under `.runtime/locks/paths/`; per-repository `based_on` in `execution.yaml`.

## Constraints and edge cases

Disjoint leases run concurrently and merge cleanly by construction; overlapping
leases serialize even absent a dependency. Leases prevent textual conflicts only,
never semantic coupling across different files (model that as one plan or a
dependency). A single-repository fan-in plan holds everything below it if a
predecessor fails — keep such a plan late and thin.

## Implementation references

- `.agents/skills/cc-run-stack/SKILL.md`, `.context-circuit/wrapper/adapters/WORKFLOW.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_plan_dependencies`, `cc_plan_dep_closure`,
  `cc_plan_is_descendant`, `cc_region_overlap`, `cc_lease_check`,
  `cc_lease_acquire`, `cc_lease_release`, `cc_plan_same_repo_preds`,
  `cc_base_prepare`, `cc_plan_ready`, `cc_run_stack_ready`
- `.context-circuit/wrapper/contracts/schemas/lease.yaml`, `plan.yaml` (`plan_dependencies`,
  `schema_version [1, 2]`), `execution.yaml` (`based_on`)
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-PLAN-05, INV-CONCURRENCY-01,
  INV-CONCURRENCY-02

## Verification

`sh test/acceptance.sh` (path leases + run-stack suites); live scenario
`agent-harness/scenarios/10-run-approved-stack` (grade.sh PASS + human-simulator
pass; ten plans built and verified in dependency order with integration bases for
the fan-ins).

## Acceptance notes

Accepted 2026-08-27.
