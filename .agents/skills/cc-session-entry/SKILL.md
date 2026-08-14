---
name: cc-session-entry
description: Enter or resume root and child agent sessions from durable workspace state.
---

# Context Circuit session entry

Use this skill when an agent starts or resumes work in a Context Circuit
workspace, or when a root session delegates a bounded child session. The
workspace filesystem is the coordination surface; do not make a user invoke a
CLI command as the entry point.

Read AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md, and the relevant
Product Knowledge before taking consequential action. Read
docs/agent-workspace-workflow.md for behavior and docs/runtime-contract.md for
runtime record fields and ownership rules.

## Identify the session

Determine whether the session is:

- a root session owning the human request;
- a child session with a parent and root session; or
- a read-only verifier/reviewer.

Never infer a child assignment from conversation history. Require a delegation
packet with explicit scope, permissions, plan/task, acceptance criteria, stop
conditions, and handoff schema. If the packet or session record is absent,
report a blocker.

## Root entry and resume

1. Inspect .runtime/sessions/ and identify only records relevant to the
   request. There is no global current-session file.
2. Inspect active plans, dependencies, leases, worktrees, and handoffs.
3. Report the route: orienting, gathering-context, planning,
   awaiting-approval, executing, verifying, blocked, or handoff.
4. For an approved plan without a live owner, claim execution only within its
   lease and exclusive worktree. Different plans may proceed concurrently.
5. Resume from the durable session record and latest handoff, not from assumed
   conversation state.

For fresh work, gather source evidence and draft Product Knowledge or a plan
as appropriate. Do not silently approve plans or change canonical statuses.

## Child entry and handoff

Before working, verify that the child packet matches the session record and
contains:

- session_id, parent_session_id, root_session_id, kind, and role;
- objective, scope, non-goals, plan/task, and context references;
- repository/worktree, write permissions, acceptance criteria, and stop
  conditions;
- the expected handoff schema.

Work only in the assigned worktree. Do not modify another session's runtime
records, plan/task statuses, wrapper context, publication data, or external
systems. A verifier is read-only.

At handoff, report:

- session and parent IDs;
- objective, delegated scope, and evidence inspected;
- decisions, assumptions, changed files, and tests/verification;
- questions, blockers, limitations, and next action.

If interrupted, preserve the worktree and runtime records. A new session may
resume or take over only with explicit ownership evidence; never overwrite a
live lease or silently replace a stale session.
