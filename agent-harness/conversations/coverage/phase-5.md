# Phase 5 coverage cross-check — candidate, completion & delivery

How the six phase-5 plots cover the existing `scenarios/` cases and the harness gaps
they surface. Descriptive only: every rule cited is owned by
`wrapper/contracts/invariants.yaml` or `test/acceptance/criteria-map.yaml`. This phase
is the back half of the lifecycle: candidate honesty, the two completion modes, the
delivery gate, and the change-set path (happy and blocked).

## Existing cases → plots

| existing case | plot | relationship |
| --- | --- | --- |
| `06-repair-then-complete` | `critical-repair-then-complete` | **evolves** — same seed/budget; Critical explicit completion |
| `08-delivery-boundary` | `delivery-boundary-block-no-remote` | **evolves** — same seed/budget; Gate 2 block |
| `17-stale-candidate-requires-reverification` | `stale-candidate-refuse-complete` | **evolves** — same seed/budget; AC corrected (below) |
| `18-standard-delivery-inferred-completion` | `standard-inferred-completion` | **evolves** — same seed/budget |
| `21-change-set-one-verification` | `change-set-one-verification` | **evolves** — same seed/budget |
| *(none)* | `change-set-base-unbuildable` | **new** — the honest-block branch of the change-set path |

## AC mapping drift (surfaced, not hidden)

- **`17` → `[AC-14, AC-32]` (was `[AC-14, AC-15, AC-28]`).** Case 17 is the
  candidate-staleness scenario: a change after the check yields a new candidate that
  voids prior evidence, and completion refuses stale evidence — this is **AC-32**
  exactly ("any new commit or criteria change yields a new candidate that voids prior
  evidence and acceptance (completion refuses stale evidence)"). The historical tags
  are weaker fits: **AC-15** is *explicit human completion after verification* (the
  happy completion, case 06), and **AC-28** is *the pull-request/delivery-blocks
  criterion* (case 08). The plot keeps **AC-14** (verified is not complete, genuinely
  demonstrated) and adds **AC-32**. **Recommendation:** correct case 17's AC list to
  `[AC-14, AC-32]` at regeneration; AC-32's suite (`test/candidate/test-candidate.sh`)
  remains the deterministic owner.
- **`06`, `08`, `18`, `21` mappings unchanged** — all four were already correct
  (AC-15/22 completion, AC-16/28 delivery, AC-35 inferred completion and change set).

## The new plot: `change-set-base-unbuildable`

The negative branch of the change-set path (design conversation C's tail): two
individually-verified members whose combination won't build. The mechanic is owned by
**INV-CONCURRENCY-02** ("a base that cannot be built cleanly is a blocked execution,
not a worker failure") plus **INV-DELIVER-01** and **INV-PRESERVE-01**; **AC-35** is
the change-set family anchor. Reported as an honest block — nothing merged, nothing
completed, both members preserved, offer to split or reorder — never an invented
worker failure or a silent partial ship.

## Harness capability gaps surfaced this phase

Flagged inline in the affected `generated/*.case.yaml`. None blocks authoring.

1. **New fault: `change-set-base-unbuildable`** — forces `BASE_UNBUILDABLE` at
   `change-set-prepare`, mirroring case 09's `fault: verifier-unavailable`. Lets the
   probe exercise the block without hand-crafting a genuinely unbuildable merge.
2. **New post-condition predicates** (`change-set-base-unbuildable`):
   `change_set_base_unbuildable`, `change_set_blocked_not_worker_failure`,
   `no_change_set_completed`, `members_preserved`. These assert the *conversational /
   state* outcome of the block; the deterministic change-set mechanics are owned by
   `test/completion/test-inferred.sh` and the candidate suite.
3. **No drift in seed shapes.** The five rewires reuse the exact seed states the
   existing cases already define (`verified-after-repair`, `verified`,
   `verified-stale`, `verified-accepted`, two `verified` members); no new seed
   capability is needed for them.
