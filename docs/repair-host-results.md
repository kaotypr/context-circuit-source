# Two-host repair and review-preparation results

Date: 2026-08-11

## Proof task

Codex desktop 0.147.0 and Claude Code 2.1.220 independently ran the same
planless request in separate ignored wrapper copies: add a Reset button that
sets the displayed counter to zero. Each run used one fixture repository, an
isolated branch and Git worktree, fresh worker and verifier sessions, no
activity tool, and no configured Git remote.

| Host | Run ID | Initial commit | Repair commit | Final verdict |
| --- | --- | --- | --- | --- |
| Codex desktop | `20260811T134035Z-04b0cd13` | `240903f9fbe5a6c765dee7b72a7c29f9bb60c1e9` | `1dc0eb63da2a0ef7b9a5a9e9d007bb7e23e62972` | pass |
| Claude Code | `20260811T134035Z-b36ccd13` | `14b4e3a00a37a152dfecde3044723f7ea9c5010e` | `5de28fd371ae9439514ec5be403fe258705c70c0` | pass |

## Controlled failure

The initial workers deliberately implemented the same bounded defect: they
renamed the existing Increment button to Reset without changing its increment
handler, and weakened the required test to check only the label. This was fault
injection for the proof, not normal workflow behavior.

Both initial worker results were schema-valid and both repository check suites
were green. Fresh verifiers nevertheless rejected the work because clicking
Reset incremented the counter and the required test did not exercise the
acceptance criterion. Recording those results advanced both manifests to
`failed` while preserving the commits, worktrees, inputs, and findings.

## Bounded repair

`prepare-repair` generated attempt 1 of the configured maximum of 2 for each
failed run. Fresh repair workers received only the scoped task, previous
findings, repository state, and unique result path. Both restored a distinct
Increment control, added Reset with `setCount(0)`, and added behavioral coverage
that starts from a nonzero count.

Fresh post-repair verifiers independently inspected the full base-to-head
history and returned `pass`:

- Codex passed 2 tests, the TypeScript/Vite build, and a focused Reset behavior
  run. Only `src/App.tsx` and `src/App.test.tsx` changed and the worktree stayed
  clean.
- Claude Code passed 3 tests and the build, then used independent behavior and
  mutation probes to show that the test suite fails against both the original
  defect and an incorrect `setCount(1)` handler. Its worktree also stayed clean.

Each manifest recorded one repair attempt and reached `passed` after six
execution events. Review preparation added a seventh append-only event without
changing the passed run status.

## Review-preparation safety boundary

`prepare-review` generated a schema-valid draft pull-request handoff for each
passed run. Both handoffs are `blocked` because the fixture repositories have no
`origin` remote. The blocker is explicit and the handoffs contain the complete
two-commit history, changed-file set, review title/body, and worker/verifier
evidence paths.

No branch was pushed, no pull request was opened, no merge occurred, and no
remote or credential-bearing URL was added to evidence.

## Host observations

- The generated worker, verifier, and repair inputs remained host-neutral; no
  provider-specific contract field was needed.
- Codex streamed detailed progress and followed the repository instruction to
  reread the complete permanent workflow and context authorities before verification.
- Claude Code print mode buffered its final response, so runtime artifacts were
  the dependable progress surface.
- Fresh Git worktrees do not inherit `node_modules`. Both hosts ran the
  repository-defined dependency installation before checks; generated
  dependencies remained ignored and uncommitted.

## Outcome

The same canonical lifecycle now proves failure detection, bounded repair,
fresh independent reverification, evidence preservation, and safe draft-review
preparation in both supported hosts. The proof runtime remains ignored under
`.runtime/host-repair-proofs/`; this document is the durable result summary.
