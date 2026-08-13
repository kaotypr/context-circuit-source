# Product Knowledge

Product Knowledge is a lean, generated, and maintained view of the business,
role, domain, and workflow context of a product. It gives humans and agents
task-relevant business context without loading or modeling the entire
application. It is optional and adopted incrementally: a wrapper with no Product
Knowledge stays valid, and coverage can grow one page at a time.

## Layout

Product Knowledge lives as human-reviewable Markdown under `context/`:

- `context/PROJECT.md` — the compact product map. It routes a reader to the right
  role or domain without loading the whole tree.
- `context/GLOSSARY.md` — business terms.
- `context/roles/README.md` and `context/roles/<role>.md` — one page per business
  role, owning the cross-domain user perspective and end-to-end role story.
- `context/domains/<domain>/README.md` — a domain summary that owns business-area
  routing.
- `context/domains/<domain>/workflows/<workflow>.md` — the authoritative page for
  exact current behavior, variations, and business rules.

Authority never overlaps: role pages own the cross-domain story, domain summaries
own routing, and workflow pages own exact current behavior and rules. Role,
domain, and product-map pages may not contain a `## Current flow` or
`## Business rules` section — those belong to workflow pages. Validation enforces
required sections, resolvable relative references, present owners and sources, and
this ownership split. The templates under `.agents/templates/product-knowledge/`
are a valid worked example.

Each page carries small YAML frontmatter (kind, title, owners, sources, review
date, references, and known gaps). Product Knowledge describes current effective
behavior; plans describe proposed behavior, and implementation observations never
silently redefine business truth.

## Lifecycle

1. **Baseline.** A fresh `configure-workspace` bootstrap request may include
   `context.product_knowledge` to create a minimal reviewed baseline — product
   purpose, major roles, major domains with any known workflows, sources, and
   explicit unknowns. Nothing is invented; unknowns stay explicit.
2. **Discovery.** `gather-context` traverses the product map to only the relevant
   pages for a request and records credential-free candidates, contradictions, and
   gaps as ignored runtime evidence under `.runtime/product-knowledge/`. Discovery
   never edits canonical pages.
3. **Planning.** A plan may declare `product_knowledge` — the role, domain, and
   workflow pages it references, an `impact` (`none`, `documentation-correction`,
   `implementation-only`, `behavior-change`, `new-workflow`, `retired-workflow`),
   and a `proposed_change`. Planning describes proposed behavior only.
4. **Execution.** `run-task` carries a compact, immutable task-context package —
   the bounded resolved pages plus a content digest and revision — into the task
   brief, so a worker and its independent verifier share the same business
   baseline without the whole tree. Planless requests may declare the same
   package.
5. **Evidence.** Worker and verifier results report whether Product Knowledge
   impact is absent, matches the declared change, is unexpectedly broader, or
   contradicts current context. Closeout records the impact; unexpected impact
   withholds synchronization for human review.
6. **Effective synchronization.** Canonical current-behavior pages change only
   through `sync-context` after a dedicated, workspace-configured
   `product_knowledge.confirming_role` — distinct from the per-page owners and the
   merge-gate human — declares the behavior effective. Merged but disabled or
   unreleased behavior stays proposed or recorded in contribution evidence.
7. **Onboarding.** `onboarding-pack --roles <role[,role...]>` generates a
   revision-stamped, source-listed onboarding view. The pack is reproducible from a
   single revision and clearly marked as a generated view, not a source of truth.

## Migration for existing wrappers

Existing wrappers need no migration. Product Knowledge is additive: without a
`context/roles/` or `context/domains/` tree, validation treats Product Knowledge as
absent and the workspace stays valid. To adopt it, add pages incrementally — start
with `PROJECT.md` as a product map and one role or one workflow — and let
just-in-time discovery grow coverage. To enable effective synchronization, add
`product_knowledge.confirming_role` to `workspace.yaml`.

The canonical rules live under `.agents/`; Codex and Claude host adapters remain
thin delegates.
