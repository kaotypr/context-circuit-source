# Intent i0002 — make the conversation-spec library live

_Status: draft (awaiting Gate 1). Tier: standard. Scope: `template-harness/`.
Depends on: i0001-conversation-spec-library._

## The bigger picture

The conversation library under `template-harness/conversations/` — authored under
**i0001** (39 plots, 39 generated cases, a coverage doc per phase, a flow guide) — is
still a **proposal that sits parallel to** the 22 live scenario cases. Until each
scenario case is *generated from* a plot, the plots are documentation that can drift
from the harness, which is exactly the failure this library exists to prevent.

This intent is the one decision that closes that gap: **adopt the library as the
generative source and wire it in**, so a plot is the source and its case is a build
artifact, and editing the plot first is the only way to change a case.

## What a correct change achieves

- A deterministic **plot → case generator** that provably *derives* cases from plots
  (not echoes committed ones), is mutation-sensitive, and is byte-deterministic.
- The **harness-capability gaps** the library surfaced are closed — first
  **id-agnostic post-conditions** (without which the flagship whole-arc cases cannot
  run at all), plus the new seed states, faults, and predicates the coverage docs
  enumerate — each actually *exercised* and *falsifiable*, never a no-op that always passes.
- The **22 existing cases are regenerated** from their plots with seeds and budgets
  provably unchanged, assertion power non-regressed, and the acceptance suite green.
- The **five acceptance-criteria corrections** the library caught (cases 01, 16, 17,
  19, 20) are mechanically confirmed in the live case files.
- Every **frozen net-new case runs** against the live coordinator to a recorded grade.
- A **bidirectional drift guard** gates the suite: a hand-edited case fails, and a plot
  edited without regeneration fails — the edit-the-plot-first guarantee becomes real.

## Deliberately out of scope

No product behavior, runtime, contract, or invariant change (the library stays
descriptive). The deterministic five-family suite (AC-36) stays the owner of the
end-to-end families. The 22 cases are evolved, not recreated. The file-rename /
diagram / flow-metadata polish is deferred. Delivery (merge/push/publish) is a
separate Gate 2 action. No typed-language runtime port.

## Shape of the work (one intent, several plans)

One decision realized as a small stack of plans, each naming this intent and staying
inside `template-harness/`: (1) the generator + fidelity/mutation/determinism checks;
(2) the harness-capability gaps, id-agnostic post-conditions first, with negative
fixtures; (3) regenerate the 22 cases + mechanically confirm the five corrections;
(4) run the frozen net-new cases + add the bidirectional drift guard. The exact plan
breakdown and the frozen net-new plot list are authored under `cc-plan` after approval.

## Assurance

Standard tier: a single repository, reversible, test-infrastructure change with real
novelty (a new generator and grader capabilities) — not Explore-eligible, so it gets
an independent verifier, but it carries no security, money, migration, production, or
irreversibility signal, so it is not Critical.

## Spec adversary

An independent spec adversary challenged the first draft of the criteria and returned
`criteria_sound: needs-work` with material findings (a circular/single-instance
fidelity check, "reaches a grade" measuring liveness not assertion power, an undefined
"id-agnostic", a drift guard that excluded the very seed/budget region it should
freeze, and the plot→case direction left unguarded). The criteria above were rewritten
to close all of them; see `adversary.md`. Because the criteria changed, a fresh
adversary pass on the revised criteria is required before Gate 1 (INV-CANDIDATE-01 /
the adversary-record rule) — that re-challenge is the immediate next step.

## Dogfooding note

Authored source-only as draft artifacts for review; the runtime was not run, so no id
was machine-allocated, no `contract_digest` was frozen, and no `.runtime` state exists.
Approval (Gate 1) is yours to give, after the re-challenge; freezing and plan
derivation happen only then.
