I'll execute plan p1000. First I'm reading the execution skills and the plan so the work follows this workspace's process.Plan p1000 is implemented on branch `cc/p1000/api` in the api worktree, commit `a24c254` (`feat(invoice): add CSV export of invoice totals`).

`WriteCSV` in `internal/invoice` writes a header and one row per invoice: number, customer, status, and total in cents. The plain-text summary is unchanged. `TestWriteCSVIssuedInvoice` covers a single issued invoice.

`go test ./...` passed (`ok example.com/api/internal/invoice`).

The style extension's after-implement hook failed and asked only to report it: "style check: 3 warnings in internal/invoice". It returned no further instructions. The branch is local and has not been pushed.

Context Circuit plugin 0.5.0 is available (this session has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.