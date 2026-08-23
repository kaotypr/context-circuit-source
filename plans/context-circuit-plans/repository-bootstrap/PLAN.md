# Portable repository metadata and explicit repository bootstrap

Status: draft  
Repository: context-circuit-source  
Source: sources/context-circuit-design/DESIGN-SPEC.md

## Review summary

This plan dogfoods the redesigned Context Circuit workflow on the Context
Circuit source repository. It adds a portable shared repository identity,
root-level per-user bindings, an optional gitignored `repositories/` location,
and an explicit clone/bootstrap journey for new team members.

The plan keeps repository location separate from repository identity. Team
members may use different filesystem paths, while `workspace.yaml` remains
portable and reviewable.

## What approval authorizes

Approval authorizes bounded implementation of the tasks and acceptance
criteria in `plan.yaml`. Approval does not start execution, claim a lease,
create a worktree, change Git, clone a repository, deliver, publish, deploy,
merge, or clean runtime.

## Scope and non-goals

The implementation owns the repository contract, host-local binding file,
optional convenience directory, explicit bootstrap route, safety checks,
runtime integration, migration behavior, release exclusions, documentation,
and semantic tests.

It does not add automatic external activity, credential storage, provider
payload persistence, or a second routing/lifecycle authority.

## Proposed solution

`workspace.yaml` will contain only shared repository identity:

```yaml
repositories:
  app:
    canonical_url: https://github.com/acme/app.git
    default_branch: main
```

Each developer may create the ignored root file `repositories.local.yaml`:

```yaml
repositories:
  app:
    path: /home/alice/projects/app
    remote: git@github.com:acme/app.git
```

If the developer chooses the convenience layout, the path may instead be
`repositories/app`. Context Circuit never assumes that directory exists and
never includes its contents in a release artifact.

A clone request produces a confirmation card listing the logical repository,
canonical URL, selected transport, branch, destination, and safety checks.
Only a current explicit confirmation may create the destination or invoke
Git. Existing directories, dirty sources, unsafe paths, unavailable
credentials, and offline providers remain blocked or report an offline
fallback. Credentials are supplied by the host Git configuration or SSH
agent and are never recorded.

After binding, execution prepares the isolated worktree under
`.runtime/worktrees/<repository-key>/<plan-id>/`. The bound source checkout is
never used as the writer worktree and is never silently modified.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| RB-001 | Shared metadata, local binding, ownership, and release contracts | — |
| RB-002 | Safe resolution for external and convenience paths | RB-001 |
| RB-003 | Explicit clone/bootstrap route and evidence | RB-001, RB-002 |
| RB-004 | Runtime, migration, release, and documentation integration | RB-002, RB-003 |
| RB-005 | Independent semantic and self-hosted verification | RB-001–RB-004 |

## Acceptance criteria

- RB-AC-01: Shared repository identity is portable and credential-free.
- RB-AC-02: Root local bindings support external, relative, and convenience paths.
- RB-AC-03: Clone/bootstrap requires an exact human confirmation.
- RB-AC-04: Unsafe, dirty, existing, offline, and credential failures are safe.
- RB-AC-05: Writers use isolated runtime worktrees and preserve source checkouts.
- RB-AC-06: Upgrade, rollback, recovery, and release boundaries preserve and
  exclude the correct data.

## Verification

The verification IDs and canonical commands are owned by `plan.yaml`. The
independent verifier must inspect the changed contracts, exercise success and
failure fixtures, rebuild the release artifact, and confirm that the source
plan, local bindings, and repository contents are absent from it.

## Risks, assumptions, and open decisions

- Git transport selection must not become an authorization grant.
- A canonical URL is identity metadata; the local remote may differ by host
  transport but must remain inspectable.
- Cloning private repositories depends on the user’s existing Git credential
  helper or SSH agent; Context Circuit must not request or persist secrets.
- A repository source that is dirty remains preserved and blocked for isolated
  execution until the user resolves it independently.
- The root `plans/` directory is maintainer-only and must remain outside the
  release staging allowlist.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

The plan is grounded in the active design specification, the wrapper contract
and invariant owners, the current runtime worktree primitive, release
assembler, migration boundary, and semantic acceptance suites.
