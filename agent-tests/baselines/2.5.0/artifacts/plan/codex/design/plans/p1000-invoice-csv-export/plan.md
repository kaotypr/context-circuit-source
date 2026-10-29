---
id: p1000
created_by: maya
created_at: 2026-10-29T06:09:52Z
intent: i100
repositories:
- api
---

# Invoice CSV export

## Outcome and repository scope

Add a CSV export in `api/internal/invoice` so callers can write a collection of
invoices with the invoice number, customer, status, and total in cents. The
existing plain-text summary remains available, and invoice voiding and crediting
behavior does not change. All implementation and tests stay in the `api`
repository; no third-party dependency is added.

## Approach

Add an export function beside `TextSummary` that accepts an `io.Writer` and a
slice of invoices, then writes a header followed by one record per invoice.
Use Go's standard `encoding/csv` package so commas, quotes, and newlines in
customer names are escaped correctly. Format the invoice number and total in
cents as decimal values, flush the CSV writer, and propagate write errors to the
caller. Empty input produces the header only. Keep the export independent of
invoice lifecycle operations and do not change the existing summary API.

## Tasks and order

1. Add the CSV writer API in `internal/invoice/export.go`, emitting columns for
   invoice number, customer, status, and total cents in that order.
2. Add focused tests for the header and row values, multiple invoices, CSV
   escaping, empty input, and writer errors, including errors surfaced on flush.
3. Run `go test ./...` from the `api` repository.

## Risks and checks

The output column order and meaning must remain stable for callers. Using
`encoding/csv` avoids hand-written escaping; tests should include a customer
value containing a comma or quote. The success check is `go test ./...` from the
`api` repository.

