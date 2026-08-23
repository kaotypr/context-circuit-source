---
schema_version: 2
id: AGF-002
plan: approval-gate-ux-performance
status: done
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - agents/coordinator.md
  - wrapper/adapters/WORKFLOW.md
  - docs/gates.md
  - docs/planning.md
  - test/lifecycle/test-lifecycle.sh
  - test/approval/test-approval-flow.sh
  - test/ownership/test-ownership.sh
depends_on: [AGF-001]
acceptance: [AGF-AC-02, AGF-AC-03, AGF-AC-04, AGF-AC-06, AGF-AC-05]
verification: [AGF-VT-02, AGF-VT-03, AGF-VT-04, AGF-VT-05]
expected_evidence:
  - Confirmed approval calls the canonical transition primitive once and projects all included task statuses.
  - Before/after fixture diff proving only intended status lines changed and file bodies/endings are preserved.
  - On product-source, the same turn then presents the existing commit-approved-plan card without committing.
  - On instantiated or wrapped workspaces, the handoff names Run approved plan as the next explicit request.
  - Coordinator handoff describing approval as status-only with no lease, worktree, commit, or execution side effects.
stop_conditions:
  - A confirmed request edits task bodies, normalizes unrelated formatting, or retries through manual repair patches.
  - The transition is partial, non-idempotent, or can run without current confirmation.
  - Confirm approval commits Git, folds the commit into the approval confirmation phrase, or starts any writer, verifier, lease, worktree, delivery, or external activity.
  - Product-source still waits until Run approved plan to present the commit card when the dirty set is the exact approval projection.
---

# Make confirmed approval a single bounded status transition

## Objective

Reduce confirmed-turn work to one canonical, status-only transition that is
safe, formatting-preserving, and easy to verify. On `product-source`, follow
that transition immediately with the existing maintainer commit card so run
does not become the place the dirty-base stop is discovered.

## Work

Integrate the confirmed approval path with the existing plan lifecycle helper
instead of hand-editing the plan and each task independently. Preserve the
current lifecycle checks, confirmation requirement, task projection mapping,
and `INV-OWN-07` dirty-delta classifier. After a successful status transition
on `product-source`, if the dirty set is exactly the approval projection,
present the current `commit-approved-plan` card from `docs/gates.md`. Do not
run `commit-approved-plan` on this turn. Instantiated or wrapped workspaces
keep today's next action: `Run approved plan <id>`. Add only procedural
coordinator and planning/gate wording needed to invoke the owner primitive,
present that follow-on card, and report bounded evidence.

## Non-goals

Do not redesign the lifecycle, make approval atomic with commit or execution,
weaken `DIRTY_BASE_BLOCKED` for unrelated files, auto-commit, create a
worktree from dirty files, suppress required verification, or promise a
fixed remote-host wall-clock duration.

## Verification

Use AGF-VT-02, AGF-VT-03, AGF-VT-04, and AGF-VT-05.

## Expected evidence

Lifecycle and ownership test output, byte/diff comparison of task fixtures,
one transition record, a product-source follow-on commit card with no Git
commit on the confirm-approval turn, and confirmation that no runtime, lease,
worktree, or execution side effect was introduced.

## Stop conditions

Stop on formatting drift, partial projection, repeated verification/repair
loops, missing confirmation, a Git commit on confirm approval, or scope
expansion into execution.
