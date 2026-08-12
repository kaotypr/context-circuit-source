# Workflow continuity delivery record

Durable, append-only record of the workflow-improvement work delivered on
`codex/improve-context-circuit-workflow`. This is maintainer evidence for the
Context Circuit source repository; it is intentionally excluded from the shipped
wrapper template. The temporary `PLAN.md` that scoped this work was removed at
closeout; this record is its durable replacement.

## Outcome

Context Circuit now operates as one coherent journey from workspace setup and
planning through implementation, review, human merge, closeout, and the next
recommendation. An approved plan stays linked to its execution and completion
evidence, so the wrapper never recommends already-completed work. All existing
safety properties are preserved: clean base repositories, isolated worktrees,
independently scoped workers and read-only verifiers, explicit publication
authorization, a human merge gate, recoverable closeout, and no implicit
deletion of user work or runtime evidence.

## Delivered areas

1. **Plan-linked execution and outcome contracts.** Approved plan items execute
   under their stable plan work IDs and approval digest instead of being
   normalized to `ADHOC-*`. Draft, stale, unknown, or dependency-blocked plan
   inputs fail before any branch or worktree is created. Direct planless
   requests still produce `ADHOC-*` briefs unchanged.
   (`e6ef72c`, `5fcf2af`, `ca89940`, `d183a10`)
2. **`whats-next` reconciliation.** A read-only work-state projection combines
   configured activity facts, validated plan-linked runtime manifests, closeout
   records, and durable completed-work contributions. Durable outcomes take
   precedence over closeout, runtime, and activity evidence, but contradiction
   yields an explicit reconciliation action rather than silently choosing a
   state or mutating any source. Completed and cancelled work is excluded.
   (`bc646b5`, `6a4ff49`)
3. **Repository identity separated from descriptive area.** Plan work items carry
   an explicit repository key validated against `workspace.yaml`; `area` remains
   a human description only. Unknown repository keys fail plan validation before
   approval; descriptive area text never resolves a repository.
   (`54773c4`)
4. **Continuous review, publication, merge, and closeout handoff.** Distinct
   `ready-for-local-review`, `ready-for-publication`, `published-for-review`,
   `merge-confirmation-required`, and closeout-ready states replace a generic
   missing-remote blocker. Local review emits exact diff, commit-inspection,
   test, and human-run merge commands. Remote publication requires explicit
   authorization and records only confirmed results. `confirm-merge` is
   read-only with respect to Git and requires the exact reviewed head and
   reported merge to be reachable from the configured default target before
   closeout is ready. No agent path merges, deploys, or deletes the branch.
   (`6c4a320`, `74e98e2`)
5. **Workspace configuration and recoverable transactions.** `configure-workspace`
   is the primary flow for fresh and existing wrappers, with bootstrap retained
   as an explicit internal phase and `initialize-workspace` kept as a routed
   compatibility alias. Fresh configuration is idempotent and requires exact
   initial-commit authorization; team reconfiguration produces reviewable
   changes and never commits directly to `main`. Cross-file configuration
   updates are applied as a recoverable transaction, and authoritative
   context sources are captured as cited, non-overriding references.
   (`29ef0d4`, `cfaf00e`, `457c215`)
6. **Project-facing documentation, context intake, retention, and packaging.**
   The wrapper README leads with workspace identity, registered repositories,
   and common next actions, with framework documentation kept secondary.
   Authoritative PRD/architecture/context sources are recorded as stable,
   source-cited references that cannot override workspace instructions. Cleanup
   returns ordered closeout blockers, retention/pruning stays a separate bounded
   operation, lifecycle-hook versus human-gate versus new-stage extensibility is
   documented, and every advertised skill and host adapter exists in source and
   in the built template.
   (`c8eb18c`, `2a2d564`, `13a74ce`)

## Affected areas

- `.agents/` contracts, skills, and bundled command; `scripts/` and
  `scripts/lib/` workflow logic; `docs/` guides; `.codex/` and `.claude/` thin
  adapters; `README.md`, `WORKFLOW.md`, `workspace.yaml`, and `context/SOURCES.md`.
- 16 commits ahead of `main` at `13a74ce`.

## Verification

- `node --import tsx --test --test-concurrency=1 test/*.test.ts` — 137 pass, 0 fail.
- `npm run typecheck` (`tsc --noEmit`) — clean.
- `npm run validate -- --check-paths --check-documents` — valid workspace.
- `git diff --check main...HEAD` — no whitespace or conflict markers (16-commit diff).
- `npm run build:template` — emits version-matched archive, staging directory,
  and bundled `.agents/bin/cc.mjs`.
- `node --import tsx --test --test-concurrency=1 test/packaging.test.ts test/release-artifact.test.ts` — 9 pass, 0 fail.
- End-to-end journey (`test/end-to-end.test.ts`) drives configure → approve plan
  → execute plan item → record worker and verifier results → local review →
  human merge → verified `confirm-merge` → safe closeout → `whats-next` excludes
  the completed item. The verifier helper receives only the workspace root and
  the emitted verifier-input path; it validates that input and the referenced
  task brief before any inspection, derives HEAD independently from Git and base
  from the verifier input, runs every configured verification command, enforces
  the exact `base...HEAD` scope, keeps worker and verifier identities distinct,
  and rejects a tampered verifier input before any result or manifest is written.
- An independent read-only verification pass reproduced all of the above in the
  isolated worktree with no fixes required.

## Decisions and deviations

- The durable completed-work record for framework development lives in this
  maintainer `docs/` file rather than `contributions/general/`, because the
  latter is copied whole into the shipped wrapper template and its file
  inventory is baked into the bundle. Wrapper users' own runs still use
  `contributions/general/` and the `finish-work` append-only flow unchanged.

## Remaining risks and follow-up

- **Configuration crash consistency is bounded, not total.** Caught
  configuration errors, permission failures, and injected rename failures roll
  back fully. A hard `SIGTERM` or process crash during the rename phase can
  leave recoverable sibling `.stage` and `.backup` artifacts. The next
  configuration run detects this residue and refuses with manual recovery
  guidance. This is an accepted limitation: do not claim full crash consistency.

## Candidate durable learnings

- Recorded in `context/DECISIONS.md` and `context/PROJECT.md` as part of this
  closeout. No further canonical-context change is asserted automatically.
