# Getting started as an AI agent

Context Circuit is entered through a host agent session. The user should be
able to say:

> Start or resume work in this workspace.

## 1. Enter and orient

Read AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md, relevant
Product Knowledge, active runtime sessions, plans, repository instructions,
and observed Git state.

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

## 2. Initialize the workspace

Initialization asks one compact set of core identity questions:

1. Is this a solo or team workspace?
2. Which repositories or project items are known, if any? A zero-repository
   answer is valid for a new project or idea.
3. For each repository, what is its role and which branch should be the default
   active branch? If `development` exists, recommend it; let the user choose a
   different branch.

Record the confirmed identity in `workspace.yaml`. A repository entry has a
path, mode, role, agent, and `default_branch`. Do not turn the current working
branch into a delivery policy.

When no repository is registered, offer the user a short Idea Brief
conversation, selected source intake, both, or deferment. Do not treat the
no-repository path as an error or as a reason to invent a repository.

Initialization stops after core identity. It does not ask about delivery,
commits, pushes, merges, publication, deployment, or external activity tools;
those are optional later configuration.

## 3. Start with a PRD or source

When the user provides a PRD, source document, or repository evidence:

1. Read the source.
2. Separate stated requirements, current behavior, assumptions, unknowns, and
   evidence.
3. Draft concise Product Knowledge with source references.
4. Present the draft for human confirmation.
5. Do not silently replace accepted context.

If no useful source or repository exists, ask focused discovery questions before
creating Product Knowledge or a plan.

## 4. Create and approve work

When context is sufficient:

1. Propose one or more plans by coherent domain or repository boundary.
2. Include dependencies, acceptance criteria, implementation scope, test scope,
   and verification commands.
3. Present plans as drafts.
4. Optionally review with `cc-review-plan`. Review stays read-only.
5. Wait for explicit human approval through `cc-approve-plan`.
6. Do not infer approval from conversation tone, tests, or agent output.
   Approval does not start execution.

## 5. Execute with sessions

After approval, `cc-run-plan` directs the root to:

- claim a plan execution and create or reuse its exclusive worktree;
- spawn a writer child through the host child-session primitive;
- spawn a later independent verifier child (`write_worktree: false`);
- record completion.yaml and ask for `cc-finish-plan`.

Sequential tasks share one writer child. Independent plans get separate
children and worktrees. If the host cannot spawn a child, report the missing
host primitive to the human.

Every child session receives a bounded delegation packet and returns a structured
handoff to its parent.

## 6. Resume, finish, and optional cleanup

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
