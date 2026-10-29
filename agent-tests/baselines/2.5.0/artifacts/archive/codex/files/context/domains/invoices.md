# Invoices

An invoice is voided only while it is a draft; once issued it is credited with a credit note. Invoices can be exported as a plain-text summary or as CSV rows.

## What does CSV export contain?

CSV export writes a header and one row per invoice. Each row contains the invoice number, customer, status, and total in cents.

Owner:

- `api@internal/invoice/` — voiding, credit notes, and the plain-text summary.
