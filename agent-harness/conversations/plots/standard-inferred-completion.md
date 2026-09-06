# Standard delivery records Gate 2 only — done requires an explicit mark-done

Proves that recording delivery does not finish a Standard change. The checked
result is accepted and its pull request is opened; the coordinator records that
delivery and does not mark the plan done, start a knowledge update, or infer
completion. Mark-done remains a separate explicit ask.

## Spec
```yaml
id: standard-inferred-completion
title: Record Gate 2 delivery without marking a Standard plan done
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
      seed_state: verified-accepted
  state: seeded:verified-accepted-standard
persona: >
  A solo maker who reviewed the checked result and opened its pull request. They want
  the coordinator to record that delivery. Has never heard of Context Circuit or its internals.
human_turns:
  - "I reviewed and accept the checked result, and I opened the pull request. Please record that delivery."
  - "Is the change done?"
  - "Okay, thanks."
reactions:
  approves: true
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
  - id: do-not-mark-done
    when: delivery has been recorded
    coordinator_must: >
      not mark the plan done and not infer completion from delivery. Done requires
      an explicit mark-done ask.
  - id: no-knowledge-update-from-delivery
    when: the human asks whether the change is done
    coordinator_must: >
      say the work is delivered but not marked done, and that recording delivery
      did not update what the project knows.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-standard-export
  - candidate_bound: 0001-standard-export
  - human_acceptance_recorded: 0001-standard-export
  - delivery_recorded: 0001-standard-export
  - plan_status: 0001-standard-export:draft
  - knowledge_debt_pending: 0
  - product_knowledge_unchanged_silently: true
access_discipline:
  complete:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** I reviewed and accept the checked result, and I opened the pull request.
> Please record that delivery.
>
> **C:** Recorded — you've accepted this exact result and opened the pull request.
> Nothing has been marked done. `[decision_point: record-delivery-only-on-event]` `[decision_point: do-not-mark-done]`
>
> **H:** Is the change done?
>
> **C:** The pull request is recorded. The change is not marked done until you ask
> me to mark it done. Recording delivery did not update what the project knows.
> `[decision_point: no-knowledge-update-from-delivery]`
>
> **H:** Okay, thanks.

Delivery records Gate 2 only. A plan becomes done only on an explicit mark-done.
