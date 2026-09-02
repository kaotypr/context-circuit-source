# Upgrade boundary

An upgrade replaces template-owned files (the shipped `wrapper/`, `.agents/skills/`,
`agents/`, and `docs/`) and never overwrites workspace-owned material.

Preserve on upgrade:

- project identity (`workspace.yaml`);
- accepted Product Knowledge (`context/`);
- sources and provenance (`sources/`);
- plans and plan status (`plans/`, including `plans/archive/`);
- local repository bindings (`repositories.local.yaml`);
- connected repositories (`repositories/`);
- runtime evidence and active worktrees (`.runtime/`).

The authoritative split is declared in `wrapper/manifest.yaml` under
`upgrade_boundary`. If a template change would change the meaning of a plan,
context unit, or runtime record, the upgrade must report migration-needed and
preserve the old state rather than silently rewriting it.

The v1.0.0 branch rename changes the host-local binding key
`repositories.local.yaml` `anchor_branch` to `base_branch`. Schema 1 bindings
remain readable for compatibility; an explicit migration rewrites the
workspace-owned file and records schema 2:

```text
sh wrapper/runtime/engine.sh repository-binding-migrate <workspace-root>
```

The migration refuses a binding that contains both names with conflicting
values. Runtime execution records likewise read legacy `anchor_branch` fields
and write canonical `base_branch` fields.

The canonical archive directories are `plans/archive/` and `sources/archive/`.
An upgrade from a workspace containing the previous hidden archive directories
must report migration-needed and preserve them until the user explicitly
authorizes their move; it must not silently rewrite workspace-owned material.
