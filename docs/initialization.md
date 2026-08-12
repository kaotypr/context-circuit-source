# Configure and bootstrap a workspace

The canonical configurator supports both a neutral release archive and an already
configured wrapper. The human-facing walkthrough is in
[getting started](getting-started.md).

## Neutral template

An unused template has `repositories: {}`, an empty managed ignore block, no
`repositories/` directory, and no technology-specific domain agent. This state
is valid for inspection but cannot run work until at least one repository is
registered.

## Bootstrap contract

After interviewing the human, the host writes an ignored, validated
`.runtime/bootstrap/request.json` containing the final workspace configuration,
wrapper Git intent, repository source actions, and explicit initial-commit
authorization. It presents those actions before executing:

```bash
node .agents/bin/cc.mjs configure-workspace --request .runtime/bootstrap/request.json
```

Repository sources are `new`, `clone`, `existing`, and `submodule`. New and
cloned paths must not exist. Existing paths must already be exact Git roots.
Submodules are tracked and never ignored. New repositories receive an authorized
empty base commit so worktree execution has a stable base. Existing and cloned
repositories receive no artificial commit.

The command configures exact ignore paths before repository creation and makes
the fully configured wrapper state its first commit. It preserves partial state
if an external clone or Git operation fails; it never cleans or deletes recovery
evidence automatically.

## Existing wrapper

For inspection, run:

```bash
node .agents/bin/cc.mjs configure-workspace --check-only
```

For approved configuration edits in an existing wrapper, edit the reviewable
files and run:

```bash
node .agents/bin/cc.mjs configure-workspace --request .runtime/bootstrap/request.json
node .agents/bin/cc.mjs validate --check-paths --check-documents
git diff --check
```

Reruns are idempotent. The configurator does not commit changes in an existing
wrapper and never pushes, creates remotes or pull requests, mutates activity,
merges, or deploys.

`initialize-workspace` is retained as a compatibility command. It detects the
Git state, reports whether it routed to fresh bootstrap or existing inspection,
and delegates legacy `--bootstrap` requests to configuration. Bootstrap remains
an explicit internal first-time phase; configuration is not an upgrade workflow.
