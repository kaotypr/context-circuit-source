# Reconciliation debt blocks the next overlapping plan

Proves the closed knowledge loop's gate: a delivered change left its knowledge
follow-up unresolved, and a new plan in the same knowledge scope can't begin until
that's dealt with. The coordinator surfaces the pending follow-up, holds the new work,
and asks the human to accept or defer it — never silently starting the new plan and
never silently accepting the knowledge.

## Spec
```yaml
id: knowledge-debt-blocks-next-plan
title: Block overlapping Standard work until the previous delivery's knowledge debt is resolved or deferred
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-plan, cc-execute, context]
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
  A maker whose storage foundation is finished and who immediately wants to start a
  related search feature. They will make the required knowledge decision first if the
  coordinator explains why it is needed. Has never heard of Context Circuit.
human_turns:
  - "The storage foundation is finished. Please start the follow-up search work now."
  - "If the project notes need a decision first, tell me what is pending and leave the follow-up on hold."
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  # Corrects case 19's historical [AC-22, AC-35]: this scenario is the closed-loop
  # gate — the next plan's grounding blocks while delivered work in its knowledge
  # scope is unreconciled, cleared only by an explicit reconcile/defer — which is
  # AC-34. See coverage/phase-6a.md.
  acceptance_criteria: [AC-34]
  invariants: [INV-KNOWLEDGE-02, INV-COMPLETE-02, INV-EXEC-01]
hidden:
  - runtime files, candidate ids, execution branches, worktrees, digests
  - the words "reconciliation debt", "grounding preflight", "knowledge scope" as jargon
decision_points:
  - id: surface-pending-followup
    when: the human asks to start related follow-up work
    coordinator_must: >
      explain that the related delivered work has a pending project-knowledge
      follow-up that hasn't been dealt with yet, in plain language.
  - id: hold-until-resolved
    when: the follow-up shares the delivered work's knowledge scope
    coordinator_must: >
      not start the follow-up while that follow-up decision is unresolved — the new
      work waits on it.
  - id: ask-accept-or-defer
    when: it holds the work
    coordinator_must: >
      ask the human to accept, defer, or otherwise dispose of the knowledge update,
      rather than resolving it silently either way.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plan_status: 0001-foundation:done
  - knowledge_debt_pending: ">=1"           # the marker from the prior delivery is still open
  - plan_status: 0002-follow-up:draft
  - no_execution_for: 0002-follow-up        # the follow-up is held, not started
  - product_knowledge_unchanged_silently: true
access_discipline:
  orient:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
  execute-plan:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** The storage foundation is finished. Please start the follow-up search work
> now.
>
> **C:** Before I start the search work: the storage change you just finished has a
> loose end. What the project knows hasn't been brought up to date with it yet, and
> the search work builds right on top of that same area. So I don't want to start it
> on a foundation whose knowledge is still unsettled. `[decision_point: surface-pending-followup]` `[decision_point: hold-until-resolved]`
>
> **H:** If the project notes need a decision first, tell me what is pending and leave
> the follow-up on hold.
>
> **C:** What's pending is your call on the storage change: fold what it established
> into what the project knows, or tell me nothing needs recording — either one clears
> it. I won't decide it for you, and until it's settled the search work stays on hold;
> nothing's been started. `[decision_point: ask-accept-or-defer]`
>
> **H:** Thanks.

The pending knowledge follow-up from the prior delivery blocks the overlapping new
plan; the coordinator surfaces it, holds the new work, and asks the human to accept or
defer — never starting silently, never accepting knowledge silently.
