---
kind: domain
status: accepted
title: Host adapters
slug: host-adapters
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: e268297
    basis: current-wrapper
generated_at: 2026-09-08T00:00:00Z
review_date: 2026-12-08
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
  - .context-circuit/wrapper/adapters/AGENTS.md
---

# Host adapters

## Summary

Codex CLI, Claude Code, and Cursor Agent CLI enter the same Context Circuit
workflow as thin transports. Each host also discovers that workflow through
its committed native project folders — `.claude/`, `.codex/`, `.cursor/` —
which route to Context Circuit owners rather than copying policy. Host
identity, version, capability, permission mode, and provider status are
bounded, provider-neutral evidence only. Route host-discovery, native-folder
routing, child-delegation, permission-mode, and missing-child questions here.

## Scope

Inside: provider-neutral `host_evidence`, the shared `AGENTS.md` instruction
surface, the `CLAUDE.md` Claude Code adapter, the `CURSOR.md` Cursor Agent
adapter, committed host-native folders (`.claude/`, `.codex/`, `.cursor/`) as
routes into owners, native-child mapping to the worker, verifier, or planner
packet, and `host-blocked` missing-child behavior.

Outside: embedding a host CLI or SDK, storing credentials/transcripts/auth
state, treating permission mode or child features as authorization, and adding a
second router or lifecycle.

## Behavior

All three hosts enter through workspace-root instruction adapters **and**
committed native project folders. `.agents/` stays at the workspace root.
Product adapter text lives under `.context-circuit/wrapper/adapters/` and is
what release assembly copies onto a new workspace. That copy must not
overwrite this maintainer source checkout: root `AGENTS.md` / `WORKFLOW.md`
stay source-specific, and root `CLAUDE.md` / `CURSOR.md` import workspace-root
`AGENTS.md` (`@AGENTS.md`) so they do not pull the shipped product
`AGENTS.md` from the nested adapters. In an instantiated workspace, the root
files *are* those adapter copies.

Host-native folders are the committed integration surface each host already
searches. They are thin routes: they name the owner and the host-only
frontmatter or TOML the host requires. They do not copy role bodies or invent
a second authorization policy.

| Host | Instruction surface | Native project tree | Native child mapping |
| --- | --- | --- | --- |
| Codex CLI | `AGENTS.md`, then `.agents/skills/cc-*` | `.codex/agents/*.toml` (no `.codex/rules/`) | `spawn_agent` → worker, verifier, or planner packet |
| Claude Code | root `CLAUDE.md` (`@AGENTS.md`) | `.claude/agents/`, `.claude/rules/`, `.claude/skills/cc-*/SKILL.md` routes to `.agents/skills/cc-*` | Task/subagent → same packet |
| Cursor Agent CLI | root `CURSOR.md` (`@AGENTS.md`) | `.cursor/agents/`, `.cursor/rules/`; skills via `.agents/skills/` | Task/subagent if available; otherwise host-blocked |

The product host set is worker, verifier, and planner stubs in each host's
format, plus standing-rule stubs on Claude and Cursor (role-tiering spawn,
commit convention; Cursor also keeps the GitHub unsandboxed `gh` rule).
Claude skill discovery uses thin `SKILL.md` routes under `.claude/skills/`; Codex and
Cursor already scan `.agents/skills/`. Maintainer-only extras
(`.claude/agents/cc-human-simulator.md`, `.claude/skills/cc-test-case/`) stay
in this source checkout and are not the shipped set.

Host identity, version, capability, permission mode, and provider status are
bounded provider-neutral evidence; they never authorize a route, role, lease,
gate, verification, or completion (INV-HOST-01).

A native child maps only to the bounded worker, independent read-only
verifier, or planner packet. If a required child is unavailable, the route
stays read-only and reports `host-blocked`; the host must never self-verify or
downgrade a verifier into a worker.

Across hosts, the coordinator's resumable session may have a session or thread
id, but that is only the root conversation transport. Required worker, verifier,
and planner roles must be native child agents attached to that root, using the
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

Each host obtains the per-role `(model, effort)` config by invoking the
workspace runtime, never by reading, finding, or grepping the config file
directly: a print-only runtime verb resolves the workspace-root local
override when present, else a committed fallback, and reports which source it
used. This exists because a host's own file-search tooling can be rewritten by
that host's local configuration to silently exclude gitignored paths, so a
failed search is not evidence the config is absent. The verb only prints; it
never selects a host group or applies a default itself (INV-RUNTIME-01), and
the resulting `(model, effort)` stays bounded host evidence that authorizes
nothing (INV-HOST-01).

## Workflows

- Shared host contract and safety spine: `.context-circuit/wrapper/adapters/AGENTS.md`
- Claude Code adapter: `.context-circuit/wrapper/adapters/CLAUDE.md`
- Cursor Agent adapter: `.context-circuit/wrapper/adapters/CURSOR.md`

## Interfaces

- Shared instructions: root `AGENTS.md` / `WORKFLOW.md`
- Claude adapter: root `CLAUDE.md`; shipped copy `.context-circuit/wrapper/adapters/CLAUDE.md`
- Cursor adapter: root `CURSOR.md`; shipped copy `.context-circuit/wrapper/adapters/CURSOR.md`
- Host-native routes: `.claude/`, `.codex/`, `.cursor/` (product stubs and links)
- Coordinator role: `.context-circuit/agents/coordinator.md`
- `host_evidence` shape owned by `.context-circuit/wrapper/contracts/schemas/`
- Per-role config access rule and fallback behavior: `.context-circuit/docs/role-tiering.md`

## Constraints and edge cases

Optional live host probes are explicitly pass, unavailable, or blocked; never a
false success. Host-native agent stubs stay pointers: `.context-circuit/agents`
remain the role owners. Standing host rules cite the owning invariant and
restate the clauses that must be true at commit or spawn time — they do not
invent a second policy. The commit-convention rule restates INV-COMMIT-01's
no-attribution clause (including stripping a host-injected trailer) because a
worker or host git wrapper still commits without opening `invariants.yaml`.
Cursor spawn rules point at workspace-root `CURSOR.md`, not a second copy of
the spawn rule, and must not import the nested product adapter in this source
checkout. Codex has no documented project rules tree — do not create
`.codex/rules/`; Codex standing instructions stay in this adapter and agent
TOML. Personal host state (`~/.claude/`, `~/.codex/`, `~/.cursor/`,
`settings.local.json`, transcripts, credentials) is never workspace state.

## Adapter duties and limits

The adapter's positive duties: pass the execution brief, ask the runtime to
record commits and evidence, coordinate repair attempts, and return child
results. It must provide actual read-only capability for the verifier child or
report blocked.

The adapter must not bypass a missing approval, must not change the plan to done
after a verification pass, and must not self-verify when the verifier child is
unavailable.

## Implementation references

- `.context-circuit/wrapper/adapters/AGENTS.md`, `.context-circuit/wrapper/adapters/CLAUDE.md`,
  `.context-circuit/wrapper/adapters/CURSOR.md`, `.context-circuit/wrapper/adapters/WORKFLOW.md`,
  `.context-circuit/wrapper/adapters/README.md`
- `.claude/`, `.codex/`, `.cursor/` (committed host-native routes)
- `.context-circuit/agents/coordinator.md` (routing owner; there is no `routes.yaml`)
- `.agents/skills/cc-pair/SKILL.md` (the worker-child mapping for direct collaboration)
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-HOST-01, INV-PAIR-01

## Acceptance notes

Accepted 2026-08-24. Corrected against the shipped wrapper: the invariant family
consolidated to a single `INV-HOST-01`; the offline-fallback rule is no longer a
separate invariant id; `.context-circuit/wrapper/contracts/routes.yaml`,
`.context-circuit/docs/host-capabilities.md`, and `.context-circuit/docs/agent-workspace-workflow.md` no longer
exist. The seven shipped skills are `cc-workspace`, `cc-plan`, `cc-execute`,
`cc-verify`, `cc-complete`, `cc-archive`, `cc-deliver`.

Extended 2026-09-03: the worker packet also serves direct collaboration
([cc-pair](../direct-collaboration/README.md), the Explore tier), which launches
no verifier and fails closed to `host-blocked` with no coordinator write
fallback. INV-HOST-01 is unchanged.

Extended 2026-09-08: committed `.claude/`, `.codex/`, and `.cursor/` trees are
the native integration surface. They route to existing owners; they are not
optional host-local convenience and are not a second policy copy.

Extended 2026-09-09: the commit-convention host rule restates INV-COMMIT-01's
no-attribution clause (inspect and strip a host-injected trailer) instead of
only pointing at the invariant, so a worker or host git wrapper still sees it
at commit time.

Extended 2026-09-10: a host obtains the per-role config only by invoking the
runtime, never by reading the config file directly, because a host's own
file-search tooling can be configured to silently exclude gitignored paths.
