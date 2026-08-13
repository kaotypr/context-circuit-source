# Execute an approved plan

`execute-plan` is the core implementation workflow for one approved numbered
plan. It keeps one cumulative branch and worktree for each affected repository,
executes the authoritative task graph in dependency order, and requires a
holistic verifier before one review handoff per repository.

```sh
node .agents/bin/cc.mjs execute-plan --plan plans/<repository-key>-plans/<number>-<slug> --version <version> --approved-digest <sha256:...>
```

If the plan changes before merge, explicitly revise and reapprove it, then
resume the existing runtime:

```sh
node .agents/bin/cc.mjs execute-plan --resume-run <run-id> --plan <reference> --version <version> --approved-digest <sha256:...> --reason "..."
```

The append-only revision record preserves prior evidence, re-executes changed
dependency closures, and invalidates stale final/review/merge evidence.
