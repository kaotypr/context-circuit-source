# Agent behavior

This document defines how agents behave in a v0.5 workspace. It covers the
root conversational agent and the worker and verifier roles.

The agent should spend its context understanding the project and doing project
work. Product Knowledge retrieval is a first-class behavior: the agent should
find the relevant knowledge units by concept, domain, repository, decision, or
constraint before drafting a plan. It must not spend normal turns re-learning
the Context Circuit runtime.

## 1. Common operating loop

For every request:

1. Identify the workspace root and host role.
2. Read the small workspace entry files and the relevant Product Knowledge index.
3. Identify the named plan, repository, task, or source file if one exists; if
   no plan is named, use the active `plans/INDEX.md`.
4. Read only the active plan and context references needed for the requested
   action; execution, verification, repair, and resume use the immutable
   execution snapshot when one exists.
5. Ask one concise question only when a missing fact would change the action or
   create unsafe ambiguity.
6. Perform the requested read or write through the appropriate role.
7. Report what changed, what was verified, and what remains for the human.

The agent should not:

- read all Context Circuit contracts because a request mentions a plan;
- search sibling workspaces or unrelated plans for examples;
- infer approval, execution, completion, or delivery from a vague statement;
- use a worker claim as verifier evidence;
- broaden repository or path scope to avoid a focused question;
- load engine.sh as a substitute for an execution brief.

`plans/.archived/` is outside normal agent context. The agent must not read or
traverse it for orientation, plan discovery, review, execution, or context
gathering. An explicit “restore plan X” request is the only operation that may
access the named archived directory, and it moves the plan back to the active
plan area before normal reading resumes.

## 2. Context loading order

The root agent loads context in this order:

1. Host and system instructions.
2. Workspace AGENTS.md and WORKFLOW.md.
3. workspace.yaml identity and repository keys.
4. context/INDEX.md and relevant identity and project pages.
5. The named active plan's PLAN.md, plan.yaml, and required task files, or the
   immutable execution snapshot for an existing execution.
6. Product Knowledge references named by the plan or directly relevant to the
   request.
7. Repository-local instructions in the affected repository or worktree.
8. Runtime state only when the request concerns execution, repair, recovery,
   verification, or completion.
9. Request-scoped source files only when the human names them or the plan
   explicitly includes them.

The agent does not read all of context/ by default. The context index is a
retrieval catalog with the metadata needed to select relevant knowledge units;
linked pages are the evidence.

An active execution is self-contained after its snapshot is created. If its
plan is archived, the worker, verifier, repair flow, and recovery flow continue
from that snapshot and runtime evidence; they do not read the archived plan
directory.

## 2.1 Agent retrieval contract

Every Product Knowledge unit should expose enough metadata for an agent to find
it without guessing filenames:

- stable context ID or topic;
- concise summary;
- domain and repository associations;
- keywords and common aliases;
- applicable constraints and decisions;
- freshness and status;
- provenance references.

The agent first forms a retrieval question from the plan request, then uses the
catalog to select knowledge units. It records why each selected unit matters to
the plan. A broad directory scan is not a substitute for a useful retrieval
index.

## 3. Behavior by action

| Request | Agent behavior | Mutation |
| --- | --- | --- |
| Explain workspace | Summarize identity, repositories, Product Knowledge areas, plans, and current execution state if relevant. | None. |
| Gather context | Read named or selected evidence and propose a context update with provenance and conflicts. | Only after human acceptance if accepted knowledge changes. |
| Connect repository | Register a logical repository and bind an existing local checkout, including its user-local `anchor_branch`. | Workspace identity and local binding only. No clone or initialization unless separately requested. |
| Clone or initialize repository | Perform the explicitly requested clone or `git init` setup under the default ignored `repositories/` path, then record the `anchor_branch`. | External Git setup and local binding. No plan execution. |
| Create plan | Read relevant Product Knowledge and repositories, then draft PLAN.md, plan.yaml, tasks, and the active index entry. | Creates or updates a draft plan and `plans/INDEX.md`. |
| Review plan | Check scope, repository mapping, dependencies, acceptance, verification, and risks. Resolve the human's questions and corrections in the draft when requested. | Draft plan content only; never plan status, approval, execution, or completion. |
| Approve plan | Validate readiness and change draft to approved after explicit human approval. | Plan status and active index entry. |
| Execute plan | Validate approved status, prepare worktrees, launch one worker, and coordinate verification. | Runtime and isolated worktrees. |
| Verify plan | Launch or resume the independent read-only verifier. | Verifier evidence only. |
| Repair execution | Pass verifier failure evidence to the worker, allow a new commit, and verify again. | Worker repair commit and runtime evidence. |
| Mark complete | Confirm passed verification, record the accepted implementation, and reconcile its Product Knowledge impact after explicit human request. | Plan status, active index entry, implementation completion evidence, and pending context-impact result; accepted knowledge still requires a separate decision. |
| Archive plan | Move the named active plan to `plans/.archived/` and remove it from `plans/INDEX.md` without checking plan or execution status. | Plan files move; status and runtime evidence are preserved. |
| Restore plan | Move the named archived plan back to `plans/`, preserve its files and status, and re-add it to the active index. | Plan files move; no approval, execution, or completion. |
| Open pull request | Use each execution branch as the source and its recorded repository `anchor_branch` as the target. | Only after a separate human request. |
| Deliver or merge | Report the separate action and its effects. | Only after a separate human request. |

## 4. Orientation

For a read-only orientation request, answer from workspace identity and the
Product Knowledge index. State uncertainty instead of inventing project facts.

A useful orientation reports:

- workspace name and purpose;
- registered repository keys and whether bindings are available;
- important Product Knowledge areas;
- draft, approved, and active plans listed in `plans/INDEX.md`;
- current verified or failed executions when the user asks about state;
- one useful next action, without treating that recommendation as authorization.

Orientation does not create an execution record.

## 5. Gathering or refreshing context

When asked to gather context, define the evidence boundary first:

- named source files;
- a named repository and bounded paths;
- an existing plan;
- a Product Knowledge topic;
- or a focused project question.

Read the smallest useful evidence set, compare it with accepted context, and
write a proposal containing:

- proposed fact or change;
- affected context page;
- evidence references;
- observed date or revision;
- confidence and unresolved ambiguity;
- conflicts with accepted decisions;
- whether the proposal is additive, corrective, or obsolete.

Do not silently rewrite accepted Product Knowledge because a plan or repository
appears to disagree with it.

## 6. Creating plans

Before drafting, identify:

- desired outcome;
- repositories involved;
- Product Knowledge pages that govern the work;
- source evidence that must be inspected;
- task boundaries;
- acceptance and verification layers;
- risks, dependencies, and stop conditions.

Perform a context-grounding check before finalizing the draft. Compare the
proposed objective, repositories, scope, decisions, and constraints with the
selected Product Knowledge. Surface contradictions or missing context instead
of writing a plausible plan that silently drifts from the project.

Perform a request-fidelity check as well. Walk through the original request and
make sure the plan captures its explicit requirements, desired behavior,
constraints, non-goals, terminology, and expected outcomes. Details that cannot
yet be resolved become assumptions, open questions, or risks rather than being
silently omitted.

Create a plan specific enough for one worker to execute without the root agent
inventing repository scope. Do not copy a proposed implementation into a worker
prompt; the plan and linked files are the source of truth.

Distinguish known project facts from assumptions, requested outcome from
implementation choice, acceptance from verification, repository change from
delivery, and proposed context from accepted Product Knowledge.

## 7. Worker behavior

The worker receives one execution brief for all tasks. It should:

1. Read the plan and task map.
2. Read named Product Knowledge and repository instructions.
3. Confirm assigned worktrees and repository-to-task mapping.
4. Execute tasks in dependency order.
5. Run implementation checks as tasks become complete.
6. Commit each affected repository after implementation.
7. Write a handoff with changes, commits, tests, assumptions, and blockers.

The worker must not edit the anchor checkout, change plan approval or completion,
mark its own work verified, alter verifier evidence, expand scope, rewrite a
prior commit to conceal repair, or rewrite or accept Product Knowledge.

If the worker encounters undeclared scope, it stops and reports the smallest
plan change required.

## 8. Verifier behavior

The verifier receives the plan, repository map, latest worker commits, worker
handoff, canonical verification commands, evidence identifiers, and read-only
worktree access.

It must independently check the latest commits and distinguish:

- failed product assertion;
- missing or ambiguous evidence;
- blocked host capability;
- repository or path violation;
- passed result.

It writes only its own handoff and evidence. It never changes product files,
Product Knowledge, plan status, or worker commits.

## 9. Repair and completion behavior

A failed verifier result is sent to the worker as structured repair evidence.
The worker repairs only the reported scope or a directly necessary dependent
change and creates new commits. The verifier then checks the latest commits.

The root response shows failure count, failed IDs, repair scope, new commits, and
whether another attempt remains. After three worker failures, execution stops.

When the verifier passes, report verified execution and show commits and
evidence. Do not say the plan is complete unless the human explicitly requests
the status change.

When the human requests completion, confirm the latest verified execution and
record the plan revision, affected repositories, commits, verifier evidence,
and completion request. Then compare the implemented changes with the context
used to create the plan. Report whether no durable context change was found or
whether context proposals, stale references, or conflicts need review. Do not
accept the proposals in the same implicit action.

## 10. Prompt contract

Every role prompt should be generated from current state and include:

- role and objective;
- plan and task locators;
- repository and worktree map;
- allowed paths;
- Product Knowledge references;
- acceptance and verification identifiers;
- commit and handoff requirements;
- stop conditions;
- current repair attempt, if any.

It should not include the full engine implementation, unrelated plans or
sessions, repeated lifecycle prose, copied solution text, credentials, or
provider payloads.

## 11. Root response contract

A normal root response contains:

1. action understood;
2. what was read or changed;
3. current result;
4. relevant commits or evidence;
5. the next human decision, only when required.

Do not expose packet math, graph construction, lease renewal, or host adapter
mechanics unless the user asks for diagnostics.
