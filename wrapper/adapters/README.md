# Context Circuit

A universal project workspace for AI-assisted work. It holds agent-oriented
Product Knowledge and readable plans so an agent can understand your project,
create a grounded plan, and execute it safely across one or more Git
repositories.

## Get started

Talk to the workspace in ordinary language:

> What is this workspace?

To set up a project:

> Initialize this workspace for <project>, and connect the <repo> repository.

To do work:

> Create a plan for <feature>.
> Review plan <plan-id>.
> Approve plan <plan-id> and execute it.
> Mark <plan-id> complete.

## What you control

Approval, execution, completion, archive, restore, and delivery (pull request,
merge, push, publication, deployment) are separate explicit actions you request.
Creating or reviewing a plan never approves or executes it. Verification produces
evidence; only you decide when a plan is complete.

## Human surfaces

- `plans/<plan-id>/PLAN.md` — the readable plan.
- `plans/INDEX.md` — active plans.
- the worker handoff and the independent verifier result after an execution.

The workspace works offline and stores no credentials. Branches, worktrees, and
runtime records are managed for you.
