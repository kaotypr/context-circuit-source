# Plan 0026 — Context files hold only durable product knowledge

**Intent:** i019-durable-context-knowledge  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** draft

## Objective

When an agent generates or updates Product Knowledge, live `context/` files contain
only durable knowledge about the product. They never name a particular plan, intent
file, or sources file, and never name a `sources/` path. `DECISIONS.md` records
what is now true about the product — not edited file paths and not the ephemeral
artifact behind the change. Existing maintainer context and the blank template seed
are rewritten; automated checks prove the rule.

## Grounding (HEAD 9ddcc59)

INV-KNOWLEDGE-01/02 govern retrieval and in-place updates but do not ban ephemeral
citations. Today:

- `context/sources.yaml` exists (retired ledger); `context/SOURCES.md` and the
  template copy instruct provenance recording there and name `sources/` paths.
- Twenty-six live context files contain `intent/`, `plans/`, or `sources/` strings —
  mostly Provenance footnotes, `DECISIONS.md` history, and `INDEX.md` pointers.
- `scripts/release-manifest.txt` still requires `context/sources.yaml` in the
  shipped seed.
- `test/knowledge/` covers INV-KNOWLEDGE-02 debt only; nothing scans context for
  forbidden citations.

i015 (draft) explicitly non-goals rewriting Product Knowledge pages; this plan is
the dedicated rewrite.

## Decisions

**One plan, four tasks** — contract/procedure, maintainer rewrite, template/release,
acceptance. Stacked plans would split one verification boundary without independent
value.

- **INV-KNOWLEDGE-03** — new invariant for durable-only citations and DECISIONS shape.
- **Delete** `context/sources.yaml` and `template/context/sources.yaml`; rewrite
  `SOURCES.md` without `sources/` path literals.
- **Retire** domain Provenance footnotes that name particular plans or sources paths;
  keep Implementation references to shipped contract paths.
- **Path patterns** like `intent/<id>/` for explaining product structure remain allowed;
  particular ids like `intent/i019-…` or `plans/0006-…` do not.

## Tasks

### DCK-001 — Add durable-only invariant and write paths

Add INV-KNOWLEDGE-03; update `product-knowledge.md`, `context-index.yaml`,
`domain-context.md` template, `CONVENTIONS.md` (maintainer + template), coordinator,
and `cc-complete`.

**Done when:** invariant and owner docs exist; coordinator and cc-complete carry the
same rule for gather-context and mark-done reconcile.

### DCK-002 — Rewrite maintainer live context

Delete `context/sources.yaml`. Rewrite `SOURCES.md`, `DECISIONS.md`, `INDEX.md`,
terminology, roles, and all domain pages to remove particular plan/intent/sources
cites and provenance footnotes naming `sources/` paths.

**Done when:** no maintainer context file names a forbidden cite; DECISIONS records
knowledge changes; concept explanations remain.

### DCK-003 — Template seed and release manifest

Remove `sources.yaml` from template and release manifest; align template `SOURCES.md`
and `CONVENTIONS.md`.

**Done when:** release assembly and template seed match the durable-only shape.

### DCK-004 — Acceptance checks

Add `test/knowledge/test-durable-context.sh`; wire into `test/acceptance.sh` and
contract roster.

**Done when:** suite fails on regressions; full semantic acceptance stays green.

## Risks

- **Over-scrubbing** — removing path patterns like `intent/<id>/` that explain product
  structure. Mitigation: verification asserts concept pages still mention Gate 1 /
  passive source material.
- **DECISIONS history loss** — rewriting entries may drop useful rationale.
  Mitigation: preserve decision/rationale/consequence prose; drop only ephemeral pointers.
- **Release break** — manifest still requires `sources.yaml`. Mitigation: DCK-003
  updates manifest and release tests in the same plan.

## Verification summary

| Check | Command |
| --- | --- |
| Invariant present | `rg INV-KNOWLEDGE-03 invariants.yaml product-knowledge.md` |
| No forbidden cites | `test/knowledge/test-durable-context.sh` |
| Release seed | `sh test/release/test-release.sh` |
| Full suite | `sh test/acceptance.sh` |

## Product Knowledge impact

Reconcile on completion: `CONVENTIONS.md`, `DECISIONS.md`, `INDEX.md`, `SOURCES.md`,
domain pages with retired provenance footnotes, and `TERMINOLOGY.md`.
