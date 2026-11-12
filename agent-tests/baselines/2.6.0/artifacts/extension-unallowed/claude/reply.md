Plan p1000 is implemented and committed on its own branch. Nothing has been pushed and no plan record was changed.

**What changed** (in the `api` repository, branch `cc/p1000/api`, worktree `.worktrees/p1000/api`):
- `internal/invoice/export.go`: new `CSVExport([]Invoice) ([]byte, error)`. It writes a header `number,customer,status,total_cents` and one row per invoice, with the total in cents. It sits beside `TextSummary`, which is unchanged.
- `internal/invoice/export_test.go`: new `TestCSVExportIssuedInvoice`. It checks the exact output for one issued invoice with two lines (100 + 25), expecting `125` in the total column.
- Commit `93c1d21 feat(invoice): add CSV export`. The worktree is clean.

**Checks run:** `gofmt -l internal/` (no output), `go vet ./...` and `go test ./...`. All passed (`ok example.com/api/internal/invoice`).

**Choices I made:**
- The name `CSVExport` and the column name `total_cents` were my picks; the plan didn't name them.
- The writer takes a slice, so it covers several invoices, while the test covers one issued invoice as the plan asked.

**Extensions:** the `before implement` hook returned nothing. The `after implement` hook skipped the `notify` extension because it isn't allowed on this machine. Running `extension allow --id notify` is your decision, and I haven't run it. Nothing in the plan depended on it.

**Left for you:** pushing `cc/p1000/api`, opening a PR against `main`, and marking p1000 completed all need your go-ahead. I haven't done any of them.