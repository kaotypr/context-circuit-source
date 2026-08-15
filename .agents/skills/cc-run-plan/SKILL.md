---
name: cc-run-plan
description: Execute an approved plan through bounded solo or delegated sessions, exclusive worktrees, verification, and durable handoffs.
---

# Run an approved plan

Use this skill when the user explicitly asks to execute an approved plan or
when root-session routing identifies an approved, dependency-ready plan.
`cc-run-plan` is the sole standard plan-execution capability. There is no
user-facing `cc-run-task` workflow.

## Preflight

Read the canonical `plan.yaml`, all task contracts, relevant Product Knowledge,
repository-local instructions, current repository state, active sessions,
leases, worktrees, and the latest handoffs. Confirm:

- the plan status is exactly `approved`;
- every declared dependency is approved or done as required by the plan;
- task projections are reconciled to `ready` without rerunning checks;
- the repository base is clean and the target worktree is available;
- no other writing session owns the plan or target worktree;
- the requested work remains within scope and its assumptions are safe.

An unapproved, unknown, contradictory, or ownership-conflicted plan is
blocked. Do not create a hidden replacement plan, silently steal a lease, or
change canonical status to make execution possible.

## Claim and route

The root coordinator creates or resumes an explicit session record, atomically
claims `.runtime/plans/<plan-id>/lease.lock/`, records the owner and lease,
then creates or reuses the exclusive repository worktree and branch. The
lease, session, and worktree must name the same plan, session, repository, and
path.

Choose the shortest safe route:

- **Solo:** perform a small bounded task directly in the plan worktree while
  recording the same evidence, scope, verification, and handoff fields.
- **Delegated:** create an explicit child packet for each bounded writer or
  verifier. Packets include identity, objective, plan/task, paths, non-goals,
  context references, repository/worktree, permissions, acceptance criteria,
  stop conditions, and `session-handoff-v1`.

Delegated writers use separate worktrees when they write concurrently. A
verifier is independent and read-only for implementation, plan, lease,
worktree, and activity state; it may write only its own handoff.

## Execute and recover

Run dependency-ready tasks in order. Preserve task evidence, questions,
blockers, dirty work, and handoffs in durable records. A worker stops and
returns to the root when it reaches an unassigned path, a scope change, a
contradictory source, an ownership conflict, or missing required evidence.

On interruption, resume from the explicit session record and latest handoff.
Inspect heartbeat, lease, worktree, and dirty state first. A stale or missing
record requires a visible recovery decision; it is not permission to invent
progress or silently take over. A human-authorized takeover creates a new
session that names the replaced session, reason, and preserved evidence.

## Verify and hand off

Run the plan's verification commands and the independent verifier's checks.
Failed verification remains failed and blocks completion until repaired within
scope or escalated for human judgment. Do not merge, publish, deploy, or mark
the plan done from this skill.

Before requesting completion, record evidence for every task, a passing
independent verifier handoff, remaining limitations, and the `status-change`
human gate in `.runtime/plans/<plan-id>/completion.yaml`. That evidence record
does not change `plan.yaml`; the human controls final plan and task status.

## Output

Report the session and plan, route chosen, worktree and ownership evidence,
tasks executed, files changed, tests and verification, decisions, assumptions,
blockers, and next safe action. If completion is ready, ask for the human
status-change gate rather than implying that tests or Git state completed the
plan.
