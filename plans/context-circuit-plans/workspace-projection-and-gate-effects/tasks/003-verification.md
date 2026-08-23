---
schema_version: 2
id: WPE-003
plan: workspace-projection-and-gate-effects
status: done
repository: context-circuit-source
paths:
  - wrapper/migrations/README.md
  - wrapper/migrations/upgrade.sh
  - docs/configuration.md
  - docs/getting-started.md
  - test/contracts/test-contracts.sh
  - test/routing/test-router.sh
  - test/upgrades/test-upgrades.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [WPE-002]
acceptance: [WPE-AC-02, WPE-AC-05, WPE-AC-06]
verification: [WPE-VT-01, WPE-VT-03, WPE-VT-05, WPE-VT-06, WPE-VT-07]
expected_evidence:
  - Upgrade and rollback fixtures that insert or restore the identity region without rewriting Product Knowledge or accepted identity.
  - Documentation matching identity-region and effect ownership, including displayed defaults.
  - Full semantic acceptance and release-boundary results.
stop_conditions:
  - Migration rewrites accepted identity or authored Product Knowledge without a human gate.
  - Release assembly includes maintainer plans, local bindings, or runtime state.
  - Any existing gate, repository, host, or lifecycle suite regresses.
---

# Integrate migration, documentation, and semantic verification

## Objective

Prove identity-region and gate-effect behavior across existing workspaces,
release assembly, and the full semantic contract.

## Work

Add migration and rollback that insert a missing identity region into legacy
`WORKSPACE.md`, `PROJECT.md`, and `INDEX.md` without rewriting authored
Product Knowledge and without changing accepted identity. Document that
`workspace.yaml` is the identifier, `context/` is Product Knowledge, omitted
gate fields display as proposed defaults, and effect identifiers are
descriptive. Extend contract, routing, upgrade, release, and acceptance
fixtures. Verify repository registration updates shared identity and every
identity region without exposing local bindings.

## Non-goals

Do not migrate unrelated mutable context, publish an artifact, activate
create-empty, or repair failures outside this plan's declared paths.

## Verification

Use WPE-VT-01, WPE-VT-03, WPE-VT-05, WPE-VT-06, and WPE-VT-07.

## Expected evidence

Migration round trips, release artifact inspection, complete semantic results,
and an independent handoff mapping each acceptance criterion.

## Stop conditions

Stop on lossy migration, identity mutation without confirmation, Product
Knowledge rewrite, release leakage, or a regression requiring scope beyond
the declared plan.
