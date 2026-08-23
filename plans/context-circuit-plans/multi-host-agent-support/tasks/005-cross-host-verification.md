---
schema_version: 2
id: MHS-005
plan: multi-host-agent-support
status: done
repository: context-circuit-source
paths:
  - test/hosts/test-host-adapters.sh
  - test/hosts/fixtures/
  - test/routing/fixtures.yaml
  - test/behavior-matrix/matrix.yaml
  - test/behavior-matrix/fixture-pairs.tsv
  - test/ownership/test-ownership.sh
  - test/context-budget/test-budgets.sh
  - test/security/test-boundaries.sh
  - test/acceptance.sh
depends_on: [MHS-002, MHS-003, MHS-004]
acceptance: [MHS-AC-01, MHS-AC-02, MHS-AC-03, MHS-AC-04, MHS-AC-05, MHS-AC-06, MHS-AC-08]
verification: [MHS-VT-02, MHS-VT-03, MHS-VT-04, MHS-VT-05, MHS-VT-06, MHS-VT-07]
expected_evidence:
  - One host fixture matrix covering root, writer, verifier, resume, offline, and blocked paths.
  - Static checks and complete offline suite results for all three hosts.
  - Explicitly labeled optional live results or unavailable/blocked evidence.
stop_conditions:
  - A live provider call is required to pass offline acceptance.
  - A host smoke captures credentials, provider payloads, or uncontrolled transcripts.
  - Failed verification is repaired by the verifier or represented as success.
---

# Add cross-host semantic and optional live verification

## Objective

Prove that the three host adapters produce the same observable Context Circuit
behavior across success, failure, recovery, isolation, budget, and gate cases.

## Work

Add a host test harness and fixtures that inspect instruction surfaces,
capability descriptors, delegation packets, role permissions, resume records,
and normalized route outcomes. Cover at least:

- root entry and read-only next-action recommendation;
- plan review, approval, execution, verification, and finish gates;
- writer-only worktree writes and independent read-only verifier behavior;
- missing child primitive, host unavailable, provider disabled, and offline fallback;
- interrupted session resume with receipt, wrapper, Git, and ownership checks;
- context budget ceilings and forbidden credential/provider-payload fields.

Support an explicit live mode only when the user opts in and the required host
binary is already installed and authenticated. Run against a disposable
fixture, sanitize output, preserve no host state in the repository, and report
`pass`, `unavailable`, or `host-blocked` with the executable/version evidence.

## Non-goals

Do not launch external agents from default CI, create credentials, or let live
output redefine canonical acceptance. Do not make a verifier repair a failed
implementation.

## Verification

Use MHS-VT-02, MHS-VT-03, MHS-VT-04, MHS-VT-05, MHS-VT-06, and MHS-VT-07.

## Expected evidence

Host fixture inventory, normalized outcomes, test output, context-budget
measurements, sanitized optional live summaries, and explicit unavailable or
blocked records where live access is absent.

## Stop conditions

Stop on any mismatch in canonical reason codes, missing independent verifier
evidence, credential-shaped output, path escape, hidden external mutation, or
scope expansion into provider infrastructure.

