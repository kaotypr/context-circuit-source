# Proposal 0016 — add domain: run-stack

- id: 0016-domain-run-stack
- target_context_unit: context/domains/run-stack/README.md
- operation: add
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [74eb510]
- evidence_refs:
  - sources/system-design/context-circuit/v0.6/run-stack/
  - wrapper/contracts/invariants.yaml (INV-CONCURRENCY-01, INV-CONCURRENCY-02, INV-PLAN-05, INV-DELIVER-01 drift-guard clause)
  - wrapper/contracts/schemas/lease.yaml
  - wrapper/contracts/schemas/plan.yaml (plan_dependencies; schema_version [1, 2])
  - wrapper/contracts/schemas/execution.yaml (based_on)
  - wrapper/runtime/engine.sh (cc_plan_dependencies, cc_plan_dep_closure, cc_region_overlap, cc_lease_{check,acquire,release}, cc_base_prepare, cc_plan_ready, cc_run_stack_ready, cc_delivery_{drift,rebase})
  - .agents/skills/cc-run-stack/SKILL.md
  - wrapper/adapters/WORKFLOW.md ("Execute plans X…Z / run the ready stack")
  - test/concurrency/test-leases.sh, test/run-stack/test-run-stack.sh
  - template-harness/scenarios/10-run-approved-stack/case.yaml

## Statement

v0.6 ships **run-stack**: executing a *set* of approved plans in one request (a
"plan stack"), where any conflict between them becomes a scheduling decision made
*before* a worker runs, not a merge collision discovered afterward. It adds no new
authority — approval, verification, completion, and delivery gates are unchanged.
This is a bounded domain distinct from single-plan execution and should have its
own retrieval unit; [plan-execution](../plan-execution/README.md) is updated
separately (proposal 0018) to cross-reference it.

## Proposed unit (summary to author on acceptance)

- **Inter-plan dependencies.** `plan.yaml` gains optional `plan_dependencies`
  (`{id, reason}`); declaring them requires `schema_version: 2` (INV-PLAN-05;
  manifest accepts plan `[1, 2]`). Distinct from a task's intra-plan `depends_on`.
- **Path leases (INV-CONCURRENCY-01).** The one-writer lock generalizes to
  `(repository, path-region)` scope under `.runtime/locks/paths/`. Overlap =
  equal / path-prefix ancestor / repository-wide `.`; a non-descendant holder
  blocks (`LEASE_CONFLICT`); a declared descendant is exempt; held from execution
  start until delivery; release preserves the record.
- **Execution bases (INV-CONCURRENCY-02).** A dependent plan's base per repository
  is the anchor tip (no same-repo predecessor), the single predecessor branch
  (stack), or a runtime-authored integration merge (≥2 predecessors), recorded as
  `base_commit` with `based_on`; the kept ref lives at `refs/cc-base/<plan>/<repo>`
  (never nested under a branch ref). An unbuildable base is *blocked*, not a worker
  failure; a predecessor repair invalidates and rebuilds a stale base.
- **Readiness + scheduling.** `cc_plan_ready` is a dependency AND-join plus a lease
  gate; `cc_run_stack_ready` partitions a set into verified / ready / waiting /
  failed / blocked / refused (terminal or in-progress plans are not runnable). The
  runtime detects readiness deterministically; the coordinator (`cc-run-stack`
  skill) decides how many ready plans to launch — no scheduler heuristic in the
  runtime (INV-RUNTIME-01).
- **Failure containment.** A failed or blocked plan holds only its descendants;
  unrelated verified/ready plans are unaffected.
- **Delivery drift guard (INV-DELIVER-01, extended).** On delivery, if the recorded
  base has diverged from the current anchor tip, the plan is rebased onto the tip
  and re-verified before its pull request (`cc_delivery_drift` / `cc_delivery_rebase`,
  `DELIVERY_REBASE_CONFLICT`).

Per-plan discipline is unchanged: one worker, one independent verifier, the
three-failure limit, human completion, separate delivery.

## Verification

`sh test/acceptance.sh` (path leases + run-stack suites); live case
`10-run-approved-stack` (grade.sh PASS + human-simulator pass).
