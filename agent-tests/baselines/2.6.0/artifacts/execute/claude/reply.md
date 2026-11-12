Plan p1000 is implemented and committed on the worktree branch `cc/p1000/api`. It is not pushed or delivered.

- **Change:** `internal/invoice/export.go` now has `WriteCSV(w io.Writer, invoices []Invoice) error`. It writes a header `number,customer,status,total_cents` and one row per invoice, with the total in cents. The plain-text `TextSummary` is unchanged.
- **Test:** `internal/invoice/export_test.go` has `TestWriteCSVIssuedInvoice`, which checks the exact output for a single issued invoice.
- **Commit:** `f3157eb feat(invoice): add CSV invoice export`, on top of `main` at `d82563a`. The working tree is clean.
- **Checks run:** `go vet ./...`, `gofmt -l` (no output), and `go test ./...` all pass in the worktree. The test suite passed before the change too.
- **Worktree:** `ws/.worktrees/p1000/api`.

I chose the function name, the `io.Writer` signature, and the column name `total_cents`, since the plan left them open. I put the writer in the existing `export.go` and its test in `export_test.go` so they sit next to the summary, as the plan asks.

The plan record and the workspace records are unchanged, as the workspace rules require. Pushing, opening a PR, and marking p1000 completed each need your explicit go-ahead.