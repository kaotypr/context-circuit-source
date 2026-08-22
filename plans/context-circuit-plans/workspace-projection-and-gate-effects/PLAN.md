# Atomic workspace projections and explicit gate effects

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

This plan isolates the unresolved projection and confirmation-integrity work
from the broader source PRD. It makes shared workspace summaries mechanically
consistent with `workspace.yaml` and makes gate effects complete before a
human confirms them.

## What approval authorizes

Approval authorizes bounded changes to the canonical workspace schema, route
and invariant owners, projection engine operations, gate presentation,
mutable seed, migration, documentation, and semantic fixtures named by this
plan.

Approval does not start execution, claim a lease, create a worktree, change
Git, deliver, publish, deploy, merge, or clean runtime.

## Scope and non-goals

`plan.yaml` owns the exact scope. This plan does not add repository creation,
replace repository bootstrap, change lifecycle status, generate runtime child
records, or add a second effect-policy owner.

## Proposed solution

Define a canonical projection from `workspace.yaml` into the three required
human-facing summaries and provide one atomic render/validate operation.
Define gate effects beside the canonical route and gate owners, separating
what the current confirmation changes from what a later action may change.
Cards must display every confirmed value or reject the request as incomplete.
Delegations must remain a subset of the approved effect set.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| WPE-001 | Projection and effect contracts have one owner and stable failures | — |
| WPE-002 | Engine and cards render, validate, and present those contracts atomically | WPE-001 |
| WPE-003 | Migration, release, documentation, and semantic suites prove the behavior | WPE-002 |

The completed `repository-bootstrap` plan is a prerequisite because repository
registration already owns shared identity and host-local binding boundaries.

## Acceptance criteria

- WPE-AC-01: One operation renders or validates every required summary.
- WPE-AC-02: Projection disagreement blocks entry or write preflight.
- WPE-AC-03: Cards separate immediate effects from later authorized effects.
- WPE-AC-04: Missing fields and undeclared delegation effects fail closed.
- WPE-AC-05: Repository registration updates shared projections atomically.
- WPE-AC-06: Existing gated flows remain valid.

## Verification

WPE-VT-01 through WPE-VT-07 cover contracts, gates, routing, recovery,
upgrades, release assembly, and full semantic acceptance. Executable commands
remain canonical in `plan.yaml`.

## Risks, assumptions, and open decisions

- Projection rendering must preserve mutable user-authored sections or define
  the three summaries as wholly generated; execution must resolve that format
  before writing a migration.
- Effect vocabulary belongs with existing route and invariant owners. A new
  parallel policy file is not authorized unless contract ownership is updated
  explicitly.
- The dependent create-empty plan may consume a reserved effect identifier,
  but this plan does not activate repository creation.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Must
requirements 2 through 4 and the related acceptance criteria. Current
repository evidence includes CC-001 repository ownership, CC-004 bounded
approval behavior, the workspace schema, route/invariant owners, and the
required mutable context seed. The source PRD remains passive and is not
Product Knowledge.
