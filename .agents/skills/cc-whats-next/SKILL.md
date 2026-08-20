---
name: cc-whats-next
description: Recommend the smallest safe next action from plans, sessions, dependencies, leases, and human gates.
---

# Choose the next action

Use this skill during root-session entry, resume, or when the user asks what
should happen next.

Inspect only the relevant runtime records, approved plans, dependencies,
Product Knowledge warnings, repository state, worktrees, and pending human
gates. Inspect optional `archive.yaml` sidecars and omit a plan whose latest
event is `archived`; keep its original path available only for explicit
historical inspection. A malformed archive record is a blocker. Do not infer a current session from conversation history or a global
pointer.

Recommend one smallest safe action at a time. These names are discoverability
aids, not mandatory ceremonies for every session:

- gather context or draft a plan;
- review a draft, or approve a coherent draft through `cc-approve-plan`;
- execute an approved dependency-ready plan through `cc-run-plan`;
- execute connected approved unimplemented plans, or resume an interrupted
  `.runtime/stacks/<stack-id>/` run, through `cc-run-stack`;
- resume a session from its latest handoff;
- verify or repair within approved scope;
- finish a plan through `cc-finish-plan` when completion evidence is ready;
- clean local execution state through `cc-cleanup-runtime` after a plan is
  done or when the user asks to clear `.runtime/`, still mentioning dirty or
  unpushed work when present;
- request a human decision for a contradiction, ownership conflict, scope
  change, publication, merge, or deployment.
- archive or restore one named plan through `cc-archive-plan` only after its
  archive preflight and explicit human `archive` confirmation.

The recommendation is read-only until the root session or human explicitly
performs a consequential action. Report evidence, assumptions, blockers, and
the next safe action separately. For an approved dependency-ready plan,
recommend `cc-run-plan` so the root creates a writer child and a verifier
child rather than skipping children for small work. For connected approved
unimplemented plans, or a resumable stack run, recommend `cc-run-stack`.
Refuse to treat a stack run as one `cc-run-plan`. A one-plan request still
enters `cc-run-plan`.
