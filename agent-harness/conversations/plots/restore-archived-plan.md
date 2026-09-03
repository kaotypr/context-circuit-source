# Restore an archived plan by name — bring it back active, nothing lost

The inverse of `archive-plan`, split out because an archive→restore round-trip in one
conversation is undetectable from final state (active == never-archived). Here the plan
starts archived and must end active: the coordinator finds it by name, brings it back to
the active list with its status preserved, and validates nothing.

## Spec
```yaml
id: restore-archived-plan
title: Restore an archived plan by name — returns it active, status preserved, nothing lost
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-archive]
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
      objective: Add an export.py command that exports the notes to a chosen output file.
      repository: notes
      seed_state: archived               # STARTS archived -> the round-trip is detectable from final state
  state: seeded:one-archived-plan
persona: >
  A solo maker who earlier set a plan aside and now wants it back. Remembers it was
  about the export command. Has never heard of Context Circuit or its internals.
human_turns:
  - "Can you bring back that export plan I set aside earlier? I want it active again."
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-24]
  invariants: [INV-ARCHIVE-01, INV-ARCHIVE-02]
hidden:
  - runtime files, the archived-area path, plan file names
  - lifecycle jargon ("reindex", "status projection") as jargon
decision_points:
  - id: find-and-restore
    when: the human asks to bring back the set-aside plan
    coordinator_must: >
      find the set-aside plan by name and return it to the active list, confirming in
      plain terms.
  - id: status-preserved
    when: it restores
    coordinator_must: >
      keep the plan's status exactly as it was and lose nothing — restore validates or
      changes nothing.
  - id: no-overclaim
    when: it confirms
    coordinator_must: >
      not claim to have approved, completed, delivered, or run anything — it's simply
      back on the active list.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plan_status: 0001-md-export:draft     # restored active with status preserved
  - plan_not_archived: 0001-md-export     # moved back out of the archive
  - plan_indexed: 0001-md-export          # active index row restored
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  restore:
    # Restore is the ONE operation that may act on the archived plan; it MOVES it via
    # the runtime rather than READing it, so reading the archive is still forbidden.
    required: []
    allowed: ["plans/**"]
    forbidden:
      - "wrapper/runtime/engine.sh"
      - "sources/**"
```

## Dialogue

> **H:** Can you bring back that export plan I set aside earlier? I want it active
> again.
>
> **C:** Done — the export plan is back on your active list, just as it was when you
> set it aside. Nothing was lost and nothing changed about it; it's simply active
> again. `[decision_point: find-and-restore]` `[decision_point: status-preserved]` `[decision_point: no-overclaim]`
>
> **H:** Thanks.

The set-aside plan is found by name and returned to the active list with its status
preserved and nothing lost — no approval, completion, or execution implied.
