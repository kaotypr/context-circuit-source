# Releasing

@kaotypr is the only maintainer and the only person who merges. Contributors open pull requests into `development/v2`. The repository allows merge commits only, so every merge below is a merge commit. [CONTRIBUTING.md](../CONTRIBUTING.md) is the same rule for contributors.

## Release flow

1. The maintainer cuts `development/X.Y.Z` from `development/v2`.
2. The maintainer merges the work pull requests for that release into `development/X.Y.Z`. Contributors open their pull requests into `development/v2`.
3. The maintainer opens a release-prep pull request. It sets `VERSION`, `CLI_VERSION`, or both, and adds the release request at `release/requests/cli/X.Y.Z.md`, `release/requests/template/X.Y.Z.md`, or both. [Publish](../.github/workflows/publish.yml) runs on a push to `main` that touches `release/requests/**`, `VERSION`, or `CLI_VERSION`. The version that publishes is the text of the version file, unless a manual run names one, and `scripts/resolve-release-version.sh` requires `release/requests/<cli|template>/<version>.md` when that version is not out yet.
4. The maintainer opens a pull request from `development/X.Y.Z` into `development/v2` and merges it.
5. The maintainer opens a pull request from `development/v2` into `main` and merges it.
6. The merge to `main` runs Publish when those paths changed. The CLI job runs first, then the template job.

Before the release-prep pull request, from the repository root:

```sh
gofmt -l .
sh scripts/check-release.sh "$(mktemp -d)/out"
```

`gofmt -l .` prints nothing when the Go files are formatted. The directory passed to `check-release.sh` must not exist yet. The script runs `go test ./...` and `go vet ./...`, then builds both products and runs the publication and version checks.

### Cutting 2.5.1

| Product | File to set | Release request to add |
| --- | --- | --- |
| CLI | `CLI_VERSION`, contents `2.5.1` | `release/requests/cli/2.5.1.md` |
| Template | `VERSION`, contents `2.5.1` | `release/requests/template/2.5.1.md` |

Set both rows to release both products. Set one row to release that product only. The other job resolves whatever its version file still says, and when that release is already out it does nothing. Run the two checks above first. The note in the new request follows [Release notes](#release-notes).

## Checks

[Check](../.github/workflows/check.yml) decides by the branch a pull request targets which jobs run. A job that does not apply is skipped, and a skipped job counts as passed.

| Pull request into | Go tests and Filesystem clone tests | Release build | Required to merge |
| --- | --- | --- | --- |
| `development/v2` | run | skipped | the five test jobs: Go tests on Ubuntu, macOS and Windows, Filesystem clone tests on Ubuntu and Windows |
| `development/X.Y.Z` | run | skipped | nothing; the results are advisory |
| `main` | skipped | runs | Release build |

A push to `development/v2` or `development/X.Y.Z` runs the tests and skips Release build. A push to `main` runs Release build and skips the tests. The pull request into `main` is the release pull request from `development/v2`, whose commits already ran the tests there, so it runs only Release build. A newer push to a pull request cancels the older run.

Only `main` and `development/v2` are protected branches: a pull request is required, force-pushes and deletion are blocked, and the repository administrator can bypass the rules. `development/X.Y.Z` branches are not protected.

## Branches

`development/X.Y.Z` is not protected, so it can be deleted without bypassing anything. The repository has a setting that can delete a pull request's branch when it merges. It is off while 2.5.0 is released and will be turned on again afterwards. When it is on, merging `development/X.Y.Z` into `development/v2` deletes it, and its commits are then on `development/v2`. Topic branches are deleted on merge in the same way when the setting is on. While it is off, the maintainer deletes branches by hand.

## Agent-test baseline

Once the agent-test baseline exists under `agent-tests/baselines/`, refresh it after each template release and whenever a pinned model or host version changes; run the previous template on the new host or model first, to separate model drift from template change.

## Release notes

A release request is `release/requests/<cli|template>/X.Y.Z.md`. The 2.5.0 requests keep their current headings. Starting with the first release after 2.5.0, the body uses these headings, in this order:

1. Upgrade notes
2. Breaking
3. Added
4. Fixed
5. Contributors

Each change links its pull request. A CLI request on this repository uses `#N`. A template request uses `kaotypr/context-circuit-source#N`, as `release/requests/template/2.5.0.md` does, because the note is published on the template repository. Contributors are credited by GitHub handle.

A heading with nothing under it is omitted. Upgrade notes stays, and says `None.` when it is empty.

| Field | Who reads it | Values |
| --- | --- | --- |
| `changelog:` | `scripts/publish-template.sh`, on a template request | `keep` or `reset`. Omitted means `keep`. |
| `level:` | A person. No file under `scripts/` or `.github/` reads it. | Existing requests use `stable`, `minor`, and `prerelease`. Nothing checks the word. |
| `reason:` | A person. No file under `scripts/` or `.github/` reads it. | One sentence for the person cutting the release. |

`keep` prepends the new entry to the published template changelog. `reset` starts that changelog over; `release/requests/template/2.0.0.md` is the request that sets it. Any other value fails template publication. The CLI half of Publish strips the frontmatter and publishes the body. It leaves `changelog:` unread.

A note names the change and links the pull request. It carries no plan or intent ID, no record path, and no private workspace name. It carries no AI-assistance credit.

## Compatibility

| Surface | Breaking | Not breaking |
| --- | --- | --- |
| CLI | A removed or renamed command or flag. A changed exit code. A removed or renamed JSON field. A new refusal of something that used to work. | A new command. A new optional flag. A new JSON field. Clearer output text. |
| Template | A change that makes an existing workspace's records invalid, or that needs its owner to do something to keep working. | New skill text. A reworded instruction. A new optional file. |
| Workspace schema | Any change that needs a migration. The schema number goes up. | Documentation of existing fields. |

[product/docs/commands.md](../product/docs/commands.md) is the command surface: `--json` for structured fields, and a nonzero exit status on errors (an unknown flag exits 2). [release/cli/README.md](cli/README.md) is the package page. CLI v2 writes workspace schema 3 and still reads schema 2. [product/docs/workspace.md](../product/docs/workspace.md) lists the files. The schema number is `version` in `workspace.yaml`. Schema 3 is the per-member ID layout. Moving a schema 2 workspace there is a migration, and one person on the team has to make that split, commit it, and push it before anyone else changes anything. Writing down a field that is already in the layout leaves the number unchanged.

A breaking change ships in a new major version and is listed under Breaking in its release note. Separately, the pin guard protects a workspace from the wrong CLI. It applies when the running CLI is not the pinned one and it did not hand off. It refuses a writing change that would move the workspace to schema 3, `init` included, and it refuses any writing change when the running CLI and the pin are on different release lines. [internal/cli/release_line.go](../internal/cli/release_line.go) calls `X.Y.Z` stable and `X.Y.Z-<label>.N` a variant line for that label. Two stable versions share the stable line, so a new major is still stable: the guard does not refuse it for the number alone. A version of neither shape is an unknown line. Reserved labels (`alpha`, `beta`, `rc`, `dev`, `test`, `pre`, `preview`, `snapshot`) keep a candidate on the unknown line. An unknown line is not a different line by itself; the schema refusal is what still applies. `CC_CLI_VERSION=local` never lifts a refusal. Read-only commands, a preview such as `record delete` without `--confirm`, and other changes on the same release line only warn. The 2.5.0 CLI request states this refusal.

Nothing in this repository deprecates a feature today. From here on, a feature is deprecated in the notes of one minor release, under Upgrade notes, before it is removed in a later release.

Which versions receive fixes is in [SECURITY.md](../SECURITY.md).

### Variant forks

A published stable tag is left as published. `scripts/resolve-release-version.sh` treats a version that already has a release as out and does not publish it again. Template publication stops when the tag `v<version>` already exists.

`X.Y.Z-<label>.N` is reserved for forks. Apart from candidates under a reserved label such as `alpha`, this repository publishes no version of that form. Publish does not set `CC_VERSION_SUFFIX`, and this tree has no `release/variant.env`, so a public checkout prints `CLI_VERSION` or `VERSION` unchanged. A fork adds `release/variant.env`, which is the counter file `scripts/release-version.sh` reads, and sets `CC_VERSION_SUFFIX` to the label. The label is 1–20 lowercase letters, and it is not a reserved label. The base it builds on is a stable `2.Y.Z`. The composed version is `<base>-<label>.<counter>`. With `CC_VERSION_SUFFIX` set, naming a different version on a manual run is refused. A variant CLI installs as `<label>-context-circuit-cli`. `CC_VERSION_SUFFIX` carries the label.

## Publish

[Publish](../.github/workflows/publish.yml) runs three jobs:

| Job | When it runs | What it does |
| --- | --- | --- |
| `gate` | First, on Linux, macOS, and Windows | `go test ./...` and `go vet ./...` |
| `cli` | After `gate` | Publishes the CLI when `CLI_VERSION`, or the named version, is not already out |
| `template` | After `gate` and `cli` | Publishes the template when `VERSION`, or the named version, is not already out |

The template job is second because the template release names the CLI release. When the template job is going to publish, it runs `scripts/check-release.sh` before it assembles the template.

A product is published only when its version file names a release that is not out yet. `workflow_dispatch` can name a version: `cli_version` blank means `CLI_VERSION`, and `template_version` blank means `VERSION`. A version that already has a release is left alone, and the other product can still publish.

Publication is the maintainer's explicit act. The [README](../README.md#release-assembly) says commits, tags, pushes, release publication, and deployment require an explicit maintainer request, and are never part of validation. Merging the version files and requests to `main`, or starting Publish by hand, is that act. A pull request does not publish.
