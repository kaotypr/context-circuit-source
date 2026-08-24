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
