---
name: cc-review-plan
description: Review a plan for evidence, scope, dependencies, risks, contradictions, and human decisions before execution.
---

# Review a plan

Use this skill when the user asks whether a plan is ready, what needs review,
or what decisions remain before execution.

## Inspect

Read the selected plan, its tasks, referenced Product Knowledge, source
provenance, repository-local instructions, relevant implementation evidence,
and active runtime ownership. Check that:

- `plan.yaml` is the canonical lifecycle record and has a known status;
- scope, non-goals, repository boundaries, dependencies, task paths, and stop
  conditions are explicit;
- acceptance criteria map to test and verification evidence;
- leases, worktrees, sessions, and handoffs do not show an ownership conflict;
- contradictions, stale context, missing evidence, and delivery actions remain
  visible.

## Report

Return the structured outcome from `docs/plan-review.md`. Distinguish routine
checks from human decisions. A review may recommend approval, revision, or a
blocker, but never changes plan status or claims completion.

If the plan is approved and dependency-ready, recommend `cc-run-plan`. Do not
create a separate task-execution ceremony or direct the user to `cc-run-task`.
