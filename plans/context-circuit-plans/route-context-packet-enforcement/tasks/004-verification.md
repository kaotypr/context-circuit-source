---
schema_version: 2
id: RCP-004
plan: route-context-packet-enforcement
status: done
repository: context-circuit-source
paths:
  - docs/agent-workspace-workflow.md
  - docs/host-capabilities.md
  - wrapper/migrations/README.md
  - wrapper/migrations/upgrade.sh
  - test/routing/test-router.sh
  - test/context-budget/test-budgets.sh
  - test/hosts/test-host-adapters.sh
  - test/security/test-boundaries.sh
  - test/upgrades/test-upgrades.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [RCP-002, RCP-003]
acceptance: [RCP-AC-05, RCP-AC-06, RCP-AC-07]
verification: [RCP-VT-02, RCP-VT-03, RCP-VT-04, RCP-VT-05, RCP-VT-06, RCP-VT-07, RCP-VT-08]
expected_evidence:
  - Fixtures fail redundant Tier 0, broad discovery, sibling reads, and undeclared sources.
  - Upgrade and release install canonical packets and preserve compatible receipts.
  - Full semantic acceptance passes across every supported host.
stop_conditions:
  - Live-overread checks rely only on a static byte ledger.
  - Upgrade silently rewrites incompatible receipt evidence.
  - Existing route, gate, host, lifecycle, or release behavior regresses.
---

# Fail live overread and complete semantic verification

## Objective

Prove that selected packets, not agent search behavior, determine live context
across hosts, upgrades, and released workspaces.

## Work

Add semantic fixtures for repeated Tier 0 loads, broad wrapper/docs discovery,
sibling workspaces, undeclared sources, and measured overruns. Update
documentation and migration behavior, then run route, budget, host, security,
upgrade, release, and full acceptance suites.

## Non-goals

Do not claim enforcement that a host cannot provide, require live providers in
offline CI, or repair unrelated semantic failures.

## Verification

Use RCP-VT-02 through RCP-VT-08.

## Expected evidence

Failing overread cases, passing bounded cases, cross-host outcomes, migration
round trips, release inspection, and independent acceptance mapping.

## Stop conditions

Stop if overread can still pass behind static totals, a host limitation is
hidden, or a regression requires scope beyond this plan.
