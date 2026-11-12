I wrote the intent as **i100, "Export invoice ke CSV lewat API"**. Nothing is approved, planned, or built yet, and the record is uncommitted. It's in Indonesian because Maya, the only member on the roster, writes in Indonesian. I can give you an English copy if you need to read it directly.

**What it proposes:** Customers can already get a plain-text summary of each invoice. This adds a CSV download of their own invoices through the API. Only the API endpoint is in scope, not a web UI button. Amounts stay in cents, like the rest of the API.

I'd like answers to four questions. Each has my recommendation:

1. **One row per line item, or one per invoice?** I'd go per line item, so customers get the detail they need for their books. Invoice totals still add up from the lines. Per invoice is simpler but hides what was sold.
2. **Date range and limit?** Customers pick a date range, and one request is capped at twelve months so large accounts don't time out. Twelve months, or a different limit?
3. **Drafts and credit notes?** I'd leave drafts out, since they were never issued. Credit notes appear as their own rows pointing at the invoice they credit, so the CSV matches what the account shows.
4. **Who can export?** I'd allow anyone who can already view invoices on that account, not just the owner. Or should it be owner-only?

If you don't answer, nothing moves. Planning and building wait for approval.

I only read the invoice package in the `api` repo to ground this. I didn't check the routes or auth code, so the endpoint details will be settled in the plan after you approve.

To approve, reply with your answers to the four questions and say yes. I'll record the approval and then ask whether to commit it.