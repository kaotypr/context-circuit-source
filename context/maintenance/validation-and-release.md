# Validation and release checks

How a change here is proven before it ships, and what stays unproven.

## What a change here runs

Format changed Go files, then the test suite, then vetting, then the release
check. Tests use temporary directories and disposable Git repositories, so they
establish real file and Git behavior rather than mocked behavior.

The release check is the broad one: it builds and exercises a native binary,
verifies the embedded seed inventory against the manifest, cross-compiles every
supported target, and tests the publication guards and commit and tag behavior
in disposable fixtures. Passing it a new output directory retains both sets of
checked assets.

Continuous integration runs the native suite — including installer upgrade and
rollback fixtures — on Linux, macOS, and Windows, plus the release check on
Linux. The matrix exists because **cross-compilation alone is not proof of
behavior on another operating system**, and the filesystem cloning this product
depends on is exactly where that bites.

Owner:

- `context-circuit-source@scripts/check-release.sh` — the broad check a change
  here runs before it is published.
- `context-circuit-source@.github/workflows/` — the matrix that proves behavior
  per operating system, and the publication runs.

## Build outputs

With no output argument each build clean-rebuilds its own directory under the
distribution folder, named for the product and version, so repeated runs replace
rather than fail and neither build removes the other's assets. An explicitly
passed output directory must be new: builds never replace existing output there.

## Publication

The two products publish separately. The workspace template publishes from its
version file to the published template repository's release tags; the executable
publishes from its own version file to this repository's prefixed tags. Changing
a version file does not publish — publication is explicitly invoked and requires
authorization.

Release assets carry dependency licenses. Dependencies are deliberately few: a
document-preserving YAML library and a portable file-locking library. Go is a
contributor requirement only; a user needs the executable for their platform and
installed Git, with no other language runtime.

## What validation does not establish

Markdown scenarios are validated against the shared entry instruction by
reading, not by execution. Deterministic tests establish file and Git behavior;
they do not establish whether a coding host loads or follows the shared
instruction. **Host behavior is reported as unverified unless it was actually
exercised.**

This is the honest limit of the whole test story, and it is stated in the
shipped release notes too, because a green pipeline would otherwise imply it.
The instruction surface is the product, and it is unverified in the way that
matters most.
