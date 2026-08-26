# System-design stage — schema and contracts

Each change names its canonical owner; this scope specifies intent, the owner
files carry the rule. No rule is duplicated (v0.5 principle 5.7).

| Owner | Change |
| --- | --- |
| `wrapper/contracts/schemas/system-design.yaml` **(new)** | the machine-readable half of a design: stable design id, human status (`draft`/`accepted`/`delivered`), version+scope, Product Knowledge references it builds on, and the context units it expects to establish |
| `wrapper/contracts/schemas/design-acceptance.yaml` **(new)** | the acceptance record: design id, `accepted_by`, `accepted_at`, the context proposals it emitted, and the scope accepted (per-scope acceptance) |
| `wrapper/contracts/invariants.yaml` | add INV-DESIGN-01 (a design is human-accepted; authoring/review change nothing) and INV-DESIGN-02 (accepting a design emits context proposals and enables plan grounding, but never auto-accepts Product Knowledge); add owner-map entries for the design artifact and the design-acceptance record |
| `wrapper/adapters/WORKFLOW.md` | add the conversational actions: *design the system for X*, *review the design*, *accept the design*, *design status* |
| `.agents/skills/cc-system-design/SKILL.md` **(new, shipped)** | the product authoring skill for system designs — scaffolds the three-tier `design/<version>/<scope>/` layout, grounds the draft in Product Knowledge + named sources, embeds mermaid, and produces a *draft* (never accepts). Triggered by the *design the system* action or a host `/cc-system-design` slash command (INV-SKILL-01) |
| `wrapper/manifest.yaml` | no schema/version change is *owned by this scope*; the plan schema `[1, 2]` and `runtime_version 0.6.0` are set by the run-stack scope in the same release |

## Plan grounding (no new plan field)

This scope adds **no field to `plan.yaml`**. A design's accepted decisions flow
into Product Knowledge (see [product-knowledge-boundary.md](./product-knowledge-boundary.md)),
and plans ground on them through the **existing** `product_knowledge` references —
so a plan implementing a design is already grounded in it, transitively, with no
second channel. If design→plan traceability is wanted, it lives on the **design
side** (the design, or its acceptance record, lists the plans that slice it),
never as a per-plan field.

## Invariants (intent)

- **INV-DESIGN-01** — A System Design has human-controlled status
  (`draft → accepted → delivered`). Authoring and review change nothing; only an
  explicit human acceptance moves it to `accepted`.
- **INV-DESIGN-02** — Accepting a design **emits** context proposals for the
  durable facts it establishes and **enables** planning against the design once
  its Product Knowledge is accepted. It does not
  itself write Product Knowledge; those proposals are accepted separately through
  the existing context path. A design is not execution authority.

## What this scope does NOT change

- No change to worker/verifier behavior, the repair loop, or the completion gate.
- No change to how Product Knowledge is stored or accepted — the design merely
  *feeds* the existing proposal path.
- No change to `plan.yaml`: plans ground in Product Knowledge exactly as in v0.5;
  this scope adds no plan field and no plan schema change.
- No new execution authority. The stage produces grounding, not permission to run.
