# Engine-owned runtime records and child packets

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

The engine already owns leases, binding evidence, atomic writes, and record
validation, while coordinators still handcraft sessions, receipts,
delegations, and completion evidence. This plan closes that narrower gap
without adding a scheduler or host-specific runtime.

## What approval authorizes

Approval authorizes atomic runtime constructors, ownership-graph validation,
minimal launch projections, a handoff skeleton, lifecycle/resume integration,
compatibility handling, and bounded host, recovery, security, release, and
semantic fixtures.

Approval does not start execution, claim a real plan lease, create a real
worktree, launch a provider, change Git, deliver, publish, deploy, merge, or
clean runtime.

## Scope and non-goals

`plan.yaml` owns exact scope. Existing schemas and ownership invariants remain
authoritative. Validation tokens prove current graph validation only; they
never authorize execution or replace a human gate.

## Proposed solution

Add engine operations that build each record from validated inputs, compute
real metadata, stage and validate output, then publish atomically. Validate the
entire ancestry, lease, repository, worktree, host evidence, receipt,
delegation, and handoff graph before returning a short-lived launch result.
Project only locators and role/root information into the host launch request.
Include a bounded handoff skeleton so children do not search other sessions.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| ERR-001 | Constructors and ownership relationships are schema-defined | — |
| ERR-002 | Engine generates and validates complete runtime graphs | ERR-001 |
| ERR-003 | Lifecycle, resume, and hosts consume constructors consistently | ERR-002 |
| ERR-004 | Interruption, compatibility, security, release, and acceptance pass | ERR-003 |

The completed repository and host plans plus draft route-context packet plan
are prerequisites.

## Acceptance criteria

- ERR-AC-01: Engine operations generate complete records with real metadata.
- ERR-AC-02: Partial records never become authoritative.
- ERR-AC-03: Child start requires one validated ownership graph.
- ERR-AC-04: Launch projection cannot disagree with or expand the packet.
- ERR-AC-05: Children use packet and handoff skeleton without repository search.
- ERR-AC-06: Resume and completion reject stale or placeholder evidence.
- ERR-AC-07: Legacy evidence remains explicitly recoverable and never ships.

## Verification

ERR-VT-01 through ERR-VT-09 cover contracts, runtime, ownership, recovery,
hosts, security, upgrades, release assembly, and full semantic acceptance.

## Risks, assumptions, and open decisions

- Atomic publication across multiple files requires a defined staged
  transaction boundary and recovery marker; sequential renames alone must not
  expose a partially valid ownership graph.
- Host launch remains outside the engine. The engine validates and returns
  bounded evidence; adapters transport it.
- Existing legacy records may remain readable but must not gain new ownership
  or authorization merely through migration.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Should
requirements 11 and 12 and related acceptance criteria. Existing lease and
binding constructors are treated as implemented groundwork. The source PRD
remains passive and is not Product Knowledge.
