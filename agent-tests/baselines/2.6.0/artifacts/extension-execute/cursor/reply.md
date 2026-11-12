I'll execute plan p1000. First I'll read the execution procedure and the plan itself so the work follows this workspace's gates.Plan p1000 is implemented on branch `cc/p1000/api` in the api worktree (commit `618ffe2`).

`WriteCSV` in `internal/invoice` writes a header and one row per invoice: number, customer, status, and total in cents. `TestWriteCSVIssuedInvoice` covers a single issued invoice. The plain-text summary is unchanged.

`go test ./...` passed. The working tree is clean. Nothing was pushed.

## Parity

- The plain-text summary still lists every invoice: **kept**. `TextSummary` is unchanged.
- Totals stay in cents: **kept**. The CSV total is `Invoice.Total()`, which sums line amounts in cents.

Context Circuit plugin 0.5.0 is available (this environment has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.