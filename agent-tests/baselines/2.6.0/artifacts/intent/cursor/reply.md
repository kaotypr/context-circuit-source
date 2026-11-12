I'll write this as an intent you can approve before any build starts. First I'll read the intent skill and the workspace setup so the record matches how this project works.I've written it up as an invoice CSV download (i100). A customer can call the API and get their own invoices as a CSV a spreadsheet can open. Creating, issuing, voiding, and crediting invoices stay as they are, and there is no new screen. Today the only export is a plain-text summary, so the CSV has to sit next to that rule rather than invent another way to close an invoice.

Three things need your call. Here's what I'd pick, so tell me if you'd go another way:

1. Keep the plain-text summary and add the CSV beside it. Replacing the summary would break anyone who already downloads it.
2. One row per invoice, with the same facts the summary already shows. A row per line item is a different report, and this request is a download of the invoices.
3. Issued invoices only, including ones that already have a credit note. Drafts and voided drafts aren't something a customer is billed for.

If those work for you, say so and I'll record the approval. The plan isn't written yet. Once you've approved, ask me and I'll write it.

The Context Circuit plugin in this session is 0.5.0, and 0.2.2 is what is installed. You can update it with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.