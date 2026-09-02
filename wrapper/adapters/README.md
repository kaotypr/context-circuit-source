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

> Work with me directly on <small change> in <repository>.
> I want to build <feature> — <what "correct" means>.
> Approve this intent.
> Execute it, then ship it.

## What you control

Approval, execution, completion, archive, restore, and delivery (pull request,
merge, push, deployment) are separate explicit actions you request.
Direct collaboration is human-supervised and does not produce verified status.
Creating or reviewing a plan never approves or executes it. Verification produces
evidence; only you decide when a plan is complete.

## Human surfaces

- `plans/<plan-id>/PLAN.md` — the readable plan.
- `plans/INDEX.md` — active plans.
- the worker handoff and the independent verifier result after an execution.

The workspace works offline and stores no credentials. Branches, worktrees, and
runtime records are managed for you.
