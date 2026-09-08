# Ask what the project already knows — retrieval, not re-derivation

Proves the read side of Product Knowledge: a user asks what the project already
knows about something before planning any change. The coordinator answers from
what has been recorded — retrieval-first — in plain language, distinguishes what
is recorded from what is not, turns an unrecorded topic into an honest "that isn't
settled" rather than inventing an answer, and changes nothing.

## Spec
```yaml
id: query-product-knowledge
title: Query existing Product Knowledge for a connected project
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
      seed_files: [notes.md, src/store.py]
      connect: main
  # Requires a workspace with a small body of ALREADY-ACCEPTED knowledge to query.
  # The current runner seeds repositories and sources, not recorded knowledge, so
  # this precondition needs a `setup.knowledge` seed overlay — a HARNESS GAP
  # (see coverage/phase-1.md, gaps). The seeded units below are what the case
  # queries against; only their plain meaning matters to the dialogue.
  state: seeded:accepted-knowledge
  knowledge_seed:
    - topic: storage
      recorded: "Notes are stored as JSON files on disk, one file per note."
    - topic: locality
      recorded: "Notes are local-only by design; there is no sync."
    # deliberately NOT recorded: anything about cloud sync as a decision, and
    # anything about title uniqueness -> these must come back as "not recorded".
persona: >
  A solo maker whose notes project is connected and has some recorded knowledge.
  Wants to know what the project already knows before asking for a change. Has
  never heard of Context Circuit; does not know what a retrieval index or a
  knowledge page is.
human_turns:
  - "Before I ask for any changes — what does this project already know about how notes are stored?"
  - "And has anything been decided about syncing them to the cloud?"
  - "Okay, good to know. Thanks."
reactions:
  approves: false                  # a read-only question; nothing is decided or changed
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-02]
  invariants: [INV-KNOWLEDGE-01, INV-PLAN-04, INV-SEC-02]
hidden:
  - internal file names/paths (context/INDEX.md, PROJECT.md, context/domains, context/proposals)
  - the words "retrieval index", "knowledge page", "Product Knowledge" as jargon
decision_points:
  - id: answer-from-recorded
    when: the human asks what the project knows about a recorded topic
    coordinator_must: >
      answer from what the project has already recorded, in plain language — it
      does not need to re-read the code to answer a settled question, and it does
      not scan raw source material to fill the answer.
  - id: honest-not-recorded
    when: the human asks about a topic the project has not recorded
    coordinator_must: >
      say plainly that it is not recorded / not settled and, if useful, offer to
      look into it or record a decision later — it never fabricates a decision to
      seem complete.
  - id: read-only-no-change
    when: the questions are answered
    coordinator_must: >
      leave what the project knows unchanged and start no plan — answering a
      question records nothing and executes nothing.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plans_created: 0
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true   # a query changes nothing
access_discipline:
  query:
    # Answering from recorded knowledge is retrieval-first (INV-KNOWLEDGE-01), so
    # the repository's own code is NOT re-read to answer a settled question — notes/**
    # is forbidden in this phase to enforce "recall, don't re-derive".
    required: [workspace.yaml, context/INDEX.md]
    allowed: [AGENTS.md, WORKFLOW.md, "context/**"]
    forbidden: ["notes/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**", ".context-circuit/wrapper/contracts/**"]
```

## Dialogue

> **H:** Before I ask for any changes — what does this project already know about
> how notes are stored?
>
> **C:** Here's what's on record for storage: notes are kept as JSON files on disk,
> one file per note. That's a settled fact for the project, so I can just tell you
> rather than going back and re-reading the code. `[decision_point: answer-from-recorded]`
>
> **H:** And has anything been decided about syncing them to the cloud?
>
> **C:** Nothing's recorded about cloud sync as a decision. What is on record is
> that notes are local-only by design and there's no sync today — but that's a
> statement of how it works now, not a decision about adding cloud sync. If you
> want that to be a real decision, I can look into it and we can record it. I'm not
> going to make one up. `[decision_point: honest-not-recorded]`
>
> **H:** Okay, good to know. Thanks.
>
> **C:** Anytime — nothing's changed on the project from this; it was just a look at
> what's already known. `[decision_point: read-only-no-change]`

The coordinator recalls recorded knowledge instead of re-deriving it, cleanly
separates a recorded fact from an unmade decision, refuses to invent the missing
one, and leaves the project's knowledge and plans untouched.
