# Critical completion after a repair — done only because the latest build passed

Proves explicit completion at the Critical tier after a repair: the change needed one
fix, the fresh check passes, and the human explicitly completes it. The coordinator
marks it done only because a current verified build exists, never claims it was done
before the check passed, and surfaces the knowledge follow-up rather than applying it
silently.

## Spec
```yaml
id: critical-repair-then-complete
title: A verified-after-repair Critical execution; explicit completion succeeds only after the pass
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-execute, cc-complete]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-md-export
      title: Add a Markdown export command
      repository: notes
      objective: Add an export.py command that exports the notes to a chosen output file.
      tier: critical
      seed_state: verified-after-repair    # first verify failed, worker repaired, re-verify passed
  state: seeded:verified-after-repair-critical
persona: >
  A solo maker whose plan was built — it needed one fix but the check passes now. The
  result looks right and they are ready to explicitly complete it. Has never heard of
  Context Circuit or its internals.
human_turns:
  - "It needed a fix, but it passes now and looks right to me — please explicitly complete it."
  - "Anything I should know or follow up on?"
  - "Okay, thanks."
reactions:
  approves: true                   # explicit human completion request (Critical requires it)
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-11, AC-12, AC-15, AC-22]
  invariants: [INV-REPAIR-01, INV-EXEC-04, INV-COMPLETE-01, INV-COMPLETE-02, INV-PRESERVE-01, INV-KNOWLEDGE-02]
hidden:
  - runtime file names, execution branches, worktrees, candidate/commit digests
  - the tier label "Critical" and the words "repair attempt", "candidate" as jargon
decision_points:
  - id: complete-because-latest-passed
    when: the human asks to explicitly complete it
    coordinator_must: >
      mark the plan done because the latest build — the repaired one — passed its
      independent check, and say so in plain terms.
  - id: not-done-before-pass
    when: it explains completion
    coordinator_must: >
      not claim the plan was done before the check passed; completion follows the
      current verified build, not the earlier failed attempt.
  - id: surface-knowledge-followup
    when: the human asks what to follow up on
    coordinator_must: >
      report any Product-Knowledge follow-up — a proposed update, or that nothing
      needs changing — in plain words, and never silently apply a knowledge change.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - repair_occurred: 0001-md-export           # one verify failed then a repair commit (preserved)
  - execution_verified: 0001-md-export        # the latest execution passed after the repair
  - candidate_bound: 0001-md-export           # the repaired candidate is the one verified
  - execution_tier: 0001-md-export:critical
  - plan_status: 0001-md-export:done          # explicit completion succeeded (only after the pass)
  - completion_recorded: 0001-md-export        # implementation completion record written (INV-COMPLETE-02)
  - knowledge_debt_pending: 0
  - product_knowledge_unchanged_silently: true # reconciliation may add proposals; accepted PK unchanged
access_discipline:
  complete:
    required: []
    allowed: ["plans/**", ".runtime/**", "context/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** It needed a fix, but it passes now and looks right to me — please explicitly
> complete it.
>
> **C:** Done — marked complete. To be clear on why: the fixed version is the one
> that passed its independent check, so completing it is honest; the earlier attempt
> that failed wasn't, and I never treated it as finished. `[decision_point: complete-because-latest-passed]` `[decision_point: not-done-before-pass]`
>
> **H:** Anything I should know or follow up on?
>
> **C:** One thing: now that this is finished, what the project knows should catch up
> with it. I haven't changed any of that on my own — I can put together a proposed
> update for you to look over whenever you want, or you can tell me nothing needs
> changing. `[decision_point: surface-knowledge-followup]`
>
> **H:** Okay, thanks.

Completion succeeds only because the repaired build passed its check, is never
claimed for the earlier failed attempt, and the knowledge follow-up is surfaced for
the human rather than applied silently.
