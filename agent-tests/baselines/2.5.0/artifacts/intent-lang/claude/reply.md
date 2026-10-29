I wrote the intent for CSV invoice download and stopped there. Nothing is approved and nothing is built.

The record is `intent/i100-invoice-csv-download.md`. Its body is in Indonesian because that is Maya's recorded language. The headings stay in English. It is not committed.

**Intent: customers can download their invoices as CSV from the API**

- **Goal:** Customers get their invoices as a CSV file directly from the API. Today invoices export only as a plain-text summary, so customers copy data by hand to get it into a spreadsheet or accounting tool.
- **Non-goals:** Changing or removing the plain-text summary. Bulk export for internal staff. Sending CSV by email or to third parties. Changing the void and credit-note rules.
- **Constraints:**
  - An invoice can be voided only while it is a draft. Once issued, it is corrected with a credit note. The CSV must show each invoice's real status.
  - A customer can download only their own invoices. No customer may see another customer's data.
  - The export code lives in `api@internal/invoice/`, so the change stays in the `api` repository.
- **Success criteria:**
  - An authenticated customer downloads a CSV with all their invoices and no one else's.
  - Each invoice shows its correct status: draft, issued, or credited.
  - The file opens cleanly in a spreadsheet, with UTF-8 encoding, agreed columns, and agreed number and date formats.
- **Repository scope:** `api` only.

**Open questions for you to answer by number:**

1. **Relation to the plain-text summary:** Does CSV sit alongside the existing summary, or replace it?
2. **Row granularity:** One row per invoice, or one row per line item?
3. **Which invoices:** Do customers download draft invoices, or only issued ones? Do credit notes appear as their own rows?
4. **Columns:** Which columns do you want? My suggested minimum is invoice number, issue date, due date, status, currency, subtotal, tax, and total.
5. **Range and size:** Is one download of all invoices enough, or do you need a date filter? If a customer has many invoices, do you need pagination or an asynchronous export?
6. **Number and date formats:** Which decimal and thousands separators? Is the currency written on every row? Are dates ISO 8601 (`YYYY-MM-DD`)?
7. **Customer authentication:** How do customers authenticate to the API today? Is there already a per-customer token or API key this endpoint can use?

Once you answer these, I'll record your approval in your own words. Planning is the next step after that, and I won't start it until you ask.