---
id: p1000
created_by: maya
created_at: 2026-11-11T16:45:14Z
intent: i100
repositories:
  - api
---

# Invoice CSV export

## Outcome and repository scope

A caller in the api repository can write a list of invoices as CSV, with one row per invoice and columns for the invoice number, the customer, the status, and the total in cents. The change stays in `internal/invoice`, beside the plain-text summary. Voiding a draft and crediting an issued invoice stay as they are, and the module keeps using only the standard library.

## Approach

`TextSummary` in `internal/invoice/export.go` already prints one invoice as a single line from `Number`, `Customer`, `Status`, and `Total()`. Add a function next to it that takes a slice of invoices and an `io.Writer`, and writes them with `encoding/csv`. That package is in the standard library, so it meets the constraint against third-party dependencies.

The first row is a header: `number`, `customer`, `status`, `total_cents`. Each following row is one invoice. The total is the existing `Total()` sum of line amounts, written as a decimal integer. `encoding/csv` quotes a customer name that contains a comma or a quote. An empty slice writes the header and no data rows. The function returns the error from the writer.

`Void` and `Credit` in `internal/invoice/invoice.go` are left unchanged.

## Tasks and order

1. Add the CSV writer in `internal/invoice/export.go`, next to `TextSummary`.
2. Add tests in `internal/invoice/export_test.go` for a couple of invoices, a customer name that needs quoting, and an empty slice.
3. Run `go test ./...` in the api repository.

## Risks and checks

A customer name with a comma would split a hand-built row into extra columns, so the writer has to go through `encoding/csv` rather than `fmt`. The total has to be `Total()`, because the struct doesn't store a total of its own.

`go test ./...` in the api repository is the check the intent names, and it's the repository's own rule before a commit.
