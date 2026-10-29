I did not change anything in this workspace, because the fix belongs upstream rather than in our `context/`.

**Where the claim lives:** the rule "invoice numbers are global and never restart" comes from the `handbook` knowledge repository, in `numbering.md` and its `INDEX.md` entry. The handbook is mounted in `workspace.yaml` under `knowledge_repositories`, so it is borrowed and read-only here. Our `context/` has no numbering note, and nothing local repeats the claim. The local `context/domains/invoices.md` covers only voiding, credit notes, and the plain-text summary.

**Why I didn't record a correction here:** `AGENTS.md` says a correction to borrowed knowledge is never recorded in `context/`, and the handbook itself must not be edited, cloned, or pushed unless you ask for that as its own task. A local "numbers restart yearly" note would be exactly that kind of correction, so I left `context/` alone. I also did not commit or push anything.

**Evidence check:** your `api` checkout has `Number` as a field on `Invoice` (in `internal/invoice/invoice.go` and `export.go`), but no code I found shows how numbers are assigned, so the code neither confirms nor contradicts the rule. The fix needs your knowledge of the yearly reset, or the code that assigns numbers.

**What I'd suggest next:**
1. Send the handbook's maintainers a change to `numbering.md` (and the `INDEX.md` line) saying numbers restart each year, through that repository's own review. I can draft that text for you to submit.
2. Until it lands, anyone reading the handbook will still get "global". If you want a local warning that doesn't count as a correction, tell me the wording and I'll record it as a plain project note.
3. If you want the handbook unmounted from this workspace in the meantime, say so and I'll make that change as its own task.

Which of these would you like?