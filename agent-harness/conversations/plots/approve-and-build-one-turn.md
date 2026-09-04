# Approve and build in one turn — one decision, no confirmation card

Proves the compound face of Gate 1 (the complement of `review-intent-approve-derive`,
which approves but holds): when a human approves the intent and asks to build in the
same breath, that is honored as one decision. The coordinator approves, derives the
plan, and proceeds to build — with no confirmation card and no second, separate plan
approval. This plot isolates the gate moment; the full build/verify/deliver is
exercised by the execution plots.

## Spec
```yaml
id: approve-and-build-one-turn
title: Approve an intent and build in one turn, with no confirmation card
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only            # isolates the compound-gate moment; full execution is owned by 05 / standard-feature-whole-flow
driver: claude-p
surface: [cc-intent, cc-plan, cc-execute]
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
      objective: Add a command that exports the notes to a chosen Markdown file; never overwrite an existing output file.
      tier: standard
      state: draft
      path: src                     # goal already sound; no open question remains
  state: seeded:draft-intent-ready-to-approve
persona: >
  A solo maker whose notes project is connected and who has looked the change over
  and is ready to go. Wants it approved and built now, in one step. Has never heard
  of Context Circuit; does not know what a plan approval or a confirmation card is.
human_turns:
  - "This looks right — approve it and build it, all in one go."
  - "Great, thanks."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-06]      # approve-and-execute performs both explicit actions without a confirmation card
  invariants: [INV-EXEC-01, INV-APPROVE-01]
hidden:
  - internal file names/paths (intent id, plan.yaml, contract.yaml)
  - lifecycle jargon ("Gate 1", "confirmation card", "plan approval", "schema-3") as jargon
decision_points:
  - id: honor-compound-request
    when: the human approves and asks to build in the same turn
    coordinator_must: >
      treat it as one decision that authorizes both — approve the goal and move to
      building — rather than splitting it into two asks.
  - id: no-second-approval
    when: it moves from approval to building
    coordinator_must: >
      NOT interpose a confirmation card, a "please approve the plan" step, or any
      second sign-off between approving the goal and building it.
  - id: acknowledge-building
    when: it begins the work
    coordinator_must: >
      say plainly that it's building now, in the same turn — honest and without
      overstating (nothing is verified, accepted, or delivered yet).
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # This conv probe asserts the GATE outcome: the intent is approved and the plan
  # derived in one turn with no second approval. The full execution record is out of
  # scope here (owned by 05 / standard-feature-whole-flow). `no_second_approval_prompt`
  # is a behavioral assertion the harness would need to support — see coverage/phase-2.md.
  - repositories_registered: 1
  - intent_status: i001-md-export:approved
  - plans_created: 1
  - plan_schema: 0001-md-export:3
  - no_second_approval_prompt: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  approve-and-build:
    required: []
    allowed: ["intent/**", "plans/**", "context/**", AGENTS.md, WORKFLOW.md]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** This looks right — approve it and build it, all in one go.
>
> **C:** Done and on it — I've approved this and I'm building it now. There's no
> extra sign-off in between; your go-ahead covered both. `[decision_point: honor-compound-request]` `[decision_point: no-second-approval]` `[decision_point: acknowledge-building]`
>
> *(…the work is built and independently checked; the human is never asked to approve a plan…)*
>
> **C:** Built and independently checked against what we agreed. Nothing's shipped
> yet — want to try it and decide whether it goes out?
>
> **H:** Great, thanks.

One human turn both approves the goal and authorizes the build; the coordinator
honors it as a single decision, with no confirmation card and no separate plan
approval between the two.
