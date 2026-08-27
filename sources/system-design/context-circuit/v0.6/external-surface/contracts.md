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
| `publication_record` | `wrapper/contracts/schemas/publication-record.yaml` (new — describes `published/<plan-id>.yaml`) |
| publication-kind trigger | INV-SKILL-01 (existing) — kinds ship as `cc-<kind>` at `.agents/skills/cc-<kind>/SKILL.md` (first: `cc-publish-plan`) |

The config schema owns the `config.yaml` shape; the record schema owns the per-plan
record shape; the trigger is **not** a new authority — it is an ordinary skill under
the existing INV-SKILL-01, surfaced as a slash command by the host adapter as an
optional convenience.

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
v0.5-shaped workspace plus the availability of the adapter skills.

## The "publish" wording guardrail

`publish-plan` shares the word "publish" with the v0.5 delivery boundary, where
"publication" means **git** publication (INV-DELIVER-01, and the enumerated
`…/push/publish/deploy/…` list in INV-RUNTIME-01). The senses are separated by
object — a *plan* is published to a tracker; a *branch / pull request* is published
in git — but to keep the record unambiguous, every core-delivery mention stays
qualified as **git publication**. This is a wording discipline on the
delivery/runtime owners, not a new rule of its own.
