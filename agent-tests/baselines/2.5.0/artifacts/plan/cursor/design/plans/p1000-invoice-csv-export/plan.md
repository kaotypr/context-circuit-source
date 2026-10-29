---
id: p1000
created_by: maya
created_at: 2026-10-29T06:12:04Z
intent: i100
repositories:
- api
---

# Invoice CSV export

## Outcome and repository scope

Callers in the api repository can write a list of invoices as CSV. The work stays in `internal/invoice`, beside the existing plain-text summary. Voiding a draft and crediting an issued invoice stay as they are. No format besides CSV and the plain-text summary is added, and the module takes no third-party dependency.

## Approach

`TextSummary` in `internal/invoice/export.go` renders one invoice as one line: number, customer, total in cents from `Invoice.Total()`, and status. CSV belongs in the same file, using the standard library `encoding/csv` writer.

Add `WriteCSV(w io.Writer, invoices []Invoice) error`. It writes a header row `number,customer,status,total_cents`, then one row per invoice in the given order. Number and total are decimal integers. Customer and status are written through the CSV writer so commas and quotes in a customer name stay in one field. An empty list writes only the header. The function does not filter by status and does not change the invoice.

## Tasks and order

1. Add `WriteCSV` in `internal/invoice/export.go`.
2. Add tests in `internal/invoice/export_test.go`: one invoice, several invoices, an empty list, and a customer name that contains a comma.
3. Run `go test ./...` in the api repository.

## Risks and checks

- A hand-rolled join would split a customer name that contains a comma. Use `encoding/csv`.
- Total must be `Invoice.Total()`, the sum of line amounts, not a field on the struct.

Check: `go test ./...` in the api repository.
