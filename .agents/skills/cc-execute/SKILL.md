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
   acquires the one-worker lock, snapshots the plan, creates one branch and
   worktree per affected repository, and discovers each repository's own agent
   guidance from the prepared worktree (recorded as a grounding manifest).
2. Assemble the worker brief with `worker-brief-assemble` (per affected
   repository), adding only a one-line task focus. The runtime fills the brief
   from the grounding manifest and the plan and refuses a brief missing its
   repository-grounding section (preflight). Deliver the assembled brief
   verbatim — do not author or omit the repository-grounding facts, and do not
   read the runtime implementation to compose them (INV-GROUND-01/03).
3. Launch exactly one worker with that brief and the assigned worktrees
   (see `agents/worker.md`). The worker reads and honors the repository's own
   agent guidance, executes all tasks in dependency order, and commits each
   affected repository. Record each commit with `worker-commit-record` and the
   handoff with `worker-handoff-record` (including any `repository_friction`).
4. Launch one independent, read-only verifier (see `agents/verifier.md`) after
   `verifier-prepare`. It inspects the latest commit of every affected
   repository. Record its outcome with `verifier-result-record`.

Do not require confirmation for individual tasks, branches, worktrees, commits,
verifier steps, or repairs. The approved plan is the scope.

## Runtime actions — invoke, never read the engine

You never need to open the runtime implementation; invoke each action as
`sh wrapper/runtime/engine.sh <action> <args>` from the workspace directory.
Reading `wrapper/runtime/engine.sh` itself is out of scope for the coordinator.

- `execution-begin . <plan-id> <owner>` — preflight, snapshot, worktree(s), and
  repository-grounding discovery; prints `execution_id`.
- `worker-brief-assemble . <execution-dir> <repo> "<task focus>"` — assemble the
  grounded worker brief (fills the grounding directive + environment from the
  manifest); it preflights the required grounding slot.
- `attempt-begin <execution-dir>` · `worker-commit-record <execution-dir> <repo> implementation|repair`
  · `worker-handoff-record <execution-dir> <handoff-file>`.
- `verifier-prepare <execution-dir>` · `verifier-result-record <execution-dir> <attempt> passed|failed|blocked`.
- `repair-allowed <execution-dir>`.

The execution directory is `.runtime/executions/<plan-id>/<execution-id>/`; the
assigned worktree(s) and the grounding manifest are named in that record.

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

## Report — plain language only

Report by effect, in plain project language. Never expose internal mechanism to a
lay user: no runtime file names, no execution branches (`cc/...`), no worktrees, no
base-commit SHAs, and no preflight/verifier internals — even in a summary table.
Say "I built it and it was independently checked, and the check passed; nothing is
marked complete yet — that's your call" — not the branches, worktrees, or commits
behind it. Refer to a plan by its title and the branch the user works from by its
plain name. Reveal runtime records or branch mechanics only if the user explicitly
asks for diagnostics (`docs/terminology.md` is the internal→user-facing mapping).

## Boundaries

Never auto-approve, auto-complete, merge, push, publish, deploy, or discard.
Verification produces `verified` evidence; it never marks the plan `done`.
