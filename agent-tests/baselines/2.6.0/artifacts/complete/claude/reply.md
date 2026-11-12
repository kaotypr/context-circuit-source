p1000 (invoice CSV export) is marked complete, and its knowledge is reconciled.

- **Completion:** `record complete` stamped `completed_at` on `plans/p1000-invoice-csv-export/plan.md` and appended a note stating that you reported the merge into `main` in the api repo and asked for completion.
- **Knowledge reconciled:** The catalog offered two candidates. The `invoices.md` note said invoices "export only as a plain-text summary," which the merged code makes false. I confirmed on the api repo's `main` that `internal/invoice/export.go` has the CSV export. I changed that sentence to cover both exports and updated the catalog entry's terms and reviewed date to 2026-11-11. `check` now passes.
- **Not touched:** The other candidate, `domains/billing/invoice-lifecycle.md`, does not exist in this workspace and is not in the catalog, so there was nothing to judge. Worth checking whether the CLI's candidate list is pointing at a different workspace.

Three files are uncommitted: `plans/p1000-invoice-csv-export/plan.md`, `context/domains/invoices.md`, and `context/INDEX.md`. Should I commit and push them, or commit only? Nothing has been staged or pushed yet.

Two other finished records of yours are not archived. I can move them to `_archive/` if you want, keeping their IDs.