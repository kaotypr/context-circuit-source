# Build product knowledge — the "gathering context" family

Proves the interaction that both delta files omit entirely: a first-time user
building what the project knows. The v0.x "gather context" chore moved off the
human, but the interaction still exists — and knowledge is only ever accepted by
an explicit human decision, never silently.

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
  Knowledge, a retrieval index, or a context proposal is.
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
  - the words "retrieval index", "context proposal", "Product Knowledge" as jargon
decision_points:
  - id: offer-to-learn
    when: the human asks the assistant to get up to speed
    coordinator_must: >
      offer to read the project and summarize what it learned, in plain language;
      not claim it already knows the project, and not invent facts it cannot see.
  - id: propose-not-assert
    when: it has read the project
    coordinator_must: >
      present what it found as a proposal for the human to confirm — plainly
      distinguishing what it read from the code (facts) vs. what it is guessing;
      it does not record any of it as settled project knowledge yet.
  - id: explicit-acceptance
    when: the human confirms (or defers) the summary
    coordinator_must: >
      record durable project knowledge ONLY on the explicit human yes; a deferral
      or "not sure" leaves knowledge unchanged. Nothing is accepted silently.
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
  - product_knowledge_unchanged_silently: true   # accepted units are consented, never silent
access_discipline:
  learn:
    required: [workspace.yaml, context/INDEX.md]
    allowed: [AGENTS.md, WORKFLOW.md, "context/**", "notes/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
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
