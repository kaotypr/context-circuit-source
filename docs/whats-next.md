# Whats next

`whats-next` is read-only. It ignores archived plans, considers only plans with `status: approved` and unfinished tasks, and requires every declared plan dependency to be done. The result recommends a whole plan, lists its remaining tasks, explains why it is ready, and includes the plan and Product Knowledge references.

If a source in `context/sources.yaml` has changed, the result warns that Product Knowledge may need a refresh. If nothing is executable, it recommends reviewing a draft plan or creating one. It never claims, publishes, or updates work.
