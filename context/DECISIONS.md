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

## 2026-08-17 — Additive stack execution is runtime, not a plan type

`cc-run-plan` remains the sole standard single-plan execution skill and still
creates worktrees from the repository default or active branch.
`cc-run-stack` executes a connected set of already-approved plans in one root
session. The plan tree lives in `.runtime/stacks/<stack-id>/graph.yaml` and
the resume cursor in `progress.yaml`. There is no scheduler, no durable
`stack.yaml` under `plans/`, and no `plans/<repository-key>-stacks/` layout.
Implemented is runtime evidence; canonical plan status stays `draft`,
`approved`, and `done`. `cc-finish-plan` remains the per-plan human gate.

No project-specific product decisions have been accepted yet.

Record accepted and superseded project decisions here after human review.
Keep this file as this project's decision log; do not treat other template or
framework history as this project's product decisions.
