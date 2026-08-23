# Full development support for Codex, Claude Code, and Cursor CLI

Status: done  
Repository: context-circuit-source  
Source: repository evidence and current host documentation

## Review summary

This plan makes the existing host-neutral Context Circuit workflow usable and
verifiable from Codex CLI, Claude Code, and Cursor Agent CLI. It is one
connected plan bundle with separate host work packages, not three competing
implementations.

The current repository already provides a two-stage router, filesystem runtime
records, `cc-*` skill adapters, role documents, and a high-level host matrix.
The remaining work is to make the host contract explicit, give Claude a
shipped `CLAUDE.md` bridge, verify Cursor CLI behavior, map native child
facilities to the existing writer/verifier packets, and add cross-host
acceptance evidence.

The design uses `AGENTS.md` as the shared instruction surface for Codex and
Cursor CLI. Claude Code gets a thin `CLAUDE.md` adapter that imports
`AGENTS.md`; it adds only host-specific loading and delegation notes. The
router, lifecycle, ownership rules, acceptance IDs, verification IDs, and
human gates remain owned by the existing wrapper contracts.

## What approval authorizes

Approval authorizes bounded implementation of the tasks and acceptance
criteria in `plan.yaml`. Approval does not start execution, claim a lease,
create a worktree, change Git, invoke a host provider, change credentials,
deliver, publish, deploy, merge, or clean runtime.

## Scope and non-goals

The implementation covers host capability evidence, root and child session
mapping, instruction discovery, resume and permission boundaries, offline
fallback, semantic fixtures, optional live smoke checks, release assembly,
migration, and documentation.

It does not embed a host CLI or provider SDK, create a scheduler, store host
authentication, or make a host's child-agent feature a reason to weaken
independent verification. The core remains complete when a host is missing,
unauthenticated, offline, or unable to create a required child.

## Proposed solution

### One workflow, three thin host adapters

| Host | Primary project instructions | Child mapping | Safe CLI evidence |
| --- | --- | --- | --- |
| Codex CLI | `AGENTS.md`, then `.agents/skills/cc-*` | Native subagent → delegated writer or verifier packet | Interactive `codex`; optional `codex exec` read-only probe |
| Claude Code | `CLAUDE.md` imports `AGENTS.md` | Task/subagent → delegated writer or verifier packet | `claude`; optional `claude -p` read-only probe |
| Cursor Agent CLI | Root `AGENTS.md`; `CLAUDE.md` is also readable | Task/subagent if available → same packet; otherwise host-blocked | `cursor-agent`; optional `-p`/`resume` read-only probe |

The exact host command is a verification probe, not a Context Circuit command.
Each host must begin with the human-facing request “Start or resume work in
this workspace.” The host then loads the entry spine and selects the one
bounded probe from the existing router. A host may supply a native child
primitive, but it must pass the same identity, scope, permissions, acceptance,
verification, evidence, and stop-condition fields.

Host identity and capability information is evidence, not authorization. A
missing child primitive emits the existing host-blocked action; the host may
not silently self-verify, downgrade a verifier into a writer, or skip a gate.
Provider failure remains `disabled`, `denied`, or `unavailable` and does not
block filesystem-only planning, execution evidence, or recovery.

### Instruction and release boundary

The source checkout keeps its maintainer instructions in the root adapters.
The released workspace receives the wrapper-owned root adapters from
`wrapper/adapters/`. The Claude task therefore adds and stages a
`wrapper/adapters/CLAUDE.md`; release assembly must copy it to the artifact
root and test it. Cursor CLI needs no `.cursorrules` or mandatory
`.cursor/rules` file for this plan because its CLI can read the root
`AGENTS.md` and `CLAUDE.md`. A future scoped Cursor rule must remain a thin
adapter and cannot duplicate canonical policy.

### Child, resume, and verification behavior

The root coordinator still owns the request, route, preflight, lease,
worktree, delegation, cards, and evidence consolidation. The writer remains
the only implementation writer in its assigned worktree. The verifier remains
a separate read-only child. Host-specific resume commands only re-enter the
same filesystem session after receipt, wrapper, Git, and ownership checks.

Live host probes are opt-in and use an isolated disposable fixture. They may
use credentials already managed by the host, but never print, capture, or
persist credentials, provider payloads, or transcripts in the repository.
Unavailable live binaries or authentication are recorded as unavailable or
host-blocked; the offline suite remains the release gate.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| MHS-001 | One host capability contract and cross-host acceptance matrix | — |
| MHS-002 | Codex CLI entry, child, and resume adapter | MHS-001 |
| MHS-003 | Claude Code import bridge, child mapping, and released adapter | MHS-001 |
| MHS-004 | Cursor Agent CLI entry, permission, child, and resume adapter | MHS-001 |
| MHS-005 | Offline semantic matrix and guarded live host probes | MHS-002–MHS-004 |
| MHS-006 | Release, migration, documentation, and independent verification | MHS-005 |

## Acceptance criteria

- MHS-AC-01: All three hosts enter the same route and lifecycle.
- MHS-AC-02: Host capability and session evidence is provider-neutral.
- MHS-AC-03: Codex maps native subagents to bounded roles and blocks safely.
- MHS-AC-04: Claude loads shared instructions through the shipped bridge.
- MHS-AC-05: Cursor CLI works in interactive, print, and resume probes.
- MHS-AC-06: Isolation, verification, recovery, budgets, and gates are shared.
- MHS-AC-07: Release and upgrade boundaries preserve the right state.
- MHS-AC-08: Offline acceptance is complete and live evidence is honest.

## Verification

`plan.yaml` owns the canonical commands. The independent verifier must inspect
the host contract, exercise every host fixture including missing-child and
offline paths, validate writer/verifier permissions, rebuild the release
artifact, confirm that credentials and provider payloads are absent, and run
the complete semantic acceptance suite. Live CLI evidence is supplementary;
it cannot replace the offline contract suite.

## Risks, assumptions, and open decisions

- Host CLIs change faster than the filesystem contract. Live probes must record
  the executable and version, while the core suite asserts stable behavior.
- Claude Code reads `CLAUDE.md`, not `AGENTS.md`, so the import bridge is a
  release requirement. Cursor CLI reads both root files; the plan must check
  for accidental duplicate or contradictory instructions.
- Cursor print mode has broad tool access. Any mutating probe requires a
  disposable fixture and an explicit human gate; default probes are read-only.
- Native subagents may be disabled or unavailable. That is a host capability
  blocker, not permission to self-verify or change the canonical route.
- Host-local permission and authentication files remain outside committed
  workspace configuration. No plan task may add example secrets or provider
  payload fixtures.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Local evidence: `docs/host-capabilities.md`,
`docs/agent-workspace-workflow.md`, `docs/runtime-contract.md`,
`wrapper/contracts/invariants.yaml`, `wrapper/contracts/routes.yaml`,
`wrapper/contracts/context-sets.yaml`, `wrapper/manifest.yaml`, the current
release assembler and semantic suites, and the existing host adapters.

Current host references used to shape the adapter boundary:

- [Codex CLI project instructions and subagents](https://learn.chatgpt.com/docs/agent-configuration/agents-md)
- [Codex CLI subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents)
- [Claude Code memory and `AGENTS.md` bridge](https://code.claude.com/docs/en/memory)
- [Claude Code CLI reference](https://docs.anthropic.com/en/docs/claude-code/cli-usage)
- [Cursor CLI usage and project instructions](https://docs.cursor.com/en/cli/using)
- [Cursor CLI parameters](https://docs.cursor.com/en/cli/reference/parameters)

