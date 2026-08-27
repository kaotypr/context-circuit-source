# 0006 — Independent verification and repair

- **Plan ID:** `0006-verification-and-repair`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0005-approval-and-execution`
- **Owns (invariants):** INV-VERIFY-01, INV-VERIFY-02, INV-REPAIR-01

## Original request

Retroactive plan for the verification + repair increment, built as if from an
empty repo. Source design:
`sources/system-design/context-circuit/v0.5/core/05-planning-and-execution.md`
(verification, repair) and `06-runtime-engine.md`.

## Objective and desired behavior

- One independent verifier checks the latest commit of every affected
  repository, strictly read-only with respect to product files. It never
  repairs, modifies product files, changes plan status, or self-verifies.
- If the host cannot create an independent read-only verifier, the execution is
  **blocked** — the worker/coordinator must not self-verify as a substitute.
- A failed verification routes back to the same worker as a **new** commit. The
  failure counter increments on each rejection (the initial implementation is
  attempt one) and stops at three, preserving all evidence.

## Constraints and non-goals

- Non-goal: human completion (0007). A passing verifier never marks a plan done.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `architecture`
  (`context/ARCHITECTURE.md`), `decisions` (`context/DECISIONS.md`).

## Tasks

1. **VER-001** — verifier-result schema, verifier role, read-only enforcement.
2. **VER-002** — repair loop + three-failure counter.
3. **VER-003** — `cc-verify` skill.

## Acceptance & verification

- Verifier cannot modify product files; missing verifier → blocked; repair is a
  new commit; stop at third failure with evidence preserved.
- `sh test/execution/test-execution.sh`.

## Assumptions, open questions, risks

- Risk: conflating a host block with a worker failure — separated explicitly by
  INV-REPAIR-01.

## Expected commits and delivery notes

Repair attempts add commits to the execution branches; nothing is delivered here.
