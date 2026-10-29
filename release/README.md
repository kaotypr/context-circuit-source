# Releasing

@kaotypr is the only maintainer and the only person who merges. Contributors open pull requests into `main`. The repository allows merge commits only, so every merge below is a merge commit. [CONTRIBUTING.md](../CONTRIBUTING.md) is the same rule for contributors.

## Release flow

`main` is the trunk and never publishes. Each minor release line has a release branch, `release/X.Y`, which publishes that minor and its patches.

1. Work merges into `main`. A pull request that changes something users see adds a change file under `.changelog/`; see [Release notes](#release-notes).
2. At a minor release's freeze, the maintainer cuts `release/X.Y` from `main`. The agent-test baseline runs there while `main` keeps taking work. Cutting the branch publishes nothing, because its version files still name releases that are already out.
3. A fix for the release lands on `main` first. The maintainer then opens a cherry-pick pull request into `release/X.Y`, which carries the fix's change file with it.
4. The maintainer opens a release-prep pull request into `release/X.Y`. It sets `VERSION`, `CLI_VERSION`, or both, and adds the release request of each product it publishes, `release/requests/cli/X.Y.Z.md`, `release/requests/template/X.Y.Z.md`, or both. Drafting and editing those requests is in [Release notes](#release-notes).
5. Merging the release-prep pull request runs [Publish](../.github/workflows/publish.yml) on `release/X.Y`. The CLI job runs first, then the template job. The version that publishes is the text of the version file, unless a manual run names one, and `scripts/resolve-release-version.sh` requires `release/requests/<cli|template>/<version>.md` when that version is not out yet. It also refuses a version that is not on its own line's branch, such as `2.6.0` on `release/2.5`, or anything on `main`.
6. After a stable release, the maintainer merges `release/X.Y` back into `main` from a topic branch cut from `release/X.Y`, so delete-on-merge removes the topic branch and not the release branch. `main` then carries the release's version files and requests, and the removal of the change files it shipped.

A release candidate, such as `2.5.0-rc.1`, follows steps 4 and 5 on the same branch and is published as a prerelease. A patch release, such as `2.5.1`, follows steps 3 to 6 on the existing `release/X.Y`.

Before the release-prep pull request, from the repository root:

```sh
gofmt -l .
sh scripts/check-release.sh "$(mktemp -d)/out"
```

`gofmt -l .` prints nothing when the Go files are formatted. The directory passed to `check-release.sh` must not exist yet. The script runs `go test ./...` and `go vet ./...`, then builds both products and runs the publication, version, and release-note checks.

### Cutting 2.5.1

On a branch from `release/2.5`:

| Product | File to set | Release request to add |
| --- | --- | --- |
| CLI | `CLI_VERSION`, contents `2.5.1` | `sh scripts/release-notes.sh cli 2.5.1`, then edit `release/requests/cli/2.5.1.md` |
| Template | `VERSION`, contents `2.5.1` | `sh scripts/release-notes.sh template 2.5.1`, then edit `release/requests/template/2.5.1.md` |

Set both rows to release both products. Set one row to release that product only. The other job resolves whatever its version file still says, and when that release is already out it does nothing. Run the two checks above first, and open the pull request into `release/2.5`.

## Checks

[Check](../.github/workflows/check.yml) runs on every pull request, and on pushes to `main` and `release/**`. A job that does not apply is skipped, and a skipped job counts as passed.

| Pull request into | Go tests and Filesystem clone tests | Release build | Change files |
| --- | --- | --- | --- |
| `main` | run | skipped | runs |
| `release/X.Y` | run | runs | runs |

A push to `main` or `release/**` runs the tests and Release build. A newer push to a pull request cancels the older run. [Change files](../.github/workflows/change-files.yml) is its own workflow, so editing a pull request's description re-runs it without cancelling the tests. It is skipped for Dependabot.

`main` and `release/*` are protected branches. A pull request is required, and six checks are required to merge: Go tests on Ubuntu, macOS and Windows, Filesystem clone tests on Ubuntu and Windows, and Change files. Force-pushes and deletion are blocked, and the repository administrator can bypass the rules.

## Branches

The repository deletes a pull request's branch when it merges. Topic branches go that way. A release branch is never the head of a pull request: the merge back into `main` comes from a topic branch, as step 6 says. Release branches stay for their line's patches.

## Agent-test baseline

The agent-test baseline is committed under `agent-tests/baselines/`, one directory per template version. Refresh it after each template release and whenever a pinned model or host version changes; run the previous template on the new host or model first, to separate model drift from template change.

## Release notes

A release request is `release/requests/<cli|template>/X.Y.Z.md`. It is drafted from change files and edited by the maintainer.

**Change files.** A pull request that changes something users see adds `.changelog/<slug>.md`, whose frontmatter names its `product` (`cli`, `template`, or `both`), its `kind` (`added`, `fixed`, `changed`, or `breaking`), an optional one-line `upgrade` note, and optional `authors`. Its body is the entry. [.changelog/README.md](../.changelog/README.md) is the format. A pull request that changes nothing users see writes `None` under Release note instead. The Change files check, `scripts/check-change-files.sh`, fails a pull request that changes CLI or template files with neither. For this check, the template is `product/`, `template/` and `release/template-repo/`, and the CLI is `cmd/`, `internal/`, `assets.go`, `go.mod`, `go.sum`, `release/cli/` and `scripts/build-cli.sh`, without Go tests and testdata.

**Drafting.** In the release-prep branch, `sh scripts/release-notes.sh <cli|template> <version>` writes the request from every change file for that product. The body uses these headings, in this order:

1. Upgrade notes
2. Breaking
3. Added
4. Fixed
5. Changed
6. Contributors

A heading with nothing under it is omitted. Upgrade notes stays, and says `None.` when it is empty. Each entry links the pull request whose merge brought its file to the release branch: `#N` in a CLI request, and `kaotypr/context-circuit-source#N` in a template request, because that note is published on the template repository. A cherry-picked fix therefore links its cherry-pick pull request, unless its change file names the original fix in `pr:`. A file that reached the branch without a pull request merge gets no link, and the script says so. Contributors are each file's `authors`, or else the account its pull request came from, by GitHub handle, without bot accounts.

The draft is a starting point. Edit it in the release-prep pull request as freely as a request written by hand: fill in `reason:`, merge entries, or rewrite them. What merges is what publishes.

**A release note describes the difference from the previous release, not the pull requests in between.** Read the draft against the last release a user had, not against the merge history. An entry a later change made irrelevant is dropped: a feature reworked before release keeps one entry describing where it ended, a change added and reverted keeps none, and a fix for something that never shipped is not a fix. A pull request that reworks or undoes an unshipped change file is expected to edit or delete that file itself, as [.changelog/README.md](../.changelog/README.md) asks, so most of this is already done when the draft is made; the release-prep edit is where what remains is settled. A file's `pr:` field names the pull requests an entry describes when they are not the one that added it, for example a change built in several pull requests, or a cherry-pick naming the original fix. `--stdout` prints a draft and changes nothing; `--force` drafts over an existing request.

**Candidates and stable releases.** A candidate is a version whose label is a reserved one, such as `2.5.0-rc.1`. Its draft leaves the change files in place, so the stable release's draft covers everything since the last stable release again, and its template request says `changelog: skip`. A stable release's draft removes the change files it used, with `git rm`, in the same pull request. A `both` file is rewritten to name the other product, so that product's next stable release still lists it. A variant fork's `-<label>.N` is a release of that variant, not a candidate.

| Field | Who reads it | Values |
| --- | --- | --- |
| `changelog:` | `scripts/publish-template.sh`, on a template request | `keep`, `reset`, or `skip`. Omitted means `keep`. |
| `level:` | A person. No file under `scripts/` or `.github/` reads it. | The draft writes `stable` or `prerelease`. Older requests also use `minor`. Nothing checks the word. |
| `reason:` | A person. No file under `scripts/` or `.github/` reads it. | One sentence for the person cutting the release. |

`keep` prepends the new entry to the published template changelog. `reset` starts that changelog over; `release/requests/template/2.0.0.md` is the request that sets it. `skip` leaves the changelog as it is, so a candidate's notes appear on its release alone. Any other value fails template publication. The CLI half of Publish strips the frontmatter and publishes the body. It leaves `changelog:` unread.

The 2.5.0 requests were written by hand before change files existed, and keep their own headings.

A note names the change and links the pull request. It carries no plan or intent ID, no record path, and no private workspace name. It carries no AI-assistance credit. A change file is published text too.

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

Publication is the maintainer's explicit act. The [README](../README.md#release-assembly) says commits, tags, pushes, release publication, and deployment require an explicit maintainer request, and are never part of validation. Merging a release-prep pull request into `release/X.Y`, or starting Publish by hand on that branch, is that act. A pull request does not publish.
