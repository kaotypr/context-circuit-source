---
schema_version: 2
id: AGF-002
plan: approval-gate-ux-performance
status: ready
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - agents/coordinator.md
  - wrapper/adapters/WORKFLOW.md
  - test/lifecycle/test-lifecycle.sh
  - test/approval/test-approval-flow.sh
depends_on: [AGF-001]
acceptance: [AGF-AC-02, AGF-AC-03, AGF-AC-05]
verification: [AGF-VT-02, AGF-VT-03, AGF-VT-04]
expected_evidence:
  - Confirmed approval calls the canonical transition primitive once and projects all included task statuses.
  - Before/after fixture diff proving only intended status lines changed and file bodies/endings are preserved.
  - Coordinator handoff describing approval as status-only with no lease, worktree, or execution side effects.
stop_conditions:
  - A confirmed request edits task bodies, normalizes unrelated formatting, or retries through manual repair patches.
  - The transition is partial, non-idempotent, or can run without current confirmation.
  - The implementation starts any writer, verifier, lease, worktree, Git, delivery, or external activity.
---

# Make confirmed approval a single bounded status transition

## Objective

Reduce confirmed-turn work to one canonical, status-only transition that is
safe, formatting-preserving, and easy to verify.

## Work

Integrate the confirmed approval path with the existing plan lifecycle helper
instead of hand-editing the plan and each task independently. Preserve the
current lifecycle checks, confirmation requirement, task projection mapping,
and maintainer dirty-delta behavior. Add only procedural coordinator guidance
needed to invoke the owner primitive and report bounded evidence.

## Non-goals

Do not redesign the lifecycle, make approval atomic with execution, suppress
required verification, or promise a fixed remote-host wall-clock duration.

## Verification

Use AGF-VT-02, AGF-VT-03, and AGF-VT-04.

## Expected evidence

Lifecycle test output, byte/diff comparison of task fixtures, one transition
record, and confirmation that no runtime, lease, worktree, or Git side effect
was introduced.

## Stop conditions

Stop on formatting drift, partial projection, repeated verification/repair
loops, missing confirmation, or scope expansion into execution.
