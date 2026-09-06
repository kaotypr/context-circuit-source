# A later plan can start without waiting for a Product Knowledge update

Proves that unupdated Product Knowledge does not block the next overlapping
plan. A prior Standard plan is already marked done; the coordinator starts the
follow-up when asked and does not hold it for a knowledge-acceptance gate.

## Spec
```yaml
id: knowledge-debt-blocks-next-plan
title: Start overlapping Standard work even if prior Product Knowledge is not yet updated
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-plan, cc-execute, context]
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
  A maker whose storage foundation is finished and who immediately wants to start a
  related search feature. Has never heard of Context Circuit.
human_turns:
  - "The storage foundation is finished. Please start the follow-up search work now."
  - "If something is waiting on project notes, don't hold this — go ahead and plan it."
  - "Thanks."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-34]
  invariants: [INV-KNOWLEDGE-02, INV-COMPLETE-02, INV-EXEC-01]
hidden:
  - runtime files, candidate ids, execution branches, worktrees, digests
  - the words "reconciliation debt", "grounding preflight", "knowledge scope" as jargon
decision_points:
  - id: start-follow-up
    when: the human asks to start related follow-up work
    coordinator_must: >
      start planning or executing the follow-up rather than holding it until
      Product Knowledge from the prior plan is updated.
  - id: no-knowledge-gate
    when: prior Product Knowledge may still be unupdated
    coordinator_must: >
      not require a separate knowledge-acceptance decision before the next plan
      can start.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plan_status: 0001-foundation:done
  - knowledge_debt_pending: 0
  - plan_status: 0002-follow-up:draft
  - product_knowledge_unchanged_silently: true
access_discipline:
  orient:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
  execute-plan:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** The storage foundation is finished. Please start the follow-up search work
> now.
>
> **C:** Starting the search work. The storage change is already marked done; I'll
> plan the search on top of it. `[decision_point: start-follow-up]`
>
> **H:** If something is waiting on project notes, don't hold this — go ahead and
> plan it.
>
> **C:** Not holding it. Project notes can be updated later; that doesn't block this
> plan. `[decision_point: no-knowledge-gate]`
>
> **H:** Thanks.

A later plan can start even if Product Knowledge from a prior plan has not been
updated.
