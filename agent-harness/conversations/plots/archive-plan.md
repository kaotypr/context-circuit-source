# Archive a plan by name — move it aside, no status check, nothing lost

Proves status-blind organization: the human sets a plan aside, and the coordinator
moves it off the active list without validating or changing its status, without
requiring it to be complete, and without deleting anything. Archiving is pure
organization — no completion, no delivery, no cleanup implied.

## Spec
```yaml
id: archive-plan
title: Archive a plan by name — moves it aside, deindexes it, preserves status, no validation
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
      seed_state: draft
  state: seeded:one-draft-plan
persona: >
  A solo maker who wants a plan moved off the active list for now. Has never heard of
  Context Circuit and does not know its internals.
human_turns:
  - "Archive the export plan for now — move it off my active list; I'm not ready to work on it."
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
  - lifecycle jargon ("deindex", "status projection") as jargon
decision_points:
  - id: archive-on-request
    when: the human asks to set the named plan aside
    coordinator_must: >
      archive the named plan and confirm in plain terms that it's off the active list.
  - id: no-status-fuss
    when: it archives
    coordinator_must: >
      not warn about or change the plan's status, and not require it to be complete
      first — archiving is status-blind.
  - id: nothing-destroyed
    when: it confirms
    coordinator_must: >
      not claim to have deleted, completed, delivered, or run anything — the plan is
      preserved, just moved aside.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plan_archived: 0001-md-export:draft     # in the archive, deindexed, status preserved (no validation)
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  archive:
    required: []
    allowed: ["plans/**"]
    forbidden:
      - "plans/archive/**"            # never READ the archived area (the runtime MOVES the plan there)
      - ".context-circuit/wrapper/runtime/engine.sh"
      - "sources/**"
```

## Dialogue

> **H:** Archive the export plan for now — move it off my active list; I'm not ready
> to work on it.
>
> **C:** Done — I've set the export plan aside, so it's off your active list. It's not
> deleted or changed in any way; it's just tucked away, exactly as it was, and you can
> bring it back whenever you like. `[decision_point: archive-on-request]` `[decision_point: no-status-fuss]` `[decision_point: nothing-destroyed]`
>
> **H:** Thanks.

The named plan is moved off the active list without any status check or change, and
nothing is deleted, completed, or run — organization only.
