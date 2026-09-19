# Security policy

## Reporting a vulnerability

Report privately through GitHub's **Report a vulnerability** button on the Security
tab of
[kaotypr/context-circuit-source](https://github.com/kaotypr/context-circuit-source/security),
where the code is maintained. That opens a private advisory visible only to you and
the maintainers. Please do not open a public issue for a suspected vulnerability.

Include what you did, what happened, and what you expected. A reproduction against
a scratch workspace — never one holding your own project's information — is more
useful than a description.

## What to expect from a workspace

These are design constraints, so a report that one of them is violated is a
security report:

- The CLI holds no model credentials and never calls a model API. Your coding
  agent's provider relationship is yours and is untouched by it.
- Installing the CLI needs no administrator access. It verifies a checksum and
  the CLI's reported version before switching the shared command, and leaves
  your existing command in place on any failure.
- Versions install side by side, so installing one for a workspace cannot change
  which version another workspace runs.
- If your organization mirrors releases into its own registry, the credential for
  it reaches only that registry.
- Files ending in `.local.yaml` hold one machine's repository paths and identity.
  They are never committed and never travel with the workspace.
- Commit, push, merge, delivery, and deletion of workspace data are gated on your
  explicit authorization and are never inferred by the agent or the CLI.

## Out of scope

- Code a coding agent writes. A workspace constrains where an agent works and what
  it must ask you about; it does not review what it produces.
- Vulnerabilities in the repositories you connect to a workspace.
- Prereleases (`-rc`, `-beta`, `-alpha`) are candidates. A fix is published as the
  next candidate rather than backported.
