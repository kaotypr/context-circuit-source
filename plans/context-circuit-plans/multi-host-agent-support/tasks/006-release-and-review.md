---
schema_version: 2
id: MHS-006
plan: multi-host-agent-support
status: ready
repository: context-circuit-source
paths:
  - scripts/release-artifact.sh
  - scripts/release-manifest.txt
  - wrapper/manifest.yaml
  - wrapper/migrations/README.md
  - wrapper/migrations/upgrade.sh
  - docs/host-capabilities.md
  - docs/getting-started.md
  - docs/configuration.md
  - docs/operations.md
  - docs/release.md
  - test/upgrades/test-upgrades.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [MHS-005]
acceptance: [MHS-AC-01, MHS-AC-02, MHS-AC-04, MHS-AC-05, MHS-AC-06, MHS-AC-07, MHS-AC-08]
verification: [MHS-VT-01, MHS-VT-06, MHS-VT-08, MHS-VT-09, MHS-VT-10]
expected_evidence:
  - Updated release and upgrade inventory with the Claude adapter included and local host state excluded.
  - Onboarding, configuration, operations, and troubleshooting documentation for all three hosts.
  - Independent verifier handoff and complete semantic acceptance output.
stop_conditions:
  - Artifact inventory contains maintainer plans, runtime records, tests, provider payloads, or credentials.
  - Upgrade or rollback can overwrite user workspace, local host configuration, or runtime evidence.
  - The verifier cannot reproduce acceptance evidence independently.
---

# Integrate release, migration, documentation, and independent review

## Objective

Ship the host adapters consistently, document how developers use each host,
and close the plan with independent evidence rather than host-specific claims.

## Work

Update the release manifest and assembler for the shipped Claude bridge and
any host contract inventory. Confirm wrapper upgrades and rollback replace
only wrapper-owned files while preserving workspace, plans, local host
configuration, registered repositories, and runtime evidence. Add concise
onboarding, configuration, operations, troubleshooting, and compatibility
guidance that distinguishes host capability from authorization and keeps the
filesystem workflow usable offline.

Run the full acceptance suite after the verifier inspects the changed files,
task-to-acceptance mapping, release inventory, migration boundaries, and
optional live evidence labels. The verifier writes only its independent
handoff; a human still controls finish, delivery, publication, and cleanup.

## Non-goals

Do not publish an artifact, commit automatically, install a host, or mark the
plan done from test output alone.

## Verification

Use MHS-VT-01, MHS-VT-06, MHS-VT-08, MHS-VT-09, and MHS-VT-10.

## Expected evidence

Release staging inventory, upgrade/rollback results, documentation review,
full acceptance output, independent verifier handoff, known limitations, and
the remaining human finish/delivery gates.

## Stop conditions

Stop on failed acceptance, missing host adapter files, artifact leakage,
unresolved migration risk, dirty/unpushed work, missing child primitives, or
any attempt to infer completion or publication authorization.

