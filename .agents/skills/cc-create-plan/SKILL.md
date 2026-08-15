---
name: cc-create-plan
description: Draft a context-grounded, human-reviewed plan from accepted intent and selected evidence.
---

# Create a plan

Use this skill when the user asks to turn an accepted Idea Brief, PRD, direct
request, or selected evidence into intended implementation work.

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

Offer `cc-review-plan` to identify decisions and missing evidence, or ask for
explicit approval when the plan is ready. Do not start execution from this
skill; approved work enters through `cc-run-plan`.
