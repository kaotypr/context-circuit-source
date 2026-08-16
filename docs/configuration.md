# Configuration

`workspace.yaml` is the canonical, portable configuration location. Its core
identity fields remain the workspace name, solo/team mode, default branch,
registered repositories, and optional source-register location. An empty
`repositories: {}` registry is valid before initialization. Repository entries
have an exact key, path, mode, role, agent, and `default_branch`; initialization
recommends `development` only when it exists and lets the user choose another
branch.

Initialization records core identity only. Delivery behavior, merge policy,
publication, deployment, host conveniences, and external activity
integrations are optional later configuration and do not belong in first-run
questions. Users request those capabilities through `cc-configure-workspace`
or an equivalent explicit natural-language request.

## Portable configuration record

The optional `configuration` block is durable intent, not a secret store or an
authorization bypass:

```yaml
configuration:
  delivery:
    scope: workspace
    policy: manual
    repositories:
      - context-circuit
    target_branches:
      context-circuit: development
    authorization: pending
    fallback: manual
  integrations:
    - id: activity-record
      enabled: false
      provider: none
      reads:
        - plan summaries
      writes: []
      authorization: not-requested
      fallback: core-filesystem
  hosts:
    codex:
      capability: cc-configure-workspace
      mode: native-or-core
```

The fields mean:

- `delivery.scope` is `workspace` or a specific repository scope.
- `delivery.policy` is `remote-review`, `local-target`, or `manual`.
  When reading an existing file, `team-review` is an alias for
  `remote-review` and `solo-local` is an alias for `local-target`. Those
  old IDs remain valid reads and are not treated as missing configuration.
  New writes use the new IDs.
- `delivery.repositories` and `target_branches` identify affected targets;
  targets must match registered repositories and their configured default
  active branches unless the user explicitly confirms a later branch.
- `delivery.authorization` records only a state such as `pending`,
  `approved`, `denied`, or `expired`; it never contains a token, credential,
  or person-specific secret.
- `delivery.fallback` is `manual` and is required whenever the selected
  capability is absent, denied, stale, or unavailable.
- Each integration is disabled unless `enabled: true` follows explicit
  opt-in. `reads`, `writes`, and provider-neutral authorization state describe
  the boundary without copying provider activity into workspace files. An
  explicitly configured provider may expose an opaque `external_status`
  annotation, but it never replaces canonical task status.
- `hosts` maps a host to the same capability name and fallback mode. It does
  not redefine intent, ownership, verification, or human gates.

Write the block only after the user confirms the requested capability, scope,
effects, authorization boundary, and fallback. A changed repository, target
branch, risk, or authorization state triggers focused confirmation rather than
silently reusing the old configuration.

Keep repository paths credential-free and specific. Use `ignored-clone` for a
separate local clone or `submodule` for a tracked submodule. Context Circuit
refuses dirty bases and never rewrites unrelated work. Credentials, tokens,
secrets, provider payloads, and external activity records remain outside
ordinary workspace state; the wrapper stores none of them.

Configuration does not change canonical plan/task status, runtime ownership,
or verification results. It only tells the coordinator what delivery or
optional adapter behavior may be proposed after the normal evidence and human
gates are satisfied. Current authorization is still required before commit,
push, or merge. Confirming configuration does not grant those actions.
