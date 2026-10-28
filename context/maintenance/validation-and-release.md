# Validation and release checks

How a change here is proven before it ships, and what the tests do not establish.

## What a change here runs

The contributor guide, `CONTRIBUTING.md`, lists the commands a change here runs.
Tests use temporary directories and disposable Git repositories, so they
establish real file and Git behavior rather than mocked behavior.

The suite also validates this project's own product knowledge with the product's
own diagnostic. One test builds a throwaway workspace around the real knowledge
tree and binds it to this checkout, so each note's code anchors resolve against
real history.

Findings split by the taxonomy the diagnostic already uses:

- one a command or an edit discharges fails the test;
- one marked as needing a person — a note to re-read against code that moved —
  is reported instead, rather than blocking a contributor on another's judgment.

Continuous integration clones full history for it, because that pass counts
commits under those anchors and a shallow clone would let it pass without
looking.

The release check is the broad one: it builds and exercises a native binary,
verifies the embedded seed inventory against the manifest, cross-compiles every
supported target, and tests the publication guards and commit and tag behavior
in disposable fixtures. Passing it a new output directory retains both sets of
checked assets.

Continuous integration runs the native suite — including installer upgrade and
rollback fixtures — on Linux, macOS, and Windows, plus the release check on
Linux for a pull request into a release branch. Every pull request also runs a
check that it carries a change file for each product it changes, or says it
changes nothing users see. The matrix exists because **cross-compilation alone is not proof of
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

The two products publish separately, each from the repository that owns it. The
workspace template publishes from its version file to the published template
repository's `v*` release tags. The CLI is tagged `cli-v*` here, where its
source is, and its release is created here too, on the commit that built it.
Both publish from a release line's branch, `release/X.Y`, when a release-prep
merge there sets a version file; the resolver refuses a version on any other
branch, `main` included. That merge, or a run started by hand, is the explicit,
authorized act — changing a version file on `main` publishes nothing.

Publishing the template replaces the destination's whole tree, so the files that
repository owns rather than a workspace — its license, conduct, contributing and
security pages — are restored over the extracted artifact and held out of the
comparison that decides whether there is anything to publish. The fixture proves
all of it: that they land, that the license shown there is byte-identical to the
one a workspace receives, that an exported workspace carries none of them at its
root, and that a `README.md` among them is refused rather than replacing the
assembled product guide.

Release assets carry dependency licenses. Dependencies are deliberately few: a
document-preserving YAML library and a portable file-locking library. Go is a
contributor requirement only; a user needs the CLI for their platform and
installed Git, with no other language runtime.

## What validation does not establish

Deterministic tests establish file and Git behavior. They do not establish
whether a coding host loads or follows the shared instruction, and a green
pipeline does not either.

That is measured separately, by the agent tests (see
[Agent tests](agent-tests.md)). A result there is a statement about one host,
one model, and one date: host H, model M, k of N runs passed, on date D. It is
produced outside CI, by whoever changes what agents see. No 2.5.0 baseline is
recorded yet, so no figures are cited here. Host behavior that was not exercised
is reported as not exercised, and a claim about a host is never stronger than
the runs behind it.
