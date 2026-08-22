---
kind: domain
status: proposed
title: Host adapters
slug: host-adapters
owners: []
sources:
  - plans/context-circuit-plans/multi-host-agent-support/plan.yaml
  - plans/context-circuit-plans/multi-host-agent-support/PLAN.md
source_revisions:
  - plan: multi-host-agent-support
    status: done
    updated_at: 2026-08-22T13:40:00Z
generated_at: 2026-08-23T00:00:00Z
review_date: 2026-09-22
freshness: proposed-from-done-plans
assumptions:
  - The three named hosts share one router and one lifecycle.
unknowns:
  - Live host smoke evidence was not re-read for this refresh.
contradictions: []
acceptance:
  state: pending
  accepted_at:
  accepted_by:
workflows:
  - docs/host-capabilities.md
  - docs/agent-workspace-workflow.md
---

# Host adapters

## Summary

Codex CLI, Claude Code, and Cursor Agent CLI enter the same Context Circuit
workflow. Host identity and capability are evidence only. Route host-discovery,
child-delegation, resume, permission-mode, or missing-child questions here.

## Scope

Inside: provider-neutral `host_evidence`, shared `AGENTS.md` discovery, the
Claude `CLAUDE.md` adapter, child mapping to writer or verifier packets,
offline fallback, and host-blocked missing-child behavior.

Outside: embedding a host CLI or SDK, storing credentials or transcripts,
treating permission mode or child features as authorization, and adding a
second router or lifecycle.

## Behavior

All three hosts use the two-stage router in `wrapper/runtime/engine.sh`. Host
fields never authorize a route, role, lease, gate, or verification result
(INV-HOST-01).

| Host | Instruction surface | Child mapping |
| --- | --- | --- |
| Codex CLI | `AGENTS.md`, then `.agents/skills/cc-*` | Native subagent → writer or verifier packet |
| Claude Code | `CLAUDE.md` imports `AGENTS.md` | Task/subagent → same packet |
| Cursor Agent CLI | Root `AGENTS.md`; `CLAUDE.md` is also readable | Task/subagent if available; otherwise host-blocked |

A missing required child emits `host-blocked` / `block-missing-child-primitive`.
The host must not self-verify, downgrade a verifier into a writer, or skip a
gate (INV-HOST-02). Provider status `disabled`, `denied`, or `unavailable`
uses the filesystem-only fallback and does not block planning, execution
evidence, or recovery (INV-OFFLINE-01).

`host_evidence` records only bounded host, version, capability, role,
permission-mode, and provider-status fields. It never stores credentials,
provider payloads, transcripts, or auth state (INV-HOST-03).

The shipped discovery list stays at seven skills: `cc-entry`, `cc-next`,
`cc-plan`, `cc-execute`, `cc-verify`, `cc-gates`, and `cc-upgrade`.

## Workflows

- Host capability contract: `docs/host-capabilities.md`
- Shared workflow: `docs/agent-workspace-workflow.md`

## Interfaces

- Shared instructions: root `AGENTS.md` / `WORKFLOW.md`
- Claude adapter: `CLAUDE.md` → `wrapper/adapters/CLAUDE.md`
- Host evidence fields owned by wrapper session, delegation, and handoff schemas

## Constraints and edge cases

Optional live host probes are explicitly pass, unavailable, or blocked; they
are never a false success. Cursor does not require `.cursor/rules` for this
behavior. A future scoped Cursor rule must remain a thin adapter.

## Implementation references

- `wrapper/adapters/AGENTS.md`, `wrapper/adapters/CLAUDE.md`, `wrapper/adapters/WORKFLOW.md`
- `wrapper/contracts/routes.yaml` host_binding
- `wrapper/contracts/invariants.yaml` INV-HOST-01 through INV-HOST-03

## Verification

Done-plan verification IDs MHS-VT-01–MHS-VT-10.

## Provenance

Read only the selected done plan `multi-host-agent-support` (CC-002). Raw
`sources/` was not scanned. HEAD at generation was `b7a11f3`.

## Acceptance notes

This page is proposed. Accepted architecture in `context/ARCHITECTURE.md`
already states that skills are thin adapters; this page adds the shipped
three-host mapping. Human context acceptance is still required.
