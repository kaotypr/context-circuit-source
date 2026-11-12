I've written the intent for customer invoice CSV download as **i100**. It is not approved yet, and nothing is planned or built.

The record is at `intent/i100-customer-invoice-csv-export.md`. It is written in Indonesian, because the active member Maya's language is set to Indonesian in `members.yaml`. If you can't read Indonesian, tell me and I'll give you an English reading before you decide.

**Goal:** Customers can download their invoices as CSV through the API. Today invoices export only as a plain-text summary, so customers copy data by hand to get it into a spreadsheet.

**Non-goals:**
- No change to void or credit-note rules. Drafts can still be voided, and issued invoices are still corrected with a credit note.
- The plain-text summary stays as it is.
- No formats beyond CSV, such as PDF or Excel.

**Constraints:**
- A customer can download only their own invoices. Requests for another customer's invoice are refused. I wrote this as a constraint rather than a question. Tell me if you want it changed.
- Changes stay in the `api` repository, under `internal/invoice/`.

**Success criteria:**
- A customer can request a CSV of their invoice through the API and open the result in a spreadsheet.
- The file's contents and layout match your answers to questions 1–4.
- A request for another customer's invoice is refused and returns no data.
- The plain-text summary behaves as it did before.

**Open questions** (answer by number):

1. Should one request download one invoice, or a list of all the customer's invoices with a date filter?
2. What goes in the file: one row per invoice (number, date, status, total), or one row per line item on each invoice?
3. Should draft invoices and invoices already corrected with a credit note be included?
4. Which CSV conventions: comma or semicolon as the column separator, and what number and date formats? Indonesia usually uses a comma as the decimal separator, which can make amount columns read wrongly in a spreadsheet.

When you approve, give your answers to the questions. I'll record your decision in the record's language. Planning comes next, and I won't start it until you ask. The new record file is uncommitted, and I'll ask about committing once it is approved.