# Gather context writes live files — no extra knowledge-acceptance gate

Replaces the removed pending-proposal consent conversation. Gathering context
edits live context files and keeps the catalog consistent. There is no sidecar
and no third gate.

## Spec
```yaml
id: accept-or-defer-context-proposal
title: Gather context by writing live Product Knowledge files
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-workspace, context]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  state: fresh
persona: >
  A solo maker who asks the assistant to learn the project. Has never heard of
  Context Circuit.
human_turns:
  - "Please get up to speed on this project and write down what it knows."
  - "Yes, notes are JSON files on disk, one per note. Save that."
  - "That's fine for now."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-34]
  invariants: [INV-KNOWLEDGE-02, INV-COMPLETE-02]
hidden:
  - internal file names/paths
  - the words "context proposal", "knowledge unit", "reconciliation debt", "candidate" as jargon
decision_points:
  - id: gather-in-place
    when: the human asks to get up to speed and write it down
    coordinator_must: >
      read the project and write live context files, keeping the catalog
      consistent, rather than staging a sidecar for later acceptance.
  - id: no-third-gate
    when: knowledge is written
    coordinator_must: >
      not ask for a separate accept-the-update decision after gathering.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plans_created: 0
  - no_execution_records: true
access_discipline:
  learn:
    required: [workspace.yaml, context/INDEX.md]
    allowed: [AGENTS.md, WORKFLOW.md, "context/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** Please get up to speed on this project and write down what it knows.
>
> **C:** I'll read the project and update what it knows in place.
> `[decision_point: gather-in-place]`
>
> **H:** Yes, notes are JSON files on disk, one per note. Save that.
>
> **C:** Saved: notes are JSON files on disk, one per note. That's now part of
> what the project knows. `[decision_point: no-third-gate]`
>
> **H:** That's fine for now.

Gathering writes live files. There is no extra knowledge-acceptance gate.
