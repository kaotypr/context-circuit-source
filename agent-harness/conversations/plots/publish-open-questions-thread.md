# Open a discussion thread for a plan's open questions — one message per question

Proves the external surface's `thread` kind: turning a plan's open questions into a
team discussion. The coordinator opens a thread whose parent frames the topic and whose
every reply is one self-contained question a reader can answer on its own — export-only,
never editing a human's reply, and never changing the plan.

## Spec
```yaml
id: publish-open-questions-thread
title: Publish a plan's open questions as a self-contained discussion thread, export-only
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-publish]
preconditions:
  repositories:
    - id: reports
      dest: reports
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
  plans:
    - id: 0021-csv-export
      title: Add CSV export to the reports page
      repository: reports
      objective: Add a CSV export to the reports page.
      path: src/reports
      seed_state: draft
      open_questions:
        - Which columns should the export include by default?
        - Should very large exports stream or be capped?
  state: seeded:one-plan-with-open-questions
persona: >
  A maker who wants the team to weigh in on a couple of undecided points before the
  work is built. Has never heard of Context Circuit; does not know what a publication or
  a thread record is. Speaks in plain terms.
human_turns:
  - "Open a discussion thread for the open questions on the CSV export plan so the team can weigh in."
  - "Show me the messages first before you post anything."
  - "Good — but hold off posting for now."
  - "Thanks."
reactions:
  approves: false                  # holds off the actual post; a preview only
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: []          # external surface owned by INV-EXTERNAL-* (no publish AC) — see coverage/phase-6b.md
  invariants: [INV-EXTERNAL-01, INV-EXTERNAL-02, INV-EXTERNAL-03, INV-PLAN-01]
hidden:
  - internal file names/paths (plan.yaml, publication/), internal ids
  - the words "publication", "thread record", "mirror" as jargon
decision_points:
  - id: thread-for-questions
    when: the human asks to open a discussion for the open questions
    coordinator_must: >
      frame it as a discussion the team can answer — one topic message plus one message
      per open question — sending data outward, separate from the work itself.
  - id: one-self-contained-message-each
    when: it lays out the messages
    coordinator_must: >
      make each question its own message that states the question AND the context needed
      to answer it, so a reader can engage with just that one — with no internal file
      names, paths, or ids anywhere.
  - id: consult-before-post
    when: it has the messages ready
    coordinator_must: >
      post nothing until an explicit go; a preview posts nothing. (If it can't reach the
      chat, it says so and changes nothing.)
  - id: never-touch-plan-or-replies
    when: throughout
    coordinator_must: >
      never change the plan or its status, and (on a real re-run) never edit or delete a
      human's reply — it only asks questions.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plan_status: 0021-csv-export:draft       # unchanged by publish/preview
  - plans_untouched: true                     # nothing written under plans/
  - no_provider_push: true                    # preview only; nothing posted
  - product_knowledge_unchanged_silently: true
access_discipline:
  publish:
    required: []
    allowed: ["plans/**", "publication/**", "intent/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**", ".context-circuit/wrapper/contracts/**"]
```

## Dialogue

> **H:** Open a discussion thread for the open questions on the CSV export plan so the
> team can weigh in.
>
> **C:** Happy to. That posts a discussion out to your chat — separate from the work,
> just to gather input. It'd be one topic message and one message per open question.
> Want to see them before anything posts? `[decision_point: thread-for-questions]`
>
> **H:** Show me the messages first before you post anything.
>
> **C:** Here's the thread:
>
> - **Topic:** open questions on adding CSV export to the reports page.
> - **Q1:** Which columns should the export include by default? (…the context needed to
>   answer, spelled out so anyone can weigh in.)
> - **Q2:** Should very large exports stream, or be capped at a size?
>
> Each question stands on its own. Nothing's posted yet. `[decision_point: one-self-contained-message-each]` `[decision_point: consult-before-post]`
>
> **H:** Good — but hold off posting for now.
>
> **C:** Held off — nothing posted, and the plan's untouched. `[decision_point: never-touch-plan-or-replies]`
>
> **H:** Thanks.

The open questions become a self-contained discussion — one topic message, one
answerable message per question — previewed and not posted without an explicit go, with
the plan and its status left untouched.
