# Plan review

Plan review is a focused readiness check, not an automatic approval step. The
reviewer gathers evidence from the plan, selected Product Knowledge, source
provenance, repository instructions, and current implementation state, then
separates routine evidence from decisions that belong to the human.

## Readiness checklist

Review the following before execution:

- objective, source, repository scope, implementation scope, and non-goals;
- relevant Product Knowledge and explicitly selected raw sources;
- task boundaries, dependencies, stop conditions, and expected handoffs;
- acceptance criteria, test scope, verification commands, and completion gate;
- assumptions, unresolved questions, contradictions, overlapping paths, and
  ownership or worktree risks;
- delivery boundaries, including what is deliberately outside the plan.

For a newly generated bundle, confirm that `plan.yaml`, one concise
`overview.md`, and validated `tasks/*.md` are sufficient. Check that every
specialist companion is justified by an explicit complexity signal in the
overview and points back to canonical plan fields instead of copying
acceptance, verification, or lifecycle values. The absence of optional
companions is not a defect in a simple plan.

For an existing bundle, discover companions rather than requiring the compact
layout. Historical `requirements.md`, `solution.md`, `risks.md`,
`delivery.md`, `acceptance.md`, `verification.md`, and other existing
companions remain readable without a migration or rewrite. Their content is
review rationale; `plan.yaml` remains authoritative.

Before reporting a newly generated task as ready, apply the host-provided
deterministic validation contract to its initial front matter: delimiters,
YAML, required fields, enum values, nested shapes, and the
`context-circuit.task` schema. A failed check is a readiness defect to repair,
not a plan-lifecycle result. Passing validation does not approve, execute,
verify, or complete the plan.

The review should identify the smallest next decision. It should not ask for a
separate approval of every task when the whole approved plan is ready to run.

## Review outcomes

Use one of these outcomes:

- `ready-for-approval`: the plan is coherent but still requires human plan
  approval; the next action is `cc-approve-plan`;
- `needs-revision`: scope, evidence, dependencies, or acceptance is incomplete;
- `blocked`: a contradiction or ownership decision prevents safe progress;
- `approved-for-execution`: the observation that the plan is already approved
  and may enter `cc-run-plan`, or that a connected approved set may enter
  `cc-run-stack`.

The review never writes `plan.yaml`. It records evidence and recommendations;
`cc-approve-plan` is the named plan-approval skill. The human controls
approval, material scope changes, and completion. Connected approved plans
may later enter `cc-run-stack`, which freezes a runtime `graph.yaml` and
resumes from `progress.yaml`. That is not a scheduler and not a second
approval of a stack artifact.

## Evidence format

A concise review records:

```text
Plan: <plan path>
Outcome: <review outcome>
Evidence: <context, repository, and plan files inspected>
Routine checks: <checks that need no human decision>
Decisions: <questions that change scope, intent, ownership, or delivery>
Contradictions: <source or artifact conflicts, or none>
Next action: <revise, cc-approve-plan, cc-run-plan, or cc-run-stack>
```

If a source or accepted Product Knowledge page changes during execution, stop
and surface a context-refresh decision. Do not silently rewrite the plan or
accepted context.
