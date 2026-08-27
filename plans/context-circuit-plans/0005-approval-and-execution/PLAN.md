# 0005 — Approval and multi-repository execution

- **Plan ID:** `0005-approval-and-execution`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0004-planning`
- **Owns (invariants):** INV-APPROVE-01, INV-EXEC-01/02/03/04, INV-PRESERVE-01 (cites INV-OWN-01, INV-COMMIT-01)

## Original request

Retroactive plan for the approval + execution increment, built as if from an
empty repo. Source design:
`sources/system-design/context-circuit/v0.5/core/05-planning-and-execution.md`
(execution) and `06-runtime-engine.md`.

## Objective and desired behavior

- Approval is an explicit conversational gate that only flips `draft → approved`
  after deterministic readiness checks — no confirmation card, no hidden token.
- Only an approved plan executes; execution is a separate authorization unless
  one request explicitly asks for both.
- One worker runs every task of one approved plan in dependency order across all
  mapped repositories, each in its own branch `cc/<plan-id>/<repo-id>` and
  worktree from the captured anchor tip; the anchor checkout is never written.
- The worker commits each changed repository before verification; repairs are
  new commits, never amendments.
- At most one active writer (atomic lock); failure or interruption preserves all
  evidence and never cleans up silently.

## Constraints and non-goals

- Non-goal: verification, repair loop counting (0006), completion (0007),
  delivery (0009). This plan produces committed, unverified work.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `architecture`
  (`context/ARCHITECTURE.md`), `decisions` (`context/DECISIONS.md`).

## Tasks

1. **EXE-001** — approval gate + execution/handoff schemas.
2. **EXE-002** — execution-begin, worktrees, worker role, commit-before-verify.
3. **EXE-003** — recovery, one-writer ownership, preservation.
4. **EXE-004** — `cc-execute` approve/execute/repair skill.

## Acceptance & verification

- No card at approval; unapproved cannot execute; isolated worktrees; commit
  before verify; lock never stolen; interruption preserves evidence.
- `sh test/execution/test-execution.sh`.

## Assumptions, open questions, risks

- Risk: a dirty anchor corrupting the base — mitigated by rejecting dirty anchors
  at execution-begin.

## Expected commits and delivery notes

One commit per affected repository per attempt. Delivery is separate (0009).
