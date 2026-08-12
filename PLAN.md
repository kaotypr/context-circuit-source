# Context Circuit workflow improvement plan

Status: Draft  
Branch: `codex/improve-context-circuit-workflow`  
Source: Review of Codex task `019ff2f6-8807-7752-9ee8-d2a51315a3d9`

> Temporary planning artifact: delete this file after every implementation and
> verification item below is complete, and before the branch is merged to
> `main`. The final durable outcome belongs in an append-only contribution and,
> when appropriate, curated `context/` updates.

## Outcome

Make Context Circuit operate as one coherent journey from workspace setup and
planning through implementation, review, human merge, closeout, and the next
recommendation. An approved plan must remain linked to its execution and
completion evidence so the wrapper never recommends already completed work.

The implementation must preserve the existing safety properties: clean base
repositories, isolated worktrees, independently scoped workers and verifiers,
explicit publication authorization, a human merge gate, recoverable closeout,
and no implicit deletion of user work or runtime evidence.

## Work plan

### 1. Establish plan-linked execution and outcome contracts

- Add a first-class input path for executing one approved plan item or an
  explicitly selected dependency-ordered set of plan items.
- Validate the plan status, version, approval digest, work IDs, dependencies,
  affected repositories, scope, test expectations, and acceptance evidence
  before preparing a run.
- Populate task briefs with `source.kind: plan`, the plan reference, approval
  state, and stable plan work IDs instead of normalizing approved-plan work as
  `ADHOC-*` direct requests.
- Define how one run records outcomes for multiple selected plan items without
  fabricating completion for later dependencies.
- Keep the existing direct-request path unchanged for genuinely planless work.

Likely areas:

- `.agents/contracts/run-task-request.schema.json`
- `.agents/contracts/task-brief.schema.json`
- `.agents/contracts/runtime-manifest.schema.json`
- `scripts/run-task.ts`
- `scripts/lib/run-task.ts`
- `scripts/lib/types.ts`
- `.agents/skills/run-task/SKILL.md`
- `docs/run-task.md`

Acceptance:

- A selected approved item produces a plan-linked task brief using its stable
  work ID and approval evidence.
- A direct request still produces an `ADHOC-*` task brief.
- Draft, stale, unknown, or dependency-blocked plan inputs fail before branch or
  worktree creation.

### 2. Reconcile `$whats-next` with execution and closeout evidence

- Introduce a read-only work-state projection that combines configured activity
  facts, active runtime manifests, verification/review state, closeout records,
  and durable completed-work evidence.
- Define explicit precedence and contradiction handling. Conflicting evidence
  must yield a reconciliation action rather than duplicate implementation.
- Exclude completed and cancelled plan items and rank active, failed,
  review-ready, or cleanup-ready work appropriately.
- Preserve useful behavior when no activity provider is configured.
- Avoid placing mutable live status inside approved plan documents.

Likely areas:

- `scripts/lib/whats-next.ts`
- `scripts/whats-next.ts`
- `.agents/contracts/work-candidate.schema.json`
- `.agents/contracts/whats-next-result.schema.json`
- `.agents/skills/whats-next/SKILL.md`
- `docs/whats-next.md`

Acceptance:

- A passed or closed plan-linked run is not recommended for execution again.
- A review-ready run is recommended as review/merge work, not implementation.
- Contradictory activity and local evidence is reported explicitly and does not
  mutate either source.
- Recommendations remain deterministic and read-only.

### 3. Separate repository identity from descriptive work area

- Add an explicit repository key to plan work items; retain `area` for a human
  description if it remains useful.
- Validate repository keys against `workspace.yaml` during plan creation,
  validation, approval, publication, execution, and next-work discovery.
- Provide a compatibility migration or clear validation error for existing plans
  whose `area` currently contains the repository key.
- Never silently guess a repository from free-form area text.

Likely areas:

- `.agents/contracts/plan-draft-request.schema.json`
- `.agents/contracts/plan-work-breakdown.schema.json`
- `scripts/lib/plans.ts`
- `scripts/lib/whats-next.ts`
- plan publication contracts and code
- planning documentation and examples

Acceptance:

- Descriptive values such as `application foundation` no longer break repository
  resolution.
- Unknown repository keys fail plan validation before approval.
- Existing valid plans have a documented, tested compatibility path.

### 4. Create a continuous review, publication, merge, and closeout handoff

- Distinguish `ready-for-local-review`, `ready-for-publication`,
  `published-for-review`, `merge-confirmation-required`, and closeout-ready
  states rather than treating a missing remote as a generic blocker.
- For local review, emit exact diff, commit inspection, test, and human-run merge
  commands using the recorded source and target branches.
- For remote review, retain explicit authorization before pushing or opening a
  draft PR/MR and record only confirmed publication results.
- After the human reports a merge, verify reachability and give a direct,
  unambiguous handoff to `finish-work`.
- Keep merge execution as a human gate.

Likely areas:

- `.agents/contracts/review-preparation.schema.json`
- `.agents/contracts/review-publication-record.schema.json`
- `scripts/lib/review-lifecycle.ts`
- `scripts/prepare-review.ts`
- `.agents/skills/run-task/SKILL.md`
- `.agents/skills/finish-work/SKILL.md`
- `docs/review-lifecycle.md`
- `docs/finish-work.md`

Acceptance:

- A repository without `origin` can be ready for local review.
- A repository with a remote can be published only after explicit authorization.
- Exact merge evidence is verified before a merged closeout can proceed.
- No agent path merges, deploys, or deletes the feature branch.

### 5. Replace recurring initialization UX with workspace configuration

- Add `$configure-workspace` as the primary workflow for both fresh and existing
  wrappers.
- Keep bootstrap as an explicit internal phase for initial Git and empty base
  commits.
- Retain `$initialize-workspace` as a compatibility alias that routes based on
  detected state.
- Ask about wrapper and product remotes, review mode, authoritative PRD/context
  sources, repository roles, activity behavior, and lifecycle policies.
- Add an explicit upgrade boundary for future schema/template migrations rather
  than treating upgrades as ordinary configuration.

Acceptance:

- Fresh configuration remains idempotent and requires exact initial-commit
  authorization.
- Reconfiguration of a team wrapper produces reviewable changes and never
  commits directly to `main`.
- Existing active runs retain their captured configuration.
- The old initialization invocation remains compatible and clearly routed.

### 6. Generate project-facing wrapper documentation and context intake

- Generate or reconcile the wrapper README so it leads with workspace identity,
  purpose, registered repositories and roles, common next actions, and links to
  durable context and plans.
- Keep Context Circuit framework documentation as a secondary section.
- During initial configuration, discover authoritative PRDs, architecture
  documents, issues, or repository documentation and record stable source
  references.
- Preserve the boundary that `gather-context` is read-only and `sync-context`
  owns reviewable durable-context writes; consider a higher-level refresh flow
  that orchestrates both.

Acceptance:

- A configured wrapper README describes the actual project and product paths.
- Unknown architecture, conventions, or decisions remain explicitly unknown.
- External content is source-cited and cannot override workspace instructions.

### 7. Improve cleanup, retention, lifecycle extensibility, and packaging

- Detect cleanup blockers early and return an ordered closeout checklist.
- Recommend the exact wrapper commit/review step and cleanup rerun when closeout
  evidence is not yet durable.
- Add a separate, conservative retention/pruning policy for closed runtime
  evidence; never overload normal closeout with wholesale deletion.
- Document and test the distinction between actions on existing lifecycle hooks,
  configurable human gates, and versioned new workflow stages.
- Resolve the missing thin Codex adapter for `gather-context` or stop advertising
  the absent adapter path.

Acceptance:

- Dirty or unrecorded work remains preserved with precise blockers.
- Successful cleanup removes only the registered clean worktree.
- Runtime pruning is explicit, bounded, and recoverable where practical.
- Every advertised skill path exists in the built template and installed skill
  catalog.

## Verification strategy

- Add focused unit tests for every new contract and state transition.
- Add an end-to-end fixture that initializes a wrapper, creates and approves a
  plan, executes a plan item, records worker and verifier results, prepares local
  review without a remote, verifies a human merge, closes the run, and confirms
  `$whats-next` does not recommend the completed item.
- Add a parallel fixture for remote review preparation/publication records
  without performing a real push.
- Cover compatibility for existing plans and `$initialize-workspace` callers.
- Run:
  - `npm test`
  - `npm run build`
  - `npm run validate -- --check-paths --check-documents`
  - `git diff --check`
  - release archive/package validation

## Delivery order

1. Contracts and plan-linked run normalization.
2. Work-state projection and `$whats-next` reconciliation.
3. Explicit repository mapping and plan compatibility.
4. Review/merge/closeout states and handoffs.
5. Workspace configuration workflow and compatibility alias.
6. Project-facing README and context-source intake.
7. Cleanup retention, lifecycle documentation, and skill packaging.
8. Full regression, host-adapter, and release-artifact verification.
9. Record the durable contribution/context updates.
10. Delete `PLAN.md`, rerun the full verification suite, and only then present
    the branch for merge into `main`.

## Non-goals

- Automatic merge or deployment.
- Silent publication, remote creation, or activity mutation.
- A central task database, background service, or workflow UI.
- Copying external documents directly into canonical context without review.
- Deleting dirty, unpushed, failed, blocked, or otherwise unrecorded work.

