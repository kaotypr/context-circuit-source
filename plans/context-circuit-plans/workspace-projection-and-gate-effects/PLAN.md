# Bounded identity projections and explicit gate effects

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

This plan isolates identity-agreement and confirmation-integrity work from
the broader source PRD. `workspace.yaml` stays the workspace identifier.
`context/` stays Product Knowledge. A bounded identity region in the three
Tier 0 summaries must agree with that identifier. Gate cards show proposed
defaults and separate immediate effects from later authorized effects.

## What approval authorizes

Approval authorizes bounded changes to the canonical workspace schema, route
and invariant owners, identity-region engine operations, identity-acceptance
and repository-registration card templates, other named gate presentations,
mutable seed, migration, documentation, and semantic fixtures named by this
plan.

Approval does not start execution, claim a lease, create a worktree, change
Git, deliver, publish, deploy, merge, or clean runtime.

## Scope and non-goals

`plan.yaml` owns the exact scope. This plan does not generate Product
Knowledge from `workspace.yaml`, add repository creation, replace repository
bootstrap, change lifecycle status, generate runtime child records, or add a
second effect-policy owner.

## Proposed solution

Keep identity metadata in `workspace.yaml`. Keep authored Product Knowledge
outside a contract-owned delimited identity region in `WORKSPACE.md`,
`PROJECT.md`, and `INDEX.md`. One engine operation may rewrite only that
region, and entry or write preflight fails if the region disagrees with
`workspace.yaml`. Registering a logical repository updates shared identity
and all three regions in the same operation; host-local paths stay in
`repositories.local.yaml`.

Add identity-acceptance and repository-registration card templates to
`docs/gates.md`. Omitted confirmation fields appear as proposed defaults on
the current card; confirmation records those displayed values. Fields with
no default reject as incomplete. No field is invented after confirmation.

Lock effect identifiers on the existing route and gate owners. They describe
effects; they do not authorize a route or skip confirmation:

| Identifier | This confirmation may change now | Later authorized action |
| --- | --- | --- |
| `workspace.accept_identity` | Accepted identity fields and identity regions | — |
| `workspace.register_repository` | Shared logical key, optional URL, optional branch, identity regions | clone, create-empty, or execute |
| `repository-bootstrap` | Clone destination after its own gate | execute |
| `repository-create-empty` | Reserved; not activated by this plan | — |
| `git.commit` | Isolated worktree only | delivery |
| `delivery.push` | — | delivery gate |

`git.init` and `git.clone` are not `execute-plan` effects. Delegated effects
must be a subset of approved intent.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| WPE-001 | Identity-region, default-display, and effect contracts have one owner | — |
| WPE-002 | Engine and cards render, validate, and present those contracts atomically | WPE-001 |
| WPE-003 | Migration, release, documentation, and semantic suites prove the behavior | WPE-002 |

The completed `repository-bootstrap` plan is a prerequisite because repository
registration already owns shared identity and host-local binding boundaries.

## Acceptance criteria

- WPE-AC-01: One operation renders or validates only the identity region.
- WPE-AC-02: Identity-region disagreement blocks entry or write preflight.
- WPE-AC-03: Cards separate immediate effects from later authorized effects.
- WPE-AC-04: Omitted fields are displayed defaults; undeclared delegation fails closed.
- WPE-AC-05: Repository registration updates shared identity and identity regions atomically.
- WPE-AC-06: Existing gated flows remain valid.

## Verification

WPE-VT-01 through WPE-VT-07 cover contracts, gates, routing, recovery,
upgrades, release assembly, and full semantic acceptance. Executable commands
remain canonical in `plan.yaml`.

## Risks, assumptions, and open decisions

- Exact delimiter syntax for the identity region is owned by the workspace
  schema during WPE-001; the region count is one per file and the engine
  cannot rewrite authored Product Knowledge outside it.
- `roles` is a displayed identity-acceptance field. If it is missing from
  the workspace schema, WPE-001 adds a canonical optional field rather than
  allowing post-confirmation invention.
- Effect identifiers belong with existing route and invariant owners. A new
  parallel policy file is not authorized.
- `repository-create-empty` is reserved for the dependent create-empty plan
  and is not activated here.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Must
requirements 2 through 4, refined by the 2026-08-23 named-plan review.
Locked review decisions: `context/` remains Product Knowledge;
`workspace.yaml` remains the workspace identifier; omitted fields display as
proposed defaults then confirm; identity-acceptance and
repository-registration card templates are in-scope named paths.

Current repository evidence includes CC-001 repository ownership, CC-004
bounded approval behavior, the workspace schema, route/invariant owners, and
the required mutable context seed. The source PRD remains passive and is not
Product Knowledge.
