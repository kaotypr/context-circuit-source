# Two changes, one combined delivery — existing checks, one acceptance, both complete

Proves the change-set path: two verified stacked changes that belong together are
delivered as one pull request from the covering execution branch. The coordinator
does not spawn a second verifier. It records one acceptance, and after the single
pull request completes both members together.

## Spec
```yaml
id: change-set-one-verification
title: Deliver two verified plan members as one change set with one acceptance and no delivery-time check
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
    - { id: 0001-export-core, title: Export core, repository: notes, objective: Add the core export implementation., path: src/export-core, seed_state: verified }
    - { id: 0002-export-cli, title: Export command line interface, repository: notes, objective: Add the command line interface for the export implementation., path: src/export-cli, deps: [0001-export-core], seed_state: verified }
  state: seeded:two-verified-members
persona: >
  A maker with two checked changes that belong in one release. They want one combined
  delivery and one clear acceptance decision, not two separate pull requests. Has never
  heard of Context Circuit or its internals.
human_turns:
  - "These two checked changes belong together. Deliver them as one pull request, and I accept it."
  - "I opened the one pull request for the combined result. Complete both changes together and tell me what follow-up remains."
  - "Thanks."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-14, AC-16, AC-22, AC-35]
  invariants: [INV-CANDIDATE-01, INV-VERIFY-01, INV-DELIVER-01, INV-COMPLETE-01, INV-COMPLETE-02]
hidden:
  - candidate ids, integration-tip internals, execution branches, worktrees, runtime files
  - the words "change set", "integration candidate", "integration tip" as jargon
decision_points:
  - id: one-combined-delivery
    when: the human asks to deliver two related changes together
    coordinator_must: >
      prepare one combined delivery of the two changes rather than opening two
      separate deliveries, using the covering branch that already contains both.
  - id: reuse-existing-checks
    when: the combined delivery is prepared
    coordinator_must: >
      reuse the existing independent checks and record one human acceptance —
      not spawn a second verifier, and not one acceptance per change.
  - id: deliver-without-mark-done
    when: the human opens the single pull request
    coordinator_must: >
      after the one explicit delivery event, record delivery for the covering
      tip without marking members done and without starting a knowledge update.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-export-core
  - execution_verified: 0002-export-cli
  - change_set_completed: 2                     # both members delivered from one acceptance; not marked done
  - plan_status: 0001-export-core:draft
  - plan_status: 0002-export-cli:draft
  - knowledge_debt_pending: 0
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh"]
  deliver:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** These two checked changes belong together. Deliver them as one pull
> request, and I accept it.
>
> **C:** They already sit on one branch that contains both, and each already has
> its independent check. I've recorded your acceptance of that combined result,
> once — no second check. `[decision_point: one-combined-delivery]` `[decision_point: reuse-existing-checks]`
>
> **H:** I opened the one pull request for the combined result. Complete both changes
> together and tell me what follow-up remains.
>
> **C:** Both changes are completed together off that one delivery. The one follow-up:
> this delivered work hasn't been folded into what the project knows yet — I'll draft a
> proposed update for you when you want, and won't change anything on my own. `[decision_point: complete-both-on-one-delivery]`
>
> **H:** Thanks.

Two related stacked changes are delivered as one pull request from the covering
branch, accepted once without a delivery-time verifier, and both complete together
from the single delivery.
