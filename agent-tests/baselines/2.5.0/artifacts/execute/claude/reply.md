Plan p1000 (Invoice CSV export) is implemented and committed on the worktree branch `cc/p1000/api`. The change is not delivered: nothing is pushed, and no pull request is open.

**What changed** (in `api`, commit `859b442 feat(invoice): csv export of invoices`):
- `internal/invoice/export.go`: new `WriteCSV(w io.Writer, invoices []Invoice) error`. It writes the header `number,customer,status,total_cents` and one row per invoice. The total is `Invoice.Total()` in cents. It uses `encoding/csv`, so customer names with commas or quotes are escaped. `TextSummary` is unchanged.
- `internal/invoice/export_test.go`: new `TestWriteCSV`, which covers one issued invoice with two lines (100 + 25). It expects `number,customer,status,total_cents\n7,Ada,issued,125\n`.

**Checks I ran:**
- `gofmt -l internal/` printed nothing.
- `go vet ./...` passed.
- `go test ./...` passed, and the verbose run showed `TestWriteCSV` and `TestTextSummary` passing.

**Choices I made:**
- The plan does not name a function signature, so I chose a writer that takes `io.Writer` and `[]Invoice`.
- I used `encoding/csv` rather than hand-built strings, for correct quoting.

**Git state:** The worktree is at `.worktrees/p1000/api` and is clean. The branch `cc/p1000/api` is one commit ahead of `main`. The plan record is unchanged, as the workflow requires.

**Next steps for you:** Delivering the branch as a pull request against `main`, or marking p1000 completed, each needs your explicit go-ahead.