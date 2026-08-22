# Human gates and cards

`wrapper/runtime/engine.sh` provides card formatting; this document owns the
human-facing effect language. A card is valid only for its current session,
action, target, observed state, and listed effects. A changed target or state
requires a new card.

Every consequential mutation names what will change and what will not. The
separate actions are identity acceptance, context acceptance, plan approval,
    execution trigger, repository bootstrap, material scope change, status change, takeover, delivery,
archive/restore, publication/deployment/merge, and destructive cleanup.
These are separate explicit gates, not one reusable confirmation.

Approval changes plan status and task projections only. Execution starts only
from a separate explicit request. A verifier can produce completion evidence but
cannot finish a plan. Cleanup first reports dirty and unpushed work and needs an
additional discard confirmation before anything destructive.

## Plan approval card

`Approve plan <id>` is a pre-confirmation request. It presents a current
session-bound card and mutates nothing: no plan, task, runtime, Git, lease, or
worktree state. Confirmation does not commit Git and does not start
`Run approved plan`.

```text
Action: present-approval-card
Target: <plan-id>
Observed state: current session; nothing has changed yet; plan remains draft and included tasks remain draft
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
persistence is allowed.

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
Will change: commit the listed plan.yaml and task frontmatter status changes
Will not change: implementation files, registered repositories, runtime state,
  external remotes, or delivery state
Risks/open decisions: this is a maintainer-source commit; inspect the file list
Confirmation requested: Confirm commit of the approved plan state.
```

The agent never commits this state automatically. Any other dirty path remains
`DIRTY_BASE_BLOCKED` and must be preserved for human resolution.
