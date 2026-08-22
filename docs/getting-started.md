# Getting started

In a released workspace, say “Start or resume work in this workspace.” The
agent reads the compact entry spine and reports identity, repositories,
runtime state, the selected probe, and one safe next action.

The same request works from Codex CLI, Claude Code, and Cursor Agent CLI.
Codex and Cursor read the shared root `AGENTS.md`; Claude Code reads the
shipped `CLAUDE.md`, which imports that same contract. Host versions and
capabilities are evidence only. They never replace the router, a human gate,
the exclusive writer, or the independent verifier.

## Host entry and resume

The same request is used in interactive mode for Codex CLI, Claude Code, and
Cursor Agent CLI. Print or non-interactive mode is a read-only probe unless the
same Context Circuit human gate is already present. Resume re-reads the
session receipt, latest handoff, wrapper version, Git state, lease, and
worktree before re-entering the route.

Codex native subagents, Claude Task/subagents, and Cursor Task/subagents map to
the bounded writer or independent read-only verifier packet. If the required
child is unavailable, the result is host-blocked and the root does not
self-verify. If a provider is disabled, denied, or unavailable, continue with
the filesystem-only workflow.

## New workspace

Until `workspace.yaml` records `identity.status: accepted`, requests that
would create or change anything are routed back to initialization. This
includes broad requests such as “help me build this”; the agent must not infer
identity, create a plan, or write implementation files. Read-only orientation
remains available.

## Bind a repository

Shared `workspace.yaml` metadata identifies a logical repository without a
machine-specific path:

```yaml
repositories:
  app:
    canonical_url: https://github.com/acme/app.git
    default_branch: main
```

Each host may create the ignored root file `repositories.local.yaml`:

```yaml
repositories:
  app:
    path: /home/alice/projects/app
    remote: git@github.com:acme/app.git
```

The binding path may also be workspace-relative, such as `projects/app`, or
use the optional convenience location `repositories/app`. The path is
explicit; Context Circuit never scans for repositories or stores credentials.
A missing, unsafe, identity-mismatched, or dirty source is reported and
remains untouched.

If `repositories.local.yaml` is missing, that is expected on a fresh or newly
cloned workspace: the file is host-local, gitignored, and never shipped.
Context Circuit reports the missing binding and does not scan the filesystem,
invent a path, or create the file on its own. Create the ignored root file
with an explicit `path` for the logical repository already named in
`workspace.yaml`. If the checkout does not exist yet, say
`Bootstrap repository <key>` instead of guessing a destination. Keep
credentials in host Git configuration; never put secrets or machine-specific
paths into shared workspace files.

## Bootstrap a repository

Say “Bootstrap repository app” to request a clone. The agent presents a
repository-bootstrap card containing the logical repository, canonical URL,
selected remote, branch, exact destination, and existing-path check. Only a
current confirmation of that exact card may create the destination or invoke
Git. Host Git configuration or an SSH agent supplies authentication; secrets
are never requested, recorded, or copied into workspace state. Provider
failure produces an offline fallback while filesystem evidence remains
resumable.

Initialization records only mode, repositories or project items, roles, and
default branches. Zero repositories is valid. Identity acceptance is a human
gate. Later, selected evidence may produce an Idea Brief, PRD, accepted Product
Knowledge, or direct plan; no artifact is forced when the request is already
clear.

## Plan journey

Create a plan bundle with `plan.yaml`, human-facing `PLAN.md`, and task files.
Review is read-only. Approval changes `draft` → `approved` and tasks to
`ready`; it does not execute. A separate named run request creates the root
session, exclusive lease/worktree, writer delegation, and independent verifier.
After evidence is ready, a separate finish confirmation changes `approved` →
`done`. Runtime and worktrees remain until a later, separately gated cleanup.

## Interrupted work

Resume validates the session, receipt, handoff, wrapper version, primary
evidence, Git state, and ownership. A foreign live owner or ambiguous lease is
read-only and needs a human takeover decision. Dirty work is preserved.

## Human surfaces

Humans normally need only this README, a plan's `PLAN.md`, the conversational
next-action card, and the latest surfaced handoff. Machine records remain
inspectable but are not a manual editing workflow.
