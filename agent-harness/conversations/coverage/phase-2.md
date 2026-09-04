# Phase 2 coverage cross-check — the intent gate (Gate 1)

How the six phase-2 plots cover the existing `scenarios/` cases and the harness
gaps they surface. Descriptive only: every rule cited is owned by
`wrapper/contracts/invariants.yaml` or `test/acceptance/criteria-map.yaml`.

## Existing cases → plots

| existing case | plot | relationship |
| --- | --- | --- |
| `03-plan-review-and-approve` | `review-intent-approve-derive` | **evolves** — same seed + budgets; approve-but-hold (separated gate) |
| `04-refuse-unapproved-execution` | `refuse-before-gate1` | **evolves** — same seed + budget; fail-closed |
| `15-intent-adversary-revise-and-approve` | `feasibility-question-then-approve` | **evolves** — same seed + budgets |
| `16-scope-expansion-regate` | `scope-reach-feasibility-question` | **evolves** — same seed + budget; feasibility question |
| `20-tier-escalation-fails-upward` | `tier-fails-upward-refuse-explore` | **evolves** — same seed + budget; the one safety-critical automated check |
| *(none)* | `approve-and-build-one-turn` | **new** — the compound face of Gate 1 (AC-06) |

The compound and separated faces of Gate 1 are now a clean pair:
`review-intent-approve-derive` (approve, hold) and `approve-and-build-one-turn`
(approve, build) — both owned by INV-EXEC-01's two clauses.

## AC / invariant mapping drift (surfaced, not hidden)

Two of these matter more than phase 1's, because they sit where v1.0 reworked how
safety and feasibility are settled.

- **`20` → AC-33 (was `[AC-20, AC-29]`).** Case 20 is the one safety-critical
  automated check: a security surface classifies Critical and the coordinator
  refuses to drop the independent check (Explore) — exactly **AC-33** ("consequence
  tiering … refuses dropping the verifier (Explore) on a risk surface"). The
  historical tags are mismatched: **AC-20** is *host-blocked verifier* (a capability
  failure, owned by case 09), and **AC-29** is *direct-collaboration isolation*
  (owned by case 14). Neither is what this dialogue demonstrates. The plot maps
  `[AC-33]` with `INV-ASSURE-01`/`INV-RUNTIME-01`/`INV-EXEC-01`. **Recommendation:**
  correct case 20's AC list to `[AC-33]` at regeneration; confirm no coverage report
  leans on 20 to carry AC-20/AC-29 (both remain covered by 09/14).
- **`16` → AC-05 (was `[AC-04, AC-05]`).** v1.0 removed the automated scope-envelope
  check: there is no re-gate on a plan whose region exceeds the intent. Instead,
  delivering work that would reach outside the approved area is surfaced as a
  plain-language **feasibility question** (a quality gate over planning), and
  scope-safety is settled at delivery (Gate 2). Case 16 now demonstrates **AC-05**
  (a plan can't quietly execute past what was approved); **AC-04** (a plan preserves
  request detail) is not what this dialogue exercises. The plot maps `[AC-05]`.
- **`03` / `15` additions.** `review-intent-approve-derive` adds `INV-INTENT-01`
  (the dialogue *is* the intent gate) and `INV-PLAN-04` (the open question) to case
  03's invariant list; both are genuinely demonstrated. `15`'s mapping is unchanged.
- **`04` mapping unchanged.**

## Harness capability gaps surfaced this phase

Flagged inline in the affected `generated/*.case.yaml`. None blocks authoring.

1. **`no_second_approval_prompt: true`** (`approve-and-build-one-turn`) — a
   behavioral post-condition asserting the coordinator interposed no confirmation
   card / second plan-approval between approving the goal and building. The transcript
   `forbids_regex` (confirmation-card / "approve the plan") gives a first-order check
   today; a positive structural assertion would need harness support. Owned rule:
   `INV-EXEC-01`, AC-06.
2. **Scope of `approve-and-build-one-turn`.** It is authored `conversation-only` to
   isolate the *gate moment*; the actual worker build/verify/deliver is out of scope
   here and owned by case 05 and `standard-feature-whole-flow`. If the maintainer
   later wants this case to also assert a real execution record, it would move to the
   `cc-test-case` (full-execution) driver — a deliberate scoping choice, noted so it
   is not mistaken for missing coverage.
3. **Seed shapes reused as-is.** The `intents:` seed (with `open_question`, `tier`,
   `path`) and the `plans:` seed (with `intent_path`, `seed_state`) are the existing
   runner shapes from cases 03/04/15/16; no new seed capability is required for phase 2
   (unlike phase 1's knowledge/proposal overlays).
