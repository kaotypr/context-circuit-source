# 0013 — Runtime cleanup also removes Explore worktrees

Plan ID: `0013-cleanup-explore-worktrees`  
Intent: `i011-cleanup-explore-copies`  
Status: `draft`

## Original request and coverage

When you ask for runtime cleanup, Explore worktrees are included, not only plan-execution worktrees. This plan depends on `0012-explore-session-path`.

## Objective and desired behavior

An explicit runtime-cleanup request also removes worktrees under `.runtime/explore/`. Cleanup does not run by itself when an Explore session closes. A still-live Explore session is not silently deleted.

## Constraints and non-goals

- Depends on `0012-explore-session-path` for the new path.
- Do not auto-delete an Explore worktree when the session closes cleanly.
- Do not change delivery, promote, or the one-repository Explore rule.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `direct-collaboration` — close preserves the worktree until explicit cleanup.
- `plan-execution` — failed work is never silently cleaned up.

Grounding summary: derived from the approved i011 contract without a tracer, at explicit human request.

## Repositories and source evidence

`context-circuit-source`. Expected surface: pairing skill, runtime cleanup, coordinator, pairing tests.

## Tasks

1. **CLN-001** — teach cleanup about Explore worktrees. Acceptance: `CLN-AC-001`–`CLN-AC-003`. Verification: `CLN-VT-001`.
2. **CLN-002** — prove they are included. Acceptance: `CLN-AC-001`. Verification: `CLN-VT-002`.

## Acceptance criteria

- `CLN-AC-001` — Explicit cleanup removes Explore worktrees too.
- `CLN-AC-002` — Cleanup finds them under `.runtime/explore/`.
- `CLN-AC-003` — No silent delete on close or of a live session.

## Verification

- `CLN-VT-001` — `rg -n 'cleanup|runtime/explore|pair-close|worktree' .agents/skills/cc-pair/SKILL.md agents/coordinator.md wrapper/runtime/engine.sh context/domains/direct-collaboration/README.md`
- `CLN-VT-002` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: cleanup must not walk into a live Explore worktree.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place update expected on `direct-collaboration`. Review at completion.
