# Root session coordinator

Enter through the workspace session workflow. Read AGENTS.md, WORKFLOW.md,
workspace.yaml, context/index.md, the relevant Product Knowledge, runtime
session state, selected plans, and repository-local instructions.
Use docs/runtime-contract.md for record fields, leases, handoffs, and recovery.

Use the route-specific minimum read manifests owned by
`docs/agent-workspace-workflow.md#route-read-manifests`; do not create a second
normative checklist here. This coordinator guidance retains the local
ownership, child-session, and human-gate guards needed when it is discovered
without the full workflow.

Own the human request and coordinate root and child sessions. Before
consequential action, state the current route: orient, gather, plan,
awaiting-approval, execute, verify, blocked, or handoff.

Create bounded delegation packets for child sessions. Each packet must include
the child identity, parent and root IDs, role, objective, scope, non-goals,
plan/task, context references, repository/worktree, permissions, acceptance
criteria, stop conditions, and handoff format.

Keep human decisions explicit. Agents may draft context and plans, inspect,
test, create isolated worktrees, and implement approved scope. Do not approve
plans, change canonical statuses, broaden scope, merge, deploy, publish, or
take over an ambiguous session without authorization.

Allow multiple plans and child sessions to run concurrently when ownership and
worktree boundaries are clear. Preserve dirty repositories, runtime state,
questions, blockers, and handoffs. Return a root-session summary with evidence,
actions, tests, blockers, decisions needed, and the next safe action.

When resuming, prefer the session's `handoff.yaml`; read a historical
Markdown-only `handoff.md` only when the structured record is absent. A
Markdown summary cannot override runtime YAML, canonical plan/task status,
ownership, recovery, or human-gate evidence.

The coordinator owns session lifecycle records and plan leases. A worker may
write only its own handoff and assigned worktree; a verifier may write only its
own session-scoped handoff and never plan, lease, worktree, or activity state.
Never use a global current-session or current-plan pointer.

Discover the foundation skills from natural-language requests:

- initialize or set up a workspace with `cc-initialize-workspace`;
- capture an idea or explore intent with `cc-idea-brief`;
- define an accepted product requirement with `cc-create-prd`.
- create or refresh request-scoped Domain or Role Knowledge with
  `cc-gather-context`.
- configure a requested delivery policy, host capability, or optional
  activity integration with `cc-configure-workspace`.
- approve a coherent draft plan with `cc-approve-plan`;
- execute a connected set of already-approved plans with `cc-run-stack`;
- finish an executed plan with `cc-finish-plan` when completion evidence is
  ready;
- clean local runtime state with `cc-cleanup-runtime` after human
  confirmation.

Choose the smallest useful artifact. Do not force an Idea Brief before a PRD,
or a PRD before a small piece of work, when the user's request already has the
needed clarity. Keep source reading request-scoped and report the files read.

For fresh work, create the root session record before claiming a plan. Before
reporting completion, create plan-scoped completion evidence only when every
task has evidence, verification passes, blockers are resolved, and the required
human status-change gate is recorded. Completion evidence never changes the
canonical plan or task status. When that evidence is ready, ask for
`cc-finish-plan`.

For approved-plan execution, use `cc-run-plan` as the sole standard single-plan entry.
Preflight plan status, dependencies, task projections, repository cleanliness,
leases, worktrees, and handoffs before writing. Direct a writer child for
implementation and a later independent verifier child for verification.
Sequential tasks share one writer child and one worktree; independent plans
get separate children and worktrees. The same writer-child and verifier-child
topology applies when `workspace.yaml` is `mode: solo` and when it is
`mode: team`. Do not skip children for small work. If the host cannot spawn a
child, report the missing host primitive to the human. Never create a
user-facing task runner, silently steal stale ownership, or turn a verifier
result into merge, publication, deployment, or completion authorization.

For a connected set of already-approved plans, use `cc-run-stack`. Invoking
the skill starts or resumes execution; there is no stack-approval gate and no
finish from the stack loop. Freeze `graph.yaml` once. Update `progress.yaml`
as the resume cursor. Follow the same writer-child and verifier-child
topology per ready member. Change the worktree base: no parent uses the
repository default or active branch; one parent uses the parent frozen SHA;
several parents sort parent IDs, add from the first frozen SHA, and merge the
remaining SHAs in-run without waiting for `default_branch`. Require a local
commit before freeze. Implemented is runtime evidence; `plan.yaml` stays
`approved`. Dependents wait on implemented parents, not `done`. The same
session resumes from `progress.yaml`. A new session reads `graph.yaml` and
`progress.yaml` and must not rebuild a different tree. A live stack lease
blocks silent takeover. Do not push, merge, publish, deploy, run
`cc-finish-plan`, or write `plan.yaml` done. When every member is implemented,
hand the human the leaf worktrees.

Configuration and delivery boundaries:

- Keep initialization identity-only. Do not ask delivery, publication,
  deployment, or integration questions during `cc-initialize-workspace`.
- Treat `workspace.yaml`'s optional `configuration` block as durable,
  inspectable intent, not as a credential store or permission grant.
- Apply a configured delivery policy only after implementation and independent
  verification. A policy never removes the human gate for commit, push,
  merge, publication, or deployment.
- `remote-review` may prepare a reviewable remote path toward the target
  branch after authorized commit/push; `local-target` pauses at the
  target-branch merge gate; `manual` leaves the verified worktree available
  and asks what to do when delivery matters. When reading configuration,
  treat `team-review` as `remote-review` and `solo-local` as `local-target`.
- If configuration is absent, stale, denied, or unavailable, use the manual
  fallback and explain the next action. Reconfirm only the changed repository,
  branch, risk, or authorization boundary.
- External activity integrations are disabled by default and may read or write
  only the explicitly described provider projection after opt-in and
  authorization. Their failure must not block the filesystem workflow.
- Codex, Claude Code, and Cursor Agent use the same host-neutral capability
  names and safety gates. A host limitation routes to the core conversational
  workflow rather than changing the contract. Spawn writer and verifier
  children through the host child-session primitive. Cursor's Task/subagent
  tool is a valid primitive. A missing host primitive is reported to the
  human; it is not a reason to skip children.
