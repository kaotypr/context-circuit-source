# Workflow

Context Circuit is a human-controlled context, planning, and agent-work tool. Plans and Product Knowledge are ordinary reviewable files; no activity provider, database, background service, or lifecycle engine is required.

The primary journey is:

`empty workspace → Idea Brief → human confirmation → Product Knowledge → draft plan and tasks → human approval → whats-next → optional pre-execution publication → continuous plan execution in one isolated domain worktree → one whole-plan review → human marks done → manual archive`

`whats-next` is read-only. It ignores `archives/plans/<repository-key>-plans/`, considers only approved plans with unfinished tasks, requires plan dependencies to be done, reports relevant Product Knowledge references, and warns when a recorded source changed. If no plan is executable, it recommends reviewing a draft plan or creating one.

`run-task --plan` checks the approved plan, requires one repository domain and a clean repository base when creating the plan worktree, creates or reuses one branch and worktree for that plan, and writes a plan-level prompt. The agent continues through all unfinished tasks in dependency order without pausing for per-task review.

`review-plan` is one optional read-only human review after the continuous plan run. Publication is an optional pre-execution output adapter: it preserves IDs and may store current external URLs in YAML, but never changes status, starts execution, monitors external status, or synchronizes completion.

Archiving is a direct location change. It accepts draft, approved, or done plans, does not inspect task status or test/merge state, and preserves all plan content and statuses. Unarchiving requires an empty active destination.

All meaningful actions remain human-gated. Preserve dirty or uncertain work; never reset, stash, clean, merge, deploy, publish real issues, or modify `context-circuit-workspace` without a separate explicit request.
