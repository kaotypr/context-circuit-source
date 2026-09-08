# Phase 6b coverage cross-check — peripheral surfaces & whole-flow

How the four phase-6b plots cover the surface and the harness gaps they surface.
Descriptive only: every rule cited is owned by `.context-circuit/wrapper/contracts/invariants.yaml` or
`test/acceptance/criteria-map.yaml`. This phase completes the whole-surface catalog:
the two external-publication kinds, the system-design→intents path, and a live case for
the whole-flow plot.

## Plots → coverage

| plot | status | proves |
| --- | --- | --- |
| `publish-plan-to-external` | **new** | cc-publish `plan` kind: export-only, self-contained, preview-first, orthogonal |
| `publish-open-questions-thread` | **new** | cc-publish `thread` kind: one self-contained message per open question |
| `author-system-design-spawns-intents` | **new** | cc-system-design drafts source, spawns one intent per concern |
| `standard-feature-whole-flow` | **case for existing plot** | the whole lay-user flow across both gates, end to end |

No existing `scenarios/` case maps to any of these — the external surface, the
system-design path, and the whole-flow-as-one-live-conversation had no live case.

## AC gaps (honest)

- **The two publish plots map `acceptance_criteria: []`.** The external surface is a
  v0.7 addition owned by **INV-EXTERNAL-01/02/03** (plus **INV-PLAN-01** for "never
  changes plan status"); the v0.5-rooted `criteria-map.yaml` has no publish AC. Same
  honest pattern as `lease-ownership-conflict` and `accept-or-defer-context-proposal` —
  invariant-owned, no discrete AC. If the maintainer wants an AC handle for the external
  surface, that is a criteria-map decision, not a library one.
- **`author-system-design-spawns-intents` maps `[AC-30]`** — the "genuinely separate
  decisions become separate intents" property is the first-class-intent criterion; the
  system-design authoring itself is source drafting (INV-SEC-02 / INV-SKILL-01), not a
  gated lifecycle step.

## The whole-flow case

`standard-feature-whole-flow` is the *(opt)* whole-flow plot authored early as a worked
example; this phase generates its `case.yaml`. It is the one case whose post-conditions
are **id-agnostic** — a fresh whole-arc conversation authors its intent/plan ids live,
so it cannot use the id-keyed assertions the seeded cases use. This is the harness
capability the plot itself flags (and which `explore-promote-to-standard` also needs).

## Harness capability gaps surfaced this phase

Flagged inline in the affected `generated/*.case.yaml`. None blocks authoring.

1. **id-agnostic post-conditions** (`standard-feature-whole-flow`): `intents_approved`,
   `plans_created`, `all_plans_done`, `deliveries_recorded`, `completion_kind`,
   `knowledge_debt_pending` asserted without fixed ids. The single biggest recurring
   harness gap across the library (also whole-flow, `explore-promote-to-standard`).
2. **Preview-without-a-provider** (both publish plots): the cases are authored as
   preview/consult conversations so they need no live external provider; `no_provider_push`
   and `plans_untouched` assert nothing was sent and nothing was written under `plans/`.
   The deterministic external-record mechanics are owned by the publication schemas /
   any external suite; these probes cover the lay-user-facing discipline (export-only,
   self-contained, orthogonal, consult-first).
3. **New authoring/state predicates**: `system_design_authored` and `intents_created`
   (count) for `author-system-design-spawns-intents`; `no_provider_push` /
   `plans_untouched` for the publish plots.
4. **Sources authoring is a legitimate write** (`author-system-design-spawns-intents`):
   its `design` phase writes under `sources/system-design/**` — the one place a plot's
   access policy *allows* writing under `sources/`, because authoring a design there is a
   normal write (the passive-`sources/` rule, INV-SEC-02, is about reads/scans).
