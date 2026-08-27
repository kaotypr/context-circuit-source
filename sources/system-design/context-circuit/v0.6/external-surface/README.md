# External-surface — design

The v0.6 scope for a **manually-triggered, config-driven way to publish Context
Circuit data to external systems** (task trackers, chat, docs) — fully orthogonal
to the core workflow. A configured pipeline is a **publication**; publications live
under `publication/`. v0.6 is a **delta on v0.5** — see [../README.md](../README.md)
for the version index and [../../v0.5/](../../v0.5/) for the base design.

## Reading order

1. [design.md](./design.md) — the overview: the capability, the orthogonality
   principle, the backbone (the `publication/` folder, the trigger, the isolation
   contract), boundaries, fixed decisions, and implementation order.
2. [configuration-and-records.md](./configuration-and-records.md) — the
   `publication/<name>/` folders, `config.yaml`, create-on-first-use, and the
   per-plan records under `published/` (with cross-plan lookup produced on demand).
3. [publish-plan.md](./publish-plan.md) — the `plan` kind: mapping a plan and its
   tasks onto ClickUp / Jira / GitHub / Notion work items, self-contained,
   one-way, and idempotent, with a worked trace.
4. [thread.md](./thread.md) — the `thread` kind: a plan's open questions as a chat
   discussion (Slack first) — a `[thread]`-prefixed parent plus one fully-described
   reply per question, discussion-safe on re-run.
5. [usage.md](./usage.md) — how a template/wrapper user configures and publishes, in
   plain language (the prompts they use); ships in the getting-started docs.
6. [contracts.md](./contracts.md) — the proposed invariants and owner-map
   additions, and why this scope needs **no core contract bump**.

## Authority

This scope adds no owner and duplicates no rule; it specifies intent and points to
the canonical owners under `wrapper/`. The contract additions this capability will
require are summarized in [contracts.md](./contracts.md).
