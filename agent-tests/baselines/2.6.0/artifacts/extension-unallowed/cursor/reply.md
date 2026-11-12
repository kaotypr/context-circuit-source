I'll execute plan p1000. First I'll read the execution skills and the plan so the work starts from the recorded outcome.The plan is a single CSV writer in the api repository. I'll prepare its worktree and run the implement hooks before changing code.p1000 is implemented on `cc/p1000/api` in the api worktree at `.worktrees/p1000/api`.

`WriteCSV` in `internal/invoice` writes a header and one row per invoice: number, customer, status, and total in cents. The plain-text summary is unchanged. `TestWriteCSVIssuedInvoice` covers a single issued invoice.

`go test ./...` passed (`ok example.com/api/internal/invoice`). The work is committed as `6840798` — `feat(invoice): export invoices as CSV`. The working tree is clean.

The after-implement hook skipped `notify` because this machine has not allowed it. Allowing it is `extension allow --id notify`, which shows what the scripts do first. Nothing else from the hooks is waiting on you.

The plan is not marked complete. Say when you want that, or when you want this branch delivered.

Context Circuit plugin 0.5.0 is available (installed 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.