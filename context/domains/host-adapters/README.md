---
kind: domain
status: accepted
title: Host adapters
slug: host-adapters
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - The three named hosts share one router and one lifecycle.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - wrapper/adapters/AGENTS.md
---

# Host adapters

## Summary

Codex CLI, Claude Code, and Cursor Agent CLI enter the same Context Circuit
workflow as thin transports. Host identity, version, capability, permission
mode, and provider status are bounded, provider-neutral evidence only. Route
host-discovery, child-delegation, permission-mode, and missing-child questions
here.

## Scope

Inside: provider-neutral `host_evidence`, the shared `AGENTS.md` instruction
surface, the `CLAUDE.md` Claude Code adapter, native-child mapping to the writer
or verifier packet, and `host-blocked` missing-child behavior.

Outside: embedding a host CLI or SDK, storing credentials/transcripts/auth
state, treating permission mode or child features as authorization, and adding a
second router or lifecycle.

## Behavior

All three hosts enter the same coordinator, lifecycle, and runtime library.
Host identity, version, capability, permission mode, and provider status are
bounded provider-neutral evidence; they never authorize a route, role, lease,
gate, verification, or completion (INV-HOST-01).

| Host | Instruction surface | Native child mapping |
| --- | --- | --- |
| Codex CLI | `AGENTS.md`, then `.agents/skills/cc-*` | subagent → writer or verifier packet |
| Claude Code | `wrapper/adapters/CLAUDE.md` imports `AGENTS.md` | Task/subagent → same packet |
| Cursor Agent CLI | root `AGENTS.md` (`CLAUDE.md` also readable) | Task/subagent if available; otherwise host-blocked |

A native child maps only to the bounded writer or independent read-only
verifier packet. If a required child is unavailable, the route stays read-only
and reports `host-blocked`; the host must never self-verify or downgrade a
verifier into a writer.

`host_evidence` records only bounded host, version, capability, role,
permission-mode, and provider-status fields. It never stores credentials,
provider payloads, transcripts, or auth state.

## Workflows

- Shared host contract and safety spine: `wrapper/adapters/AGENTS.md`
- Claude Code adapter: `wrapper/adapters/CLAUDE.md`

## Interfaces

- Shared instructions: root `AGENTS.md` / `WORKFLOW.md`
- Claude adapter: `wrapper/adapters/CLAUDE.md`
- Coordinator role: `agents/coordinator.md`
- `host_evidence` shape owned by `wrapper/contracts/schemas/`

## Constraints and edge cases

Optional live host probes are explicitly pass, unavailable, or blocked; never a
false success. Cursor does not require `.cursor/rules` for this behavior; a
future scoped Cursor rule must remain a thin adapter.

## Adapter duties and limits

The adapter's positive duties: pass the execution brief, ask the runtime to
record commits and evidence, coordinate repair attempts, and return child
results. It must provide actual read-only capability for the verifier child or
report blocked.

The adapter must not bypass a missing approval, must not change the plan to done
after a verification pass, and must not self-verify when the verifier child is
unavailable.

## Implementation references

- `wrapper/adapters/AGENTS.md`, `wrapper/adapters/CLAUDE.md`,
  `wrapper/adapters/WORKFLOW.md`, `wrapper/adapters/README.md`
- `agents/coordinator.md` (routing owner; there is no `routes.yaml`)
- `wrapper/contracts/invariants.yaml`: INV-HOST-01

## Provenance

Re-grounded on the current wrapper at HEAD `4b8ac0b`. The previous-version
`multi-host-agent-support` plan that seeded this page was deleted in `4b8ac0b`;
its provenance was retired. Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-08-24. Corrected against the shipped wrapper: the invariant family
consolidated to a single `INV-HOST-01`; the offline-fallback rule is no longer a
separate invariant id; `wrapper/contracts/routes.yaml`,
`docs/host-capabilities.md`, and `docs/agent-workspace-workflow.md` no longer
exist. The seven shipped skills are `cc-workspace`, `cc-plan`, `cc-execute`,
`cc-verify`, `cc-complete`, `cc-archive`, `cc-deliver`.
