I wrote up the intent as record `i100` for you to review. Nothing is approved yet, and nothing has been built or planned. The file is `intent/i100-customer-invoice-csv-download.md`, and it is uncommitted.

**Customers can download their invoices as CSV from the API**

**Goal:** A customer of the invoicing API can download their own invoices as a CSV file, for use in spreadsheets and accounting tools.

**Non-goals:** No CSV import, no change to invoice creation, voiding, or credit notes, no web UI change, no cross-customer bulk export, no scheduled or emailed exports.

**Constraints:**
- Drafts can be voided. Issued invoices are credited with a credit note, not voided, so the CSV must show that.
- Invoices currently export only as a plain-text summary. This adds a CSV form next to it.
- Customers see only their own invoices, under the API's existing authentication.
- The same invoice must give the same values on every download.

**Success criteria:**
- An authenticated customer gets a CSV with the agreed columns that opens in a spreadsheet.
- Another customer's invoice gets the same refusal as an unknown invoice, with no data.
- A credited invoice shows its credit note in the agreed form.
- The plain-text summary is unchanged, unless question 4 changes it.
- The `api` test suite passes, with new tests for CSV output and access.

**Repository scope:** `api`, under `api@internal/invoice/`, which already owns the plain-text summary. No other repository is expected to change.

**Open questions** (answer by number; my proposal is in each one):

1. **Which invoices are included?** Proposed: issued invoices and their credit notes, no drafts. Or include drafts?
2. **How is a credited invoice shown?** Proposed: the original row stays and each credit note is its own row with a negative total, so totals reconcile. Or show the credit only as a status column?
3. **What is one row?** Proposed: one row per invoice with totals. Or one row per line item?
4. **Does the plain-text summary stay?** Proposed: yes, CSV is added beside it. Or should CSV replace it?
5. **Which columns and formats?** Proposed: invoice number, issue date (`YYYY-MM-DD`), status, currency, subtotal, tax, total, customer reference. UTF-8 with a header row, `.` as the decimal separator.
6. **How does a customer ask for CSV?** Proposed: a `format=csv` option on the existing invoice listing. Or a separate `/invoices.csv` endpoint?
7. **Date range or row limit?** Proposed: optional start and end issue dates, one file, no pagination. Or paginate or cap large exports?
8. **Access rule beyond current authentication?** Proposed: none. Confirm that no separate permission or API key is intended.

Question 4 needs a decision from you, not a guess, because the existing invoice note says invoices export only as a plain-text summary. Once your answers are in, I'll update the intent with them. Approval comes after that. Your approval should be in your own words, and I'll record them as given. Planning only starts when you ask for it.