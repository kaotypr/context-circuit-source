# Intent

An **intent** is the first-class decision for one change — the bigger picture a
human approves *before* it fragments into plans and tasks (Context Circuit v1.0).
Approving an intent is the single upstream human gate (Gate 1); it also confirms the
coordinator understood the ask, which is what lets the tracer read the real code next.
Plans then derive from the approved intent automatically, with no separate plan gate
and no automated scope gate (scope-safety is settled at delivery, Gate 2).

```
intent/
├── INDEX.md                # active intent index (this workspace)
├── <intent-id>/
│   ├── INTENT.md           # readable: goal, shape, what's out of scope
│   ├── contract.yaml       # canonical: outcome criteria, coarse optional scope, tier, frozen digest
│   └── trace/              # the tracer's per-repository manifests (recorded on approval)
│       └── <repository>.yaml
└── archive/<intent-id>/    # archived intents; never normal agent context
```

- Intent IDs are stable and take the form `i<NNN>-<kebab-slug>` (for example
  `i001-checkout-retries`) — distinct from plan ids so an id reads as an intent
  on sight. The sequence is never reused.
- Intent status is human-controlled: `draft → approved`. Approval freezes
  `contract_digest`, the identity of the criteria the change is later proven
  against; changing a criterion after approval is a new decision that re-enters
  the gate.
- `intent/INDEX.md` lists only active intents. Archived intents are excluded from
  normal agent reads until an explicit restore.

Authoring or approving an intent never creates a plan, executes, verifies,
completes, or delivers. Contract owner:
`.context-circuit/wrapper/contracts/schemas/intent-contract.yaml`.
