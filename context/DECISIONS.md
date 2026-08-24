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
