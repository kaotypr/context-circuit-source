# Conventions

- Keep normative behavior in AGENTS.md and WORKFLOW.md.
- Keep host adapters thin and expose one workspace-entry workflow.
- Keep Product Knowledge concise, scoped, and source-cited.
- Use ordinary Markdown and YAML for context, plans, session records, and
  handoffs.
- Keep plan and task statuses exactly draft, approved, or done.
- Keep runtime session state separate from plan and task status.
- Give every child session an explicit parent, objective, scope, permissions,
  and handoff format.
- Give every writing session an exclusive worktree.
- Never infer completion from Git, tests, worktrees, external systems, or agent
  output.
- Preserve dirty or uncertain work; never reset, stash, clean, merge, or deploy
  implicitly.
- Keep .runtime/ private and preserve it until a human chooses cleanup.
- Never store credentials or external activity records.
