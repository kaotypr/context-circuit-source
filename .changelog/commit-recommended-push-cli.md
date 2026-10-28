---
product: cli
kind: changed
pr: 209
---

The `commit_recommended` note now tells the agent to ask the person whether to
commit and push the listed files, or commit only, and carries what to do when
that push is refused: pull with rebase, resolve only an ID file conflict by
keeping both sides' entries, stop with the rebase in progress for any other
conflict, even one where both sides could be kept, and never force-push. It
still says nothing was staged, committed or pushed; the CLI never stages,
commits, or pushes.
