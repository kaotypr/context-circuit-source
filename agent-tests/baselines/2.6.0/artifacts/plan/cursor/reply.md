I'll plan the approved intent i100. First I'll read the planning skill and the intent so the plan follows this workspace's process.The invoice CSV export plan (p1000) is written for the approved intent. Customers' invoices can be written as CSV from the api, next to the existing plain-text summary. Say when you want it implemented.

The writer goes in `internal/invoice/export.go`, beside `TextSummary`, and uses `encoding/csv` from the standard library. Each invoice is one row: number, customer, status, and total in cents from `Total()`. A header row comes first, and a name that contains a comma or a quote is quoted by the CSV writer. Voiding a draft and crediting an issued invoice stay as they are.

The check is `go test ./...` in the api repository.

The plan is at `plans/p1000-invoice-csv-export/plan.md`. There are no extra files beside it.

Context Circuit plugin 0.5.0 is available (this workspace has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.