# Phase 6a coverage cross-check — closed knowledge loop & organization

How the five phase-6a plots cover the existing `scenarios/` cases and the harness gaps
they surface. Descriptive only: every rule cited is owned by
`.context-circuit/wrapper/contracts/invariants.yaml` or `test/acceptance/criteria-map.yaml`.

## Existing cases → plots

| existing case | plot | relationship |
| --- | --- | --- |
| `19-knowledge-debt-blocks-next-plan` | `knowledge-debt-blocks-next-plan` | **evolves** — same seed/budget; AC corrected (below) |
| *(none)* | `reconcile-and-proceed` | **new** — the happy complement (reconcile → clear → unblock) |
| `07-archive-plan` | `archive-plan` | **evolves** — same seed/budget |
| `07b-restore-archived-plan` | `restore-archived-plan` | **evolves** — same seed/budget |
| *(none)* | `archive-restore-intent` | **new (opt)** — intent-side status-blind archive (part of AC-30) |

## AC mapping drift (surfaced, not hidden)

- **`19` → `[AC-34]` (was `[AC-22, AC-35]`).** Case 19 is the closed-loop gate: a
  prior delivery's reconciliation debt blocks the next plan's grounding in the same
  knowledge scope until an explicit reconcile/defer clears it — **AC-34** exactly. The
  historical tags are misfits: **AC-22** is *marking done reconciles PK* (the
  completion moment, cases 06/18), and **AC-35** is *inferred completion / change set*
  (cases 18/21). This dialogue demonstrates neither — the foundation was completed at
  seed time; the dialogue is the block. **Recommendation:** correct case 19's AC list
  to `[AC-34]` at regeneration; AC-34's suite (`test/knowledge/test-debt.sh`) remains
  the deterministic owner.
- **`07` / `07b` mappings unchanged** — `[AC-24]` / `INV-ARCHIVE-01/02` are correct.

## The two new plots

- **`reconcile-and-proceed`** (AC-34, INV-KNOWLEDGE-02, INV-COMPLETE-02) — the positive
  branch of the same AC-34 clause case 19 exercises negatively: the human reconciles
  (accepts) the pending unit, the debt clears, and the previously-held follow-up is
  unblocked. Conversation-only; execution of the unblocked plan is owned by the phase-4
  plots. Distinct from phase-1's `accept-or-defer-context-proposal` (the standalone
  consent gate) — this one is specifically the debt marker clearing so the next plan
  proceeds.
- **`archive-restore-intent`** (opt; AC-30, INV-ARCHIVE-01/02) — completes the
  Organization section by showing an intent archives the same status-blind way a plan
  does (cc-intent ships `intent-archive`/`intent-restore`). Authored as the archive
  direction (detectable: starts active, ends archived); restore is the symmetric
  inverse, exactly as `restore-archived-plan` is to `archive-plan`.

## Harness capability gaps surfaced this phase

Flagged inline in the affected `generated/*.case.yaml`. None blocks authoring.

1. **New post-condition predicates.**
   - `knowledge_debt_resolved`, `knowledge_unit_accepted`, `grounding_unblocked`
     (`reconcile-and-proceed`) — assert the debt cleared by explicit decision and the
     follow-up unblocked. The deterministic debt-clear mechanic is owned by
     `test/knowledge/test-debt.sh`.
   - `intent_archived` (`archive-restore-intent`) — the intent analogue of
     `plan_archived`, with the approved state preserved.
2. **No new seed shapes.** `reconcile-and-proceed` reuses case 19's exact seed
   (`completed-standard` foundation + `draft` follow-up); the archive plots reuse cases
   07/07b's seeds; `archive-restore-intent` reuses the phase-2 `intents:` seed shape
   with `state: approved`.
