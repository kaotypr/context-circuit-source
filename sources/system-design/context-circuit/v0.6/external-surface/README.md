# External-surface — design

The v0.6 scope for a **manually-triggered, config-driven way to push Context
Circuit data to external systems** (task trackers, chat, docs) — fully orthogonal
to the core workflow. v0.6 is a **delta on v0.5** — see [../README.md](../README.md)
for the version index and [../../v0.5/](../../v0.5/) for the base design.

## Reading order

1. [design.md](./design.md) — the overview: the capability, the orthogonality
   principle, the backbone (config surface, trigger, isolation), boundaries,
   fixed decisions, and implementation order.
2. [configuration-and-records.md](./configuration-and-records.md) — the
   user-owned `external-targets/` config folder, the `external-target.yaml`
   shape, create-on-first-use, and where the records live (authoritative per-plan
   mapping vs derived catalog).
3. [publish-plan.md](./publish-plan.md) — the first consumer: mapping a plan and
   its tasks onto ClickUp / Jira / GitHub / Notion work items, one-way and
   idempotent, with a worked trace.
4. [contracts.md](./contracts.md) — the proposed invariants and owner-map
   additions, and why this scope needs **no core contract bump**.

## Authority

This scope adds no owner and duplicates no rule; it specifies intent and points to
the canonical owners under `wrapper/`. The contract additions this capability will
require are summarized in [contracts.md](./contracts.md).
