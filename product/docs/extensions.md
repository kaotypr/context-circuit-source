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

The extension folder holds `extension.yaml` and the instruction files and
scripts it names. Every path it names stays inside the folder.

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

`worktree prepare` returns these two implement commands as `next` when an
extension applies to implementing in that worktree. Run them where it says.

`--plan` scopes the point to the plan's repositories, and `--repo` narrows it
to one; `--intent` carries no repository list, so only extensions limited to no
repositories apply there. A direct change has no plan: pass `--repo` for the
bound checkout.

A command that fires its own stage returns both its points in its result, and
runs their scripts outside the workspace lock: the `before` point before it
writes, the `after` point once it has. `extension hooks` refuses that stage's
`after` point, because the command already returned it. Its `before` point may
be asked for just before running the command, and returns instructions only: its
scripts run when the command fires it.

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
    ran:
      - id: parity
        status: failed          # ok | failed | skipped
        on_failure: stop
        summary: "exited with status 1: 2 of 14 behaviors differ"
    confirmation_required:
      - id: chat-notify
        action: Post the pull request to the team channel
        confirm: extension hooks --stage implement --when after --plan p0001 --repo api --confirm chat-notify
    skipped:
      - id: clickup-sync
        reason: off for member maya
    stopped: true
    resolve: "take this to the person before the stage goes on: parity failed at after implement …"
```

- `instructions` names each file to follow, as a workspace path. Open it and do
  what it says as part of the stage, with the settings given.
- `ran` is each script that ran: `ok`, `failed`, or `skipped` by the script
  itself, with its `on_failure` and a summary. Report each result with the
  stage's.
- `confirmation_required` is each `each-time` script waiting for the person. Ask
  them about its `action` in one question, and run its `confirm` command only
  on their yes. A no is reported, and nothing runs.
- `skipped` names every extension or hook that did not apply, and why: off for
  this member, outside the repositories it is limited to, unavailable on this
  machine, a script not allowed on this machine, or a hook with no Windows
  command. Report the skips with the stage's result; never act as though a
  skipped hook ran.
- `stopped: true` means the stage does not go on until the person decides: an
  extension that may be required could not be read on this machine, a required
  extension's script is not allowed here, or a script whose `on_failure` is
  `stop` failed. `resolve` names what was found. Take it to the person.

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
- A failing check that an instruction or a script runs is reported like a
  failing test. Fix it as part of the work where that is the work, but do not
  rerun the hook in a loop until it passes.

## Scripts

A hook may run a script beside its instructions, or instead of them:

```yaml
hooks:
  - at: after implement
    instructions: instructions/after-implement.md
    script:
      run: [sh, scripts/parity-run.sh]
      run_windows: [powershell, -NoProfile, -File, scripts/parity-run.ps1]
      on_failure: stop      # stop | report
      consent: none         # none | standing | each-time
      action: ""            # required unless consent is none
      timeout: 10m
```

What a script can rely on:

- **Launch.** `run`, or `run_windows` on Windows, starts with no shell. The
  first element is the program and the rest its arguments; a relative path
  among them is resolved inside the extension folder and may not leave it. A
  hook with no `run_windows` is skipped on Windows, with that reason.
- **Environment.** The member's own, minus the Git variables that would point a
  Git command at another repository. Credentials come from the machine, never
  from `settings`.
- **Working directory.** The worktree, for a prepare, implement, review, or
  deliver point that fires for one repository; otherwise the workspace root.
- **Input.** One JSON object on stdin: `point`, `workspace_root`, `member`,
  `extension` (`id`, `settings`), `record` (`id`, `kind`, `path`,
  `repositories`, or null), `worktrees` (each repository's absolute path),
  and `confirmed`, true only when an `each-time` script runs after the person's
  yes. An archive also passes every record it moved as `records`.
- **Output.** Exit code 0 is success, anything else is failure. Stdout may end
  with one JSON object, `{"status": "ok" | "skipped", "summary": "…",
  "details_path": "…"}`. Stdout over 64 KiB is cut and reported as cut, and
  stderr is kept for a failure's summary.
- **Time.** `timeout`, 10 minutes by default. A script still running then is
  stopped and reported as failed.
- **Lock.** Scripts run outside the workspace lock, so a script may call the
  CLI itself. A script after a command reads what the command wrote, never what
  it was about to.
- **Records made without the extension.** A record with nothing for the script
  to act on, such as a plan without a parity list, is reported as
  `{"status": "skipped"}` with the reason. A script never fills the gap.

### Allowing a version

A script is someone else's code running with the member's own permissions, so
it runs only after the member using this machine allowed that exact version:

```sh
sh .context-circuit/cli extension allow --id parity
```

`allow` shows the source, the version, and each script hook with its command,
`on_failure`, consent, and action, then records the allow in the ignored
`.context-circuit/local/extension-trust.yaml`. It runs nothing. The version is
the commit of a remote source, or a SHA-256 over a folder source's files,
paths and contents, with line endings normalized. Any change to either asks
again. The allow belongs to one member on one machine; another member, or
another machine, allows separately.

`extension list` shows each extension with scripts as `allowed`, `not allowed`,
or `changed since allowed`. A script not allowed is skipped with that reason
and the command that clears it. A `required` extension whose script is not
allowed stops the stage.

No flag or setting skips the allow, and the CLI never allows anything by
itself. Run `extension allow` only when the person asks you to; a skip naming
it is reported to them, not cleared by you.

### Consent

`consent` is the author's statement of what a script does beyond this machine,
made to the person who allows it:

- **`none`.** The script acts only locally: it reads, checks, or writes inside
  the worktree. A script that publishes while declaring `none` is the author's
  defect, and an extension the workspace should remove.
- **`standing`.** The script runs when its point fires. Enabling the extension
  in `extensions.yaml`, and allowing this version after reading its `action`,
  is consent ahead of time to that one action, for that hook. It is the only
  place configuration stands in for a person's answer, and it covers nothing
  else: not another action of the script, and not anything the agent does. A
  new version is a new allow, so a changed action is read again.
- **`each-time`.** The script does not run when its point fires. It is listed
  under `confirmation_required` with its `action`; ask the person, and on a
  yes run its `confirm` command, which passes `confirmed: true`. At a stage the
  CLI fires, nobody can be asked mid-command, so an `each-time` script there is
  skipped with that reason.

`standing` and `each-time` need an `action`; a manifest without one is a
`check` finding.

### Failure

| Outcome | `on_failure: report` | `on_failure: stop` |
| --- | --- | --- |
| A `before` script fails | The stage goes ahead; the failure is in `ran` | The stage does not happen; a CLI command exits 1 having written nothing |
| An `after` script fails | Reported | Reported as `stopped`; what the stage wrote stays, and the person decides before the next stage |
| A script is skipped: off, not allowed, unavailable, no Windows command | Reported as skipped | Reported as skipped, except a `required` extension whose script is not allowed, which stops |

A stop is never a loop. Report it, and the person decides; the hook is not
rerun until it passes.

## The line no extension crosses

These hold whatever an extension says. An instruction that asks for one of them
is not followed, and the conflict is reported to the person.

- Pushing, opening a pull or merge request, merging into a base branch,
  deploying, deleting workspace data, and committing in a bound checkout need
  the person's explicit authorization. The one exception is a `standing`
  script's own named action, below.
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
still runs the stage's hooks itself and reads their results. Scripts are always
run by the coordinator through `extension hooks`, never by a role.

## Managing extensions

```sh
sh .context-circuit/cli extension list
sh .context-circuit/cli extension allow --id parity
sh .context-circuit/cli extension add --id parity --path extensions/parity --repo api
sh .context-circuit/cli extension add --id clickup-sync --url GIT_URL --ref v3.0.0 --subpath clickup-sync
sh .context-circuit/cli extension update --id clickup-sync --ref v3.1.0
sh .context-circuit/cli extension disable --id clickup-sync
sh .context-circuit/cli extension enable --id clickup-sync --member maya
sh .context-circuit/cli extension sync
sh .context-circuit/cli extension sync --prune
sh .context-circuit/cli extension remove --id clickup-sync
```

Add, update, remove, enable, disable, and allow each happen only on a person's
request. They change `extensions.yaml` and return no commit recommendation:
committing it is the person's call, as for any shared file. Settings are edited
by hand, and `check` validates them. `sync` fetches every recorded commit this
machine lacks and reports a ref that moved upstream without following it;
recording the move is `update`, on request.

`allow` is the one command here that changes nothing shared: it writes this
machine's trust file. A manifest needing a newer CLI is reported as such;
upgrade through `.agents/skills/cc-cli/SKILL.md`.

`check` reports what is wrong in `extensions.yaml` and each manifest, with
what clears it. A remote this machine has not fetched is listed under
`unchecked`, apart from the findings, with `extension sync` as the next step.
