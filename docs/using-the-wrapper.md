# Using the wrapper

The wrapper owns Idea Briefs, context, plans, whole-plan selection, and safe worktree preparation. Product repositories own code and repository-local conventions. Codex and Claude adapters point to the same canonical behavior under `.agents/`.

Keep Idea Briefs, Product Knowledge, plans, and task files reviewable in Git. Keep `.runtime/` private and preserved while plan worktrees or prompts are still useful. Human decisions—brief confirmation, plan approval, pre-execution publication, review, status changes, and archiving—remain explicit.
