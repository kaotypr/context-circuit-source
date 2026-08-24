# Human interaction

Context Circuit v0.5 uses ordinary conversation as its control surface. Files
and runtime state provide deterministic evidence, but humans do not operate a
visible checklist of internal lifecycle steps.

## 1. Conversation rules

The conversational adapter must:

- recognize a named plan or repository;
- distinguish inspect from mutate;
- distinguish approval from execution;
- allow explicit compound requests such as “approve and execute”;
- explain refusals in project language;
- show the smallest missing decision;
- never infer a consequential action from “okay,” “looks good,” or an old reply;
- never require a hidden confirmation token.

Project language means the effect, not the mechanism. The canonical mapping from
each internal term to its user-facing phrasing is doc 08 §7; the shipped
projection an instantiated workspace reads is `docs/terminology.md`. Internal
mechanism is revealed only when the user explicitly asks for diagnostics.

If two plans match, ask which plan. Do not inspect unrelated bundles to guess.

## 2. Normal request vocabulary

| Human says | Meaning |
| --- | --- |
| “What is this workspace?” | Read-only orientation. |
| “Gather context about checkout.” | Propose or update checkout context. |
| “Connect the api repository.” | Register the logical repository and bind an existing checkout; no clone or initialization is implied. |
| “Clone the api repository.” | Explicitly clone to `repositories/api/` by default and set its local `anchor_branch`. |
| “Initialize an empty api repository.” | Explicitly create `repositories/api/`, run `git init`, set its local `anchor_branch`, and establish the initial commit needed for execution. |
| “Create a plan for this feature.” | Draft a plan from current knowledge and evidence. |
| “Review plan 0003-checkout-v2.” | Non-executing discussion of request coverage, plan details, open questions, and risks; human-requested corrections may update the draft plan. |
| “Approve plan 0003-checkout-v2.” | Explicit approval gate. |
| “Execute plan 0003-checkout-v2.” | Execute only if already approved. |
| “Approve plan 0003-checkout-v2 and execute it.” | Approve, then execute if checks pass. |
| “What happened with 0003-checkout-v2?” | Summarize execution evidence. |
| “Repair the failed 0003-checkout-v2 verification.” | Start another worker attempt if allowed. |
| “Mark 0003-checkout-v2 complete.” | Human-controlled completion status change. |
| “Review context updates for 0003-checkout-v2.” | Discuss Product Knowledge impact found after the implementation was marked done. |
| “Accept the context update for 0003-checkout-v2.” | Explicitly accept a proposed Product Knowledge change after review. |
| “Open a pull request for 0003-checkout-v2.” | Separate delivery action using each execution branch as source and its repository `anchor_branch` as target. |
| “Archive plan 0002-saved-checkout-sessions.” | Move the named plan out of the active plan area without checking its status. |
| “Restore plan 0002-saved-checkout-sessions.” | Return the named archived plan to the active plan area without changing its status. |
| “Merge the checkout branch.” | Separate delivery or merge action. |

Natural variations are allowed, but the resulting action must map to the same
contract.

## 2.1 Plan review discussion

“Review plan X” starts a focused conversation about the plan rather than a
single pass/fail response. The agent should walk through, as needed:

- the original request and how each important detail is represented;
- objective, desired behavior, constraints, and non-goals;
- Product Knowledge and repository evidence used to ground the plan;
- repository and task mapping;
- task order, dependencies, and expected changes;
- acceptance and verification evidence;
- assumptions, open questions, and risks;
- delivery effects intentionally outside execution.

The human may resolve an open question, correct a requirement, change scope, or
ask for a task to be more specific. When the human does so, the agent updates
the draft plan and continues the discussion until the details are clear. Review
may change draft content, but it never changes plan status, approves, or
executes the plan. Approval remains a separate explicit request.

## 3. Approval interaction

Approval is explicit but conversational. The agent may say:

> Plan 0003-checkout-v2 is ready for approval. It covers the API and dashboard
> repositories and has two unresolved risks. Reply “approve plan 0003-checkout-v2” to
> approve it.

This is guidance, not a confirmation card. A clear equivalent is sufficient.

If the user says “execute 0003-checkout-v2” while the plan is draft, respond:

> I cannot execute 0003-checkout-v2 because it is still draft. You can approve it and
> then execute it by saying “approve plan 0003-checkout-v2 and execute it.”

No branch, worktree, worker, or verifier is started by the refusal.

## 4. Compound approval and execution

Handle a compound request sequentially:

1. Resolve the named plan.
2. Verify that it is a valid draft and ready for approval.
3. Change draft to approved.
4. Re-read the current approval state.
5. Run execution preflight.
6. Start the isolated execution.

If approval or preflight fails, stop and report the failure. Do not claim that
execution occurred merely because approval succeeded.

## 5. Execution interaction

Before execution starts, give a short summary of:

- plan name and objective;
- repositories and branches involved;
- number of tasks;
- worker and verifier roles;
- known execution limits.

The user does not need to confirm each repository, task, worktree, packet, lease,
or repair. The plan is the approved scope.

During execution, provide event-based updates:

- worktrees prepared;
- worker implementation committed;
- verifier started;
- verifier failed with specific evidence;
- repair committed;
- verifier passed;
- execution stopped after the third failure.

Do not narrate every shell command or internal record write.

## 6. Failure interaction

Every failure response says:

- what failed;
- whether product work changed;
- latest commit(s);
- whether another repair is allowed;
- what the human can decide next.

Example:

> Verification failed for BILL-VT-002 because the dashboard journey could not
> find the expected subscription state. Repair attempt 1 of 3 is available. The
> API commit is unchanged; the dashboard worktree contains the failed worker
> commit.

A verifier capability failure is different:

> The verifier could not be created with read-only capability. No product work
> was verified, and no self-verification was attempted. The execution is
> blocked; existing worktrees and commits are preserved.

## 7. Completion interaction

After a pass:

> Execution of 0003-checkout-v2 is verified. The latest commits are ... The plan is
> still approved. Say “mark 0003-checkout-v2 complete” when you decide the plan itself
> is complete.

On explicit completion:

> Plan 0003-checkout-v2 is now complete. Verification passed before the status change.
> The branches remain available; no merge or publication was performed. The
> implementation record found one Product Knowledge update candidate. Say
> “review context updates for 0003-checkout-v2” to discuss it.

Marking the plan complete records which verified implementation the human
accepted. It also starts knowledge reconciliation against the actual commits,
the verifier evidence, and the context used to create the plan. It does not
silently change accepted Product Knowledge. The agent reports either that no
durable update was found, or that proposals, stale references, or conflicts are
available for review.

The human may complete a plan while a context proposal is pending. Completion
means the human accepted the implementation outcome; accepting a Product
Knowledge update is a separate decision. A relevant pending proposal is
surfaced during future plan creation so the agent does not plan from silently
outdated context.

## 7.1 Archive and restore interaction

Archiving uses the plan ID, not a title fragment or a repository name. On
“archive plan X,” the agent moves the exact active directory from
`plans/X/` to `plans/.archived/X/` and removes its row from `plans/INDEX.md`.
It does not inspect or validate plan status, execution status, verification,
commits, branches, worktrees, open questions, or context impact. It does not
stop work, mark the plan done, or delete evidence.

The archived plan is unavailable to normal agent reading. On “restore plan X,”
the agent moves the exact directory back to `plans/X/`, restores its index row,
and preserves all files and status. Restore is not approval, execution,
completion, verification, or validation.

If the user asks to complete an unverified plan:

> I cannot mark 0003-checkout-v2 complete because its latest execution has not passed
> independent verification. The current result is failed or blocked.

## 8. Delivery interaction

Delivery is separate from completion. For a pull-request request, the
execution branch is the source and the repository's recorded `anchor_branch`
is the target by default:

> The verified work is on the api and dashboard execution branches. Opening
> pull requests would target the api and dashboard anchor branches. I have not
> opened them. If you want that action, ask me to open pull requests for the
> named plan.

For a merge request:

> The verified work is on the api and dashboard branches. Merging them will
> change the target branches. I have not merged them. If you want that action,
> ask me to merge the named branches.

The Context Circuit runtime never interprets verification as merge authorization.

## 9. Human review surfaces

The human should be able to inspect:

- plans/<plan-id>/PLAN.md;
- plans/<plan-id>/plan.yaml;
- task files;
- latest worker handoff;
- verifier result and evidence;
- repository branches and commit history;
- current execution summary.

The human should not manually edit worktree ownership, attempt counters,
verifier read-only records, runtime locks, or host capability evidence.

## 10. Ambiguity and clarification

Ask a focused question when:

- a plan name is missing or ambiguous;
- a repository key has no local binding;
- a task names an unknown repository;
- two branches could be the intended anchor branch;
- a context update conflicts with an accepted decision;
- a delivery target or effect is unclear.

An open question or risk is not a failure when the plan makes it visible and
the human can resolve it during review. It becomes a blocker when leaving it
unresolved would change the intended product behavior, repository scope, or
verification result.

Do not ask when a safe default is recorded by the plan or workspace. Show the
default in the response and continue when the human has already authorized it.

## 11. Human authority

Humans explicitly control:

- acceptance of Product Knowledge updates;
- plan approval;
- execution if it was not included in the approval request;
- scope changes;
- completion status;
- pull-request creation, merge, push, publication, deployment, archive, and
  cleanup.

The agent controls only bounded implementation and evidence work authorized by
the current plan and action.
