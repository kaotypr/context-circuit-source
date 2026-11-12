---
id: p1000
created_by: maya
created_at: 2026-11-12T09:46:42Z
intent: i100
repositories:
- api
---

# Invoice CSV export

## Outcome and repository scope

Callers can write a collection of invoices as CSV with the invoice number, customer, status, and total in cents. The change stays in `api@internal/invoice` and uses only the Go standard library.

## Approach

Add a writer alongside `TextSummary` that accepts an `io.Writer` and the invoices to export. Use `encoding/csv` so customer values with commas, quotes, or newlines are escaped correctly; write a header and one row per invoice, deriving the total from `Invoice.Total()`. Propagate write and flush errors to the caller. Keep the existing summary and invoice lifecycle behavior unchanged.

## Tasks and order

1. Implement the CSV writer in `internal/invoice`, with columns for invoice number, customer, status, and total cents.
2. Add package tests that parse the output and verify the header, values, totals, multiple invoices, and CSV escaping for customer names.
3. Run `go test ./...` in the API repository.

## Risks and checks

CSV consumers rely on a stable header and column order, so tests should assert both explicitly. `Invoice.Total()` returns an `int64` in cents; format it as a base-10 integer without currency conversion. Check the writer's flush error as well as immediate write errors.

Run `go test ./...` in the API repository.

## Details

No supporting files.
