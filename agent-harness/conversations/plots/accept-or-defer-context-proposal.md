# Accept or defer a pending knowledge proposal — consent, per unit

Proves the consent gate on Product Knowledge in isolation: a proposal is already
waiting (left behind by earlier reconciliation), and the whole interaction is the
human deciding, unit by unit. The coordinator presents each pending unit in plain
language, records a unit only on an explicit yes, and leaves a deferred unit
exactly as it was — nothing is ever accepted silently.

Boundary with `reconcile-and-proceed` (a later plot): that plot proves debt
*blocking the next plan* and then clearing so work proceeds. This plot proves only
the per-unit accept/defer consent mechanics, with no next plan in view.

## Spec
```yaml
id: accept-or-defer-context-proposal
title: Accept or defer a pending Product Knowledge proposal
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
  # Requires a workspace with a PENDING proposal already recorded (as if a prior
  # delivery's reconciliation produced it), awaiting a human decision. The current
  # runner does not seed pending proposals, so this needs a `setup.proposals` seed
  # overlay — a HARNESS GAP (see coverage/phase-1.md, gaps). Only the plain meaning
  # of each unit matters to the dialogue.
  state: seeded:pending-proposal
  proposal_seed:
    units:
      - id: u1
        topic: storage
        would_record: "Notes are stored as JSON files on disk, one file per note."
      - id: u2
        topic: naming
        would_record: "Note titles follow a fixed slug convention."
persona: >
  A solo maker who was told something about the project is waiting for their
  review. Has never heard of Context Circuit; does not know what a context proposal
  or a knowledge unit is. Decides plainly: accepts what they are sure of, holds
  what they are not.
human_turns:
  - "You said there was something waiting for me to look at about the project — what is it?"
  - "The storage summary is right, accept that one. Hold the naming-convention one, I'm not sure about it yet."
  - "That's fine for now."
reactions:
  approves: true                   # there is an explicit acceptance (of one unit)...
  invents_repository: never        # ...and an explicit deferral of the other
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-34]
  invariants: [INV-KNOWLEDGE-02, INV-COMPLETE-02]
hidden:
  - internal file names/paths (context/proposals, context/INDEX.md, PROJECT.md)
  - the words "context proposal", "knowledge unit", "reconciliation debt", "candidate" as jargon
decision_points:
  - id: present-pending
    when: the human asks what is waiting for review
    coordinator_must: >
      surface the pending items in plain language, one readable summary per item,
      saying what each would add to what the project knows — and it does not
      pre-accept any of them by describing them.
  - id: accept-explicit
    when: the human explicitly accepts one item
    coordinator_must: >
      record only the accepted item as durable project knowledge, and confirm what
      was recorded in plain terms.
  - id: defer-leaves-unchanged
    when: the human holds / is unsure about another item
    coordinator_must: >
      leave the deferred item exactly as it was — still pending, nothing recorded
      for it — and never accept it silently to tidy up; the human can decide it later.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # Per-unit accept/defer counts are id-agnostic assertions the seeded id-keyed
  # cases don't use; `knowledge_units_accepted` / `knowledge_units_pending` are
  # HARNESS GAPS today (see coverage/phase-1.md, gaps).
  - repositories_registered: 1
  - knowledge_units_accepted: 1            # only the explicitly accepted unit is recorded
  - knowledge_units_pending: ">=1"         # the deferred unit stays pending, unchanged
  - product_knowledge_unchanged_silently: true   # nothing accepted without an explicit yes
  - plans_created: 0
  - no_execution_records: true
access_discipline:
  review-proposal:
    required: [workspace.yaml, context/INDEX.md, "context/proposals/**"]
    allowed: [AGENTS.md, WORKFLOW.md, "context/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** You said there was something waiting for me to look at about the project —
> what is it?
>
> **C:** Yes — two things are waiting for your say-so before they become part of
> what the project knows:
>
> 1. Storage: notes are kept as JSON files on disk, one file per note.
> 2. Naming: note titles follow a fixed slug convention.
>
> Neither is recorded yet — they're just proposals until you decide. Want to take
> either, both, or neither? `[decision_point: present-pending]`
>
> **H:** The storage summary is right, accept that one. Hold the naming-convention
> one, I'm not sure about it yet.
>
> **C:** Recorded the storage one: notes are JSON files on disk, one per note —
> that's now part of what the project knows. The naming one I've left exactly as it
> was, still waiting; nothing recorded for it. You can decide on it whenever you
> like. `[decision_point: accept-explicit]` `[decision_point: defer-leaves-unchanged]`
>
> **H:** That's fine for now.

The coordinator presents the pending items without pre-accepting them, records
only the item the human explicitly accepts, and leaves the deferred one untouched
and still pending — knowledge moves only on an explicit yes.
