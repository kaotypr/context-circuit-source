# worker-role-naming — design

## Capability

Give the plan-execution implementing role **one name — `worker`** — everywhere it
appears, retiring `writer` as a second name for the same role. One word for one
concept, so a reader never has to learn that "writer" and "worker" mean the same
thing.

## Problem

Context Circuit's execution model has exactly one implementing role and one
independent checking role: a **worker** commits the changes, a **verifier**
checks them. But the implementing role is named **two ways** across the product,
and the two names are used interchangeably for the same thing:

- The **role file is `agents/writer.md`**, and the invariants owner map records it
  as `worker_role: agents/writer.md` — the value already calls it "worker" while
  the path still says "writer".
- **Terminology defines it circularly.** `context/TERMINOLOGY.md`:
  "Worker / verifier | The single **writer** role for an execution …";
  `docs/terminology.md`: "Worker | The single **writer** role for one execution."
  The canonical glossary needs the other term to define the term.
- **Invariants prose is split.** INV-OWN-01/02 speak of the "one-**writer** lock"
  and "at most one active **writer**"; INV-GROUND-01/02/03 say "the **writer**
  reads and honors …". The same rules elsewhere say "one **worker** executes"
  and "the **worker** commits".
- **Runtime and adapters carry `writer`.** The brief is the "**writer** brief"
  (`writer-brief.md`, `cc_writer_brief_assemble`, the `writer-brief-assemble`
  command); the schema beside it is `worker-handoff.yaml`. Two spellings for two
  halves of the same handoff.

The dominant, user-facing term is already **worker** (the schema is
`worker-handoff.yaml`; the spoken form in `docs/terminology.md` is "the worker" /
"I made the changes"). `writer` is the minority spelling and the source of the
confusion. Unify on `worker`.

## Principles

- **One name per concept.** A glossary that defines "worker" as "the writer role"
  fails its one job. The fix is to pick the dominant term and remove the other.
- **Worker is already the canonical term.** It is the term in the schema name, the
  execution invariants' common phrasing, and the conversational voice. Renaming
  *toward* worker moves the fewest concepts and breaks no established user-facing
  word.
- **Rename, do not re-scope.** This changes *what the role is called*, never what
  it does, who may act, or any gate. No behavior, authority, or contract semantics
  move.
- **The brief is named after the role.** "Writer brief" embeds the old role name;
  unifying the role means the brief, its file, its engine function, and its
  command follow — otherwise the retired word survives in the most-used runtime
  surface.

## Fixed decisions

1. **`worker` is the sole name** of the implementing execution role. `writer` is
   retired as a synonym and removed from roles, contracts, runtime, tests, and
   Product Knowledge prose.
2. **Rename the role file** `agents/writer.md` → `agents/worker.md`, and update the
   invariants owner-map value `worker_role:` to point at it. The owner-map *key*
   was already `worker_role`; only the path changes.
3. **Reword the split invariant prose** to `worker`: INV-OWN-01/02 ("one-worker
   lock", "at most one active worker") and INV-GROUND-01/02/03 ("the worker reads
   and honors …"). Rule **IDs and semantics are unchanged** — only the noun.
4. **Rename the brief to the "worker brief"**, following the role: the concept in
   terminology, the file `writer-brief.md` → `worker-brief.md`, the engine
   function `cc_writer_brief_assemble` → `cc_worker_brief_assemble`, and the
   `writer-brief-assemble` command → `worker-brief-assemble`. The brief's content,
   assembly, and preflight are unchanged; only its name is.
5. **Fix the glossaries** so each defines `worker` directly: drop "the writer role"
   phrasing in `context/TERMINOLOGY.md` and `docs/terminology.md`, and rename the
   "writer brief" glossary rows to "worker brief".
6. **No core contract bump and no new owner.** The invariants owner map is edited
   in place (one path value, one term); no rule is added, split, or moved to a new
   owner.

## Shape of the change (defers *how* to the owners)

Each surface is owned elsewhere and changed through its owner's normal action:

- **Roles** (`agents/writer.md` → `agents/worker.md`; `agents/coordinator.md`,
  `agents/verifier.md` prose that names "the writer").
- **Contracts** (`wrapper/contracts/invariants.yaml`): the `worker_role` and
  `writer_brief` owner-map values and the INV-OWN / INV-GROUND prose nouns.
- **Runtime** (`wrapper/runtime/engine.sh`): the brief function and any `writer`
  identifiers/messages; the `writer-brief-assemble` CLI verb.
- **Adapters / templates** (`wrapper/adapters/writer-brief.md` and the shipped
  copy): file rename to `worker-brief.md` and any in-file self-reference.
- **Skills** (`.agents/skills/cc-execute`, `cc-run-stack`): the assemble command
  name and any "writer" prose.
- **Release** (`scripts/release-artifact.sh`, `scripts/release-manifest.txt`): the
  staged brief path and the required-role path.
- **Product Knowledge** (`context/TERMINOLOGY.md`, `context/DECISIONS.md`,
  `context/domains/*`, `context/roles/maintainer.md`) and **docs**
  (`docs/terminology.md`): prose nouns and the glossary rows.
- **Tests / harness** (`test/**`, `template-harness/**`): fixture paths, function
  names, and asserted strings that carry `writer`.

## Constraints, edge cases, and sequencing

- **Overlap with `writer-brief-placement`.** That scope *relocates* the brief
  (`writer-brief.md` → `wrapper/runtime/`); this scope *renames* it
  (`writer-brief.md` → `worker-brief.md`). They touch the same file for different
  reasons. If both land, the result is
  **`wrapper/runtime/worker-brief.md`** (this scope owns the name, that scope owns
  the location). Land them in one change, or land placement first and let this
  rename follow the moved file; either order converges on the same path.
- **The source of truth stays under `wrapper/adapters/`** — renamed to
  `worker-brief.md` there too; ownership (INV-GROUND-03) is unchanged beyond the
  filename.
- **No retroactive history rewrite.** Completed maintainer plans and prior version
  designs that say "writer" are historical record and are not edited.
- **Grep the aliases out.** A residual "writer" anywhere in a live surface
  reintroduces the split, so the acceptance for the eventual plan is that no live
  product surface names the role "writer".

## How it feeds the rest — unchanged

This is source. It reaches Product Knowledge and a plan through the normal path:
the coordinator gathers context from this named design source, proposes the
terminology and contract wording changes through the existing context-proposal
path, a human accepts, and a plan grounds in the result. There is no separate
design-acceptance gate.
