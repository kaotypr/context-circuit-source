I'll execute plan p1000. First I'm reading the execution procedure and the plan itself so the work starts from the recorded outcome.Plan p1000 is ready to implement, and the worktree is prepared, but a required extension stopped the work before any code change.

The api worktree is at `.worktrees/p1000/api` on branch `cc/p1000/api`, started from `main` at `d82563a2`. The plan adds a CSV invoice export in `internal/invoice` (header plus one row per invoice: number, customer, status, total in cents) and a test for a single issued invoice, then `go test ./...`.

The `contract` hook at **before implement** failed and its failure policy is **stop**:

> contract check failed: 2 of 14 orders-service behaviors differ from the recorded contract

That hook is not re-run until you decide how to proceed. Nothing in the worktree was edited or committed.

The context-circuit plugin here is 0.5.0 (installed as 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.