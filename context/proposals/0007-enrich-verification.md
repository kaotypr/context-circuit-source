---
id: 0007-enrich-verification
target_context_unit: context/domains/verification/README.md
operation: change
statement: >
  Enrich the verification page with the shipped per-ID verifier report structure,
  the bounded-additional-checks allowance, and the verifier read-only
  host-capability-or-block rule. Evidence layers stay OUT of the page: they are a
  design requirement not implemented in the shipped runtime (logged in 0013),
  and the accepted policy is shipped-is-truth.
evidence_refs:
  - .agents/skills/cc-verify/SKILL.md
  - agents/verifier.md
  - wrapper/adapters/AGENTS.md   # verifier read-only capability or block
  - wrapper/contracts/schemas/verifier-result.yaml
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich verification with shipped report structure and read-only-or-block

Additions (shipped):
- Per-ID report contents: for each acceptance/verification id, the observed
  evidence, passed/failed/blocked outcome, the exact failure and a repair
  recommendation, and the repository and commit checked.
- The verifier runs each canonical verification command and bounded additional
  checks needed to establish the declared acceptance.
- Verifier read-only host-capability rule: the host adapter must provide actual
  read-only capability; if it cannot guarantee it, the adapter reports blocked.
  A prompt saying "do not edit" is not enough.

Explicitly NOT added: "evidence layers". The design (05 §10-11) requires an
evidence layer per acceptance criterion, but the shipped verifier does not
implement it. Recorded as a design-vs-implementation delta in `0013`; the page
keeps describing the shipped verifier.

## Acceptance action

Fold the report structure and read-only-or-block rule into Behavior; keep the
provenance note that evidence layers are a design requirement not yet shipped,
now cross-referencing the delta log.
