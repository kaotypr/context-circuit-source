---
id: 0027-system-design-grouping
target_context_unit: system-design-authoring
operation: change
statement: >
  The system-design source path's middle segment is a grouping dimension, not a
  strict version. The v0.6.1 design-layout-grouping scope generalized
  sources/system-design/<product>/<version>/<scope>/ to
  sources/system-design/<product-or-project>/<grouping>/<scope>/, where the
  grouping defaults to a version (now 3-number semver with a v prefix going
  forward, e.g. v0.6.1) but may be a named grouping (phase-2, Q1, a milestone) as
  an explicit author choice. When the author does not specify one, the agent
  derives the grouping from the version the design targets — the next version for
  forthcoming design, the current version for as-built state — read from the
  project's own source of truth, never invented. Legacy v0.5/v0.6 folders stay
  as-is (no retroactive rename). The cc-system-design skill already carries the
  generalized rule; the system-design-authoring domain page still shows the old
  <version> path and should be updated to match.
evidence_refs:
  - sources/system-design/context-circuit/v0.6.1/design-layout-grouping/design.md
  - context/domains/system-design-authoring/README.md:47
  - context/domains/system-design-authoring/README.md:68
  - .agents/skills/cc-system-design/SKILL.md
affected_repositories:
  - context-circuit
affected_commits:
  - c67dd23  # docs(cc-system-design): generalize design-source path segment to grouping
related_plan: 0023-design-layout-grouping
confidence: high
status: review-needed
---

# Generalize the system-design path segment to a grouping

## What is stale

The `system-design-authoring` domain page pins the path's middle segment to a
strict version:

- line 47 — "A system design lives at
  `sources/system-design/<product>/<version>/<scope>/` …"
- line 68 — "Output: authored files under
  `sources/system-design/<product>/<version>/<scope>/`"

## What changed (v0.6.1 · design-layout-grouping, plan 0023)

Two gaps surfaced in use: the middle level isn't always a version (a team may
organize by delivery phase, quarter, or milestone), and version depth was
inconsistent (2-number `v0.5`/`v0.6` folders vs. 3-number semver on every shipped
surface). The convention was refined to:

```
sources/system-design/<product-or-project>/<grouping>/<scope>/
```

- **`<grouping>`** — the organizing dimension. Default: a version, in 3-number
  semver with a `v` prefix going forward (`v0.6.1`). Optionally a named grouping
  (`phase-2`, `Q1`, a milestone) as a deliberate author choice; pick one grouping
  dimension per product rather than mixing.
- **Agent default derivation** — when no grouping is given, derive it from the
  version the design targets (next for forthcoming, current for as-built), read
  from the project's own source of truth (`workspace.yaml` `template_version`,
  release manifests); surface the choice when genuinely ambiguous rather than
  guessing.
- **`<product-or-project>`** may collapse to one segment in a single-project
  workspace; still never one folder per git repository.
- **No retroactive rename** — legacy `v0.5`/`v0.6` stay; the 3-number rule binds
  new groupings only, so both depths may coexist during transition.

The `cc-system-design` skill already carries this generalization (commit
`c67dd23`); the domain page is the surface that lags.

## Proposed change

Update the two path references (and the "A source, not a stage" / "Interfaces"
prose) to use `<product-or-project>/<grouping>/<scope>/`, and add a short
sentence naming the version-default, 3-number-going-forward rule, the named-
grouping option, and the agent's current/next default derivation. Keep the
three-tier layout, detail-per-file rubric, and scope-by-concern rule unchanged.

## Why not more

This refines a layout convention only; it adds no lifecycle, status, gate,
runtime, schema, or new owner, and no core contract bump. On acceptance, refresh
the domain's `source_revisions` commit pin (currently `cb84870`) to the revision
carrying `c67dd23`.
