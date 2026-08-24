# Decisions

## 2026-08-21 — wrapper/template split

Decision: keep shipped wrapper ownership under `wrapper/` and the blank mutable
seed under `template/`.

Rationale: source identity must not be confused with an instantiated workspace.
Consequence: release assembly overlays root adapters and template state while
preserving no user data.

## 2026-08-21 — separate lifecycle gates

Decision: review, approval, execution, completion, delivery, archive, takeover,
and cleanup remain separate human actions.

Rationale: eligibility is not authorization and runtime evidence is not Done.
Consequence: the router emits an exact gate or read-only recommendation.

## 2026-08-24 — built-template harness location

Decision: the source-only built-template behavior laboratory lives at top-level
`template-harness/`, relocated and renamed from `test/template-runtime/`.

Rationale: it assembles and drives the released product, distinct in kind from
the `test/` engine and contract suites; the name reflects its assemble → drive →
grade role.
Consequence: the release allowlist still excludes it (it cannot ship);
`test/lib/assert.sh` resolves the repo root git-based so suites may live at any
nesting. Accepted from proposal `0002-template-harness-and-terminology-decision`.

## 2026-08-24 — terminology as a context category

Decision: Terminology is a standard context category. A per-project glossary
seed lives at `context/TERMINOLOGY.md`; the canonical Context Circuit glossary
and its internal → user-facing translation live in `docs/terminology.md`,
projected from design chapter 08, and the coordinator references that projection
instead of hardcoding the term list.

Rationale: one owner per rule — chapter 08 owns term meaning and translation,
`docs/terminology.md` is the shipped projection, the coordinator role references
it rather than duplicating policy.
Consequence: the shipped template gains `context/TERMINOLOGY.md`,
`docs/terminology.md`, and `docs/templates/terminology-context.md`. Accepted from
proposal `0002-template-harness-and-terminology-decision`.

## 2026-08-24 — current-state context refresh

Decision: after `plans/context-circuit-plans/` was deleted in commit `4b8ac0b`
as obsolete previous-version plans, re-ground the entire `context/` Product
Knowledge on the current wrapper (skills, `engine.sh`, schemas, `invariants.yaml`,
docs/adapters) rather than on the deleted plans, and accept the result.

Rationale: the domain/role pages were reverse-engineered from previous-version
plans and had drifted from the shipped contract (consolidated invariant ids,
renamed skills, deleted `routes.yaml`/`context-sets.yaml`/`docs/gates.md`, and an
approval card/token mechanism that contradicted `INV-APPROVE-01`). The current
wrapper is the authoritative present-day evidence.

Consequence: the domain set is now the nine domains that mirror the shipped
lifecycle — repository-binding (broadened with orientation), plan-review
(broadened to planning), plan-approval, plan-execution, verification,
completion, plan-organization, delivery, and host-adapters — plus the refreshed
maintainer role, INDEX owner pointers, and terminology authority pointer.
`context/sources.yaml` provenance was retired (the wrapper is not a `sources/`
read). The proposals for this refresh were consumed from `context/proposals/`
on acceptance, as usual.

## 2026-08-24 — design↔context reconciliation

Decision: reconcile the accepted `context/` Product Knowledge against the v0.5
design under `sources/context-circuit-v0.5-design/`, under the policy
"shipped is truth; log deltas". Keep "wrapper" as an accepted synonym for the
universal project workspace product and update the design source to accept it
(rather than aligning context to the design's deprecation).

Rationale: the current shipped state is built toward that design, so context
should describe shipped behavior while design↔implementation divergences stay
visible instead of silently overwriting either side.

Consequence: accepted the reconciliation proposals. Fixed contradictions in
ARCHITECTURE.md (removed the two-stage-router-as-canonical framing and the
"lease" wording; added the execution model), PROJECT.md (stated the universal
multi-repository workspace identity; removed the "operating system" metaphor and
router framing), and INDEX.md (retrieval-catalog role). Enriched the lifecycle domain pages from the shipped wrapper. Added the
`source-release-and-upgrade` domain, a `context/DESIGN-DELTAS.md` log (evidence
layers, repair-limit, plan-id reuse, terminology authority), and CONVENTIONS
policy-change escalation. The design source (`08-terminology.md`,
`09-source-and-template.md`) was updated to accept "wrapper", and the glossaries
gained the term.
