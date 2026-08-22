# Fast and explicit plan approval confirmation flow

Status: approved  
Repository: context-circuit-source  
Source: repository-evidence

## Review summary

The observed approval session behaved safely but felt slow and ambiguous. The
initial `Approve plan <id>` request spent time checking state before returning a
card, and the confirmed turn took nearly three minutes while manually changing
one plan status and six task projections, running acceptance checks, and
repairing unintended Markdown end-of-file drift.

On this `product-source` checkout the worktree is cut from the current
workspace, not from a bound product repository. Confirming approval therefore
leaves a dirty base, and `Run approved plan <id>` stops on
`MAINTAINER_APPROVAL_COMMIT_REQUIRED` unless the exact status projection is
committed first. Wrapped product workspaces do not have that stop: their
worktree comes from a different clean repo.

This plan keeps the two-step approval gate, the separate maintainer commit
gate, and the separate execution trigger. It makes the approval card explicit,
routes confirmed approval through one canonical status transition, presents the
existing commit card immediately on `product-source` instead of discovering it
at run time, and adds regression coverage for those failure modes.

## What approval authorizes

Approval authorizes bounded implementation of the tasks and acceptance
criteria in `plan.yaml`. It does not authorize execution, lease acquisition,
worktree creation, Git changes, delivery, publication, deployment, merge, or
cleanup. On `product-source`, Git commit of the exact approval projection
remains a separate current-session confirmation.

## Scope and non-goals

The implementation is limited to the approval route/card, the confirmed
status-transition path, the `product-source` follow-on presentation of the
existing commit card, coordinator procedure, human-facing wording, and tests.
The canonical router, lifecycle, ownership classifier, and human-gate owners
stay unchanged. Existing dirty or untracked files, including unrelated plans
and source fixtures, remain untouched.

## Proposed solution

Keep `Approve plan <id>` as a pre-confirmation request, but make its response
state plainly that nothing changed and show the exact confirmation needed. The
confirmed request should call the existing canonical transition primitive once
for the plan and all task projections, preserving Markdown bodies and file
formatting. The coordinator guidance should describe this as a bounded
status-only operation and should not encourage a sequence of hand-edits.

On an instantiated or wrapped workspace, the confirm-approval handoff says
approval is complete, execution has not started, and `Run approved plan <id>`
is the separate next request.

On `product-source`, after that same status-only transition, if the dirty set
is exactly the approval projection, present the existing
`commit-approved-plan` card in the same session. Do not commit, do not create a
worktree from dirty files, and do not wait for `Run approved plan <id>` to
surface `MAINTAINER_APPROVAL_COMMIT_REQUIRED`. The exact confirmation remains
`Confirm commit of the approved plan state.` Unrelated dirty files stay
`DIRTY_BASE_BLOCKED`. After that commit, `Run approved plan <id>` is the next
explicit request.

Tests should assert the absence of pre-confirmation mutation, the exact
post-confirmation status diff, the `product-source` follow-on commit card, and
that wrapped workspaces do not gain a commit card, then run the complete
semantic suite.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| AGF-001 | Clear, session-bound approval card and two-stage interaction | — |
| AGF-002 | Single bounded confirmed transition, then product-source commit card | AGF-001 |
| AGF-003 | Regression suite and independent acceptance evidence | AGF-001, AGF-002 |

## Acceptance criteria

- AGF-AC-01: The first approval request is mutation-free and clearly asks for confirmation.
- AGF-AC-02: Confirmation produces only the intended plan/task status projection.
- AGF-AC-03: The confirmed path avoids worktree/lease/implementation work and repeated repair churn.
- AGF-AC-04: The confirm-approval response distinguishes approval from execution.
- AGF-AC-05: Existing semantic and release boundaries remain valid.
- AGF-AC-06: On `product-source`, the existing commit card is presented immediately after confirmed approval when the dirty set is the exact status projection.

## Verification

The verification IDs and canonical commands are owned by `plan.yaml`. An
independent verifier must inspect the route decision, current gate state,
working-tree diff, task file bytes, and complete acceptance output. Timing
should be evaluated through bounded operation counts and absence of redundant
work; a hard wall-clock threshold must not make remote-host verification
flaky.

## Risks, assumptions, and open decisions

- Removing the confirmation turn would violate the consequential-gate contract and is out of scope.
- Folding the maintainer commit into `Confirm approval of plan <id>` is out of scope; the commit stays a separate exact confirmation.
- Auto-commit and creating a worktree from a dirty `product-source` base are out of scope.
- Remote host latency varies, so the plan measures structural work rather than promising a fixed duration.
- The final implementation must use existing owner files and must not add approval policy to a skill or role file.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.
The maintainer commit of an approved plan's status projection remains a
separate gate from this plan's own delivery.

## Provenance

Prepared from the observed Codex session, the shipped approval route,
transition primitive, and maintainer-commit gate, the plan/lifecycle/gate
contracts, coordinator guidance, and the existing routing, lifecycle, gate,
ownership, and semantic acceptance suites.
