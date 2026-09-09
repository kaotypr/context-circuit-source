# Plan 0030 — Keep the workspace agent loop; remove discovery waste

**Intent:** i015-lean-workspace-context  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done

## Objective

Keep today's correct agent loop and remove three discovery wastes: catalog
selection without tree search, one home per procedure fact with pointers
elsewhere, and skills that name files and runtime invokes without restating
contract schemas. Ship the leaner behavior in the template seed and product
adapters against the nested home already in place.

## Grounding (HEAD 36d2907)

The loop, nested home, host-native routes, durable-only knowledge, and band
allocation are already landed. The waste is instructional:

- `template/context/INDEX.md` has no exemplar entry shape or empty-catalog stop.
- `context/domains/README.md` is shipped and reads like a second index; agents
  must not fall through to it when `INDEX.md` is empty or unmatched.
- `AGENTS.md`, `WORKFLOW.md`, `coordinator.md`, and skills retell Gate 1,
  tracer, verifier, Explore, durable-only, bands, and invoke-not-read.
- `cc-intent` enumerates `intent-contract.yaml` fields as procedure.

`plan-allocate-id` returns `PLAN_PREFIX_COLLISION` on historical archive duplicate
prefixes; this plan uses **0030** as the next free in-band id.

## Decisions

**One plan, four tasks** — catalog, one-home, skill invokes, proof. Tasks 2–3
both touch skills but share one execution and verification boundary.

- **Do not** fill maintainer `context/INDEX.md` rows or rewrite domain pages.
- **Keep** `context/domains/README.md` as domain-structure knowledge; add a
  not-the-catalog banner and enforce the negative in procedure and checks.
- **Add** `test/knowledge/test-lean-context.sh` and `AC-LEAN-01..06` criteria-map
  entries; existing loop suites must stay green.

## Tasks

### LWC-001 — Make the catalog actually select

Update `template/context/INDEX.md` with entry shape and empty-catalog stop.
Align `product-knowledge.md`, coordinator, `AGENTS.md`, and retrieval skills
with catalog-only selection and explicit no-grep/no-list/no-fall-through
negatives. Banner on `context/domains/README.md` (source + template).

**Done when:** template shows entry shape; agent-facing files forbid tree search
and domains/README fall-through.

### LWC-002 — Give each procedure fact one home

Dedupe `AGENTS.md`, `WORKFLOW.md`, and `coordinator.md` by concern. Convert
skill and host-native lifecycle restatements to pointers. Host stubs stay thin
routes.

**Done when:** each listed fact has one home; pointers remain; spine is not gutted.

### LWC-003 — Make skills name files and invokes

Refactor `cc-intent` to template + validate, not schema transcription. Audit
and fix invoke lines on touched skills against `engine.sh`.

**Done when:** authoring skills avoid schema-as-procedure; invoke lines match runtime.

### LWC-004 — Prove improvement without thinning the loop

Add `test/knowledge/test-lean-context.sh`, criteria-map entries, and
acceptance wiring. Extend harness access discipline where needed. Confirm release
assembly ships the lean surface.

**Done when:** new checks fail on regressions; full `test/acceptance.sh` passes.

## Risks

- Over-trimming coordinator or AGENTS.md could drop safety spine — mitigated by
  pointer-only cuts and existing semantic suites.
- Slashing cc-intent could remove needed approval procedure — keep route-specific
  steps, drop neighbor-route and schema essays only.
- Historical archive prefix collisions block future `plan-allocate-id` until
  reconciled — documented in finding; does not block this plan id.

## Intent acceptance mapping

| Intent criterion | Plan task |
| --- | --- |
| AC-LOOP-PRESERVED | LWC-004 (existing suites + no gate drops) |
| AC-CATALOG-SELECTS | LWC-001, LWC-004 |
| AC-ONE-HOME | LWC-002, LWC-004 |
| AC-SKILL-INVOKE | LWC-003, LWC-004 |
| AC-NO-DOWNGRADE | LWC-002, LWC-004 |
| AC-PRODUCT-WORKSPACE | LWC-001, LWC-004 (template + release) |
