---
schema_version: 2
id: LBT-001
plan: local-binding-troubleshoot
status: done
repository: context-circuit-source
paths:
  - docs/getting-started.md
  - test/contracts/test-contracts.sh
depends_on: []
acceptance: [LBT-AC-01, LBT-AC-02, LBT-AC-03]
verification: [LBT-VT-01, LBT-VT-02, LBT-VT-03]
expected_evidence:
  - Short troubleshooting paragraph in docs/getting-started.md for a missing repositories.local.yaml.
  - Optional contains assertion that the guidance remains present.
  - No credentials, machine-specific shared identity, or sibling maintainer-plan edits.
stop_conditions:
  - The change would create repositories.local.yaml, scan for checkouts, or invent a path.
  - The change would modify wrapper/runtime/engine.sh or bootstrap confirmation behavior.
  - The change would edit plans/context-circuit-plans/repository-bootstrap/ or plans/context-circuit-plans/multi-host-agent-support/.
---

# Add the missing-binding troubleshooting paragraph

## Objective

Tell a human the next safe action when `repositories.local.yaml` is missing,
in one short paragraph on the getting-started page.

## Work

Insert the troubleshooting paragraph into the existing Bind a repository
section of `docs/getting-started.md`. Keep surrounding examples. Optionally
add one `contains` assertion in `test/contracts/test-contracts.sh` so the
guidance cannot silently disappear.

Intended meaning, not a second authority over `plan.yaml`: missing binding
file is expected on a fresh clone; report and stop; create an explicit
host-local path for a logical key already in `workspace.yaml`, or request
`Bootstrap repository <key>` if the checkout does not exist.

## Non-goals

Do not change binding resolution, clone/bootstrap gates, `workspace.yaml`
identity, `.gitignore`, or the two existing maintainer plans. Do not write
`repositories.local.yaml` and do not store credentials.

## Verification

Reference verification IDs; do not copy canonical commands.

## Expected evidence

Changed `docs/getting-started.md`, optional contract assertion, passing
LBT-VT-01 through LBT-VT-03, and a handoff that records the docs-only
limitation.

## Stop conditions

Scope expansion into runtime or bootstrap behavior, contradiction with
INV-REPO-01..05, ownership conflict, or missing troubleshooting language
after the edit.
