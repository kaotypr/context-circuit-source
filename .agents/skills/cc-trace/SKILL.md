---
name: cc-trace
description: On intent approval, spawn one planner child per repository in scope to read the real code and write that repository's plan or plans, then run the feasibility check so the coordinator can publish without rewriting.
---

This skill is the approval-spawn packet. The child it launches is the
**planner** (`.context-circuit/agents/planner.md`).

It runs **after Gate 1**. At **Explore** no planner spawns. At **Standard** and
**Critical** it is mandatory. `contract.yaml` approved cannot skip or stand in
for the planner.

## Spawn immediately

After `intent-approve`, spawn before loading plan templates or reading the
target repository. The coordinator must **not** glob, grep, or walk that
repository. The planner is the code read.

Spawn each planner at the concrete `(model, effort)` configured for the
`planner` role in the host-local role-tiering config (`.context-circuit/docs/role-tiering.md`),
with adapter defaults when that host has no planner entry. Read
`role-tiering.local.yaml` from the **workspace root** when present — the same
directory as `repositories.local.yaml`, never a repository working copy. A
missing file in an isolated working copy is not an absent config.

For each repository in the intent's scope, spawn **one planner child per repository**,
in parallel when each child writes only its own repository's
plan or plans. If two children would allocate plan ids at the same time, spawn them
sequentially instead. Deliver the assembled
`.context-circuit/agents/planner-brief.md` with intent id, repository
id, tier, and workspace root filled in. Do not paste the contract. Do not add extra classification homework, naming notes, or a request to reconcile old product names.

If current plans already exist for this intent and repository and the finding
records the same revision, do not spawn a second planner.

## Feasibility check

Once each child has returned, run the feasibility check on `intent/<id>/finding.yaml`
before treating any plan as published and before any plan is written by the coordinator. This is coordinator judgment, not an
engine verb:

- **Feasible** with no intent-level question → update `INTENT.md` with
  `intent-human-status . <id> "approved, look complete, feasible"` and hand
  `cc-plan` the already-written plan files **in the same turn**. Do not rewrite
  them. Do not read the repository.
- **Intent revision required** → write no further plan. Relay the finding,
  update the intent, take Gate 1 again if the approved decision changes.
- **Not feasible** → **stop.** Explain the blocker. The planner wrote no plan.
- Classify every question before feasibility: intent-revision, plan-resolution,
  or already-answered.

An approved intent with an unresolved intent-level question is not ready for
publication. A required change that must *modify* a repository or area beyond a
bound scope is a question only when the plain request does not already authorize
that area.

## Plan boundary

The planner earns runnable done-checks against the current code and carries them into the plan. The coordinator never invents those commands.

The planner's files already name one repository. This child may write more than
one plan for that repository. Keep **one plan with
embedded tasks** when they share one bounded execution and verification
boundary. Use **stacked** plans, each one repository, with `plan_dependencies`
when partitions are independent, ordered, or have distinct failure surfaces.
A two-repository intent still needs one child per repository. The coordinator may add a `plan_dependencies` edge across
children; it must not collapse two repositories into one plan and must not
rewrite task bodies.

When `intent/<id>/detail/` exists, it is confirmed topic shape for the planner,
not a tracing output and not a second gate.

## Reuse

Do not rebuild a structural map cache as a substitute for the planner. Reused
grounding is the existing plan plus `finding.yaml` at the recorded revision.
Uncertain drift spawns the planner again.

## Host blocked

If the host cannot create the planner child, report `host-blocked` and write
no plan. Never author the plan in the coordinator session as a stand-in, and
never claim a look at the code that did not happen.

## Report

The human never hears "the planner" or "the finding." Relay what they must
decide, then the breakdown the child already wrote.
