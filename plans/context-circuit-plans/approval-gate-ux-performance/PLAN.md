# Fast and explicit plan approval confirmation flow

Status: draft  
Repository: context-circuit-source  
Source: repository-evidence

## Review summary

The observed approval session behaved safely but felt slow and ambiguous. The
initial `Approve plan <id>` request spent time checking state before returning a
card, and the confirmed turn took nearly three minutes while manually changing
one plan status and six task projections, running acceptance checks, and
repairing unintended Markdown end-of-file drift.

This plan keeps the two-step human gate and the separate execution trigger,
while making the card explicit, routing the confirmed action through one
canonical status transition, and adding regression coverage for the exact
failure mode.

## What approval authorizes

Approval authorizes bounded implementation of the tasks and acceptance
criteria in `plan.yaml`. It does not authorize execution, lease acquisition,
worktree creation, Git changes, delivery, publication, deployment, merge, or
cleanup.

## Scope and non-goals

The implementation is limited to the approval route/card, the confirmed
status-transition path, coordinator procedure, human-facing wording, and
tests. The canonical router, lifecycle, ownership, and human-gate owners stay
unchanged. Existing dirty or untracked files, including unrelated plans and
source fixtures, remain untouched.

## Proposed solution

Keep `Approve plan <id>` as a pre-confirmation request, but make its response
state plainly that nothing changed and show the exact confirmation needed. The
confirmed request should call the existing canonical transition primitive once
for the plan and all task projections, preserving Markdown bodies and file
formatting. The coordinator guidance should describe this as a bounded
status-only operation and should not encourage a sequence of hand-edits.

The final handoff should distinguish `approved` from `executing` and point to
`Run approved plan <id>` as the separate next action. Tests should assert both
the absence of pre-confirmation mutation and the exact post-confirmation diff,
then run the complete semantic suite.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| AGF-001 | Clear, session-bound approval card and two-stage interaction | — |
| AGF-002 | Single bounded confirmed transition with no formatting drift | AGF-001 |
| AGF-003 | Regression suite and independent acceptance evidence | AGF-001, AGF-002 |

## Acceptance criteria

- AGF-AC-01: The first approval request is mutation-free and clearly asks for confirmation.
- AGF-AC-02: Confirmation produces only the intended plan/task status projection.
- AGF-AC-03: The confirmed path avoids worktree/lease/implementation work and repeated repair churn.
- AGF-AC-04: The handoff distinguishes approval from execution.
- AGF-AC-05: Existing semantic and release boundaries remain valid.

## Verification

The verification IDs and canonical commands are owned by `plan.yaml`. An
independent verifier must inspect the route decision, current gate state,
working-tree diff, task file bytes, and complete acceptance output. Timing
should be evaluated through bounded operation counts and absence of redundant
work; a hard wall-clock threshold must not make remote-host verification
flaky.

## Risks, assumptions, and open decisions

- Removing the confirmation turn would violate the consequential-gate contract and is out of scope.
- Remote host latency varies, so the plan measures structural work rather than promising a fixed duration.
- Maintainer approval still creates a predictable source diff and may still require the separate maintainer commit gate.
- The final implementation must use existing owner files and must not add approval policy to a skill or role file.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Prepared from the observed Codex session, the shipped approval route and
transition primitive, the plan/lifecycle/gate contracts, coordinator guidance,
and the existing routing, lifecycle, gate, and semantic acceptance suites.
