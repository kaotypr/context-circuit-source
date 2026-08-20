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

## Bundle generation

New plan generation starts with the smallest compatible bundle:

1. `plan.yaml` records canonical lifecycle, repository scope, dependencies,
   acceptance criteria, test scope, and verification commands.
2. `overview.md` records concise intent, evidence, decisions, trade-offs, and
   the next decision, using references such as
   `plan.yaml#acceptance_criteria` instead of copied authority.
3. `tasks/` contains bounded Markdown task contracts with valid YAML front
   matter and the `context-circuit.task` schema.

The writer records a complexity signal in the overview before adding a
specialist companion. Use `requirements.md` for multiple stakeholder outcomes
or unresolved requirement boundaries; `solution.md` for multi-component,
multi-repository, or material design alternatives; `risks.md` for
irreversible, security, data, operational, or unresolved dependency risk;
`delivery.md` for multiple delivery targets or explicit branch/publication/
deployment boundaries; `acceptance.md` for several acceptance scenarios that
need a review map; and `verification.md` for multiple independent checks,
environments, or recovery paths. These files explain rationale and point back
to canonical plan fields; they do not restate or replace them.

An uncomplicated plan omits specialist companions. A reader must still accept
historical bundles with the full companion set or other existing companions.
Historical approved and done plans are not rewritten as part of compact
generation.

Before a new or repaired task is presented as ready, the active host applies
the approved deterministic validation capability: opening and closing
front-matter delimiters, YAML parsing, required fields, enums, nested shapes,
and the task schema. Repair reruns the complete sequence. This validation
boundary does not approve a plan, start execution, produce verification
evidence, or complete a plan.

## Prototype comparison and duplication evidence

The bounded prototype evidence in
`test/fixtures/document-system/plans/prototype-comparison.md` compares a
modeled full companion bundle with the compact default. Its inventory counts
9 before files versus 3 after files, with 2/3/3 copied lifecycle, acceptance,
and verification surfaces before and zero copies after. The companion
`duplication-inventory.yaml` records the same counts as fixture evidence.

These counts measure duplicated authority surfaces, not the amount of rationale
that may be removed. The after prototype retains one human-readable overview,
bounded tasks, and a safety record requiring rationale retention and
complexity-triggered specialist companions. A plan is not considered better
merely because it has fewer files.

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
