## Verdict

Mixed: one behavior is kept; one is unclear.

## Findings

- **Plain-text summary still lists every invoice — unclear.** `internal/invoice/export.go` leaves `TextSummary` unchanged from `main`, but it formats a single invoice. I found no list-level summary to confirm that it lists every invoice.
- **Totals stay in cents — kept.** `internal/invoice/export.go` writes `item.Total()` as an integer in the CSV row. `internal/invoice/export_test.go` checks that an invoice with a 100-cent line produces `100`.

## Not checked

None.