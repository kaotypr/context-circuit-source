---
name: cc-create-plan
description: Draft a context-grounded, human-reviewed plan from accepted intent and selected evidence.
---

# Create a plan

Use this skill when the user asks to turn an accepted Idea Brief, PRD, direct
request, or selected evidence into intended implementation work.

## Planning-only progressive reads

This sequence begins only after existing entry classification has selected
planning. It does not change entry classification or any non-planning route.
Start with the initial planning bundle: AGENTS.md, WORKFLOW.md, workspace.yaml,
context/INDEX.md, context/WORKSPACE.md, context/PROJECT.md, the direct request
or accepted Idea Brief/PRD, only explicitly selected source evidence, this
planning contract, current repository instructions and Git state, directly
relevant plan/dependency/archive metadata, and scoped session or ownership
records. The bundle must retain authority, intent, dependencies, relevant
Product Knowledge navigation, safety, and ownership.

Retrieve exact deeper evidence before drafting when intent is incomplete,
disputed, source-based, or changed; when a domain, role, workflow, decision, or
implementation assumption is unresolved; when the request names, replaces,
depends on, archives, resumes, or conflicts with another plan; when safe
drafting depends on live ownership, recovery, record shape, or worktree state;
or when a safety or lifecycle rule remains ambiguous or contradictory. Use the
relevant Idea Brief/PRD and selected provenance, linked Product Knowledge and
minimum repository evidence, directly relevant plan/task/archive/dependent/
handoff/completion evidence, or exact runtime records as the trigger requires.

If a required fact is missing, ambiguous, contradictory, or malformed after
targeted retrieval, use the existing complete-context fallback: load the full
workflow and runtime contracts and all relevant Product Knowledge, plan, task,
archive, session, handoff, lease, worktree, repository, and explicitly
selected-source evidence required by current behavior. Never infer a fact to
stay within the smaller bundle. Do not scan the source inbox as a side effect.

## Read and ground

Read the smallest relevant set of Product Knowledge, the accepted product
artifact if one exists, repository-local instructions, current implementation
evidence, and explicitly selected raw sources. Do not scan the source inbox as
a side effect. State which sources or repository files were read and why.

Separate user intent, observed implementation, accepted decisions,
assumptions, unknowns, proposals, and contradictions. Record provenance in the
plan without copying raw source text.

## Draft the artifact

Create one coherent plan under `plans/<repository-key>-plans/` with:

- canonical `plan.yaml` metadata and `status: draft`;
- objective, source, repository scope, Product Knowledge references,
  implementation scope, and non-goals;
- dependencies and connections to other plans or repositories;
- bounded tasks with explicit paths, dependencies, acceptance, verification,
  expected evidence, and stop conditions;
- acceptance criteria, test scope, verification commands, risks, assumptions,
  open questions, and human gates.

Use `plans/README.md` and the existing plan documents as the contract. Keep
the plan useful for a writer child: bounded tasks, explicit paths, acceptance,
and verification. Do not present root implementation as the small-work path.

## Gate

The result remains a draft until the user explicitly approves it. A complete
document, passing checks, or a positive review does not approve the plan.
Unresolved contradictions and material scope decisions stay visible.

## Next action

Offer `cc-review-plan` to identify decisions and missing evidence, or
`cc-approve-plan` when the plan is ready. Only `cc-approve-plan` writes
`draft` → `approved` after explicit human confirmation. Do not start
execution from this skill; approved work enters through `cc-run-plan`.
