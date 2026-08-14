# Conventions

- Keep canonical behavior under `.agents/` and host adapters thin.
- Use ordinary YAML and Markdown with small mechanical lint checks.
- Keep plan and task statuses exactly `draft`, `approved`, or `done`.
- Never infer status from Git, tests, worktrees, external systems, or agent output.
- Keep Product Knowledge concise and source-cited; refresh it only after human review.
- Preserve dirty or uncertain work; never reset, stash, clean, merge, or deploy implicitly.
- Keep `.runtime/` private and preserve it until a human chooses cleanup.
- Never store credentials or external activity records.
