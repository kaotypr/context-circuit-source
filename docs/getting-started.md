# Getting started

Context Circuit is a universal project workspace. You talk to it in ordinary
language; it holds Product Knowledge and readable plans and executes approved
plans safely across one or more Git repositories.

## 1. Initialize

> Initialize this workspace for <project name and purpose>.

This records project identity in `workspace.yaml` and creates the minimal
structure. It does not invent Product Knowledge or create plans.

## 2. Connect repositories

> Connect the api repository at ../commerce-api on branch development.

Registration records portable logical identity in `workspace.yaml` and a
host-local binding (path + `anchor_branch`) in `repositories.local.yaml`. New
project repositories can be cloned or initialized under the git-ignored
`repositories/<id>/`. The workspace root itself may be bound as the reserved
`workspace` repository. `anchor_branch` — not `default_branch` — is the
execution base and default pull-request target.

## 3. Gather context (optional)

> Gather context about billing from sources/billing-requirements.md.

Context proposals are staged and require explicit human acceptance.

## 4. Create and review a plan

> Create a plan for adding recurring billing.
> Review plan 0001-billing-v2.

Creating or reviewing a plan never approves or executes it. Review is a
discussion that can update draft content.

For a small change you want to judge live, you can work directly instead:

> Work with me on tightening the checkout screen in the web repository.

Direct collaboration uses one worker in a separate working copy of one connected
repository. It creates no plan or independent check; the result is
human-supervised, never verified, and delivery remains separate.

## 5. Approve and execute

> Approve plan 0001-billing-v2 and execute it.

One worker implements every task in dependency order across the mapped
repositories and commits each one; an independent read-only verifier checks the
latest commits and the worker repairs failures with new commits (stopping after
three failures).

## 6. Complete

> Mark 0001-billing-v2 complete.

Completion is human-controlled and allowed only after a verifier pass. It records
the implementation and reconciles it against Product Knowledge, surfacing
proposals you accept or defer separately.

## 7. Deliver (separate)

> Open a pull request for 0001-billing-v2.

Pull request, merge, push, publication, deployment, archive, and cleanup are
separate explicit actions. A pull request uses each execution branch as source
and the recorded `anchor_branch` as the default target.
