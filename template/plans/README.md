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
- Plan status is human-controlled: `draft → approved → done`.
- `plans/INDEX.md` lists only active plans. Archived plans are excluded from
  normal agent reads until an explicit restore.

Approval, execution, completion, archive, restore, and delivery are separate
explicit human actions. Creating or reviewing a plan never approves or executes
it. Contract owner: `wrapper/contracts/schemas/plan.yaml`.
