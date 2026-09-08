# 0012 — Human-named Explore worktrees under .runtime/explore

Plan ID: `0012-explore-session-path`  
Intent: `i010-explore-session-path`  
Status: `draft`

## Original request and coverage

When you ask to start Explore, you choose the name, and that Explore worktree lives under `.runtime/explore/<name>`, not mixed with plan-execution worktrees under `.runtime/worktrees/cc-pair/`.

## Objective and desired behavior

Starting Explore asks you for a short name and uses it. The worktree lives at `.runtime/explore/<human-name>`. The agent does not invent the folder name. Explore still uses one repository, you still judge the result live, and promote still works.

## Constraints and non-goals

- Do not change Explore's one-repository, no-verifier rules.
- Do not change plan-execution worktree layout except to stop mixing Explore into it.
- Cleanup of existing copies is the following plan.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `direct-collaboration` — Explore isolation.
- `plan-execution` — `.runtime/worktrees/` for plan execution.

Grounding summary: derived from the approved i010 contract without a tracer, at explicit human request.

## Repositories and source evidence

`context-circuit-source`. Expected surface: `cc-pair`, pairing runtime, pairing schema, direct-collaboration domain, pairing tests.

## Tasks

1. **EXP-001** — ask for the name and keep Explore worktrees under `.runtime/explore`. Acceptance: `EXP-AC-001`–`EXP-AC-003`. Verification: `EXP-VT-001`.
2. **EXP-002** — prove the split. Acceptance: `EXP-AC-002`. Verification: `EXP-VT-002`.

## Acceptance criteria

- `EXP-AC-001` — Human-chosen name; agent does not invent it.
- `EXP-AC-002` — Explore worktrees live under `.runtime/explore/<human-name>`.
- `EXP-AC-003` — Explore meaning, one-repository rule, and promote are unchanged.

## Verification

- `EXP-VT-001` — `rg -n 'runtime/explore|cc-pair|pair-begin|worktree' .agents/skills/cc-pair/SKILL.md wrapper/runtime/engine.sh context/domains/direct-collaboration/README.md`
- `EXP-VT-002` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: existing pairing-session pointers that record `.runtime/worktrees/cc-pair/` need a migration or ignore path.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place update expected on `direct-collaboration`. Review at completion.
