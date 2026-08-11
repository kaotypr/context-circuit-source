# Finish-work host compatibility results

Date: 2026-08-11

## Proof shape

Codex CLI 0.147.0 and Claude Code 2.1.220 each ran the same canonical
`finish-work` workflow through its thin adapter in a separate ignored wrapper
copy. Every wrapper contained two clean, independently verified fixture runs:

- a feature commit already merged into the fixture repository's local `main`;
- a deliberately abandoned feature commit with no merge and no remote ref.

The human proof prompt explicitly declared both outcomes, authorized local
commits of generated contribution files, and authorized the cleanup attempts.
It did not authorize source edits, new remotes, pushes, pull requests, merges,
branch deletion, or runtime-evidence pruning.

## Results

| Host | Outcome | Run ID | Product head | Contribution commit | Closeout |
| --- | --- | --- | --- | --- | --- |
| Codex | merged | `20260811T150000Z-c10ce809` | `2f0c12c53037808f55dcc1df0d20bbc5e63ecb98` | `5aed7a5f6a3b5e9fda9eb1f9b43b98dbaa640339` | closed |
| Codex | abandoned | `20260811T151000Z-abadc67e` | `53f8ddecd647cb93874ff4be142fc06705c5baca` | `b418b92147feec81f87be2aafd2fc2a01ad1c0b8` | blocked |
| Claude Code | merged | `20260811T150000Z-c1a0d956` | `29f9da48ee1992929c7b976c29bb546f36415d64` | `574d076123e98a5f0ef5d7206a601fd14a4746cf` | closed |
| Claude Code | abandoned | `20260811T151000Z-abad5450` | `7b308d6c49c621c8f3da971d21f1a1529667ae94` | `f394eaaffb7099d7060daf093800e8a09bb12a35` | blocked |

For each merged run, the host prepared and validated one append-only
contribution, committed only that file in the isolated wrapper, and explicitly
reran closeout with `--cleanup`. The registered worktree was removed, while its
feature branch, runtime manifest, closeout record, worker result, verifier
result, and task brief remained. Both manifest and closeout status became
`closed`, with no blockers.

For each abandoned run, the host prepared, validated, and locally committed one
contribution before requesting cleanup. Cleanup returned the same blocker:

> Abandoned branch contains commits that are neither merged nor preserved by a remote ref.

The manifests remained `closing`, the closeout records became `blocked`, and the
worktrees stayed registered, clean, on their recorded branches and exact heads.

## Independent reconciliation

After both host sessions exited, a deterministic reconciliation checked:

- all four closeout records against the closeout JSON Schema;
- all four runtime manifests against the runtime-manifest JSON Schema;
- all four contribution filenames, required sections, run references, and
  credential patterns;
- wrapper, fixture base, and preserved worktree Git cleanliness;
- exact branch-to-head preservation;
- merged-worktree absence and abandoned-worktree presence;
- wrapper diffs containing exactly two contribution files per host; and
- absence of product remotes.

Every check returned zero errors. No source or workflow file changed inside the
proof wrappers.

## Host observations

- Codex streamed its inspection and closeout progress. Its workspace sandbox
  requested an additional reviewed permission for the local contribution commit,
  then continued without broadening the workflow authorization.
- Claude Code print mode buffered output until the session completed, matching
  the behavior observed in earlier compatibility proofs.
- Codex supplied the local merge commit explicitly. Claude Code omitted it and
  exercised the alternative proof that the product head was already reachable
  from `main`. Both paths produced the same closed state.
- Both hosts independently revalidated generated contributions and closeout
  records before committing them. No provider-specific contract field or
  workflow branch was needed.
- With `activity.provider: none`, both recorded `task.completed` or
  `task.cancelled` as skipped rather than fabricating an external update.

## Outcome

Human-invoked, activity-free closeout now satisfies the two-host compatibility
proof for both safe cleanup and interrupted-closeout preservation. Ignored proof
state remains under `.runtime/host-closeout-proofs/20260811-increment/`; this
document is the durable result summary.
