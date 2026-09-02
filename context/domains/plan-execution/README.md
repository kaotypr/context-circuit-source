---
kind: domain
status: accepted
title: Plan execution
slug: plan-execution
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - One worker executes one intent-authorized plan in one bounded execution.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - docs/getting-started.md
---

# Plan execution

## Summary

Implementing an intent-authorized plan with one bounded worker. Execution's
authorization derives from the approved intent's scope envelope (re-checked at
execution start), not a separate plan-approval step. Route "execute plan `<id>`",
repair, and resume/recovery requests here. Owned by the `cc-execute` skill and the
`agents/worker.md` worker role.

## Scope

Inside: execution begin from the base tip, isolated worktrees and branches,
per-repository worker commits, the worker handoff (a claim), the exclusive-create
ownership lock, the three-failure repair counter, and preserved, resumable
runtime records.

Outside: authorization itself ([plan-authorization](../plan-authorization/README.md),
derived from the intent's approved envelope), the read-only check
([verification](../verification/README.md)), completion
([completion](../completion/README.md)), and any pull request/merge/push
([delivery](../delivery/README.md)).

## Behavior

Only a plan within its approved intent's scope envelope may execute (INV-EXEC-01).
One worker executes every task of one intent-authorized plan in one bounded
execution, in dependency order, across all mapped repositories (INV-EXEC-02).
Execution creates exactly one deterministic
branch `cc/<plan-id>/<repo-id>` and one isolated worktree per affected
repository, from the captured base-branch tip; the base checkout is never
written (INV-EXEC-03).

A plan with **same-repo predecessors** (see
[run-stack](../run-stack/README.md)) begins base-aware: its worktree is prepared
on the predecessor branch (stack) or a runtime-authored integration merge, and the
repository record gains `based_on`. A plan with **no** dependency is unchanged — it
still branches from the captured base tip (INV-EXEC-03). Execution also discovers
the target repository's own agent guidance and delivers a grounded, runtime-
assembled brief rather than a free-composed one (see
[repository-grounding](../repository-grounding/README.md)). The plan contract uses
`schema_version: 3` for every v1.0 plan, including plans without dependencies;
legacy plan shapes are not accepted.

The worker commits each changed repository before verification; every repair
creates a new commit and a prior commit is never amended to conceal a repair
(INV-EXEC-04). At most one active worker owns a plan execution, enforced by an
atomic exclusive-create lock; a competing worker gets a read-only or blocked
result and a live lock is never silently stolen (INV-OWN-01).

The worker-failure counter increments on each verifier rejection including the
initial implementation, to a maximum of three, after which execution stops
(INV-REPAIR-01). Failure, interruption, or blocking preserves branches,
worktrees, commits, handoffs, verifier evidence, repair attempts, and runtime
records; failed work is never silently cleaned up (INV-PRESERVE-01). Runtime
records are written atomically; a partial or contradictory record cannot grant
ownership, resume a worker, prove verification, or authorize completion
(INV-RUNTIME-01, INV-RUNTIME-02).

## Execution brief, snapshot, and repair discipline

Before the worker runs, execution preflights the affected repositories (the
intent-envelope check, dependency order, clean base checkouts, bounded scope) and
generates a bounded
execution brief plus an immutable plan-snapshot — a serialized copy of `PLAN.md`,
`plan.yaml`, and all task files that the worker and verifier prompts read.

Repair is disciplined: a repair may not redesign the plan. A repair continues in
the same execution only when the intended scope and acceptance are unchanged; if
repair requires new scope, execution stops and the human is asked to change the
plan, and the old-execution → new-revision relationship is recorded.

Resume is allowed only when the plan revision, the configured base branches and
captured base commits, the worktree paths, and ownership all still match.

## Workflows

- Approve the intent and execute, then inspect results: `docs/getting-started.md`

## Interfaces

- Human request: "Execute plan `<id>`" / "Repair the failed `<id>` verification"
- Execution branch: `cc/<plan-id>/<repo-id>`; worktrees under `.runtime/worktrees/…`
- Records: `execution.yaml`, `handoff.md`

## Data

`execution.yaml` holds `execution_id`, a `plan_revision` digest, `owner`,
`status` (`running`/`verifying`/`repairing`/`verified`/`failed`/`blocked`),
`worker_failures` 0..3, and per-repository worktree/branch/base/latest records
(a dependent plan's record also carries `based_on`). Each execution records a
grounding manifest under `grounding/<repo>.yaml`.
The worker handoff is a claim (commits, tasks done, checks, assumptions, limits,
suggested verifier focus, and whether a repair is being performed), not
verification proof.

## Constraints and edge cases

The worker never edits the base checkout, changes completion status,
marks its own work verified, claims independent verification, rewrites or accepts
Product Knowledge, expands scope silently, or merges/pushes/publishes.

## Implementation references

- `.agents/skills/cc-execute/SKILL.md`, `agents/worker.md`
- `wrapper/runtime/engine.sh`: `cc_execution_begin` (base-aware for dependents),
  `cc_execution_next_id`, `cc_exec_set`, `cc_attempt_begin`,
  `cc_worker_commit_record`, `cc_worker_handoff_record`, `cc_lock_acquire`,
  `cc_lock_owner`, `cc_lock_release`, `cc_repair_allowed`, `cc_recovery_inspect`;
  base selection and grounding: `cc_base_prepare`, `cc_discover_repo_grounding`,
  `cc_worker_brief_assemble` (see the run-stack and repository-grounding domains)
- `wrapper/contracts/schemas/execution.yaml`,
  `wrapper/contracts/schemas/worker-handoff.yaml`
- `wrapper/contracts/invariants.yaml`: INV-EXEC-01, INV-EXEC-02, INV-EXEC-03,
  INV-EXEC-04, INV-OWN-01, INV-REPAIR-01, INV-PRESERVE-01, INV-RUNTIME-01,
  INV-RUNTIME-02

## Verification

`sh test/acceptance.sh` (execution suite). Independent verification of results
is owned by [verification](../verification/README.md).

## Provenance

Authored from the current wrapper at HEAD `4b8ac0b`. Raw `sources/` was not
scanned. This is not the deleted previous-version "engine-owned-runtime-records"
plan; it describes the runtime records that ship today.

## Acceptance notes

Accepted 2026-08-24 from proposal `0011-domain-plan-execution`. Extended
2026-08-27 from proposal `0018-change-plan-execution` (v0.6 base-aware begin for
dependents, the grounded brief, and plan `schema_version [1, 2]`).
