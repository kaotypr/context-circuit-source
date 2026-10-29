p1000 is marked complete, and its knowledge reconciliation is done. Nothing is committed or pushed yet.

**What was recorded**
- `record complete --id p1000` stamped `completed_at` and appended your words as the note: "p1000 has been merged into main in the api repo. Please mark it complete."
- I confirmed in the api repo that main has the CSV export (`internal/invoice/export.go`, merge commit `7e967a3`).
- The reconciliation candidate `context/domains/invoices.md` said invoices export only as plain-text. That is no longer true, so I updated it to cover the CSV export and its columns. I also updated its catalog entry in `context/INDEX.md` (new terms `csv, export`, reviewed date 2026-10-29). The void/credit rule is unchanged.
- `check` returns `ok: true` with no issues.

**Needs your decision**
- **Commit and push?** The uncommitted files are `plans/p1000-invoice-csv-export/plan.md`, `context/domains/invoices.md`, and `context/INDEX.md`. Should I commit and push them, or commit only? Other clones can't see the completion until they're pushed.
- **Archive?** There are 2 other finished records of yours not archived. Should I archive them with `record archive`? I haven't run it.

**Noticed, not changed**
- `context/INDEX.md` lists an "Invoice lifecycle" note at `domains/billing/invoice-lifecycle.md` (twice, lines 14 and 109), but that file doesn't exist. `check` didn't flag it, so the broken entry may need a look.