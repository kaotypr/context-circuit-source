Plan p1000 is executed. The CSV export is committed on the worktree's own branch, and the Go tests pass.

**What changed** (commit `769c367` on `cc/p1000/api`, worktree `.worktrees/p1000/api`):
- `internal/invoice/csv.go`: `WriteCSV(w io.Writer, invoices []Invoice) error` writes a header `number,customer,status,total_cents`, then one row per invoice. Total is in cents. It uses `encoding/csv`, so fields are quoted correctly.
- `internal/invoice/csv_test.go`: `TestWriteCSVSingleIssued` checks the exact output for one issued invoice with two lines (100 + 25), expecting `125`.

I chose the function name, signature, and column name `total_cents`. The plan didn't specify them.

**Checks I ran:**
- Baseline `go test ./...` before the change: pass.
- After the change: `gofmt -l` clean, `go vet ./...` clean, `go test ./...` pass.
- `TestWriteCSVSingleIssued` and `TestTextSummary` both pass when run verbosely.

**Parity** (from the plan's `## Parity` section, per the parity extension's after-implement instruction):
- The plain-text summary still lists every invoice: **kept**. `export.go` has no diff from the base commit, and `TextSummary` still passes its test.
- Totals stay in cents: **kept**. The CSV writes `Total()` as an integer count of cents.

**Not done, waiting on you:**
- Nothing is pushed, and no pull request is open. Both need your authorization.
- The `cc/p1000/api` branch is local only.
- The plan record is unchanged, as the execution rules require. Completing p1000 is a separate request.

The worktree has no uncommitted changes.