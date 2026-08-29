# 0021 — External-service reference knowledge category

- **Plan ID:** `0021-context-references`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0003-context-lifecycle`
- **Owns:** the `context/references/` Product Knowledge category and its convention

## Original request

Retroactive plan for the v0.6.1 `context-references` scope, authored as if from
an empty repo. Source design:
`sources/system-design/context-circuit/v0.6.1/context-references/design.md`.

## Objective and desired behavior

- Introduce a first-class, plural `context/references/` category for reference
  knowledge about **external services the workspace consumes but does not own**
  (for example a third-party API a dependency repository wraps or gates).
- One entry per external service, mirroring `domains/`: a per-service
  sub-directory with a `README.md`. The directory is created when the first
  service is documented; it need not exist empty.
- Register the route in the layout index (and the blank template seed) and state
  the owned-vs-external rule once in the conventions, so an agent reads the home
  instead of improvising one.

## Constraints and non-goals

- `references/` is **plural**, per the collection-directory convention.
- Owned knowledge is **unchanged**: a wrapper/gateway repository's own knowledge
  stays in `domains/`, tagged `repositories: [<key>]`. Only knowledge about the
  external service *itself* goes under `references/`.
- **No core contract bump**; the retrieval-catalog shape
  (`wrapper/contracts/schemas/context-index.yaml`) is unchanged.
- **Non-goal:** an `external_services` index-metadata field and a
  wrapper→service binding link — both deferred; the directory stands alone.

## Product Knowledge grounding

- `conventions` (`context/CONVENTIONS.md`) — the owned-vs-external rule.
- `context-index` (`wrapper/contracts/schemas/context-index.yaml`) — confirmed
  unchanged by the directory approach.

## Tasks

1. **CTXREF-001** — register the route in the layout index and template seed;
   state the convention.

## Acceptance & verification

- The layout index and its seed list the external-references route; the
  conventions state the owned-vs-external rule.
- `sh test/acceptance.sh`.

## Assumptions, open questions, risks

- Formalizes the folder instinct rather than adding a parallel tagging scheme;
  the deferred metadata field can layer on later without reshaping this category.

## Expected commits and delivery notes

Index + conventions edits in the source and the template seed; no schema change.
