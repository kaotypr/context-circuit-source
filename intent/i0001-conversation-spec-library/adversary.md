# Spec adversary — i0001-conversation-spec-library (merged intent)

- contract_digest: (unset — draft, re-gated; the runtime was not run)
- criteria_sound: needs-work. Five passes ran. Passes 1–4 folded every finding; the pass-5
  confirming re-run (Sonnet 5, criteria read cold) found five fresh reachable loopholes —
  two high. The human elected to **accept all five as known residuals — none folded**. The
  criteria therefore still carry these five satisfy-and-still-wrong loopholes by conscious
  decision. They must surface at Gate 1 as **accepted-risk open questions**, never presented
  as certified-clean. Reaching a mechanically-clean `criteria_sound: yes` would require
  folding findings 1–5 (below) and one further confirming pass; the human declined that cost.
- human-gate residual: the semantic ADEQUACY/DEPTH of each plot's post-conditions (versus
  the mechanical non-vacuity that AC-FALSIFIABLE-GRADES / AC-ID-AGNOSTIC enforce) cannot be
  mechanically certified. Surface it at Gate 1 as a human-judgment item, not as certified.

This intent merged the former i0001 (authoring) and i0002 (go-live) into one decision.
Four independent adversary passes ran; every finding changed or added a criterion (high
change-rate, well-calibrated). The signal across passes is convergence: 12 → 12 → 6 → 4
findings, severities falling, the pass-4 findings being relocated-loophole refinements
rather than new structural holes.

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

## Pass 3 — confirming pass on the 15 merged criteria (6 findings, all folded)

Fresh independent child, criteria read cold (no implementation, no prior record). It
confirmed the mechanics were non-gameable — AC-OVERLAY-BOUNDED, AC-WIRED-SOURCE,
AC-FALSIFIABLE-GRADES, AC-REGEN-EQUIV attacked and unbroken — but found the criteria
certified internal self-consistency rather than the goal's two loudest claims
(whole-surface completeness, cannot-drift):

- **high — AC-SURFACE-COVERAGE certified inventory↔plot consistency, not surface
  completeness.** A self-drawn narrow inventory passed while the real surface was larger:
  "certifying a subset as the whole." → AC-SURFACE-COVERAGE now ANCHORS the inventory to an
  external product enumeration (`.agents/skills/*` verbs + WORKFLOW.md phases + AC-*), and a
  narrow inventory FAILS.
- **medium — cited INV/AC ids were not required to resolve.** A plot could cite a bogus or
  renamed id and pass — a silent-drift vector. → AC-DESCRIPTIVE-ONLY now requires every
  cited id to resolve to a live definition (method static → command).
- **medium — inventory↔plot drift was one-time-manual, not gated.** → AC-DRIFT-BIDIRECTIONAL
  now gates all three layers (inventory↔plot↔case); an added plot with no inventory entry,
  or vice versa, FAILS.
- **medium — "load-bearing field" self-defined + Dialogue never tied to Spec.** →
  AC-GEN-DERIVES now runs mutation over the FULL non-decorative Spec field set
  (propagate-or-decorative, no silently-ignored field); new **AC-DIALOGUE-FAITHFUL** fails a
  Dialogue that contradicts its Spec.
- **medium — "expects failure" hatch unbounded + live execution only for net-new.** →
  AC-NEW-CASES-RUN sets a live-PASS floor for successful-lifecycle flows and requires an
  "expects failure" plot to assert its SPECIFIC failure mode; AC-REGEN-EQUIV now names the
  live-executed regenerated subset and asserts each reaches its expected live grade.
- **low — AC-CORRECTIONS-MECHANICAL checked against an unstated oracle.** → now compares
  against a fixed retired-tag→corrected-id mapping recorded in the coverage docs.

Result: 15 criteria → 16 (added AC-DIALOGUE-FAITHFUL). A re-run on this set is pending
before the record can read `criteria_sound: yes`.

## Pass 4 — re-run on the 16 folded criteria (4 findings, all folded)

Fresh independent child, criteria read cold. Confirmed gaps #2 (citation resolution) and #3
(inventory↔plot drift) cleanly closed. Found that three folds RELOCATED a loophole rather
than removing it — each still delegated an author-controlled choice — plus one residual on
the completeness anchor's breadth:

- **medium — AC-REGEN-EQUIV made live execution of existing cases optional/author-named with
  no floor.** A successful-lifecycle existing case could silently regress live (gap #5
  reappearing on the existing half). → the live-executed set is now non-minimizable: it must
  include every successful-completion existing case and may not shrink below the pre-change
  live set; a reclassification dodge FAILS.
- **medium — AC-GEN-DERIVES replaced author-chosen "load-bearing" with author-chosen
  "decorative" — same freedom inverted.** An ignored Spec field could hide behind a
  decorative label (gap #4 relocated). → a reverse mutation test now PROVES a decorative
  field inert (mutating it must be byte-identical + no grade change); otherwise FAILS.
- **medium — AC-CORRECTIONS-MECHANICAL's oracle was a mapping authored inside this same
  change, not the externally-owned surface** (gap #6 relocated). → corrected ids must now
  resolve (read-only) against the external acceptance surface and carry provenance citing the
  owner; a self-consistent-but-wrong id FAILS.
- **medium — AC-SURFACE-COVERAGE's external anchor was real but its breadth unverified**
  (residual of gap #1): cross-cutting outcomes (refusal, host-blocked, clarification) fit no
  skill/phase/AC and would escape the net. → the enumeration now explicitly spans those
  cross-cutting conversational outcomes.

Bottom line (child's own): none theatre — each was a reachable satisfy-and-still-wrong. It
also named the one thing the criteria cannot mechanically own — the semantic depth of the
post-conditions — as a human-gate responsibility (recorded above).

Result: 16 criteria (no new id; four strengthened). A re-run on this set is pending.

## Pass 5 — confirming re-run on the 16 folded criteria (Sonnet 5, 5 findings, NONE folded — human accepted all as residuals)

Fresh independent child on Sonnet 5, criteria read cold (only the role file and the
contract; no implementation, harness, sources, or prior adversary record). It confirmed
the pass-1..4 folds hold — AC-GEN-DERIVES's reverse mutation test, AC-DRIFT-BIDIRECTIONAL's
content-hash-not-mtime, and AC-NEW-CASES-RUN's bounded "expects failure" hatch were all
attacked and unbroken — but found five reachable satisfy-and-still-wrong loopholes the
current wording still permits. Verdict: `criteria_sound: needs-work`.

The human reviewed these and **elected to accept all five as known residuals — none was
folded and the contract was not edited.** They are recorded here verbatim so Gate 1 can
present them as accepted-risk items.

- **high (F1) — AC-SURFACE-COVERAGE's "justified gap" hatch has no quality bar or cap.**
  An implementer can enumerate the full external surface, write one real plot, and label
  every other entry point a "gap: deferred" with one-line free text; the mechanical check
  (a gap record exists and points at the external list) passes while whole-surface coverage
  is defeated. Nothing bounds gap count, restricts the gap reason to a closed enum, or
  freezes the gap list the way the net-new plot set is frozen. Suggested fix (NOT applied):
  require each gap to cite a closed reason enum (free text FAILS) and freeze the full gap
  list at plan-of-record time.
- **high (F2) — AC-DESCRIPTIVE-ONLY's "defines no rule of its own" is not machine-checkable
  as written.** The command check only validates that cited INV-*/AC-* ids resolve. A plot
  can cite one valid id anywhere to satisfy that, then elsewhere write its own free-standing
  policy ("retries 3 times", "times out after 90s") sourced from no citation — passing the
  literal check while leaking a second, uncontrolled rule source (defeats non-goal #1,
  one-rule-one-owner). Suggested fix (NOT applied): extend the AC-REPORTING-RULES lint to
  flag normative/imperative language (must/shall/always/never/numeric thresholds) not
  structurally bound to an adjacent resolved citation.
- **medium (F3) — AC-REGEN-EQUIV's pre-change live-execution floor is self-attested.** The
  contract does not say where the authoritative record of "which cases were live-executed
  pre-change" comes from; the same actor authoring the regeneration can author the "before"
  baseline and under-declare it, shrinking the real floor while the AC still reads satisfied.
  Suggested fix (NOT applied): capture the pre-change live-executed set as a hash/timestamp-
  frozen artifact generated from the current suite before the plan of record, and diff
  against that.
- **medium (F4) — the constraint "avoid the noun 'contract' in the library" has no enforcing
  criterion.** AC-REPORTING-RULES lints file-name/path/id/digest/branch/model patterns and
  overstated assurance, but not the word "contract"; every AC can pass while the word is used
  throughout the library. The constraint exists in prose only. Suggested fix (NOT applied):
  fold a literal-word denial into the AC-REPORTING-RULES lint (outside a literal
  `contract.yaml` path reference).
- **low (F5) — AC-REGEN-EQUIV's "outside itemized intended edits" exception on seed/budget
  preservation is unbounded.** Any seed/budget drift can be waved through as an "intended"
  edit with no requirement that it fix a demonstrated calibration defect (F1's pattern at
  smaller scale). Suggested fix (NOT applied): require each itemized edit to cite the
  specific prior-case defect it corrects and cap edits to those named in the plan of record.

Bottom line: two high-severity findings (F1 completeness-hatch, F2 one-rule-one-owner leak)
are genuine integrity gaps; F3/F5 are freeze-discipline consistency; F4 is a stated
constraint with no teeth. All accepted unfolded by human decision. The criteria are NOT
mechanically clean — Gate 1 must show these as accepted residuals, not certified.
