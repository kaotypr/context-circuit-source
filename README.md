# Context Circuit

Context Circuit is a small, human-controlled context, planning, and agent-work tool. It keeps Product Knowledge, numbered plans, task dependencies, and isolated plan worktrees in ordinary YAML and Markdown.

The normal journey is:

1. If the workspace is empty, use `$cc-idea-brief` to discuss the idea and create `context/IDEA-BRIEF.md`.
2. After human confirmation, import useful, source-cited facts into Product Knowledge.
3. Create a structured draft plan and task breakdown.
4. Human explicitly approves the plan.
5. `whats-next` recommends an approved plan with unfinished work.
6. Optionally publish the plan and its tasks before execution.
7. Run the plan; the agent works continuously through its tasks in one isolated domain worktree.
8. Human reviews the whole plan once, then explicitly marks tasks and the plan done.
9. Archive the plan manually whenever useful.

Plans and task statuses are exactly `draft`, `approved`, or `done`. They are not inferred from Git, tests, publication, external issues, or agent output.

## Commands

The bundled command is `node .agents/bin/cc.mjs`.

- `configure-workspace` — configure or inspect the wrapper.
- `import-product-knowledge` — import source-cited knowledge and record provenance.
- `refresh-product-knowledge` — detect changed sources and propose a reviewable refresh.
- `create-plan` — create numbered draft plans and tasks.
- `whats-next` — read-only dependency-aware recommendation.
- `cc-idea-brief` — discuss an early idea and create a human-reviewed Idea Brief.
- `run-task --plan <ref>` — run the entire approved plan continuously in its isolated domain worktree.
- `review-plan` — optional read-only review of the whole plan after execution.
- `publish-plan` — optional publication of the plan and its tasks before execution.
- `archive-plan` / `unarchive-plan` — manual location-based plan archival.
- `set-plan-state` / `set-task-state` — explicit human status changes.
- `validate` / `validate-plan` — lightweight YAML, Markdown, reference, and graph checks.

See [getting started](docs/getting-started.md) and [the command reference](docs/command-reference.md).
