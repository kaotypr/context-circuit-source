# Route-selected context packet enforcement

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

Context sets and static budgets already exist, but current routes can emit
probe names that are not registered set IDs and live agents are not constrained
to the static ledger. This plan turns the existing contract into an enforced
route-to-packet boundary.

## What approval authorizes

Approval authorizes bounded route/context alignment, complete initialization
and plan-draft packets, a host-neutral loader with path and byte enforcement,
canonical shipped fixtures, receipt integration, and semantic overread tests.

Approval does not start execution, claim a lease, create a worktree, change
Git, deliver, publish, deploy, merge, or clean runtime.

## Scope and non-goals

`plan.yaml` owns exact scope. This plan keeps one router and existing
authorization. It does not generate the runtime ownership graph, persist
provider payloads, inspect sibling workspaces, or use maintainer plans as
released examples.

## Proposed solution

Map every route to one registered context set rather than copying the Stage A
probe name into `context_set`. Publish a packet manifest with exact paths,
budget, and selected-evidence slots. Load only through a host-neutral operation
that measures bytes and writes receipt references and digests. Ship canonical
plan, task, gate, and packet fixtures so drafting and initialization do not
search tests or other workspaces. Add semantic fixtures that fail observed
overread patterns.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| RCP-001 | Routes and context contracts use one valid ID and allowlist | — |
| RCP-002 | Loader enforces paths, bytes, and receipt evidence | RCP-001 |
| RCP-003 | Released packets contain canonical bounded fixtures | RCP-001 |
| RCP-004 | Live-overread, host, upgrade, release, and semantic tests pass | RCP-002, RCP-003 |

Completed host, approval, and interactive-review plans are prerequisites whose
routes and gates must remain unchanged.

## Acceptance criteria

- RCP-AC-01: Every route emits a registered context-set ID and exact packet.
- RCP-AC-02: Initialization and drafting need no sibling or maintainer examples.
- RCP-AC-03: Ordinary turns invoke rather than read router implementation.
- RCP-AC-04: The loader rejects undeclared paths and measured overruns.
- RCP-AC-05: Product-source and named-plan routes preserve existing behavior.
- RCP-AC-06: Semantic fixtures fail live overread patterns.
- RCP-AC-07: Release and upgrade install only canonical packet evidence.

## Verification

RCP-VT-01 through RCP-VT-08 cover contracts, routing, measured budgets, host
adapters, security, upgrades, release assembly, and complete acceptance.

## Risks, assumptions, and open decisions

- A host-neutral loader can enforce only reads performed through its adapter.
  The implementation must define how undeclared direct host reads become a
  failing fixture or host-blocked outcome rather than claiming impossible
  sandbox guarantees.
- Packet examples must be released product fixtures, not copies of current
  maintainer plans.
- The source PRD's request to add `plan-draft` and `initialization` set IDs is
  refined here: route phases and context-set IDs require an explicit mapping,
  not automatic name equality.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Should
requirements 7 through 10 and 13, Could requirement 16, and associated
acceptance criteria. Current repository evidence shows existing static context
sets, receipts, and budgets, so this plan targets live enforcement rather than
recreating those contracts. The source PRD remains passive and is not Product
Knowledge.
