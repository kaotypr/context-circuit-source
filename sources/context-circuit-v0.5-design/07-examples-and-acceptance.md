# Examples and acceptance

This document gives concrete examples for implementation and review. Examples
are illustrative unless the main v0.5 design marks the behavior as fixed.

## 1. Minimal workspace

~~~text
workspace.yaml
context/INDEX.md
context/WORKSPACE.md
context/PROJECT.md
plans/INDEX.md
repositories.local.yaml
~~~

A workspace can start with no repositories. The agent can create Product
Knowledge and a plan before a repository is connected, but execution must wait
until every planned repository is registered and bound.

## 2. Multi-repository plan

~~~yaml
plan: 0003-checkout-v2
status: approved
repositories:
  - id: api
    purpose: Checkout API
  - id: web
    purpose: Checkout UI
  - id: contracts
    purpose: Shared checkout schema
tasks:
  - id: CONTRACT-001
    repositories: [contracts]
    paths: [checkout]
    depends_on: []
  - id: API-001
    repositories: [api]
    paths: [src/checkout, test/checkout]
    depends_on: [CONTRACT-001]
  - id: WEB-001
    repositories: [web]
    paths: [app/checkout, test/checkout]
    depends_on: [CONTRACT-001, API-001]
execution:
  worker: one
  independent_verifier: required
  max_worker_failures: 3
~~~

## 3. Context brief used to create the plan

~~~yaml
topic: checkout-v2
knowledge:
  - context/domains/checkout/README.md
  - context/ARCHITECTURE.md
  - context/DECISIONS.md#contract-ownership
sources:
  - sources/checkout-requirements.md
repositories:
  api:
    evidence: src/checkout at abc123
  web:
    evidence: app/checkout at xyz789
  contracts:
    evidence: checkout at qrs456
open_questions:
  - guest checkout expiration behavior
~~~

The root/coordinator should surface the open question in PLAN.md instead of silently
choosing an expiration policy.

## 4. Conversation: create and approve

~~~text
Human: Create a plan for adding saved checkout sessions.

Agent: I found the API, web, and contracts repositories. The current Product
Knowledge says the API owns checkout state and contracts owns shared schemas.
One question remains about guest-session expiration. I drafted plan
0002-saved-checkout-sessions with that question recorded.

Human: Review the plan.

Agent: The plan maps schema changes to contracts, API behavior to api, and the
UI to web. It preserves the requested saved-session behavior, but guest
expiration remains an open question. The plan also has a risk that API and web
deployment order affects the user journey. We can resolve those during review.

Human: Approve plan 0002-saved-checkout-sessions.

Agent: The plan is approved. It is not executing yet.

Human: Execute plan 0002-saved-checkout-sessions.

Agent: I am preparing one isolated worktree per repository and starting one
worker for all plan tasks. An independent verifier will inspect the commits
after implementation.
~~~

## 4.1 Conversation: archive and restore

~~~text
Human: Archive plan 0002-saved-checkout-sessions.

Agent: Plan 0002-saved-checkout-sessions was moved to `plans/.archived/` and removed from
the active plan index. Its files and status were preserved. No plan-status or
execution validation was performed.

Human: Restore plan 0002-saved-checkout-sessions.

Agent: Plan 0002-saved-checkout-sessions was restored to the active plan area and
added to the plan index. Its files and status are unchanged.
~~~

The archive operation does not inspect whether the plan is draft, approved,
done, running, failed, or verified. The agent cannot read the archived plan
again until the explicit restore action completes.

## 5. Conversation: combined request

~~~text
Human: This plan looks good; then approve plan 0002-saved-checkout-sessions and
execute it.

Agent: I will approve the named draft plan, then execute it if approval
preflight succeeds.

Agent: The plan is approved. Three repository worktrees are ready. The worker
has started.
~~~

No confirmation card is required, but the request is still explicit.

## 6. Conversation: unapproved execution

~~~text
Human: Execute plan 0002-saved-checkout-sessions.

Agent: I cannot execute it because the plan is draft. Say “approve plan
0002-saved-checkout-sessions and execute it” to approve and run it, or approve it
first and review the result.
~~~

No worktree or worker is created.

## 7. Successful execution

~~~text
Worker:
- contracts commit 111aaa
- api commit 222bbb
- web commit 333ccc
- implementation checks passed

Verifier:
- CONTRACT-AC-001 passed
- API-AC-001 passed
- WEB-AC-001 passed
- latest commits independently checked

Agent:
Execution is verified. The plan remains approved. Say “mark
0002-saved-checkout-sessions complete” when you decide the plan itself is complete.

Human: Mark 0002-saved-checkout-sessions complete.

Agent:
Plan 0002-saved-checkout-sessions is now done. The implementation record links
commits 111aaa, 222bbb, 333ccc and the independent verifier result. Product
Knowledge reconciliation found a proposed update to checkout ownership. Say
“review context updates for 0002-saved-checkout-sessions” to discuss it.
~~~

## 8. Repair execution

~~~text
Verifier:
WEB-VT-001 failed. The checkout UI does not display the API's expired-session
state.

Agent:
Verification failed. Worker failure 1 of 3 is available. I am sending the
failure to the worker for repair.

Worker:
web repair commit 444ddd

Verifier:
WEB-VT-001 passed against 444ddd.

Agent:
Execution is verified after one repair commit. The plan remains approved.
~~~

## 9. Third failure

~~~text
Verifier:
WEB-VT-001 failed again.

Agent:
Worker failure 3 of 3 reached. Execution has stopped. The latest web commit is
preserved, along with all previous commits and verifier evidence. The plan was
not marked complete.
~~~

## 10. Context update from a plan

~~~text
Human: What did we learn about checkout ownership?

Agent: The verified execution shows that the API owns expiration state and the
web application renders it. I can propose an update to
context/domains/checkout/README.md citing plan
0002-saved-checkout-sessions and commits 222bbb and 444ddd.
~~~

The proposal requires human acceptance before changing accepted Product
Knowledge.

## 11. Context update from a repository

~~~text
Human: Refresh the API architecture context.

Agent: I will inspect the API architecture files and tests at the current
revision, compare them with context/ARCHITECTURE.md, and propose only stable
facts with repository paths and revision evidence.
~~~

The agent does not copy the repository into context or update context silently.

## 12. Acceptance scenarios

The semantic acceptance suite should cover at least:

### Workspace and context

- initialize an empty workspace;
- bind the Git workspace root as logical repository `workspace` with its own
  local anchor branch;
- register one repository;
- register three repositories;
- clone a repository to the ignored `repositories/<repository-id>/` default
  path;
- initialize an empty repository there with `git init`, establish its anchor
  branch, and require an initial commit before execution;
- reject a missing or mismatched binding;
- gather context from named sources;
- propose a context update from a plan;
- reconcile a completed plan's commits and verifier evidence against Product
  Knowledge;
- propose a context update from repository evidence;
- detect a conflict with accepted decisions;
- use the accepted context brief to draft a plan;
- locate Product Knowledge by concept, domain, repository, decision, and
  constraint without scanning every context file;
- detect a plan objective or repository scope that drifts from accepted context;
- retain a grounding summary and context references in the plan.

### Plans and conversation

- assign new plans consistent, stable, human-friendly IDs;
- create a readable multi-repository plan;
- preserve the important details of the original request in the plan and tasks;
- create explicit open questions, assumptions, and risks when detail is not
  settled;
- review without approval or execution, while allowing human-requested draft
  content updates;
- use “review plan” to discuss and resolve plan questions before approval;
- reject execution of a draft plan;
- approve through normal conversation;
- approve and execute in one explicit request;
- preserve plan status after verification;
- complete only after explicit human request;
- record implementation completion when a plan becomes done;
- surface context proposals, stale references, or no-update-needed results after
  completion;
- accept or defer context proposals explicitly without changing plan completion;
- archive a named plan without validating its status;
- remove archived plans from `plans/INDEX.md` and exclude them from normal
  agent context;
- restore a named archived plan without changing its status.
- continue an already-started execution from its immutable snapshot after the
  plan is archived.

### Execution

- create one branch/worktree per repository;
- create each worktree from the repository's captured personal or team anchor
  branch;
- execute all tasks with one worker;
- commit each affected repository before verification;
- verify latest commits independently;
- repair and create a new commit;
- stop on third worker failure;
- preserve failed worktrees and evidence;
- use each repository's anchor branch as the default pull-request target when a
  human explicitly requests pull requests;
- block pull-request delivery when the source branch, configured provider, or
  target branch is unavailable instead of inferring a remote, pushing silently,
  or substituting `default_branch`;
- block when no independent verifier is available.

### Safety and recovery

- reject dirty anchor checkouts;
- reject undeclared paths and repositories;
- reject verifier product writes;
- preserve state after interruption;
- prevent concurrent writers for one plan;
- keep pull-request creation, merge, push, publication, deployment, and cleanup
  separate.

## 13. Implementation review checklist

A reviewer should be able to answer yes to all of these:

- Does the human see a simple conversation rather than a process checklist?
- Can one plan map changes to multiple repositories?
- Does the plan remain easy to read?
- Can the agent find the relevant Product Knowledge without scanning the whole
  context directory?
- Does the agent use Product Knowledge to ground the objective and repository
  scope when creating the plan?
- Does the plan expose context drift instead of silently choosing an
  implementation that conflicts with accepted knowledge?
- Does marking a plan done create a durable implementation record tied to the
  verified commits and trigger knowledge reconciliation?
- Are context updates proposed with provenance rather than silently applied?
- Are pending context impacts findable during future plan creation?
- Does one worker handle all tasks?
- Are branches and worktrees hidden but real?
- Does the worker commit before verification?
- Does every repair create a new commit?
- Is the verifier truly independent and read-only?
- Does the third worker failure stop the execution?
- Does verification leave plan completion to the human?
- Is engine.sh small, host-neutral, and unnecessary to load for project work?
- Are old runtime records and worktrees preserved on failure?
