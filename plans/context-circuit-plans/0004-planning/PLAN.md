# 0004 — Grounded planning

- **Plan ID:** `0004-planning`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0002-workspace-repository-model`, `0003-context-lifecycle`
- **Owns (invariants):** INV-PLAN-01/02/03/04

## Original request

Retroactive plan for the planning increment, built as if from an empty repo.
Source design: `sources/system-design/context-circuit/v0.5/core/05-planning-and-execution.md`.

## Objective and desired behavior

Expand a request into a readable, grounded plan and a canonical machine half.

- `plan.yaml` owns the stable id (`NNNN-<kebab-slug>`) and the human status
  (`draft → approved → done`); `PLAN.md` is the readable companion; task status
  is a synchronized projection, never a second lifecycle authority.
- Every task names the repositories it may change and bounded paths (or explicit
  repo-wide scope), with explicit dependencies and distinct acceptance /
  verification ids.
- Plans are grounded against Product Knowledge, source and repository evidence;
  anything missing or contradictory becomes an open question, assumption, or
  risk — never a silently invented decision.
- Creating or reviewing a plan never approves or executes it.

## Constraints and non-goals

- Non-goal: approval and execution (0005). This plan only drafts and reviews.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `conventions`
  (`context/CONVENTIONS.md`), `decisions` (`context/DECISIONS.md`).

## Tasks

1. **PLN-001** — plan + task schemas and readable templates.
2. **PLN-002** — `plan-allocate-id`, `plan-validate`, active `plans/INDEX.md`.
3. **PLN-003** — `cc-plan` create + review skill.

## Acceptance & verification

- Canonical id/status; distinct acceptance/verification ids; unique never-reused
  ids; undeclared repo / unknown dep rejected; review is status-preserving.
- `sh test/plans/test-plans.sh`.

## Assumptions, open questions, risks

- Risk: an agent inventing a decision to fill a gap — prevented by INV-PLAN-04
  and the grounding drift / request-fidelity checks in `cc-plan`.

## Expected commits and delivery notes

One commit for the source repository. Delivery is separate (0009).
