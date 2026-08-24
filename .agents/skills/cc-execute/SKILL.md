---
name: cc-execute
description: Approve and execute an approved plan with one worker and one independent verifier, coordinating repair within the three-failure limit.
---

## Approval

Approval is an explicit conversational human gate, not a confirmation card. On
"approve plan X", run the runtime `plan-approve` (draft → approved after
readiness checks). If the user asks to execute a draft plan, refuse plainly:
the plan must be approved first. Support the compound request "approve plan X
and execute it" by approving, re-reading the approved status, then executing.

## Execute

Execute only an approved plan. Give a short summary (plan, objective,
repositories/branches, task count, worker and verifier roles, failure limit),
then:

1. Run `execution-begin`: it validates approval and repository bindings,
   validates and captures each `anchor_branch` tip, rejects dirty anchors,
   acquires the one-writer lock, snapshots the plan, and creates one branch and
   worktree per affected repository.
2. Launch exactly one worker with the execution brief and assigned worktrees
   (see `agents/writer.md`). The worker executes all tasks in dependency order
   and commits each affected repository. Record each commit with
   `worker-commit-record` and the handoff with `worker-handoff-record`.
3. Launch one independent, read-only verifier (see `agents/verifier.md`) after
   `verifier-prepare`. It inspects the latest commit of every affected
   repository. Record its outcome with `verifier-result-record`.

Do not require confirmation for individual tasks, branches, worktrees, commits,
verifier steps, or repairs. The approved plan is the scope.

## Repair

On a verifier failure, pass the failure evidence back to the same worker within
the same execution. Check `repair-allowed`, begin a new attempt, let the worker
create a new commit for every repository it changes, and verify again. The
worker-failure counter increments on each rejection (including the first); at
three failures execution stops and all evidence is preserved.

## Blocked verifier

If the host cannot create an independent read-only verifier, the execution is
blocked (`host-blocked`). Do not self-verify and do not downgrade the evidence
requirement. Preserve worktrees and commits.

## Boundaries

Never auto-approve, auto-complete, merge, push, publish, deploy, or discard.
Verification produces `verified` evidence; it never marks the plan `done`.
