---
name: cc-run-plan
description: Execute an approved plan through directed writer and verifier children, exclusive worktrees, verification, and durable handoffs.
---

# Run an approved plan

Use this skill when the user explicitly asks to execute an approved plan or
when root-session routing identifies an approved, dependency-ready plan.
`cc-run-plan` is the sole standard single-plan execution capability. It is not
a stack runner. Connected approved plans enter `cc-run-stack`.

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
then creates or reuses the exclusive repository worktree from the repository
default or active branch. The
lease, session, and worktree must name the same plan, session, repository, and
path.

After preflight, the expected execution records are a writer child packet and
a later independent verifier child packet. Preferred prompt `route` is
`delegated`. A prompt may still say `writer-count: 1` when sequential tasks
share one worktree. Sequence is not a reason to skip children.

```text
Root session
  -> claim lease and create exclusive worktree
  -> writer child (host subagent, write_worktree true)
  -> verifier child (different host subagent, write_worktree false)
  -> root records completion.yaml and asks for cc-finish-plan
```

For more than one approved dependency-ready plan:

```text
Root session
  -> plan A writer + plan A verifier
  -> plan B writer + plan B verifier
  -> report overlapping paths before merge or publication
```

- Create an explicit writer child packet with `role: implementer` and
  `write_worktree: true`. Sequential tasks in one plan go to one writer child
  and one worktree.
- Create a later independent verifier child packet with `role: verifier` and
  `write_worktree: false`. Independent verification is a different child, not
  a verifier session record plus the root re-running checks.
- Independent approved dependency-ready plans go to separate children and
  worktrees; they may run concurrently. Overlapping paths are reported, not
  absorbed into one root implementation pass.

The same writer-child and verifier-child topology applies when
`workspace.yaml` is `mode: solo` and when it is `mode: team`. Workspace mode
does not select an execution topology.

Spawn each child through the host child-session primitive. The filesystem
packet remains the coordination record; the host supplies the child execution
context. Cursor's Task/subagent tool is a valid child-session primitive. If
the host cannot spawn a child, the missing host primitive is reported to the human;
ask how to proceed. Do not quietly skip children.

Packets include identity, objective, plan/task, paths, non-goals, context
references, repository/worktree, permissions, acceptance criteria, stop
conditions, and `session-handoff-v1`.

Delegated writers use separate worktrees when they write concurrently. A
verifier is independent and read-only for implementation, plan, lease,
worktree, and activity state; it may write only its own handoff.

## Execute and recover

Run dependency-ready tasks in order through the writer child. Preserve task
evidence, questions, blockers, dirty work, and handoffs in durable records. A
worker stops and returns to the root when it reaches an unassigned path, a
scope change, a contradictory source, an ownership conflict, or missing
required evidence.

On interruption, resume from the explicit session record and latest handoff.
Inspect heartbeat, lease, worktree, and dirty state first. A stale or missing
record requires a visible recovery decision; it is not permission to invent
progress or silently take over. A human-authorized takeover creates a new
session that names the replaced session, reason, and preserved evidence.

## Verify and hand off

Direct independent verification to the verifier child and run the plan's
verification commands. Failed verification remains failed and blocks
completion until repaired within scope or escalated for human judgment. Do
not merge, publish, deploy, or mark the plan done from this skill.

Before requesting completion, record evidence for every task, a passing
independent verifier handoff, remaining limitations, and the `status-change`
human gate in `.runtime/plans/<plan-id>/completion.yaml`. That evidence record
does not change `plan.yaml`. When completion evidence is ready, ask for
`cc-finish-plan` rather than an unnamed status-change request. The human
controls final plan and task status through that named skill. This skill
must not mark the plan done.

## Output

Report the session and plan, delegated route, worktree and ownership evidence,
tasks executed, files changed, tests and verification, decisions, assumptions,
blockers, and next safe action. If completion is ready, ask for `cc-finish-plan`
rather than implying that tests or Git state completed the plan.
