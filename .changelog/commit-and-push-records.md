---
product: template
kind: changed
pr: 209
---

Where `record approve`, `record complete`, or `record archive` return
`commit_recommended`, the agent now asks one question: commit and push the
records and ID files it names, or commit only. The person's answer authorizes
that commit and a push of the workspace's base branch to its upstream, that
once; after commit only, the agent says other clones cannot see the records yet,
and with no remote or upstream it commits and says why it could not push. A push
refused because the remote moved on is brought up to date with `git pull
--rebase`, which replays only the person's unpushed record commits: what Git
combines is kept and reported, an ID file conflict keeps the union of both
sides' entries, and any other conflict — the same record changed on both sides,
a record deleted on one side and changed on the other, or a conflict in a file
that is not a record — goes to the person with both versions, the rebase left in
progress. The agent never edits a conflicted record itself, not even to keep
both sides: which version stands is the person's call. `check` follows; a number
reserved twice goes to the person. The agent pushes once more, reports a second
refusal or a refusal for any other reason, and never force-pushes. Teammates' ID
file conflicts from upgrading at the same time are resolved the same way, by the
union and `check`. Delivering code is unchanged: drift there is still explained,
not rebased.
