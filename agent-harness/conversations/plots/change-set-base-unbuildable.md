# Two parallel changes that don't share one branch — stop and say so honestly

Proves the honest-block branch of the change-set path: two individually-checked
changes are asked to ship as one pull request, but neither execution branch
contains the other. Delivery does not merge them. The coordinator reports plainly
that they cannot ship as one and it stopped — nothing merged, nothing completed —
and offers to split or execute an integration.

## Spec
```yaml
id: change-set-base-unbuildable
title: A change set whose members do not share one covering branch is reported blocked, not completed
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
fault: change-set-no-single-tip
surface: [cc-deliver]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - { id: 0001-export-core, title: Export core, repository: notes, objective: Add the core export implementation., path: src/export, seed_state: verified }
    - { id: 0002-export-rewrite, title: Export rewrite, repository: notes, objective: Rework the same export surface a second, incompatible way., path: src/export, seed_state: verified }
  state: seeded:two-verified-members-incompatible
persona: >
  A maker with two checked changes they assume ship together. They're fine to split or
  reorder once told they don't sit on one branch. Has never heard of Context Circuit or
  its internals.
human_turns:
  - "Both of those are checked — combine them and ship them as one, and complete them."
  - "Ah — okay, I'll split them then. Good that you didn't force it."
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-35]
  invariants: [INV-DELIVER-01, INV-CONCURRENCY-02, INV-PRESERVE-01]
hidden:
  - candidate ids, integration-tip internals, execution branches, worktrees, runtime files
  - the words "change set", "CHANGE_SET_NO_SINGLE_TIP", "integration tip" as jargon
decision_points:
  - id: report-doesnt-combine
    when: the members do not share one covering branch
    coordinator_must: >
      say plainly the two changes don't sit on one branch and it stopped rather than
      merging them at delivery — honest about the block, not dressed up as a code failure
      of either change on its own.
  - id: nothing-merged-or-completed
    when: it reports the block
    coordinator_must: >
      not merge, partially ship, or complete either change; both remain as they were,
      individually checked and preserved.
  - id: offer-split-or-reorder
    when: the human decides what to do
    coordinator_must: >
      offer a real path — split them into separate deliveries, or execute an integration
      so they share one branch — leaving the choice to the human.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-export-core
  - execution_verified: 0002-export-rewrite
  - change_set_base_unbuildable: true
  - change_set_blocked_not_worker_failure: true
  - plan_status: 0001-export-core:draft
  - plan_status: 0002-export-rewrite:draft
  - no_change_set_completed: true
  - members_preserved: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  deliver:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** Both of those are checked — combine them and ship them as one, and complete
> them.
>
> **C:** They don't sit on one branch, so they can't ship as a single pull request.
> Each one is fine on its own. I stopped rather than merge them together or ship
> half of it. `[decision_point: report-doesnt-combine]`
>
> Nothing's merged and neither is marked complete — they're both exactly as they were,
> still individually checked. `[decision_point: nothing-merged-or-completed]`
>
> Two ways forward: ship them as two separate changes instead, or build an
> integration so they share one branch and then combine. Your call. `[decision_point: offer-split-or-reorder]`
>
> **H:** Ah — okay, I'll split them then. Good that you didn't force it.
>
> **C:** Sounds right — say the word and I'll take them one at a time.
>
> **H:** Thanks.

The members don't share one covering branch, so the coordinator reports an honest
block — not a worker failure and not a silent partial ship — merges and completes
nothing, preserves both members, and offers to split or integrate.
