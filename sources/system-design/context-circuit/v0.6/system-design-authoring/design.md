# Context Circuit v0.6 — System-Design Authoring (overview)

Status: authoritative source design for the v0.6 system-design-authoring scope (delta on v0.5)
Revision: 2 — 2026-08-27

This is the overview of the system-design-authoring scope: a **skill** that helps
the root session author and structure a **system design** as source material.
Read [README.md](./README.md) first, and [../../v0.5/core/design.md](../../v0.5/core/design.md)
for the base lifecycle this rides on unchanged.

## The one capability

A larger change — a new product, a new version, a cross-cutting feature — is
easier to plan when someone has first written down the **shape of the system**:
its topology, domains, cross-cutting flows, and the decisions behind them. That
write-up is a **system design**. Today an author improvises its structure (sndp
put one under `sources/system-design/sndp/v0.1/` by hand); the result is only as
coherent as the author's memory of the layout convention.

v0.6 ships a **`cc-system-design` skill** that teaches the agent to author one
**well** — the right place, the three-tier layout, the right amount of detail per
file, and clean separation of scopes by concern. That is the whole scope.

## What this is NOT

- **Not a lifecycle stage.** There is no new position between Product Knowledge
  and Plans, no acceptance gate, and no status. A system design is a **source**,
  and `sources/` is passive raw evidence (v0.5 INV-SEC-02); it carries no status
  the way a plan carries `draft/approved/done`.
- **Not runtime.** No engine function, no `.runtime` record, no new schema, no new
  invariant. The skill is the only owner.
- **Not a second Product Knowledge.** A system design is source material the agent
  reads to *gather context*; the durable accepted facts still live in `context/`.
- **Not a new plan field or grounding channel.** Plans ground in Product Knowledge
  exactly as in v0.5.

## Where it sits — a source, in the existing flow

```mermaid
flowchart LR
  SD["sources/system-design/…<br/>authored via cc-system-design"] -->|"gather context about X<br/>(reads the named design)"| PK["context/<br/>Product Knowledge"]
  PK --> PLN["plans/"] --> EX["execution"] --> CMP["completion"]
```

Nothing in this diagram is new except *how the design source is authored*. The
coordinator reads the named design source during ordinary context gathering and
proposes context units through the **existing** context-proposal path; plans then
ground in that Product Knowledge as they already do. No stage, no gate, no hook.

## Why a skill (and not just prose in a doc)

The value is entirely in getting the **structure and altitude** right —
consistently, by any author, on any host. That is exactly what a read-as-procedure
skill packet is for (INV-SKILL-01): it encodes the layout convention, the
detail-per-file rubric, and the scope-separation rule so the author does not have
to reconstruct them from memory. The maintainer repo already follows this
convention by hand; the skill makes it repeatable and dogfoods the same rules.

## What changes relative to v0.5

| Area | v0.5 | v0.6 |
| --- | --- | --- |
| Authoring a system design | ad hoc, under `sources/`, no guidance | a shipped **`cc-system-design`** skill with a layout + detail + scope rubric |
| Location | ad hoc | **stays** under `sources/system-design/<product>/<version>/<scope>/` |
| Lifecycle / status / gate | — | **none added** (a design is a source) |
| Product Knowledge / plans | sources → context → plans | **unchanged** |

Everything else in v0.5 is unchanged.

## Principles

- **A system design is a source.** It is authored, then read like any other named
  source; it has no status and no acceptance gate of its own.
- **The skill owns structure, not authority.** It shapes files; it never approves,
  accepts, plans, or executes. It grants no route or role (INV-SKILL-01).
- **Detail is scale-triggered.** Keep a concern inline until it grows; split a file
  or a sub-folder only when it earns one — the same discipline the maintainer
  design set follows.
- **Scope by concern, never by repository.** A system design is inherently
  cross-repo; its coherence dies if it is fragmented one-folder-per-git-repo.
- **Dogfood.** The skill encodes the exact convention this repo's own
  `sources/system-design/` already uses; there is one convention, not two.

## Detailed design

- [authoring-rubric.md](./authoring-rubric.md) — location, three-tier layout, the
  detail-per-file rubric, and scope separation (the center of the scope).
- [product-knowledge-relationship.md](./product-knowledge-relationship.md) — a
  design is a source that feeds Product Knowledge via the existing flow.
- [contracts.md](./contracts.md) — the shipped skill and its allowlist wiring; no
  runtime change.
- [examples.md](./examples.md) — a worked example and acceptance criteria.

## Final design decisions

- A system design is **source material** under `sources/system-design/`, not a
  lifecycle stage; it has **no status, no acceptance gate, no reconciliation**.
- The scope ships exactly **one product skill, `cc-system-design`**, which is
  **purely skill-invoked** (`/cc-system-design`); there is no WORKFLOW action.
- **No engine change, no runtime record, no new schema, no new invariant, no
  `plan.yaml` change, no first-class `design/` area.**
- The skill's core payload is the **detail-per-file** and **scope-by-concern**
  rubric; the top level keys the product/initiative, never a single repository.
- The maintainer source repo **dogfoods** the same skill and convention.
