# Plans

Plans are human-reviewed definitions of intended work. They connect an intent
or accepted product artifact to the smallest useful Product Knowledge set,
repositories, bounded tasks, acceptance evidence, verification, risks, and
human gates.

## Plan artifact

Each plan lives under `<repository-key>-plans/<number>-<slug>/` and has a
canonical `plan.yaml` lifecycle record. The companion Markdown files explain
the objective, requirements, solution, delivery, acceptance criteria,
verification, risks, and task contracts.

The minimum metadata contract is:

```yaml
id: example-plan
number: 7
title: Deliver one bounded outcome
status: draft
source:
  kind: accepted-prd
  reference: contributions/prds/example.md
repositories:
  - example-repository
product_knowledge:
  references:
    - context/PROJECT.md
implementation_scope:
  - The files or behavior this plan may change
non_goals:
  - Explicitly excluded behavior
dependencies:
  - A plan or decision that must be ready first
acceptance_criteria:
  - Observable evidence required for review
test_scope:
  - Fixtures or checks that prove the criteria
verification_commands:
  - git diff --check
```

`plan.yaml` is the canonical lifecycle authority. Plan status is exactly
`draft`, `approved`, or `done`; included task status is a synchronized
projection of that value: `draft` → `ready` → `done`. Approval and completion
are human gates. Tests, Git changes, or an agent handoff never silently
approve or complete a plan.

## Task contract

Every task has a stable ID, bounded implementation paths, explicit
dependencies, acceptance criteria, verification commands, and stop conditions.
Task status is not a second approval gate or an execution lease. The
coordinator reconciles task projections in one idempotent operation when a plan
is approved, completed, or resumed.

## Planning boundaries

Plans may be drafted from an accepted Idea Brief, PRD, selected raw sources,
durable Product Knowledge, repository evidence, or a clear direct request.
Source paths and repository evidence are recorded as provenance. Unresolved
contradictions, assumptions, ownership conflicts, and missing acceptance
evidence remain visible for human review.

Execution uses `cc-run-plan` as the sole standard plan-execution entry. There
is no user-facing `cc-run-task` workflow. The root directs a writer child and
an independent verifier child; every writer gets an exclusive worktree and
every verifier remains independent and read-only. Sequential tasks share one
writer child rather than skipping children.
