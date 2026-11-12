I'll write this as an intent you can approve before any build starts. First I'm reading how this workspace records intents and what it already says about invoices.Knowledge says invoices only export as a plain-text summary, so the CSV download changes that. I'll check workspace setup, then write the intent.maya, customers who can already see their invoices will be able to download them as CSV from the API. I wrote that outcome as **Customer invoice CSV download** (i100) for you to approve before anyone builds it.

Today an invoice only exports as a plain-text summary. Void and credit stay as they are: a draft is voided, and an issued invoice is credited with a credit note. The CSV is a download of those same facts, one row per invoice, in the `api` repository. The plain-text summary stays unless you say otherwise.

Four things still change what we build. My recommendations:

1. Keep the plain-text summary beside the CSV, so existing callers do not have to switch.
2. One row per invoice, matching the summary. One row per line item is a larger export than this request describes.
3. Include only issued invoices. Drafts are still voidable and are not a customer's bill yet.
4. Anyone who can already see the invoice can download it, not only the account owner.

Reply with the numbers you want changed, or say they are all fine. Once you approve this, tell me and I will write the plan. Nothing gets built until then.

The Context Circuit plugin here is 0.2.2; 0.5.0 is available. Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.