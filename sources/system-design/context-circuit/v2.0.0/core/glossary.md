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

**Local binding** — the machine-specific half of identity: checkout paths, the
branch this machine starts work from, and the active member. Never shared.

**Workspace repository** (`workspace_repository`) — the Git repository carrying
the workspace itself. Described separately from `repositories`, because every
consumer of that map reads an entry as somewhere work happens.

**Knowledge repository** — a repository of knowledge this workspace reads and
never owns, mounted read-only and retrieved beside its own catalog. Records no
base branch, takes no worktree, is never named by a plan, and is never
reconciled here. Also *borrowed knowledge*.

**Language** — the language a member's intents and plans are written in, recorded
in `members.yaml`. Unset means English. Knowledge, headings, and field names stay
English regardless.

**Tone** — one line of prose recording the register that language is written in.
Quoted into the instruction that writes a record, never parsed.

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

**Completion** — an explicit human-requested act that appends a note, stamps the
`completed_at` instant, and returns catalog candidates for reconciliation. Not
delivery, and not automatic.

**Instant** — a canonical ISO 8601 UTC timestamp. The form of `created_at`,
`approved_at`, and `completed_at`. Every other date is an ISO 8601 calendar date.

**Direct change** — a change a person asks for without an intent, a plan, or a
worktree, made in a bound checkout. Their choice to make; it relaxes no
authorization and skips no knowledge reconciliation.

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

**`Owner:` block** — the one place a note's anchors live, at its end. An anchor
written into a sentence is a `check` finding.

**Actor note** — one note per person, team, or external system the project deals
with, recording what that actor needs from it today. What an intent is grounded
against.

**Reviewed date** — the date a catalog entry's note was last confirmed against
the code. Read back by `check` against the commits under that note's anchors.

**Readability finding** — a `check` finding about the shape of a note rather than
its links: an anchor in a sentence, a paragraph hiding a list, a catalog heading
with no line under it, a diagram fence that will not render.

**Resolution** — the `resolve` every `check` finding carries: the command, the
edit, or `needs a person` where no command discharges it.

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
to `inherit`. `role-tiering.local.yaml` overrides one host/role pair on this
machine.

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

These appear in v1 and pre-v1 material, which this repository keeps in its
history rather than its tree. None of them exists in v2; see
[retired-machinery.md](./retired-machinery.md).

consequence tier · Explore / Standard / Critical · `contract_digest` · candidate ·
acceptance record · execution record · verification record · host evidence ·
path lease · repair loop · promotion · reconciliation debt · `cc-publish` ·
publication record · *mandatory* member band · invariant (`INV-*`) ·
contract schema ·
wrapper · runtime engine · host adapter · plan status `draft` ·
`record note` and the `## Update` section (removed in 2.0.0-rc.6) ·
the `completed` date field and the bare approval date (superseded by instants in
2.0.0-rc.4)
