# Agent tests

A contributor tool that runs a coding host against a fresh fixture workspace.
It is a separate Go module. It is never shipped, and CI does not run it.

A change is in scope when an agent in a generated workspace could see different
bytes or different CLI output after it. That covers the files `init` writes,
the layout of those files, and what the CLI prints.

Run `surface-diff` first. A non-empty result needs a design case and agent-test
evidence in the pull request. An empty result needs neither. This is not CI
and not a release gate. You run it yourself, on your own accounts, before you
open the pull request.

The Design case and the Agent tests sections of the pull request description
are described in `CONTRIBUTING.md` under "Changes that affect agents".

## Quickstart

From a clean checkout of this repository, with Go 1.25 or newer:

```sh
cd agent-tests
go test ./...
go build -o /dev/null ./cmd/agent-tests
```

`go test ./...` needs no model. It builds the real CLI, seeds each scenario,
and grades a scripted solution.

Compare your branch with the branch the pull request targets. That branch is
`development/v2` unless you were asked to target another:

```sh
go run ./cmd/agent-tests surface-diff --base development/v2
```

The command prints the report and writes it under `agent-tests/out/`, which is
gitignored. The directory name is a UTC timestamp, the base commit, and the
branch commit, for example `out/<timestamp>-<base>-<branch>/surface/`. When
the two refs are the same commit, that directory holds only `surface.json` and
`surface.md`. When the refs differ, it also holds `base/`, `branch/`, and
`outputs/`. The command exits 0 whether or not anything differs. Flags and
what is compared are in "surface-diff".

When the report says `No agent tests needed: surface diff is empty`, stop.
Paste that sentence in the pull request. You need no design case and no host
run.

When the report names scenarios, read "Trigger map" and "Hosts" below, then
preview the plan. `--dry-run` prints the plan and starts nothing:

```sh
go run ./cmd/agent-tests run --dry-run
```

A full dry run lists every scenario on every host in `hosts.yaml`, one
combination per line (`scenario x host x side x rN`), and notes each pinned
model. Narrow it the same way a live run is narrowed:

```sh
go run ./cmd/agent-tests run --scenarios orient --hosts claude --runs 1 --base development/v2 --dry-run
```

With `--base`, run 1 is the base ref and then your branch. The model note says
the pin is provisional until a baseline run confirms it.

`probe` runs `<command> --version` for each host. It does not call a model.
It prints the configured version, the command's version, and whether the
pinned model is set. It then lists every configured auth variable name, and
prints `auth: ready` or `auth: not ready` with the missing names. It does not
list which names are set:

```sh
go run ./cmd/agent-tests probe
```

It then prints an older overhead measurement (the figures are in "Tokens and
cost") and says the isolated figure is lower. It does not start a prompt.

If the surface report says the change looks larger, open an issue with the
design case and wait for the maintainer, who reviews, before any `--yes` run.
What counts as larger is in "Trigger map".

A live run needs `--yes`. `probe --yes` starts the overhead prompt and spends
money. Without `--yes`, `run` prints the plan and this sentence on stdout and
again on stderr, then exits 1:

```text
this run would use model usage on your account; re-run with --yes to start
```

One scenario on one host, with the base ref run on the same machine, looks
like this. Replace `orient` and `claude` with a scenario and a host the
surface report names. Repeat for each required pair. `--yes` starts the host
and spends on your account:

```sh
go run ./cmd/agent-tests run --scenarios orient --hosts claude --runs 1 --base development/v2 --yes
```

Each combination is written under `agent-tests/out/<timestamp>/`. The run id
is a UTC timestamp. `runs.jsonl` sits in that directory. One run's files sit
in `runs/<scenario>/<host>/<side>/r<N>/`.

After the runs finish, write the pull-request tables. Neither command starts
a host. `report` reads the newest directory under `agent-tests/out` that
contains `runs.jsonl`:

```sh
go run ./cmd/agent-tests report
go run ./cmd/agent-tests report --design
```

`report` writes `report.md` and `bundle.zip` in that run directory. `report
--design` also writes `design.md` there. Paste `report.md` into the Agent
tests section. Paste the before/after from `design.md` into the Design case.
Attach `bundle.zip`. You still write why the change is needed, what is better,
the catch, and the alternatives. `CONTRIBUTING.md` describes those sections
under "Changes that affect agents". What each flag writes is in "report".

## Trigger map

`triggers.yaml` maps a changed source path to scenarios, hosts, and a
before/after kind. `surface-diff` applies that file. An empty `hosts` list
means any one host. A non-empty list means exactly those hosts. The report
prints the scenario names, the host names, and the kinds. Read "Hosts" and
run what they name.

`product/AGENTS.md.in` always adds `orient`. The heading you edited adds the
scenarios in this table. "Independent review on request" has no scenario in
the catalogue.

`product/workspace-README.md` runs `orient`. Its kind is `always-loaded`. It
is written to `README.md` at init and is not a release-manifest line.

| Section | Scenarios |
|---|---|
| Orient and initialize | `orient`, `join` |
| Updating the workspace template | `orient`, `join` |
| Who is working | `orient`, `join` |
| Language | `intent-lang` |
| Shared knowledge | `orient`, `borrowed`, `complete` |
| Intent and planning | `intent`, `plan`, `direct` |
| Worktrees and implementation | `execute` |
| Independent review on request | none |
| Delivery, completion, and organization | `complete`, `execute` |

Skill files map like this. A `route.md` uses the same scenarios as that
skill's `SKILL.md`. Its hosts list is `claude` and `cursor` only. Codex does
not read route files.

| Skill | Scenarios |
|---|---|
| `cc-workspace` | `join` |
| `cc-intent` | `intent` |
| `cc-plan` | `plan` |
| `cc-dispatch` | `plan`, `execute` |
| `cc-worktree` | `execute` |
| `cc-knowledge` | `complete`, `borrowed` |
| `cc-complete` | `complete` |
| `cc-direct` | `direct` |
| `cc-cli` | `orient` |

Shipped docs: `commands.md` runs `intent`; `working.md` runs `intent`, `plan`,
and `complete`; `workspace.md` runs `join` and `borrowed`; `worktrees.md` runs
`execute`; `agents.md` runs `plan` and `execute`; `how-it-works.md` runs
`orient`.

Template data: `INDEX.md` and `glossary.md` run `orient`, `complete`, and
`borrowed`; `intent/README.md` runs `intent`; `plans/README.md` runs `plan`;
`sources/README.md` runs `orient` and `borrowed`; `gitignore.in` runs `join`
and `execute`; `workspace.yaml` and `members.yaml` run `join` and `orient`.

CLI output: `topics.go` runs `orient`, and a changed help topic adds the
scenario for that command (`record` adds `intent` and `plan`, `worktree` adds
`execute`, `knowledge` adds `borrowed`, `workspace`, `repo`, `member`, and
`agent` add `join`, `context` adds `complete`). `commands.go` is an upper
bound (`borrowed`, `complete`, `execute`, `intent`, `join`, `plan`); narrow
it to the command that changed. `cli.go` runs `complete`, `intent`, and
`plan`. `records.go` runs `intent` and `plan`. `inspect.go`,
`readability.go`, and `review.go` run `complete`. `worktrees.go` runs
`execute`. `readme.go` runs `join`.

`role-tiering.yaml` and `internal/workspace/agents.go` run `plan` and
`execute` on any one host, and the row also names the host whose generated
definition changed. `claude-settings.json` runs `orient` and `execute` on
Claude Code only. The launcher runs `orient` and `join`. The shipped installer
runs `orient`. The trigger notes call that a faithful install through the
shipped installer. The runner has no separate faithful mode, so those
scenarios run the same way as the others. The Windows launcher is for when
the maintainer, who reviews, asks. `scripts/release-manifest.txt` and
`assets.go` are a structure change: every scenario on any one host, plus
`orient` on `claude`, `codex`, and `cursor`.

| Kind | Before / after to show |
|---|---|
| `always-loaded` | The file diff and its size from `surface-diff`. From `orient`, the reads list and the reply, plus the record from the section's scenario. |
| `skill-text` | The skill diff and its size. The files that scenario produces (`design_artifacts` in the scenario file). |
| `skill-route` | The route diff. The reads list: whether the host loaded the skill, and at which step. |
| `shipped-docs` | The doc diff. The reads list, plus the record from the mapped scenario. |
| `structure` | The `base/` and `branch/` trees under the surface directory. `orient`'s reads list on each host. |
| `cli-action` | The command output, text and JSON, from `surface-diff`. The reply and the files the run wrote. |
| `record-template` | The fresh skeleton from `record create`. The record the agent filled in. |
| `check-findings` | `check` text and JSON on the same fixture. The `complete` reply and the edits that cleared the finding. |
| `role-definitions` | The generated definition files. The plan folder or the commits the run produced. |
| `launcher` | The command output captured with the run. |
| `template-data` | The file diff in the fresh workspace. The record or file the mapped scenario writes. |
| `host-settings` | The settings diff. A run cannot show which permission prompts a person sees. |

When a change spans kinds, show each kind's artifact. The deterministic half
is the surface diff. The agent half comes from the runs you already have to
make. `report --design` pairs them.

These paths need no scenarios when they are the whole change, and no design
case either. A normal pull-request description is enough:

- `**/*_test.go`, `internal/cow/**`, and `internal/installer/**`
- `.github/**`, `scripts/*.sh`, and `release/**`
- `README.md`, `CONTRIBUTING.md`, `WORKFLOW.md`, `CLI.md`, `SECURITY.md`, and `CODE_OF_CONDUCT.md`
- `context/**`, `.agents/**`, and `agent-tests/**`
- `product/README.md`, `product/assets/**`, and `product/LICENSE`
- `VERSION` and `CLI_VERSION`

`surface-diff` lists those paths under "Ignored" when they changed beside
something that does reach a workspace.

A larger change is agreed in an issue first. The stop before any `--yes` run
is in "Quickstart". Larger means any of these:

- layout: a shipped file added, removed, or moved
- a skill added or removed, or a change to the shape of what a skill produces (headings, frontmatter, a new or removed file in a plan's folder, or a record skeleton's headings or frontmatter keys)
- a CLI action: help text changed, JSON keys changed, or check findings changed
- a gate: an added or removed line matches a keyword (`push`, `commit`, `commit_recommended`, `reconcile_required`, `attribution`, `approval`, which also matches `approve` and `approved`, the planning stop, which also matches `planning_recommended`, `sources`, or borrowed knowledge)
- migration: the `version` field in `workspace.yaml` changed, or the ids file layout changed

`surface-diff` reports this as "Looks larger" and tags each reason exact or
heuristic. Layout, and a skill added or removed, are exact. Skill shape, a
CLI action, a gate, and migration are heuristic. A body-only edit of a skill
is not a shape change. A `VERSION` or `CLI_VERSION` bump is not migration.
A person decides whether the change is actually larger.

Borderline cases:

1. An empty surface diff needs no scenarios, even for a file under `product/`.
2. A changed path with no row is printed as "No row: the author picks the closest scenario." Pick the closest scenario and say why in the pull request.
3. When nothing fits, write a one-off scenario in `scenarios/` in the same shape as the files already there, include it in the pull request, and run it. `cc-review`, `cc-deliver`, `cc-stacked`, and `cc-system-design` have no row. An unknown check name fails, so use checks the runner already has.

The maintainer, who reviews, may ask for more scenarios, hosts, or runs, and
may waive a trivially safe text change. The maintainer, who reviews, decides
disagreements and any change to the table.

## Hosts

`hosts.yaml` names three hosts: `claude` (Claude Code), `codex` (Codex), and
`cursor` (Cursor, command `cursor-agent`). Each has one pinned model.

An empty hosts list means any one host. A non-empty list means exactly those
hosts, and nothing else.

- `product/CLAUDE.md.in` and `product/host/claude-settings.json` list `claude` only.
- `product/CURSOR.md.in` and `product/cursor.mdc.in` list `cursor` only.
- A skill `route.md` lists `claude` and `cursor` only. Codex does not read route files.
- A role-definition change has an empty hosts list, so `plan` and `execute` run on any one host. The note on that row also names the host whose generated definition changed.
- A structure change is every scenario on any one host, plus `orient` on `claude`, `codex`, and `cursor`.

The maintainer, who reviews, may ask for more hosts. When you lack a host the
change needs, say so. The maintainer, who reviews, runs it.

Run one time per scenario per required host. Run more when the maintainer,
who reviews, asks because a result is borderline. `--runs` defaults to 1.

## What "before" is

The before side is a committed baseline when your host version and model match
`hosts.yaml`. No baseline is committed yet. There is no baseline directory in
this repository.

Until one is committed, run the base ref and your branch alternately on the
same machine. Pass `--base` with the pull request's target branch. With
`--runs 1` the order is base, then branch. On a later repetition, an even run
number starts with the branch and an odd run number starts with the base, so
the same side is not always second.

`report --baseline DIR` can read a directory that holds `runs.jsonl`, or
`summary.json` (`{"rows":[...]}`) when the log is absent. A summary baseline
row is labelled `(summary)`. Without `--baseline`, rows in the same
`runs.jsonl` with `"side":"base"` are the before side.

## Safety

Runs do not wait for a permission prompt, and they are graded from the end
state: files, Git refs, and `check --json`.

Claude Code is started with `--dangerously-skip-permissions`. Cursor is
started with `--force` and `--trust`. Codex is started with `--sandbox
workspace-write`. Because those commands run with your rights, run someone
else's branch only in a container or a VM.

Each host gets an empty config home under the run root. None of the three is given an MCP config file:

- Claude Code uses an empty `CLAUDE_CONFIG_DIR` and an empty `HOME`, and `--strict-mcp-config` with no MCP config file. Auth comes from `ANTHROPIC_API_KEY` or `CLAUDE_CODE_OAUTH_TOKEN`.
- Codex uses an empty `CODEX_HOME` and an empty `HOME`, and `--ignore-user-config`. The API key is passed only on the stdin of `codex login --with-api-key`, and is removed from the child environment. Auth comes from `OPENAI_API_KEY`.
- Cursor uses an empty `HOME` and an empty `XDG_CONFIG_HOME`. Auth comes from `CURSOR_API_KEY`.

The child environment is an allowlist. It sets `PATH` to the decoy directory
and then your `PATH`, `TMPDIR` and `GOCACHE` under the run root, `LANG` to
`C.UTF-8`, `GIT_CONFIG_NOSYSTEM` to `1`, `GIT_CONFIG_GLOBAL` to the run's
gitconfig, `GIT_TERMINAL_PROMPT` to `0`, and `NO_COLOR` to `1`. It copies the
auth variables named in `hosts.yaml` when they are set. It does not copy your
`HOME`, `XDG_*`, `SSH_AUTH_SOCK`, `SSH_AGENT_PID`, `GITHUB_TOKEN`, `GH_TOKEN`,
`AWS_*`, or `KUBECONFIG`. The host extras above are applied last, which is how
a host receives its own `HOME`.

The run root is a temporary directory outside any repository. The checkout is
read through git and is not written. Auth comes only from the environment
variables named in `hosts.yaml`. Nothing is read from or written to the real
home, the real workspace, or user config. The one shared place is the Go
module cache: the caller's `GOMODCACHE` when that is set, otherwise the path
`go env GOMODCACHE` prints. Building the CLI uses that cache with
`GOFLAGS=-mod=readonly` and `GOTOOLCHAIN=local`, so it reads modules already
there. Scenario setup commands run as `sh -c` in the fixture workspace, so
only run scenarios from a checkout you trust.

You pay for runs on your own account. Use a spend-limited key. Claude Code
also receives `--max-budget-usd 3` as a guard against a runaway process.

## Tokens and cost

The primary metric is total input tokens. Output tokens are reported
separately.

- Claude Code sums input, cache-creation, and cache-read tokens, and sets that total only when all three figures are present. When the result includes `modelUsage`, the sum across the models named there replaces the total only when the input, cache-read, and cache-creation figures are all present. Those models include subagents. The dollar figure is `total_cost_usd` as the host reported it. The report calls the total a list-price equivalent. On a subscription you may not be billed that amount.
- Codex sums `input_tokens` over every `turn.completed` event. Its cost comes from the `gpt-6-sol` row in `prices.yaml` (cost basis `price-table`). That row has a long-context tier. The uncached split treats cached and cache-write tokens as already inside `input_tokens`. That split is an assumption and has not been confirmed on a live run. Whether reasoning tokens are already inside output tokens is the same kind of assumption.
- Cursor sums `inputTokens`, `cacheReadTokens`, and `cacheWriteTokens`. The pinned model is `auto`, which has no price row, so the cost is left unset.

A missing token or cost is `n/a`. The report never writes 0 for a missing
number. A null total has no percentage.

A failed check, a lower result than before, or more than 15% more total input
tokens is a warning. The maintainer, who reviews, weighs it. It does not
block the pull request. An increase of exactly 15% stays under the line. The
warning shows the change to one decimal place, so a rise of 15.1% does not
display as 15%. Pass ranks above inconclusive, and inconclusive ranks above
fail, which is how "lower" is decided.

The model pins in `hosts.yaml` are provisional. `probe` and a dry run say so:
`provisional: not yet confirmed by a baseline run`. Codex is pinned at effort
`medium`. Claude Code and Cursor omit the effort flag, and the report shows
that as `default`.

`probe` without a live prompt prints an older measurement of a fixed "Reply
with exactly: OK", taken with a developer's own host config loaded. The
isolated figure is lower:

- `claude`: 42500 input tokens, $0.2519 list
- `codex`: 15974 input tokens
- `cursor`: 21366 input tokens (2 uncached + 21364 cache write), 73-75 output

A scenario is many times that. The bands in "Scenarios" are guesses.

## Scenarios

One fixed prompt each, run in a fresh fixture workspace built by the CLI from
the ref under test. State is seeded through the real CLI instead of chaining
turns, and the result is graded from what the agent leaves behind: files, Git
refs, and `check --json`. Nothing depends on exact reply wording. Reply checks
are patterns.

Each scenario's human items live in its yaml `human` list. The table lists
covers, checks, and token bands. The bands are a guess of total input tokens
until a baseline run exists.

| Scenario | What it covers | Scenario checks | Input tokens (a guess until a baseline run exists) |
|---|---|---|---|
| `orient` | always-loaded files, retrieval through the knowledge index, `sources/` staying passive | `orient-reply`, `orient-clean-trees` | 60k-250k |
| `join` | a new member on a fresh clone: member, band, role, local setup | `join-members`, `join-local-files` | 150k-600k |
| `intent` | an intent from a request, then stop for approval | `intent-record`, `intent-no-plan` | 120k-500k |
| `intent-lang` | the same in the member's language | `intent-record`, `intent-no-plan`, `intent-lang-body` | 120k-500k |
| `plan` | planning an approved intent from real code, then stop | `plan-record`, `plan-linked` | 300k-1.5M |
| `execute` | a plan run in a worktree on its own branch | `execute-worktree`, `execute-records-untouched`, `execute-tests-pass` | 0.5M-3M |
| `complete` | completion plus knowledge reconciliation, ending check-clean | `complete-plan-done`, `complete-note` | 200k-900k |
| `direct` | a direct change in a bound checkout without records | `direct-readme`, `direct-no-records` | 60k-250k |
| `borrowed` | borrowed knowledge stays read-only and is not copied | `borrowed-untouched`, `borrowed-reply` | 80k-400k |

The overhead measurement (a fixed reply in an empty directory and in the
fixture workspace) is `probe`, not a scenario file. The older figures are in
"Tokens and cost".

Seven checks run in every scenario:

- `launcher-only`: every CLI call went through the launcher. A decoy log that cannot be read fails.
- `nothing-pushed`
- `no-stray-commits`: checked-out HEADs are unchanged, and only branches the scenario lists may move.
- `no-attribution`: no agent name or trailer in new commit messages.
- `check-clean`: `check --json` issues are exactly the scenario's `expected_issues`, as a multiset. A document with no `issues` field fails.
- `sources-passive`: the stream canary. It fails when the stream contains the canary text. It passes when a reads list exists and the canary is absent. It is inconclusive when the host exposes no reads list.
- `no-outside-writes`: nothing changed in the run root outside the active workspace, build, host-state, and output directories.

An empty remote, ref, HEAD, or run-digest snapshot fails those checks.
`join` also snapshots the original checkout, under the keys `ws0` and
`ws0/<repo>`, so a commit left there fails `no-stray-commits`, and a file
written there fails `no-outside-writes`.

`join-local-files` requires the workspace binding to resolve to the active
workspace directory and the api binding to resolve to `repositories/api`
inside it. It requires the four role files explorer, planner, worker, and
reviewer for the host that ran when that host is known: `claude` (and the
CLI's `claude-code`) uses `.claude/agents`, `codex` uses `.codex/agents`, and
`cursor` uses `.cursor/agents`. When the check has no host name, or a name
with no role directory, any one of those complete sets is enough.
`execute-worktree` requires a Go function in `internal/invoice` that calls
`encoding/csv`. `complete-note` requires the invoices note to record that a
CSV export exists, and accepts the run's UTC date or the day on either side.

A run ends as `pass`, `fail`, `inconclusive` (the host did not expose the
evidence a check needs, for example tool output; never counted as a pass), or
`harness-error` (setup, auth, or host crash; excluded from a count and rerun).
An unknown check name fails. A scenario's hash covers its file and the
fixtures, so a changed scenario starts its own history.

Each run directory holds the raw stream, the reply, the reads list, snapshots,
`check.json`, `check-before.json`, the files the agent changed, the scenario's
design artifacts, and Git state per repository under `artifacts/git/`. For
`execute`, that also includes `api-diff-stat.txt`, a `git diff --stat` from
`main` to the plan branch, and `api-log.txt`, that branch's commit messages
from `git log` with `--format=%B`. `reproduce.txt` is a date comment plus the
`agent-tests run` line for that scenario, host, and ref. A base-side line also
has `--base`. The line does not include `--yes`. Add `--yes`, and run it as
`go run ./cmd/agent-tests` from `agent-tests/`, when you repeat it:

```sh
go run ./cmd/agent-tests run --scenarios execute --hosts claude --runs 1 --yes
```

The tests under `go test ./...` need no model. For every scenario they build
the real CLI, seed the starting state, replay a scripted reference solution
through the normal grading path and expect `pass`, then damage the end state
in the one way each check exists to catch (a push, an extra commit, an
attribution trailer, a write outside the workspace, the source canary in the
host stream, and so on) and expect `fail`. The slow ones skip under `-short`.
What a scripted pass does not show is in "What it does not prove".

## surface-diff

`surface-diff` compares two git refs without a coding host and without a key.
It is free and offline. For each ref it archives the tree, builds the CLI,
initializes a workspace, and runs one short script: full help and the topics
`init`, `status`, `check`, `member`, `knowledge`, `repo`, `record`, `context`,
`worktree`, `template`, `agent`, `workspace`, and `version`; `record create`
for an intent and a plan; `check` and `status`, each as text and as JSON; and
`agent setup`. The script uses the same fixture names on both sides. When the
two refs are the same commit, it writes the empty message and does not build.

```sh
go run ./cmd/agent-tests surface-diff --base <ref> [--ref HEAD] [--checkout PATH] [--out DIR]
```

`--base` is required. `--ref` defaults to `HEAD`. `--out` defaults to
`agent-tests/out` inside the checkout. That directory is gitignored. The
command exits 0 whether or not anything differs.

It compares three things:

- workspace files, as added, removed, changed, or moved (a removed file and an added file with the same bytes are a move)
- CLI outputs, as a unified diff of the normalized text
- the size of the always-loaded files (`AGENTS.md`, `CLAUDE.md`, `CURSOR.md`, `.cursor/rules/context-circuit.mdc`, `README.md`) and of each skill (`.agents/skills/<name>/SKILL.md`), in bytes and in approximate tokens (bytes/4)

Each side's CLI is stamped with the `CLI_VERSION` in that ref's archive.
Before comparison, the normalizer rewrites the `id`, `created_at`, and
`approved_at` fields the fixed script writes, temporary paths, and that CLI
version. Other ids, dates, and timestamps stay, so an example id in help text
still differs. JSON key sets are compared on the raw stdout.

The report is written under `out/<id>/surface/`. A same-commit run writes only
`surface.json` and `surface.md`. When the refs differ, that directory also
holds `base/`, `branch/`, and `outputs/`. When the files, the outputs, and
the sizes are all equal, the report says `No agent tests needed: surface diff
is empty`.

A non-empty diff names the scenarios, hosts, and artifact kinds from
`triggers.yaml`. Which hosts that means is in "Hosts". A changed path with
no row is reported as "No row: the author picks the closest scenario." When
the change looks larger, the report says so and tags each reason exact or
heuristic. What counts as larger, and who decides, is in "Trigger map".

## report

`report` reads a finished run and writes `report.md`. It does not start a
host. With no `runs.jsonl` under the output directory it exits 1 and says
`report: no run with runs.jsonl under the output directory`.

```sh
go run ./cmd/agent-tests report [--run ID] [--out DIR] [--baseline DIR] [--design] [--upgrade-check] [--no-zip]
```

`--run` is a run id under `--out`, or a path to a run directory. `--out`
defaults to `agent-tests/out` inside the checkout, which is gitignored. The
default run is the newest directory under that folder that contains
`runs.jsonl`. `--baseline` selects the before side. How that directory is
read, and the `(summary)` label, are in "What "before" is".

`report.md` is the Agent tests section: change kind and surface counts when
`surface/surface.json` is present, the scenarios the trigger map required, the
host table, and one row per scenario and host. Several runs of the same
scenario and host are averaged, and the row says how many. A missing token or
cost is `n/a`, never 0. A null total has no percentage. Paste this file into
the Agent tests section. `CONTRIBUTING.md`, under "Changes that affect
agents", describes that section.

A warning in that file is for the pull request to explain. How a warning is
chosen is in "Tokens and cost". It is not a block.

The generated checklist is for the maintainer, who reviews. Its larger-change
checkbox asks for the issue where the design was agreed. Its hosts checkbox
records hosts the contributor lacked. Who runs those hosts is in "Hosts".

`bundle.zip` holds every run's JSON line, stream, reply, and produced files.
Before anything is written, the zip builder rewrites a home prefix to `$HOME`,
a run-root prefix to `$RUN`, and the account name, as a whole word in text and
in entry names, to `$USER`. A longer word that merely contains that name is
left as written. It drops `auth.json`, `.credentials*`, `*.token`, `.netrc`,
`.env`, `id_rsa*`, `*.pem`, and the contents of a per-host auth directory
(`claude` and `codex` config directories, `cursor-home`, and `.claude`,
`.codex`, `.cursor`) before it scans for keys. The host id in
`runs/<scenario>/<host>/` is kept. Symlinks are skipped, and entry names stay
relative. A key inside a dropped auth file does not refuse the zip. If a
key-shaped string remains in anything that would be zipped, the command stops,
names the file, and leaves no zip behind, including a partial one. `--no-zip`
skips the archive. Attach the zip to the pull request.

`--design` pairs one base run with one branch run of the same scenario and
host and writes `design.md`. The run on a side is the passing run with the
median total input tokens, and the lowest run number breaks a tie. When none
passed, the median of the remaining runs is used. The file diffs
`artifacts/files` and `artifacts/design`, summarizes files, headings, and line
counts, diffs `reads.txt` (or says the reads list is not available), quotes
the first 12 lines of each reply, and gives the token change. Paste that
before/after into the Design case.

`--upgrade-check` takes the template delta from `surface/base` and
`surface/branch` under the run directory, or rebuilds those two init
workspaces from the recorded refs. It skips paths a workspace never receives
from the template repository (`README.md`, `LICENSE`, `CHANGELOG.md`,
`CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`, `SECURITY.md`, and `assets/` at the
root), `.context-circuit/ids.yaml`, and a `version` change in `workspace.yaml`.
The rest is applied to a copy of the base tree, and the branch CLI's `check
--json` runs on that copy. The section in `design.md` says this is a
file-level delta, not the host's judgement. It needs a source checkout so the
branch CLI can be built. It does not start a host.

## What it does not prove

Permission prompts. A headless run cannot measure which prompts a person sees.

Whether a real host satisfies the checks. The scripted tests show that the
grader catches a damaged end state. Only a live run shows that a host passes.

The token bands in "Scenarios". They are guesses. See "What "before" is".

A host that exposes no tool output. Cursor records a reads list only when the
stream contains a tool call. Until then the reads list is not available, and
`sources-passive` is inconclusive.

Codex's cached-token split, and whether reasoning tokens are already inside
output tokens. See "Tokens and cost".

Live runs support Unix-like systems. On Windows a timed-out host is stopped,
but its child processes may keep running.
