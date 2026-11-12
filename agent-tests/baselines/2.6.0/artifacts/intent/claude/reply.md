The intent is written as **i100**, `intent/i100-customer-invoice-csv-export.md`. It is not approved, and nothing is being built yet.

**Goal:** A customer can download their own invoices as a CSV through the API. Today the invoice export is only a plain-text summary.

**Non-goals:** No changes to issuing, voiding, or credit-note behavior. No web UI export. No admin-wide export. The plain-text summary stays unless you decide otherwise (question 5).

**Constraints:**
- Customers see only their own invoices. I haven't checked the existing access rules in the code yet, so that is to be confirmed after approval.
- The CSV should come from the same data as the plain-text summary, so they can't disagree. Also to be confirmed in the code.
- `context/domains/invoices.md` says invoices export only as plain text. That note will be wrong once this ships, and I'll update it as part of the work.

**Success criteria:**
- The API returns a CSV with a header row and one row per invoice in the selected range.
- Values match what the API returns for the same invoice elsewhere.
- A test confirms a customer cannot receive another customer's invoice.
- Works with zero, one, and many invoices.

**Repository scope:** `api`, in `internal/invoice/`, alongside the existing summary.

**Open questions** (I've put a recommendation on each):

1. **Which invoices?** Recommended: issued invoices only, with an optional `from`/`to` date range. Drafts stay out.
2. **Credit notes?** Recommended: each is its own row with a negative amount, matching how credits are already recorded.
3. **One row per what?** Recommended: one row per invoice for now. Line items can be a separate export later.
4. **Endpoint shape?** Recommended: a new `GET /invoices/export.csv`, so existing endpoints don't change.
5. **Keep the plain-text summary?** Recommended: yes, keep it alongside the CSV.
6. **Direct response or background job?** Recommended: direct response with a row cap. You set the cap value, and an over-cap request asks the customer to narrow the date range.
7. **Formats?** Recommended: decimal amounts with a separate currency column, `YYYY-MM-DD` dates, and UTF-8 with a byte-order mark so Excel opens it cleanly.

Please reply with an answer to each question by number, or "accept the recommendations." Once you approve, I'll record the decision. Planning comes next, and I'll wait for you to ask for it. The intent and its ID file are uncommitted. I'll ask whether to commit them when the approval is recorded.