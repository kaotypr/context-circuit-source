# Intent

An **intent** is the first-class decision for one change — the bigger picture a
human approves *before* it fragments into plans and tasks (Context Circuit v1.0).
Approving an intent is the single upstream human gate (Gate 1); plans then derive
from it automatically within its scope envelope.

```
intent/
├── INDEX.md                # active intent index (this workspace)
├── <intent-id>/
│   ├── INTENT.md           # readable: goal, shape, what's out of scope
│   ├── contract.yaml       # canonical: criteria, scope envelope, tier, frozen digest
│   └── adversary.md        # the spec adversary's findings on the criteria
└── archive/<intent-id>/    # archived intents; never normal agent context
```

- Intent IDs are stable and take the form `i<NNNN>-<kebab-slug>` (for example
  `i0001-checkout-retries`) — distinct from plan ids so an id reads as an intent
  on sight. The sequence is never reused.
- Intent status is human-controlled: `draft → approved`. Approval freezes
  `contract_digest`, the identity of the criteria the change is later proven
  against; changing a criterion after approval is a new decision that re-enters
  the gate.
- `intent/INDEX.md` lists only active intents. Archived intents are excluded from
  normal agent reads until an explicit restore.

Authoring or approving an intent never creates a plan, executes, verifies,
completes, or delivers. Contract owner:
`wrapper/contracts/schemas/intent-contract.yaml`.
