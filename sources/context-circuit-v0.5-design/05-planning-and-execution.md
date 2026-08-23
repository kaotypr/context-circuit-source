# Planning and execution details

This document connects agent-oriented Product Knowledge retrieval, readable
plans, one-worker execution, independent verification, commits, and repair.

## 1. Plan creation brief

Before writing a plan, the root/coordinator creates a short internal brief:

~~~yaml
request: Add recurring billing
objective: Customer can create, view, and cancel a subscription
repositories:
  - id: api
    reason: Owns billing state and API
  - id: web
    reason: Owns customer billing UI
product_knowledge:
  - id: domain.billing
    path: context/domains/billing/README.md
    reason: Governs billing behavior.
  - id: architecture.system
    path: context/ARCHITECTURE.md
    reason: Maps service ownership.
  - id: conventions.project
    path: context/CONVENTIONS.md
    reason: Defines testing and implementation conventions.
context_grounding:
  summary: The API owns billing state and the web application presents it.
  constraints: [billing-events-append-only]
  decisions: [api-owns-subscription-state]
sources:
  - sources/billing-requirements.md
known_decisions:
  - context/DECISIONS.md#api-owns-subscription-state
open_questions:
  - trial conversion event behavior
~~~

The brief is not the plan and does not need to become a durable runtime record.
It prevents the root/coordinator from writing a plan before understanding scope.

The root/coordinator should retrieve knowledge by the request's concepts, domains,
repositories, decisions, and constraints. It should not treat the existence of
Markdown files as proof that it understands the project.

## 1.1 Context grounding and drift check

Before the plan becomes ready for review, compare the draft against the
selected Product Knowledge:

- objective against product purpose and domain rules;
- repository map against ownership and architecture;
- tasks against accepted constraints and decisions;
- acceptance against promised behavior;
- verification against the required evidence layer;
- implementation assumptions against open questions.

The result is recorded as stable context IDs, paths, reasons, and a short
grounding summary.
The plan also records request coverage: how the original requirements,
constraints, desired behavior, and non-goals are represented by the plan and
tasks. The root/coordinator does not discard detail merely to make the plan shorter. It
removes repetition, not meaning.
Contradictions become visible risks or questions. A plan must not be approved
as if it were aligned when its objective or scope conflicts with accepted
project context.

## 2. Readable plan requirements

PLAN.md should contain:

1. Title and objective.
2. Current status and execution summary.
3. Repositories and the behavior each owns.
4. Product Knowledge and source references.
5. Scope and explicit non-goals.
6. Task sequence and dependencies.
7. Acceptance criteria.
8. Verification approach and evidence layer.
9. Risks, assumptions, and open decisions.
10. Expected commits and delivery notes.
11. Human decisions still required.
12. Expected Product Knowledge impact, or an explicit statement that no durable
    context change is currently expected.

The plan should be as detailed as necessary to preserve the request and make
execution unambiguous. Detail means concrete behavior, repository ownership,
task inputs and outputs, dependencies, acceptance, verification, and stop
conditions; it does not mean repeating the same explanation in every task.

The human should not need plan.yaml to understand the intended work.

## 3. Machine-readable plan requirements

plan.yaml owns:

- stable plan ID;
- title and objective;
- human status;
- repository map;
- context references;
- task IDs and dependencies;
- allowed paths;
- declared effects;
- acceptance and verification IDs;
- repair limit;
- plan-level stop conditions;
- expected Product Knowledge impact and context units to reassess at completion.

Task Markdown files may explain rationale and implementation notes, but they do
not own lifecycle status or duplicate canonical identifiers.

## 4. Task design

A good task is independently understandable without duplicating the whole plan.

~~~yaml
id: API-001
title: Add subscription lifecycle API
repositories: [api]
paths:
  - src/subscriptions
  - test/subscriptions
depends_on: [CONTRACT-001]
context:
  - context/domains/billing/README.md
changes:
  - Add create, retrieve, and cancel operations.
acceptance:
  - id: API-AC-001
    statement: Active subscriptions can be created and retrieved.
verification:
  - id: API-VT-001
    command: test/subscriptions.sh
effects:
  - modify-source
  - run-tests
stop_conditions:
  - Do not change payment provider configuration.
~~~

Split tasks when they have different repository ownership, evidence layers,
dependencies, or human decisions. Do not split merely to create more workers.

Each task should be detailed enough that the worker can understand what to
change and how the result will be judged without reconstructing the request
from conversation history. A task records its intended behavior, concrete
change, affected surfaces, inputs, outputs, dependencies, context references,
acceptance, verification, risks, open questions, and stop conditions. It may
leave implementation choices to the worker when the plan does not need to
prescribe them.

## 5. Plan review

Review checks:

- objective is concrete;
- every important detail of the original request is represented or explicitly
  listed as an open question, assumption, or risk;
- repositories are registered and mapped;
- each task has paths or explicit scope;
- dependencies form a valid order;
- Product Knowledge references exist and are relevant;
- the plan records why the selected knowledge grounds its objective and scope;
- no unresolved context contradiction is hidden in an implementation choice;
- acceptance describes the desired result;
- verification can observe the required evidence layer;
- effects and stop conditions are clear;
- delivery is not hidden inside implementation;
- open questions are visible.

Review is an interactive plan discussion. The agent should be able to explain
why each task exists, which context supports it, what remains uncertain, and
what would change if an open question were answered differently.

Review reports risks and gaps. It does not edit or approve the plan.

## 6. Approval readiness

A plan is ready for approval when:

- the human objective is represented accurately;
- repository mappings are explicit;
- no task depends on an unknown repository;
- context conflicts are resolved or clearly accepted as risks;
- acceptance and verification are testable;
- worker and verifier scope is bounded;
- status is draft.

Approval changes only status. It does not create worktrees or start agents.

## 7. Execution brief

The host adapter constructs a compact execution brief from the approved plan:
The repository `anchor_branch` and its captured base commit are resolved from
the user's local binding at execution start, not authored as portable plan
scope. The same recorded anchor branch is available as the default pull-request
target for a later explicit delivery request.

~~~yaml
execution_id: exec-0001-billing-v2-001
plan: 0001-billing-v2
plan_snapshot: .runtime/executions/0001-billing-v2/exec-0001-billing-v2-001/plan-snapshot.yaml
role: worker
repositories:
  api:
    worktree: .runtime/worktrees/0001-billing-v2/api
    branch: cc/0001-billing-v2/api
    anchor_branch: development
    base_commit: abc123
    allowed_paths: [src/billing, test/billing]
  web:
    worktree: .runtime/worktrees/0001-billing-v2/web
    branch: cc/0001-billing-v2/web
    anchor_branch: development
    base_commit: xyz789
    allowed_paths: [app/billing, test/billing]
tasks: [BILL-001, BILL-002]
context:
  - context/domains/billing/README.md
  - context/ARCHITECTURE.md
acceptance: [BILL-AC-001, BILL-AC-002]
verification: [BILL-VT-001, BILL-VT-002]
attempt: 1
worker_failures: 0
commit_required: true
stop_conditions:
  - undeclared repository or path
  - unsafe external effect
  - missing required context
~~~

The brief is generated data. It is not a new policy owner or replacement for
the plan. `plan-snapshot.yaml` is an immutable serialized copy of PLAN.md,
plan.yaml, and all task files captured before execution. Worker and verifier
prompts use this snapshot rather than requiring the active plan directory to
remain in place.

## 8. Worker prompt

The worker prompt should be equivalent to:

~~~text
You are the sole writer for approved plan 0001-billing-v2.

Read:
- the immutable plan snapshot named by the execution brief;
- listed Product Knowledge pages;
- repository instructions in each assigned worktree.

Work:
- execute all tasks in dependency order;
- modify only listed paths in assigned worktrees;
- run declared implementation checks;
- commit each changed repository before handoff;
- write the worker handoff with files, commits, tests, assumptions, and blockers.

Do not:
- modify the anchor checkout;
- modify plan approval or completion status;
- rewrite or accept Product Knowledge;
- expand repository or path scope;
- claim independent verification;
- merge, push, publish, deploy, or delete work.

Stop and report if the plan or scope is insufficient.
~~~

The host may add provider-native details, but the product contract is the same.

## 9. Commit and handoff

The worker handoff records:

- execution and attempt IDs;
- repository commits;
- tasks completed;
- tests run and results;
- assumptions;
- known limitations;
- suggested verifier focus;
- whether a repair is being performed.

The worker commits before writing the final handoff. The runtime records hashes
and verifies that they belong to assigned branches.

## 10. Verifier prompt

The verifier prompt should be equivalent to:

~~~text
You are an independent read-only verifier for approved plan 0001-billing-v2.

Inspect:
- latest worker commit for each assigned repository;
- plan and task acceptance criteria from the immutable execution snapshot;
- worker handoff as a claim;
- declared verification commands and evidence layers;
- relevant Product Knowledge.

Run each canonical verification command and bounded additional checks needed to
establish the declared acceptance.

Report each acceptance and verification ID, observed evidence and layer, passed,
failed, or blocked outcome, exact failure and repair recommendation, and the
repository and commit checked.

Do not modify product files, plan files, Product Knowledge, or worker commits.
~~~

The verifier's handoff is the independent result. It must not merely repeat that
the worker said tests passed.

## 11. Verification outcomes

A verifier returns:

- passed: required evidence observed at the required layer;
- failed: evidence contradicts or does not satisfy the requirement;
- blocked: required observation could not be performed;
- waived: a human limitation decision exists, but it is not an unqualified pass.

Only passed satisfies execution verification.

## 12. Repair attempt

A repair begins with failed verification IDs, exact evidence or command output,
latest commits, allowed repair paths, and remaining failure count.

The worker must not use a repair failure as permission to redesign the plan. If
repair requires new scope, execution stops and the human is asked to change the
plan.

## 13. Three-failure sequence

The execution record should make this sequence obvious:

~~~yaml
attempts:
  - number: 1
    worker:
      status: committed
      commits: {api: abc123, web: xyz789}
    verifier:
      outcome: failed
      failures: [BILL-VT-002]
    worker_failures: 1
  - number: 2
    worker:
      status: committed
      commits: {web: def456}
    verifier:
      outcome: failed
      failures: [BILL-VT-002]
    worker_failures: 2
  - number: 3
    worker:
      status: committed
      commits: {web: ghi789}
    verifier:
      outcome: failed
      failures: [BILL-VT-002]
    worker_failures: 3
execution_status: failed
~~~

At three failures, no automatic fourth attempt is allowed.

## 14. Human completion

A successful verifier creates completion evidence:

~~~yaml
execution_status: verified
plan_status: approved
human_completion: pending
latest_commits:
  api: abc123
  web: xyz789
~~~

When the human asks to complete the plan, the runtime checks the latest
execution ID and commit set. It changes the plan status to done, preserves the
execution evidence, and creates an implementation completion record.

The workspace then reconciles the actual commits and verifier evidence against
the Product Knowledge used by the plan. The result is either no durable update
needed, or one or more proposals for changed, new, retired, stale, or
conflicting context. The plan's completion does not automatically accept those
proposals.

~~~yaml
execution_status: verified
plan_status: done
human_completion: accepted
implementation_completion:
  execution_id: exec-0001-billing-v2-001
  commits: {api: abc123, web: xyz789}
knowledge_impact:
  status: review-needed
  proposals: [context-impact-0001-billing-v2-001]
~~~

The human can say “review context updates for 0001-billing-v2” and explicitly accept
or defer each proposal. A deferred proposal remains visible to future plan
creation when its context is relevant.

## 15. Plan changes after failure

If the human changes scope, repository map, acceptance, or verification after a
failure, the existing execution cannot silently continue as if nothing changed.
The workspace runtime records the relationship between old execution and new plan
revision. A new execution may be required.

A repair may continue in the same execution only when intended scope and
acceptance remain unchanged.
