---
id: p1000
created_by: maya
created_at: 2026-11-11T07:26:33Z
intent: i100
repositories:
- api
---

# Invoice CSV export

## Outcome and repository scope

Invoices can be written as CSV, one row per invoice, with the invoice number, the customer, the status, and the total in cents. The work is in the `api` repository and is confined to `internal/invoice`, beside `TextSummary` in `export.go`. Nothing else in the repository changes, and no third-party module is added.

The intent's goal says customers can download their invoices. This repository has no HTTP handler or download route, so this plan delivers the CSV writer that such a route would call. The route is outside the approved scope and is not part of this plan (see the questions at the end).

## Approach

Add a `WriteCSV` function in a new file, `internal/invoice/csv.go`. It takes an `io.Writer` and a slice of `Invoice` values, and writes a header row followed by one row per invoice. Using `io.Writer` keeps the function testable without a file and lets a future download handler stream the output directly.

The standard library `encoding/csv` does the quoting, so a customer name containing a comma or a quote is written correctly. The total comes from the existing `Invoice.Total()`, so the cents value matches what `TextSummary` already reports. Columns are fixed in this order: `invoice_number`, `customer`, `status`, `total_cents`.

The function returns the error from the csv writer (`Flush` then `Error`) and does not swallow write failures.

## Tasks and order

1. Add `internal/invoice/csv.go` with `WriteCSV(w io.Writer, invoices []Invoice) error`, using `encoding/csv`.
2. Add `internal/invoice/csv_test.go` covering: the header row and column order; a customer name with a comma and a quote, quoted correctly; an invoice with no lines, which writes a total of `0`; an empty slice, which writes only the header; and a writer that fails, whose error is returned.
3. Run `go test ./...` in `repositories/api` and confirm it passes.

Task 1 must land before task 2. There are no dependencies on other plans.

## Risks and checks

- **Scope creep into a download route.** The temptation is to add an HTTP handler to make the goal literal. The constraint forbids it here, so the route stays out and is raised as a question.
- **Total drift.** If the CSV total were computed separately from `Total()`, it could disagree with the plain-text summary. Using `Total()` directly avoids that.
- **Header choice.** The intent does not say whether a header row is wanted. This plan includes one, because a CSV without column names is hard to read in a spreadsheet. A person can ask for it to be dropped.
- Checks: `go test ./...` in `repositories/api`. Also `gofmt -l internal/invoice` should print nothing, and `go vet ./...` should pass.
