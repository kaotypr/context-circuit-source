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

An existing wrapper must be clean before reconfiguration (an ignored request
under `.runtime/` is permitted). The command completes schema, semantic, path,
repository, source, managed-README, and output preflight before its first atomic
write. An exact rerun of the successful fresh-bootstrap request is a read-only
success; a materially different bootstrap-shaped request requires explicit
reviewable-change authorization.

If a release archive was initialized with `git init` but still has no `HEAD`,
the neutral extracted-template inventory is accepted with
`wrapper.initialize_git: false` and exact initial-commit authorization. Extra
authored paths stop bootstrap without creating a commit. The standalone engine
contains the trusted release inventory and checks both the editable manifest and
bundled-command digest against it; the manifest cannot whitelist extra files.

Existing-wrapper output uses a cross-file transaction. Every sibling temporary
file is written and synced before any managed target changes. Existing targets
then move to recoverable backups; a rename failure restores all backups and
removes staged artifacts before returning an error.

An operating-system shutdown or process kill can interrupt those renames before
rollback runs. A later `configure-workspace` invocation detects managed-file
siblings ending in `.stage` or `.backup` and refuses to change configuration.
It does not guess which copy is authoritative or delete recovery evidence.
Inspect the named target and sibling files, restore exactly one known-good
target manually, preserve uncertain copies elsewhere, and rerun the same
configuration request. Hard interruption cannot be made fully transactional;
this conservative check keeps recoverable copies from being overwritten.

`initialize-workspace` is retained as a compatibility command. It detects the
Git state, reports whether it routed to fresh bootstrap or existing inspection,
and delegates legacy `--bootstrap` requests to configuration. Bootstrap remains
an explicit internal first-time phase; configuration is not an upgrade workflow.
