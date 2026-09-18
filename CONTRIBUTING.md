# Contributing to Context Circuit

This repository is the maintainer source for two separately versioned products:
the workspace template under `product/` and `template/`, and the native CLI under
`cmd/` and `internal/`. Changes to either are made here and never by patching a
generated workspace.

## Before you write code

Decide which side of the seam your change belongs on. The
[README](README.md#what-belongs-in-the-executable) states the rule: anything whose
answer depends on reading the actual project belongs to the coding agent, and
anything that must return the same result every time belongs to Go. The executable
holds no model credentials and never calls a model API, so a feature needing
judgment is a skill or an instruction change rather than a command.

[CLI.md](CLI.md) describes the executable boundary and packaging.
[WORKFLOW.md](WORKFLOW.md) describes source ownership and validation.

## Making a change

- Edit the source trees, not a generated workspace. `scripts/release-manifest.txt`
  maps every shipped file to its destination; keep it and `assets.go` consistent
  with anything you add or remove.
- Product decisions live in `context/`. Retrieve them through `context/INDEX.md`
  rather than scanning, and keep a note and its catalog entry in the same change.
- `sources/` and `publication/` are passive history. They are not the current
  specification, and they are never shipped.
- Preserve unrelated working-tree changes.

## Validating

```sh
gofmt -w path/to/changed.go
go vet ./...
go test ./...
sh scripts/check-release.sh
```

`scripts/check-release.sh` is the full product check: it tests the Go packages,
builds the workspace archive and all six CLI targets, verifies checksums, and
exercises installation, initialization, and the publication safeguards. Run it in
fresh temporary directories; release assembly never overwrites existing output.

Cross-compilation proves a binary builds for another platform, not that it behaves
there. CI supplies native Linux, macOS, and Windows coverage.

## Opening a pull request

Use `type(scope): imperative subject` for commit subjects, with the changed
component as the scope. Describe what the change decides, not only what it edits —
a reviewer needs to see which behavior moved and why.

A change to shipped instruction or documentation needs the same care as code. The
test suite enforces which layer a rule belongs in and which audience a shipped
document is written for; if a documentation test fails, the placement is the
finding, not the test.

Publication, tagging, and deployment are separate maintainer acts and never part
of a pull request.
