# Context Circuit v1.0.0

Source design for Context Circuit v1.0.0 — a **proposed major evolution** on
today's v0.7 line. Read the preceding designs ([../v0.7.0/](../v0.7.0/),
[../v0.6/](../v0.6/), [../v0.5/](../v0.5/)) first; this grouping is an argument
to test, not accepted Product Knowledge. This README is the version index; each
scope owns its own design.

## Scopes

- [core/](./core/) — the intent-gated, candidate-proven trust-core study:
  what "correct" means (intent), what ships (delivery), and the mechanical
  middle. Start at [core/design.md](./core/design.md).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. A **scope** is a bounded
area of the design. The version folder holds one scope folder per area it
covers; this version README indexes them and carries the reading order.

Every folder has a **`README.md`** as its landing and index. Each scope's
normative design is **`design.md`**; the scope's `README.md` states purpose and
reading order.

## Source boundary

These documents are maintainer design material. They have no status, no
authority, and change nothing on their own. They are not copied into an
instantiated workspace as Product Knowledge and are not normal agent context.
