# Two changes that don't combine cleanly — stop and say so honestly

Proves the honest-block branch of the change-set path (the tail of design conversation
C): two individually-checked changes are asked to ship together, but their combination
doesn't build. The coordinator reports plainly that they don't combine cleanly and it
stopped — nothing merged, nothing completed — and offers to split or reorder. A base
that can't be built is a block for the human to resolve, never an invented worker
failure and never a silent partial ship.

## Spec
```yaml
id: change-set-base-unbuildable
title: A change set whose combined base will not build is reported blocked, not completed
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
fault: change-set-base-unbuildable
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
    # Two individually-verified members whose combination does not build (they touch
    # the same surface in incompatible ways). The `fault` forces BASE_UNBUILDABLE at
    # change-set-prepare — a NEW fault, mirroring case 09's verifier-unavailable fault
    # (see coverage/phase-5.md).
    - { id: 0001-export-core, title: Export core, repository: notes, objective: Add the core export implementation., path: src/export, seed_state: verified }
    - { id: 0002-export-rewrite, title: Export rewrite, repository: notes, objective: Rework the same export surface a second, incompatible way., path: src/export, seed_state: verified }
  state: seeded:two-verified-members-incompatible
persona: >
  A maker with two checked changes they assume ship together. They're fine to split or
  reorder once told the combination doesn't work. Has never heard of Context Circuit or
  its internals.
human_turns:
  - "Both of those are checked — combine them and ship them as one, and complete them."
  - "Ah — okay, I'll split them then. Good that you didn't force it."
  - "Thanks."
reactions:
  approves: false                  # the block prevents the combined completion they asked for
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  # The change-set family (AC-35) with the negative branch owned by INV-CONCURRENCY-02
  # ("a base that cannot be built cleanly is a blocked execution, not a worker failure")
  # and INV-DELIVER-01/INV-PRESERVE-01. See coverage/phase-5.md.
  acceptance_criteria: [AC-35]
  invariants: [INV-DELIVER-01, INV-CONCURRENCY-02, INV-PRESERVE-01]
hidden:
  - candidate ids, integration-tip internals, execution branches, worktrees, runtime files
  - the words "change set", "BASE_UNBUILDABLE", "integration tip" as jargon
decision_points:
  - id: report-doesnt-combine
    when: the combined result will not build
    coordinator_must: >
      say plainly the two changes don't combine cleanly and it stopped rather than
      forcing them together — honest about the block, not dressed up as a code failure
      of either change on its own.
  - id: nothing-merged-or-completed
    when: it reports the block
    coordinator_must: >
      not merge, partially ship, or complete either change; both remain as they were,
      individually checked and preserved.
  - id: offer-split-or-reorder
    when: the human decides what to do
    coordinator_must: >
      offer a real path — split them into separate deliveries, or reorder so they
      combine — leaving the choice to the human.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-export-core          # each member is individually verified...
  - execution_verified: 0002-export-rewrite
  - change_set_base_unbuildable: true             # ...but the combination will not build (NEW predicate, harness gap)
  - change_set_blocked_not_worker_failure: true   # a block, never an invented worker failure (INV-CONCURRENCY-02)
  - plan_status: 0001-export-core:draft           # neither completed
  - plan_status: 0002-export-rewrite:draft
  - no_change_set_completed: true                 # nothing merged or completed
  - members_preserved: true                       # both members preserved as-is (INV-PRESERVE-01)
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
> **C:** I tried to prepare the two as one combined change, and they don't combine
> cleanly — put together, the result doesn't build. Each one is fine on its own; it's
> the combination that doesn't work. So I stopped rather than force something together
> or ship half of it. `[decision_point: report-doesnt-combine]`
>
> Nothing's merged and neither is marked complete — they're both exactly as they were,
> still individually checked. `[decision_point: nothing-merged-or-completed]`
>
> Two ways forward: ship them as two separate changes instead, or rework one so they
> fit together and then combine. Your call. `[decision_point: offer-split-or-reorder]`
>
> **H:** Ah — okay, I'll split them then. Good that you didn't force it.
>
> **C:** Sounds right — say the word and I'll take them one at a time.
>
> **H:** Thanks.

The combined base won't build, so the coordinator reports an honest block — not a
worker failure and not a silent partial ship — merges and completes nothing, preserves
both members, and offers to split or reorder.
