# Command reference

All commands are available through `node .agents/bin/cc.mjs <command>`.

| Command | Purpose |
| --- | --- |
| `configure-workspace --request <file>` | Configure a wrapper or inspect it with `--check-only`. |
| Skill `$cc-idea-brief` | Discuss an early idea and create a human-reviewed `context/IDEA-BRIEF.md`. |
| `import-product-knowledge --source <path>` | Record a source and create missing Product Knowledge pages. |
| `refresh-product-knowledge` | Print a proposal for changed sources. |
| `create-plan --input <file>` | Create one or more numbered draft plans. |
| `validate-plan <path>` | Check one plan, its tasks, references, and dependency graph. |
| `validate` | Check workspace, Product Knowledge, all active and archived plans, and collisions. |
| `whats-next` | Read-only recommendation from approved plans. |
| `run-task --plan <ref>` | Run the whole approved plan continuously in its plan-scoped domain worktree. |
| `review-plan --plan <ref> --worktree <path>` | Optional read-only review summary for the whole plan after execution. |
| `set-plan-state --plan <ref> --status <status>` | Explicitly set `draft`, `approved`, or `done`. |
| `set-task-state --plan <ref> --task <id> --status <status>` | Explicitly set a task status. |
| `publish-plan --plan <ref> --references <file>` | Optionally store plan/task external URLs before execution, without status changes. |
| `archive-plan <ref>` | Move one exact plan to `archives/plans/<repository-key>-plans/`. |
| `unarchive-plan <ref>` | Restore one archived plan to `plans/<repository-key>-plans/` after collision checks. |

Commands do not create external activity records or require generated schema files. Publication providers are optional adapters and external URLs are the only publication state stored in plan/task YAML.
