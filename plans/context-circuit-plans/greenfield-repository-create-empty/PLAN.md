# Explicit greenfield repository creation

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

CC-001 supports existing bindings and confirmed cloning. It intentionally does
not provide a local greenfield `git init` path. This plan adds that optional
capability as a separate action after projection and effect integrity exists.

## What approval authorizes

Approval authorizes a distinct create-empty route and card, a bounded engine
operation that initializes and binds a local repository after confirmation,
and the required documentation, migration, recovery, security, release, and
semantic fixtures.

Approval does not start execution, claim a lease, create an execution
worktree, change Git, deliver, push, create a remote repository, publish,
deploy, merge, or clean runtime.

## Scope and non-goals

`plan.yaml` owns the exact boundaries. Repository bootstrap remains the clone
path. Execute remains unable to create, clone, or initialize product Git. This
plan does not require the wrapper workspace itself to be versioned.

## Proposed solution

Add `repository-create-empty` as a named action beside repository bootstrap.
The current card presents the logical key, safe destination, initial branch,
existing-path result, and effects. Exact confirmation stages an empty Git
repository, validates it, writes the ignored local binding, then commits shared
identity projections atomically through CC-006 behavior. Any partial operation
remains non-authoritative and recoverable.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| GRE-001 | The action, card, effects, and blockers are contract-owned | — |
| GRE-002 | Confirmed creation and binding are bounded and atomic | GRE-001 |
| GRE-003 | Recovery, security, compatibility, and release behavior are proven | GRE-002 |

The completed `repository-bootstrap` plan and draft
`workspace-projection-and-gate-effects` plan are explicit dependencies.

## Acceptance criteria

- GRE-AC-01: Greenfield creation is a distinct, fully described gate.
- GRE-AC-02: Only exact confirmation creates and binds the repository.
- GRE-AC-03: Unsafe, ambiguous, existing, or incomplete targets fail safely.
- GRE-AC-04: Execute never provisions product Git.
- GRE-AC-05: Shared projections and local bindings update atomically.
- GRE-AC-06: Existing bootstrap and isolation behavior remains unchanged.

## Verification

GRE-VT-01 through GRE-VT-08 cover contracts, routing, ownership, recovery,
security, upgrades, release assembly, and complete semantic acceptance.

## Risks, assumptions, and open decisions

- This plan is needed only if Context Circuit supports local greenfield
  products. Approval should confirm that product decision.
- Initial branch naming must use configured identity or an explicit card field;
  the operation must not silently depend on host Git defaults.
- A remote-repository provider action remains out of scope.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Must
requirements 1, 5, and 6. CC-001 is accepted as completed groundwork rather
than reimplemented. The source PRD remains passive and is not Product
Knowledge.
