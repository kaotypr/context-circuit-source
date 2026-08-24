---
kind: domain
status: accepted
title: Workspace orientation and repository binding
slug: repository-binding
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - Portable identity lives in workspace.yaml; host-local paths in repositories.local.yaml.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - docs/getting-started.md
---

# Workspace orientation and repository binding

## Summary

Orienting a workspace and connecting its repositories. Shared repository
identity stays in `workspace.yaml`; host-local paths stay in the ignored root
file `repositories.local.yaml`. Route orient, register, connect, clone, `git
init`, missing-binding, and repository-resolution requests here. Owned by the
`cc-workspace` skill.

## Scope

Inside: workspace orientation (read-only), portable logical repository ids,
credential-free canonical URLs, `default_branch` as portable clone guidance,
host-local bindings (`path` + user-selected `anchor_branch`), the reserved
`workspace` id at path `.`, and fail-closed binding resolution.

Outside: the execution worktree lifecycle (see [plan-execution](../plan-execution/README.md)),
credential storage, filesystem scanning for checkouts, and using a bound source
checkout as the writer worktree.

## Behavior

`cc-workspace` initializes or orients a workspace and never creates plans or
executes work. Orientation reports whether each repository's local binding is
available.

`workspace.yaml` holds only portable identity: a logical repository key, an
optional credential-free canonical URL, and an optional `default_branch`. It
never contains a machine-specific path or credentials (INV-REPO-01).

Each host may create `repositories.local.yaml` binding a key already named in
`workspace.yaml` to a concrete `path` and a user-selected `anchor_branch`.
`anchor_branch` is the required execution base and default pull-request target;
`default_branch` is only portable clone/setup guidance and is never inferred as
the anchor (INV-REPO-02). A missing binding file is expected on a fresh clone:
the engine reports `BINDING_MISSING`; the agent does not scan, invent a path, or
create the file — the human supplies an explicit `path` or connects/clones the
repository.

`repositories/<key>` is the gitignored default destination for cloned or
initialized repositories. When the workspace root is itself a Git repository it
binds as the reserved logical id `workspace` at path `.` and is never placed
under `repositories/` (INV-REPO-03).

Binding resolution is explicit and bounded. Missing, ambiguous, non-Git,
traversal, unsafe-symlink, or identity-mismatched bindings fail closed without
scanning the filesystem or substituting a similarly named path (INV-REPO-04).
Workspace-relative paths reject traversal and unsafe symlinks; credentials
remain in host Git configuration or the host agent (INV-SEC-01).

## Workflows

- Orient, connect, clone, or initialize a repository: `docs/getting-started.md`

## Interfaces

- Shared identity: `workspace.yaml` `repositories.<key>`
- Host binding: `repositories.local.yaml` (`path`, `anchor_branch`)
- Convenience directory: `repositories/<key>` (gitignored; never assumed present)
- Reserved id: `workspace` at path `.`
- Missing-binding signal: `BINDING_MISSING`

## Constraints and edge cases

Credentials stay in host Git configuration or the SSH agent. Relative paths
resolve from the workspace root and reject traversal. A dirty or mismatched
bound source fails closed rather than being repaired implicitly.

## Registration and connected-repository states

Portable identity in `workspace.yaml` carries `schema_version`, `workspace`,
`title`, and `purpose`; the reserved `workspace` entry is omitted when the root
is unversioned. Portable identity must never contain local machine paths, access
tokens, private keys, provider payloads, or credentials.

Registration is two parts: record the portable logical identity and add a
host-local binding. Registration does not clone; a clone or `git init` first
reports the source URL, destination, branch, and external Git effect. A `git
init` creates the directory, initializes with the anchor as the initial branch,
records identity and anchor, and makes an initial (optionally empty) anchor
commit; a repository with no commit is registered but not ready, because it
cannot provide a worktree base.

A connected repository moves through states — registered, bound, unavailable,
mismatched, dirty, ready, active — reported in project terms rather than raw
record names. `anchor_branch` is the user's actual active branch (for example
`development` or `kao/development/v0.5`), not the same as `default_branch`;
different users may set different anchors for the same logical repository.

## Implementation references

- `wrapper/runtime/engine.sh`: `cc_repository_register`, `cc_binding_field`,
  `cc_repo_resolve`, `cc_repo_anchor_commit`, `cc_repo_clean`,
  `cc_repository_preflight`, `cc_worktree_prepare`, `cc_workspace_validate`,
  `cc_workspace_init`
- `wrapper/contracts/schemas/workspace.yaml`,
  `wrapper/contracts/schemas/repositories-local.yaml`
- `wrapper/contracts/invariants.yaml`: INV-REPO-01, INV-REPO-02, INV-REPO-03,
  INV-REPO-04, INV-SEC-01
- `.agents/skills/cc-workspace/SKILL.md`

## Provenance

Re-grounded on the current wrapper at HEAD `4b8ac0b`. The previous-version
maintainer plans that first seeded this page were deleted in `4b8ac0b`; their
provenance was retired (see the accepted retirement of `context/sources.yaml`).
Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-08-24. Consolidated from the former `repository-binding` draft,
re-grounded on the shipped wrapper, and broadened to include `cc-workspace`
orientation. The worktree isolation lifecycle moved to
[plan-execution](../plan-execution/README.md).
