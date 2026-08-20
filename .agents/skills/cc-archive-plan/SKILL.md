---
name: cc-archive-plan
description: Archive or restore a plan through an explicit human gate without changing lifecycle status or deleting evidence.
---

# Archive or restore a plan

Use this skill only when the human explicitly asks to archive or restore one
named plan. It satisfies a human `archive` gate. Archive eligibility is separate from the canonical plan lifecycle:
the skill appends evidence to `archive.yaml`; it does not edit `plan.yaml`,
task status, runtime records, Git history, or the plan bundle path.

## Read and validate

Read the named plan's `plan.yaml`, optional `archive.yaml`, included tasks,
declared dependencies and dependents, active sessions, leases, stack records,
worktree records, Git state, and relevant compatibility instructions. A missing
`archive.yaml` is the active-compatible state. When present, require:

- `schema_version: 1`;
- `plan` exactly matching the canonical plan ID;
- a non-empty, append-only `events` sequence;
- every event to include `action` (`archived` or `restored`), `at`, `actor`,
  `reason`, `observed_status` (`draft`, `approved`, or `done`), and
  `replacements` (possibly empty).

The latest event is the current state. Reject a malformed, empty,
identity-mismatched, or unknown-event record instead of guessing. Historical
inspection may always read the bundle at its unchanged original path.

## Preflight

Before appending either action, block on any of the following:

- a live lease, live writing session, active stack membership, or ambiguous
  session/worktree ownership for the plan;
- dirty or unpushed work in the assigned plan worktree, or an uncertain
  worktree/branch relationship;
- an archive request for an already archived plan, or a restore request for an
  already active plan;
- a non-archived dependent that would be left with an unresolved archive
  target. An archived `done` target remains inspectable and satisfies a
  dependency as done; an archived `draft` or `approved` target blocks rather
  than disappearing;
- on restore, incompatible current repository branch, worktree, dependency, or
  compatibility state.

There is no cascade. A replacement reference records context but does not edit
another plan's dependencies or authorize its archive or restore. Archive and
restore are separate from approval, execution, finishing, cleanup, merge,
publication, and delivery.

## Confirm and append

Report the target, canonical status that will remain unchanged, current archive
state, preflight results, and the exact reason and replacement references that
will be recorded. Append an event only after explicit current-session human
confirmation of the requested `archive` action.

Use an atomic same-directory write. Preserve every existing event unchanged;
do not rewrite history. Append exactly one event with the current UTC time,
acting session or human identity, reason, canonical status observed at action
time, and replacement list. The action is `archived` or `restored`.

On archive, ordinary discovery must stop selecting the bundle. On restore,
ordinary discovery may select it again, but normal approval and execution
preflight still apply. A restore never re-approves a plan.

## Route guards

Normal session entry and `cc-whats-next` omit archived plans. `cc-approve-plan`,
`cc-run-plan`, `cc-run-stack`, and `cc-finish-plan` reject an archived plan,
including an explicitly named plan or stack member. Those routes report the
historical path and recommend explicit inspection or a separate restore gate;
they must not silently restore the plan.

## Safety

Do not delete, move, or clean a plan bundle, `.runtime/`, a worktree, a
branch, or Git history. Do not change lifecycle status or task projection, take
over ownership, merge, push, publish, deploy, or create a replacement plan.
Only the root coordinator after the human's current-session confirmation may
append an archive record. A child worker or verifier cannot archive or restore.
