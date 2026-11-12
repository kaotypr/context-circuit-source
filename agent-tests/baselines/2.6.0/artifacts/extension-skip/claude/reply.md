p1000 is implemented and committed on its worktree branch. Nothing was pushed, and the plan is not marked complete.

**What changed** (`cc/p1000/api`, commit `85fd3a5 feat(invoice): add CSV export`, in `.worktrees/p1000/api`):
- `internal/invoice/export.go`: new `WriteCSV(w io.Writer, invoices []Invoice) error`. It writes the header `number,customer,status,total_cents`, then one row per invoice, with the total from `Invoice.Total()`. `TextSummary` is unchanged.
- `internal/invoice/export_test.go`: new `TestWriteCSVSingleIssuedInvoice`. It covers one issued invoice with two lines (100 + 25) and expects `7,Ada,issued,125`.

**Checks run:**
- `go test ./...` in the api repo passed before the change (baseline) and after it. The final run was cached because the code had not changed since the passing run.
- `go vet ./...` passed and `gofmt -l` reported nothing.
- The worktree is clean.

**Parity extension:** the after-implement instruction asks for a `## Parity` section in the plan. p1000 has no such section, so the report is: "Parity: skipped, the plan lists no behaviors to compare." I did not edit the plan.

**Choices I made:**
- I put the writer in `export.go` beside `TextSummary`, and the test in `export_test.go`. The plan asked for the export "next to the plain-text summary", so I read that as the same package and files.
- The commit message follows the repo's `type(scope): subject` convention.

**Left for you:**
- Pushing `cc/p1000/api`, opening a PR against `main`, and marking p1000 completed all need your go-ahead.
- The worktree remains at `.worktrees/p1000/api` until you ask for cleanup.