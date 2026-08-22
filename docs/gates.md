# Human gates and cards

`wrapper/runtime/engine.sh` provides the generic `cc_confirmation_card`
template and the status-transition primitive. This document owns specialized
human-facing card wording. A card is valid only for its current session,
action, target, observed state, and listed effects. A changed target or state
requires a new card.

Every consequential mutation names what will change and what will not. The
separate actions are identity acceptance, context acceptance, plan approval,
    execution trigger, repository bootstrap, material scope change, status change, takeover, delivery,
archive/restore, publication/deployment/merge, and destructive cleanup.
These are separate explicit gates, not one reusable confirmation.

Effect identifiers on these cards are descriptive metadata owned by
`wrapper/contracts/routes.yaml`. They never grant authorization, select a
route, or skip confirmation. Each current confirmation lists immediate
effects and later authorized effects. Omitted fields appear as proposed
defaults on the current card; confirmation records those displayed values.
Fields with no default remain incomplete. No field is invented after
confirmation.

Approval changes plan status and task projections only. Execution starts only
from a separate explicit request. A verifier can produce completion evidence but
cannot finish a plan. Cleanup first reports dirty and unpushed work and needs an
additional discard confirmation before anything destructive.

## Identity-acceptance card

Identity acceptance records `workspace.yaml` identity metadata. It does not
generate Product Knowledge. The bounded identity region in `WORKSPACE.md`,
`PROJECT.md`, and `INDEX.md` is updated to agree with that identifier.
Authored Product Knowledge outside the region is unchanged.

```text
Action: identity-acceptance
Target: workspace identity
Observed state: current session; identity remains uninitialized until confirmation
Proposed defaults:
  mode: solo
  repositories: none
  roles: none
  default branches: main
Immediate effects: workspace.accept_identity
Later authorized effects: none from this confirmation
Will change after confirmation: accepted identity fields and identity regions in WORKSPACE.md, PROJECT.md, and INDEX.md
Will not change: Product Knowledge outside the identity region, Git, leases, clone, create-empty, execution, or delivery
Risks/open decisions: omitted fields are recorded as the displayed defaults; fields with no default remain incomplete
Confirmation requested: Confirm identity acceptance for this workspace in the current session.
```

## Repository-registration card

Repository registration is the shared-identity card consumed by a later
repository-bootstrap or reserved create-empty gate. It is not a router action
and does not clone or create a repository. Machine-specific paths remain only
in `repositories.local.yaml`.

```text
Action: repository-registration
Target: <logical-key>
Observed state: shared identity card for a later clone or create-empty gate; no destination is created now
Proposed defaults:
  logical key: <required; incomplete if omitted>
  canonical URL: none
  default branch: main
Immediate effects: workspace.register_repository
Later authorized effects: repository-bootstrap; repository-create-empty (reserved, not activated); execute-plan
Will change after confirmation: shared logical key, optional URL, optional branch, and identity regions
Will not change: repositories.local.yaml, clone, create-empty, execution, delivery, or authored Product Knowledge outside the identity region
Risks/open decisions: this card does not authorize clone or create-empty; those remain later gates
Confirmation requested: Confirm repository registration for this logical key in the current session.
```

## Plan approval card

`Approve plan <id>` is a pre-confirmation request. It presents a current
session-bound card and mutates nothing: no plan, task, runtime, Git, lease, or
worktree state. Confirmation does not commit Git and does not start
`Run approved plan`.

```text
Action: present-approval-card
Target: <plan-id>
Observed state: current session; nothing has changed yet; plan remains draft and included tasks remain draft
Immediate effects: plan.yaml status draft→approved; included task projections draft→ready
Later authorized effects: git.commit (isolated worktree only, after a separate execution request); delivery.push (delivery gate)
Will change after confirmation: plan.yaml status draft→approved; included task projections draft→ready
Will not change: Git, leases, worktrees, execution, delivery, publication, or runtime state
Risks/open decisions: confirmation does not commit Git and does not start Run approved plan
Confirmation requested: Confirm approval of plan <id>.
```

The exact confirmation `Confirm approval of plan <id>` is a separate
`approve-plan` action. It is status-only. It does not fold in the maintainer
commit or the execution trigger.

## Repository bootstrap card

Repository bootstrap is a separate gate from workspace initialization and plan
execution. Its card must name the logical repository, credential-free canonical
URL, selected remote, branch, exact destination, and observed existing-path
check. Confirmation is exact and session-bound. Without it, no destination,
parent directory, network request, or Git clone is started. Existing
destinations, unsafe paths, dirty sources, unavailable credentials, and offline
providers remain safe failures; no stash, reset, overwrite, or credential
persistence is allowed. `git.init` and `git.clone` are not `execute-plan`
effects.

```text
Action: repository-bootstrap
Repository: <logical-key>
Canonical URL: <credential-free URL or none>
Selected remote: <inspectable remote>
Branch: <branch>
Destination: <exact destination>
Existing-path check: <absent|exists>
Immediate effects: repository-bootstrap
Later authorized effects: execute-plan
Not execute-plan effects: git.init, git.clone
Will change: create the exact destination and clone the selected remote after confirmation.
Will not change: existing paths, bound repositories, credentials, or delivery state.
Risks/open decisions: host Git credentials and provider availability remain external.
Confirmation requested: Confirm repository-bootstrap for this exact target in the current session.
```

## Maintainer approval commit

When the `product-source` maintainer checkout is dirty only because approval
changed the selected plan from `draft` to `approved` and its tasks from `draft`
to `ready`, present this existing commit card in the same session immediately
after confirmed approval. Do not wait for `Run approved plan <id>` to discover
the stop. Instantiated or wrapped workspaces do not receive this card; their
next action remains `Run approved plan <id>`. The agent never commits this
state on the confirm-approval turn. Execution still reports
`MAINTAINER_APPROVAL_COMMIT_REQUIRED` if someone runs before that commit; that
is not an execution exemption. Any other dirty path remains `DIRTY_BASE_BLOCKED`.

```text
Action: commit-approved-plan
Target: <plan-id>
Observed state: exact approval projection is the only dirty source change
Immediate effects: git.commit (isolated worktree or maintainer source only)
Later authorized effects: delivery.push
Will change: commit the listed plan.yaml and task frontmatter status changes
Will not change: implementation files, registered repositories, runtime state,
  external remotes, or delivery state
Risks/open decisions: this is a maintainer-source commit; inspect the file list
Confirmation requested: Confirm commit of the approved plan state.
```

The agent never commits this state automatically. Any other dirty path remains
`DIRTY_BASE_BLOCKED` and must be preserved for human resolution.
