---
id: i100
created_by: maya
created_at: 2026-11-11T07:22:08Z
---

# Customers can download their invoices as CSV from the API

## Goal

A customer can download their own invoices as a CSV file through the API, so they can take them into their own spreadsheets and accounting tools. Today the invoice export is only a plain-text summary, which customers cannot load into those tools.

## Non-goals

- Changing invoice issuing, voiding, or credit-note behavior.
- A CSV export in the web UI. This intent covers the API only.
- Exporting other customers' data, or any admin-wide export.
- Replacing the existing plain-text summary, unless question 5 is answered that way.

## Constraints

- A customer sees only their own invoices. The access rules the API already applies to invoice reads are the ones to follow; confirm them in the code after approval.
- The CSV should be built from the same invoice data the plain-text summary uses, so the two cannot disagree about amounts or statuses. Confirm this in the code after approval.
- The `context/domains/invoices.md` note records that invoices export only as plain text. That note becomes wrong once this ships, and it is updated as part of the work, not left to drift.

## Success criteria

- A customer calls the API and receives a CSV file with a header row and one row per invoice in the selected range.
- Every column has a value that matches the invoice as the API returns it elsewhere. Amounts and dates appear in one stated format.
- A customer who requests the export cannot receive an invoice that belongs to another customer. This is checked by a test.
- The export runs with no invoices, one invoice, and a large number of invoices, and each returns a well-formed file.

## Repository scope

- `api` — the invoice export endpoint and the CSV writer, alongside the existing invoice summary in `internal/invoice/`.

## Open questions

1. **Which invoices are included?** Options: (a) all of the customer's invoices, (b) only issued invoices, leaving out drafts, (c) a date-range filter such as `from` and `to`. Recommended: (b) plus an optional date range, so drafts never reach a customer's accounting.

2. **How are credit notes shown?** Options: (a) a credit note is its own row with a negative amount, (b) credit notes are left out, (c) the credited invoice's row shows its credited amount. Recommended: (a), because it matches how the invoice lifecycle already records credits.

3. **What is one row?** Options: (a) one row per invoice, (b) one row per line item, with invoice fields repeated on each row. Recommended: (a) for the first version. Line items can come later as a separate export.

4. **What is the endpoint shape?** Options: (a) a new endpoint such as `GET /invoices/export.csv`, (b) the existing invoice list endpoint, selected by a format parameter or an `Accept: text/csv` header. Recommended: (a), so the existing endpoints keep their contract.

5. **Does the plain-text summary stay?** Options: (a) keep it alongside the CSV, (b) replace it. Recommended: (a), because the request adds a format and removes nothing. Confirm this, because replacing it would be a change customers notice.

6. **Is the export returned directly, or built in the background?** Options: (a) returned directly in the response, (b) a background job that the customer collects later. Recommended: (a) for now, with a row cap and an error that says to narrow the date range when the cap is hit. The cap value is for you to set.

7. **What format do amounts, dates, and encoding take?** Options: amounts with a decimal point or with the currency code in its own column, dates as `YYYY-MM-DD` or ISO 8601 with a time zone, and UTF-8 with or without a byte-order mark for spreadsheet compatibility. Recommended: decimal amounts with a separate currency column, `YYYY-MM-DD` dates, and UTF-8 with a byte-order mark so Excel opens the file without garbled characters.
