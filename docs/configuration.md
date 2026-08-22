# Configuration

Configuration is optional provider-neutral intent. Initialization does not ask
delivery or integration questions. A human may later confirm one policy:

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
reconfirmation. Provider failure falls back to the filesystem workflow with
`disabled`, `denied`, or `unavailable` evidence.
