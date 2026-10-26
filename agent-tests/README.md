# Agent tests

A contributor tool that runs a coding host against a fresh fixture workspace.
It is a separate Go module. It is never shipped, and CI does not run it.

## How to run

From this directory:

```sh
cd agent-tests
go test ./...
go run ./cmd/agent-tests probe
go run ./cmd/agent-tests run --dry-run
```

`probe` runs `<command> --version` for each host (no model is called) and prints its version, whether its pinned model is set, and which
auth variable names are present. `run --dry-run` prints the plan and does not
start a host. A live run needs `--yes`.

## What a run touches

The run root is a temporary directory outside any repository. The checkout is
read through git and is not written. The environment given to a host is an
allowlist. Auth comes only from the environment variables named in
`hosts.yaml`. Nothing is read from or written to the real home, the real
workspace, or user config. The one shared place is the Go module cache: building the
CLI uses the cache that `go env GOMODCACHE` reports, and may download a missing module
into it. Scenario setup commands run as `sh -c` in the fixture workspace, so only run
scenarios from a checkout you trust.

## Scenarios

One fixed prompt each, run in a fresh fixture workspace built by the CLI from the ref under test. State is seeded through the real CLI instead of chaining turns, and the result is graded from what the agent leaves behind: files, Git refs and `check --json`. Nothing depends on exact reply wording; reply checks are patterns, and what a person must still judge is listed in each file's `human` field.

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

The overhead measurement (a fixed reply in an empty directory and in the fixture workspace) is `probe`, not a scenario file.

Seven checks run in every scenario: `launcher-only` (every CLI call went through the launcher), `nothing-pushed`, `no-stray-commits` (checked-out HEADs unchanged, and only branches the scenario lists may move), `no-attribution` (no agent name or trailer in new commit messages), `check-clean` (`check --json` issues are empty, or equal the scenario's `expected_issues`), `sources-passive` (the stale brief under `sources/` was not read) and `no-outside-writes` (nothing changed in the run root outside the workspace, build, host-state and output directories).

A run ends as `pass`, `fail`, `inconclusive` (the host did not expose the evidence a check needs, for example tool output; never counted as a pass) or `harness-error` (setup, auth or host crash; excluded from a count and rerun). A scenario's hash covers its file and the fixtures, so a changed scenario starts its own history.

Each run directory holds the raw stream, the reply, the reads list, snapshots, `check.json`, the files the agent changed, the scenario's design artifacts, Git state per repository under `artifacts/git/`, and `reproduce.txt`, the command that repeats the run:

```sh
go run ./cmd/agent-tests run --scenarios execute --hosts claude --runs 1 --yes
```

The tests under `go test ./...` need no model. For every scenario they build the real CLI, seed the starting state, replay a scripted reference solution through the normal grading path and expect `pass`, then damage the end state in the one way each check exists to catch (a push, an extra commit, an attribution trailer, a write outside the workspace, a read of `sources/`, and so on) and expect `fail`. The slow ones skip under `-short`. Whether a real host satisfies the checks is only shown by a live run.

## What it costs

Every live run uses model usage on the account whose key it is given. A fixed
reply of OK measured 15k to 45k input tokens per host with a developer's
config loaded. A scenario is many times that. Use a spend-limited key. Run
someone else's branch in a container or a VM, because host permissions are
skipped inside the run.

## What it does not prove

Permission prompts. A headless run cannot measure them.

Live runs support Unix-like systems. On Windows a timed-out host is stopped, but its child processes may keep running.
