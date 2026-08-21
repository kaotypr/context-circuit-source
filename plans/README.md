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
  reference: sources/example-prd.md
repositories:
  - example-repository
product_knowledge:
  references:
    - context/WORKSPACE.md
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

## Archive eligibility

Archival is a separate, durable eligibility decision; it is not a fourth plan
status and never rewrites `plan.yaml` or task projection. An unarchived plan
has no `archive.yaml` sidecar. Its absence preserves the historical active
behavior and needs no migration.

An archived bundle stays at its original path and adds `archive.yaml`:

```yaml
schema_version: 1
plan: plan-archive-capability
events:
  - action: archived
    at: 2026-08-20T18:00:00Z
    actor: sess-example
    reason: superseded by a replacement plan
    observed_status: approved
    replacements:
      - plans/context-circuit-plans/0025-progressive-context-routing
  - action: restored
    at: 2026-08-21T09:00:00Z
    actor: sess-example
    reason: the replacement was withdrawn
    observed_status: approved
    replacements: []
```

The final event is the current eligibility: `archived` excludes the bundle
from ordinary discovery; `restored` makes it discoverable again. Events are
append-only. Archive and restore retain the plan identity, decision actor,
time, reason, status observed at the decision, and replacement references;
they do not remove prior archive or restore evidence. A malformed sidecar, an
unknown action, a mismatched plan identity, or a record with no events is a
blocker rather than a reason to infer eligibility.

Archived bundles remain explicitly readable at their original paths for
history and dependency analysis. They are excluded from ordinary session
entry, next-action recommendations, approval, `cc-run-plan`, `cc-run-stack`,
and finishing. `cc-archive-plan` is the only capability that changes this
eligibility, and each archive or restore requires its own explicit current
session `archive` confirmation.

Dependency resolution remains visible. An archived `done` dependency is
inspectable and still satisfies a dependent as done. An archived `draft` or
`approved` dependency cannot satisfy a dependent or disappear from its
dependency check; it blocks until restored and normally resolved, or until a
separately human-approved dependency revision replaces it.

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

Execution uses `cc-run-plan` as the sole standard single-plan execution entry.
The root directs a writer child and an independent verifier child; every writer
gets an exclusive worktree and every verifier remains independent and
read-only. Sequential tasks share one writer child rather than skipping
children. Standalone `cc-run-plan` uses the repository default or active
branch. Connected already-approved plans enter `cc-run-stack`, which freezes a
runtime `graph.yaml`, resumes from `progress.yaml`, and bases dependent
worktrees on parent frozen SHAs, including in-run joins. A stack run is
runtime state under `.runtime/stacks/<stack-id>/`, not a plan type. There is
no `plans/<repository-key>-stacks/` layout, no durable `stack.yaml`, and no
scheduler. Canonical plan status stays `draft` / `approved` / `done`;
implemented is runtime evidence, not a plan status.
