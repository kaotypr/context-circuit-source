# Security policy

## Reporting a vulnerability

Report privately through GitHub's **Report a vulnerability** button on this
repository's Security tab. That opens a private advisory visible only to you and
the maintainers. Please do not open a public issue for a suspected vulnerability.

Include what you did, what happened, and what you expected — a reproduction
against a scratch workspace is more useful than a description. You will get an
acknowledgement, and a fix or an explanation of why the behavior is intended.

## Supported versions

Fixes land in the current release line. Prereleases (`-rc`, `-beta`, `-alpha`)
are candidates, and a fix for one is published as the next candidate rather than
backported.

## What the products do and do not touch

These properties are design constraints, so a report that one of them is violated
is a security report:

- The CLI holds no model credentials and never calls a model API.
- Installation needs no administrator access, verifies a checksum and the CLI's
  reported version before switching the shared command, and leaves the existing
  command in place on any failure.
- A registry credential reaches only the registry in use.
- Files ending in `.local.yaml` hold one machine's paths and identity. They are
  never committed and never travel with the workspace.
- Commit, push, merge, publication, deployment, and deletion of workspace data
  are gated on explicit human authorization and are never inferred.

## Out of scope

- Anything a coding agent decides to write. The workspace constrains where an
  agent works and what it must ask about; it does not review generated code.
- Vulnerabilities in a project you connect to a workspace.
- Missing hardening in `sources/`, which is passive design history and is
  never shipped.
