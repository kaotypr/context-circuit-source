# Plans

Plans are human-reviewed definitions of intended work. They connect an intent
or accepted product artifact to the smallest useful Product Knowledge set,
repositories, bounded tasks, acceptance evidence, verification, risks, and
human gates.

## Plan artifact

Each plan lives under `<repository-key>-plans/<number>-<slug>/` and has a
canonical `plan.yaml` lifecycle record. The default generated bundle is
deliberately small:

```text
<plan>/
  plan.yaml
  overview.md
  tasks/
    <task-id>.md
```

`plan.yaml` is the sole authority for the plan's lifecycle, repository scope,
dependencies, acceptance criteria, and verification commands. Its status is
exactly `draft`, `approved`, or `done`; included task status remains the
synchronized projection `draft`, `ready`, or `done`. No Markdown companion,
task front matter, schema result, runtime record, or verifier handoff can
approve, execute, verify, or complete a plan.

`overview.md` is concise human rationale. It explains intent, evidence,
decisions, trade-offs, and the next decision while referring to canonical
fields such as `plan.yaml#acceptance_criteria` and task IDs. It must not copy
the plan's lifecycle, acceptance, or verification authority.

Every task remains a bounded Markdown contract with YAML front matter. New or
repaired task output is not ready until the host-provided deterministic
front-matter, YAML, and `context-circuit.task` schema checks pass. Task status is
not a second approval gate, execution lease, or verification result.

## Optional specialist companions

The writer adds a specialist companion only when the overview records the
complexity signal that requires it. A companion adds rationale or detail; it
references canonical plan fields instead of copying them.

| Companion | Add when the plan has this signal |
| --- | --- |
| `requirements.md` | Multiple stakeholder outcomes or unresolved requirement boundaries need focused discussion. |
| `solution.md` | Multiple components, repositories, interfaces, or material design alternatives require a trade-off record. |
| `risks.md` | Irreversible, security, data-integrity, operational, or unresolved dependency risk needs explicit treatment. |
| `delivery.md` | More than one delivery target or an explicit branch, publication, deployment, or merge boundary needs explanation. |
| `acceptance.md` | Acceptance evidence has several scenarios that need a readable review map beyond the canonical list. |
| `verification.md` | Verification needs multiple independent checks, environments, or recovery paths beyond the canonical commands. |

The companion's headings should point to the relevant `plan.yaml` field and
task IDs. A companion never introduces a lifecycle value, replaces
`acceptance_criteria`, or changes `verification_commands`. An uncomplicated
plan omits all specialist companions; a complex plan may include only those
justified by its recorded signals.

Readers must continue to accept older bundles containing any or all of
`requirements.md`, `solution.md`, `risks.md`, `delivery.md`, `acceptance.md`,
and `verification.md`, as well as other existing plan companions. Historical
approved or done plans are read as authored and are not rewritten merely to
match the compact default.

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

## Task contract

Every task has a stable ID, bounded implementation paths, explicit
dependencies, acceptance criteria, verification commands, and stop conditions.
Task status is not a second approval gate or an execution lease. The
coordinator reconciles task projections in one idempotent operation when a plan
is approved, completed, or resumed, preserving metadata such as
`external_status`.

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
