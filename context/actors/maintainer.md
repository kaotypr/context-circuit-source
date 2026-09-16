# Maintainer

Owns this checkout and decides what reaches a workspace.

## [Source layout and what ships](../maintenance/source-layout.md)

1. **As a maintainer, I want one explicit mapping from source file to shipped
   path,** so that what a workspace receives is readable in a single file rather
   than inferred from a build.
2. **As a maintainer, I want maintainer material to be structurally unable to
   ship,** so that design history and evidence cannot leak into a release by
   being in the wrong directory.

## [Validation and release checks](../maintenance/validation-and-release.md)

1. **As a maintainer, I want a new command to fail the build until its
   documentation exists,** so that the shipped reference cannot drift from the
   implemented surface.
2. **As a maintainer, I want the release check to build and exercise a real
   binary,** so that a candidate is proven by running rather than by compiling.
3. **As a maintainer, I want host behavior reported as unverified when it was
   not exercised,** so that a claim about a coding host is never stronger than
   the evidence behind it.

## [Two products and their version lines](../product/two-products-and-versioning.md)

1. **As a maintainer, I want the template and the executable to version
   independently,** so that a workspace can pin one without being dragged by the
   other's release cadence.
2. **As a maintainer, I want the embedded seed checked against the manifest,**
   so that the blank workspace inside the binary cannot quietly diverge from the
   one in the tree.

Owner:

- `context-circuit-source@scripts/release-manifest.txt` — the mapping this actor
  reads to know what ships.
- `context-circuit-source@scripts/` — the build, the checks, and the explicitly
  invoked publication.
- `context-circuit-source@internal/cli/documentation_test.go` — the coupling
  between a command and its reference.
