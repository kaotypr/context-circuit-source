I'll execute plan p1000. First I'm checking workspace setup and reading how execution is supposed to run.The invoice type is in place. I'll add a CSV writer beside the plain-text summary and cover one issued invoice with a test.Plan p1000 is implemented on branch `cc/p1000/api` in `.worktrees/p1000/api` (commit `3448d21`).

`WriteCSV` in `internal/invoice` writes a header (`number`, `customer`, `status`, `total`) and one row per invoice, with the total in cents from `Invoice.Total`. `TestWriteCSVIssuedInvoice` covers a single issued invoice. The plain-text summary is unchanged. `go test ./...` passed.

The csv-check repair loop ran after the commit. The verifier's verdict was clean after 1 round, with no findings and no further commits.

The branch is local only. Say if you want it pushed or opened as a pull request.

Context Circuit plugin 0.5.0 is available (this session has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.