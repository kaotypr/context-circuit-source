I'll write this as an intent you can approve before any building starts. First I'm reading how intents are recorded here and checking who is working.**i100 — Download invoices as CSV** is written for your approval. Nothing has been built.

**Goal.** A customer can download their invoices from the API as a CSV file.

**Non-goals.** No change to when an invoice is voided or credited, and no download surface outside the API.

**Constraint.** Recorded knowledge says invoices export only as a plain-text summary. That rule and a CSV download cannot both stand as written.

**Success.** A customer requests their invoices from the API and receives a CSV. The plain-text summary stays or goes according to your answer below.

**Repository.** `api` only.

**Open questions**

1. **Should CSV replace the plain-text summary, or sit beside it?**
2. **Which invoices can a customer download: issued invoices only, or drafts, voided drafts, and credit notes as well?**

Reply with the answers by number, and say whether you approve this outcome. Approval is the next step. Planning starts only after you ask for it.

Context Circuit plugin 0.2.4 is available (installed 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.