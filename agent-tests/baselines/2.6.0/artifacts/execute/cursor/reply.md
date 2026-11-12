I'll execute plan p1000. First I'm reading the plan and the workspace rules that govern how execution starts.The api worktree is on `cc/p1000/api`. I'll read the invoice package and add the CSV export there.Plan **p1000** is implemented on branch `cc/p1000/api` in the api worktree.

`WriteCSV` in `internal/invoice` writes a header and one row per invoice: number, customer, status, and total in cents. The plain-text summary is unchanged. `TestWriteCSVIssuedInvoice` covers a single issued invoice.

`go test ./...` passed. The change is committed as `d710991` (`feat(invoice): export invoices as CSV`). The working tree is clean. The plan is not marked complete.

The context-circuit plugin here is 0.5.0 (installed 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.