# Decisions

## Accepted starter schema

Workspace identity and project identity are separate Product Knowledge pages.

- `context/WORKSPACE.md` describes this workspace.
- `context/PROJECT.md` describes the project being built, or records that it
  is uninitialized.
- `context/PRODUCT-DIRECTION.md` is not part of the starter schema.

This is a template file-schema decision, not a product-direction dump for a
specific project.

## 2026-08-16 — Delivery policy IDs are path names

Live `delivery.policy` values are `remote-review`, `local-target`, and
`manual`. These names describe the path to prepare, not workspace identity
and not an authorized git action.

When reading `workspace.yaml`, `team-review` is an alias for
`remote-review` and `solo-local` is an alias for `local-target`. An old ID
is not missing configuration and does not fall back to `manual`. New writes
use the new IDs.

`workspace.yaml` `mode` remains identity and stays independent of delivery
policy. Configuration does not grant commit, push, PR, or merge.

No project-specific product decisions have been accepted yet.

Record accepted and superseded project decisions here after human review.
Keep this file as this project's decision log; do not treat other template or
framework history as this project's product decisions.
