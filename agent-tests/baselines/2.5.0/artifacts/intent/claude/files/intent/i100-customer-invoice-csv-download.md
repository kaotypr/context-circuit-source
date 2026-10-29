---
id: i100
created_by: maya
created_at: 2026-10-29T06:01:18Z
---

# Customers can download their invoices as CSV from the API

## Goal

A customer of the invoicing API can download their own invoices as a CSV file, so they can load them into their spreadsheets and accounting tools without parsing the existing plain-text summary.

The intended reader is the customer's finance or accounting side, who needs one row per invoice (or per line item, see question 3) with stable column names.

## Non-goals

- No CSV import, and no change to how invoices are created, issued, voided, or credited.
- No change to the web UI.
- No bulk export across customers. Each request returns only the calling customer's invoices.
- No scheduled or emailed exports.

## Constraints

- Invoices are voided only while they are drafts; once issued they are credited with a credit note (`context/domains/invoices.md`). The CSV must reflect that: a credited invoice is not presented as if it were unchanged.
- Today invoices export only as a plain-text summary (`context/domains/invoices.md`). This intent adds a CSV form; whether the plain-text summary stays is question 4.
- A customer must never receive another customer's invoice data. The access rule follows the API's existing customer authentication, which this intent does not change.
- Money values, dates, and identifiers must be stable across downloads of the same invoice, so a re-download matches the first one.

## Success criteria

- An authenticated customer requests their invoices in CSV and receives a file that opens in a spreadsheet with the agreed columns from question 5.
- A request for an invoice belonging to another customer returns the same refusal an unknown invoice returns, with no data.
- A credited invoice appears with its credit note reflected in the agreed form (question 2).
- The existing plain-text summary behaves exactly as before, or changes only as question 4 decides.
- The `api` test suite passes with new tests covering the CSV output and the access rule.

## Repository scope

- `api` — invoice export under `api@internal/invoice/`, which already owns the plain-text summary and voiding (`context/domains/invoices.md`). The new CSV form, its endpoint or format parameter, and its tests are expected to live here.
- No other repository is expected to change. If the work reaches another repository, that will be shown before it is planned.

## Open questions

1. **Which invoices are included?**
   Proposed: issued invoices and their credit notes only, excluding drafts. Drafts are not final, but a customer may want to see them. Include drafts, or exclude them?

2. **How is a credited invoice shown?**
   Proposed: the original invoice row stays, and each credit note is its own row with a negative total, so the totals reconcile. Or should the credit be shown only as a status column on the original?

3. **What is one row?**
   Proposed: one row per invoice, with totals. Or one row per line item, with the invoice fields repeated on each row?

4. **Does the plain-text summary stay?**
   Proposed: keep it unchanged and add CSV beside it. Or should CSV replace it?

5. **Which columns, and in what format?**
   Proposed: invoice number, issue date (`YYYY-MM-DD`), status, currency, subtotal, tax, total, and customer reference, in UTF-8 with a header row, `.` decimal separator, and no thousands separator. Confirm or change the column list and formats.

6. **How does a customer ask for CSV?**
   Proposed: a format option on the existing invoice listing (for example `format=csv`), rather than a new endpoint. Or a separate `/invoices.csv` endpoint?

7. **Is there a date range or a row limit?**
   Proposed: a customer may pass a start and end issue date, and the export returns all matching rows in one file with no pagination. Or should large exports be paginated or capped at a size?

8. **Does the access rule need anything beyond the existing customer authentication?**
   Proposed: none. The CSV uses the same authentication and customer scoping as the invoice listing. Confirm that no separate permission or API key is intended.
