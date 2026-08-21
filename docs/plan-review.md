# Review Card

Plan review is a read-only evidence check. Inspect the named plan bundle,
declared accepted context and provenance, repository instructions/evidence,
dependencies, archive state, active ownership, and delivery boundary.

Return:

```text
Plan: <id>
Outcome: ready-for-approval | needs-revision | blocked | approved-for-execution
Summary: <what changes and why>
Approval would authorize: <bounded implementation intent>
Approval would not authorize: execution, Git delivery, publication, deployment
Evidence inspected: <paths/revisions>
Tasks: <id, outcome, dependency>
Acceptance mapping: <criterion → task → verification evidence>
Risks/assumptions: <material items only>
Human decisions: <0–3 focused choices>
Contradictions/blockers: <explicit or none>
Next action: revise, approve, run, or resolve blocker
```

A review never changes plan status, task projections, leases, runtime, Git, or
external state.
