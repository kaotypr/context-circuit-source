# Runtime engine and host integration

This document defines the internal runtime boundary for v0.5. It answers what
engine.sh does, when it is used, and what remains outside it.

## 1. Runtime boundary

The runtime is deterministic support for workspace and Git operations. It is
not the product's conversation, planning intelligence, or provider scheduler.

~~~text
conversation
  → host adapter resolves requested action
  → runtime validates state and prepares bounded resources
  → host adapter launches worker or verifier
  → runtime records revisions and evidence
  → host adapter coordinates repair or reports the result
~~~

The runtime may be implemented as engine.sh plus smaller modules. The filename
is not a product requirement; the boundary is.

## 2. Current monolith

The current engine is approximately 123KB and 2,523 lines. It combines:

- conversation probing and action selection;
- context packet and byte-budget logic;
- repository identity and binding;
- worktree preparation;
- plan lifecycle transitions;
- gates and confirmation behavior;
- runtime records and graphs;
- evidence matching;
- leases and recovery.

This makes it difficult to reason about and tempts agents to load it as a
substitute for a task-specific execution brief.

## 3. v0.5 runtime contents

The reduced runtime may contain:

### Workspace and repository primitives

- workspace identity validation;
- safe identifier and path validation;
- local repository binding resolution;
- Git repository identity checks;
- clean-anchor-checkout checks;
- deterministic branch and worktree creation;
- branch and worktree revision inspection.

### Plan and execution primitives

- plan structure validation;
- approved-status validation;
- atomic draft-to-approved transition;
- execution snapshot creation;
- task and repository scope validation;
- attempt and failure counter management;
- commit capture and branch revision validation;
- verifier result validation;
- completion eligibility validation;
- safe active-plan index maintenance and exact plan archive/restore moves,
  without reading or validating plan status for those moves;
- implementation completion record creation linking the verified execution,
  plan revision, commits, and human request;
- recording the context-impact handoff or proposal references without
  interpreting or accepting Product Knowledge.

### Recovery primitives

- atomic record writes;
- simple one-writer ownership locks;
- interrupted execution inspection;
- resume eligibility;
- preservation and cleanup safety checks.

## 4. Runtime exclusions

The runtime must not contain:

- provider-specific agent launch code;
- model prompts;
- Product Knowledge interpretation;
- plan-writing behavior;
- a broad conversation router as a policy owner;
- context packet combinatorics or model-token budgets;
- confirmation-card rendering;
- product test semantics;
- plan archive/restore status gates or verification;
- pull-request creation, merge, push, publication, deployment, or automatic
  completion.

The runtime may execute a declared repository command only when it is part of a
bounded verification request. It does not decide which product tests are
correct; the plan and verifier do.

## 5. Conceptual operations

Names are illustrative; the contract matters more than shell function names.

~~~text
workspace_validate(root)
repository_resolve(root, repository_id)
repository_preflight(root, plan_id)
worktree_prepare(root, execution_id, repository_id)
plan_validate(plan_dir)
plan_approve(plan_dir)
plan_archive(root, plan_id)
plan_restore(root, plan_id)
execution_begin(plan_dir)
attempt_begin(execution_dir, attempt)
worker_commit_record(execution_dir, repository_id, revision)
verifier_result_record(execution_dir, attempt, result)
repair_allowed(execution_dir)
completion_ready(execution_dir)
recovery_inspect(execution_dir)
~~~

These operations return structured results and reason codes. They do not print
long process instructions for the agent.

`plan_archive` and `plan_restore` acquire a short-lived plan-organization lock,
check only exact source/target paths and index collisions, then update the plan
directory and `plans/INDEX.md` as one recoverable operation. A failed move or
index write leaves the original directory and index entry intact. They never
read plan status or execution status.

## 6. Logical use times

The runtime is logically used at these points:

1. Workspace entry or resume when state must be checked.
2. Workspace initialization when deterministic files must be created.
3. Repository registration and local binding validation.
4. Plan creation or review for schema and repository-map checks.
5. Active plan index maintenance and explicit archive/restore file moves. These
   operations do not inspect plan or execution status.
6. Plan approval for the atomic status transition.
7. Execution preflight for approval, dependencies, clean anchor checkouts, and
   scope.
8. Branch and worktree preparation for each affected repository.
9. Worker start and handoff record creation.
10. Commit capture after worker implementation.
11. Verifier preparation and read-only scope validation.
12. Repair recording and three-failure enforcement.
13. Final execution result and evidence recording.
14. Human completion status transition.
15. Implementation completion record and context-impact handoff recording.
16. Recovery, interruption, and cleanup inspection.

The runtime is not logically used to make the agent read Product Knowledge,
choose a product solution, write a plan, launch a provider child, or decide
that the human's plan is complete. The root/coordinator performs knowledge
reconciliation from the recorded evidence; the runtime only preserves its
completion and proposal references.

## 7. Host adapter boundary

The host adapter:

- receives normal-language requests;
- resolves plan IDs and action intent;
- asks runtime for current state;
- prepares provider-native worker and verifier launch;
- passes the execution brief;
- returns child results;
- asks runtime to record commits and evidence;
- coordinates repair attempts.

The adapter must not bypass missing approval, change the plan to done after a
verification pass, or self-verify when the verifier child is unavailable.

The runtime must not contain Codex, Claude, Cursor, MCP, model, or provider
authentication code.

## 8. Records

A minimal execution record has:

~~~yaml
schema_version: 1
execution_id: exec-0001-billing-v2-001
plan: 0001-billing-v2
plan_revision: sha256:...
status: verifying
worker_failures: 1
repositories:
  api:
    worktree: .runtime/worktrees/0001-billing-v2/api
    branch: cc/0001-billing-v2/api
    anchor_branch: development
    base_commit: abc123
    latest_commit: def456
attempts:
  - number: 1
    worker_commit_record: attempts/001/worker.yaml
    verifier_result: attempts/001/verifier.yaml
~~~

Records are written atomically. A partial record cannot grant ownership, resume
a writer, prove verification, or permit completion.

## 9. Read-only verifier enforcement

The runtime declares verifier scope and validates that the requested role does
not include product writes. The host adapter must provide actual read-only
capability when the host supports it.

If the host cannot guarantee verifier read-only behavior, the adapter reports
blocked. A prompt saying “do not edit” is not enough for a required verifier.

The verifier may write its own evidence outside product worktrees when the host
supports that isolated activity. It may never write product files, plan intent,
or worker commits.

## 10. Locking and concurrency

The runtime needs enough ownership state to prevent two active writers from
changing the same plan execution. A simple atomic lock can contain:

- execution ID;
- owner/session ID;
- repository set;
- acquired timestamp;
- current attempt;
- recovery state.

Archive and restore use a separate short-lived plan-organization lock. It
protects the directory move and active-index rewrite only; it does not inspect
or change plan status, execution status, or execution ownership.

A competing writer receives a read-only or blocked result. The runtime never
silently steals a live lock.

Different plans may use separate worktrees and run concurrently. Shared
repository overlap is reported as a later integration risk.

## 11. Recovery

On interruption, preserve:

- execution record;
- current attempt;
- worker and verifier handoffs;
- branches and worktrees;
- commit history;
- failure count.

Resume is allowed only when plan revision, configured anchor branches and
captured anchor commits, worktree paths, and ownership still match. Otherwise
execution is read-only and the human must decide whether to recover or start a
new execution.

## 12. Engine size and split rule

The initial engineering target for the core runtime is approximately 20–40KB.
This is not a model context budget; it is a maintainability signal.

If engine.sh grows beyond that because it contains unrelated concerns, split the
runtime by responsibility instead of adding compressed process prose. Likely
modules are:

- engine.sh: portable safety and file/Git primitives;
- execution.sh: execution records, attempts, and repair state;
- repository.sh: binding and worktree operations.

Agents still receive generated execution briefs, not these implementation files.

## 13. Tests for the runtime boundary

Runtime tests must prove:

- an unapproved plan cannot start;
- a dirty or mismatched repository cannot start;
- each repository gets its own worktree and branch;
- undeclared paths fail;
- worker commits are recorded before verification;
- a repair creates a new commit record;
- the third worker failure stops execution;
- verifier writes are rejected;
- verification does not change plan status;
- explicit completion requires passed evidence;
- explicit completion records the accepted commits and starts a context-impact
  handoff without silently changing Product Knowledge;
- archive moves preserve plan files and status, remove the active index entry,
  update the index atomically, and do not validate plan or execution status;
- archive/restore rejects only an exact-path or index collision and leaves all
  files unchanged when the move cannot complete;
- interruption preserves all state;
- provider-specific launch logic is absent from the runtime.
