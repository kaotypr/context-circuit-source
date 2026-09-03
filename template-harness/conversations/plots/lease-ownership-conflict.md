# Someone's already working that area — read-only, never a silent steal

A conversational probe over a deterministically-owned mechanic (the path lease /
one-worker lock, owned by the engine tests): a second request tries to touch an area
another execution already holds. This plot proves how that is *explained to a lay
user* — the coordinator says the area is currently in progress, won't step on it,
and offers to wait or work elsewhere. A live claim is never silently taken over.

## Spec
```yaml
id: lease-ownership-conflict
title: A second request to touch in-progress work is held read-only, never steals the lease
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-execute, cc-run-stack]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-md-export
      title: Add a Markdown export command
      repository: notes
      objective: Add an export command under the export area.
      path: src/export
      seed_state: executing          # actively executing, holding a lease on notes:src/export
  # A seeded in-progress execution holding a live lease on notes:src/export. Seeding an
  # active lease is a NEW capability (harness gap) — see coverage/phase-4.md.
  state: seeded:execution-holding-lease
persona: >
  A solo maker who forgot a change is already being built and asks to start another
  one that touches the same area. Has never heard of Context Circuit; does not know
  what a lease or a lock is. Fine to wait once it's explained.
human_turns:
  - "While that's going, can you also start building another tweak in the export area right now?"
  - "Ah, right — okay, wait until the current one's done then. Don't force it."
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  # No dedicated acceptance criterion: the path lease / one-worker ownership is owned
  # at the invariant level (INV-CONCURRENCY-01 / INV-OWN-01) and by the engine tests;
  # the v0.5-rooted criteria map has no lease AC. See coverage/phase-4.md.
  acceptance_criteria: []
  invariants: [INV-CONCURRENCY-01, INV-OWN-01, INV-PRESERVE-01]
hidden:
  - runtime file names, execution branches, worktrees, the words "lease"/"lock" as jargon
  - internal mechanism ("path region", "exclusive-create") as jargon
decision_points:
  - id: recognize-in-progress
    when: the human asks to start work on an area already being built
    coordinator_must: >
      recognize that the area is currently being worked by an active run and say so in
      plain language — that work is in progress right now.
  - id: refuse-to-step-on-it
    when: it has recognized the conflict
    coordinator_must: >
      decline to start a second change that would touch the same area — read-only /
      held — and explicitly not take over or interrupt the running work.
  - id: offer-wait-or-elsewhere
    when: it declines
    coordinator_must: >
      offer a real path: wait until the current work is done, or work on a different
      area that isn't in progress — the choice is the human's.
  - id: preserve-running-work
    when: the human says to wait
    coordinator_must: >
      leave the running work untouched and start nothing — the live work is never
      disturbed.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - second_execution_blocked_read_only: notes:src/export   # NEW predicate: the conflicting request was held (harness gap)
  - lease_not_stolen: 0001-md-export                        # the live lease is intact, never taken over
  - running_work_preserved: 0001-md-export                  # the in-progress execution is undisturbed (INV-PRESERVE-01)
  - product_knowledge_unchanged_silently: true
access_discipline:
  handle-conflict:
    required: []
    allowed: ["plans/**", ".runtime/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** While that's going, can you also start building another tweak in the export
> area right now?
>
> **C:** That same area is being worked right now by the build already in progress,
> so I won't start a second change on top of it — I don't want to step on work that's
> underway or risk clobbering it. `[decision_point: recognize-in-progress]` `[decision_point: refuse-to-step-on-it]`
>
> Two options: wait until the current one finishes and then I'll do this next, or I
> can pick up something in a different area that isn't in progress. Which do you
> prefer? `[decision_point: offer-wait-or-elsewhere]`
>
> **H:** Ah, right — okay, wait until the current one's done then. Don't force it.
>
> **C:** Will do — I've left the running work exactly as it is and started nothing
> new. Once it's finished, I'll pick this up. `[decision_point: preserve-running-work]`
>
> **H:** Thanks.

The coordinator recognizes the area is under an active claim, holds the second
request read-only rather than stealing the lease, offers to wait or work elsewhere,
and leaves the running work undisturbed.
