# Invoices

An invoice is voided only while it is a draft; once issued it is credited with a credit note. An invoice can be exported as a plain-text summary or as CSV.

## What does each export contain?

| Export | Contents |
| --- | --- |
| Plain-text summary | One line: number, customer, total in cents, and status |
| CSV | A header and one row per invoice: number, customer, status, and total in cents |

Owner:

- `api@internal/invoice/` — voiding, credit notes, the plain-text summary, and the CSV export.
