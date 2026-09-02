# Plans

Each plan is a workspace-level bundle so one plan can map changes across
multiple repositories:

```
plans/
├── INDEX.md              # active plan index (this workspace)
├── <plan-id>/
│   ├── PLAN.md           # readable plan
│   ├── plan.yaml         # canonical status, ids, repository map, tasks
│   └── tasks/            # per-task detail files
└── archive/<plan-id>/    # archived plans; never normal agent context
```

- Plan IDs are stable and take the form `NNNN-<kebab-slug>` (for example
  `0001-billing-v2`). The sequence is never reused.
- Plan status is `draft → done`: a plan derives from an approved intent within its
  scope envelope (no separate plan approval); Standard `done` follows candidate
  acceptance plus delivery, while Critical `done` requires explicit completion.
- `plans/INDEX.md` lists only active plans. Archived plans are excluded from
  normal agent reads until an explicit restore.

Intent approval and delivery are the two explicit human gates. Execution follows
an approved intent within its envelope; Standard completion is inferred from
candidate acceptance plus delivery, while Critical completion is explicit.
Explore is planless. Creating or reviewing a plan never approves or executes it,
and verification never implies completion. Contract owner:
`wrapper/contracts/schemas/plan.yaml`.
