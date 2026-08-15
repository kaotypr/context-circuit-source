# Conventions

- Keep normative behavior in AGENTS.md and WORKFLOW.md.
- Keep host adapters thin and expose one workspace-entry workflow.
- Keep Product Knowledge concise, scoped, and source-cited.
- Keep raw source files in `sources/`; never copy them into accepted context.
- Read only selected source files for a source-based request and record the
  source path, reason, and revision or freshness when known.
- Store Idea Briefs in `contributions/idea-briefs/` and PRDs in
  `contributions/prds/`; drafts remain visibly unaccepted until a human gate.
- Use ordinary Markdown and YAML for context, plans, session records, and
  handoffs.
- Keep plan status exactly `draft`, `approved`, or `done`; task status is the synchronized projection `draft`, `ready`, or `done`.
- Treat plan approval and completion as the only canonical lifecycle gates for
  ordinary whole-plan execution. Do not add a second task approval gate.
- Keep runtime session state separate from plan and task status.
- Give every child session an explicit parent, objective, scope, permissions,
  and handoff format.
- Give every writing session an exclusive worktree.
- Use `cc-run-plan` to direct a writer child and an independent verifier
  child. `workspace.yaml` mode is identity, not an execution-topology
  selector. `solo-local` and `team-review` remain delivery policies.
- Never infer completion from Git, tests, worktrees, external systems, or agent
  output.
- Keep provider-specific activity state separate under `external_status` when an
  explicitly configured adapter needs it; never store credentials or external
  activity records in ordinary workspace state.
- Preserve dirty or uncertain work; never reset, stash, clean, merge, or deploy
  implicitly.
- Keep .runtime/ private and preserve it until a human chooses cleanup.
- Never store credentials or external activity records.
