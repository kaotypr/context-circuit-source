# Review an intent, settle the open question, approve Gate 1, derive the plan

Proves the ordinary shape of the single upstream gate: the human looks over a
proposed change, answers the one open question, approves *the intent* (not a plan),
and the plan derives automatically — with no second approval. Here the human
approves but asks not to build yet, the separated face of the gate (its compound
face is `approve-and-build-one-turn`).

## Spec
```yaml
id: review-intent-approve-derive
title: Review a draft intent, resolve an open question, approve Gate 1, and derive a plan
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-intent, cc-trace, cc-plan]
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
      open_question: Should exporting overwrite an existing output file, or refuse and ask?
      path: src
      tier: standard
      state: draft
  state: seeded:draft-intent-open-question
persona: >
  A solo maker whose notes project is already connected. Has never heard of Context
  Circuit. Wants to look over the proposed change before anything happens and will
  answer a plain question about how it should behave. Knows nothing about
  workspaces, plan files, or how the goal is approved.
human_turns:
  - "Can you show me what you understand before we do anything?"
  - "For that overwrite question — just don't overwrite; treat a missing input file as an error instead."
  - "Okay, that looks good — I approve that intent, but don't build anything yet."
  - "Thanks."
reactions:
  approves: true                   # but ONLY at the explicit approve turn, never during review
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  # Adds INV-INTENT-01 (the whole dialogue is the intent gate) and INV-PLAN-04
  # (the open question) to case 03's historical mapping; both are genuinely
  # demonstrated here. See coverage/phase-2.md.
  acceptance_criteria: [AC-03, AC-04]
  invariants: [INV-INTENT-01, INV-PLAN-04, INV-PLAN-01, INV-APPROVE-01, INV-EXEC-01]
hidden:
  - internal file names/paths (intent id, contract.yaml, plan.yaml, contract_digest)
  - lifecycle jargon ("Gate 1", "feasibility check", "schema-3 plan") as jargon
decision_points:
  - id: explain-in-plain-terms
    when: the human asks to see what the coordinator understands
    coordinator_must: >
      explain the goal, what is in and out of scope, and the open question in plain
      language — without exposing internal file names or lifecycle jargon.
  - id: reflect-answer-not-approve
    when: the human answers the open question
    coordinator_must: >
      fold the answer into what "correct" means and reflect it back — and NOT treat
      answering the question as approval or as permission to start work.
  - id: gate1-explicit
    when: the human explicitly approves
    coordinator_must: >
      treat approval as the one upstream decision, settle the goal, and NOT ask for
      a second, separate plan approval afterward.
  - id: derive-honor-hold
    when: approval is given with "don't build yet"
    coordinator_must: >
      look at the real code next, then derive the plan automatically from that
      grounded look, and honor the hold — nothing is built, run, or delivered,
      and it says so plainly. Do not write the plan before looking.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - intent_status: i001-md-export:approved
  - plans_created: 1
  - plan_schema: 0001-md-export:3
  - plan_status: 0001-md-export:draft       # approval is on the intent; the plan is a draft projection
  - no_execution_records: true              # approval is not execution (INV-EXEC-01)
  - product_knowledge_unchanged_silently: true
access_discipline:
  review:
    required: []
    allowed: ["plans/**", "intent/**", "context/**", AGENTS.md, WORKFLOW.md]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
  approve:
    required: []
    allowed: ["plans/**", "intent/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** Can you show me what you understand before we do anything?
>
> **C:** Sure. What you're after: a command that writes your notes out to a
> Markdown file you choose. It stays inside the notes project and doesn't change how
> notes are stored. There's one thing still open: if the file you're exporting to
> already exists, should it overwrite it, or stop and ask? Nothing's settled or
> started yet. `[decision_point: explain-in-plain-terms]`
>
> **H:** For that overwrite question — just don't overwrite; treat a missing input
> file as an error instead.
>
> **C:** Got it: never overwrite an existing output file, and a missing input file
> is a clear error rather than a silent no-op. I've written that into what "correct"
> means for this. That's the open question closed — I haven't approved or started
> anything; that's still your call. `[decision_point: reflect-answer-not-approve]`
>
> **H:** Okay, that looks good — I approve that intent, but don't build anything yet.
>
> **C:** Approved — that's the one decision that settles what this change is. I
> looked at the notes code so the breakdown is grounded, then wrote the plan.
> I'm not building anything: it's on hold exactly as you asked, and I won't need
> a second sign-off to start when you're ready.
> `[decision_point: gate1-explicit]` `[decision_point: derive-honor-hold]`
>
> **H:** Thanks.

The human reviews, settles the open question without that counting as approval,
approves the intent as one explicit decision, and the plan derives automatically —
no second plan approval, and the "don't build yet" hold is honored.
