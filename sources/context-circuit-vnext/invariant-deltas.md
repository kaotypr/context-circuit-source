# Invariant deltas

Exact changes to `wrapper/contracts/invariants.yaml`, keeping the "one rule, one
owner" discipline. Grouped as NEW / REWORKED / KEPT-VERBATIM. IDs below reuse the
real current ids where a rule is reworked.

## NEW invariants

- **INV-INTENT-01 — the correctness gate.** Every writing request is anchored to a
  first-class `intent/<id>/` whose `contract.yaml` holds goal, non-goals,
  constraints, acceptance criteria (each executable or explicitly manual), scope,
  and tier. An independent **spec adversary** challenges the acceptance criteria
  before approval. Approval is the single upstream human gate and freezes
  `contract_digest`. Owner: `wrapper/contracts/schemas/intent-contract.yaml`
  (new) + this file.
- **INV-INTENT-02 — the scope envelope.** A plan declares its parent intent and may
  touch only repositories and path regions within that intent's declared `scope`.
  A plan (or candidate) that exceeds the envelope, or relies on a changed
  `contract_digest`, is held and re-gated to a human; it never proceeds on the
  original approval. The envelope check fails upward on any ambiguity. Owner: this
  file + `wrapper/contracts/schemas/plan.yaml` (adds required `intent`).
- **INV-CANDIDATE-01 — candidate-bound evidence.** A candidate is a deterministic
  digest over the per-repository commit map, the selected base commits, and the
  intent's `contract_digest`. Independent-check results and human acceptance bind
  to a candidate; any new commit or criteria change yields a new candidate and
  voids all prior evidence and acceptance for the old one. No role decides whether
  old evidence is "close enough." Owner:
  `wrapper/contracts/schemas/candidate.yaml` (new) + `verifier-result.yaml`.
- **INV-ASSURE-01 — tiered assurance.** Assurance is composed from parts and
  summarized by tier (Explore / Standard / Critical), declared on the intent from
  transparent risk signals and raisable by the human. Explore is human-supervised
  with no independent verifier and is never labeled "verified"; Standard and
  Critical require an independent read-only verifier bound to the current
  candidate. The runtime is model-blind: it enforces the verifier floor per tier
  but never selects a tier (INV-RUNTIME-01, INV-HOST-01). Tiering fails upward.
  Owner: this file (policy) + `wrapper/contracts/schemas/intent-contract.yaml`.

## REWORKED invariants (same id, changed rule)

- **INV-APPROVE-01** — Approval is the explicit conversational human gate on the
  **intent** (draft→approved on the contract), not on each plan. It follows the
  spec-adversary pass and freezes `contract_digest`. No confirmation card, no
  hidden token. Plan readiness within an approved envelope is automatic
  (INV-INTENT-02), not a second human gate.
- **INV-EXEC-01** — A plan derived from an approved intent and confirmed within its
  envelope may execute; execution authorization comes from the intent's approval,
  not a separate per-plan approval. (Combining a request that approves intent and
  asks to execute in one turn remains allowed.)
- **INV-PLAN-01** — `plan.yaml` names its parent `intent` and carries a status that
  is a **projection**: `approved` follows automatically from envelope confirmation,
  `done` is inferred from candidate acceptance + delivery (explicit at Critical).
  Task status remains a synchronized projection, never a second authority.
- **INV-VERIFY-01** — At Standard and Critical, one independent verifier checks the
  current **candidate** of every affected repository, strictly read-only, and its
  result binds to that candidate. At Explore there is no independent verifier and
  the result is never "verified." It never repairs, modifies product files, changes
  status, or self-verifies.
- **INV-COMPLETE-01** — Completion is **inferred** from candidate acceptance +
  delivery at Explore/Standard; at Critical an explicit human completion is
  required. Verification alone never completes anything, at any tier.
- **INV-COMPLETE-02** — Completion (or delivery) records the accepted candidate and
  emits a **reconciliation-debt marker** that starts, and blocks on, Product
  Knowledge reconciliation for the affected scope.
- **INV-KNOWLEDGE-02** — Knowledge acceptance is still an explicit, separate human
  decision (never auto-accepted by completion, verification, or code changes). In
  addition, a new plan's grounding **blocks** (Standard/Critical) or **loudly
  warns** (Explore) while delivered work in its knowledge scope remains
  unreconciled. The gate refuses to let debt be forgotten; it never auto-accepts.
- **INV-PAIR-01** — Retained in full, but **relocated as the Explore tier of
  INV-ASSURE-01** rather than a standalone mode. All of its content stands (one
  repository, fresh isolated branch/worktree, human oracle, never verified, commit
  on request, drift-blocks, light pointer, host-blocked fallback), plus an explicit
  **promote** step: attaching an intent and raising the tier turns a session into a
  candidate-bearing change without restart. Owner remains the pairing skill,
  cross-referenced from INV-ASSURE-01.

## KEPT VERBATIM (the mechanics and safety spine — unchanged)

INV-KNOWLEDGE-01 (retrieval-first knowledge), INV-PLAN-02/03/04/05 (task scope,
stable ids, grounding, inter-plan dependencies), INV-EXEC-02/03/04 (one worker,
deterministic branch/worktree from anchor tip, commit-before-verify + new-commit
repairs), INV-COMMIT-01 (Conventional Commits + no AI attribution), INV-VERIFY-02
(host cannot make an independent verifier → blocked, never self-verify),
INV-REPAIR-01 (three-failure counter; a tier may set the *budget* but the counting
and blocked-vs-failure honesty are unchanged), INV-PRESERVE-01 (preserve on
failure), INV-ARCHIVE-01/02 (archive/restore as status-blind moves), INV-REPO-01/02/
03/04 (portable identity, anchor vs default branch, gitignored repositories/, fail-
closed binding), INV-SEC-01/02 (credentials out of workspace, passive sources),
INV-RUNTIME-01/02 (model-blind runtime, atomic records), INV-OWN-01 (one-worker
lock), INV-CONCURRENCY-01/02 (path leases + base selection/integration merge),
INV-GROUND-01/02/03 (repository grounding + brief assembly), INV-HOST-01 (bounded
host evidence, incl. per-role model/effort tiering — orthogonal to assurance
tiers), INV-SKILL-01 (skills ship at `.agents/skills/<name>/SKILL.md`),
INV-DELIVER-01/02 (delivery separateness + drift guard — Gate 2),
INV-EXTERNAL-01/02/03 (publication orthogonal, export-only, self-contained).

## Owner-map additions

New rows for `owners:` in `invariants.yaml`:

```
intent_contract:     wrapper/contracts/schemas/intent-contract.yaml
intent_gate:         wrapper/contracts/invariants.yaml
scope_envelope:      wrapper/contracts/invariants.yaml
candidate_identity:  wrapper/contracts/schemas/candidate.yaml
assurance_tiering:   wrapper/contracts/invariants.yaml
human_acceptance:    wrapper/contracts/schemas/candidate.yaml
reconciliation_debt: wrapper/contracts/schemas/context-impact.yaml
```
