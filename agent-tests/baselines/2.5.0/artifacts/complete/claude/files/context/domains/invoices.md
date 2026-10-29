# Invoices

An invoice is voided only while it is a draft; once issued it is credited with a credit note. Invoices export as a plain-text summary or as CSV, one row per invoice with its number, customer, status, and total in cents.

Owner:

- `api@internal/invoice/` — voiding, credit notes, the plain-text summary, and the CSV export.
