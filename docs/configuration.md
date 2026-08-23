# Configuration

Configuration is optional provider-neutral intent. Initialization does not ask
delivery or integration questions. A human may later confirm one policy:

`workspace.yaml` is the workspace identifier. `context/` is Product Knowledge.
A bounded identity region in `WORKSPACE.md`, `PROJECT.md`, and `INDEX.md` must
agree with that identifier; authored Product Knowledge outside the region is
not a projection input or output. Omitted confirmation fields are shown as
proposed defaults on the current card; confirmation records those displayed
values. Effect identifiers are descriptive metadata and never authorize a
route or gate.

Repository identity and host binding are separate configuration layers. Shared
`workspace.yaml` may contain only a logical repository key, an optional
credential-free `canonical_url`, and an optional `default_branch`. A developer
may create the ignored root `repositories.local.yaml` with an explicit `path`
and inspectable `remote`. The path may be absolute, workspace-relative, or
`repositories/<key>`; it is never copied into shared identity or a release.

```yaml
configuration:
  delivery:
    policy: manual | remote-review | local-target
    repositories: [app]
    target_branches: {app: main}
    authorization: pending | confirmed | stale
    fallback: manual
  integrations: []
  hosts: []
```

Configuration never stores credentials or provider payloads and never grants
commit, push, merge, publication, deployment, or external activity authority.
Changed repository, branch, risk, or stale authorization requires focused
reconfirmation. Host entries are optional provider-neutral observations, for
example `host_id`, `observed_version`, `instruction_surface`, capability
labels, and `provider_status`; do not copy host-local settings, auth state, MCP
configuration, or transcripts into workspace files. Provider failure falls
back to the filesystem workflow with `disabled`, `denied`, or `unavailable`
evidence. A missing required child is `host-blocked`, never an authorization
grant.
