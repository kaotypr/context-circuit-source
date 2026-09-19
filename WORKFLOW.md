# Source workflow

Maintain v2 directly on the active branch under AGENTS.md. The user's request
defines the change. Product instructions apply only to generated workspaces.

## Owners

- product/AGENTS.md.in: shared agent behavior and responsibility boundaries.
- product/docs/: workspace files, agent-facing commands, and worktree guidance.
- internal/workspace/: Go bookkeeping, YAML edits, Git and worktree operations.
- internal/cli/: the command surface and human/JSON output.
- cmd/context-circuit/: CLI entry point.
- template/: blank workspace files, embedded with the product instruction.
- scripts/release-manifest.txt: exact source-to-output file mapping.
- assets.go: embeds only product assets and materializes the manifest.
- VERSION: workspace template version.
- CLI_VERSION: independent CLI version; each file drives its own publication.
- product/skills/: the workspace skills, including cc-system-design, which
  authors design material under sources/system-design/ and ships from here.
- internal/cow/: native filesystem cloning with independent-copy fallback.
- scripts/: build, checks, and publication.
- release/requests/cli/, release/requests/template/: one request per product
  per version; the request body is that release's notes.
- LICENSE: Apache-2.0, covering this checkout and the CLI.
- product/LICENSE: 0BSD, covering everything a workspace receives, shown only at
  the published template repository's root.
- product/README.md, product/assets/readme/: the published repository's guide
  and artwork; restored at publication and shipped to no workspace.
- product/workspace-README.md: the front page initialization writes; not
  exported, because a blank seed has no workspace to name.
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
never ships. sources/ and release requests are passive maintainer history and
never ship. The source skills under .agents/skills/ — cc-source-develop and
cc-system-design — are maintainer tooling, excluded from the binary and the
exported workspace. cc-system-design states the same design rubric that
product/skills/cc-system-design/SKILL.md ships to a workspace, differing only
where a maintainer checkout genuinely differs: the version source, how knowledge
is written, and the absence of records. Change the layout, altitude, splitting,
or scope-separation rules in one and change them in the other. .claude/skills is
a symlink to that directory, so a host reading either path finds the same
skills.

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
export-ignore. tembiter creates a project with `git archive`, which honors that
file, so a new workspace starts with no front page about Context Circuit, no
release notes for the template, and no LICENSE at its root. A tag-to-tag sync
reads the git tree rather than an archive and cannot see the mark, so
product/AGENTS.md.in names the same paths for the agent performing that update.

## How context/ is validated

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

It was registered as a workspace before this, which put the same checks behind
two gitignored files: they ran on one machine and never in CI, so a fresh clone
got silence instead of a warning. The job that runs `go test` clones full
history, because the reviewed-date pass counts commits under those anchors and a
shallow clone would make it pass without looking.

## Validation

Run gofmt on changed Go files, go test ./..., go vet ./..., and
sh scripts/check-release.sh. Tests use temporary directories and disposable Git
repositories. The release check builds and exercises a native binary, verifies
the embedded seed inventory, and cross-compiles the supported binary targets.
It also tests publication guards and commit/tag behavior in disposable fixtures.
An optional new output directory retains the checked release assets.
CI runs native tests, including installer upgrade/rollback fixtures, on Linux, macOS, and Windows. Cross-compilation alone is not
proof of behavior on another operating system.

Validate Markdown scenarios against the shared entry instruction. Deterministic
tests establish file and Git behavior, not whether a coding host loads or follows
instructions. Report host behavior as unverified unless actually exercised. Do
not revive the retired v1 acceptance harness or run a separate agent review
unless the user requests one.
