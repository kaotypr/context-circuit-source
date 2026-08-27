# External-surface — contracts

Continues from [design.md](./design.md). This file summarizes the **contract
additions** this scope requires and the **canonical owners** that would hold them.
A system design specifies intent and points to the owners under `wrapper/`; it does
not itself become an owner or define a parallel copy of a rule (INV-SKILL-01, and
the one-owner rule in `AGENTS.md`).

## Proposed invariants

Three new invariants, owned by `wrapper/contracts/invariants.yaml`, capture the
orthogonality, the data boundary, and the self-contained-artifact rule so they are
enforceable rather than merely described.

- **INV-EXTERNAL-01 (isolation) — sketch.** The external surface is a set of
  opt-in, manually-triggered publications that are **orthogonal to the core
  workflow**. No core-workflow phase (plan, approve, execute, verify, deliver),
  runtime action, coordinator route, or role references, triggers, waits on, or is
  affected by any publication, and none runs except on an explicit human invocation
  of its adapter skill. A publication is a command, not a standing relationship; it
  observes the workspace as it finds it at invocation time. There is no automatic,
  scheduled, event, or workflow-driven trigger.

- **INV-EXTERNAL-02 (data boundary) — sketch.** A publication reads only the
  workspace artifact types its `config.yaml` names, and writes only its own records
  under `publication/<name>/published/`. **Nothing is written under `plans/`**, so
  no planning, execution, verification, or delivery step encounters a publication
  record and none updates external state or status. It never mutates core Context
  Circuit state or plan status (INV-PLAN-01), and it is **export-first**: data flows
  Context Circuit → outward only. External side effects go through the host / MCP
  layer, never the runtime (INV-RUNTIME-01); configuration and records are
  credential-free (INV-SEC-01). A record's lifecycle is independent of a plan's
  archive state, keyed to the stable plan id. Any future import must pass through
  the normal plan-authoring gate (INV-APPROVE-01), never around it.

- **INV-EXTERNAL-03 (self-contained external artifacts) — sketch.** Every external
  artifact a publication creates or updates (task, document, message) is
  understandable to a reader with no workspace access: it describes the plan, task,
  or document in plain language and never contains a workspace file name, a
  workspace-internal path, a workspace-internal id, or internal Context Circuit
  mechanism/vocabulary. A stable plan or task identifier may appear (normally in a
  title, as `[<plan-number>]`) for human cross-reference. This is the external
  mirror of the product's existing "never expose internal mechanism" report hygiene.

These restate, for this surface, guarantees the product already makes elsewhere;
none grants a new authority.

## Proposed owner-map additions

Added under the `owners:` map in `wrapper/contracts/invariants.yaml`:

| Concern | Owner |
| --- | --- |
| `external_surface` (isolation, data boundary, self-contained artifacts) | `wrapper/contracts/invariants.yaml` (INV-EXTERNAL-01/02/03) |
| `publication_config` | `wrapper/contracts/schemas/publication-config.yaml` (new — describes `config.yaml`) |
| `publication_record` | `wrapper/contracts/schemas/publication-record.yaml` (new — the `plan` kind's record) |
| `publication_thread_record` | `wrapper/contracts/schemas/publication-thread-record.yaml` (new — the `thread` kind's record) |
| publish trigger | INV-SKILL-01 (existing) — the `cc-publish` skill at `.agents/skills/cc-publish/SKILL.md` publishes a publication per its `kind` |

The config schema owns the `config.yaml` shape; each kind's record schema owns that
kind's record shape (record shapes are per-kind); the trigger is **not** a new
authority — `cc-publish` is an ordinary skill under the existing INV-SKILL-01,
surfaced as a slash command by the host adapter as an optional convenience.

`config.yaml` carries an optional **`language`** (default `en`) that applies to every
kind: `cc-publish` authors all external text in that language. It changes no workspace
text and no owner; it is an authoring choice recorded in `publication-config.yaml`.

## No core contract bump

Unlike run-stack and repository-grounding, this scope requires **no coordinated
contract bump**:

- `plan.yaml` is unchanged — the external item id is a record under `publication/`,
  never a plan field (INV-PLAN-01 stays intact).
- `execution.yaml` is unchanged; nothing in execution, verification, or delivery is
  touched.
- `runtime_version` need not bump — the runtime library gains no provider or network
  code (INV-RUNTIME-01); it at most reuses its existing atomic write for the record.

The scope ships a **skill + two schemas + a config convention**, comparable to how
`system-design-authoring` ships only a skill and needs no contract change. It is
fully additive and opt-in: a workspace that configures no publication is a
v0.5-shaped workspace plus the availability of the `cc-publish` skill.

## Vocabulary: "publish" is external; git delivery is "push / open a pull request"

The product reserves "publish" and "publication" for the external surface, and
git delivery speaks in its own terms — "push the branch", "open a pull request",
"deliver" — and never "publish." So the word carries one meaning everywhere, and no
per-mention qualifier is ever needed.

Implementing this scope therefore includes a **wording-only** edit to the git-side
invariants that vacates the word: INV-DELIVER-01 and INV-RUNTIME-01 (and the
`runtime.excludes` list in `wrapper/manifest.yaml`, and acceptance criterion AC-16)
change "publication"/"publish" to "push" / "open a pull request". This changes no
behavior, version, or authority — git delivery still requires the same separate
human actions; only the word changes. It is the one edit this scope makes outside
`publication/`, and it is what lets the external surface own "publish" cleanly.
