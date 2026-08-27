# External-surface — contracts

Continues from [design.md](./design.md). This file summarizes the **contract
additions** this scope requires and the **canonical owners** that would hold them.
Per the authoring convention, a system design specifies intent and points to the
owners under `wrapper/`; it does not itself become an owner or define a parallel
copy of a rule (INV-SKILL-01, and the one-owner rule in `AGENTS.md`).

## Proposed invariants

Two new invariants, owned by `wrapper/contracts/invariants.yaml`, capture the
orthogonality and the data boundary so they are enforceable rather than merely
described.

- **INV-EXTERNAL-01 (isolation) — sketch.** The external surface is an opt-in,
  manually-triggered set of targets that are **orthogonal to the core workflow**. No
  core-workflow phase (plan, approve, execute, verify, deliver) references, triggers,
  waits on, or is affected by any target, and no target is invoked except by an
  explicit human action. A target is a command, not a standing relationship; it
  observes the workspace as it finds it at invocation time. There is no automatic,
  scheduled, or workflow-driven trigger.

- **INV-EXTERNAL-02 (data boundary) — sketch.** A target reads only the workspace
  artifact types its declaration names, and writes only its own records — the
  authoritative per-plan mapping inside the plan directory, and derived
  catalog/sync state under `external-targets/`. It never mutates core Context
  Circuit state or plan status (INV-PLAN-01), and it is **export-first**: data flows
  Context Circuit → outward only. Its external side effects go through the host /
  MCP layer; configuration and records are credential-free (INV-SEC-01). Any future
  import must pass through the normal plan-authoring gate (INV-APPROVE-01), never
  around it.

Both restate, for this surface, guarantees the product already makes elsewhere;
neither grants a new authority.

## Proposed owner-map additions

Added under the `owners:` map in `wrapper/contracts/invariants.yaml`:

| Concern | Owner |
| --- | --- |
| `external_surface` (isolation + data boundary) | `wrapper/contracts/invariants.yaml` (INV-EXTERNAL-01/02) |
| `external_target_config` | `wrapper/contracts/schemas/external-target.yaml` (new) |
| `external_target_mapping` | `wrapper/contracts/schemas/external-target-mapping.yaml` (new) |
| target-kind trigger | INV-SKILL-01 (existing) — kinds ship at `.agents/skills/<kind>/SKILL.md` |

The config schema owns the `external-target.yaml` shape; the mapping schema owns the
per-plan record shape; the trigger is **not** a new authority — it is an ordinary
skill under the existing INV-SKILL-01, surfaced as a slash command by the host
adapter as an optional convenience.

## No core contract bump

Unlike run-stack and repository-grounding, this scope requires **no coordinated
contract bump**:

- `plan.yaml` is unchanged — the external item id is a sibling record, never a plan
  field (INV-PLAN-01 stays intact).
- `execution.yaml` is unchanged; nothing in execution, verification, or delivery is
  touched.
- `runtime_version` need not bump for this scope — the runtime library gains no
  provider or network code (INV-RUNTIME-01); it at most reuses its existing atomic
  write for the mapping record.

The scope ships a **skill + two schemas + a config convention**, comparable to how
`system-design-authoring` ships only a skill and needs no contract change. It is
fully additive and opt-in: a workspace that configures no target is a v0.5-shaped
workspace plus the availability of the adapter skills.

## The "publish" wording guardrail

`publish-plan` shares the word "publish" with the v0.5 delivery boundary, where
"publication" means **git** publication (INV-DELIVER-01, and the enumerated
`…/push/publish/deploy/…` list in INV-RUNTIME-01). The senses are separated by
object — a *plan* is published to a tracker; a *branch / pull request* is published
in git — but to keep the record unambiguous, every core-delivery mention must stay
qualified as **git publication**, so bare "publish" never sits next to "plan" with a
delivery meaning. This is a wording discipline on the delivery/runtime owners, not a
new rule of its own.
