# Plan 0029 — Align contracts, docs, and checks with band allocation

**Intent:** i021-id-number-blocks  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done  
**Depends on:** 0028-band-allocation

## Objective

Product rules, schemas, invariants, skills, docs, template seed, and acceptance
checks match implemented band allocation — short ids unchanged, branches stay
`cc/<plan-id>/<repo>`, overlapping bands and duplicate prefixes fail closed.

## Grounding (HEAD 7db954f)

Contracts and docs today describe a single workspace-wide sequence. Plans 0027–0028
implement roster, identity, and band-scoped runtime behavior. This plan aligns
every owned product surface and maps AC-BAND-01 through AC-BAND-09 to runnable checks.

## Decisions

**Stack plan 3 of 3** — docs and contracts follow working runtime so verification
proves product behavior, not aspirational wording.

- Update INV-INTENT-01 and INV-PLAN-03 in place; add roster invariant if needed.
- Refresh skills, getting-started, plan-review, and repository-binding domains.
- criteria-map.yaml gains AC-BAND entries pointing at test/bands/ suites.

## Tasks

### BC-001 — Invariants and contract schemas

Revise invariants and intent/plan schema notes for band allocation semantics.

**Done when:** test/contracts/test-contracts.sh passes with updated roster references.

### BC-002 — Skills, docs, template, and Product Knowledge

Align cc-intent, cc-trace, cc-plan, cc-workspace, coordinator, engine-and-seam,
domain pages, and template seed.

**Done when:** docs-skills band contract check passes; no block-number prompts documented.

### BC-003 — Criteria map and full acceptance

Map AC-BAND-01–09 to runnable commands; wire band suites; full semantic acceptance green.

**Done when:** criteria-map complete; test/acceptance.sh and test/release/test-release.sh pass.

## Risks

- Stale global-sequence wording in context/ could contradict new invariants if any
  page is missed during the doc pass.
- Release manifest must list members.yaml and deny member.local.yaml explicitly.

## Acceptance criteria mapping

| Intent criterion | Task evidence |
| --- | --- |
| AC-BAND-01–09 | BC-003 criteria-map and band contract suite |
| AC-BAND-03 | BC-002 skills and getting-started |
| AC-BAND-07 | BC-002 branch guidance unchanged |
