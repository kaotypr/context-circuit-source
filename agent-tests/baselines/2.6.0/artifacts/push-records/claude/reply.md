i100 is approved, and the records are committed and pushed to `origin/main` as `1aa4c5a`.

- **Approval:** recorded with the note "User approved the intent as written: customers can download their invoices as CSV from the api, in internal/invoice."
- **Push:** the first push was refused because the remote had one new commit, `a111c09 chore: plan invoice numbering`. It adds plan `p1000` and also changes `.context-circuit/ids/maya.yaml`.
- **Conflict:** the only conflict was in `maya.yaml`. I kept both sides' entries, so the file now reserves intent `i100` and plan `p1000`. I did not touch the plan file that came in.
- **Check:** `check` reports no issues.
- **Files:** only `.context-circuit/ids/maya.yaml` and `intent/i100-invoice-csv-export.md` were committed. Nothing else was staged.

Planning has not started. The next step is creating the plan for i100 when you ask for it.