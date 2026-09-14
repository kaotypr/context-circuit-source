# Glossary

v2 vocabulary, and the v1 terms that no longer exist.

## The workspace

**Workspace** — the durable knowledge and coordination layer for one product
across one or more repositories. Not a synonym for a repository.

**Member** — a person sharing the workspace. Appears as `created_by` and nowhere
else. The roster is `members.yaml`; the active member on this machine is the
gitignored `member.local.yaml`.

**Logical repository ID** — the shared name for a repository (`api`, `web`),
recorded in `workspace.yaml` with its default base branch and its relationships.
Bound to a machine-local path in the gitignored `repositories.local.yaml`.

**Relationship** — a recorded, described link between two logical repositories
("web consumes the API from api").

**Local binding** — the machine-specific half of identity: checkout paths and the
active member. Never shared.

## Records

**Intent** (`iNNN-slug.md`) — one human decision about what "correct" means for a
change: goal, non-goals, constraints, observable success criteria, rough
repository scope, and the recorded approval. Written **before** detailed code
investigation.

**Plan** (`pNNNN-slug.md`) — how an approved intent maps to real code, plus what
happened. May span several repositories. Written **after** approval, from the
code.

**Band** — a member's optional allocation block index. Band N takes intents from
`N*100` and plans from `N*1000`. Unbanded members take what no band claims. A
band prevents offline collisions; it is not a namespace and not ownership.

**Reservation** — a permanently allocated ID in `.context-circuit/ids.yaml`.
Never reused, including after archival or deletion.

**Completion** — an explicit human-requested act that appends a note, stamps a
`completed` date, and returns catalog candidates for reconciliation. Not
delivery, and not automatic.

## Knowledge

**Note** — one durable unit of project knowledge in `context/`. Describes the
project, never the machinery that produced it.

**Catalog** (`context/INDEX.md`) — the retrieval index. One note is one unwrapped
entry carrying its link, repositories in braces, the question it answers, search
terms, and a reviewed date.

**Entry** — one catalog line. Must be unwrapped, because retrieval matches whole
lines.

**Glossary** (`context/glossary.md`) — the project's vocabulary mapped to the code
identifiers implementing it.

**Repository anchor** — a durable code reference written with the logical
repository ID, exact or patterned: `api@internal/billing/dunning/`.

**Reconciliation** — judging which durable knowledge a completed plan changed,
and editing the note and its catalog entry in one pass.

**Source** (`sources/`) — passive raw evidence. Read only when a request names the
exact file; never scanned, and never cited from a note.

## Working

**Worktree** — an isolated Git working copy. Default branch
`cc/<plan-id>/<repository-id>`, default path
`.worktrees/<plan-id>/<repository-id>`.

**Reuse** — cloning ignored `node_modules` and `.env` entries from the bound
checkout into a new worktree, by filesystem CoW where available.

**Copy mode** — `auto` (CoW, copy fallback), `required` (CoW or fail), `copy`
(force copies), `off`.

**Dependency inputs** — the tracked manifests, lockfiles, workspace definitions,
and Node version files compared before dependencies are reused.

**Wave** — one layer of plans that may run concurrently, derived from recorded
dependencies and completion.

**Linear chain** — every plan stacked on the previous one; strictly serial, no
integration merge.

**Integration merge** — the local merge that assembles a dependent plan's base at
a fan-in. Implementation, not delivery.

**Start reference** — the base (and any merges) a plan begins from, per
repository, as reported by `record order`.

**Conflict surface** — a coarse signal of how much undelivered work shares a
repository. Repositories, not paths.

## Agents

**Coordinator** — the main session. Owns interpretation, task selection,
dispatch, integration, and reporting.

**Explorer / planner / worker / reviewer** — the four role shapes. Capability
descriptions, never authority levels.

**Role tiering** — the per-host `(model, effort)` pair for each role, defaulting
to `inherit`.

**Dispatch specification** — what `agent dispatch` returns, carrying
`launch_required: true`. Not evidence that anything ran.

**Sole ownership / shared ownership** — the two worker briefs. One worker per
plan owns its worktree; `--shared` is for several workers inside one.

**Independent review** — a manually requested, read-only inspection by a fresh
context. Never automatic, never blocking.

## Distribution

**Template version** (`VERSION`, `.context-circuit/VERSION`) — the workspace
template's version, released as `v*`.

**CLI version** (`CLI_VERSION`, `.context-circuit/CLI_VERSION`) — the CLI version
a workspace pins, released as `cli-v*`.

**Version store** — the directory holding side-by-side CLI installs at
`<store>/<version>-<os>-<arch>/`.

**Seed** — the blank workspace embedded in the CLI, used by `init` and
`template export`.

## Retired — terms with no v2 meaning

These appear in v1 material and in `sources/CORE_CONCEPT-v0.6.md`,
`-v0.7.md`, and `-v1.md`. None of them exists in v2; see
[retired-machinery.md](./retired-machinery.md).

consequence tier · Explore / Standard / Critical · `contract_digest` · candidate ·
acceptance record · execution record · verification record · host evidence ·
path lease · repair loop · promotion · reconciliation debt · `cc-publish` ·
publication record · *mandatory* member band · invariant (`INV-*`) ·
contract schema ·
wrapper · runtime engine · host adapter · plan status `draft`
