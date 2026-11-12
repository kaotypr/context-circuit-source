I'll execute plan p1000. First I'll read the plan and the worktree and execution steps this workspace requires.Plan p1000 is implemented on branch `cc/p1000/api` in `.worktrees/p1000/api`, commit `aaa6552` (`feat(invoice): export invoices as CSV`).

`WriteCSV` in `internal/invoice` writes a header and one row per invoice: number, customer, status, and total in cents. `TestWriteCSVIssuedInvoice` covers a single issued invoice. The plain-text summary is unchanged. `go test ./...` passed (`ok example.com/api/internal/invoice`).

The after-implement hook asks to push that branch to origin. A push still needs your explicit go-ahead, so the branch is local only until you say to push it.

Context Circuit plugin 0.5.0 is available (this workspace has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.