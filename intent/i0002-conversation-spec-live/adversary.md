# Spec adversary — i0002-conversation-spec-live

- contract_digest: (unset — draft; the runtime was not run, so no digest is frozen)
- challenged: the FIRST draft of the acceptance criteria
- criteria_sound: needs-work
- status: findings folded into a criteria revision; a fresh challenge on the revised
  criteria is REQUIRED before Gate 1 (the adversary record must bind to the criteria
  about to be frozen, and criteria_sound must be `yes` to approve).

## Findings (first pass)

1. **high — AC-GEN-FIDELITY is circular and single-instance.** "generator output ==
   committed case file" is satisfiable by a generator that reads/embeds the committed
   case, and it was checked for one plot only. → Folded into **AC-GEN-DERIVES**:
   reads no committed case, fidelity for ALL plots, mutation-sensitivity on ≥3 named
   plots, plus determinism on re-run.
2. **high — "reaches a recorded grade" measures liveness, not assertion power.** A case
   with vacuous post-conditions runs and grades while asserting nothing. → New
   **AC-FALSIFIABLE-GRADES**: every post-condition must be flippable to FAIL by a
   negative fixture.
3. **high — "id-agnostic" undefined; admits count-only / content-free matching.** →
   **AC-ID-AGNOSTIC** now defines it as id-string-insensitive but bound to entity kind,
   count, and terminal state, with a wrong-entity/kind/count negative fixture.
4. **high — the drift guard's exclusion zone is its own bypass, and gating unspecified.**
   Excluding the overlay (budgets/host lane) left the calibrated seeds/budgets
   unguarded, and the guard might run out-of-band. → **AC-DRIFT-BIDIRECTIONAL** runs
   inside the suite and covers the full case incl. seed/budget fields.
5. **high — editing a PLOT was never tied to forced revalidation** (only the case→plot
   direction was guarded). → **AC-DRIFT-BIDIRECTIONAL** also fails when a plot is newer
   than / inconsistent with its generated case.
6. **medium — "suite stays green" satisfiable by lowering assertion power.** → **AC-
   REGEN-EQUIV**: case count unchanged (or removals itemized), and every retired inline
   expectation maps 1:1 to an equal-or-stronger generated post-condition.
7. **medium — "seeds/budgets preserved exactly" not actually checked** (green ≠
   byte-identical). → **AC-REGEN-EQUIV** adds a mechanical zero-diff on seed/budget fields.
8. **medium — the "environmental overlay" exclusion is unbounded.** → Constraint added:
   the overlay is a FIXED enumerated set (per-host budgets + host lane), everything else
   is plot-derived and drift-guarded.
9. **medium — AC-NEW-FIXTURES tested "accept", not "exercise".** → **AC-NEW-FIXTURES-
   EXERCISED**: each seed/fault/predicate exercised by ≥1 case with a negative fixture.
10. **medium — AC-CORRECTIONS rested on subjective manual review; could leave dangling
    references.** → **AC-CORRECTIONS-MECHANICAL**: grep for zero retired-tag references +
    assert the five cases map to corrected ids.
11. **low — net-new set not frozen; "preserved failure" an open escape hatch.** →
    Constraint freezes the net-new set at plan authoring; **AC-NEW-CASES-RUN** restricts
    "preserved failure" to plots whose Spec expects failure.
12. **low — generator determinism under-specified** (single run vs two agreeing runs).
    → Folded into **AC-GEN-DERIVES** (two consecutive generations byte-identical).

Change rate: 12/12 findings changed a criterion — the pass was well-calibrated, not
theatre.
