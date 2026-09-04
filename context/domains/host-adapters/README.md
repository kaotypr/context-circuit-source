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
surface, the `CLAUDE.md` Claude Code adapter, native-child mapping to the worker
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
| Codex CLI | `AGENTS.md`, then `.agents/skills/cc-*` | subagent → worker or verifier packet |
| Claude Code | `wrapper/adapters/CLAUDE.md` imports `AGENTS.md` | Task/subagent → same packet |
| Cursor Agent CLI | root `AGENTS.md` (`CLAUDE.md` also readable) | Task/subagent if available; otherwise host-blocked |

A native child maps only to the bounded worker or independent read-only
verifier packet. If a required child is unavailable, the route stays read-only
and reports `host-blocked`; the host must never self-verify or downgrade a
verifier into a worker.

Across hosts, the coordinator's resumable session may have a session or thread
id, but that is only the root conversation transport. Required worker, verifier,
and tracer roles must be native child agents attached to that root, using the
host's supported child primitive. A separate top-level task, peer thread, or
resumed root session is not a child and does not satisfy the role requirement;
if native child creation is unavailable, the route is `host-blocked`.

The worker packet also carries [direct collaboration](../direct-collaboration/README.md)
(`cc-pair`, the Explore tier): a native child maps to the same worker role for a
live pairing turn just as for plan execution. Pairing never launches a verifier
(the Explore tier has none, INV-ASSURE-01), so the only child it needs is the
worker. If the host cannot create that worker child (or its isolated worktree),
the outcome is `host-blocked` and read-only — there is no coordinator write
fallback, because the coordinator never performs the worker's edits itself
(INV-PAIR-01). This stays within INV-HOST-01: child capability is bounded
evidence and a missing required capability fails closed.

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
- `.agents/skills/cc-pair/SKILL.md` (the worker-child mapping for direct collaboration)
- `wrapper/contracts/invariants.yaml`: INV-HOST-01, INV-PAIR-01

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

Extended 2026-09-03 from proposal `0030-change-host-adapters-for-pairing`: the
worker packet also serves direct collaboration ([cc-pair](../direct-collaboration/README.md),
the Explore tier), which launches no verifier and fails closed to `host-blocked`
with no coordinator write fallback. INV-HOST-01 is unchanged.
