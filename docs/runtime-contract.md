# Runtime contract

The record shapes are owned by `wrapper/contracts/schemas/`; this page defines
how the records relate.

```text
.runtime/
├── .transactions/<transaction-id>/
│   ├── transaction.yaml
│   ├── records/
│   │   ├── session.yaml
│   │   ├── context-receipt.yaml
│   │   ├── delegation.yaml
│   │   ├── child-start.yaml
│   │   ├── handoff.md
│   │   └── completion.yaml
│   ├── launch.yaml
│   └── commit.marker
├── sessions/<session-id>/
│   ├── session.yaml
│   ├── delegation.yaml
│   ├── context-receipt.yaml
│   └── handoff.md
├── plans/<plan-id>/
│   ├── lease.lock/owner.yaml
│   ├── lease.yaml
│   ├── completion.yaml
│   ├── reviews/
│   ├── evidence/
│   └── handoffs/
├── stacks/<stack-id>/
│   ├── graph.yaml
│   ├── progress.yaml
│   └── lease.lock/owner.yaml
└── worktrees/<repository-key>/<plan-id>/
```

There is no global current-session, current-plan, or current-stack pointer.
Optional host binding enables lookup but never authorizes mutation.

The transaction directory is engine-owned staging evidence. Readers require a
validated `transaction.yaml` and `commit.marker` before treating the graph as
authoritative. An interrupted or unmarked transaction is preserved for
diagnosis but cannot grant ownership, resume, launch, verification, or
completion eligibility.

## Repository bindings and worktrees

`workspace.yaml` is portable repository identity. The ignored root
`repositories.local.yaml` is the only host-local binding source and may point
to an explicit external path, workspace-relative path, or
`repositories/<repository-key>`. Resolution records the selected path,
inspectable Git identity, and clean/blocked result; it never scans for a
repository or persists credentials. A valid binding is prepared only at
`.runtime/worktrees/<repository-key>/<plan-id>/`. The bound source checkout is
not the writer worktree and dirty sources remain preserved and blocked.

## Sessions and delegation

A root session has no parent. A child has parent and root IDs, role, objective,
bounded paths, non-goals, context set/receipt, repository/worktree, permissions,
acceptance, verification, expected evidence, stop conditions, and handoff
schema. Missing fields block child entry. Workers do not change canonical
status. A verifier's implementation, plan, lease, worktree, and activity writes
are false.

## Lease and worktree

The first atomic `mkdir lease.lock` wins. `owner.yaml` and `lease.yaml` must
agree on plan, session, repository, worktree, and branch. Contenders remain
read-only. Only the matching owner releases. Stale or ambiguous state requires
human-gated takeover evidence; heartbeat expiry is not ownership.

## Receipts and handoffs

A receipt contains wrapper version, route digest, context set, evidence paths,
revisions/bytes, and invariant IDs. It is a cache key, not truth. On resume,
changed or newly required primary evidence is reloaded; an invalid receipt
cannot authorize a write.

A handoff records observed state, route/reason, evidence read, action, changed
state, verification, blockers/human decisions, and one next safe action. It may
locate evidence but cannot override the plan or accepted context.

Resume consumes the current primary records through
`cc_runtime_graph_authoritative`, `cc_validate_runtime_graph`, and receipt
validation. Changed references, stale digests, foreign leases, mismatched
ancestry, or incomplete child evidence fail closed; resume never reconstructs
records from a template or prior prompt.

## Evidence layers

`wrapper/contracts/schemas/plan.yaml` owns `required_layer` and
`produced_layer`. Runtime records carry `observed_layer`, `outcome`, and
`evidence_ref` through delegation, handoff, and completion. The engine compares
required and observed layers with exact match only. Writer-recorded mappings
are claims and cannot satisfy independent verification. Only `passed` satisfies
completion; `failed`, `blocked`, and `waived` remain non-passing. Missing host
capability is recorded as a non-passing limitation and is never converted into
a pass.

## Completion and stacks

Completion evidence is runtime only: all task evidence, independent verifier
pass, and clean committed worktree. `completion.yaml` can request a human
status-change, but it cannot mark the plan done. A stack freezes `graph.yaml`
once and uses `progress.yaml` as its resume cursor. Parent frozen SHAs are
joined deterministically; implemented parents unblock dependents; no scheduler
or auto-finish exists.
