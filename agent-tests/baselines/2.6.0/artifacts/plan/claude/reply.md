Plan **p1000 "Invoice CSV export"** is written for intent i100 and linked to it. Nothing is executed yet.

**Plan:** `plans/p1000-invoice-csv-export/plan.md`

- **Outcome:** add `WriteCSV(w io.Writer, invoices []Invoice) error` in a new `internal/invoice/csv.go`. It writes a header row, then one row per invoice with `invoice_number`, `customer`, `status`, and `total_cents`.
- **Approach:** `encoding/csv` handles quoting, so commas and quotes in customer names are safe. The total comes from the existing `Invoice.Total()`, so it matches `TextSummary`. The standard library only, so no new dependencies.
- **Tasks:** (1) add `csv.go`; (2) add `csv_test.go` covering column order, quoting, a zero total, an empty slice, and a failing writer; (3) run `go test ./...` in `repositories/api`.
- **Checks:** `go test ./...`, `go vet ./...`, and `gofmt -l internal/invoice`. I ran `go test ./...` on the untouched repo and it passes, so the baseline is green.

**Two things for you to decide:**

1. **The download route.** i100's goal says customers can download invoices, but the repo has no HTTP handler, and the constraints confine the work to `internal/invoice`. The plan delivers only the writer a download route would call. If you want the route, it needs its own intent or a change to i100's constraints.
2. **Header row.** The intent doesn't say whether to include one. The plan includes it. Say if you'd rather not have it.

**Record changes:** `plans/p1000-invoice-csv-export/plan.md` is new. `intent/i100-invoice-csv-export.md` gained a `plans: [p1000]` link, and `.context-circuit/ids/maya.yaml` reserves `p1000`. These are uncommitted. `check` reports no issues.

Should I commit the records, and push them, or commit only? Or leave them uncommitted while you read the plan?