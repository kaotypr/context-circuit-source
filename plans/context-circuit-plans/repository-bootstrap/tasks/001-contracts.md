---
schema_version: 2
id: RB-001
plan: repository-bootstrap
status: draft
repository: context-circuit-source
paths:
  - workspace.yaml
  - template/workspace.yaml
  - repositories.local.yaml
  - .gitignore
  - wrapper/contracts/schemas/workspace.yaml
  - wrapper/contracts/invariants.yaml
  - wrapper/manifest.yaml
depends_on: []
acceptance: [RB-AC-01, RB-AC-02]
verification: [RB-VT-01, RB-VT-06]
---

# Define portable repository and local-binding contracts

## Objective

Define one shared repository identity contract and one root-level host-local
binding contract without storing credentials or machine-specific paths in
committed workspace state.

## Work

Add canonical URL and default-branch semantics, local binding fields, explicit
ownership rules, ignored-path rules, and release-boundary assertions. Document
that `repositories/` is optional and that the root maintainer `plans/` tree is
source-only.

## Non-goals

Do not clone, fetch, modify repositories, or change lifecycle status.

## Verification

Use RB-VT-01 and RB-VT-06.

## Expected evidence

Updated schemas, invariants, manifests, templates, ignore rules, and contract
tests with no credential-shaped values.

## Stop conditions

Stop on conflicting ownership, ambiguous path authority, or any request to
persist credentials.
