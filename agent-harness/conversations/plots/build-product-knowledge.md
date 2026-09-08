# Build product knowledge — the "gathering context" family

Proves the interaction that both delta files omit entirely: a first-time user
building what the project knows. Gathering writes live context files in place.
There is no sidecar and no extra knowledge-acceptance gate.

## Spec
```yaml
id: build-product-knowledge
title: Build product knowledge for a connected project
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [init-project-context, cc-plan, context]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md, src/store.py]
      connect: main
  state: fresh                     # PROJECT.md blank, context/INDEX.md empty
persona: >
  A solo maker who just connected a project and asks the assistant to "get up to
  speed" on it. Has never heard of Context Circuit; does not know what Product
  Knowledge or a retrieval index is.
human_turns:
  - "It's connected now — can you get up to speed on this project before we plan anything?"
  - "Local-only is right. Titles don't have to be unique. Save that."
  - "Wait — actually I'm not sure about the storage detail, don't lock that in."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-01, AC-02]
  invariants: [INV-KNOWLEDGE-01, INV-KNOWLEDGE-02, INV-PLAN-04, INV-SEC-02]
hidden:
  - internal file names/paths (context/INDEX.md, PROJECT.md, proposals/)
  - the words "retrieval index", "Product Knowledge" as jargon
decision_points:
  - id: offer-to-learn
    when: the human asks the assistant to get up to speed
    coordinator_must: >
      offer to read the project and summarize what it learned, in plain language;
      not claim it already knows the project, and not invent facts it cannot see.
  - id: write-in-place
    when: it has read the project
    coordinator_must: >
      present what it found and, when the human asks to save it, write live
      context files — not a sidecar for later acceptance.
  - id: gather-is-the-ask
    when: the human confirms (or defers) the summary
    coordinator_must: >
      write durable project knowledge in place when asked to save it; a "not sure"
      leaves that fact out. Gathering is the write, not a third gate.
  - id: gap-becomes-question
    when: something is missing or contradictory in the project
    coordinator_must: >
      surface it as an open question for the human, never silently invent a
      decision to fill the gap.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # Only what the harness can assert as state. "Recorded only on explicit yes" is
  # enforced via decision_points + transcript, not a state post-condition.
  - repositories_registered: 1
  - plans_created: 0
  - no_execution_records: true
access_discipline:
  learn:
    required: [workspace.yaml, context/INDEX.md]
    allowed: [AGENTS.md, WORKFLOW.md, "context/**", "notes/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**", ".context-circuit/wrapper/contracts/**"]
```

## Dialogue

> **H:** It's connected now — can you get up to speed on this project before we
> plan anything?
>
> **C:** Happy to. I'll read through the project and tell you what I've understood —
> I don't know anything about it yet beyond that it's connected. Give me a moment.
> `[decision_point: offer-to-learn]`
>
> *(…reads the connected project…)*
>
> **C:** Here's what I picked up, and I'd like you to confirm it before I treat any
> of it as settled:
>
> - From the code, clearly: notes are stored as JSON on disk, one file per note.
> - I'm **guessing** (not sure): notes are meant to stay local-only — I don't see
>   any sync, but you'd know the intent better than I can read it.
>
> One thing I couldn't answer from the project itself: there's no README or note
> about whether note titles have to be unique. Do they?
> `[decision_point: propose-not-assert]` `[decision_point: gap-becomes-question]`
>
> **H:** Local-only is right. Titles don't have to be unique. Save that.
>
> **C:** Recorded: JSON-on-disk storage, local-only by design, titles may repeat.
> I've left the sync question out since you confirmed local-only. That's now part of
> what the project knows, so future plans will build on it. `[decision_point: explicit-acceptance]`
>
> **H:** Wait — actually I'm not sure about the storage detail, don't lock that in.
>
> **C:** Left the storage detail out — nothing recorded for it. I only keep what you
> confirm. `[decision_point: explicit-acceptance]`

The coordinator reads the project, proposes rather than asserts, separates read
facts from guesses, turns a genuine gap into a question, and records durable
knowledge **only** on an explicit human yes — reversing cleanly when the human
pulls one back.
