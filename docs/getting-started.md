# Getting started as an AI agent

Context Circuit is a workspace for AI agents, not a command console for users.
The user should be able to say:

> Start or resume work in this workspace.

Durable shared context lives in ordinary Markdown and YAML: Product Knowledge,
decisions, source provenance, plans, and task definitions. Private runtime
state lives under `.runtime/`: session records, delegation handoffs, plan
leases, prompts, and isolated worktrees. Runtime state is resumable but is not
authoritative over context, plans, human decisions, or repository
instructions. Preserve it until a human chooses cleanup through
`cc-cleanup-runtime`. That skill inspects every runtime worktree for
uncommitted and unpushed work, stops with a confirmation list when risk
exists, and deletes `.runtime/` only after the human chooses cleanup.
Cleanup is workspace-wide for `.runtime/` and does not delete Git branches or
modify the base checkout.

A writing session modifies only its assigned exclusive worktree. Different
plans may use different worktrees concurrently. Research and verification
sessions are read-only unless their delegation explicitly grants write access.

Agents may inspect, draft, test, delegate, create isolated worktrees, and
implement approved scope. Humans control Product Knowledge acceptance, plan
approval through `cc-approve-plan`, material scope changes, merge, publication,
deployment, completion through `cc-finish-plan`, runtime cleanup through
`cc-cleanup-runtime`, and ambiguous session takeover.

Host integrations should expose the same workspace entry and delegation
behavior. Host-specific commands are adapters and must not become a second
workflow or source of truth.

The workspace is instruction- and filesystem-driven. It has no Node or
JavaScript command layer. Verify behavior with `sh test/acceptance.sh`.

## 1. Enter and orient

Read AGENTS.md, WORKFLOW.md, workspace.yaml, context/index.md,
context/WORKSPACE.md, context/PROJECT.md, relevant Product Knowledge, active
runtime sessions, plans, repository instructions, and observed Git state.

Identify whether this is:

- a new workspace or project;
- a new PRD or source;
- an unfinished root session;
- a delegated child session;
- an approved plan ready to execute;
- a verification, review, or blocked session.

Explain the selected route before consequential action.
For a fresh root route, create a root session record with an explicit session ID
before delegating or claiming work. For resume, use the existing session record
and latest handoff; never infer a current session from conversation history.

Choosing the next action is part of this root entry, not a separate command.
Inspect approved plans, declared dependencies, active leases and session
ownership, source freshness, repository cleanliness, worktrees, blockers, and
pending human gates. Recommend or claim only work that is dependency-ready,
explicitly scoped, and not already owned by another writing session. When no
work is executable, explain whether the session needs context, a draft plan,
human approval, a review, or a decision about a blocker.

## 2. Initialize workspace identity when needed

Initialization asks one compact set of core identity questions:

1. Is this a solo or team workspace?
2. Which repositories or project items are known, if any? A zero-repository
   answer is valid for a new project or idea.
3. For each repository, what is its role and which branch should be the default
   active branch? If `development` exists, recommend it; let the user choose a
   different branch.

Record the confirmed identity in `workspace.yaml` and summarize it in
`context/WORKSPACE.md`. Do not write workspace identity into
`context/PROJECT.md`. A repository entry has a path, mode, role, agent, and
`default_branch`. Do not turn the current working branch into a delivery
policy.

When no repository is registered, offer the user a short Idea Brief
conversation, selected source intake, both, or deferment. Do not treat the
no-repository path as an error or as a reason to invent a repository.

Initialization stops after core identity. It does not ask about delivery,
commits, pushes, merges, publication, deployment, or external activity tools;
those are optional later configuration.

## 3. Capture the right artifact

Choose the smallest useful next artifact. Do not force every request through
an Idea Brief and then a PRD.

- Capture an Idea Brief for uncertain intent.
- Capture a PRD for durable requirements.
- Use a selected source or repository evidence when that is already enough.

When the user provides a PRD, source document, or repository evidence, read
that evidence and keep it as the input for Product Knowledge. Do not invent a
missing Idea Brief or PRD.

When intent is still uncertain, draft an Idea Brief. An Idea Brief is enough
to continue; do not require a PRD after it.

If no useful source, repository, or brief exists, ask focused discovery
questions before creating Product Knowledge or a plan.

## 4. Ground or refresh Product Knowledge

When the chosen artifact or evidence is enough:

1. Separate stated requirements, current behavior, assumptions, unknowns, and
   evidence.
2. Draft concise Product Knowledge with source references.
3. Present the draft for human confirmation.
4. Do not silently replace accepted context.

## 5. Draft and approve a plan

When context is sufficient:

1. Propose one or more plans by coherent domain or repository boundary.
2. Include dependencies, acceptance criteria, implementation scope, test scope,
   and verification commands.
3. Present plans as drafts.
4. Optionally review with `cc-review-plan`. Review stays read-only.
5. Wait for explicit human approval through `cc-approve-plan`.
6. Do not infer approval from conversation tone, tests, or agent output.
   Approval does not start execution.

## 6. Execute through `cc-run-plan`

After approval, `cc-run-plan` directs the root to:

- claim a plan execution and create or reuse its exclusive worktree from the
  repository default or active branch;
- spawn a writer child through the host child-session primitive;
- spawn a later independent verifier child (`write_worktree: false`);
- record completion.yaml and ask for `cc-finish-plan`.

Sequential tasks share one writer child. Independent plans get separate
children and worktrees. If the host cannot spawn a child, report the missing
host primitive to the human.

Every child session receives a bounded delegation packet and returns a structured
handoff to its parent.

A one-plan request still enters `cc-run-plan`. Connected approved unimplemented
plans, or an interrupted `.runtime/stacks/<stack-id>/` run, enter
`cc-run-stack`. Refuse to treat a stack run as one `cc-run-plan`. Invoking
`cc-run-stack` starts or resumes execution; there is no stack-approval gate
and no scheduler. `cc-run-stack` is not a second way to run one plan.

## 7. Resume from the handoff

At the end of each session, record:

- session and parent identifiers;
- objective and current route;
- evidence and decisions;
- changed files and tests;
- blockers and open questions;
- next safe action.

The next session reads this handoff instead of relying on conversation history.
Human review and explicit status changes remain required before work is
considered complete. A completion evidence record may request `cc-finish-plan`,
but it must not change plan or task status itself. After a plan is `done`, or
when the user asks to clear local execution state, `cc-cleanup-runtime` may
delete `.runtime/` only after inspecting dirty or unpushed work and receiving
explicit confirmation.

## 8. Roll out document-system compatibility safely

When adopting the document-system contracts, use this order:

1. declare the participating OKF bundle roots and reserved files;
2. deploy readers that accept canonical `type` plus legacy `kind` and preserve
   unknown fields;
3. switch live writers and templates to canonical `type`, list-shaped
   `verified`, and lowercase `index.md`;
4. validate newly generated and changed knowledge, plan/task, skill, and
   runtime artifacts at their owning schema boundary;
5. migrate live templates and current accepted knowledge;
6. keep historical approved artifacts and Markdown-only handoffs readable
   through explicit compatibility routes;
7. collect route, duplication, schema, resume, and human-review evidence;
8. propose fallback removal only after the supported-workspace migration
   evidence is complete.

The compatibility window covers legacy `kind`, mapping-shaped `verified`,
historical index names, and Markdown-only handoffs. A removal proposal requires
no supported workspace dependency, a historical safe-reading route, a complete
regression run, and a separate human `compatibility-window-decision`. No
fallback is removed by passing tests alone.

The `metric-threshold-decision`, `status-change`, `publication`, and release
gates remain separate and unsatisfied until a human explicitly decides them.
Conformance evidence, a local commit, or a staged artifact is not permission
to merge, publish, tag, release, or mark a plan done.

## Which document

Use this map only to find a retained contract. It is not a second workflow.

| Need | Read |
| --- | --- |
| Session behavior | `docs/agent-workspace-workflow.md` |
| Runtime records | `docs/runtime-contract.md` |
| Idea Brief / PRD | `docs/idea-brief.md`, `docs/prd.md` |
| Product Knowledge | `docs/product-knowledge.md` |
| Plan artifact | `plans/README.md` |
| Plan lifecycle | `docs/planning.md` |
| Plan review | `docs/plan-review.md` |
| Optional config | `docs/configuration.md` |
