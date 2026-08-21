---
name: cc-session-entry
description: Enter or resume root and child agent sessions from durable workspace state.
---

# Context Circuit session entry

Use this skill when an agent starts or resumes work in a Context Circuit
workspace, or when a root session delegates a bounded child session. The
workspace filesystem is the coordination surface; do not make a user invoke a
CLI command as the entry point.

Read AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md,
context/WORKSPACE.md, context/PROJECT.md, and the relevant Product Knowledge
before taking consequential action. Read docs/agent-workspace-workflow.md for
behavior and docs/runtime-contract.md for runtime record fields and ownership
rules. The planning-only read sequence below may defer planning-irrelevant
parts of those complete references, but only after normal entry classification
has selected the planning route; the complete current behavior remains the
fallback whenever targeted reads are insufficient.

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

## Plan-aware routing

### Planning-only progressive reads

Existing entry classification is unchanged. Apply this read sequence only
after it selects planning, including direct plan drafting or a request to
create a plan. It is a request-scoped read order in the existing planning
surfaces, not a new route manifest, catalog, schema, payload, or generalized
routing framework. All other routes retain their current reads and behavior.

Start with the initial planning bundle:

- Authority and universal safety: host and system instructions, AGENTS.md,
  WORKFLOW.md, and the source boundary they establish.
- Workspace and project identity: workspace.yaml, context/INDEX.md,
  context/WORKSPACE.md, and context/PROJECT.md.
- Accepted intent: the direct request or accepted Idea Brief or PRD, plus only
  evidence explicitly selected for the request. Do not scan the source inbox.
- Planning contract: cc-create-plan, plans/README.md, docs/planning.md, the
  current repository-local instructions, branch and status, and the proposed
  implementation surface.
- Scoped planning state: the next plan ID, directly named or connected plan
  metadata, dependency status, archive eligibility, and only the current
  session or ownership records relevant to drafting safely.
- Product Knowledge navigation: the context index and project identity above,
  with linked domain, role, workflow, architecture, convention, and decision
  pages retrieved only when needed.

The initial bundle must preserve six categories: authority, intent,
dependencies, relevant Product Knowledge, safety, and scoped ownership. Read
the exact deeper evidence when a trigger is present:

- Intent is incomplete, disputed, source-based, or changed. Retrieve the
  accepted Idea Brief or PRD and only the selected source or provenance
  records.
- A domain, role, workflow, decision, or implementation assumption is unresolved.
  Retrieve the exact linked Product Knowledge page and the minimum repository
  evidence needed to resolve it.
- The request names, replaces, depends on, archives, resumes, or conflicts with another plan.
  Retrieve the directly relevant plan.yaml, task, archive.yaml, dependent,
  handoff, and completion evidence.
- Safe drafting depends on live ownership, recovery, record shape, or worktree state.
  Retrieve the exact session, lease, stack, worktree, or handoff record and the
  relevant runtime-contract section.
- A safety or lifecycle rule remains ambiguous or contradictory after targeted retrieval.
  Retrieve the complete agent-workspace workflow and runtime contract.

If any required authority, intent, dependency, Product Knowledge, safety, or
ownership fact remains missing, ambiguous, contradictory, or malformed after
targeted retrieval, use the existing complete-context fallback. That fallback
loads the complete workflow and runtime contracts plus all relevant Product
Knowledge, plan, task, archive, session, handoff, lease, worktree, repository,
and explicitly selected-source evidence required by current behavior. Never
guess to avoid the fallback.

When the request names a plan, inspect its canonical `plan.yaml`, task
projections, declared dependencies, active lease, assigned worktree, and latest
handoff. Also inspect the optional `archive.yaml`: absence is active-compatible;
a latest `archived` event excludes the plan from ordinary routing while leaving
its original path readable for history. A malformed archive record is blocked,
not inferred. An approved dependency-ready active plan with no live writing owner routes to
`cc-run-plan`; a connected set of approved unimplemented plans, or an
interrupted `.runtime/stacks/<stack-id>/` run, routes to `cc-run-stack`; a
coherent draft routes to `cc-approve-plan` when the user is
asking to approve, or otherwise to review; an approved plan with ready
completion evidence routes to `cc-finish-plan`; a human cleanup request routes
to `cc-cleanup-runtime`; a plan with a live owner routes to resume or
coordination; and a contradictory, stale, or ambiguous record routes to a
visible recovery decision. Refuse to treat a stack run as one `cc-run-plan`.
Invoking `cc-run-stack` starts or resumes execution; there is no stack-approval
gate.

Do not recommend, approve, execute, finish, or include an archived plan in a
stack. When a human asks to change its eligibility, route to
`cc-archive-plan`; archive and restore each require a separate explicit human
confirmation and do not change canonical lifecycle status.

Approved-plan execution through `cc-run-plan` directs a writer child and a
later independent verifier child. Do not treat small or sequential work as a
root-implements shortcut. Sequential tasks share one writer child and one
worktree. Spawn children through the host child-session primitive; Cursor's
Task/subagent tool is a valid primitive. If the host cannot spawn a child,
report the missing host primitive to the human and ask how to proceed.

On approval or resume, reconcile included task status in bulk from the plan
status. The reconciliation is idempotent metadata repair and never reruns
implementation or verification. Keep the plan status and task projection
separate from session execution status.

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
