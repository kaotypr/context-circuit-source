# Files and folder structure

The installed workspace layout under v1.0, the conventions every folder and file
follows, and — specifically — how `sources/` is organized and how it relates to the
new `intent/` object. Additive: v1.0 adds `intent/` and a few runtime records; it
renames nothing and keeps every existing convention.

## The installed workspace tree

```
<workspace root>/
  workspace.yaml                 # portable identity (credential-free)          [committed]
  repositories.local.yaml        # host-local bindings (machine paths)          [host-local]
  role-tiering.local.yaml        # optional per-role (model, effort)            [host-local]

  context/                       # durable Product Knowledge                    [committed]
    INDEX.md                     #   retrieval catalog (INV-KNOWLEDGE-01)
    PROJECT.md ARCHITECTURE.md CONVENTIONS.md DECISIONS.md TERMINOLOGY.md
    WORKSPACE.md SOURCES.md
    domains/  roles/  references/
    proposals/                   #   pending knowledge changes (human-gated)
    sources.yaml

  intent/                        # NEW — the decision for each change           [committed]
    i0007-checkout-retries/
      INTENT.md                  #   human-facing: the bigger picture
      contract.yaml              #   frozen criteria + scope + tier (M1), schema_version 2
    INDEX.md                     #   active intents catalog
    archive/                     #   archived intents (status-blind move)

  plans/                         # grounded plans (derive from an intent)       [committed]
    0012-add-retry/
      PLAN.md  plan.yaml         #   plan.yaml gains `intent: 0007-...`
      tasks/                     #   (or steps/) per-task detail
    INDEX.md
    archive/

  sources/                       # passive raw evidence + structured design      [committed]
    <request-named files>        #   read only when explicitly named (INV-SEC-02)
    system-design/               #   cc-system-design output (three-tier)
      <product>/<grouping>/<scope>/{README.md, design.md, <concern>.md}
    archive/
    README.md

  publication/                   # external export surface (orthogonal)          [never-ship]
    <name>/{config.yaml, intent/, published/}

  repositories/                  # gitignored clone destination                  [host-local]
    <bound project repos>

  .runtime/                      # private runtime state                         [never-ship]
    executions/<plan>/<exec>/
      execution.yaml  snapshot/  attempts/  repositories/  grounding/
      candidate.yaml             #   NEW — current candidate identity (M2)
      human-acceptance.yaml      #   NEW — candidate-bound acceptance (M2)
      context-impact.yaml        #   + reconciliation-debt marker (M4)
    worktrees/                   #   isolated worktrees (incl. cc-pair/)
    pairing/<session>/           #   Explore session pointers
    trace/<intent>/          #   NEW — per-repo trace manifests, recorded
                                  #     grounding evidence, freshness-checked (M1)
    knowledge-debt/              #   NEW — delivered-but-unreconciled markers (M4)
    locks/                       #   worker locks + path leases
```

The only structural addition is **`intent/`** (a committed, first-class tree) plus a
handful of runtime records. Everything else is today's layout.

## How `sources/` is organized (and how it relates to `intent/`)

`sources/` is **passive raw evidence** — the same role it has today (INV-SEC-02): only
files a request explicitly names are read; the workspace never scans all of `sources/`
or a sibling workspace to fill a gap. v1.0 does not change that. Two kinds of content
live there:

- **Ad-hoc source files** — a pasted spec, a transcript, an exported doc — dropped in
  by the human and referenced by exact name when authoring an intent or plan.
- **Structured design** under `sources/system-design/<product>/<grouping>/<scope>/`,
  authored by `cc-system-design` in the fixed three-tier layout: `README.md`
  (orientation), `design.md` (the normative overview a reviewer can stop at), and one
  `<concern>.md` per concern at implementation depth. `<grouping>` defaults to a
  `v`-prefixed semver; scopes cut by concern, never per repository.

**The relationship to `intent/`** — this is the seam the study's front door adds:

```
sources/  (raw evidence, structured design)   ── grounds ──▶   intent/  (the decision)
                                                                  │  approved + frozen
                                                                  ▼
                                                               plans/   (derivation)
```

- `sources/` is *material you might build from*; it has **no status, no gate, no
  runtime record** (unchanged).
- `intent/` is the *decision to build a specific thing* — first-class, gated,
  frozen at approval. An intent may be **grounded in** named `sources/` files (and in
  Product Knowledge), the way a plan is grounded today.
- So a `cc-system-design` document under `sources/system-design/` now has a natural
  next step it lacked before: it feeds an **intent** (the human's approved decision),
  instead of only feeding context proposals. `sources/` stays passive; `intent/`
  carries the authority.

## File and folder conventions

- **Stable ids — distinct for intents and plans.** Plans use `<NNNN>-<kebab-slug>`
  (for example `0012-add-retry`). **Intents use an `i`-prefixed form,
  `i<NNNN>-<kebab-slug>`** (for example `i0007-checkout-retries`), so an id reads as an
  intent on sight and never collides with a plan id when it appears alone (as in a
  plan's `intent:` field). Each tree keeps its own never-reused four-digit sequence
  (the next after the highest ever allocated, archived included), reusing the existing
  plan-id allocation logic. Slugs are lowercase kebab-case, no dates, no underscores.
- **`.md` + `.yaml` pairing.** Every first-class object pairs a human-facing Markdown
  file with a machine record: `INTENT.md`/`contract.yaml`, `PLAN.md`/`plan.yaml`. The
  human reads the `.md`; the engine reads the `.yaml`. Neither duplicates the other's
  authority.
- **INDEX catalogs.** `intent/INDEX.md`, `plans/INDEX.md`, and `context/INDEX.md` are
  retrieval-first catalogs (INV-KNOWLEDGE-01) — navigable by a human, but structured so
  an agent locates the right unit without scanning the whole tree.
- **Archive is a status-blind move.** `intent/archive/` and `plans/archive/` receive a
  whole `<id>/` directory unchanged; archiving performs no status or completion
  validation and implies nothing (INV-ARCHIVE-01/02). The normal agent does not
  traverse `archive/`; an explicit restore returns an item first.
- **Three-tier design docs.** `sources/system-design/` follows README → design →
  concern altitude, one concern per file, diagrams inline as fenced ` ```mermaid `.
- **This study is itself source material.** `sources/context-circuit-v1.0/` is
  exactly the "structured design under `sources/`" convention applied to the
  product's own evolution — no status, no authority.

## Ownership classes (what ships, what stays local)

| Class | Paths | Rule |
| --- | --- | --- |
| **Committed / portable** | `workspace.yaml`, `context/`, `intent/`, `plans/`, `sources/` | credential-free; travels with the workspace |
| **Host-local (gitignored)** | `repositories.local.yaml`, `role-tiering.local.yaml`, `repositories/`, `.runtime/` | machine paths, clones, runtime state; never portable |
| **Never-ship** | `publication/`, `.runtime/` | excluded from any template/release (manifest `never_ship`) |
| **Credentials** | — | never in any workspace file or runtime record (INV-SEC-01); stay in host Git config / the host agent |

`intent/` joins the **committed/portable** class: an intent is credential-free
decision material that travels with the workspace, exactly like `plans/` and
`context/`. Its runtime by-products (candidate, acceptance, debt markers) live under
the host-local, never-ship `.runtime/`.

## Maintainer-source note

In the maintainer source repository, the product layout above is what
`template/` seeds and `wrapper/` governs; the `intent/` tree, its schema
(`wrapper/contracts/schemas/intent-contract.yaml`), and the `cc-intent` skill
(`.agents/skills/cc-intent/`) are the additions delivered into that shipped layout.
The maintainer-source `sources/` (where this study lives) is a source-only surface,
never part of `context-circuit-template`, unchanged by v1.0.
