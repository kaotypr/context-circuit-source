# Mark-done already happened — the next plan can proceed

The complement of the inverted next-plan plot: a prior plan is marked done.
The coordinator starts the follow-up when asked. In-place knowledge updates,
when they happen, are the same gather-style edits — not a sidecar and not a
block on the next plan.

## Spec
```yaml
id: reconcile-and-proceed
title: Start the next plan after a prior plan was marked done
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [context, cc-plan]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-foundation
      title: Add notes storage
      repository: notes
      objective: Add a small storage foundation for the notes tool.
      path: src/storage
      seed_state: completed-standard
    - id: 0002-follow-up
      title: Add note search
      repository: notes
      objective: Add search on top of the notes storage foundation.
      path: src/search
      seed_state: draft
  state: seeded:completed-standard-follow-up-draft
persona: >
  A maker whose storage change is marked done and who wants the search work to
  go ahead. Has never heard of Context Circuit.
human_turns:
  - "The storage change is done. Please start the search work."
  - "If the project notes still need an update, don't wait — start the search plan."
  - "Thanks."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-34]
  invariants: [INV-KNOWLEDGE-02, INV-COMPLETE-02]
hidden:
  - runtime files, candidate ids, execution branches, digests
  - the words "reconciliation debt", "marker", "grounding preflight", "knowledge scope" as jargon
decision_points:
  - id: proceed
    when: the human asks to start the follow-up
    coordinator_must: >
      start the follow-up plan rather than waiting for a Product Knowledge update.
  - id: notes-can-wait
    when: project notes may still be unupdated
    coordinator_must: >
      not treat unupdated notes as a gate on the next plan.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plan_status: 0001-foundation:done
  - knowledge_debt_pending: 0
  - plan_status: 0002-follow-up:draft
  - no_execution_for: 0002-follow-up
access_discipline:
  reconcile:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", ".context-circuit/wrapper/contracts/**"]
```

## Dialogue

> **H:** The storage change is done. Please start the search work.
>
> **C:** Starting the search plan. `[decision_point: proceed]`
>
> **H:** If the project notes still need an update, don't wait — start the search
> plan.
>
> **C:** Not waiting. Notes can be updated in place later; that doesn't block this.
> `[decision_point: notes-can-wait]`
>
> **H:** Thanks.

A later plan can start after mark-done without waiting for a knowledge update.
