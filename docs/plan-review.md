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

## Optional host question prompt

When the Human decisions line has one to three focused choices, the current
session may also offer those same choices through the current host's optional
native question-prompt primitive, after showing the card text:

| Host | Optional primitive |
| --- | --- |
| Cursor Agent | `AskQuestion` |
| Claude Code | `AskUserQuestion` |
| Codex CLI | `request_user_input` when listed |

The prompt is optional evidence of the current host, not a capability the
product requires. A missing, denied, or failed prompt falls back to the card
text; it is not `host-blocked` and must never be retried as if it were a
required writer or verifier child. Do not record the prompt, the chosen
option, or any transcript in `host_evidence` or any other durable record —
the choice stays in the conversation.

A next-action label such as `Revise this plan` or `Approve this plan`,
whether typed or chosen from the prompt, may continue only into that route's
existing first card in the current session (for example the Approval Card
owned by `cc-gates` and `docs/gates.md`). Choosing a label is never itself a
confirmation and never writes `plan.yaml` or task status.
