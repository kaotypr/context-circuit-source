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

[Source layout and ownership](#source-layout-and-ownership) says which files own
what. [release/cli/README.md](release/cli/README.md) is the page shipped as the
README inside every CLI archive, so edit it as user documentation.

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
     better; the costs (tokens, results on each host, existing workspaces and
     customized templates, longer output for readers); and the
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
behaves there. For every pull request, CI runs `go test` and `go vet` natively
on Linux, macOS and Windows, `gofmt -l .` on Linux, the filesystem clone tests on
Linux and Windows, and the Change files check. These must pass before it merges.
`check-release.sh` runs on Linux for a pull request into a release branch.
[release/README.md](release/README.md) lists which checks run where.

Tests use temporary directories and disposable Git repositories. The release
check builds and exercises a native binary, verifies the embedded seed
inventory, and cross-compiles the supported binary targets. It also tests
publication guards and commit/tag behavior in disposable fixtures. An optional
new output directory retains the checked release assets. CI runs native tests,
including installer upgrade/rollback fixtures, on Linux, macOS, and Windows.
Cross-compilation alone is not proof of behavior on another operating system.
The runners' default volumes cannot clone, so the native-clone job mounts an XFS
reflink volume on Linux and a ReFS VHDX on Windows and runs
TestNativeCloneOnCapableVolume there with CC_COW_NATIVE_DIR; without that
variable the test skips.

Deterministic tests establish file and Git behavior, not whether a coding host
loads or follows instructions. That is measured by the agent tests under
agent-tests/, which run a real host and grade what it leaves behind, outside CI
and on the contributor's own account; see agent-tests/README.md and
context/maintenance/agent-tests.md. Report host behavior as measured on a named
host, model, and date, and as not exercised when it was not. Do not run a
separate agent review unless the user requests one.

Where this block has been checked:

| Platform | Result |
| --- | --- |
| macOS | both commands run and pass |
| Linux | `gofmt -l .` and `go test ./...` run in CI; `check-release.sh` runs in CI on pull requests into release branches and on pushes to `main` |
| Windows | `go test ./...` and `go vet ./...` run in CI; every script under Git Bash or WSL, and the symlink setting above, untested |

Which scripts are yours to run:

| Script | Who runs it |
| --- | --- |
| `check-release.sh`, `build-dist.sh`, `build-cli.sh` | contributors |
| `check-change-files.sh` | contributors, to see what the Change files check will say: `PR_BODY="$(cat body.md)" sh scripts/check-change-files.sh main HEAD` |
| `release-artifact.sh`, `check-publication.sh`, `check-release-version.sh`, `check-release-notes.sh`, `release-version.sh` | called by the scripts above |
| `release-notes.sh` | the maintainer, in a release-prep pull request |
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

## Source layout and ownership

Maintain v2 directly on the active branch under AGENTS.md. The user's request
defines the change. Product instructions apply only to generated workspaces.

- product/AGENTS.md.in: shared agent behavior and responsibility boundaries.
- product/docs/: workspace files, agent-facing commands, and worktree guidance.
- internal/workspace/: Go bookkeeping, YAML edits, Git and worktree operations.
- internal/cli/: the command surface and human/JSON output.
  - internal/cli/pin.go: the pin switch and guard: hands off to the pinned CLI
    when it is installed, and otherwise refuses a migrating or cross-line change.
- cmd/context-circuit/: CLI entry point.
- template/: blank workspace files, embedded with the product instruction.
- scripts/release-manifest.txt: exact source-to-output file mapping.
- assets.go: embeds only product assets and materializes the manifest.
- VERSION: workspace template version.
- CLI_VERSION: independent CLI version; each file drives its own publication.
- product/skills/: the workspace skills, including cc-system-design, which
  authors design material under a workspace's sources/system-design/ and
  ships from here.
- internal/cow/: native filesystem cloning with independent-copy fallback.
- scripts/: build, checks, and publication.
- release/README.md: how releases are cut, checked, and published.
- release/cli/README.md: the README shipped inside every CLI archive.
- .changelog/: one change file per user-visible change, waiting for the next
  release; .changelog/README.md is the format.
- release/requests/cli/, release/requests/template/: one request per product
  per version, drafted from the change files by scripts/release-notes.sh and
  edited by the maintainer; the request body is that release's notes. A
  template request's `changelog:` says whether the published changelog keeps
  its earlier entries (`keep`, the default), starts over (`reset`), or is left
  alone for a candidate (`skip`).
- LICENSE: Apache-2.0, covering this checkout and the CLI.
- product/LICENSE: 0BSD, covering everything a workspace receives, shown only at
  the published template repository's root.
- product/README.md, product/assets/readme/: the published repository's guide
  and artwork; restored at publication and shipped to no workspace.
- product/workspace-README.md: the front page initialization writes; not
  exported, because a blank seed has no workspace to name.
- product/launcher/: the workspace launcher, which runs exactly the pinned CLI.
- product/host/: host settings the template ships.
- scripts/release-version.sh: the effective version of each product; plain
  here, because release/variant.env is absent here, and a downstream variant
  supplies it with `CC_VERSION_SUFFIX`.
- release/template-repo/: landing-page files and the .gitattributes that keeps
  every repository-owned path out of the archive a project is created from.

The Go CLI handles workspace mechanics. The agent owns interpretation, planning,
implementation, application-specific setup, subagent dispatch, and
user-requested review/delivery. CLI worktree preparation reuses ignored runtime
files through CoW where available. Role settings materialize as native host
files. No product execution state machine or mandatory child-agent workflow is
required.

Source context/ is live knowledge about the current product, catalogued by
context/INDEX.md and validated by the knowledge checks the diagnostic runs; it
never ships. Release requests are passive maintainer history and never ship.
The source skills under .agents/skills/ are maintainer tooling: cc-source-develop
for source maintenance, and audit-cc-skills, which audits the shipped cc-* skills
and edits nothing. They are excluded from the binary and the exported
workspace. .claude/skills is a symlink to that directory, so a host reading
either path finds the same skills.

release/template-repo/ is a fourth class: not shipped to a workspace and not
history either. Publication restores those files over the extracted artifact and
holds them out of the comparison that decides whether there is anything to
publish, so the published repository keeps its own conduct, contributing and
security pages while no workspace ever receives them. The guide, its artwork and
the license are restored the same way from product/, so a README.md or LICENSE
defined in template-repo/ is refused: those have one source each.

A dist/ build is therefore the workspace only, not a preview of the published
repository: everything in this class joins the tree at publication.

The .gitattributes restored with them marks every repository-owned path
export-ignore. A workspace created from a template tag is extracted with
`git archive`, which honors that file, so it starts with no front page about
Context Circuit, no release notes for the template, and no LICENSE at its root.
A tag-to-tag sync reads the git tree rather than an archive and cannot see the
mark, so product/AGENTS.md.in names the same paths for the agent performing
that update.

agent-tests/ is a fifth class: a separate Go module that is not shipped and not
history either. It is a contributor tool that runs coding hosts against a
fixture workspace, so it is excluded from every build, archive, and the
manifest, and CI does not run it.

## Branches and pull requests

- Branch from `main` and open your pull request against it. Release branches,
  `release/X.Y`, take only the maintainer's release-prep and cherry-pick pull
  requests.
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
  command. The bug report asks for the CLI and template versions, your operating
  system and coding host, what you did, what you expected and what happened, and
  an excerpt of the workspace `check` output.
- For a bug fix, documentation or tests, open a pull request straight away.
- Usage questions go to
  [GitHub Discussions, Q&A](https://github.com/kaotypr/context-circuit-source/discussions/categories/q-a).
  Issues are for defects.
- Report a vulnerability as [SECURITY.md](SECURITY.md) says, never in a public
  issue or discussion.

## Review and merge

- @kaotypr is the only maintainer, and the only person who merges.
- The repository allows merge commits only, so your commits land as you wrote
  them.
- CI must pass, and no approval count is required. GitHub enforces the test
  checks and the Change files check on `main`.

Publication, tagging, and deployment are separate maintainer acts and never part
of a pull request.

## Release notes

- If your change is one someone using Context Circuit would notice, add a
  change file, `.changelog/<slug>.md`, in the same pull request.
  [.changelog/README.md](.changelog/README.md) is the format: which product,
  what kind of change, an optional upgrade note, and the entry itself.
- If users see no change, add no file, and write `None` under "Release note"
  in the pull request description. The Change files check fails a pull request
  that changes CLI or template files with neither.
- Do not edit `release/requests/`, `VERSION` or `CLI_VERSION`. The maintainer
  drafts the requests from the change files when a release is cut.
- The release-note format and the compatibility policy are in
  [release/README.md](release/README.md).

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

This checkout is not a workspace. It carries no `workspace.yaml`, no roster, no
`.context-circuit/`, and no records; the product's lifecycle is for generated
workspaces, and AGENTS.md bars it here.

`context/` is still checked by the product's own diagnostic.
`internal/cli/knowledge_tree_test.go` builds a throwaway workspace, copies the
real tree into it, and binds it to this checkout so each note's code anchors
resolve against real history. Findings split by the product's own taxonomy: one
a command or an edit discharges fails the test, and one marked `needs a person:`
— a note to re-read against code that moved — is reported for a maintainer
rather than blocking a contributor on somebody else's judgment.

The job that runs `go test` clones full history, because the reviewed-date pass
counts commits under those anchors and a shallow clone would make it pass
without looking.
