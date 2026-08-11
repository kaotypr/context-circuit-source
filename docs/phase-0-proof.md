# Phase 0 manual proof

## Prepare the fixture

Install wrapper dependencies and create a real ignored clone from the tracked
fixture template:

```bash
npm install
npm run fixture:create
npm run check
```

Prepare a planless run:

```bash
npm run run-task -- \
  --request "Add a reset button to the counter" \
  --repository frontend \
  --scope "src/App.tsx" \
  --test-scope "src/App.test.tsx" \
  --acceptance "A Reset button sets the displayed count to zero" \
  --verify "npm test" \
  --verify "npm run build"
```

Record the printed worker-input and verifier-input paths. The activity warning is
expected and is also recorded in the manifest. Do not reuse a prepared run for
both roles.

Before launching the worker, record `worker-started`:

```bash
npm run record-result -- --run-id <run-id> --repository frontend --stage worker-started
```

## Codex proof

1. Invoke `$run-task` or ask Codex to follow `.codex/skills/run-task/SKILL.md`.
2. Give a fresh worker sub-agent only the generated worker-input JSON. Confirm it
   edits only the emitted worktree, runs fixture checks, commits, and writes a
   schema-valid worker result. Confirm the implementation and its required test
   both appear in the commit.
3. Record `worker-result`, then give a different fresh verifier sub-agent only
   verifier-input JSON. Confirm it remains read-only, inspects the base-to-branch
   diff, runs verification, and writes a schema-valid verifier result.
4. Record `verifier-result` and confirm the manifest status matches the verdict.

## Claude Code proof

1. Invoke `/run-task`, which delegates to the canonical skill.
2. Start a fresh Claude Code worker session with only worker-input JSON and its
   referenced files. Require the same commit and worker-result behavior.
3. Record `worker-result`, then start a separate Claude Code verifier session
   with only verifier-input JSON.
   Require read-only behavior and a verifier-result file.
4. Record `verifier-result` and confirm the manifest status matches the verdict.

The proof passes when both hosts use the same generated artifacts and schemas,
the worker and verifier are fresh and isolated, and no host-specific workflow
logic is added. Launch mechanics are intentionally manual in Phase 0 because
host session behavior is the architecture risk being measured.

The first completed compatibility proof and its observed adapter behavior are
recorded in `phase-0-host-results.md`.

## Safe recovery

If preparation stops, inspect `.runtime/runs/<run-id>/manifest.json`. Never delete
a worktree with changes or commits that have not been recorded elsewhere. This
slice does not implement `finish-work`; cleanup is manual and requires checking
`git status`, branch commits, and worktree registration first.
