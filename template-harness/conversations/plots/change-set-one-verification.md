# Two changes, one combined delivery — one check, one acceptance, both complete

Proves the change-set path: two verified changes that belong together are delivered as
one. The coordinator prepares one integrated result, checks the combination once,
records one acceptance, and after the single pull request completes both members
together — one deterministic candidate, verified and accepted once, not two separate
deliveries.

## Spec
```yaml
id: change-set-one-verification
title: Deliver two verified plan members as one change set with one integration check and one acceptance
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-deliver, cc-verify, cc-complete]
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
    - { id: 0002-export-cli, title: Export command line interface, repository: notes, objective: Add the command line interface for the export implementation., path: src/export-cli, seed_state: verified }
  state: seeded:two-verified-members
persona: >
  A maker with two checked changes that belong in one release. They want one combined
  delivery and one clear acceptance decision, not two separate pull requests. Has never
  heard of Context Circuit or its internals.
human_turns:
  - "These two checked changes belong together. Prepare them as one combined delivery, check the combined result once, and I accept it."
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
      separate deliveries.
  - id: check-and-accept-once
    when: the combined result is prepared
    coordinator_must: >
      run one independent check against the combined result and record one human
      acceptance of it — not one per change.
  - id: complete-both-on-one-delivery
    when: the human opens the single pull request
    coordinator_must: >
      after the one explicit delivery event, complete both changes consistently and
      report any knowledge follow-up.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-export-core
  - execution_verified: 0002-export-cli
  - change_set_completed: 2                     # both members completed from one acceptance/delivery
  - plan_status: 0001-export-core:done
  - plan_status: 0002-export-cli:done
  - knowledge_debt_pending: ">=1"
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
  deliver:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** These two checked changes belong together. Prepare them as one combined
> delivery, check the combined result once, and I accept it.
>
> **C:** Prepared them as one combined change and checked the two together as they'll
> actually ship — that combined result passes its independent check. I've recorded
> your acceptance of the combined result, once. `[decision_point: one-combined-delivery]` `[decision_point: check-and-accept-once]`
>
> **H:** I opened the one pull request for the combined result. Complete both changes
> together and tell me what follow-up remains.
>
> **C:** Both changes are completed together off that one delivery. The one follow-up:
> this delivered work hasn't been folded into what the project knows yet — I'll draft a
> proposed update for you when you want, and won't change anything on my own. `[decision_point: complete-both-on-one-delivery]`
>
> **H:** Thanks.

Two related changes are delivered as one integrated result, checked and accepted once,
and both complete together from the single delivery — one candidate, not two separate
deliveries.
