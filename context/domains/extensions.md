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
command fires both its points around its own write, outside the workspace lock:
the setup gate and precheck, then the `before` point, then the lock and the
write, then the `after` point once the lock is released. A slow script
therefore never blocks another command, and a script may call the CLI itself.
The cost is that the precheck and the write can see different states, so an
`after` point reads only what the command returned and its command line. A
`before` point that stops leaves the command unrun, exiting 1 having written
nothing; an `after` stop is reported and what the stage wrote stays.

The other five — intent, plan, implement, review, deliver — the CLI never
sees, so the agent asks for each point with `extension hooks`. That command
refuses the `after` point of a CLI-fired stage, so an agent cannot fire it
twice or without the stage; it answers the `before` point with instructions
only, never running its scripts, for a skill that must act before running the
command.

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

A manifest declaring `requires_cli` above the running CLI is reported as
needing a newer one instead of failing to decode.

## Scripts, trust, and consent

A hook may run a script. What the script can rely on:

| Concern | What holds |
| --- | --- |
| Launch | No shell; relative paths resolved inside the extension folder; `run_windows` on Windows, and a hook without one is skipped there |
| Environment | The member's own, minus the Git variables the Git runner strips |
| Working directory | The worktree, for a repository-scoped point that fires for one repository; otherwise the workspace root |
| Input and output | A JSON object on stdin; exit code 0 is success; stdout kept to 64 KiB, with a final JSON object read from it |
| Time | Ten minutes unless the manifest names a timeout |

A script reports a record it has nothing to act on as skipped, and the CLI
carries that through.

Whether a workspace uses an extension is decided in `extensions.yaml` and
reviewed like any commit. Whether its scripts may run is decided by the member
on each machine. A script runs with that member's own permissions, and nobody
else's review can stand in for that.

So a script runs only after `extension allow` recorded this member, this
machine, and this exact version in the ignored local trust file. The version is
a remote source's commit, or a SHA-256 over a folder source's relative paths
and contents. Line endings are normalized first, so a Windows checkout does not
ask again on every run. Any change asks again.

- `allow` shows the scripts, their `on_failure`, and their consent actions, and
  runs nothing.
- No flag or setting skips it, the CLI never allows by itself, and the agent
  runs it only when the person asks.
- A required extension whose script is not allowed stops the stage rather than
  going on without its work.

Each script hook declares consent:

- `none` is the author's statement that it acts only locally.
- `each-time` never runs when its point fires. It is listed with its action
  for the person, and runs on their yes through `--confirm`, alone, told that
  it was confirmed. At a CLI-fired stage nobody can be asked mid-command, so it
  is skipped there.
- `standing` runs when its point fires. Enabling the extension, and allowing
  its version after reading the action, was the person's consent ahead of time.
  It is the one place configuration stands in for a person's answer, so it
  covers only that action, for that hook, and a new version is read again.

A manifest declaring `standing` or `each-time` without an action does not
validate, which `check` reports.

Each script hook also declares `on_failure`: `stop` halts the stage for the
person and `report` lets it go on. A stop is never a loop; the hook is not
rerun until it passes.

Owner:

- `context-circuit-source@internal/workspace/extensions.go` — the file and
  manifest formats, member switches, and hook resolution.
- `context-circuit-source@internal/workspace/extension_fetch.go` — pinning a
  ref to a commit, the cache, and sync.
- `context-circuit-source@internal/workspace/extension_edit.go` — the commands
  that change `extensions.yaml`, and the `check` findings.
- `context-circuit-source@internal/workspace/extension_run.go` — launching a
  script, its input, output, and timeout.
- `context-circuit-source@internal/workspace/extension_trust.go` — the folder
  hash, the local trust file, and `extension allow`.
- `context-circuit-source@internal/cli/extensions.go` — the stage hooks a
  command fires around its write, and `extension hooks`.
- `context-circuit-source@product/docs/extensions.md` — the reference an agent
  follows.
