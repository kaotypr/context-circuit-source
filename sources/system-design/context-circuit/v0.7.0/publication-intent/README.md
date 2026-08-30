# v0.7.0 — Publication field metadata (index)

Scope landing for **publication field metadata and the consult-before-publish
session** — a delta on the v0.6 external surface
([../../v0.6/external-surface/](../../v0.6/external-surface/)). Read that scope's
[design.md](../../v0.6/external-surface/design.md) first; this scope changes only
what it names.

Reading order:

- [design.md](./design.md) — the normative overview: the gap (published records
  store identity, not the field values that were pushed), the three-layer model
  (intent / last-published snapshot / external current), the fixed decisions, and
  the shape of the solution. Stop here for the whole design.
- [intent-and-record.md](./intent-and-record.md) — the user-owned `intent/`
  field layer, the `estimate_minutes` canonical unit and its `"2h 30m"` human I/O
  format, and the last-published `fields:` snapshot added to the publication
  record.
- [consult-and-preview.md](./consult-and-preview.md) — the manual
  consult-before-publish preview mode, the plain-language diff it renders, and the
  one contract delta: a display-only external drift read as a bounded
  clarification of INV-EXTERNAL-02.
