# Choosing the next action

Choosing the next action is part of root-session entry, not a separate command
the user must operate.

The root session inspects:

- approved plans with unfinished tasks;
- declared plan dependencies;
- active plan leases and session ownership;
- source freshness and Product Knowledge warnings;
- repository cleanliness and worktree availability;
- blockers and pending human gates.

It may recommend or claim only work that is dependency-ready, explicitly scoped,
and not already owned by another writing session. When no work is executable,
the session should explain whether it needs context, a draft plan, human
approval, a review, or a decision about a blocker.

The recommendation is evidence-backed and read-only until the root session or
human explicitly performs the next consequential action.
