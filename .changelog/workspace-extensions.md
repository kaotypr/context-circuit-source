---
product: both
kind: added
upgrade: Nothing to do unless you want extensions; a workspace without `extensions.yaml` behaves as before. Extensions need a CLI with the `extension` commands; install the version the template pins.
---

Workspaces can now use extensions: folders of instructions that add work before and after each stage — writing an intent, approval, planning, preparing a worktree, implementing, review, delivery, completion, and archiving — for the whole workspace or for chosen repositories, without changing the template. List them in `extensions.yaml` with `extension add`, from a folder in the workspace or a Git repository pinned to the commit its tag or branch names; each member can switch one off unless its author made it required. The agent runs `extension hooks` at the stages it carries out, `record approve`, `worktree prepare`, `record complete`, and `record archive` return the instructions for their own stage, and dispatched planners, workers, and reviewers receive theirs in their brief. Every hook that does not apply is reported with its reason, and no extension can authorize a push, a merge, or anything else that needs your say. `status` counts the extensions on for you and `check` validates them. Scripts an extension declares are validated and reported as not run.
