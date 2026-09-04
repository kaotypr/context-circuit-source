# Archive an intent by name — the same status-blind organization, one level up

An optional plot completing the Organization section: intents archive and restore the
same status-blind way plans do (part of AC-30). Here the human sets an approved intent
aside; the coordinator moves it off the active list without validating or changing its
approved state, and without touching any plan, execution, or delivery. Restore is the
symmetric inverse (seed archived, end active), exactly as `restore-archived-plan` is to
`archive-plan`.

## Spec
```yaml
id: archive-restore-intent
title: Archive an approved intent by name — moves it aside, status-blind, nothing else touched
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-intent]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  intents:
    - id: i001-md-export
      title: Add a Markdown export command
      repository: notes
      objective: Add a command that exports the notes to a chosen output file.
      tier: standard
      state: approved                    # approved -> archiving must not change or validate that
  state: seeded:one-approved-intent
persona: >
  A solo maker who approved a change earlier but wants to shelve the whole idea for now.
  Has never heard of Context Circuit; does not know what an intent is — just "that
  export idea I signed off on".
human_turns:
  - "Shelve that export idea I approved earlier — take it off my active list for now."
  - "Thanks."
reactions:
  approves: false                  # already approved earlier; this is organization, not a new decision
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-30]     # intents are first-class; archive/restore are status-blind moves
  invariants: [INV-ARCHIVE-01, INV-ARCHIVE-02]
hidden:
  - runtime files, the archived-area path, intent id and contract file names
  - lifecycle jargon ("deindex", "contract digest", "status projection") as jargon
decision_points:
  - id: archive-intent-on-request
    when: the human asks to shelve the approved idea
    coordinator_must: >
      move the named idea off the active list and confirm in plain terms, treating an
      intent's organization exactly like a plan's.
  - id: approval-untouched
    when: it archives
    coordinator_must: >
      not change or re-validate the fact that it was approved, and not require anything
      of it first — archiving is status-blind for intents too.
  - id: nothing-else-touched
    when: it confirms
    coordinator_must: >
      not start, plan, complete, deliver, or delete anything — only the idea's place on
      the active list changes.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - intent_archived: i001-md-export:approved   # archived, deindexed, approved-state preserved (NEW predicate, harness gap)
  - plans_created: 0
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  archive:
    required: []
    allowed: ["intent/**", "plans/**"]
    forbidden:
      - "intent/archive/**"          # never READ the archived area (the runtime MOVES the intent there)
      - "plans/archive/**"
      - "wrapper/runtime/engine.sh"
      - "sources/**"
```

## Dialogue

> **H:** Shelve that export idea I approved earlier — take it off my active list for
> now.
>
> **C:** Done — the export idea is set aside, off your active list. The fact that you'd
> approved it is untouched; nothing's been started, changed, or removed — it's just
> shelved, and you can bring it back whenever. `[decision_point: archive-intent-on-request]` `[decision_point: approval-untouched]` `[decision_point: nothing-else-touched]`
>
> **H:** Thanks.

An approved intent is shelved the same status-blind way a plan is — its approved state
preserved, nothing else touched — confirming intents are first-class and archive/restore
apply to them identically.
