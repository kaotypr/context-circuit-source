I'll plan the approved intent i100. First I'm reading the planning skill and the intent itself so the plan follows this workspace's process.**p1000** plans the approved intent **i100** (Invoice CSV export). Implementation waits for a separate request to execute it.

The api repository already renders one invoice as plain text in `internal/invoice/export.go` via `TextSummary`. CSV stays in that same file, using the standard library `encoding/csv` writer, so the module still has no third-party dependencies. Voiding and credit notes are left unchanged.

`WriteCSV` writes a header (`number`, `customer`, `status`, `total_cents`) and one row per invoice, in the order given. The total is `Invoice.Total()` (the sum of line amounts in cents). An empty list writes only the header. Customer names that contain commas go through the CSV writer so they stay in one field.

Tasks, in order:

1. Add `WriteCSV` in `internal/invoice/export.go`.
2. Test one invoice, several invoices, an empty list, and a customer name that contains a comma, in `internal/invoice/export_test.go`.
3. Run `go test ./...` in the api repository.

The plan is `plans/p1000-invoice-csv-export/plan.md`. It has no supporting files.

Context Circuit plugin 0.2.4 is available (0.2.2 is installed). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.