---
id: 0006-enrich-plan-execution
target_context_unit: context/domains/plan-execution/README.md
operation: change
statement: >
  Enrich plan-execution with shipped facts and fix the "lease" wording to the
  shipped "ownership lock". Add the execution brief + immutable plan-snapshot,
  the repair-scope stop rule, the plan-changes-after-failure rule, the full
  resume-match set, scope/dependency preflight, and the missing worker do-not /
  handoff items.
evidence_refs:
  - .agents/skills/cc-execute/SKILL.md
  - agents/writer.md
  - wrapper/runtime/engine.sh   # cc_execution_begin (snapshot), cc_lock_*, cc_recovery_inspect, cc_repository_preflight
  - wrapper/contracts/schemas/execution.yaml
  - wrapper/contracts/schemas/worker-handoff.yaml
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich plan-execution and correct lease->lock wording

Corrections:
- Replace "exclusive writer lease" with "exclusive-create ownership lock"
  (INV-OWN-01; engine `cc_lock_*`). "Lease" is legacy.

Additions (shipped):
- Execution brief + immutable `plan-snapshot` (serialized PLAN.md/plan.yaml/tasks)
  that worker and verifier prompts read.
- Scope/dependency preflight before execution (approval, dependencies, clean
  anchor, bounded scope) — `cc_repository_preflight`.
- Repair-scope stop rule: a repair may not redesign the plan; if repair needs new
  scope, execution stops and the human is asked to change the plan.
- Plan-changes-after-failure: a repair continues in the same execution only when
  intended scope and acceptance are unchanged; otherwise a new execution/revision
  and the old->new relationship are recorded.
- Full resume-match set: resume only when plan revision, configured anchor
  branches and captured anchor commits, worktree paths, and ownership all match.
- Worker do-not additions: never claim independent verification; never rewrite or
  accept Product Knowledge.
- Handoff additions: suggested verifier focus; whether a repair is being performed.

## Acceptance action

Fold into Behavior/Interfaces/Data/Constraints; ground in the cited files.
