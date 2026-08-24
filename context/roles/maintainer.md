---
kind: role
status: accepted
title: Product-source maintainer
slug: maintainer
role_type: business
owners: []
domains:
  - context/domains/repository-binding/README.md
  - context/domains/host-adapters/README.md
  - context/domains/plan-review/README.md
  - context/domains/plan-approval/README.md
  - context/domains/plan-execution/README.md
  - context/domains/verification/README.md
  - context/domains/completion/README.md
  - context/domains/plan-organization/README.md
  - context/domains/delivery/README.md
related_workflows:
  - docs/release.md
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - This checkout remains identity.kind product-source, not an instantiated workspace.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
---

# Product-source maintainer

## Role definition

The person maintaining the Context Circuit source repository. This checkout is
the product, not a customer workspace. The role ships wrapper contracts,
adapters, skills, docs, and tests, and keeps `template/` as the blank seed.

## Primary outcomes

- Keep one coordinator, one lifecycle, and one owner per rule.
- Land bounded changes without turning this repo into an instantiated workspace.
- Preserve dirty unrelated work, local bindings, and user data across release
  and upgrade.

## Product surfaces

- Maintainer plans under `plans/context-circuit-plans/` (the namespace exists
  per `plans/README.md`; it is currently empty).
- Wrapper contracts and `wrapper/runtime/engine.sh`.
- Shipped adapters and the seven `cc-*` skills: `cc-workspace`, `cc-plan`,
  `cc-execute`, `cc-verify`, `cc-complete`, `cc-archive`, `cc-deliver`.
- Release assembly under `scripts/`, described in `docs/release.md`.

## Cross-domain perspective

Repository binding keeps registered product repos portable; this source checkout
is already identified. Host adapters are how Codex, Claude Code, and Cursor enter
the same lifecycle as evidence-only transports. Planning and review are
non-mutating; approval is an explicit conversational gate; execution, independent
verification, completion, archive/restore, and delivery are each separate
explicit actions.

## Related domains and workflows

- [Workspace orientation and repository binding](../domains/repository-binding/README.md)
- [Host adapters](../domains/host-adapters/README.md)
- [Planning and plan review](../domains/plan-review/README.md)
- [Plan approval](../domains/plan-approval/README.md)
- [Plan execution](../domains/plan-execution/README.md)
- [Verification](../domains/verification/README.md)
- [Completion](../domains/completion/README.md)
- [Plan organization](../domains/plan-organization/README.md)
- [Delivery](../domains/delivery/README.md)
- Release assembly: `docs/release.md`

## Role-specific behavior and limitations

An explicit user request may authorize a source-only commit. Registered product
repositories still need their own delivery and publication gates. Do not scan
`sources/` unless the request names exact files. Do not treat this checkout as a
user project or restore an obsolete lifecycle.

## Provenance

Re-grounded on the current wrapper at HEAD `4b8ac0b`. The five previous-version
maintainer plans that first seeded this page were deleted in `4b8ac0b`; their
provenance was retired. Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-08-24. Agent execution roles remain under `agents/` (coordinator,
writer, verifier) and are not replaced by this page.
