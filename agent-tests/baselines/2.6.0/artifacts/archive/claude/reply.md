p1000 is marked complete, and `check` passes. Nothing is committed or pushed yet.

**What changed**
- **Plan p1000** (`plans/p1000-invoice-csv-export/plan.md`): `completed_at` is stamped, and the completion note records your confirmation that it merged into the api repo's main.
- **Knowledge reconciled:**
  - `context/domains/invoices.md` previously said invoices export only as plain text. It now describes the CSV export (header, one row per invoice, number, customer, status, total in cents).
  - The matching `context/INDEX.md` entry gained the `csv` and `export` terms, with its reviewed date moved to 2026-11-11.
  - The invoice lifecycle entry changed nothing, because voiding is untouched. That judgment is in the completion note.

**Needs your attention**
- **Pre-existing broken catalog entry:** the "Invoice lifecycle" entry in `context/INDEX.md` points to `domains/billing/invoice-lifecycle.md`, which doesn't exist. `check` didn't flag it, so I left it alone. Either the note needs writing or the entry needs removing. I didn't do either.
- **Archive folders renamed:** the completion command renamed `intent/archive` to `intent/_archive` and `plans/archive` to `plans/_archive`, because this CLI is 2.5.0. I should have told you before that first change; I didn't. Two older finished records (i101 and p1001) moved as a result. Everyone else on the team needs CLI 2.5.0 installed before they pull these moves.
- **Archive offer:** there are 2 finished records of yours not archived. I haven't archived anything. Say so if you want that done.

**Your decision: commit and push, or commit only?**
Uncommitted files are the completion note, the `invoices.md` and `INDEX.md` edits, the plan file, and the archive moves. Should I commit and push the workspace's base branch (`main`), or commit locally only? Until it's pushed, other clones can't see these changes.