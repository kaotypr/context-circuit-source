# Agent-facing commands

Run the CLI version this workspace pins in `.context-circuit/CLI_VERSION`. Versions
install side by side, so the shared command on PATH may select a different one; the
`cc-cli` skill resolves the pinned path. The CLI warns on stderr when it is run
against a workspace pinning another version.

Use `context-circuit-cli --workspace PATH --json COMMAND ...` for structured output.
Global flags precede the command. Run `context-circuit-cli help` for the full
surface, and `help` with a command name for that command's signatures and what it
decides, refuses, or leaves to the caller. A fully resolved command also takes
`--help` for its flags; a group such as `member` needs its subcommand first, and
reports which option it could not take. Errors go to stderr with a nonzero exit
status.
Normal output is readable YAML; `--json` gives the agent structured values.
The examples assume the current directory is the initialized workspace.

```sh
context-circuit-cli status
context-circuit-cli member add --id alex --name Alex --band 2
context-circuit-cli member band --id alex --band 3
context-circuit-cli member use --id alex
context-circuit-cli member list
context-circuit-cli repo connect --id api --path ../api --base main
context-circuit-cli repo clone --id web --path repositories/web --base main --url GIT_URL
context-circuit-cli repo init --id docs --path repositories/docs --base main
context-circuit-cli repo base --id api --branch development
context-circuit-cli repo remote --id api --url GIT_URL --default-branch main
context-circuit-cli repo relate --from web --to api --description 'Consumes the API'
context-circuit-cli repo inspect --id api
context-circuit-cli repo fetch --id api --remote origin
context-circuit-cli context find --query billing
```

Records:

```sh
context-circuit-cli record create --kind intent --slug add-billing --title 'Add billing'
context-circuit-cli record approve --id i001 --text 'User approved the billing outcome'
context-circuit-cli record create --kind plan --slug billing-api --title 'Billing API' \
  --intent i001 --repo api
context-circuit-cli record create --kind plan --slug billing-web --title 'Billing UI' \
  --intent i001 --repo web --depends-on p0001
context-circuit-cli record show --id p0001
context-circuit-cli record list
context-circuit-cli record dependencies --id p0002 --depends-on p0001
context-circuit-cli record complete --id p0001 --text 'User requested completion.'
context-circuit-cli record order --intent i001 --mode waves
```

`--band` is optional. A member holding band N allocates intents from `N*100` and
plans from `N*1000`, so band 2 writes `i200` and `p2000`; members without one
allocate from the numbers no band claims. `member band --band 0` clears a band.
Bands must be distinct across the roster, and every command refuses to run while
two members share one. Assigning, changing, or clearing a band never renumbers an
existing record or releases its reservation.

Use the actual returned IDs, not the example numbers. `--repo` and `--depends-on`
are repeatable. Creation reserves IDs and links a new plan from its intent.
Approval and completion notes require the actual corresponding user decision.
Approval records the instant in the intent's frontmatter `approved_at` and appends
the person's words under a heading carrying the same instant; an unapproved
intent has no such field. It records a decision and never establishes one, and
commands do not independently authorize anything. Dependency updates detect
cycles; omit `--depends-on` to clear a plan's dependencies. Completion also
records `completed_at` in the plan's frontmatter beside its written note. Every
record also carries the `created_at` instant of its creation. Instants are
canonical ISO 8601 UTC timestamps, `2026-09-15T10:53:00Z`; other dates are ISO
8601 `YYYY-MM-DD` in UTC.

`record order` derives execution shape from recorded dependencies and completion.
`--mode auto` reports a recommendation and follows it; `waves` and `linear`
select explicitly without changing the recommendation. Waves list each layer, its
concurrency, every plan's start reference per repository, any integration merge a
fan-in needs, and which plans in the layer share a repository. A linear chain
stacks each plan on the previous one and never needs a merge. Completed plans
drop out and their dependents become ready, so an unfinished plan holds
everything behind it; recompute after each wave. Start references name a
predecessor's branch — adjust them yourself once that work is merged and its
branch is gone. The report is advice: it runs, merges, and reserves nothing, and
`shares_repository_with` reflects repositories, not paths.

Worktrees:

```sh
context-circuit-cli worktree prepare --repo api --plan p0001 --copy-mode auto
context-circuit-cli worktree prepare --repo web --plan p0002 --start origin/main
context-circuit-cli worktree list --repo api
context-circuit-cli worktree inspect --repo api --path .worktrees/p0001/api
context-circuit-cli worktree move --repo api --path .worktrees/p0001/api --to ../api-billing
context-circuit-cli worktree repair --repo api --path ../api-billing
context-circuit-cli worktree remove --repo api --path ../api-billing
context-circuit-cli check
```

`--reuse` explicitly selects an existing branch/worktree, without resetting it.
`--discard` on removal explicitly authorizes disposing of its files; never add it
merely to bypass a failure. No command auto-deletes branches, commits, delivers,
runs application setup, executes a plan, or starts independent review.

`template export --path NEW_DIRECTORY` writes the exact blank embedded seed for
inspection or packaging. `init` can initialize that blank seed or a fresh directory.
It cannot migrate an existing v1 or reinitialize an active v2 workspace.

Subagent setup and dispatch:

```sh
context-circuit-cli agent settings
context-circuit-cli agent configure --host codex --role worker --model MODEL_ID --effort high
context-circuit-cli agent setup
context-circuit-cli --json agent dispatch --host codex --role planner \
  --intent i001 --path /actual/repository
context-circuit-cli --json agent dispatch --host codex --role worker \
  --plan p0001 --path /actual/worktree
context-circuit-cli --json agent dispatch --host codex --role explorer \
  --path /actual/repository --task 'Which module owns invoice voiding?'
```

Settings default to host inheritance. `agent setup` applies settings to native
agent files for every host unless `--host` narrows it, because one workspace is
opened in several; the host may need to reload them. Initialization runs it, so a
new workspace needs it again only after `agent configure` or in a fresh clone,
where the ignored definitions never arrive. `agent dispatch` returns the
resolved invocation for the cc-dispatch skill to launch through the host's tools.
`--task` is the assignment for a role no record assigns: an explorer's question,
a reviewer's diff and criteria, or the one slice a worker owns when several share
a plan. A planner is dispatched against a whole intent and refuses it. A worker
owning its whole plan needs none, and a brief carries no task section when none
was given. For a reviewer, `--review-requested` represents a real user request.
No CLI command launches an LLM or treats a dispatch specification as evidence of
completion.

A specification also reports `definition_path` and `definition_installed` for the
host's native role file, and `setup_required` naming the exact `agent setup` run
when that file is absent. Role files are host-local and gitignored, so a workspace
that never ran setup on this host names an agent type the host cannot resolve. It
is a report, not a refusal: the prompt remains complete, so a live spawn tool can
still carry it when native roles are unavailable.

The returned `prompt` is the complete brief and is launched unmodified: it opens
with the working directory, then ownership, then the quoted record, then the task,
and ends with what the role must return. Quoting the record in full is what stops
an agent from searching a repository for a file it was never handed.

`--plan` adds the plan's file path, repositories, dependencies, and full record to
the brief, and states that its dependencies' work is already in the branch's
ancestry and that a concurrent sibling plan is invisible. `--intent` does the same
for an approved intent and is how a planner is dispatched; a planner refuses
`--plan`, because the plan shape is the answer it returns rather than a number it
is handed, and it refuses an intent that is not yet approved. `--shared` is for
several workers inside one worktree: it replaces sole ownership of the directory
with a duty to preserve the other workers' edits. Omit it when each worker has its
own worktree.

Worktree `--copy-path` adds a specific ignored runtime entry. `--copy-mode` accepts
`auto`, `required`, `copy`, or `off`. The reuse report distinguishes native CoW
from fallback copying and reports existing entries or dependency mismatches.
