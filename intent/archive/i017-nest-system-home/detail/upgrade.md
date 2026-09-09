# Template seed and upgrades

How a new workspace and an existing one both end up with the nested home,
without moving the user's project.

## Mechanism

The source checkout still separates three identities: maintainer source, the
distributable template, and an instantiated workspace. Assembly still ships
template-owned files and never ships source-only trees.

What changes is the **shape of the shipped layer** inside the template:

- New workspaces are instantiated with `.context-circuit/wrapper`,
  `.context-circuit/agents`, and `.context-circuit/docs` already nested.
- The source checkout's blank seed remains `template/` at the source root
  (that is the seed's identity, not a user-facing product folder). The seed's
  *contents* follow the nested layout so assembly does not have to rewrite
  the home after copy.

An upgrade may replace template-owned files (the nested product, root host
pointers). It must preserve workspace-owned files: identity, accepted Product
Knowledge, sources, plans and status, local bindings, connected repositories,
runtime evidence, and active worktrees.

Moving `wrapper/`, `agents/`, and `docs/` is a template-owned path change.
The upgrade puts them under `.context-circuit/` and retargets the root
pointers. It does not relocate `context/`, `plans/`, `intent/`, or
`.runtime/`.

If a record's meaning depends on the old top-level `wrapper/` path, the
upgrade reports migration-needed and preserves the old state rather than
rewriting user data in place.

## Interfaces

- Release boundary and shipped/never-ship sets stay the release manifest and
  wrapper manifest — those files themselves move with wrapper.
- Assembler still builds the distributable template; it must emit the nested
  home, not the old root `wrapper/`.
- Migration notes live with the wrapper migrations, under the nested home.

## Edge cases

- A workspace that still has root `wrapper/`, `agents/`, or `docs/` after
  upgrade is unfinished.
- A workspace whose user files were copied into `.context-circuit/` is
  wrong: that is not an upgrade, it is a move of project data this intent
  forbids.
- Credentials, local bindings, and connected repositories stay unshipped and
  unmoved.

## Out of this file

What remains at the workspace root for hosts: [host-surfaces.md](host-surfaces.md).
