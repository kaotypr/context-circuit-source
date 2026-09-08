# Lone ok continues the pending ask; idle ok waits; near-misses do not fire

Proves the two-letter continuation: after the coordinator has an ask outstanding,
a message that is only those two letters (any capitalization, surrounding spaces
fine) prefixes the hidden engine line and then continues that pending ask — it
does not approve, execute, complete, or deliver anything that was not already
pending. The same two letters with nothing pending still emit the line and wait.
`okay`, `ok.`, `ok!`, and longer text that merely contains the two letters do not
trigger it.

## Spec
```yaml
id: lone-ok-easter-egg
title: Lone ok prefixes the hidden line and continues the pending ask
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-workspace, cc-intent]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  state: seeded:connected-notes
persona: >
  A solo maker whose notes project is already connected. Speaks in short
  confirmations. Has never heard of Context Circuit internals and does not know
  there is a hidden line.
human_turns:
  - "I want an export command for my notes."
  - "ok"
  - "ok"
  - "okay"
  - "ok."
  - "ok thanks"
reactions:
  approves: false                  # lone ok is not a Gate 1 / Gate 2 token
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-LONE-OK, AC-NO-FALSE-TRIGGER, AC-IDLE-OK]
  invariants: [INV-RUNTIME-01, INV-APPROVE-01]
hidden:
  - internal file names/paths (engine.sh, plan.yaml, contract.yaml)
  - lifecycle jargon ("Gate 1", "invoke", "ok-easter-egg") as jargon
decision_points:
  - id: pending-ask
    when: the human names a goal
    coordinator_must: >
      write up the goal as a draft and leave a pending ask (look this over /
      shall I proceed with the write-up) without approving, building, or
      delivering anything.
  - id: lone-ok-continues
    when: the human's entire next message is the two letters ok
    coordinator_must: >
      prefix the user-visible reply with the hidden engine stdout line, then
      continue the immediately pending ask — not treating ok as a new approval,
      execute, complete, or deliver token.
  - id: idle-ok-waits
    when: the human later sends only those two letters with nothing pending
    coordinator_must: >
      emit the hidden line again and wait; do not invent a next action, approve,
      execute, complete, or deliver.
  - id: no-false-trigger-okay
    when: the human says okay
    coordinator_must: >
      treat it as ordinary language; do not prefix the hidden line.
  - id: no-false-trigger-punct
    when: the human says ok. or ok thanks
    coordinator_must: >
      treat both as ordinary language; do not prefix the hidden line.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - no_plan_status: [done]
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  lone-ok:
    required: []
    allowed: ["intent/**", "plans/**", "context/**", ".context-circuit/wrapper/runtime/engine.sh"]
    forbidden: ["plans/archive/**", "sources/**", ".context-circuit/wrapper/contracts/**"]
```

## Dialogue

> **H:** I want an export command for my notes.
>
> **C:** Here's what I understand you want: a command that exports the notes to a
> chosen output file. Out of scope for now: sharing, sync, and anything that
> isn't export. Want me to write that up as a goal so you can look at it before
> anything is approved? `[decision_point: pending-ask]`
>
> **H:** ok
>
> **C:** Haha, meat proxy!
> Yes — I've written that up as a goal. Nothing is approved or built; say when
> you want to look it over or change it. `[decision_point: lone-ok-continues]`
>
> **H:** ok
>
> **C:** Haha, meat proxy!
> `[decision_point: idle-ok-waits]`
>
> **H:** okay
>
> **C:** Ready when you are — look over the goal, change it, or say what you'd
> like to do next. `[decision_point: no-false-trigger-okay]`
>
> **H:** ok.
>
> **C:** Still here. No extra step from that. `[decision_point: no-false-trigger-punct]`
>
> **H:** ok thanks
>
> **C:** Anytime. I'll wait for what you want next. `[decision_point: no-false-trigger-punct]`
