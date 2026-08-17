# Conventions

- Keep normative behavior in AGENTS.md and WORKFLOW.md.
- Keep host adapters thin and expose one workspace-entry workflow.
- Keep Product Knowledge concise, scoped, and source-cited.
- Keep workspace identity in WORKSPACE.md and project identity in PROJECT.md.
  Do not add a third every-session identity page.
- Keep raw source files in `sources/`; never copy them into accepted context.
- Read only selected source files for a source-based request and record the
  source path, reason, and revision or freshness when known.
- Store Idea Briefs and PRDs at a user/team-selected path under `sources/`;
  do not impose a subdirectory or filename convention. Record the exact
  chosen path in provenance. Drafts remain visibly unaccepted until a human
  gate.
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
  child. `cc-run-plan` is the sole standard single-plan execution skill and
  uses the repository default or active branch. Use `cc-run-stack` for a
  connected set of already-approved plans: freeze `graph.yaml`, resume from
  `progress.yaml`, base dependents on parent frozen SHAs, and join
  multi-parent leaves in-run. Do not add a scheduler or a durable stack
  artifact under `plans/`. `workspace.yaml` mode is identity, not an
  execution-topology selector. Delivery policy IDs are `remote-review`,
  `local-target`, and `manual`.
- Use `cc-approve-plan` for the human plan-approval gate, `cc-finish-plan`
  for the human status-change gate, and `cc-cleanup-runtime` to delete
  `.runtime/` after inspecting dirty or unpushed work. These are
  discoverable skills, not mandatory ceremonies for every session.
- Never infer completion from Git, tests, worktrees, external systems, or agent
  output.
- Keep provider-specific activity state separate under `external_status` when an
  explicitly configured adapter needs it; never store credentials or external
  activity records in ordinary workspace state.
- Preserve dirty or uncertain work; never reset, stash, clean, merge, or deploy
  implicitly.
- Keep .runtime/ private and preserve it until a human chooses cleanup via
  `cc-cleanup-runtime`.
- Never store credentials or external activity records.
