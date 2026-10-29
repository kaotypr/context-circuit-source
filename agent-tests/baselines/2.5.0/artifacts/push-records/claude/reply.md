i100 is approved, and the records are committed and pushed to `origin/main`.

- **Approval:** recorded for i100 with your words, "I approve it as written."
- **Commit:** `3217f99 chore(intent): approve i100 invoice CSV export`, containing only `intent/i100-invoice-csv-export.md` and `.context-circuit/ids/maya.yaml`.
- **Push conflict:** the first push was refused because the remote had a new commit, `5784162 chore: plan invoice numbering`. Rebasing hit a conflict in `.context-circuit/ids/maya.yaml`. I kept both sides' entries: `intents: [i100]` and `plans: [p1000]`. The rebase then completed and I pushed again.
- **Check:** `check` returned `ok: true` with no issues.

The incoming commit also changed `maya.yaml`, which reserves p1000 for `plans/p1000-invoice-numbering/plan.md`. Maya's ID file was edited from another clone, so it's worth confirming that p1000 is Maya's work and not someone else's.

Planning is the next step, and I haven't started it. Say the word and I'll create the plan(s) for i100.