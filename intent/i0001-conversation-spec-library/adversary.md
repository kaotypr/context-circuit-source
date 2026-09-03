# Spec adversary — i0001-conversation-spec-library (merged intent)

- contract_digest: (unset — draft, re-gated; the runtime was not run)
- criteria_sound: needs-work (both passes) → all findings folded; a final confirming
  pass on the folded criteria is the remaining step before Gate 1.

This intent merged the former i0001 (authoring) and i0002 (go-live) into one decision.
Two independent adversary passes ran; every finding changed or added a criterion (high
change-rate, well-calibrated).

## Pass 1 — on the pre-merge go-live criteria (12 findings, all folded)

Circular/single-instance fidelity → AC-GEN-DERIVES (no committed-case input, all-plots
fidelity, mutation-sensitivity, determinism). "Reaches a grade" measured liveness not
assertion power → AC-FALSIFIABLE-GRADES. "id-agnostic" undefined → AC-ID-AGNOSTIC bound
to kind/count/terminal-state. Drift guard excluded the seed/budget region and skipped the
plot→case direction → AC-DRIFT-BIDIRECTIONAL. Plus: overlay unbounded, "suite green"
weakenable, seeds/budgets not diff-checked, fixtures "accepted" not exercised, corrections
subjective, net-new set unfrozen, determinism single-run.

## Pass 2 — on the merged criteria (12 findings, all folded)

- **high — nothing asserted the generated cases ARE the live suite source.** An
  implementer could pass every AC with a parallel generated tree beside 22 untouched
  hand-authored live cases. → **AC-WIRED-SOURCE** (bijection; suite runs the generated
  files; orphan case or orphan plot fails).
- **high — generator could table-embed case bodies** (mutation checked for only 3 named
  plots). → AC-GEN-DERIVES now requires mutation-sensitivity for EVERY plot's load-bearing
  fields.
- **high — "outside the overlay" was an unbounded determinism escape hatch.** →
  **AC-OVERLAY-BOUNDED** (overlay is exactly the enumerated set; widening fails).
- **high — "live and growing" not guaranteed green or repeatable.** → AC-NEW-CASES-RUN
  now requires grade PASS (except spec-expects-failure) reproducibly via the suite, not a
  stored artifact.
- **medium** — AC-WHOLE-SURFACE checked arithmetic not coverage → **AC-SURFACE-COVERAGE**
  (interaction inventory). Drift by mtime unreliable → content-hash + orphan detection.
  Case-level not condition-level falsifiability → per-post-condition. "Equal-or-stronger"
  subjective → machine-checked. No scope fence → **AC-SCOPE-FENCED** (diff only
  agent-harness/, zero wrapper/template). AC-36 suite unprotected → AC-REGEN-EQUIV now
  asserts it stays unchanged and green.
- **low** — AC-ID-AGNOSTIC extended to every entity-bound flow; AC-REPORTING-RULES became
  a static lint that gates in the suite (not manual dead weight).

Result: 11 criteria → 15, each falsifiable and non-gameable. The load-bearing addition is
AC-WIRED-SOURCE — it is what makes "the plots are the living source" true rather than
aspirational.
