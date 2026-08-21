# Planning

Plans define human-reviewed intended work. They are not the workspace session
itself and they do not replace runtime execution state. Every-session Product
Knowledge includes `context/WORKSPACE.md` for workspace identity and
`context/PROJECT.md` for the project being built.

A plan should have:

- one repository domain where code changes are expected;
- a clear objective and source;
- implementation scope and non-goals;
- Product Knowledge references;
- dependencies and connections;
- acceptance criteria;
- test scope and verification commands;
- risks, assumptions, and open questions;
- tasks with explicit dependencies and bounded scopes.

Plans begin as `draft` and require explicit human approval through
`cc-approve-plan`. `plan.yaml` is the canonical lifecycle record: approval
changes the included task projections from `draft` to `ready`, and completion
through `cc-finish-plan` changes them to `done`. Task status is not a second approval or execution gate. It does not provide an execution lease or
verification result. A one-plan request enters through `cc-run-plan`. Connected
already-approved plans enter through `cc-run-stack`. There is no scheduler.

`cc-approve-plan` is the named plan-approval skill. After explicit confirmation
in the current session it changes `plan.yaml` from `draft` to `approved` and
reconciles included tasks to `ready`. Approval does not claim a lease, create a
worktree, or start `cc-run-plan`. A prior `cc-review-plan` run is not required.
Already-approved plans are reported as approved and may enter `cc-run-plan`
or `cc-run-stack` without a second status rewrite. Invoking `cc-run-stack`
starts execution; there is no stack-approval gate. Draft member plans block
the stack run.

A plan has at most one active writing owner. The runtime contract uses an
exclusive plan lease and an exclusive worktree. Two writing sessions must never
share that worktree. Completion evidence, tests, Git state, or a verifier
handoff do not change canonical plan or task status. Missing evidence or a
failed verifier blocks the completion record; the canonical status remains
unchanged.

## Planning-only progressive context

Existing entry classification remains unchanged. After the planning route is
selected, planning may begin with the following request-scoped initial bundle;
this is a read order in the existing skills, not a generalized routing
framework, route catalog, manifest, schema, or payload format:

- Authority and universal safety: host and system instructions, AGENTS.md,
  WORKFLOW.md, and the source boundary they establish.
- Workspace and project identity: workspace.yaml, context/INDEX.md,
  context/WORKSPACE.md, and context/PROJECT.md.
- Accepted intent: the direct request or accepted Idea Brief or PRD, plus only
  evidence explicitly selected for the request. The source inbox is not scanned.
- Planning contract: cc-create-plan, plans/README.md, docs/planning.md,
  repository-local instructions, current branch and status, and the proposed
  implementation surface.
- Scoped planning state: the next plan ID, directly named or connected plan
  metadata, dependency status, archive eligibility, and only the current
  session or ownership records relevant to drafting safely.
- Product Knowledge navigation: the context index and project identity, with
  linked domain, role, workflow, architecture, convention, and decision pages
  retrieved only when needed.

The bundle must retain these required planning categories:

| Category | Initial information | Deeper retrieval when needed |
| --- | --- | --- |
| Authority | Instruction precedence, repository authority, source boundary, human gates | Complete workflow or runtime contract when the rule is unclear |
| Intent | Accepted outcome, scope, non-goals, decisions, uncertainty, selected evidence | Accepted Idea Brief/PRD and selected source or provenance records |
| Dependencies | Connected plan identity, lifecycle, dependency status, archive eligibility, blockers | Relevant plan.yaml, task, archive.yaml, dependents, handoff, and completion evidence |
| Product Knowledge | Index, workspace/project identity, and relevant-page navigation | Exact linked domain, role, workflow, architecture, convention, or decision page plus minimum repository evidence |
| Safety | Dirty-work preservation, lifecycle separation, contradiction stops, fallback, repository boundaries | Exact workflow/runtime contract sections when targeted reads do not resolve the rule |
| Ownership | Relevant session, lease, worktree, handoff, branch, and repository ownership metadata | Exact session, lease, stack, worktree, or handoff records and runtime-contract fields |

Retrieve deeper context before drafting for these concrete triggers:

- Intent is incomplete, disputed, source-based, or changed.
- A domain, role, workflow, decision, or implementation assumption is unresolved.
- The request names, replaces, depends on, archives, resumes, or conflicts with another plan.
- Safe drafting depends on live ownership, recovery, record shape, or worktree state.
- A safety or lifecycle rule remains ambiguous or contradictory after targeted retrieval.

If required information is missing, ambiguous, contradictory, or malformed
after targeted retrieval, use the existing complete-context fallback. Load the
complete workflow and runtime contracts plus all relevant Product Knowledge,
plan, task, archive, session, handoff, lease, worktree, repository, and
explicitly selected-source evidence required by current behavior. Never guess
to avoid the fallback. Non-planning routes, archive behavior, compact payloads,
and the current complete-context behavior remain unchanged.

## Archive eligibility and discovery

`archive.yaml` is an optional, versioned sidecar in a plan bundle. It owns
only whether ordinary routing may select the plan. It never changes the
canonical `plan.yaml` status or the task projection. A missing sidecar means
the plan remains active-compatible. The sidecar keeps an append-only sequence
of `archived` and `restored` events, each with plan identity, actor, timestamp,
reason, observed canonical status, and replacement references. The latest
event determines eligibility. Invalid, empty, or identity-mismatched sidecars
block selection and require a human recovery decision.

Only `cc-archive-plan`, after an explicit current-session `archive` gate,
appends archive or restore evidence. It never deletes a bundle, runtime
evidence, dirty work, or Git history; it never changes a plan or task status,
cascades to related plans, approves a plan, or acts as cleanup.

Ordinary discovery excludes a plan whose latest sidecar event is `archived`.
Session entry and `cc-whats-next` omit it from recommendations. `cc-approve-plan`,
`cc-run-plan`, `cc-run-stack`, and `cc-finish-plan` reject it even if named
explicitly. Explicit historical reads remain available at the unchanged plan
path. Before archive or restore, preflight checks plan identity and sidecar
shape, live lease and writing-session ownership, stack membership, worktree
ownership, dirty and unpushed work, dependencies, compatibility, and (on
restore) current branch and worktree state. Any ambiguity blocks the action.

An archived `done` dependency remains an inspectable done dependency. An
archived unfinished dependency (`draft` or `approved`) is unresolved and
blocks its active dependents; it never silently satisfies or disappears from a
dependency graph. Restoring is a fresh eligibility action, not a replacement
for approval or execution preflight.

On approval, completion, and session entry or resume, the coordinator
reconciles every included task to the projection expected by the plan. The
operation is bulk, idempotent, preserves task metadata such as
`external_status`, and does not rerun implementation or verification checks.
Runtime session state may say that a plan is being executed, blocked, or
awaiting review, but it must never silently change the canonical plan status.

For large projects, prefer several coherent plans by domain or repository
boundary rather than one unbounded plan. Multiple approved plans may execute
concurrently when their worktrees and ownership are distinct.

Approved single-plan execution enters through `cc-run-plan`. The root directs
a writer child and an independent verifier child. Sequential tasks share one
writer child and one worktree; independent plans use separate children and
worktrees. Overlapping paths are reported before merge or publication.
Standalone `cc-run-plan` creates or reuses the exclusive worktree from the
repository default or active branch.

`cc-run-stack` executes a connected set of already-approved plans in one root
session. It interprets existing prose `dependencies` into a frozen runtime
`graph.yaml` and tracks resume state in `progress.yaml` under
`.runtime/stacks/<stack-id>/`. Do not migrate those dependencies into a new
`plan.yaml` schema. Runtime must not override `plan.yaml`. There is no
`plans/<repository-key>-stacks/` layout and no durable `stack.yaml`.

Implemented is runtime: the writer finished, the independent verifier passed,
the worktree HEAD is committed and clean, and `completion.yaml` is
`ready-for-human-status-change`, while `plan.yaml` remains `approved`.
Dependents wait on implemented parents, not `done`. Freeze the parent SHA on
`progress.yaml`. A no-parent member uses the default or active branch. A
single-parent member is based on the parent frozen SHA. A multi-parent member
joins those SHAs in-run and does not wait for `default_branch`. When every
member is implemented, stop and hand the human the leaf worktrees. Do not
mark plans done from the stack.
