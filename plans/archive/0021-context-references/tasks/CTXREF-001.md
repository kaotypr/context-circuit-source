---
id: CTXREF-001
title: Register the references category in the layout and conventions
repositories: [context-circuit-source]
paths: [context/INDEX.md, context/CONVENTIONS.md, template/context/INDEX.md, template/context/CONVENTIONS.md]
depends_on: []
acceptance: [CTXREF-AC-01, CTXREF-AC-02]
verification: [CTXREF-VT-01]
---

## Intended behavior

Product Knowledge gains a defined home for reference knowledge about external
services the workspace consumes but does not own. The category is
`context/references/`, plural like every other collection directory under
`context/` (`domains/`, `roles/`, `proposals/`). It holds one entry per external
service — a per-service sub-directory with a `README.md` — and is created when
the first service is documented rather than shipped empty.

## Concrete change and affected surfaces

Add the external-references route to the layout retrieval index
(`context/INDEX.md`) and its blank template seed (`template/context/INDEX.md`),
so the category is discoverable and ships in fresh workspaces. State the
convention once in `context/CONVENTIONS.md` (and the seed
`template/context/CONVENTIONS.md`): external-service reference knowledge lives
under `context/references/`; a wrapper/gateway repository's own knowledge stays
in `domains/`, tagged `repositories: [<key>]`. No change to the retrieval-catalog
schema.

## Inputs and outputs

Inputs: the owned-vs-external distinction and the existing by-concept layout.
Outputs: a registered, documented `references/` category. No index-metadata
field is added; the directory approach stands alone.

## Context references

- sources/system-design/context-circuit/v0.6.1/context-references/design.md

## Risks and open questions

- A deferred `external_services` metadata field could later let an ordinary unit
  cross-link a service; it must be able to layer on without reshaping the
  directory category introduced here.

## Stop conditions

- Do not add a schema field, a tagging scheme, or a core contract bump.
- Do not move any owned (`domains/`) knowledge.
