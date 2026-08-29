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
   (see `agents/worker.md`), at the worker's configured `(model, effort)` (see
   "Model & effort per role" below). The worker reads and honors the repository's
   own agent guidance, executes all tasks in dependency order, and commits each
   affected repository. Record each commit with `worker-commit-record` and the
   handoff with `worker-handoff-record` (including any `repository_friction`).
4. Launch one independent, read-only verifier (see `agents/verifier.md`) after
   `verifier-prepare`, at the verifier's configured `(model, effort)`. It inspects
   the latest commit of every affected repository. Record its outcome with
   `verifier-result-record`.
5. Record the observed inference wall-clock and the `(model, effort)` each role
   ran at as host evidence with `attempt-evidence-record` (below). This is
   evidence only; it never changes a verdict or the failure counter.

Do not require confirmation for individual tasks, branches, worktrees, commits,
verifier steps, or repairs. The approved plan is the scope.

## Runtime actions — invoke, never read the engine

You **must not read** `wrapper/runtime/engine.sh` or any runtime implementation
file (`wrapper/adapters/AGENTS.md` → Runtime owns this boundary): the actions below
and the execution brief carry everything needed to drive it. Invoke each action as
`sh wrapper/runtime/engine.sh <action> <args>` from the workspace directory.

- `execution-begin . <plan-id> <owner>` — preflight, snapshot, worktree(s), and
  repository-grounding discovery; prints `execution_id`.
- `worker-brief-assemble . <execution-dir> <repo> "<task focus>"` — assemble the
  grounded worker brief (fills the grounding directive + environment from the
  manifest); it preflights the required grounding slot.
- `attempt-begin <execution-dir>` · `worker-commit-record <execution-dir> <repo> implementation|repair`
  · `worker-handoff-record <execution-dir> <handoff-file>`.
- `verifier-prepare <execution-dir>` · `verifier-result-record <execution-dir> <attempt> passed|failed|blocked`.
- `attempt-evidence-record <execution-dir> <attempt> <key=value> ...` — record
  bounded per-attempt host evidence (`worker_wall_s`, `verifier_wall_s`,
  `worker_model`, `worker_effort`, `verifier_model`, `verifier_effort`,
  `complexity`, `escalated`). The runtime stores it and never interprets it.
- `repair-allowed <execution-dir>`.

The execution directory is `.runtime/executions/<plan-id>/<execution-id>/`; the
assigned worktree(s) and the grounding manifest are named in that record.

## Repair

On a verifier failure, pass the failure evidence back to the same worker within
the same execution. Check `repair-allowed`, begin a new attempt, let the worker
create a new commit for every repository it changes, and verify again. The
worker-failure counter increments on each rejection (including the first); at
three failures execution stops and all evidence is preserved.

When the worker's role has `escalate_on_repair: true`, launch the repair attempt
at a `(model, effort)` **raised above** the configured start (see below);
escalation changes only which model runs the attempt, never what a rejection
costs — the failure counter and the three-failure limit are untouched. Record the
raised `(model, effort)` and `escalated=true` with `attempt-evidence-record`.

## Model & effort per role

Spawn the worker and verifier at the concrete `(model, effort)` configured for
each role in the host-local role-tiering config, with adapter-shipped defaults for
any unset role (`docs/role-tiering.md` owns the shape, defaults, and
escalation ladder). This is a coordinator/host decision — the runtime is
model-blind (INV-RUNTIME-01) and `(model, effort)` authorizes nothing
(INV-HOST-01). It changes cost and speed, never meaning.

- Attempt 1 runs each role at its configured start. If the plan carries
  `complexity: high`, nudge the worker's attempt-1 start one step above the
  configured `(model, effort)`; a hard-pinned role ignores the hint.
- On repair, raise a role above its start only when its `escalate_on_repair` is
  true; a hard pin (`false`) holds the same setting at every attempt, even the
  third, and you report that pin's cost honestly rather than silently escalating.
- Verifier independence is role + read-only, never model class (INV-VERIFY-01/02):
  a smaller-model verifier — even the same model as the worker — is still a
  separate independent agent over committed state. Never collapse the two roles
  or let the verifier reuse the worker's context for cost.
- Record the `(model, effort)` used per attempt with `attempt-evidence-record`;
  never surface it to a lay user except under explicit diagnostics.

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
