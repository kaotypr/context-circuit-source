---
name: cc-finish-plan
description: Mark an executed plan done after human confirmation without deleting runtime.
---

# Finish a plan

Use this skill when the user asks to mark a named executed plan done.
`cc-finish-plan` is the named status-change skill. It satisfies the human
`status-change` gate. `cc-run-plan` remains the sole execution entry. Finish
is not a second implementation workflow.

## Route reads

Use the `resume` manifest for the runtime evidence and the `review` manifest
for the plan and verifier evidence in
`docs/agent-workspace-workflow.md#route-read-manifests`. The local guard below
still requires `completion.yaml`, task evidence, independent verification, and
current-session human confirmation before changing canonical status.

## Preflight

Read the named plan's canonical `plan.yaml`, `.runtime/plans/<plan-id>/completion.yaml`,
the independent verifier handoff (prefer `handoff.yaml`; use historical
Markdown-only `handoff.md` only when YAML is absent), task evidence, the plan
lease, and worktree Git state. Confirm:

- the plan identifier is known and not contradictory;
- plan status is exactly `approved`;
- `completion.yaml` exists with `status: ready-for-human-status-change`;
- every task has durable evidence, independent verification passed, and no
  remaining blocker;
- remaining Git risk is reported and is not cleaned away by this skill.

Report what will change and what will not:

- canonical `plan.yaml` status will become `done`;
- included task status will become `done` in one idempotent bulk operation;
- runtime sessions, handoffs, and worktrees stay;
- Git history stays; this skill does not merge, push, or delete `.runtime/`.

## Confirm

Change `approved` → `done` only after explicit human confirmation in the
current session. Tests, Git state, a verifier handoff, or completion evidence
do not mark the plan done.

If the human declines, the plan stays `approved` and completion evidence stays
`ready-for-human-status-change`.

## On confirmed success

1. Set plan status to `done`.
2. Reconcile included task status to `done` in one idempotent bulk operation.
3. Update `completion.yaml` to record that the human gate was satisfied:

```yaml
status: completed
human_gate: status-change
canonical_status_changed: true
```

4. Release `lease.lock/` if this coordinator owns that lease. A live writing
   session owned by someone else is a blocker, not a silent takeover.
5. Leave `.runtime/` sessions, handoffs, and worktrees in place.

## Refuse

- Plan is `draft` or already `done`: refuse. Do not invent status.
- Unknown or contradictory identifier: refuse.
- `completion.yaml` is missing or its status is not
  `ready-for-human-status-change`: refuse.
- Missing task evidence, independent passing verification, or a remaining
  blocker: refuse.

## Safety

Must not merge, push, publish, deploy, delete worktrees, or delete `.runtime/`.
A child worker or verifier cannot finish a plan. Only the root coordinator,
after current-session human confirmation, may write canonical plan status or
release a lease it owns.
