---
id: 0008-enrich-completion-and-knowledge
target_context_unit: context/domains/completion/README.md
operation: change
statement: >
  Enrich the completion/knowledge page with the shipped verifier-created pending
  evidence state, the impact-status state machine, the reconciliation read-set,
  deferred-proposal future visibility, and the staleness definition + the
  refresh-before-a-consequential-plan rule.
evidence_refs:
  - .agents/skills/cc-complete/SKILL.md
  - wrapper/runtime/engine.sh   # cc_completion_ready, cc_context_impact_record
  - wrapper/contracts/schemas/completion.yaml
  - wrapper/contracts/schemas/context-impact.yaml
  - wrapper/contracts/schemas/context-index.yaml   # freshness/staleness
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich completion with the shipped reconciliation and staleness model

Additions (shipped):
- Verifier-created pending evidence: a successful verifier yields
  execution_status verified / plan_status approved / human_completion pending,
  before the human completion request.
- Impact-status state machine (context-impact schema): not-assessed ->
  review-needed -> accepted / deferred / conflict, plus no-update-needed; only
  review-needed / deferred / conflict get a stored proposal file.
- Reconciliation read-set: final plan/task files + revisions, Product Knowledge
  references + grounding summary, changed paths/commits per repository, worker
  handoffs and verifier evidence, current revisions of relevant context units.
- Deferred proposals stay visible to future plan creation when relevant.
- Staleness: a page is stale when its freshness rule expired, a cited repository
  revision materially changed, or a relevant decision changed. Stale context may
  guide read-only orientation, but a consequential plan surfaces the stale
  reference and refreshes it before execution.

## Acceptance action

Fold into Behavior/Data/Constraints; ground in the completion/context-impact/
context-index schemas.
