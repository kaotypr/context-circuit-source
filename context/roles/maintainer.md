---
kind: role
status: proposed
title: Product-source maintainer
slug: maintainer
role_type: business
owners: []
domains:
  - context/domains/repository-binding/README.md
  - context/domains/host-adapters/README.md
  - context/domains/plan-review/README.md
  - context/domains/plan-approval/README.md
related_workflows:
  - docs/gates.md
  - docs/release.md
sources:
  - plans/context-circuit-plans/repository-bootstrap/plan.yaml
  - plans/context-circuit-plans/multi-host-agent-support/plan.yaml
  - plans/context-circuit-plans/local-binding-troubleshoot/plan.yaml
  - plans/context-circuit-plans/approval-gate-ux-performance/plan.yaml
  - plans/context-circuit-plans/interactive-plan-review/plan.yaml
source_revisions:
  - workspace_identity: product-source
  - HEAD: b7a11f3
generated_at: 2026-08-23T00:00:00Z
review_date: 2026-09-22
freshness: proposed-from-done-plans
assumptions:
  - This checkout remains identity.kind product-source, not an instantiated workspace.
unknowns: []
contradictions: []
acceptance:
  state: pending
  accepted_at:
  accepted_by:
---

# Product-source maintainer

## Role definition

The person maintaining the Context Circuit source repository. This checkout
is the product, not a customer workspace. The role ships wrapper contracts,
adapters, skills, docs, and tests, and keeps `template/` as the blank seed.

## Primary outcomes

- Keep one router, one lifecycle, and one owner per rule.
- Land bounded plans without turning this repo into an instantiated workspace.
- Preserve dirty unrelated work, local bindings, and user data across release
  and upgrade.

## Product surfaces

- Maintainer plans under `plans/context-circuit-plans/`
- Wrapper contracts and `wrapper/runtime/engine.sh`
- Shipped adapters and the seven `cc-*` skills
- Human gates in `docs/gates.md`, including the product-source commit card

## Cross-domain perspective

Repository binding is how registered product repos stay portable; this source
checkout itself is already identified. Host adapters are how Codex, Claude
Code, and Cursor enter the same workflow. Plan review is read-only. Plan
approval on product-source is status-only, then a separate maintainer commit
of the exact approval projection, then a later `Run approved plan`.

## Related domains and workflows

- [Repository binding and bootstrap](../domains/repository-binding/README.md)
- [Host adapters](../domains/host-adapters/README.md)
- [Named-plan review](../domains/plan-review/README.md)
- [Plan approval and product-source commit](../domains/plan-approval/README.md)

## Role-specific behavior and limitations

An explicit user request may authorize a source-only commit. Registered
product repositories still need their own delivery and publication gates.
Do not scan `sources/` unless the request names exact files. Do not treat
this checkout as a user project or restore an obsolete lifecycle.

## Provenance

Read the five done maintainer plans listed in `sources`. Raw `sources/` was
not scanned. HEAD at generation was `b7a11f3`.

## Acceptance notes

This role is proposed. Agent execution roles remain under `agents/` and are
not replaced by this page. Human context acceptance is still required.
