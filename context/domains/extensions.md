# Workspace extensions

How a workspace adds work before and after the stages an agent goes through,
without forking the template or changing a stage.

## Wrap, never replace

An extension adds work at a stage's `before` or `after` point. The stage, its
gates, and what it writes stay the core's. The rules the entry instruction
holds — authorization for outward actions, the plan record left alone during
execution, no repair loops, no attribution, no secrets — are not configurable,
and an instruction asking to cross one is ignored and reported to the person.

Hooks bind to stages rather than roles. Dispatch is optional, so a hook bound
to a role would silently not fire when the coordinator does the work itself. A
dispatched role receives its stage's instructions in its brief instead.

## What is written where

| Thing | Where | Shared? |
| --- | --- | --- |
| The extensions a workspace uses, in hook order, with source, repositories, settings, and member switches | `extensions.yaml` at the root | Yes |
| An extension: its `extension.yaml` manifest and the instruction files it names | a folder in the workspace, or a Git repository | Yes, or upstream |
| A remote source's fetched commit | this machine's ignored extension cache, one folder per extension and commit | No |

A remote source is pinned to the commit its ref named when it was added. A tag
moved upstream changes nothing until somebody records the move, so every
machine reads the same instructions; sync reports the move and never follows
it. A fetched source has its push URL disabled.

A manifest declares a member default: `required`, which cannot be switched
off, `default-on`, or `opt-in`. A member switch in `extensions.yaml` is read as
the text `on` or `off`, quoted or not, because YAML 1.1 reads either as a
boolean. Every path a manifest names stays inside its folder, the guard
borrowed knowledge uses.

## Firing a point

Four stages are CLI commands — approve, prepare, complete, archive — and the
command returns both its points' blocks after it wrote. The other five —
intent, plan, implement, review, deliver — the CLI never sees, so the agent
asks for each point with `extension hooks`. That command refuses the `after`
point of a CLI-fired stage, so an agent cannot fire it twice or without the
stage; it answers the `before` point with instructions only, for a skill that
must act before running the command.

Resolution walks the entries in file order. An extension off for the member, or
limited to repositories the point does not touch, is skipped with that reason.
A point with no structured repository list — an intent's — admits only
extensions limited to no repositories. A planner brief is the exception: it
lists repository-limited plan instructions marked with what they are for,
because the plans that would name repositories do not exist yet.

An extension this machine cannot read cannot be shown to be optional, so a
point stops for the person unless what could be read of the manifest says it
may be switched off. `status` counts the extensions on for the active member,
which is how an agent in a workspace without any knows the rule does not
apply without reading the file.

## Said, not guessed

Every hook that does not apply is reported with its reason. A record made
without an extension — a plan written while it was off — may hold nothing for
its instruction to act on, and the instruction reports a skip rather than
inventing what the record lacks. The CLI cannot see that case, so the
instruction text says it.

This release runs no code. A manifest may declare scripts for a later CLI;
they are validated and accepted, and each is reported as not run. A manifest
declaring `requires_cli` above the running CLI is reported as needing a newer
one instead of failing to decode.

Owner:

- `context-circuit-source@internal/workspace/extensions.go` — the file and
  manifest formats, member switches, and hook resolution.
- `context-circuit-source@internal/workspace/extension_fetch.go` — pinning a
  ref to a commit, the cache, and sync.
- `context-circuit-source@internal/workspace/extension_edit.go` — the commands
  that change `extensions.yaml`, and the `check` findings.
- `context-circuit-source@internal/cli/extensions.go` — the stage hooks a
  command fires, and `extension hooks`.
- `context-circuit-source@product/docs/extensions.md` — the reference an agent
  follows.
