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

The review should identify the smallest next decision. It should not ask for a
separate approval of every task when the whole approved plan is ready to run.

## Review outcomes

Use one of these outcomes:

- `ready-for-approval`: the plan is coherent but still requires human plan
  approval; the next action is `cc-approve-plan`;
- `needs-revision`: scope, evidence, dependencies, or acceptance is incomplete;
- `blocked`: a contradiction or ownership decision prevents safe progress;
- `approved-for-execution`: the observation that the plan is already approved
  and may enter `cc-run-plan`.

The review never writes `plan.yaml`. It records evidence and recommendations;
`cc-approve-plan` is the named plan-approval skill. The human controls
approval, material scope changes, and completion.

## Evidence format

A concise review records:

```text
Plan: <plan path>
Outcome: <review outcome>
Evidence: <context, repository, and plan files inspected>
Routine checks: <checks that need no human decision>
Decisions: <questions that change scope, intent, ownership, or delivery>
Contradictions: <source or artifact conflicts, or none>
Next action: <revise, cc-approve-plan, or cc-run-plan>
```

If a source or accepted Product Knowledge page changes during execution, stop
and surface a context-refresh decision. Do not silently rewrite the plan or
accepted context.
