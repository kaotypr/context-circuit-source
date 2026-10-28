# Contributing

This repository is published output. Every file in it except the changelog and
the repository's own landing-page files is assembled from
[kaotypr/context-circuit-source](https://github.com/kaotypr/context-circuit-source)
and replaced wholesale on each release, so a pull request opened here would be
overwritten by the next publication.

**Open issues and pull requests on the source repository instead.**

What belongs where:

- A problem with the workspace instructions, the skills, or the shipped
  documentation — the source repository, under `product/`.
- A problem with `context-circuit-cli` itself — the source repository, under
  `cmd/` and `internal/`.
- A question about using Context Circuit in your own project, or something the
  documentation does not answer — an issue on the source repository. An
  unanswerable question is a documentation defect.

Its [contributing
guide](https://github.com/kaotypr/context-circuit-source/blob/main/CONTRIBUTING.md)
describes the seam between the CLI and the coding agent, and how to validate a
change. A change that affects what agents read, do or produce also needs a
design case: what changes, why, a before and after, and who pays.

Releases here are tagged `v*`. The CLI is released on the source repository,
under `cli-v*` tags, together with its downloadable packages, and installing it
needs no checkout of either repository.
