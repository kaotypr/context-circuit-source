---
name: cc-system-design
description: Author a system design for this source checkout as source material under sources/system-design/ — the three-tier layout, how much detail belongs in each file, and why scopes separate by concern rather than by repository. Drafts structured source files only; never plans, executes, publishes, or writes knowledge.
---

# Design a system before building it

Use this on a request to work out the shape of something larger than a single
change to Context Circuit — a new command surface, a change to how the two
products are assembled or published, a reworking of the workspace contract.
It is for settling that shape, legibly, before any of it is implemented.

What it produces is **source material**: passive, durable, written to be read by
a person deciding whether the shape is right. It has no status, no gate, and no
runtime record. Drafting a design is not approving one.

A bounded change does not need this. Maintainer work runs directly on the
current branch from the user's request, as `AGENTS.md` sets out; a design with
five near-empty files is worse than no design at all.

## Where it goes

Author under `sources/`, which stays passive material in this checkout exactly
as it does in a generated workspace:

```
sources/system-design/<product-or-project>/<grouping>/<scope>/
```

- `<product-or-project>` names the product **as a whole — never a single
  repository**. Where a design concerns only one of this repository's two
  products, that product is the level; where the whole system is in view, this
  level may collapse to `sources/system-design/<grouping>/<scope>/`.
- `<grouping>` is the organizing dimension for a body of design work. **Default:
  a version** — 3-number semver with a `v` prefix (`v2.0.0`, not `v2.0`).
  A **named grouping** — `phase-2`, a milestone slug — is a deliberate choice
  where a body of work does not line up with a release. Pick one dimension
  within a product rather than mixing release and phase folders.
- `<scope>` is a **concern**.

Where the request does not name a grouping, derive it from the version the
design targets: the **next** version for forthcoming design, which is the normal
case, or the **current** version when documenting as-built state. Read it from
this repository's own source of truth — `VERSION` for the workspace template,
`CLI_VERSION` for the CLI. The two version lines are separate and each drives
its own publication, so a design spanning both says which one it targets rather
than implying a single number. Never invent a version; where the target is
genuinely ambiguous, surface the choice rather than guessing.

Writing these files is an ordinary authoring write. Afterwards `sources/` is
passive again: read a file in it only when a request names that file.

## The three-tier layout

Every folder carries a `README.md` index; each scope's normative content is
`design.md`; detail splits into files or sub-folders **as it grows**.

```text
sources/system-design/<product-or-project>/
  README.md                 # index of groupings (landing + reading order)
  <grouping>/
    README.md               # grouping index: the scopes and the reading order
    <scope>/
      README.md             # scope index
      design.md             # overview — NORMATIVE, readable end to end
      <concern>.md          # one concern each
      <concern>/            # a concern with parts: README + a file per part
```

## How much detail per file

Calibrate by what a reader needs from the file, not by filling it:

- **`README.md` — orientation only.** What exists here and in what order to read
  it. No design content, decisions, or mechanism.
- **`design.md` — the normative overview a reviewer can stop at.** The
  capability, the problem, the principles, the fixed decisions, and the shape of
  the whole — readable start to finish without opening a single detail file. It
  states *what* and *why*, and defers *how* to the detail files it links.
- **`<concern>.md` — one concern at implementation depth.** Its mechanism,
  interfaces, edge cases, and the diagram that explains it. Name it for the
  concern (`publication.md`), never `details.md`. It continues from `design.md`
  and never restates it.

Every fact lives once. If a reviewer must open five files to grasp the design,
too much left `design.md`; if `design.md` runs many screens of mechanism, too
much stayed.

## Split only when a concern earns it

- A one-paragraph concern is a **paragraph in `design.md`**, not a file.
- A concern that outgrows a section becomes its **own `<concern>.md`**.
- A concern with distinct parts becomes a **sub-folder** — a `README.md` index
  and one file per part.

Start minimal, with a `README.md` and a `design.md`, and grow from there. Do not
pre-fragment.

## Separate scopes by concern, never by repository

A system design earns its keep by describing how the pieces fit across the
system, so scopes cut by **concern** — the command surface, release assembly,
the workspace contract, a cross-cutting flow — matching the product rather than
the git layout.

- **Never one scope folder per repository.** That fragments exactly the
  coherence the design exists to capture.
- Where one area needs deep design detail, that is a **scope** within the
  product's design, not a separate top-level design and not a new
  `<product-or-project>`.

## Diagrams

Embed diagrams as fenced ` ```mermaid ` blocks in the file they explain. That is
the durable form: it renders on GitHub and in most editors, and it survives in
the file rather than in a tool. You may preview through a host plugin where one
exists, but never depend on one — with no plugin, still emit the fenced mermaid.

## How it reaches the rest of the checkout

A design is source material, so it feeds this repository the way all source
material does, with nothing new added.

**Knowledge.** What the design settles durably becomes a note in `context/`,
written through `.agents/skills/cc-source-develop/SKILL.md`, with its catalog
entry in `context/INDEX.md` kept consistent in the same change. The note carries
the concept in its own words and anchors to code; it never cites the design
file, because `sources/` is passive while the knowledge outlasts it. This skill
never writes knowledge itself.

**Records.** There are none here. This checkout holds no `workspace.yaml`, no
intents, and no plans, and `AGENTS.md` bars the product's lifecycle from it, so
a design is read by a person and by `context/` and by nothing else. The product
skill a generated workspace receives lets an intent cite a design path; that
does not apply here, because there are no intents to cite one.

## Relationship to the shipped skill

`product/skills/cc-system-design/SKILL.md` is the same rubric written for a
generated workspace, and it ships. This file is the source checkout's own, and
the two differ only where a workspace and a maintainer checkout genuinely
differ: the version source, how knowledge is written, and the absence of
records. Keep the layout, altitude, splitting, and scope-separation rules
identical in both when either changes — they are one rubric, and a reader
comparing them should find no daylight.

## Boundaries

- **Draft structured source files only.** Never plan, execute, commit, publish,
  or write knowledge. This skill grants no route, role, or authority.
- **No status, no gate, no runtime record.** A system design is a source.
- **One home.** Design material lives under `sources/system-design/` and
  nowhere else.
- **`sources/system-design/` here is the dogfooded example** of the layout and
  altitude above. It is still passive: read a file in it only when a request
  names that file.
