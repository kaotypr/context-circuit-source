Plan p1000 is implemented and committed on the worktree branch `cc/p1000/api`. Nothing is pushed, and the plan record is unchanged.

**What changed** (commit `51b4a23 feat(invoice): add CSV invoice export`, in `.worktrees/p1000/api`):
- `internal/invoice/csv.go`: `WriteCSV(w io.Writer, invoices []Invoice) error` writes a header `number,customer,status,total_cents` and one row per invoice. The total is in cents, from `Invoice.Total()`. It uses `encoding/csv`, so customer names with commas are quoted.
- `internal/invoice/csv_test.go`: covers a single issued invoice (two lines, total 125) and a customer name that needs quoting.
- `TextSummary` is untouched.

**Checks I ran:**
- `gofmt -l ./internal`: no output.
- `go vet ./...`: clean.
- `go test ./...`: `ok example.com/api/internal/invoice`.

**Style extension reported a failure.** The after-implement hook exited 1 with "style check: 3 warnings in internal/invoice". Its `on_failure` is `report`, so it does not block. The hook script is a stub that prints that fixed line and does not say which warnings it means. `gofmt` and `go vet` both pass, so I can't tell from here what it flagged. I did not chase it further. This is worth a look by you or whoever owns the style extension.

**Waiting on you:**
- Delivery: pushing `cc/p1000/api` and opening a PR against `main` both need your explicit go-ahead.
- Completion: marking p1000 completed is a separate request.