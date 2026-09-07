# Upgrade boundary

An upgrade replaces template-owned files (the shipped `.context-circuit/wrapper/`, `.agents/skills/`,
`.context-circuit/agents/`, and `.context-circuit/docs/`) and never overwrites workspace-owned material.

Preserve on upgrade:

- project identity (`workspace.yaml`);
- accepted Product Knowledge (`context/`);
- sources and provenance (`sources/`);
- plans and plan status (`plans/`, including `plans/archive/`);
- local repository bindings (`repositories.local.yaml`);
- connected repositories (`repositories/`);
- runtime evidence and active worktrees (`.runtime/`).

The authoritative split is declared in `.context-circuit/wrapper/manifest.yaml` under
`upgrade_boundary`. If a template change would change the meaning of a plan,
context unit, or runtime record, the upgrade must report migration-needed and
preserve the old state rather than silently rewriting it.

The v1.0.0 branch rename changes the host-local binding key
`repositories.local.yaml` `anchor_branch` to `base_branch`. Schema 1 bindings
remain readable for compatibility; an explicit migration rewrites the
workspace-owned file and records schema 2:

```text
sh .context-circuit/wrapper/runtime/engine.sh repository-binding-migrate <workspace-root>
```

The migration refuses a binding that contains both names with conflicting
values. Runtime execution records likewise read legacy `anchor_branch` fields
and write canonical `base_branch` fields.

The canonical archive directories are `plans/archive/` and `sources/archive/`.
An upgrade from a workspace containing the previous hidden archive directories
must report migration-needed and preserve them until the user explicitly
authorizes their move; it must not silently rewrite workspace-owned material.

## Nested product home (`.context-circuit/`)

An existing workspace that still has template-owned product trees at the
workspace root (the former top-level wrapper, agents, and docs directories) is
upgraded by moving those trees under `.context-circuit/` as folders (nest, do
not flatten) and retargeting root host pointers (`AGENTS.md`, `WORKFLOW.md`,
`CLAUDE.md`, `CURSOR.md`) into that nested home. `.agents/` stays at the
workspace root. Workspace-owned material is preserved in place: `context/`,
`plans/`, `intent/`, `.runtime/`, identity (`workspace.yaml`), and bindings
(`repositories.local.yaml`, `repositories/`).

This layout move is template-owned. There is no new engine verb. If a record's
meaning still depends on a former top-level wrapper path after the trees have
moved, report migration-needed and preserve the old state rather than silently
rewriting it. The workspace-owned-file precedent remains
`repository-binding-migrate`; this layout change does not migrate those files.
