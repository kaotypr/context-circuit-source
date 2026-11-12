# Invoices

An invoice is voided only while it is a draft; once issued it is credited with a credit note. Invoices can be exported as a plain-text summary or as CSV.

## What does the CSV export contain?

The CSV has one row per invoice, with its number, customer, status, and total in cents. It includes a header row and leaves the plain-text summary available.

Owner:

- `api@internal/invoice/` — voiding, credit notes, and invoice exports.
