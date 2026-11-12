I'll execute plan p1000. First I'll read the plan and the workspace rules that govern how a plan is run.The invoice type is in place. I'll add a CSV writer beside the plain-text summary and cover one issued invoice with a test.Plan p1000 is implemented on `cc/p1000/api` and pushed to origin.

The api worktree at `.worktrees/p1000/api` now has `WriteCSV` in `internal/invoice`. It writes a header and one row per invoice: number, customer, status, and total in cents. The plain-text summary is unchanged. `TestWriteCSVIssuedInvoice` covers a single issued invoice.

`go test ./...` passed (`ok example.com/api/internal/invoice`). The work is committed as `feat(invoice): CSV export of number, customer, status, and total` (`47dd390`). The after-implement hook asked for that branch to be published, so `cc/p1000/api` is on origin.

The plan record is still open. Say when you want it marked completed.

The Context Circuit plugin here is 0.5.0 (0.2.2 is installed). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.