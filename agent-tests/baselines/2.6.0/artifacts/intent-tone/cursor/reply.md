I'll write this as an intent you can approve before any build starts. First I'm checking how this workspace records intents and who is working.Knowledge says invoices only leave the API as a plain-text summary, so the intent adds a CSV download beside that. I'm writing it in your recorded language so you can approve the outcome itself.I've written it up as customer invoice CSV download (i100). A logged-in customer can download their own invoices as a CSV from the API, and the plain-text summary stays as it is. Void and credit rules don't change, and there's no new page.

Two quick calls. Here's what I'd go with:

1. One row per invoice, same shape as the summary. One row per line item suits accounting tools better, but it's a bigger file and more work.
2. Issued invoices only. Drafts are still voided, so they don't belong in a download the customer keeps.

All good? Say the word and I'll log the approval. The plan isn't written yet. Once this is approved, ask and I'll write it.

The Context Circuit plugin here is 0.5.0, with 0.2.2 installed. Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.