# Agent workspace workflow

This is the normative behavior reference for the released wrapper. The human
interface is conversational: “Start or resume work in this workspace.” The
host reads the entry adapters, asks the single router for one bounded probe,
then performs or proposes one action.

## Entry and evidence

1. Read the safety adapter, workflow adapter, manifest, compact Tier-0 index,
   workspace identity, and three context summaries.
2. Identify root/child role and any exact host binding.
3. Stage A selects one context set or one focused question.
4. Load only that set through the host-neutral packet operation; it validates
   exact paths, selected conditional evidence, measured bytes, and receipt
   digests before Stage B emits one normalized decision.

Direct host reads are not represented as sandbox enforcement when a host cannot
provide it. The semantic overread fixture must fail or produce a host-blocked
read-only outcome; it never widens the selected packet.

Separate observed state, instruction/decision, assumption, blocker, action, and
verification in every consequential handoff. Runtime evidence cannot redefine
plans or accepted context. The source inbox is passive and never scanned.

## Roles and isolation

The root coordinator owns the human request, route, preflight, lease, worktree,
delegation, cards, and evidence consolidation. The writer changes only assigned
paths in one exclusive worktree. The independent verifier is read-only for
implementation, plan, lease, worktree, and activity state and writes only its
own handoff. Missing delegation fields or child primitives block entry.

## Lifecycle and gates

Review is read-only. Approval changes plan status only. Execution requires a
separate explicit named request. Verification produces evidence. Finish is a
separate human status-change. Delivery, merge, publication, deployment,
archive, restore, takeover, and destructive cleanup each use an exact current
confirmation card. Eligibility is never authorization.

## Recovery

Interruption preserves sessions, leases, branches, worktrees, dirty state,
receipts, and handoffs. Resume validates wrapper compatibility, primary
evidence digests, Git state, and ownership. A live foreign owner or ambiguous
lease remains read-only; heartbeat expiry alone never grants takeover.

## Offline boundary

Filesystem operation remains complete when an optional provider is disabled,
denied, or unavailable. Workspace configuration stores provider-neutral intent,
never credentials or provider payloads. There is no scheduler, database,
global current pointer, implicit publication, or user-facing CLI.

Canonical owners are `wrapper/contracts/routes.yaml`,
`wrapper/contracts/invariants.yaml`, `wrapper/contracts/context-sets.yaml`,
the schemas, and `wrapper/runtime/engine.sh`.
