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

Intent approval and delivery (pull request, merge, push, deployment) are the two
explicit human gates. Execution follows an approved intent (scope-safety is settled
at delivery, Gate 2, not by an automated scope gate);
Standard and Critical completion are an explicit mark-done. Explore is planless
and human-supervised.
Archive and restore remain explicit organization actions. Creating or reviewing a
plan never approves or executes it, and verification never implies completion.

## Human surfaces

- `plans/<plan-id>/PLAN.md` — the readable plan.
- `plans/INDEX.md` — active plans.
- the worker handoff and the independent verifier result after an execution.

The workspace works offline and stores no credentials. Branches, worktrees, and
runtime records are managed for you.
