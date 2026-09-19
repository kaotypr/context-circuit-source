---
name: cc-system-design
description: Author a system design as source material under sources/system-design/ — the three-tier layout, how much detail belongs in each file, and why scopes separate by concern rather than by repository. Drafts structured source files only; never approves, plans, executes, or writes knowledge.
---

# Design a system before building it

Use this on a request to work out the shape of something larger than a single
change — "design the system for X", "write up the architecture for this
feature", "structure the design before we build it". It is for the developer who
wants the whole picture settled, and legible, before any of it is implemented.

What it produces is **source material**: passive, durable, written to be read by
a person deciding whether the shape is right. It has no status, no acceptance
gate, and no runtime record. Drafting a design is not approving one.

A bounded change does not need this. Go from the project's knowledge straight to
an intent and a plan; a design with five near-empty files is worse than no
design at all.

## Where it goes

Author under `sources/`, which the workspace treats as passive evidence:

```
sources/system-design/<product-or-project>/<grouping>/<scope>/
```

- `<product-or-project>` names the product **as a whole — never a single
  repository**. In a single-product workspace this level may collapse to
  `sources/system-design/<grouping>/<scope>/`.
- `<grouping>` is the organizing dimension for a body of design work. **Default:
  a version** — 3-number semver with a `v` prefix (`v2.0.0`, not `v2.0`).
  A **named grouping** — `phase-2`, `Q1`, a milestone slug — is a deliberate
  choice for a team that does not organize design by release. Within one
  product, pick one dimension rather than mixing release and phase folders.
- `<scope>` is a **concern**.

Where the author does not name a grouping, derive it from the version the design
targets: the project's **next** version for forthcoming design, which is the
normal case, or its **current** version when documenting as-built state. Read
that version from the project's own source of truth — whatever its repositories
actually use to version themselves. Never take it from `workspace.yaml`, which
records the Context Circuit schema and not this project's releases, and never
invent one. Where the targeted version is genuinely ambiguous, surface the
choice rather than guessing.

Writing these files is an ordinary authoring write. Afterwards `sources/` is
passive again: it is read only when a request or a record names a file in it.

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
  concern (`scheduling.md`), never `details.md`. It continues from `design.md`
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
system, so scopes cut by **concern** — a backend topology, a domain, a client
application, a cross-cutting flow — matching the product rather than the git
layout.

- **Never one scope folder per repository.** That fragments exactly the
  cross-repository coherence the design exists to capture.
- Where one repository needs deep design detail, that is a **scope** within the
  product's design — a `backend/` scope, say — not a separate top-level design
  and not a new `<product-or-project>`.

## Diagrams

Embed diagrams as fenced ` ```mermaid ` blocks in the file they explain. That is
the durable form: it renders on GitHub and in most editors, and it survives in
the file rather than in a tool. You may preview through a host plugin where one
exists, but never depend on one — with no plugin, still emit the fenced mermaid.

## How it reaches the rest of the workspace

A design is source material, so it feeds the workspace the way all source
material does, with nothing new added to the flow.

**Knowledge.** What the design settles durably becomes a note in `context/`,
written through `.agents/skills/cc-knowledge/SKILL.md`. The note carries the
concept in its own words and anchors to code; it never cites the design file,
because `sources/` is archived while knowledge outlasts it. This skill never
writes knowledge itself.

**Intents.** An intent may name a design file **when the outcome it asks for
actually depends on one** — a path under `sources/system-design/`, cited where
it bears on the goal, the constraints, or the repository scope. Only where it is
needed: most intents settle their outcome in their own words and name nothing.
An intent is still written before detailed code investigation, and a cited
design does not replace reading the real code after approval. There is no
second approval of the design, and citing it does not make it a gate.

## Boundaries

- **Draft structured source files only.** Never approve, accept, plan, execute,
  merge, deliver, or write knowledge. This skill grants no route, role, or
  authority.
- **No status, no gate, no runtime record.** A system design is a source.
- **One home.** Design material lives under `sources/system-design/` and nowhere
  else. Do not write design files under an intent's directory.
