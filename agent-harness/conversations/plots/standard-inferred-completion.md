# Standard completion is inferred — acceptance plus delivery, no second sign-off

Proves inferred completion at Standard: the checked result is accepted and its pull
request is opened, so the change is finished from the human's side without a separate
"mark it done" step. The coordinator records the delivery only because the human
supplied it, infers completion, and surfaces the reconciliation follow-up rather than
changing knowledge silently.

## Spec
```yaml
id: standard-inferred-completion
title: Infer Standard completion only after current acceptance and explicit Gate 2 delivery
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-deliver, cc-complete]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-standard-export
      title: Add a Markdown export command
      repository: notes
      objective: Add an export command for the notes.
      path: src/export
      seed_state: verified-accepted        # current candidate verified AND accepted
  state: seeded:verified-accepted-standard
persona: >
  A solo maker who reviewed the checked result and opened its pull request. They want
  the coordinator to record that delivery and explain what happens next. Has never
  heard of Context Circuit or its internals.
human_turns:
  - "I reviewed and accept the checked result, and I opened the pull request. Please record that delivery and finish the Standard change."
  - "What follow-up remains after it is finished?"
  - "Okay, thanks."
reactions:
  approves: true                   # supplies the explicit Gate 2 delivery event
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-14, AC-16, AC-22, AC-35]
  invariants: [INV-CANDIDATE-01, INV-DELIVER-01, INV-COMPLETE-01, INV-COMPLETE-02, INV-KNOWLEDGE-02]
hidden:
  - candidate ids, execution branches, worktrees, runtime files, digests
  - the tier label "Standard" and the words "inferred completion", "Gate 2" as jargon
decision_points:
  - id: record-delivery-only-on-event
    when: the human says they opened the pull request
    coordinator_must: >
      record delivery only because the human explicitly supplied that event — it does
      not open or infer a delivery on its own.
  - id: infer-not-ask-again
    when: acceptance and delivery both exist
    coordinator_must: >
      infer completion from the current acceptance plus the delivery, and NOT ask for
      a separate plan-completion approval.
  - id: surface-reconcile-followup
    when: the human asks what follow-up remains
    coordinator_must: >
      report the knowledge-reconciliation follow-up in plain terms without silently
      changing what the project knows.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-standard-export
  - candidate_bound: 0001-standard-export
  - human_acceptance_recorded: 0001-standard-export
  - delivery_recorded: 0001-standard-export
  - completion_kind: 0001-standard-export:inferred      # inferred, not a second approval
  - completion_recorded: 0001-standard-export
  - plan_status: 0001-standard-export:done
  - knowledge_debt_pending: ">=1"
  - product_knowledge_unchanged_silently: true
access_discipline:
  complete:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** I reviewed and accept the checked result, and I opened the pull request.
> Please record that delivery and finish the Standard change.
>
> **C:** Recorded — you've accepted this exact result and opened the pull request.
> That's the change accepted and delivered, so it's finished from your side; I didn't
> need a separate "mark it done" step. `[decision_point: record-delivery-only-on-event]` `[decision_point: infer-not-ask-again]`
>
> **H:** What follow-up remains after it is finished?
>
> **C:** One thing: this delivered work hasn't been folded into what the project knows
> yet. I won't change that on my own — say the word and I'll put together a proposed
> update, or you can note that nothing needs changing. `[decision_point: surface-reconcile-followup]`
>
> **H:** Okay, thanks.

Completion is inferred from the recorded acceptance plus the delivery the human
supplied — no second sign-off — and the reconciliation follow-up is surfaced rather
than applied silently.
