---
product: both
kind: added
upgrade: Needs a CLI with `repo sync` and `workspace sync`; install the version the template pins.
---

Type `/ff` to bring your local base branches up to date: the agent fast-forwards every bound repository and the workspace, or only the ones you name (`/ff api`, `/ff workspace`), and reports what moved and what it left alone and why. The workspace is skipped when it holds any uncommitted file, and the code repositories are still updated. When you tell the agent your merge or pull request was merged, it now offers in one reply to update those repositories' base branches and to mark the delivered plans completed; one yes does both, and saying "merged" alone changes nothing. In Codex, ask in words instead of typing `/ff`.
