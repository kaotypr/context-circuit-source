---
name: cc-system-design
description: Author and structure a system design as source material under sources/system-design/ — the three-tier layout, how much detail per file, and scope separation by concern. Drafts structured source files only; never approves, accepts, plans, executes, or writes Product Knowledge.
---

## When to use

On a request to design or structure the shape of a larger change — "design the
system for X", "write up the architecture", "structure the system design" — for a
new product, a new version, or a cross-cutting feature. A bounded change needs no
system design; go straight from Product Knowledge to a plan.

A system design is **one kind of source**, not a lifecycle stage. It has no
status, no acceptance gate, and no runtime record. This skill only shapes the
files; it never approves, accepts, plans, executes, or writes Product Knowledge.

## Where it goes

Author under `sources/`, which is passive source material:

```
sources/system-design/<product-or-project>/<grouping>/<scope>/
```

- `<product-or-project>` (or `<initiative>`) names the product **as a whole —
  never a single repository**. In a single-product workspace this level may
  collapse to `sources/system-design/<grouping>/<scope>/`.
- `<grouping>` is the organizing dimension for a body of design work. **Default:
  a version** — use **3-number semver with a `v` prefix going forward**
  (`v0.6.1`, not `v0.6`); legacy 2-number folders (`v0.5`, `v0.6`) stay and are
  **not** retroactively renamed. Optionally a **named grouping** — a slug such as
  `phase-2`, `Q1`, or a milestone — as a deliberate author choice when a team does
  not organize design by release. Within one product, pick one grouping dimension
  rather than mixing release and phase folders arbitrarily.
- When the author does not specify a grouping, derive it from the version the
  design targets: the project's **next** version for forthcoming design (the
  normal case), or the **current** version when documenting as-built state. Read
  the version token from the project's own source of truth (`workspace.yaml`
  `template_version`, release manifests) — never invent it; if the targeted
  version is genuinely ambiguous, surface the choice rather than guessing.
- `<scope>` is a **concern**.

Writing these files is a normal authoring write; `sources/` stays passive for
later reads (read only when a request names a file).

## The three-tier layout

Every folder has a `README.md` index; each scope's normative content is
`design.md`; detail splits into files or sub-folders **as it grows**:

```text
sources/system-design/<product-or-project>/
  README.md                 # index of groupings (landing + reading order)
  <grouping>/
    README.md               # grouping index: the scopes + the reading order
    <scope>/
      README.md             # scope landing/index
      design.md             # scope overview — NORMATIVE, readable end to end
      <concern>.md          # one concern each, split when it outgrows a section
      <concern>/            # a concern that itself has parts (README + a file per part)
```

## How much detail per file (altitude)

Calibrate by what a reader needs from the file, not by filling it:

- **`README.md` — orientation only.** What exists here and the reading order. No
  design content, decisions, or mechanism.
- **`design.md` — the normative overview a reviewer can stop at.** The capability,
  the problem, the principles, the **fixed decisions**, and the shape of the whole
  — readable start to finish to understand the design *without* opening every
  detail file. It states *what* and *why*; it **defers** *how* to detail files and
  links to them. Not exhaustive, not code-level.
- **`<concern>.md` — one concern at implementation depth.** Its mechanism,
  interfaces, edge cases, and the diagram that explains it. Name it for the concern
  (`scheduling.md`), never `details.md`. It continues from `design.md`; it never
  restates it.

Every fact lives once: `design.md` references a detail file rather than inlining
it. If a reviewer must open five detail files to grasp the design, too much left
`design.md`; if `design.md` is many screens of mechanism, too much stayed.

## Split only when a concern earns it — do not pre-fragment

- A one-paragraph concern is a **paragraph in `design.md`**, not a file.
- A concern that grows past a section becomes its **own `<concern>.md`**.
- A concern with distinct parts becomes a **sub-folder** (a `README.md` index plus
  one file per part).
- Start minimal — a `README.md` and a `design.md` — and grow. A new design with
  five near-empty detail files is worse than one honest `design.md`.

## Separate scopes by concern, never by repository

A system design's value is describing how the pieces fit across the system, so
scopes cut by **concern** — a backend topology, a domain, a SPA, a cross-cutting
flow — matching the product, not the git layout.

- **Never one scope folder per git repository**; that fragments the cross-repo
  coherence the design exists to capture.
- If one repository needs deep design detail, that is a **scope** within the
  product's design (e.g. a `backend/` scope), not a separate top-level design and
  not a new `<product>`.

## Diagrams

Embed diagrams as fenced ` ```mermaid ` blocks directly in the file they explain —
the durable form that renders on GitHub and most editors. You may preview via a
host mermaid plugin if one exists, but never depend on it: with no plugin, still
emit the fenced mermaid.

## How it feeds the rest — the existing flow, unchanged

A system design feeds Product Knowledge and plans through the **normal v0.5 path**,
with nothing new: the coordinator gathers context from the named design source,
proposes context units through the existing context-proposal path, a human accepts
those proposals, and plans ground in the resulting Product Knowledge via the
existing `product_knowledge` references. "Accepting the design's ideas" *is*
accepting those context proposals — there is no separate design-acceptance gate.

## Boundaries

- **Draft structured source files only.** Never approve, accept, plan, execute,
  merge, or write Product Knowledge. This skill grants no route, role, or authority
  (INV-SKILL-01); it is read-as-procedure guidance resolved by path.
- **No status, no gate, no runtime record** — a system design is a source.
- **Reference model:** this repository's own `sources/system-design/` is the
  dogfooded example of the layout, altitude, and scope separation above; follow
  the same convention.
