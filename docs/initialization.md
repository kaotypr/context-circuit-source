# Initialize a workspace

Create a wrapper from the template, open its root in Codex or Claude Code, and
invoke the host's `initialize-workspace` adapter. The canonical skill inspects
existing files before asking about unresolved choices. Team mode and ignored
local clones are the recommended defaults; solo mode and tracked Git submodules
remain supported.

## Repository preparation

Place every repository at its final path before deterministic validation:

- For an ignored clone, clone into a path such as `repositories/frontend`. Do
  not add that clone to the wrapper index.
- For a submodule, register it with Git at its final path so `.gitmodules` and
  the wrapper index contain the matching gitlink.

Never replace an existing directory, discard local changes, or ignore the whole
`repositories/` directory. The initializer owns only the marked exact-path block
in `.gitignore`:

```text
# kao-delivery-workspace:ignored-clones:start
repositories/frontend/
# kao-delivery-workspace:ignored-clones:end
```

## Review flow

The canonical skill updates `workspace.yaml`, workflow and instruction files,
domain agents, and verified context with reviewable edits. During development,
its deterministic checks are:

```bash
node .agents/bin/kao.mjs initialize-workspace --check-only
node .agents/bin/kao.mjs initialize-workspace
node .agents/bin/kao.mjs validate --check-paths --check-documents
git diff --check
```

The first command reports repository modes, paths, roles, agents, remotes,
default and current branches, cleanliness, repository-local instructions,
required documents, wrapper changes, and warnings without writing. For an
inspection-only request, stop after this command and the read-only validators;
do not run the applying command.

The second command reconciles only the managed ignore block. A repeat run with
unchanged configuration reports `gitignore_changed: false`.

Review the complete diff and initialization summary before approving it. The
workflow does not commit, push, create a pull request, configure credentials, or
mutate an activity system without separate human authorization.

## Maintainer source versus distributable

This source repository includes TypeScript tooling and fixtures used to verify
the release. They are maintainer inputs, not initialized-wrapper content. The
distributable contains the static wrapper structure and bundled Node command,
requires Node 22 or newer, and excludes npm metadata, dependencies, fixtures,
scripts, tests, and development planning artifacts.
