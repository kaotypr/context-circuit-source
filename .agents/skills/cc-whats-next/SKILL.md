---
name: cc-whats-next
description: Recommend the smallest safe next action from plans, sessions, dependencies, leases, and human gates.
---

# Choose the next action

Use this skill during root-session entry, resume, or when the user asks what
should happen next.

Inspect only the relevant runtime records, approved plans, dependencies,
Product Knowledge warnings, repository state, worktrees, and pending human
gates. Do not infer a current session from conversation history or a global
pointer.

Recommend one route:

- gather context or draft a plan;
- review or approve a draft plan;
- execute an approved dependency-ready plan through `cc-run-plan`;
- resume a session from its latest handoff;
- verify or repair within approved scope;
- request a human decision for a contradiction, ownership conflict, scope
  change, publication, merge, deployment, or completion gate.

The recommendation is read-only until the root session or human explicitly
performs a consequential action. Report evidence, assumptions, blockers, and
the next safe action separately. For an approved dependency-ready plan,
recommend `cc-run-plan` so the root creates a writer child and a verifier
child rather than skipping children for small work.
