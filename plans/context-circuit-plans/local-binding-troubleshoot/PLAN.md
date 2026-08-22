# Document missing repositories.local.yaml troubleshooting

Status: draft
Repository: context-circuit-source
Source: repository-evidence

## Review summary

`docs/getting-started.md` already shows how to bind a repository with
`repositories.local.yaml`, and it notes that a missing source is reported
and left untouched. It does not tell a human the next safe action when that
ignored host-local file is absent. This plan adds one short troubleshooting
paragraph in the existing Bind a repository section.

## What approval authorizes

A bounded docs change: insert a troubleshooting paragraph into
`docs/getting-started.md`, and optionally one semantic `contains` assertion
that the guidance remains present. Approval does not start execution, claim a
lease, create a worktree, change Git, deliver, publish, deploy, merge, or
clean runtime.

## Scope and non-goals

Paths and behavior boundaries are linked to `plan.yaml`. Do not change
binding resolution, bootstrap confirmation, existing maintainer plans
`repository-bootstrap` or `multi-host-agent-support`, or any file other than
the named docs path plus the optional contract assertion.

## Proposed solution

After the current Bind a repository examples, add a short paragraph that
states all of the following without expanding into a new page:

- A missing `repositories.local.yaml` is expected on a fresh or newly cloned
  workspace because the file is host-local, gitignored, and never shipped.
- Context Circuit reports the missing binding and does not scan for
  repositories, invent a path, or create the file on its own.
- Create the ignored root file with an explicit `path` (absolute,
  workspace-relative, or `repositories/<key>`) for a logical repository
  already named in `workspace.yaml`.
- If the checkout does not exist yet, say `Bootstrap repository <key>`
  rather than guessing a destination.
- Keep credentials in host Git configuration; never put secrets or
  machine-specific paths into shared workspace files.

Intended paragraph, to be edited only for fit with surrounding prose:

> If `repositories.local.yaml` is missing, that is expected on a fresh or
> newly cloned workspace: the file is host-local, gitignored, and never
> shipped. Context Circuit reports the missing binding and does not scan
> the filesystem, invent a path, or create the file on its own. Create the
> ignored root file with an explicit `path` for the logical repository
> already named in `workspace.yaml`. If the checkout does not exist yet,
> say `Bootstrap repository <key>` instead of guessing a destination. Keep
> credentials in host Git configuration; never put secrets or
> machine-specific paths into shared workspace files.

Canonical behavior remains INV-REPO-01 through INV-REPO-05 and
`cc_resolve_repository_binding` reporting `BINDING_MISSING` when the file
is absent.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| LBT-001 | Troubleshooting paragraph landed in getting started, with optional semantic lock | none |

## Acceptance criteria

- LBT-AC-01: `docs/getting-started.md` contains a short troubleshooting paragraph that tells a human what to do when `repositories.local.yaml` is missing.
- LBT-AC-02: The paragraph states the file is host-local/gitignored, is not created implicitly, and a missing binding fails without scanning or touching other paths.
- LBT-AC-03: The paragraph points the human to an explicit path binding or `Bootstrap repository <key>`, and does not request or record credentials.

## Verification

- LBT-VT-01 proves the binding file is still named in getting started.
- LBT-VT-02 proves the troubleshooting language is present.
- LBT-VT-03 proves task frontmatter and credential-free contracts still pass.

Executable commands remain canonical in `plan.yaml`.

## Risks, assumptions, and open decisions

- The paragraph must stay short; do not split it into a new troubleshooting page.
- Do not imply that the agent may write `repositories.local.yaml` during ordinary orientation or plan execution.
- No human decision is required beyond whether this docs-only draft is the right bounded change.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Exact evidence used to draft this bundle:

- workspace identity `product-source` in `workspace.yaml`
- Tier 0 summaries `context/INDEX.md`, `context/WORKSPACE.md`, `context/PROJECT.md`
- `context/CONVENTIONS.md`
- `docs/getting-started.md`, `docs/configuration.md`, `docs/runtime-contract.md`, `docs/planning.md`, `docs/templates/plan.md`, `docs/templates/task.md`
- `wrapper/contracts/invariants.yaml` INV-REPO-01..05, INV-AUTH-02, INV-AUTH-04
- `wrapper/contracts/schemas/workspace.yaml` and `wrapper/contracts/schemas/plan.yaml`
- `wrapper/runtime/engine.sh` `cc_resolve_repository_binding` / `BINDING_MISSING`
- plan field shape from sibling maintainer `plan.yaml` files, without modifying them
