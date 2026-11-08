# Workspace extensions

How a workspace adds work before and after the stages an agent goes through,
and roles, skills, and named stages beside the shipped ones, without forking
the template or changing a stage.

## Wrap, never replace

An extension adds work at a stage's `before` or `after` point. The stage, its
gates, and what it writes stay the core's. The fixed rules the entry
instruction holds — authorization for outward actions, the plan record left
alone during execution, no attribution, no secrets — are not configurable. No
repair loops is a default instead: an extension may override it, and three
other defaults, only where a person allowed it, as the customization note
describes.

An instruction asking to cross one is not followed and is reported to the
person, while the rest of its extension still applies: refusing the whole
extension would drop the checks it adds for one line it had no power to grant.
A dispatched role's brief repeats those limits, repair loops included, because
the role reads its instruction files without reading the entry instruction. An
allowed repair loop does not lift that line from a role: the coordinator runs
the rounds, and the role never loops by itself.

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
command. Because the CLI never sees implementation, `worktree prepare` returns
`next`, the two implement-point commands, where an extension applies to
implementing in that worktree: agents that read the rule in a skill still
skipped the hook, and a command in the output they act on is harder to miss.

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

## Roles an extension adds

An extension may declare agent roles and add to the four shipped ones. A
declared role is named `x-<extension>-<role>`: no shipped role starts with
`x-`, so an extension role sits beside the shipped ones and never replaces or
shadows one. Its access is `read-only`, briefed like the explorer, or `edit`,
briefed like the worker and dispatched against a plan. Its brief ends with
limits no instruction lifts — no outward action, approval, start of planning or
execution, or other agent — because the role reads its instructions without
the entry instruction. Its model and effort come from the tiering under its full
name, or from the shipped role its manifest calls it `like`; the tiering accepts
any `x-` name, since one member's extension may be off for another.

An addition to a shipped role reaches it through the brief only. Its definition
and access stay as they ship, and an added heading follows the core ones and may
not repeat one: a planner's fixed headings are transcribed into the plan record.

A declared role is a file on each machine, written by `agent setup` from the
version that machine has, for the extensions on for its member, and recorded in
the same inventory as the shipped definitions. That is the trade: roles make
definitions depend on extensions, so an extension change leaves them behind
until setup runs again. The extension commands and `status` say so as
`setup_required` rather than writing definitions themselves, and the inventory
lets setup remove a definition whose extension is gone or off only when nobody
edited it.

## Skills and named stages an extension adds

An extension may ship skills, installed by `agent setup` as
`x-<extension>-<skill>` for the extensions on for the member: a copy of the
folder under `.agents/skills/`, with its frontmatter name rewritten to the
installed one because Codex reads that folder directly, and a one-line wrapper
under `.claude/skills/` and `.cursor/skills/` in the shipped wrapper form. The
copy is shared by every host and recorded in its own inventory, so naming one
host still writes the folder its wrapper points to, and an edited copy is
protected and a stale one removed the way a role definition is.

The installed paths are gitignored rather than committed. The source of truth
is the extension's folder or pinned commit, so committing a copy would only let
it drift from the extension on another machine. No shipped skill and no
template path starts with `x-`, so an extension skill never shadows a shipped
one and a delta between two template tags never writes or removes one.

A skill runs on a person's request or when a core hook instruction names it,
and starts neither planning nor execution. It may declare a named stage, which
other extensions hook as `before|after <extension>/<stage>`. The parser accepts
that form by shape only; whether the stage exists is settled at resolution,
where a stage whose extension is not listed, unreadable, off for the member, or
not declaring it returns one skip with that reason. A manifest is validated
without the workspace, so a hook at an undeclared stage is a `check` finding
rather than an invalid manifest. Running a named stage is the declaring skill's
job: it never enters `record order`, the entry instruction, or a core skill,
because a stage the core does not know cannot be a gate the core enforces.

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
- `context-circuit-source@internal/workspace/extension_roles.go` — roles and
  role additions: validation, resolution, the brief, and `check` findings
  across extensions.
- `context-circuit-source@internal/workspace/agents.go` — writing, refreshing,
  and removing extension role definitions, and `setup_required`.
- `context-circuit-source@internal/workspace/extension_skills.go` — skills:
  validation, the copy and its wrappers, named-stage resolution, and the
  `check` findings across extensions.
- `context-circuit-source@product/docs/extensions.md` — the reference an agent
  follows.
