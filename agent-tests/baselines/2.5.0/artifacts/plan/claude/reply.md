I planned i100 as a single plan, **p1000 "Invoice CSV export"**, at `plans/p1000-invoice-csv-export/plan.md`. It covers the `api` repository only. `check` reports no issues.

The plan adds a `WriteCSV(w io.Writer, invoices []Invoice) error` function in a new `internal/invoice/csv.go`, using the standard library's `encoding/csv` so no dependency is added. Columns are `number,customer,status,total_cents`. The plan also adds `csv_test.go` and runs `go test ./...` before committing.

Decisions to review:
- **No HTTP route.** The api repo has no handler, so "customers can download" is served by a caller passing the writer an `http.ResponseWriter`. Adding a route is outside the intent's success criteria.
- **Spreadsheet formula injection.** Customer names starting with `=`, `+`, `-` or `@` can run as formulas when opened in a spreadsheet. The plan does not sanitize values, since that changes what a caller reads back. This needs your decision; a sanitizing change would be a separate small task.
- **Column names** `total_cents` and `number` become a public contract once callers depend on them.
- **Knowledge note.** `context/domains/invoices.md` says invoices export only as plain text. That becomes false on completion, so the completion step must reconcile it.

Nothing is executed yet. The plan is not committed, and `plans/` and `intent/` changes are uncommitted in the workspace. Running the plan needs a separate request. Do you want me to commit the plan record, and do you want the injection guard included?