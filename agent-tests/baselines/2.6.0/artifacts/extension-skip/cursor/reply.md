I'll execute plan p1000. First I'll read the plan and the execution skills so the work follows this workspace's process.The invoice package already has a plain-text summary. I'll add a CSV writer beside it, with a test for one issued invoice.Plan p1000 is implemented on branch `cc/p1000/api` in the api worktree at `.worktrees/p1000/api`.

`WriteCSV` in `internal/invoice` writes a header and one row per invoice: number, customer, status, and total in cents. The plain-text summary is unchanged. `TestWriteCSVIssuedInvoice` covers one issued invoice whose total is 100 cents.

`go test ./...` passed (`ok example.com/api/internal/invoice`). The change is committed as `feat(invoice): add CSV export of invoice totals` (`48bcc71`). The worktree is clean.

Parity: skipped, the plan lists no behaviors to compare.

The plan is not marked complete. Say if you want it delivered or completed.

Context Circuit plugin 0.5.0 is available (installed 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.