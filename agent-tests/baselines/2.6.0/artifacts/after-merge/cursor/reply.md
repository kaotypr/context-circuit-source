Your merge note for p1000 is a fast-forward moment. I’ll read the workspace’s `/ff` procedure and offer the next step without changing anything yet.The pull request for **p1000 — Invoice CSV export** in the **api** repository (`main`) is merged on the remote. This checkout has not moved yet, and the plan is still open.

I can do both of these in one go:

1. Fast-forward **api** `main` to the remote.
2. Mark **p1000 — Invoice CSV export** completed and reconcile the knowledge that plan changed.

One yes covers both. Say if you want only the sync or only the completion. Worktrees and `cc/*` branches stay as they are.

The Context Circuit plugin here is 0.5.0 (0.2.2 is installed). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.