---
product: both
kind: added
upgrade: Nothing to do unless an extension you use declares scripts; each member allows its scripts on their own machine with `extension allow --id ID` before they run.
---

Extensions can now run scripts at their hook points, not only give the agent instructions. A script runs only after you ran `extension allow --id ID` on your machine for that exact version — the commit of a remote extension, or the contents of a folder in the workspace — and any change to it asks again; `allow` shows each script, what happens when it fails, and what it does beyond your machine, and runs nothing. `extension list` shows whether each extension is allowed, not allowed, or changed since you allowed it. Each hook declares whether a failure stops the stage or is only reported, and whether it needs your consent: `standing` scripts take the one action their manifest names once you have allowed them, and `each-time` scripts wait for your yes every time. `record approve`, `worktree prepare`, `record complete`, and `record archive` run their own stage's scripts without holding the workspace lock, and a failing `before` script with `on_failure: stop` leaves the command undone. On Windows a hook runs its `run_windows` command, and a hook without one is skipped there.
