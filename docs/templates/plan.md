# {{title}}

Status: {{status}}
Repository: {{repository}}
Source: {{source}}

## Review summary

{{outcome_and_users}}

## What approval authorizes

{{bounded_implementation_intent}}

Approval does not start execution, claim a lease, create a worktree, change
Git, deliver, publish, deploy, merge, or clean runtime.

## Scope and non-goals

Paths and behavior boundaries are linked to `plan.yaml`.

## Proposed solution

{{reviewable_behavior_and_architecture}}

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| {{task_id}} | {{task_outcome}} | {{dependencies}} |

## Acceptance criteria

- {{acceptance_id}}: {{observable_outcome}}

## Verification

List verification IDs, what they prove, and independent-review expectations.
Executable commands remain canonical in `plan.yaml`.

## Risks, assumptions, and open decisions

Only items that can affect approval or execution.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Exact accepted context, source, and repository evidence used.
