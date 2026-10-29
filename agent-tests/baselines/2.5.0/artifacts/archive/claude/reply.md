p1000 is marked complete, and the knowledge it changed is reconciled. Nothing is committed or pushed yet.

**Completion:** `record complete` stamped `completed_at` on `plans/p1000-invoice-csv-export/plan.md` with your words as the note. I confirmed on the api repo's `main` that the CSV export (`internal/invoice/export.go`) is there.

**Knowledge reconciliation:**
- **Invoices** (`context/domains/invoices.md`): the sentence "Invoices export only as a plain-text summary" was no longer true. It now says invoices export as plain text or CSV. I also added `csv` and `export` to its catalog entry in `context/INDEX.md` and moved its reviewed date to 2026-10-29.
- **Invoice lifecycle** (`context/INDEX.md` line 109): p1000 didn't touch voiding or credit notes, so I left it alone. That entry points to `domains/billing/invoice-lifecycle.md`, which doesn't exist in `context/`. That was already broken before this change, so I didn't create or delete anything for it. Someone should decide whether to write that note or drop the entry.
- `check` reports `ok: true`.

**Decisions for you:**
1. **Commit and push?** The workspace records and the two knowledge edits are uncommitted. The CLI lists `plans/p1000…/plan.md`, `intent/_archive/i101…` and `plans/_archive/p1001…` as recommended for commit. It does not list `context/INDEX.md` or `context/domains/invoices.md`, so those would need to be added too. Do you want me to commit and push all of it, or commit only?
2. **Archive?** Two finished records of yours are unarchived. I can run `record archive` if you want.
3. **Folder rename:** the CLI renamed `intent/archive` and `plans/archive` to `_archive`. Anyone else who pulls needs CLI 2.5.0 or later first. The pinned CLI here is 2.5.0-rc.1, so it's worth checking that before the rename is pushed.