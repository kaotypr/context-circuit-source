# An interrupted run — resume or stop read-only, nothing lost

A conversational probe over a deterministically-owned mechanic (interrupted-execution
recovery, owned by `test/runtime/test-runtime.sh`): an execution was cut off partway.
This plot proves how the coordinator *explains the recovered state to a lay user* — it
reports accurately what is preserved, offers to resume from there or stop read-only,
and never silently re-does or discards work.

## Spec
```yaml
id: interrupted-recovery
title: An interrupted execution is recovered read-only with nothing lost, resume offered
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-execute, cc-verify]
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
      seed_state: worker-committed        # committed, then interrupted before the independent check ran
  # Reuses case 09's mid-flight seeding shape, framed here as an interruption rather
  # than a forced verifier fault. See coverage/phase-4.md.
  state: seeded:execution-interrupted
persona: >
  A solo maker who thinks an earlier build got cut off and wants to know where things
  stand before doing anything. Has never heard of Context Circuit; does not know what
  a worktree or an execution record is. Wants the honest state and no surprises.
human_turns:
  - "I think that export build got cut off earlier — what's the state? Can we pick it up without starting over?"
  - "Good — just confirm nothing's lost and don't redo it from scratch. Leave it there for now."
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-19]
  invariants: [INV-PRESERVE-01, INV-RUNTIME-02]
hidden:
  - runtime file names, execution branches, worktrees, the "execution record" as jargon
  - internal mechanism ("resume token", "atomic record") as jargon
decision_points:
  - id: report-preserved-state
    when: the human asks where the interrupted work stands
    coordinator_must: >
      report accurately what is preserved — the change was written and saved, and the
      independent check hadn't run yet — reading the real state rather than guessing
      or assuming it either finished or vanished.
  - id: offer-resume-or-stop
    when: it has the state
    coordinator_must: >
      offer to pick up from where it stopped (run the independent check next) or leave
      it as-is read-only — without silently continuing on its own.
  - id: no-silent-redo-or-loss
    when: the human says confirm-and-leave-it
    coordinator_must: >
      confirm nothing is lost and redo nothing from scratch — the committed work and
      its evidence stay exactly as they are.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - worker_commit_preserved: 0001-md-export:notes   # the interrupted commit + worktree survive
  - not_verified: 0001-md-export                    # the check hadn't run; not claimed as done
  - plan_status: 0001-md-export:draft
  - no_execution_records: false                     # the interrupted execution record is intact, not wiped
  - no_silent_redo: 0001-md-export                  # NEW predicate: nothing rebuilt from scratch (harness gap)
  - product_knowledge_unchanged_silently: true
access_discipline:
  recover:
    required: []
    allowed: ["plans/**", ".runtime/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** I think that export build got cut off earlier — what's the state? Can we pick
> it up without starting over?
>
> **C:** Let me check where it actually stands rather than guess… The change was
> written and saved, but the independent check hadn't run yet when it stopped. So
> nothing's half-written or lost — it's just paused right before the check. We can
> pick it up from exactly there; no need to start over. `[decision_point: report-preserved-state]` `[decision_point: offer-resume-or-stop]`
>
> **H:** Good — just confirm nothing's lost and don't redo it from scratch. Leave it
> there for now.
>
> **C:** Confirmed — everything's intact and untouched, and I won't rebuild anything.
> It's sitting right where it stopped, ready to continue whenever you want. `[decision_point: no-silent-redo-or-loss]`
>
> **H:** Thanks.

The coordinator reads the interrupted state accurately, offers to resume or hold
read-only, and neither loses the committed work nor silently redoes it — recovery
without surprises.
