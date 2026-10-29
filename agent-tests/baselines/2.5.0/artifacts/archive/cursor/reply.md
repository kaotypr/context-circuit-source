I'll mark plan p1000 complete using the workspace completion process, starting with the completion skill and the plan record.The plan is marked complete. I'll judge the billing notes against what the CSV export actually changed before deciding whether knowledge needs an update.The invoices note still says export is plain text only, which the merged CSV export changes. I'll confirm that against the api checkout before updating the note.**p1000** is complete. The completion note records your words: it has been merged into main in the api repo.

The invoices note still said export was plain text only. The merged `CSV` writer in `api@internal/invoice/` adds a second export (number, customer, status, and total in cents), so that note and its catalog entry now say both forms exist. The reviewed date is 2026-10-29. The workspace check reports no issues.

Completing the plan also renamed `intent/archive` and `plans/archive` to `intent/_archive` and `plans/_archive`. Those moves are uncommitted, along with the plan record and the knowledge edits. Other clones will not see them until they are pushed, and anyone else should be on CLI 2.5.0 or later before they pull.

You have two finished records that are not archived. I can archive them if you want; otherwise they stay where they are.

Do you want these files committed and pushed, or committed only?