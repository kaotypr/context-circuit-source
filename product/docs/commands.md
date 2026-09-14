# Agent-facing commands

Use `context-circuit --workspace PATH --json COMMAND ...` for structured output.
Global flags precede the command. Run `context-circuit help` for the full surface;
each command supports `--help`. Errors go to stderr with a nonzero exit status.
Normal output is readable YAML; `--json` gives the agent structured values.
The examples assume the current directory is the initialized workspace.

```sh
context-circuit status
context-circuit member add --id alex --name Alex
context-circuit member use --id alex
context-circuit member list
context-circuit repo connect --id api --path ../api --base main
context-circuit repo clone --id web --path repositories/web --base main --url GIT_URL
context-circuit repo init --id docs --path repositories/docs --base main
context-circuit repo base --id api --branch development
context-circuit repo relate --from web --to api --description 'Consumes the API'
context-circuit repo inspect --id api
context-circuit repo fetch --id api --remote origin
context-circuit context find --query billing
```

Records:

```sh
context-circuit record create --kind intent --slug add-billing --title 'Add billing'
context-circuit record approve --id i001 --text 'User approved the billing outcome'
context-circuit record create --kind plan --slug billing-api --title 'Billing API' \
  --intent i001 --repo api
context-circuit record create --kind plan --slug billing-web --title 'Billing UI' \
  --intent i001 --repo web --depends-on p0001
context-circuit record show --id p0001
context-circuit record list
context-circuit record note --id p0001 --text 'Implemented API; normal tests passed.'
context-circuit record dependencies --id p0002 --depends-on p0001
context-circuit record complete --id p0001 --text 'User requested completion.'
context-circuit record order --intent i001 --mode waves
```

Use the actual returned IDs, not the example numbers. `--repo` and `--depends-on`
are repeatable. Creation reserves IDs and links a new plan from its intent.
Approval and completion notes require the actual corresponding user decision.
Commands do not independently authorize anything. Dependency updates detect
cycles; omit `--depends-on` to clear a plan's dependencies. Completion also
records a `completed` date in the plan's frontmatter beside its written note.
Dates are ISO 8601 `YYYY-MM-DD` in UTC throughout.

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
context-circuit worktree prepare --repo api --plan p0001 --copy-mode auto
context-circuit worktree prepare --repo web --plan p0002 --start origin/main
context-circuit worktree list --repo api
context-circuit worktree inspect --repo api --path .worktrees/p0001/api
context-circuit worktree move --repo api --path .worktrees/p0001/api --to ../api-billing
context-circuit worktree repair --repo api --path ../api-billing
context-circuit worktree remove --repo api --path ../api-billing
context-circuit check
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
context-circuit agent settings
context-circuit agent configure --host codex --role worker --model MODEL_ID --effort high
context-circuit agent setup --host codex
context-circuit --json agent dispatch --host codex --role worker \
  --plan p0001 --path /actual/worktree --task 'Implement the assigned tasks from p0001'
```

Settings default to host inheritance. `agent setup` applies settings to native
agent files; the host may need to reload them. `agent dispatch` returns the
resolved invocation for the cc-dispatch skill to launch through the host's tools.
For a reviewer, `--review-requested` represents a real user request. No CLI command
launches an LLM or treats a dispatch specification as evidence of completion.

`--plan` adds the plan's repositories, dependencies, and record to the brief, and
states that its dependencies' work is already in the branch's ancestry and that a
concurrent sibling plan is invisible. `--shared` is for several workers inside one
worktree: it replaces sole ownership of the directory with a duty to preserve the
other workers' edits. Omit it when each worker has its own worktree.

Worktree `--copy-path` adds a specific ignored runtime entry. `--copy-mode` accepts
`auto`, `required`, `copy`, or `off`. The reuse report distinguishes native CoW
from fallback copying and reports existing entries or dependency mismatches.
