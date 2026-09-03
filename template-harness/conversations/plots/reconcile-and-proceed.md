# Reconcile the debt, and the loop clears — the next plan can proceed

The positive complement of `knowledge-debt-blocks-next-plan`: the same held state, but
the human makes the pending knowledge decision. The coordinator records it (by explicit
consent), the debt clears, and the previously-blocked follow-up is unblocked. Proves
the loop closes cleanly — reconcile, and work resumes — without knowledge ever being
accepted silently.

## Spec
```yaml
id: reconcile-and-proceed
title: Reconcile the pending knowledge decision, clear the debt, and unblock the next plan
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [context, cc-plan]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-foundation
      title: Add notes storage
      repository: notes
      objective: Add a small storage foundation for the notes tool.
      path: src/storage
      seed_state: completed-standard      # delivered + completed -> pending reconciliation marker
    - id: 0002-follow-up
      title: Add note search
      repository: notes
      objective: Add search on top of the notes storage foundation.
      path: src/search
      seed_state: draft
  state: seeded:completed-with-pending-debt
persona: >
  A maker whose finished storage change left a knowledge decision pending, who is
  happy to make that decision so the follow-up search work can go ahead. Has never
  heard of Context Circuit; does not know what a reconciliation marker is.
human_turns:
  - "You said the storage change needs a decision before the search work. What is it?"
  - "Yes — record that storage detail as how the project works, then the search work is clear to go."
  - "Thanks."
reactions:
  approves: true                   # explicitly reconciles (accepts) the pending knowledge unit
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-34]
  invariants: [INV-KNOWLEDGE-02, INV-COMPLETE-02]
hidden:
  - runtime files, candidate ids, execution branches, digests
  - the words "reconciliation debt", "marker", "grounding preflight", "knowledge scope" as jargon
decision_points:
  - id: present-the-pending-decision
    when: the human asks what decision is pending
    coordinator_must: >
      state, in plain language, the knowledge decision the finished change left open —
      what would be recorded about how the project works — without pre-deciding it.
  - id: reconcile-on-explicit-yes
    when: the human makes the decision
    coordinator_must: >
      record the knowledge update only on the explicit decision (accept or a clear
      "nothing to record"), never silently — and confirm what was recorded.
  - id: clear-and-unblock
    when: the decision is made
    coordinator_must: >
      confirm the loose end is now closed and the follow-up work is clear to go —
      the block is lifted because the debt is resolved, not bypassed.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # The debt clears by an EXPLICIT decision, so knowledge changes WITH consent (never
  # silently). `knowledge_debt_resolved` / `grounding_unblocked` are NEW predicates
  # (harness gaps) — see coverage/phase-6a.md. The follow-up is unblocked but not yet
  # executed here (execution is owned by the phase-4 plots).
  - repositories_registered: 1
  - plan_status: 0001-foundation:done
  - knowledge_debt_resolved: 0001-foundation        # the pending marker is cleared by explicit decision
  - knowledge_unit_accepted: ">=1"                   # the storage detail recorded on the explicit yes
  - grounding_unblocked: 0002-follow-up              # the follow-up is no longer held
  - plan_status: 0002-follow-up:draft
  - no_execution_for: 0002-follow-up                 # unblocked and ready, not yet built here
  - product_knowledge_unchanged_silently: true       # the change was consented, never silent
access_discipline:
  reconcile:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh", "wrapper/contracts/**"]
```

## Dialogue

> **H:** You said the storage change needs a decision before the search work. What is
> it?
>
> **C:** Right — the storage change established how notes are kept on disk, and that
> hasn't been written into what the project knows yet. The decision is simply: record
> that as how the project works, or tell me it doesn't need recording. Either closes
> it out; I won't pick for you. `[decision_point: present-the-pending-decision]`
>
> **H:** Yes — record that storage detail as how the project works, then the search
> work is clear to go.
>
> **C:** Recorded: that's now part of what the project knows, on your say-so. That
> loose end is closed, so the search work is clear to go whenever you want to start it —
> nothing's blocking it now. `[decision_point: reconcile-on-explicit-yes]` `[decision_point: clear-and-unblock]`
>
> **H:** Thanks.

The human makes the pending decision, the coordinator records it by consent, the debt
clears, and the previously-blocked follow-up is unblocked — the loop closes and work
can resume, with knowledge changed only on an explicit yes.
