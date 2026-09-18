# Source workflow

Maintain v2 directly on the active branch under AGENTS.md. The user's request
defines the change. Product instructions apply only to generated workspaces.

## Owners

- product/AGENTS.md.in: shared agent behavior and responsibility boundaries.
- product/docs/: workspace files, agent-facing commands, and worktree guidance.
- internal/workspace/: Go bookkeeping, YAML edits, Git and worktree operations.
- internal/cli/: the command surface and human/JSON output.
- cmd/context-circuit/: executable entry point.
- template/: blank workspace files, embedded with the product instruction.
- scripts/release-manifest.txt: exact source-to-output file mapping.
- assets.go: embeds only product assets and materializes the manifest.
- VERSION: workspace template version.
- CLI_VERSION: independent CLI version; changing either file does not publish.
- product/skills/: packaged CLI installation/update and subagent dispatch skills.
- internal/cow/: native filesystem cloning with independent-copy fallback.
- scripts/: build, checks, and explicitly invoked publication.
- LICENSE: Apache-2.0, covering this checkout and the executable.
- product/LICENSE: 0BSD, covering everything a workspace receives.
- release/template-repo/: landing-page files the published repository owns.

The Go executable handles workspace mechanics. The agent owns interpretation,
planning, implementation, application-specific setup, subagent dispatch, and
user-requested review/delivery. CLI worktree preparation reuses ignored runtime
files through CoW where available. Role settings materialize as native host files.
No product execution state machine or mandatory child-agent workflow is required.

Source context/ is live knowledge about the current product, catalogued by
context/INDEX.md and validated by the knowledge checks the diagnostic runs; it
never ships. sources/, publication/, and release requests are passive maintainer
history and never ship. The remaining source design skill is maintainer tooling
and is excluded from the binary and exported workspace.

release/template-repo/ is a fourth class: not shipped to a workspace and not
history either. Publication restores those files over the extracted artifact and
holds them out of the comparison that decides whether there is anything to
publish, so the published repository keeps its own license, conduct, contributing
and security pages while no workspace ever receives them. A README.md there is
refused, because the published README is the product guide the manifest
assembles.

## This checkout as a workspace

The checkout is itself a schema-2 workspace, so `check`, `status`, and
`context find` run against the real `context/` here instead of a copied fixture.
That is the only reason the workspace files exist. Source work still goes
directly on the branch under AGENTS.md: `intent/` and `plans/` stay empty, and
no ID is allocated for source changes. `.context-circuit/docs/` is deliberately
absent, because `product/docs/` owns that text and a second copy here would
drift from it.

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
