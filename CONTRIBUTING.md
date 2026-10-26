# Contributing to Context Circuit

This guide is for you and for the coding agent helping you. It says what you are
changing when you change Context Circuit, how to set up and validate a change,
and how a pull request gets merged.

## What you are changing

Context Circuit is three things that work together: an agent harness (the
instructions and skills a coding agent reads), a standard file and folder system
(the workspace layout and its records), and CLI tools. Together they make coding
agents take determined, standardized actions. It is not merely code.

So a change to a skill, an instruction file, the workspace layout, or what a CLI
command prints is a product-design change, even when the diff is a few words. A
reworded sentence in a skill changes what every agent reading it does next. Treat
it with the same care as a change to Go code, and say in the pull request which
behaviour it moves and why.

This repository is the source for two separately versioned products: the
workspace template under `product/` and `template/`, and the native CLI under
`cmd/` and `internal/`. Changes to either are made here, never by patching a
generated workspace.

Before you write code, decide which side of the seam your change belongs on. The
[README](README.md#what-belongs-in-the-cli) states the rule: anything whose
answer depends on reading the actual project belongs to the coding agent, and
anything that must return the same result every time belongs to Go. The CLI
holds no model credentials and never calls a model API, so a feature needing
judgment is a skill or an instruction change rather than a command.

[WORKFLOW.md](WORKFLOW.md) describes which files own what. [CLI.md](CLI.md) is
the page shipped as the README inside every CLI archive, so edit it as user
documentation.

## Setting up

You need:

- Go 1.25 or newer, and Git.
- A POSIX `sh` with `tar`, `gzip`, `zip`, `unzip`, and `sha256sum` or `shasum`,
  plus the usual `awk`, `sed`, `cmp`, `find`, `diff` and `mktemp`. The scripts
  under `scripts/` use them.
- **A full clone.** The knowledge test counts commits under each note's code
  anchors. A shallow clone leaves it nothing to count, so it passes without
  looking. Clone normally, without `--depth`.

On macOS the system tools are enough. On Linux, install `zip` and `unzip` if
your distribution does not have them.

On **Windows**, `go test ./...` runs natively. The scripts need Git Bash or WSL,
and Git Bash has no `zip`, so `scripts/check-release.sh` cannot finish there.
WSL can have every tool the scripts need once Go, `zip` and `unzip` are installed
there.

`.claude/skills` is a symlink to `.agents/skills`. On Windows, Git checks a
symlink out as a small text file unless symlinks are allowed (Developer Mode)
and `core.symlinks=true` is set before cloning:

```sh
git clone -c core.symlinks=true <repository URL>
```

Only a Claude Code session that reads the skills needs this. The tests do not.

## Changes that affect agents

A change affects agents when an agent in a generated workspace could read
different files, find a different layout, or see different CLI output after it.
[`agent-tests/README.md`](agent-tests/README.md) has the commands and the full
guide. This is not CI and not a release gate: you run it before you open the
pull request, on your own coding-host accounts, and you pay for the runs.

1. Run `surface-diff` from `agent-tests/`, with `--base` set to the branch your
   pull request targets:

   ```sh
   cd agent-tests
   go run ./cmd/agent-tests surface-diff --base <target branch>
   ```

   If the report says `No agent tests needed: surface diff is empty`, nothing
   more is needed.
2. If you changed a shipped skill, also read the maintainer skill
   `audit-cc-skills` under `.agents/skills/`. It is a free, static check of the
   skills that starts no coding host. The agent tests below are the check that
   runs a host.
3. If the change is larger, open an issue with the design case first, using the
   [design case form](https://github.com/kaotypr/context-circuit-source/issues/new?template=design-case.yml),
   and wait for the maintainer's answer before you pay for any agent runs. A
   change is larger when it adds, removes or moves shipped files; changes the
   shape of what a skill produces; adds or changes a CLI action or a `check`
   rule; touches a gate such as approval, commit or push; or needs existing
   workspaces to migrate.
4. Run the scenarios `surface-diff` names, once each, on one coding host of your
   choice, plus the specific host for any host-specific file. A change to the
   shipped file list or layout runs every scenario on one host and `orient` on
   all three. When your host
   version and model match a committed baseline, that baseline is your "before".
   Otherwise run the base and your branch alternately on the same machine. Run
   someone else's branch only inside a container or virtual machine.
5. In the pull request, write the design case first, then the agent tests:
   - **Design case:** what an agent will now read, do or produce differently;
     what the current design gets wrong, with evidence; a before and after of the
     agent-visible output from the same scenario, host and model; what is
     better; the catch and who pays (tokens, compliance on each host, existing
     workspaces, forks of the template, readers, maintenance); and the
     alternatives you considered, including leaving it as it is.
   - **Agent tests:** the scenarios, hosts and results. From `agent-tests/`,
     `go run ./cmd/agent-tests report` writes the tables and
     `go run ./cmd/agent-tests report --design` writes the before and after.
     You write the reasoning. Attach the cleaned archive of all runs.

A failed check, a lower result or more than 15% more input tokens is a warning,
not a block: explain it in the pull request and the reviewer decides. The
reviewer may ask for more runs, scenarios or hosts. If you do not have a host the
change needs, say so, and the maintainer, who reviews, runs it.

## Validating

Run this from the repository root:

```sh
gofmt -l .
sh scripts/check-release.sh "$(mktemp -d)/out"
```

`gofmt -l .` prints every Go file that is not formatted; it should print nothing.
`check-release.sh` already runs `go test ./...` and `go vet ./...`, then builds
the workspace archive and all six CLI targets, verifies checksums, and exercises
installation, initialization and the publication safeguards. Its output
directory must not exist yet; the script never overwrites one.

Cross-compilation proves a binary builds for another platform, not that it
behaves there. For a pull request into `development/v2`, CI runs `go test` and
`go vet` natively on Linux, macOS and Windows, `gofmt -l .` on Linux, and the
filesystem clone tests on Linux and Windows. These must pass before it merges.
`check-release.sh` runs on Linux for a pull request into `main`, which is the
release pull request from `development/v2`. [RELEASING.md](RELEASING.md) lists
which checks run where.

Where this block has been checked:

| Platform | Result |
| --- | --- |
| macOS | both commands run and pass |
| Linux | `gofmt -l .` and `go test ./...` run in CI; `check-release.sh` runs in CI on pull requests into `main` |
| Windows | `go test ./...` and `go vet ./...` run in CI; every script under Git Bash or WSL, and the symlink setting above, untested |

Which scripts are yours to run:

| Script | Who runs it |
| --- | --- |
| `check-release.sh`, `build-dist.sh`, `build-cli.sh` | contributors |
| `release-artifact.sh`, `check-publication.sh`, `check-release-version.sh`, `release-version.sh` | called by the scripts above |
| `resolve-release-version.sh`, `publish-template.sh`, `publish-gitlab-release.sh` | CI and release only (the checks run them against fixtures) |

A change to shipped instructions or documentation needs the same care as code.
The test suite enforces which layer a rule belongs in and which audience a
shipped document is written for. If a documentation test fails, the placement is
the finding, not the test.

## Making a change

- Edit the source trees, not a generated workspace. `scripts/release-manifest.txt`
  maps every shipped file to its destination; keep it and `assets.go` consistent
  with anything you add or remove.
- Preserve unrelated working-tree changes.
- Use `type(scope): imperative subject` for commit subjects, with the changed
  component as the scope. Describe what the change decides, not only what it
  edits: a reviewer needs to see which behaviour moved and why.

## Branches and pull requests

- Branch from `development/v2` and open your pull request against it.
- `main` stays the default branch. If you open a pull request against `main`,
  the maintainer retargets it to `development/v2`.
- Use a plain branch name that says what the change is, for example
  `fix/short-topic` or `docs/install-steps`.
- A pull request description stands on its own: what changed, why, and which
  checks you ran. A reader should not need anything outside the pull request to
  follow it.
- The [pull request template](.github/pull_request_template.md) asks for these,
  the linked issue, the release note and a short checklist, and, for a change
  that affects agents, the design case and the agent tests.

## Issue or pull request first

- For a behaviour change or a new command, open an issue first, so the design
  can be agreed before you build it. The same goes for a larger change to what
  agents read or see, even when the diff is only documentation; see "Changes
  that affect agents".
- [New issue](https://github.com/kaotypr/context-circuit-source/issues/new/choose)
  offers two forms: a bug report and a proposal for a behaviour change or new
  command.
- For a bug fix, documentation or tests, open a pull request straight away.
- Usage questions go to
  [GitHub Discussions, Q&A](https://github.com/kaotypr/context-circuit-source/discussions/categories/q-a).
  Issues are for defects.

## Review and merge

- @kaotypr is the only maintainer, and the only person who merges.
- The repository allows merge commits only, so your commits land as you wrote
  them.
- CI must pass, and no approval count is required. GitHub enforces the test
  checks on `development/v2`.

Publication, tagging, and deployment are separate maintainer acts and never part
of a pull request.

## Release notes

- Give your pull request description a paragraph headed "Release note" that says
  what changes for someone using Context Circuit. If nothing does, say so in one
  line.
- Do not edit `release/requests/`, `VERSION` or `CLI_VERSION`. The maintainer
  prepares those when a release is cut.
- The release-note format and the compatibility policy are in
  [RELEASING.md](RELEASING.md).

## Licensing

There is no DCO and no CLA. A contribution comes in under section 5 of the
[Apache License 2.0](LICENSE), which covers this repository and the CLI. Files a
workspace receives, under `product/` and `template/`, are under
[0BSD](product/LICENSE), and a contribution to them comes in under 0BSD. The
[README](README.md#licensing) explains why there are two licenses.

## AI assistance

You may use a coding agent. Commits and pull requests carry no AI-assistance
credit, no co-author trailer for an agent, and no "generated by" text. This rule
is for everyone. If your tool adds such text on its own, remove it before you
push.

## Knowledge

`context/` holds the product's knowledge: architecture, decisions, domain rules
and terminology, catalogued in `context/INDEX.md`. Retrieve it through that
index rather than scanning, and keep a note and its catalog entry in the same
change.

A behaviour change updates `context/` in the same pull request.

`go test` runs the knowledge checks over `context/`. A finding that only a person
can settle, one the product's `check` command marks `needs a person:`, is logged
and does not fail the test. To see these findings, run:

```sh
go test -count=1 -v -run TestTheProjectsOwnKnowledgeIsValid ./internal/cli
```

Each one reads `needs a maintainer's judgment: …`. Most often it says a note was
last confirmed on some date and the code it is anchored to has commits since. It
asks a person to read the note against that code, correct the note and its
catalog entry if they are wrong, and move the note's `reviewed` date in
`context/INDEX.md`. The maintainer handles these, so one does not block your pull
request.
