---
kind: domain
status: proposed
title: Repository binding and bootstrap
slug: repository-binding
owners: []
sources:
  - plans/context-circuit-plans/repository-bootstrap/plan.yaml
  - plans/context-circuit-plans/repository-bootstrap/PLAN.md
  - plans/context-circuit-plans/local-binding-troubleshoot/plan.yaml
  - plans/context-circuit-plans/local-binding-troubleshoot/PLAN.md
source_revisions:
  - plan: repository-bootstrap
    status: done
    updated_at: 2026-08-22T11:04:19Z
  - plan: local-binding-troubleshoot
    status: done
    updated_at: 2026-08-22T12:20:00Z
generated_at: 2026-08-23T00:00:00Z
review_date: 2026-09-22
freshness: proposed-from-done-plans
assumptions:
  - Canonical lifecycle status is plan.yaml, not PLAN.md.
unknowns:
  - Runtime completion sidecars were not present under .runtime/.
contradictions:
  - repository-bootstrap PLAN.md still says Status: draft while plan.yaml is done.
  - local-binding-troubleshoot PLAN.md still says Status: draft while plan.yaml is done.
acceptance:
  state: pending
  accepted_at:
  accepted_by:
workflows:
  - docs/getting-started.md
  - docs/gates.md
---

# Repository binding and bootstrap

## Summary

Shared repository identity stays in `workspace.yaml`. Host-local paths stay in
the ignored root file `repositories.local.yaml`. Cloning is an explicit
bootstrap gate. Route a bind, clone, missing-binding, or worktree-isolation
request here.

## Scope

Inside: logical repository keys, credential-free canonical URLs, default
branches, explicit path bindings, bootstrap confirmation, isolated execution
worktrees, and the missing-binding troubleshooting text.

Outside: automatic clone/fetch/push/merge, credential storage, filesystem
scanning for checkouts, and using a bound source repo as the writer worktree.

## Behavior

`workspace.yaml` may name repositories with an optional canonical URL and
default branch. It must not contain a machine-specific path (INV-REPO-01).

Each host may create `repositories.local.yaml` and bind a key already named in
`workspace.yaml` to an explicit absolute path, a workspace-relative path, or
`repositories/<key>` (INV-REPO-02). A missing binding file is expected on a
fresh clone: report `BINDING_MISSING`, do not scan, invent a path, or create
the file. The human either writes an explicit `path` or says
`Bootstrap repository <key>`.

`Bootstrap repository <key>` presents the logical key, canonical URL, selected
remote, branch, destination, and existing-path check. Only a current
`repository-bootstrap` confirmation may clone or create the destination
(INV-REPO-04). Existing paths, dirty sources, unsafe paths, missing
credentials, and offline providers fail without overwrite, stash, reset, or
credential persistence (INV-REPO-03, INV-REPO-05).

After binding, execution prepares
`.runtime/worktrees/<repository-key>/<plan-id>/` and leaves the bound source
untouched (INV-REPO-06). Wrapper upgrade and release preserve local bindings
and exclude them from artifacts (INV-REPO-07).

## Workflows

- Bind a repository: `docs/getting-started.md`
- Missing `repositories.local.yaml`: same getting-started troubleshooting paragraph
- Bootstrap confirmation: `docs/gates.md` repository bootstrap card

## Interfaces

- Shared identity: `workspace.yaml` `repositories.<key>`
- Host binding: `repositories.local.yaml`
- Optional convenience directory: `repositories/<key>` (never assumed to exist)
- Human request: `Bootstrap repository <key>`
- Runtime worktrees: `.runtime/worktrees/<repository-key>/<plan-id>/`

## Constraints and edge cases

Credentials stay in host Git configuration or the SSH agent. Relative paths
resolve from the workspace root and reject traversal. Dirty bound sources
block execution rather than being repaired implicitly.

## Implementation references

- `wrapper/runtime/engine.sh` repository bootstrap and binding primitives
- `wrapper/contracts/invariants.yaml` INV-REPO-01 through INV-REPO-07
- `wrapper/contracts/schemas/workspace.yaml`

## Verification

Done-plan verification IDs RB-VT-01–RB-VT-08 and LBT-VT-01–LBT-VT-03.

## Provenance

Read only the selected done plans `repository-bootstrap` (CC-001) and
`local-binding-troubleshoot` (CC-003), plus the getting-started troubleshooting
paragraph those plans landed. Raw `sources/` was not scanned. HEAD at
generation was `b7a11f3`.

## Acceptance notes

This page is proposed. It does not replace accepted decisions in
`context/DECISIONS.md`. Human context acceptance is still required.
