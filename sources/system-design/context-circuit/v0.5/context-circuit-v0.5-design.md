# Context Circuit v0.5 Design

Status: authoritative source design for v0.5
Revision: 8 — 2026-08-23

This document is the single source design for Context Circuit v0.5. It defines
the product, workspace model, readable plan format, conversational workflow,
multi-repository execution, worker/verifier behavior, runtime responsibilities,
and acceptance criteria.

It is source-repository design material. It is not a workspace plan, does not
approve or execute implementation work, and does not authorize commits, merges,
publication, deployment, or deletion outside the explicitly requested design
source cleanup.

The detailed design modules in this folder expand the overview: agent behavior,
human interaction, workspace and repository connections, Product Knowledge
lifecycle, planning and execution, runtime responsibilities, source/template
identity, terminology, and examples with acceptance scenarios. The overview
remains normative; the modules explain how to apply its decisions.

See [README.md](./README.md) for the reading order.

## 1. Product definition

Context Circuit is a universal project workspace for AI-assisted
work. A workspace can describe and coordinate one repository or multiple
repositories. It holds an agent-oriented Product Knowledge system and readable
plan files so an AI agent can find and understand the product context, create a
plan that does not drift from it, and execute an approved plan without
reconstructing Context Circuit's own lifecycle. The files are the storage and
human review surface; agent grounding is the primary purpose.

The product experience is a normal conversation:

~~~text
human request
  → agent reads relevant workspace knowledge and plan material
  → agent drafts or reviews a readable plan
  → human approves the plan in conversation
  → the workspace prepares isolated repository worktrees
  → one worker implements the whole plan and commits each repository
  → independent verifier checks the latest commits
  → worker repairs failures and creates new commits when needed
  → human decides when the plan is complete
  → workspace records implementation completion and reconciles Product Knowledge
~~~

The workspace hides branches, worktrees, runtime records, and verifier setup. It
does not hide the plan, the worker's handoff, the verifier's result, or the
human's authority over plan completion.

## 2. Problem

The current Context Circuit execution model makes agents spend substantial
time operating Context Circuit rather than doing project work. In the CC-010
trace, the request-to-verifier handoff took about 53 minutes. Roughly 32
minutes were spent on route reconstruction, packet selection, leases,
transaction graphs, child setup, isolation checks, and verifier process rules.
Only the remaining time was product implementation and verification.

The problem is not that plans, isolation, or independent verification are
useless. The problem is that the agent is asked to perform a long sequence of
process validations and confirmations before it can act on a plan.

Context Circuit v0.5 must make the project and its plan easy for an agent to
understand while keeping important safety properties behind a small,
deterministic Context Circuit runtime.

## 3. Product goals

1. Provide one universal project workspace for any project, including projects
   spread across multiple repositories.
2. Make Product Knowledge easy for an AI agent to find, retrieve, understand,
   and apply when creating plans, using structured indexes, scoped references,
   and provenance. Keep the underlying representation readable for humans.
3. Translate each request into a detailed plan and task set that preserves the
   request's requirements, constraints, desired behavior, non-goals, and
   unresolved decisions.
4. Preserve the current readable plan-file style as the primary human artifact.
5. Map every planned repository change explicitly to a registered repository.
6. Require plan approval before execution, while expressing approval as normal
   conversational intent rather than a confirmation-card ceremony.
7. Execute an approved plan with one worker handling all tasks in one bounded
   execution.
8. Create one branch and one isolated worktree per affected repository behind
   the scenes.
9. Commit worker changes to each repository before verification.
10. Keep an independent verifier that checks the worker's latest commits without
   changing product files.
11. Allow the worker to repair verifier failures within the same execution,
    creating a new commit for every repair attempt.
12. Stop after three worker failures while preserving branches, commits,
    worktrees, handoffs, and evidence.
13. Let the human decide when an execution result is sufficient to mark the
    plan complete, then record the implementation and assess its Product
    Knowledge impact.
14. Keep pull-request creation, merge, push, publication, deployment, and
    destructive cleanup as separate human-requested actions.
15. Make the Context Circuit runtime small enough that agents never need to load its
    implementation to perform project work.
16. Give every plan a consistent, easy-to-mention identifier and keep archived
    plans out of normal agent context until explicitly restored.

## 4. Non-goals

Context Circuit v0.5 does not aim to:

- expose a long lifecycle or process checklist to the human;
- require confirmation cards for ordinary plan approval or execution;
- make agents repeatedly validate Context Circuit's own packet, route, lease,
  or transaction rules in conversation;
- replace readable plan files with an opaque database or generated graph;
- let an unapproved plan execute;
- allow a worker to write into the anchor repository checkout;
- combine the worker and verifier into one actor;
- let the verifier repair product files;
- automatically mark a plan complete after verification;
- automatically accept Product Knowledge updates merely because a plan was
  implemented, verified, or marked complete;
- automatically open pull requests, merge, push, publish, deploy, archive, or
  delete work;
- embed a provider CLI, model SDK, authentication flow, or scheduler in the
  workspace runtime;
- store credentials, provider payloads, or private host authentication state;
- make this design document part of an instantiated project workspace.

## 5. Product principles

### 5.1 Product Knowledge grounds agent planning

Product Knowledge exists primarily to give the agent a reliable understanding
of the product and project before it creates or changes a plan. Its files are
not a general Markdown archive. They are structured, indexed, scoped, and
provenance-rich knowledge units that remain readable to humans.

Plans reference the knowledge units that constrain their objective, repository
scope, decisions, and acceptance. Runtime records support execution; they do
not replace the plan or Product Knowledge.

### 5.2 Conversation is the control surface

Humans use ordinary language. Internal states may be represented in files, but
the human should be able to say:

- “Create a plan for adding billing.”
- “Review plan 0001-billing-v2.”
- “Approve plan 0001-billing-v2.”
- “Then approve plan 0001-billing-v2 and execute it.”
- “Repair the failed verification for 0001-billing-v2.”
- “Mark 0001-billing-v2 complete.”
- “Archive plan 0001-billing-v2.”
- “Restore plan 0001-billing-v2.”

The adapter resolves the named plan and requested action. It must not infer
approval or completion from a vague orientation message.

### 5.3 Safety is hidden, not removed

Branches, worktrees, repository binding checks, commit tracking, verifier
independence, and recovery state remain deterministic. They are runtime
behavior, not a conversation that the agent must rediscover step by step.

### 5.4 One worker, one verifier

One worker owns the complete approved plan execution. One independent verifier
checks the result. A worker may be resumed for repair, but there is never more
than one active writer for an execution.

### 5.5 Human completion remains meaningful

Verification produces evidence that the execution is ready for a human
decision. It does not decide that the human's plan is complete.

### 5.6 Completion reconciles implementation with Product Knowledge

When a human marks a verified plan done, the workspace records an
implementation completion record linking the plan revision, execution, latest
commit for each affected repository, verifier result, and human completion
request. This makes the implementation a known project event rather than only
an old plan status.

Completion also triggers knowledge reconciliation. The agent compares the
actual committed changes and verifier evidence with the Product Knowledge that
grounded the plan and identifies durable knowledge that was added, changed,
invalidated, or left unaffected. The result is a no-update-needed result, a
context update proposal, or a stale/conflict warning.

The workspace must not silently rewrite accepted Product Knowledge from code,
worker claims, verifier success, or plan completion. A human may review and
accept a proposed context update in a separate conversational action. The plan
may be done while a knowledge proposal is pending; the pending impact remains
visible to future planning and is surfaced whenever it affects a new request.

### 5.7 One rule has one owner

The design, contracts, runtime, host adapters, worker role, verifier role, and
plan files must not duplicate policy. Each rule is defined once and other
artifacts reference its effect.

## 6. Workspace model

A workspace is a directory containing the template-provided product files, an agent-oriented Product
Knowledge index and knowledge units, plans, repository identity, and private
runtime state. It may coordinate any number of registered repositories,
including zero before project setup.

~~~text
workspace-root/
├── README.md                         # human quick start
├── .gitignore                        # ignores local repositories and runtime state
├── AGENTS.md                         # small host/source safety adapter
├── WORKFLOW.md                       # concise conversation and role guidance
├── workspace.yaml                    # canonical workspace and repository identity
├── context/                          # Product Knowledge
│   ├── INDEX.md
│   ├── WORKSPACE.md
│   ├── PROJECT.md
│   ├── ARCHITECTURE.md
│   ├── CONVENTIONS.md
│   ├── DECISIONS.md
│   ├── SOURCES.md
│   ├── sources.yaml
│   ├── domains/
│   ├── roles/
│   └── proposals/                    # pending context update proposals
├── sources/                          # passive raw evidence and briefs
├── plans/
│   ├── INDEX.md                       # active plans only
│   ├── <plan-id>/
│   │   ├── PLAN.md
│   │   ├── plan.yaml
│   │   └── tasks/
│   └── .archived/                     # archived plans; never normal context
├── .runtime/                         # private execution evidence
│   ├── executions/
│   └── worktrees/
├── repositories/                     # default ignored project-repository checkouts
└── repositories.local.yaml           # ignored host-local bindings
~~~

The portable workspace records logical repository identity. Machine-specific
paths remain in the ignored local binding file. The workspace must never scan
the filesystem to guess which repository a plan means.

### 6.1 Repository identity

workspace.yaml records a logical key, optional credential-free canonical URL,
and optional `default_branch` for every repository. `default_branch` is only a
portable clone and setup hint; it is not the user's active branch or the
execution base.

repositories.local.yaml maps a logical key to an explicit local path. The
runtime validates that the path exists, is a Git repository, matches the
declared identity when identity evidence is available, is not an unsafe
symlink or traversal path, and is not ambiguous. Each executable local binding
also sets the user's `anchor_branch`, the actual local active branch used as
the runtime worktree base. When the workspace root is a Git repository, the
reserved logical repository ID `workspace` binds that root at `.` and has its
own independent `anchor_branch`.

The default local destination for connected project repositories is the
gitignored `repositories/<repository-id>/` path. An existing checkout may use
another explicitly configured path. A new or empty repository can be created
there with an explicit `git init` setup, but it needs an initial anchor commit
before execution can create a worktree from it.

### 6.2 Product Knowledge

context/ contains accepted, durable knowledge about the project. It can include
product purpose, architecture, conventions, decisions, domain rules, roles,
and provenance. Its organization is optimized for agent retrieval: each
knowledge unit has a stable topic or ID, a short summary, applicable domains or
repositories, relevant decisions and constraints, freshness, and provenance.

The files remain readable and reviewable by humans, but the workspace is not a
Markdown storage product. The purpose of the files is to let an agent locate
and apply the right project knowledge when creating a plan.

context/INDEX.md is an agent retrieval catalog as well as human navigation. It
maps concepts, aliases, repositories, domains, decisions, and constraints to
the relevant knowledge units. The agent reads the selected units for a plan and
does not scan every context file or the Context Circuit implementation to understand
the project.

### 6.3 Source material

sources/ stores passive raw evidence and authored briefs or PRDs. Source
selection is explicit when a request names source material. The workspace must
not scan all sources as a fallback for missing Product Knowledge.

### 6.4 Source-side self-hosting and template-runtime tests

`context-circuit-source` uses nearly the same workspace-facing surfaces as the
product template: its own Product Knowledge, plans, repository mappings, and
`.runtime/`. This lets the Cc source project use Cc to plan and execute changes
to itself. Source-side knowledge and runtime evidence remain private to the
source project and are not template content.

The source repository also owns `template-harness/`, a behavior laboratory
that assembles the built `context-circuit-template` and runs realistic project
prompts and small implementation fixtures against isolated initialized
workspaces. It verifies the product experience at the template boundary,
including context gathering, detailed planning, review, approval, execution,
independent verification, repair, failure limits, and human completion. Its
fixtures and generated runtime state are source-only and must not be copied into
the product template.

## 7. Plan model

The plan remains the central execution contract. It has a human-readable
PLAN.md and a machine-readable plan.yaml. The two files describe the same
plan; plan.yaml owns canonical identifiers, status, repository mappings,
dependencies, and verification identifiers.

### 7.1 Plan directory

Every plan is workspace-level so a single plan can map changes across multiple
repositories:

~~~text
plans/
├── INDEX.md                           # active plan index
├── <plan-id>/
│   ├── PLAN.md
│   ├── plan.yaml
│   └── tasks/
│       ├── 001-<slug>.md
│       └── 002-<slug>.md
└── .archived/<plan-id>/               # inactive plan storage
    ├── PLAN.md
    ├── plan.yaml
    └── tasks/
~~~

The active index is `plans/INDEX.md`. It lists only plans directly under
`plans/`; it never lists or links archived plans. An archived plan is not part
of normal plan discovery or context loading.

Each active index entry contains the canonical plan ID, readable title, plan
status, one-line objective, affected repository IDs, and relative plan path.
The plan ID is the first and primary lookup value so an agent and human can
mention a plan consistently without relying on a long title.
Lifecycle status changes update the active row while the plan is active;
archiving removes the row instead of creating an archived-status row.

### 7.1.1 Consistent plan naming

Every plan has a stable `plan_id` used in conversation, paths, runtime records,
branches, worktrees, context proposals, and human-facing summaries. New plan
IDs use this form:

~~~text
NNNN-<plan-short-title>
~~~

The ID must be:

- exactly four decimal digits in `NNNN`, zero-padded and assigned as the next
  workspace sequence after the highest sequence ever allocated;
- a non-empty lowercase ASCII kebab-case `<plan-short-title>`;
- concise, concrete, and easy to pronounce in a sentence;
- free of spaces, underscores, punctuation, dates, timestamps, model names,
  status words, and random hashes;
- unique across both active and archived plans for the lifetime of the
  workspace;
- stable after creation, even if the plan title or status changes.

Examples include `0001-billing-v2`, `0002-saved-checkout-sessions`, and
`0003-checkout-v2`. The sequence is allocated when a draft is created and is
never reused, including after archive or restore.

The readable title may be a normal sentence, but the agent always presents and
accepts the stable ID as the canonical mention key: “review plan
0001-billing-v2.” If a request does not provide a suitable short
title, the agent proposes the next sequence and a short title while creating
the draft. It asks only when the proposed ID would be ambiguous or materially
misrepresent the request.
~~~

The human should be able to understand the outcome, scope, repository map,
task sequence, acceptance evidence, risks, and current status from PLAN.md.

### 7.2 Canonical plan fields

~~~yaml
schema_version: 1
plan: 0001-billing-v2
title: Billing v2
status: draft
objective: Add recurring billing support across the API and dashboard.
repositories:
  - id: api
    purpose: Billing API and persistence
  - id: dashboard
    purpose: Customer billing interface
product_knowledge:
  - id: domain.billing
    path: context/domains/billing/README.md
    reason: Governs subscription lifecycle and billing ownership.
  - id: architecture.system
    path: context/ARCHITECTURE.md
    reason: Identifies API and dashboard boundaries.
context_grounding:
  summary: Billing state is owned by the API and rendered by the dashboard.
  constraints: [billing-events-append-only]
  decisions: [api-owns-subscription-state]
knowledge_impact:
  expected_context_units: [domain.billing, architecture.system]
  review_on_completion: true
tasks:
  - id: BILL-001
    title: Add billing API
    repositories: [api]
    paths: [src/billing, test/billing]
    depends_on: []
    changes:
      - Add subscription lifecycle endpoints.
    acceptance:
      - id: BILL-AC-001
        statement: Subscription creation persists an active subscription.
    verification:
      - id: BILL-VT-001
        command: test/api-billing.sh
  - id: BILL-002
    title: Add dashboard billing view
    repositories: [dashboard]
    paths: [app/billing, test/billing]
    depends_on: [BILL-001]
    changes:
      - Display current subscription and cancellation state.
    acceptance:
      - id: BILL-AC-002
        statement: A customer can see the current subscription state.
    verification:
      - id: BILL-VT-002
        command: test/dashboard-billing.sh
execution:
  worker: one
  independent_verifier: required
  max_worker_failures: 3
~~~

In `plan.yaml`, the `plan` field is the canonical stable plan ID described as
`plan_id` in prose. Other files and runtime records reference that same value;
they must not invent a second display name or execution-specific plan name.

The exact field names may be refined during implementation, but these
properties are fixed:

- every task names the repository or repositories it can change;
- every task names bounded paths or an explicit repository-wide scope;
- task dependencies are explicit;
- every task describes its intended behavior, concrete change, affected
  surfaces, inputs, outputs, acceptance, verification, and stop conditions;
- acceptance and verification are distinct identifiers;
- every Product Knowledge reference has a stable ID, path, and reason it
  grounds the plan;
- the plan records a short context-grounding summary, constraints, and
  decisions;
- the plan declares expected Product Knowledge impact, or explicitly records
  that no durable impact is currently expected;
- execution policy is plan data, not hidden role prose;
- the worker failure limit is three for one execution.

### 7.3 Plan status and execution status

Plan status is human-controlled:

~~~text
draft → approved → done
~~~

Execution state is runtime evidence and does not replace plan status:

~~~text
not-started → running → verifying → verified
                         ↘ repairing → verifying
                         ↘ failed
                         ↘ blocked
~~~

After verification passes, execution becomes verified, but the plan remains
approved. Only an explicit human request to mark the plan complete may change
the plan status to done, and only when the runtime has a successful
verification result. That transition writes an implementation completion
record and starts Product Knowledge reconciliation; it does not automatically
accept a context update.

If execution fails or reaches the three-failure limit, the plan does not become
done. Its branches, commits, worktrees, and failure evidence remain available
for human review.

## 8. Conversational workflow

The workspace exposes a small conversational adapter. It can use a classifier
internally, but the user experience must not expose a two-stage router, context
packet construction, or confirmation-card protocol.

### 8.1 Initialize a workspace

When a user asks to initialize or set up a workspace, the agent creates or
updates the minimal workspace structure, records project identity, and
registers repositories the user names. It asks only for information that
cannot be safely inferred. It does not create plans or repositories the user
did not request.

### 8.2 Create a plan

The agent reads relevant Product Knowledge, repository instructions, and
request-scoped sources. It expands the request into a detailed readable plan with explicit repository
mapping, tasks, acceptance criteria, verification, risks, and stop conditions.
The plan preserves the request's requirements, desired behavior, constraints,
non-goals, assumptions, open questions, and risks instead of compressing them
into a vague implementation outline.

Before the plan is ready for approval, the agent reassures it against the
available Product Knowledge, repository evidence, and named sources. Every
important requirement is grounded by evidence, represented as an assumption,
or recorded as an open question or risk.

Creating a plan does not approve or execute it.

### 8.3 Review a plan

Review is a non-executing, status-preserving discussion. The agent explains
whether the plan is understandable, complete, internally consistent, and
executable across its mapped repositories. It walks through request coverage,
task details, context references, open questions, risks, assumptions, and
verification. When the human resolves a question, corrects a requirement, or
requests more detail, the agent may update the draft plan content. Review never
approves or executes the plan and never changes its status.

### 8.4 Approve a plan

Approval is a human gate expressed as an explicit conversational request. A
request such as “approve plan 0001-billing-v2” changes the plan from draft to
approved only after deterministic readiness checks pass.

There is no separate approval card or required confirmation token. If the user
asks to execute an unapproved plan, the agent refuses execution and explains
that the plan must be approved. A follow-up such as “then approve plan
0001-billing-v2 and execute it” is sufficient to perform the two explicit actions in
sequence, provided approval checks succeed.

### 8.5 Execute a plan

“Execute plan X” is a separate authorization from plan approval unless the
same request explicitly asks for both. Execution is allowed only when the plan
status is approved and all repository bindings and anchor branches are valid.

The workspace prepares isolated worktrees, the host launches the worker, and the
worker executes the complete plan. The user does not need to confirm each
task, packet, lease, branch, verifier, or repair step.

### 8.6 Inspect results

The agent can summarize worker commits, verifier evidence, repair attempts,
remaining failures, and the current human decision. Runtime records are
available for detailed inspection but are not required reading for normal use.

### 8.7 Mark a plan complete

The human may say “mark plan X complete” after execution is verified. The
runtime validates the latest execution evidence and changes the plan status to
done. A verification pass alone never performs this transition. After the
status change, the workspace records the implementation completion and reports
whether Product Knowledge needs review, for example:

> Plan 0001-billing-v2 is now done. Its verified implementation is recorded across
> api and web. Product Knowledge reconciliation found one proposed update to
> domain.billing and one architecture reference that is now stale. Say
> “review context updates for 0001-billing-v2” to discuss them.

### 8.8 Archive and restore a plan

“Archive plan X” is an explicit organization action. The workspace resolves X
from `plans/INDEX.md`, moves `plans/X/` to `plans/.archived/X/`, preserves every
plan file and its existing status, and removes X from the active index.

Archiving does not validate or verify plan status, execution status, commits,
branches, worktrees, open questions, or Product Knowledge impact. It does not
mark a plan done, stop an execution, clean runtime state, merge changes, or
delete evidence. It is only a human-requested move out of normal plan context.

Archive and restore are atomic plan-organization operations. The runtime takes
a short-lived plan-organization lock, checks only the exact source path and
target-path/index collision, moves the directory, and updates
`plans/INDEX.md` atomically. If any filesystem or index step fails, the source
directory and active index remain unchanged. These checks never inspect plan or
execution status.

If an execution already exists, the move does not interrupt it or rename its
runtime records. The worker and verifier continue from their execution brief
and plan snapshot. The plan must be restored before the normal agent can
inspect or discuss that archived execution through the plan files.

The normal agent must not read an archived plan or traverse
`plans/.archived/`. “Restore plan X” is the explicit exception: the workspace
moves `plans/.archived/X/` back to `plans/X/`, preserves its plan status and
files, and adds it to `plans/INDEX.md`. Restore does not approve, execute,
complete, or otherwise validate the plan.

## 9. Multi-repository execution

The plan is the authority for repository scope. The worker receives a compact
execution brief containing:

- the immutable plan snapshot and task map;
- mapped repositories and their worktree roots;
- task order and dependencies;
- Product Knowledge references;
- allowed paths and declared effects;
- acceptance and verification identifiers;
- commit expectations and stop conditions.

The worker may change only mapped paths in the assigned worktree for each
repository. A task that needs another repository must declare it in the plan;
the worker must not discover or modify undeclared repositories by searching the
workspace.

For every affected repository, the runtime creates a deterministic branch and
worktree from the validated `anchor_branch` tip. The plan has one worker, but
Git history remains repository-specific.

The worker commits each repository after implementation is complete and before
the verifier starts. If a repair changes a repository, the repair creates a new
commit on that repository branch. The worker does not amend the previous
implementation commit merely to hide repair history.

## 10. Worker behavior

The worker is one bounded writer role for the entire plan execution.

It must:

- read the plan and named Product Knowledge;
- execute tasks in dependency order;
- work only inside assigned repository worktrees;
- run the plan's implementation checks where appropriate;
- commit every affected repository after implementation;
- write a concise handoff listing changes, commits, tests, assumptions, and
  unresolved concerns;
- stop when the plan requires undeclared scope, an unsafe action, or a blocked
  prerequisite.

The worker must not:

- edit the anchor repository checkout;
- change plan approval or completion status;
- mark its own work verified;
- alter verifier evidence;
- silently expand repository or path scope;
- rewrite a prior commit to conceal a repair attempt.

The worker may be resumed for repair. “One worker” means one writer role and
one execution ownership boundary, not that a host must keep one process alive
for every attempt.

## 11. Independent verifier behavior

The verifier is a separate actor and receives the latest worker commits,
read-only worktree access, the plan, the worker handoff, and canonical
verification commands or evidence references.

It must:

- inspect the latest commit in every affected repository;
- replay the promised acceptance and verification evidence;
- check that evidence proves the required layer;
- check that repository and path scope was respected;
- report passed, failed, or blocked with evidence references;
- write only its own verifier result and handoff.

It must not:

- modify product files;
- repair the worker's implementation;
- change plan status;
- treat a worker claim as independent evidence;
- downgrade an evidence requirement because a host lacks a capability.

If the host cannot create an independent verifier with the required read-only
capability, execution is blocked. The worker or coordinator must not
self-verify as a substitute.

## 12. Repair loop and failure limit

After a verifier failure, the workspace runtime passes the failure evidence to the worker
within the same execution. The worker repairs the affected implementation and
creates a new commit for every repository it changes. The verifier then checks
the latest commits again.

The counter is worker_failures. Each worker result that the independent
verifier rejects increments it, including the initial implementation result.
The maximum is three. A verifier result blocked by host capability or missing
external evidence is a blocked execution, not an invented worker failure.

~~~text
worker commit 1 → verifier failed → worker failure 1
worker repair commit 2 → verifier failed → worker failure 2
worker repair commit 3 → verifier failed → worker failure 3 → stop
~~~

When the limit is reached:

- no further worker or verifier attempt starts automatically;
- the plan remains approved or its current human-controlled status;
- all commits, worktrees, handoffs, and verifier evidence remain preserved;
- the agent reports the smallest human decision needed to continue.

## 13. Completion, delivery, and cleanup

Verification and completion are separate concepts.

verified means the independent verifier passed the latest commits against the
plan's acceptance and verification requirements. It does not mean the human
considers the product work complete.

done means the human explicitly asked to mark the plan complete after seeing
the available evidence. The runtime must reject that request when the latest
execution has not passed independent verification.

Merge, push, pull-request creation, publication, deployment, archive, and
cleanup are separate actions. They may use verified branches and runtime
evidence, but none is implied by worker success, verifier success, or plan
completion.

When a human requests a pull request for an implemented plan, each execution
branch is the source and the affected repository's recorded `anchor_branch` is
the default target. A different target branch must be named explicitly; the
workspace never substitutes `default_branch` or silently follows a moving
remote branch.

Archiving is not delivery, completion, or cleanup. It does not require plan
verification or status validation and does not remove runtime evidence.

Cleanup must inspect dirty or unpushed work and preserve it unless the human
explicitly requests its removal. A failed execution is never cleaned up as a
side effect of reporting failure.

### 13.1 Implementation completion and knowledge reconciliation

The completion transition has two distinct outputs:

1. an implementation completion record owned by the runtime, proving which
   verified execution and repository commits the human accepted as done;
2. a knowledge-impact result owned by the context lifecycle, describing whether
   the implementation changes the durable understanding of the project.

The reconciliation input includes the final plan and tasks, the plan's
grounding references, changed paths and commits in every affected repository,
worker handoffs, verifier evidence, and relevant current context revisions. It
checks for:

- new product or architecture facts;
- changed decisions, constraints, conventions, or domain rules;
- context claims whose cited implementation is now stale or invalid;
- durable implementation patterns worth recording;
- no durable knowledge change.

Each finding is stored as a context update proposal or an explicit
no-update-needed result. Proposals identify the target context unit, proposed
operation, evidence, affected repository and commit, confidence, and any
conflict. Human acceptance updates the context page, provenance, and index
together. Deferral leaves the plan done but keeps the proposal visible to
future plan creation.

## 14. Runtime model

Runtime state is private, minimal, and execution-oriented. It is evidence, not
Product Knowledge or plan intent.

~~~text
.runtime/
├── executions/<plan-id>/<execution-id>/
│   ├── execution.yaml
│   ├── plan-snapshot.yaml          # complete PLAN.md/plan.yaml/task snapshot
│   ├── handoff.md
│   ├── completion.yaml
│   ├── context-impact.yaml
│   ├── attempts/001/
│   │   ├── worker.yaml
│   │   └── verifier.yaml
│   ├── attempts/002/
│   │   ├── worker.yaml
│   │   └── verifier.yaml
│   └── repositories/<repository-id>.yaml
├── locks/<plan-id>/owner.yaml
└── worktrees/<plan-id>/<repository-id>/
~~~

Every execution record must identify:

- plan and execution IDs;
- template/runtime version;
- repository and worktree bindings;
- anchor branch, captured base commit, and current commit revisions;
- worker and verifier attempts;
- verification result and evidence references;
- failure count;
- current execution state;
- implementation completion state and human completion request;
- Product Knowledge impact result and proposal references;
- timestamps and integrity digests where needed for recovery.

Records are written atomically. A partial or contradictory record cannot grant
ownership, resume a writer, prove verification, or authorize completion.

## 15. engine.sh and runtime refactor

The current runtime at wrapper/runtime/engine.sh is approximately 123KB and 2,523
lines. It currently combines conversation routing, action selection, context
packet loading, context budgets, repository binding, worktree handling,
lifecycle transitions, gates, runtime graphs, evidence checks, and recovery.

That monolithic shape is not the v0.5 product contract.

### 15.1 Target role

engine.sh should be a small host-neutral runtime library for deterministic
workspace, Git, and execution-state operations. Agents must never need to read
its full implementation to understand or execute a project plan.

An initial implementation target is roughly 20–40KB for the core. If the runtime
cannot remain within that shape without becoming unsafe or unreadable,
execution coordination must move to a separate runtime module. The size target
is an engineering guardrail, not a user-facing budget.

### 15.2 Engine contents

The refactored runtime may contain:

- safe workspace-relative path and identifier checks;
- atomic file writes and content digests;
- workspace and repository binding validation;
- clean-anchor-checkout, branch, and worktree preparation;
- plan structure and approval-state validation;
- execution and attempt record creation;
- repository commit capture and revision validation;
- verifier-result and evidence-record validation;
- the three-failure counter and repair eligibility;
- simple ownership locking for one active writer;
- interruption, resume, and recovery inspection;
- completion eligibility validation without automatically changing status;
- active plan-index maintenance and exact archive/restore moves without plan
  status validation;
- implementation completion records and context-impact handoff references
  without interpreting or accepting Product Knowledge.

### 15.3 Engine exclusions

The runtime must not contain:

- provider-specific Codex, Cursor, or Claude launch code;
- model prompts or Product Knowledge interpretation;
- plan-writing intelligence;
- brittle conversation keyword routing as the source of lifecycle policy;
- packet-budget exploration or large context-set construction;
- confirmation-card rendering or confirmation-token handling;
- product implementation or test-runner logic;
- pull-request creation, merge, push, publication, deployment, or destructive
  cleanup;
- automatic plan completion.

### 15.4 Logical engine use

The runtime is used at these logical points:

1. Workspace entry or resume: validate identity and existing execution state
   when the request needs more than read-only orientation.
2. Initialization: create or validate the minimal workspace structure.
3. Repository registration: resolve and validate a named local binding.
4. Plan creation or review: validate plan structure and repository mappings;
   the agent writes plan content.
5. Active plan index maintenance and exact archive/restore moves without plan
   status or execution validation.
6. Approval: atomically change draft to approved after explicit human intent
   and readiness checks.
7. Execution preflight: validate approval, dependencies, clean anchor
   checkouts, and repository mappings.
8. Worktree preparation: create and record per-repository branches and
   worktrees.
9. Worker start and handoff: create or update the execution record; the host
   adapter launches the worker.
10. Commit capture: record the commit produced by the worker for each affected
   repository.
11. Verifier preparation: validate latest revisions and read-only verifier
    scope; the host adapter launches the verifier.
12. Repair: record verifier failure, increment the worker-failure counter, and
    permit another worker attempt when the limit allows it.
13. Execution result: record verified, failed, or blocked and preserve all
    evidence.
14. Human completion: validate successful evidence and perform the explicit
    plan-status transition to done.
15. Implementation completion and context-impact handoff: record the accepted
    commits and preserve the reconciliation references.
16. Recovery and cleanup: inspect ownership and preserve work unless an
    explicit human action authorizes cleanup.

The engine may be sourced once by a host adapter or test process, but that is
an implementation detail. Logical use is limited to deterministic operations
such as the points above.

## 16. Responsibility boundaries

| Concern | Owner | Result |
| --- | --- | --- |
| Product purpose and domain rules | context/ Product Knowledge | Agent understands the project. |
| Human plan intent and status | PLAN.md and plan.yaml | Human-readable plan and canonical status. |
| Conversation interpretation | Host adapter/coordinator | Normal language becomes a requested action. |
| Workspace/repository safety | Runtime engine and workspace contracts | Explicit bindings and safe roots. |
| Branch/worktree preparation | Runtime engine | Isolated repository execution areas. |
| Active plan index and archive/restore placement | Runtime engine and plan storage | Active plans remain discoverable; archived plans remain outside normal context. |
| Worker launch | Host adapter | Provider-native child creation. |
| Product implementation | Worker | Changes and commits in assigned worktrees. |
| Independent verification | Verifier | Read-only evidence and result. |
| Repair counting | Execution runtime | Maximum of three worker failures. |
| Completion decision | Human request plus plan lifecycle | Plan changes to done only explicitly. |
| Product Knowledge reconciliation | Root conversational agent/coordinator | Proposals or no-update-needed results from completed implementation. |
| Pull-request creation, merge, push, publication, deployment | Separate delivery action | No automatic external effect; delivery must use an explicit target and available source branch. |

Host identity and provider capability are bounded execution evidence. They do
not authorize approval, execution, a role, a lease, verification, or completion.

## 17. Host adapters and agent roles

Codex, Claude, and Cursor remain transports. Their adapters provide the native
child-agent primitive and report whether the required worker or verifier can be
created.

The root conversational agent is the coordinator for normal workspace
interaction. Context gathering, Product Knowledge reconciliation, plan
creation, plan review, approval interpretation, completion handling, and
delivery discussion are root/coordinator capabilities, not separate child
agents. The only product child roles in an execution are the single worker and
the independent verifier.

The host adapter:

- receives the conversational request;
- asks the runtime for current plan/workspace state;
- launches one worker with the execution brief;
- launches the independent verifier with latest revisions;
- routes verifier failures back to the worker within the same execution;
- reports runtime results in normal language.

The adapter must not create a second product policy or silently bypass the
runtime. If a required verifier child cannot be created, the result is
blocked; the root agent cannot self-verify.

Role instructions remain thin. The plan and runtime execution brief provide the
task, repository, paths, acceptance, verification, effects, and stop
conditions. Roles explain how to behave inside that bounded scope; they do not
restate the entire lifecycle.

## 18. Failure and recovery behavior

The Context Circuit runtime fails closed for:

- missing or unapproved plan;
- missing or ambiguous repository binding;
- dirty or changed anchor checkout when a clean anchor is required;
- unsafe path or symlink traversal;
- branch/worktree creation failure;
- a second live writer for the same plan execution;
- undeclared repository or path scope;
- unavailable independent verifier;
- contradictory or stale runtime evidence;
- three worker failures.

On interruption, preserve the execution record, branches, worktrees, commits,
handoffs, and verifier results. A safe resume reuses the current execution only
when ownership and revisions still match. A new execution must not erase or
silently replace old evidence.

An execution failure is reported with:

- latest worker and verifier commits;
- failure count;
- exact acceptance or verification evidence that failed;
- whether another repair is allowed;
- smallest human decision required next.

## 19. Security and repository boundaries

- Workspace-relative paths reject traversal and unsafe symlinks.
- Credentials remain in host Git configuration or the host agent and never
  enter workspace files or runtime records.
- The anchor checkout remains untouched during worker implementation.
- Each repository worktree is assigned to one execution writer.
- Verifier access to product files is read-only.
- A task cannot expand repository or path scope without an explicit plan change
  and approval.
- Runtime files are private evidence and are not Product Knowledge.
- The workspace does not scan sibling workspaces, unrelated plans, or arbitrary
  repositories to fill context gaps.

## 20. Compatibility and migration

The current readable plan artifacts are valuable and should be retained. A
migration may add repository mappings and execution fields to existing plan
files, but must not silently change their objective, task meaning, acceptance,
or human status.

Existing runtime state may be read for migration or recovery. It must not be
treated as current authority when its ownership, revisions, or schema cannot be
validated. Migration preserves branches, worktrees, commits, handoffs, and
failure evidence.

An updated product template may replace template-owned instructions and runtime
code, but it
must not overwrite Product Knowledge, sources, plans, local repository
bindings, product repositories, or active work.

## 21. Acceptance criteria

Context Circuit v0.5 is acceptable when all of the following are demonstrated:

1. A new workspace can register zero or more repositories without embedding
   machine-specific paths in portable identity files.
2. An agent can locate and understand the Product Knowledge relevant to a plan
   by concept, domain, repository, decision, and constraint without scanning
   the entire workspace; the underlying knowledge remains human-readable.
3. A plan remains understandable from PLAN.md and maps every change to one or
   more repositories.
4. A plan and its tasks preserve the important details of the original request,
   with requirements, constraints, desired behavior, non-goals, assumptions,
   open questions, risks, and context grounding visible.
5. An unapproved plan cannot execute.
6. “Approve plan X and execute it” performs approval and execution as explicit
   conversational actions without a confirmation-card protocol.
7. Execution creates one isolated branch/worktree per affected repository.
8. One worker can implement all tasks in dependency order across repositories.
9. The worker commits each affected repository before verification.
10. An independent verifier checks the latest commits without modifying product
   files or plan status.
11. A verifier failure can initiate worker repair in the same execution.
12. Every repair produces a new commit rather than amending away prior work.
13. The third worker failure stops the execution and preserves all evidence.
14. A verifier pass produces verified execution evidence but does not mark the
   plan complete.
15. An explicit human completion request changes the plan to done only after
   successful independent verification.
16. Pull-request creation, merge, push, publication, deployment, and cleanup
    require separate human actions.
17. An agent can execute a plan without reading the full runtime
   implementation.
18. The refactored runtime contains no provider-specific child-launch logic.
19. Interrupted executions can resume or stop read-only without losing commits,
   worktrees, or evidence.
20. Missing verifier capability blocks execution rather than allowing
   self-verification.
21. The semantic acceptance suite covers multi-repository mapping, approval,
   worktree isolation, commits, verifier independence, repair, failure limit,
   human completion, and recovery.
22. Marking a verified plan done records its implementation commits and
    verifier evidence, reconciles the actual changes against Product Knowledge,
    surfaces proposals or stale/conflicting context, and never silently accepts
    a context update.
23. New plans receive consistent, stable, human-friendly IDs that are used as
    their canonical conversation and filesystem names.
24. Archiving moves a named plan to `plans/.archived/`, removes it from
    `plans/INDEX.md`, preserves its files and status, performs no plan-status or
    execution validation, and keeps it unread by agents until explicit restore.
25. An execution whose plan is archived continues from its immutable plan
    snapshot without reading the archived plan directory.
26. `default_branch` remains portable clone/setup guidance; execution validates
    the user-local `anchor_branch` and creates worktrees from its captured tip.
27. The workspace root may be registered as the logical `workspace` repository
    with its own anchor branch, while new project repositories default to the
    gitignored `repositories/<repository-id>/` path and support explicit clone
    or `git init` setup.
28. A human-requested pull request uses each repository's execution branch as
    the source and its recorded `anchor_branch` as the default target; delivery
    blocks when the source, configured provider, or target branch is unavailable
    rather than inferring a remote, pushing silently, or using `default_branch`.

## 22. Implementation order

The implementation should proceed in these bounded phases:

1. Plan and workspace contract: preserve readable plan files, add explicit
   repository mappings, and define execution and human-controlled status fields.
2. Product Knowledge context: establish the workspace context index and
   references used by plan creation and execution.
3. Runtime reduction: separate low-level workspace/Git/state primitives from
   current routing, packet, and lifecycle ceremony; shrink or split engine.sh.
4. Conversational adapter: support create, review, approve, execute, inspect,
   repair, complete, archive, and restore requests without confirmation cards.
5. Multi-repository execution: prepare branches/worktrees and launch one
   worker with a complete execution brief.
6. Independent verification: launch the read-only verifier against the latest
   repository commits and record evidence.
7. Repair loop: add new repair commits, failure counting, and the three-failure
   stop behavior.
8. Recovery and delivery boundaries: preserve interrupted work and keep merge,
   publication, deployment, and cleanup separate.
9. Semantic verification: run complete traces for one repository, multiple
   repositories, successful verification, repair, third failure, blocked
   verifier, interrupted execution, human completion, and archive/restore
   without plan-status validation.

Each phase must be independently reviewable. This design does not authorize
implementation, delivery, or publication by itself.

## 23. Final design decisions

The following decisions are fixed for v0.5:

- Context Circuit is a universal multi-repository project workspace.
- Product Knowledge is optimized first for agent retrieval and plan grounding;
  its file representation remains readable for humans.
- Plan approval is required before execution.
- Approval and combined approval/execution are normal conversation, without
  confirmation cards.
- The plan explicitly maps every change to repositories.
- The workspace root may be an affected repository with the reserved ID
  `workspace` and its own anchor branch.
- `default_branch` is only portable clone/setup guidance.
- Each executable local repository binding has a user-selected
  `anchor_branch`, which becomes the runtime worktree base and default
  pull-request target.
- New project repositories default to the gitignored `repositories/` folder;
  explicit clone and `git init` setup are supported.
- The workspace creates branches and worktrees behind the scenes.
- One worker executes all tasks in one execution.
- The worker commits each repository before verification.
- An independent verifier checks the latest commits read-only.
- Worker repairs happen in the same execution and create new commits.
- Three worker failures stop the execution.
- Verification never automatically marks a plan complete.
- The human decides when to change the plan status to complete.
- New plans use stable IDs in the form `NNNN-<plan-short-title>`, where the
  sequence is never reused and the suffix is lowercase kebab-case.
- Archiving is an explicit file move to `plans/.archived/`, removes the plan from
  the active index, does not validate plan status, and does not imply
  completion, delivery, cleanup, or deletion.
- engine.sh is refactored into a small host-neutral runtime layer and does not
  contain provider launch or conversational policy.
- Pull-request creation, merge, push, publication, deployment, and cleanup
  remain separately requested actions; a pull request defaults to the
  repository's recorded anchor branch as target.
