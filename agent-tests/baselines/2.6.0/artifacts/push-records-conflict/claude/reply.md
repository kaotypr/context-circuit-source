The approval is committed locally but not pushed. The remote has a newer commit that changes the intent I just approved, so I stopped before rebasing.

**What's stuck:** The push was refused because `origin/main` has a new commit, `77d72bd` "docs: ask whether voided invoices are exported" (author Agent Test). It changes the intent's Open questions from "None known." to:

> 1. Do voided invoices belong in the export?

You approved the intent with "None known." on it, so that approval doesn't cover this question. The question bears on the outcome. The workspace rules say a question that bears on the outcome goes back to the intent, and a person decides it. A `pull --rebase` would put your approval on top of text you didn't approve, so I left it alone. Nothing was rebased, force-pushed, or edited.

**Done and verified:**
- `record approve --id i100` ran and exited 0. It added `approved_at: 2026-11-11T07:27:43Z` to the frontmatter and appended an approval section stating "Maya approved the intent as written, after reading it."
- The change was committed locally as `faa1ccb chore: approve intent i100 for invoice CSV export`. It is the only commit ahead of `origin/main`, and the working tree is clean.

**Current state:** `origin/main` does not have the approval yet, so other clones can't see it.

**Your choice:**
1. **Read the new question and re-approve.** If you answer it, the intent should be updated with your answer, and then I'd approve that version and push.
2. **Approve as it stands and push.** I'd pull with rebase, which should combine the two changes, and push. The approval would then cover the intent with the voided-invoices question still open, and I'd flag that in the approval note.
3. **Hold.** Leave the local commit unpushed while you decide.

I haven't planned i100. Planning is the next step after approval, and I can start it once you've settled the question.