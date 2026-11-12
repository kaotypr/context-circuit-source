# Invoices

An invoice is voided only while it is a draft; once issued it is credited with a credit note. The same invoice can be written as a one-line plain-text summary or as CSV. The CSV has a header of number, customer, status, and cents, and one row per invoice.

Owner:

- `api@internal/invoice/` — voiding, credit notes, the plain-text summary, and the CSV export.
