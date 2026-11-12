Plan p1000 is implemented and committed on `cc/p1000/api`. Nothing is pushed, and the plan is not marked complete.

**What changed** (commit `5364866 feat(invoice): add CSV invoice export`, in `.worktrees/p1000/api`):
- `internal/invoice/export.go`: new `WriteCSV(w io.Writer, invoices []Invoice) error`. It writes the header `number,customer,status,total_cents`, then one row per invoice using `Invoice.Total()` for the total in cents. `TextSummary` is unchanged.
- `internal/invoice/export_test.go`: new `TestWriteCSV`. It checks the exact output for one issued invoice with two lines (100 + 25), expecting `125`.

**Checks I ran:**
- `gofmt -l internal/`: no files listed.
- `go vet ./...`: clean.
- `go test ./...` in the api worktree: `ok example.com/api/internal/invoice`.

**Extension (csv-check):** the workspace allows the `repair-loop` override for this extension. After implementation, the `after implement` hook returned that override, so I ran the CSV verifier once (round 1 of 3). It returned **Verdict: clean**, with no findings, so no repair round was needed. The verifier did not run `go test`; I ran it myself as above.

**Choices I made:**
- Put the writer in `export.go`, next to the summary, rather than a new file.
- Used `encoding/csv` and the header names `number`, `customer`, `status`, `total_cents`. The plan named the columns but not their exact spelling.

**Not done (needs you):**
- Push, pull request, and completion of p1000 are not done. Delivery and completion need your explicit go-ahead.
- The workspace has no uncommitted changes from this run. The plan record was not edited.