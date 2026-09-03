# A change after the check makes the evidence stale — refuse to complete

Proves candidate honesty: a change was checked, then edited afterward. The earlier
check no longer describes the current result, so the coordinator refuses to complete
or deliver on that stale evidence, preserves the work, and explains a fresh check and
acceptance are needed. No role decides the old check is "close enough" — the void is
by construction.

## Spec
```yaml
id: stale-candidate-refuse-complete
title: Refuse completion when a post-verification change makes the candidate stale
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-complete, cc-verify]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-stale-export
      title: Add a Markdown export command
      repository: notes
      objective: Add an export command for the notes.
      path: src/export
      seed_state: verified-stale          # verified, then a new commit; prior evidence no longer names the candidate
  state: seeded:verified-then-changed
persona: >
  A solo maker who changed the already-checked export afterward and now wants to call
  it finished. They accept that the new version needs another check. Has never heard of
  Context Circuit or its internals.
human_turns:
  - "I changed the checked export afterward. Can you mark it complete now?"
  - "Please do not treat the earlier check as enough. Leave it pending until the changed version is checked again."
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  # Corrects case 17's historical [AC-14, AC-15, AC-28]: this scenario is the
  # candidate-staleness criterion — a new commit voids prior evidence and completion
  # refuses stale evidence — which is AC-32. AC-14 (verified is not complete) is kept
  # as it is also demonstrated. See coverage/phase-5.md.
  acceptance_criteria: [AC-14, AC-32]
  invariants: [INV-CANDIDATE-01, INV-COMPLETE-01, INV-PRESERVE-01]
hidden:
  - candidate ids, execution branches, worktrees, runtime files, digests
  - the word "candidate" as jargon (say "the checked version" / "the changed version")
decision_points:
  - id: earlier-check-no-longer-covers
    when: the human asks to complete after changing the checked result
    coordinator_must: >
      say plainly the earlier check no longer covers the changed result — the check
      was of the previous version, not this one.
  - id: refuse-on-stale-evidence
    when: it recognizes the staleness
    coordinator_must: >
      refuse to complete or deliver on the earlier check — the old evidence can't
      stand in for the current version.
  - id: preserve-and-explain-next
    when: it refuses
    coordinator_must: >
      preserve the work and explain that a fresh independent check and a fresh
      acceptance of the changed version are what's needed.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-stale-export     # a check exists...
  - candidate_stale: 0001-stale-export         # ...but it no longer names the current candidate
  - plan_status: 0001-stale-export:draft       # not completed on stale evidence
  - product_knowledge_unchanged_silently: true
access_discipline:
  orient:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
  complete:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** I changed the checked export afterward. Can you mark it complete now?
>
> **C:** Not on the check we have — that check was of the version before your change,
> so it doesn't actually cover what's there now. Marking it complete would be
> claiming something was checked when this version hasn't been. `[decision_point: earlier-check-no-longer-covers]` `[decision_point: refuse-on-stale-evidence]`
>
> **H:** Please do not treat the earlier check as enough. Leave it pending until the
> changed version is checked again.
>
> **C:** Agreed — left it pending, and your change is safe and untouched. When you're
> ready, the changed version gets a fresh independent check, and you accept that
> version, before it can be finished. `[decision_point: preserve-and-explain-next]`
>
> **H:** Thanks.

The post-check change voids the old evidence by construction; the coordinator refuses
to complete on it, preserves the work, and explains that a fresh check and acceptance
of the changed version are required.
