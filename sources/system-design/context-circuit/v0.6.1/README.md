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
