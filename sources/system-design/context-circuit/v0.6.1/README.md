# Context Circuit v0.6.1

Source design for Context Circuit v0.6.1 — a **delta on v0.6**. Read the v0.6
design first ([../v0.6/](../v0.6/)) and the v0.5 base it builds on
([../v0.5/](../v0.5/)); v0.6.1 changes only what it names. This README is the
version index; each scope owns its own design.

## Scopes

- [context-references/](./context-references/) — a first-class, **plural**
  `context/references/` category for reference knowledge about **external
  services the workspace consumes but does not own** (a third-party API a
  dependency/wrapper repo gates). Distinguishes an *owned* wrapper repo's
  knowledge (stays in `domains/`, tagged by `repositories`) from the *external
  service itself* (no repo, no `repositories` tag — the void this scope fills).
  Start at [context-references/design.md](./context-references/design.md).
- [writer-brief-placement/](./writer-brief-placement/) — relocate the runtime
  `writer-brief.md` template out of the user-facing workspace root (where it
  reads as clutter next to the real entry docs) into an unobtrusive machinery
  location, re-pointing the engine lookup and release manifest. Source of truth
  and repository-grounding behavior unchanged; no core contract bump. Start at
  [writer-brief-placement/design.md](./writer-brief-placement/design.md).
- [design-layout-grouping/](./design-layout-grouping/) — refine the system-design
  layout convention: the middle path segment is a *grouping* dimension (version
  by default, optionally `phase-N`/`Q1`/milestone), version groupings use
  3-number semver going forward, and the agent defaults the grouping to the
  project's current/next version. A delta on v0.6
  [system-design-authoring](../v0.6/system-design-authoring/). Start at
  [design-layout-grouping/design.md](./design-layout-grouping/design.md).
- [worker-role-naming/](./worker-role-naming/) — unify the implementing execution
  role's name on **worker**, retiring **writer** as a synonym across roles,
  contracts, runtime, adapters, skills, and Product Knowledge (including the
  "writer brief" → "worker brief"). A pure rename: no behavior, authority, gate,
  or contract semantics change, and no core contract bump. Coordinates with
  [writer-brief-placement](./writer-brief-placement/) on the same brief file.
  Start at [worker-role-naming/design.md](./worker-role-naming/design.md).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. Every folder has a
`README.md` landing/index; each scope's normative design is `design.md`, with
detail split into files or sub-folders as it grows. See the v0.6 and v0.5 version
READMEs for the full rule.

## Authority

v0.6.1 adds no owner and duplicates no rule. It specifies intent and points to
the canonical owners under `wrapper/` and the Product Knowledge layout under
`context/`. It introduces one new context-layer category and its convention; it
mandates no core contract bump.
