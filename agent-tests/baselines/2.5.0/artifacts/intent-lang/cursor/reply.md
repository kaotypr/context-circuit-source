I'll write this as an intent you can approve before any build starts. First I'm reading the intent skill and the workspace setup so the record matches how this project works.**i100** is written for your approval. Nothing gets built until you approve it.

Customers can download their own invoices from the API as CSV. Void-versus-credit-note rules stay as they are, and no other download format or surface is in scope. The work sits in `api`.

The invoice note already says invoices export only as a plain-text summary. CSV cuts across that, so this needs your answer before planning:

1. **Does CSV replace the plain-text summary, or sit beside it so customers can still download both?**

Reply with the approval in your own words, and answer question 1 by number. Planning starts only when you ask for it.

The Context Circuit plugin here is 0.2.4 (0.2.2 is installed). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.