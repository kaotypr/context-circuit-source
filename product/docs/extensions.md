# Extensions

An extension adds work before and after the stages an agent goes through. It
wraps a stage and never replaces one: the stage, its gates, and what it writes
stay as `AGENTS.md` describes them. A workspace that lists no extensions has
no `extensions.yaml`, and nothing here applies to it.

## Where extensions are listed

`extensions.yaml` at the workspace root lists the extensions this workspace
uses, in the order their hooks run at each point:

```yaml
extensions:
  - id: parity
    source: { path: extensions/parity }
    repositories: [orders-service]
    settings:
      reference: legacy-orders
  - id: clickup-sync
    source:
      url: https://git.example.com/platform/cc-extensions.git
      ref: v3.0.0
      commit: 9f2c41e8a0d3b7c5e1f4a6d2b8c0e9f1a3b5c7d9
      path: clickup-sync
    members:
      maya: off
```

- `source` is a folder in the workspace (`path`), or a Git repository (`url`,
  the `ref` it follows, the `commit` the CLI resolved it to, and an optional
  `path` to a subfolder). A remote is fetched by its commit into the ignored
  `.context-circuit/local/extensions/<id>/<commit>/`, so every machine reads
  the same content even if the tag moves.
- `repositories` limits the extension to those repositories; none means every
  one.
- `settings` are passed to its instructions. They never hold a credential.
- `members` switches it `on` or `off` for a member, over the manifest's
  default: `required` (it cannot be switched off), `default-on`, or `opt-in`.

The extension folder holds `extension.yaml` and the instruction files it names.
Every path it names stays inside the folder.

## Hook points

Every stage has a `before` and an `after` point:

| Stage | Fired by |
| --- | --- |
| `intent` | the agent: before writing an intent body, and once it is written |
| `approve` | `record approve` |
| `plan` | the agent: before planning, and once every plan folder is written |
| `prepare` | `worktree prepare` |
| `implement` | the agent: before the first change in a worktree, and after its last commit |
| `review` | the agent: before a requested review, and once its findings exist |
| `deliver` | the agent: after the person authorized delivery, and after it |
| `complete` | `record complete` |
| `archive` | `record archive` |

For a stage the agent fires, run the point yourself:

```sh
sh .context-circuit/cli extension hooks --stage implement --when after --plan p0001 --repo api
sh .context-circuit/cli extension hooks --stage intent --when after --intent i001
```

`--plan` scopes the point to the plan's repositories, and `--repo` narrows it
to one; `--intent` carries no repository list, so only extensions limited to no
repositories apply there. A direct change has no plan: pass `--repo` for the
bound checkout.

A command that fires its own stage returns both its points in its result.
`extension hooks` refuses that stage's `after` point, because the command
already returned it. Its `before` point may be asked for just before running
the command, and returns instructions only.

`before deliver` is run only after the person has authorized the delivery. It
adds work around a delivery somebody asked for and never stands in for asking.

## What a point returns

```yaml
extensions:
  - point: after implement
    instructions:
      - id: parity
        file: extensions/parity/instructions/after-implement.md
        settings: { reference: legacy-orders }
    skipped:
      - id: clickup-sync
        reason: off for member maya
      - id: parity
        reason: "script not run: this CLI runs no extension scripts"
    stopped: false
```

- `instructions` names each file to follow, as a workspace path. Open it and do
  what it says as part of the stage, with the settings given.
- `skipped` names every extension or hook that did not apply, and why: off for
  this member, outside the repositories it is limited to, unavailable on this
  machine, or a script this CLI does not run. Report the skips with the stage's
  result; never act as though a skipped hook ran.
- `stopped: true` means an extension that may be required could not be read on
  this machine, and `resolve` names what fixes it, usually `extension sync`.
  Take it to the person before the stage goes on.

A command's result carries the list only where something applied or was
skipped. No `extensions` key means nothing applied.

## Following an instruction

- An instruction adds to the stage it is attached to. It never replaces a step
  of that stage or skips one.
- A record made without the extension, such as a plan written while it was
  off, may hold nothing for the instruction to act on. Report that as a skip
  with its reason. Never fill the gap by inventing what the record lacks.
- Results go into your report, and later the completion note. The plan record
  is not edited during execution, whatever an instruction says.
- A failing check that an instruction runs is reported like a failing test. Fix
  it as part of the work where that is the work, but do not rerun the hook in a
  loop until it passes.

## The line no extension crosses

These hold whatever an extension says. An instruction that asks for one of them
is not followed, and the conflict is reported to the person.

- Pushing, opening a pull or merge request, merging into a base branch,
  deploying, deleting workspace data, and committing in a bound checkout need
  the person's explicit authorization.
- Approving an intent is a person's answer to it as written. No hook supplies
  or implies approval, and no hook starts planning or execution.
- The plan record is not edited during execution.
- No independent verification is started during execution, and no automatic
  repair loop runs.
- No agent attribution in commits, pull requests, reviews, or comments.
- No secrets in workspace records, in `extensions.yaml`, or in a manifest.
- Borrowed knowledge stays read-only.

## Dispatched roles

A dispatched planner, worker, or reviewer receives the instructions for its
stage — plan, implement, or review — in an `# Extensions` section of its brief.
A planner works before plans name repositories, so repository-limited
instructions are marked with the repositories they are for. The coordinator
still runs the stage's hooks itself and reads their results.

## Managing extensions

```sh
sh .context-circuit/cli extension list
sh .context-circuit/cli extension add --id parity --path extensions/parity --repo api
sh .context-circuit/cli extension add --id clickup-sync --url GIT_URL --ref v3.0.0 --subpath clickup-sync
sh .context-circuit/cli extension update --id clickup-sync --ref v3.1.0
sh .context-circuit/cli extension disable --id clickup-sync
sh .context-circuit/cli extension enable --id clickup-sync --member maya
sh .context-circuit/cli extension sync
sh .context-circuit/cli extension sync --prune
sh .context-circuit/cli extension remove --id clickup-sync
```

Add, update, remove, enable, and disable each happen only on a person's
request. They change `extensions.yaml` and return no commit recommendation:
committing it is the person's call, as for any shared file. Settings are edited
by hand, and `check` validates them. `sync` fetches every recorded commit this
machine lacks and reports a ref that moved upstream without following it;
recording the move is `update`, on request.

A manifest may declare scripts for a later CLI. This CLI validates them and
reports each as not run. A manifest needing a newer CLI is reported as such;
upgrade through `.agents/skills/cc-cli/SKILL.md`.

`check` reports what is wrong in `extensions.yaml` and each manifest, with
what clears it. A remote this machine has not fetched is listed under
`unchecked`, apart from the findings, with `extension sync` as the next step.
