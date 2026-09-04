# What v1.0 preserves — and why each piece is load-bearing

This file is what separates this study from a ground-up rewrite. It is written
from a full read of `wrapper/runtime/engine.sh` (1,975 lines) and
`wrapper/contracts/invariants.yaml`. Each item below is **correctness that already
exists, is subtle, and would be at risk in a rewrite**. v1.0 keeps every one of
them unchanged and builds on top.

## The engine's boundary is already right — keep it

The runtime's own header and INV-RUNTIME-01 declare that it owns deterministic
mechanics and does **not** own provider launch, model prompts, Product Knowledge
interpretation, plan-writing intelligence, routing, confirmation cards, or
automatic pull-request/merge/complete. That is exactly the model-blind kernel
shape a trust core wants. **Preserve it as the invariant it already is.**

## The stable engine-action seam — the single most valuable thing to keep

Skills drive the engine as a CLI — `sh wrapper/runtime/engine.sh <action> <args>`
— and are forbidden from *reading* it (the "invoke, don't read" rule in
`wrapper/adapters/AGENTS.md`). The ~50-verb dispatch table at the bottom of
`engine.sh` is therefore a **stable adapter boundary**. Every v1.0 change is
expressed as new or changed verbs *behind this seam*, so:

- Two-thirds of the periphery ports unchanged (see `engine-and-seam.md`).
- A later behavior-preserving port of the mechanics to a typed language would not
  touch the skills.

**Preserve the seam and its emit format.** It is the thing that makes evolution
safe.

## The subtle mechanics — reproduce exactly, never "clean up" casually

These are the functions where a rewrite would silently lose correctness. v1.0
does not modify them.

- **Path leases with region overlap + descendant exemption** (`cc_lease_check`,
  `cc_region_overlap`, INV-CONCURRENCY-01). Overlap is defined precisely — equal,
  prefix-ancestor either way, or repository-wide `.` — and a declared dependency
  descendant is *exempt* because it builds on the holder rather than competing.
  Held from execution start until **delivery**, not merely verification. This is
  the race arbiter; it is easy to get subtly wrong.
- **Base selection with integration merge and stale-base rebuild**
  (`cc_base_prepare`, INV-CONCURRENCY-02). Base tip when no same-repo
  predecessor; the single predecessor branch when one; a **runtime-authored
  integration merge** when two or more — with abort-and-rollback on conflict
  (`BASE_UNBUILDABLE`, a blocked execution, never a worker failure). A recorded
  base at `refs/cc-base/<plan>/<repo>` is invalidated and rebuilt when a
  predecessor is repaired. This is the hardest logic in the engine.
- **Verifier read-only enforcement is structural, not asked-nicely**
  (`cc_verifier_result_record`, INV-VERIFY-01). It rejects an explicit
  `--wrote-products` flag **and** re-checks that every branch tip is unchanged
  since the worker commit (`VERIFIER_MODIFIED_PRODUCT`). v1.0 keeps this exact
  mechanic and merely *binds its result to a candidate* (M2).
- **Commit-before-verify, new-commit repairs** (`cc_worker_commit_record`,
  INV-EXEC-04). A repair must be a *new* commit (`COMMIT_NO_CHANGE` refuses a
  no-op), never an amendment that conceals an attempt.
- **The three-failure counter and waived≠passed** (`cc_verifier_result_record`,
  INV-REPAIR-01). `waived` and `blocked` never satisfy verification and never
  increment the failure counter; only an independent rejection does. v1.0 keeps
  the counter; a tier may adjust the *budget*, but the honesty rules stay.
- **Preserve-on-failure** (INV-PRESERVE-01, `cc_recovery_inspect`). Failure,
  interruption, or blocking preserves branches, worktrees, commits, handoffs,
  verifier evidence, and records. Nothing is silently cleaned up. Resume
  eligibility is computed from intact snapshot + intact worktrees + owner match.
- **Atomic writes and rollback** (`cc_atomic_write`, and the archive/restore
  functions that *roll the move back* if the index update fails). Partial records
  cannot grant ownership, resume a worker, prove verification, or authorize
  completion (INV-RUNTIME-02).
- **Worktree isolation from the base tip** (`cc_worktree_prepare`,
  INV-EXEC-03). One deterministic branch `cc/<plan>/<repo>` and one worktree per
  repository; the base checkout is never written during execution.
- **Delivery separateness + drift guard** (`cc_delivery_drift`,
  `cc_delivery_rebase`, INV-DELIVER-01/02). Pull request, merge, push, deploy are
  separate human actions; a plan whose recorded base drifted is rebased and
  re-verified before its pull request opens; the target is the recorded
  `base_branch`, never a moving remote. **v1.0 keeps delivery as the second
  gate, untouched.**

## Identity, security, and host neutrality — untouched

- **Portable identity vs. host-local binding** (INV-REPO-01/02/03/04):
  `workspace.yaml` holds credential-free logical identity; `repositories.local.yaml`
  holds machine paths; binding resolution fails closed on missing/ambiguous/
  non-Git/traversal/unsafe-symlink/identity-mismatch without scanning the
  filesystem. The `base_branch` (execution base + PR target) vs. `default_branch`
  (portable clone guidance) distinction is preserved exactly.
- **Credentials never enter workspace or runtime state** (INV-SEC-01); `sources/`
  is passive raw evidence read only by exact request-named files (INV-SEC-02).
- **Host neutrality** (INV-HOST-01): host identity, version, capability,
  permission mode, provider status, and per-role `(model, effort)` tiering are
  bounded provider-neutral `host_evidence` that authorizes nothing. v1.0's tiers
  are a *consequence* concept, orthogonal to and never conflated with this
  *host/model* tiering.
- **Two low-level helpers other subsystems depend on**: `cc_digest` and the
  atomic-write primitive. `cc-publish`'s idempotency (`synced_digest`) and record
  integrity lean on them at the schema level; preserve both or publication's
  guarantee silently changes.

## Publication and grounding — kept as-is

- **Publication is orthogonal** (INV-EXTERNAL-01/02/03): manual, one-way,
  export-only, self-contained, credential-free, writing only under
  `publication/<name>/`, referencing no core phase. v1.0 does not touch it. Its
  one nuance — an inbound flow must pass through the normal planning/approval gates
  — maps cleanly onto v1.0's intent gate.
- **Repository grounding** (INV-GROUND-01/02/03): the worker reads the target
  repo's own `AGENTS.md`/`CLAUDE.md`/rules/skills via a deterministic runtime scan
  that emits a data manifest (never prompt text); the worker brief is a fixed
  template filled by slot substitution with a required, preflighted grounding
  section. v1.0 keeps this whole subsystem; it composes with M1 (the intent's
  outcome-level criteria and discovery's grounded findings) rather than replacing
  it — the worker's execution-time read now confirms what discovery already found,
  rather than being the code's first contact.

## The principle

Everything above is *how to do the work safely*. v1.0 changes *what the human
decides and when* — the policy layer — and leaves the safety mechanics exactly
where a full read of the code shows they are already correct. That is the
difference between an evolution and a rewrite: the rewrite would have to
rediscover every item on this page.
