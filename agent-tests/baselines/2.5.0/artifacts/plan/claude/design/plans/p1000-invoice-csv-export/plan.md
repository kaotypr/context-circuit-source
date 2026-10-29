---
id: p1000
created_by: maya
created_at: 2026-10-29T06:09:17Z
intent: i100
repositories:
- api
---

# Invoice CSV export

## Outcome and repository scope

Callers of the `internal/invoice` package can write a set of invoices as CSV. Each row carries the invoice number, the customer, the status, and the total in cents. Scope is the `api` repository only, confined to `internal/invoice`. Voiding, credit notes, and the plain-text `TextSummary` stay unchanged. No third-party dependency is added to `go.mod`.

## Approach

`internal/invoice/export.go` holds `TextSummary` beside the new function. The CSV writer goes in the same package, in a new file `internal/invoice/csv.go`, so the existing export file is not reshaped:

```go
// WriteCSV writes one header row and one row per invoice to w.
func WriteCSV(w io.Writer, invoices []Invoice) error
```

- Use the standard library `encoding/csv`, which keeps the module free of third-party dependencies and handles quoting of commas, quotes, and newlines in customer names.
- Header: `number,customer,status,total_cents`. Row values: `Invoice.Number`, `Invoice.Customer`, `Invoice.Status`, `Invoice.Total()`.
- Take `[]Invoice` rather than one invoice so a caller writes a whole export in one call; a single invoice is `[]Invoice{inv}`.
- Flush the `csv.Writer` and return its error, so a failed write to `w` is reported to the caller.
- No HTTP handler exists in this repository. The "customers download" framing in the intent is served by a caller that passes `WriteCSV` an `http.ResponseWriter`; wiring a route is outside this plan's success criteria and is not added here.

## Tasks and order

1. Add `internal/invoice/csv.go` with `WriteCSV` as above.
2. Add `internal/invoice/csv_test.go`: a table of invoices (one with a comma and a quote in the customer name, one with zero lines) checked against the exact expected output, plus a test that a writer error is returned.
3. Run `go test ./...` in `repositories/api` and confirm it passes.
4. Commit on the worktree branch. The `api` AGENTS.md asks for `type(scope): subject` messages, e.g. `feat(invoice): add CSV export`.

Tasks 1 and 2 are one commit; task 3 is verification before commit.

## Risks and checks

- **Spreadsheet formula injection.** A customer name beginning with `=`, `+`, `-`, or `@` can execute as a formula when the CSV is opened in a spreadsheet. The intent does not address it and this plan does not sanitize values, because that changes the data a caller reads back. Flagged for the person to decide; if they want it, it is a separate small change.
- **Column names are a public contract.** Once callers depend on `total_cents`, renaming it breaks them. The names above follow the intent's wording ("the total in cents").
- **Stale knowledge note.** `context/domains/invoices.md` says invoices export only as a plain-text summary. That becomes false on completion; it is reconciled by the completion step, not by this plan's implementation.
- Checks: `go test ./...` in `repositories/api`, which must pass (intent success criterion). `go vet ./...` is an ordinary extra check.
