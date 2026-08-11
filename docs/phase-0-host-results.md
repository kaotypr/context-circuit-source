# Phase 0 host compatibility results

Date: 2026-08-11

## Proof task

Both hosts independently executed the same direct request against separate
branches and worktrees: add a Decrement button whose click reduces the displayed
counter by one. The task was planless and used no activity tool.

| Host | Run ID | Worker commit | Worker | Verifier |
| --- | --- | --- | --- | --- |
| Codex desktop | `20260811T105449Z-4e257800` | `d9b9b7963a052b79d0b7aa63b51d0248d3781ede` | completed | pass |
| Claude Code 2.1.220 | `20260811T110456Z-61a17800` | `5fc846ab8522e1704c9bf4a91cf148c528fceb94` | completed | pass |

Both workers produced the same three-line `src/App.tsx` change. Both worktrees
remained clean after commit and verification. All four result artifacts passed
the same host-neutral JSON Schema validation.

## Verification evidence

- Both workers ran `npm test` and `npm run build` successfully.
- Both verifiers independently inspected the recorded base-to-branch diff and
  reran the required checks without changing repository code.
- Codex verified the acceptance criterion through the locally served production
  build: one click changed `Count: 0` to `Count: -1` with no console errors.
- Claude Code used an out-of-worktree jsdom harness against the actual component,
  including a negative control, to verify exact decrement behavior.
- The two host branches have no `src/App.tsx` diff between them.

Runtime evidence is intentionally preserved under `.runtime/results/` and the
registered worktrees. It remains ignored and is not durable project knowledge.

## Required-test enforcement proof

A follow-up Codex run proved the explicit test-scope contract end to end:

| Run ID | Work ID | Worker commit | Required test path | Verdict |
| --- | --- | --- | --- | --- |
| `20260811T115954Z-1cf78501` | `ADHOC-20260811-004` | `3e351cf219979895b39ec414b0f23cca82bd6603` | `src/App.test.tsx` | pass |

The fresh worker changed only `src/App.tsx` and the required
`src/App.test.tsx`, then passed two tests and the production build. The result
recorder accepted the worker result only after confirming that both the
implementation scope and required test scope appeared in the clean committed
diff. A different fresh verifier independently reran both commands, confirmed
the reset behavior and required test change, and returned `pass`. The manifest
advanced through `prepared`, `running`, `verifying`, and `passed` with three
idempotent execution events.

## Observed adapter behavior

### Codex

- A sub-agent launched without inherited coordinator turns received sufficient
  context from the generated input and referenced instruction paths.
- Separate worker and verifier agents could access the emitted worktree and write
  only their designated runtime results.
- Compact result return to the coordinator worked without workflow-specific code
  in the Codex adapter.

### Claude Code

- Non-interactive `claude -p` with `--no-session-persistence` provided a fresh
  worker or verifier session.
- `--add-dir` was required because result contracts and output paths live in the
  wrapper while the session runs from the repository worktree.
- The positional prompt must precede variadic `--add-dir`; otherwise the option
  consumes the prompt and Claude Code exits before starting.
- Print-mode output was buffered until completion, so filesystem/runtime evidence
  is the reliable progress surface for this manual adapter.
- The verifier did not receive the Edit tool. It still needed Write permission
  for its designated result artifact and Bash permission for verification.

## Contract findings

- The host-neutral task, worker-result, verifier-result, and runtime contracts
  held across both hosts; no provider-specific field was needed.
- The original task brief allowed only `src/App.tsx`, while durable regression
  coverage belongs in `src/App.test.tsx`. Both initial verifiers therefore had to
  supply external acceptance evidence. Explicit implementation scope and test
  expectation policies now solve that gap: required test paths expand and are
  enforced as worker edit scope, omitted test scope is explicitly verifier-only,
  and the follow-up required-test proof passed end to end.
- The original proof exposed that manifests remained `prepared` after manual
  workers and verifiers. The next increment added a deterministic, validated,
  idempotent result recorder; the preserved proof manifests were subsequently
  advanced to `passed` without putting lifecycle logic in host adapters.
- Automated launch, repair, pull-request preparation, and safe closeout remain
  outside the Phase 0 foundation.

## Outcome

The Phase 0 host-neutrality exit criterion is satisfied for the single-repository
fixture proof: a known planless task was delegated to fresh repository workers
and independently verified in Codex and Claude Code without duplicating canonical
workflow logic.
