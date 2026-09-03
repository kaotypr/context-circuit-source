# Phase 3 coverage cross-check — Explore tier & promotion

How the two phase-3 plots cover the existing `scenarios/` cases and the harness gaps
they surface. Descriptive only: every rule cited is owned by
`wrapper/contracts/invariants.yaml` or `test/acceptance/criteria-map.yaml`. Explore
is the bottom rung of the one assurance ladder (INV-ASSURE-01), so this phase
continues the tier material from phase 2 (`tier-fails-upward-refuse-explore`).

## Existing cases → plots

| existing case | plot | relationship |
| --- | --- | --- |
| `14-codex-direct-collaboration` | `direct-collaboration-explore` | **evolves + host-neutralizes** — same seed/budgets; provider/model specifics removed |
| *(none)* | `explore-promote-to-standard` | **new** — the promote ramp into the trust pipeline |

## Host-neutralization (case 14)

Case 14 is titled and written around a specific host (Codex, the "Luna" worker,
`gpt-5`), and its marker content is `paired by Codex`. Per INV-HOST-01 (host identity
is bounded evidence that never surfaces) the plot removes all of it:

- marker content is provider-neutral: **`paired locally`**;
- the `forbids_regex` keeps generic model/effort/`cc-pair`/`verifier` guards but drops
  the host-specific `gpt-5|luna` tokens — a host-neutral case shouldn't hard-code one
  provider's names into the assertion;
- `pair_role_evidence_configured: worker` / `pair_role_absent: verifier` already read
  host-neutrally and are kept.

The mechanics (isolated working copy, uncommitted-on-close, no verifier, planless)
are unchanged; only the provider dressing is removed. **Recommendation:** when
regenerating, consider renaming the scenario directory from
`14-codex-direct-collaboration` to a host-neutral name so the case name itself stops
implying a provider.

## AC / invariant mapping

- `direct-collaboration-explore` keeps case 14's `[AC-29]` / `INV-PAIR-01`,
  `INV-HOST-01`, `INV-RUNTIME-01` verbatim — the mapping was already correct.
- `explore-promote-to-standard` maps `[AC-29, AC-30, AC-33]` with `INV-PAIR-01`
  (promote clause), `INV-ASSURE-01` (the independent check appears when the tier is
  raised), `INV-INTENT-01`, `INV-APPROVE-01`. The promote step is the one place all
  three of pairing (AC-29), first-class intent (AC-30), and tiering (AC-33) meet.

## Harness capability gaps surfaced this phase

Flagged inline in the affected `generated/*.case.yaml`. None blocks authoring.

1. **`promoted_in_place: true`** (`explore-promote-to-standard`) — a predicate
   asserting the authored plan of record binds to the *existing pairing-branch
   commits* (a candidate derived from work already done) rather than a fresh redo.
   This is the essence of "the cliff is a ramp" (INV-PAIR-01 promote clause) and has
   no post-condition predicate today.
2. **Live-authored intent/plan ids.** `explore-promote-to-standard` runs the whole
   arc in-session, so the intent/plan ids are allocated live; the id-keyed
   assertions assume the first-allocated `i0001-…` / `0001-…`. Same id-agnostic
   concern raised by `standard-feature-whole-flow` and phase 1 — a fresh whole-arc
   case wants id-agnostic assertions.
3. **No seed shape for a mid-pairing session.** Rather than seed pairing-branch
   commits (which the runner can't do today), the plot plays a short pairing beat
   then promotes, all in one full-execution arc. Noted so the continuous-arc shape is
   understood as deliberate, not accidental.
4. **`no_delivery_records: true`** — asserts the "don't ship it yet" hold after
   promotion (no pull request / delivery). Owned rule: `INV-DELIVER-01`. Confirm the
   grader vocabulary has (or gains) this predicate distinct from `no_execution_records`.

## Note on driver

Both phase-3 plots are `full-execution` (`cc-test-case` driver): `cc-pair` spawns a
real worker child, and promotion additionally exercises intent/adversary/plan. This is
the first phase whose plots are not conversation-only; dimension C (access discipline)
degrades to warning-only under the in-session driver, as documented in
`scenarios/README.md`.
