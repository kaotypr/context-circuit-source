---
name: cc-review-plan
description: Review a plan for evidence, scope, dependencies, risks, contradictions, and human decisions before execution.
---

# Review a plan

Use this skill when the user asks whether a plan is ready, what needs review,
or what decisions remain before execution.

## Route reads

Use the `review` manifest in
`docs/agent-workspace-workflow.md#route-read-manifests`. This skill remains
read-only: a missing contract, ownership ambiguity, or contradiction is
reported as a review blocker and never repaired by changing `plan.yaml`.

## Inspect

Read the selected plan, its tasks, referenced Product Knowledge, source
provenance, repository-local instructions, relevant implementation evidence,
and active runtime ownership. Check that:

- `plan.yaml` is the canonical lifecycle record and has a known status;
- scope, non-goals, repository boundaries, dependencies, task paths, and stop
  conditions are explicit;
- acceptance criteria map to test and verification evidence;
- leases, worktrees, sessions, and handoffs do not show an ownership conflict;
- contradictions, stale context, missing evidence, and delivery actions remain
  visible.

For a new bundle, expect `plan.yaml`, concise `overview.md`, and validated task
Markdown. Treat missing specialist companions as valid when no complexity
signal requires them. If companions exist, confirm that the overview records
their signals and that they reference, rather than copy, canonical acceptance,
verification, and lifecycle fields.

For historical bundles, discover the available companions and keep them
readable. Existing requirements, solution, risk, delivery, acceptance, and
verification documents are compatible inputs even when the bundle predates the
compact layout; do not propose a rewrite merely to remove duplication.

Check new task output with the approved host-provided deterministic validation
capability before calling it ready: front-matter delimiters, YAML, required
fields, enums, nested shapes, and the `context-circuit.task` schema. A failed
task check is a readiness defect. A passing check is not approval, execution,
verification, or completion.

## Report

Return the structured outcome from `docs/plan-review.md`. Distinguish routine
checks from human decisions. A review may recommend approval, revision, or a
blocker, but never changes plan status or claims completion. This skill remains
read-only and never writes `plan.yaml`.

If the outcome is `ready-for-approval`, name `cc-approve-plan` as the next
action. If the plan is already approved and dependency-ready, the
`approved-for-execution` outcome observes that fact and may recommend
`cc-run-plan`. Do not create a separate task-execution ceremony.
