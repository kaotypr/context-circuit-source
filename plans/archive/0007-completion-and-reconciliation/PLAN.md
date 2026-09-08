# 0007 — Human completion and knowledge reconciliation

- **Plan ID:** `0007-completion-and-reconciliation`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0006-verification-and-repair`, `0003-context-lifecycle`
- **Owns (invariants):** INV-COMPLETE-01, INV-COMPLETE-02

## Original request

Retroactive plan for the completion increment, built as if from an empty repo.
Source design: `sources/system-design/context-circuit/v0.5/core/05-planning-and-execution.md`
(completion) and `04-context-lifecycle.md` (reconciliation).

## Objective and desired behavior

- Only an explicit human request changes plan status to `done`, and only when the
  latest execution passed independent verification. Verification alone never
  completes a plan.
- Marking done records an implementation completion record (plan revision,
  execution, per-repository commits, verifier result, human request) and starts
  Product Knowledge reconciliation.
- Reconciliation produces proposals (or an explicit no-update statement); a plan
  may be done while a proposal is still pending, and accepting a proposal is a
  separate human decision.

## Constraints and non-goals

- Non-goal: delivery (0009) and archive (0008) — independent later actions.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `decisions`
  (`context/DECISIONS.md`).

## Tasks

1. **CMP-001** — completion schema, eligibility, record.
2. **CMP-002** — knowledge reconciliation into proposals.
3. **CMP-003** — `cc-complete` skill.

## Acceptance & verification

- Done requires verified evidence; completion record has full provenance;
  reconciliation never silently accepts an update.
- `sh test/completion/test-completion.sh`.

## Assumptions, open questions, risks

- Risk: treating a green verifier as completion — prevented by INV-COMPLETE-01.

## Expected commits and delivery notes

Completion writes runtime records only; no product commit is required to mark done.
