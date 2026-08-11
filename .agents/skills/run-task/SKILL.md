---
name: run-task
description: Prepare and coordinate an explicitly requested, planless, single-repository task in an isolated Git worktree. Use when a human asks to execute known work without requiring a plan or activity tool and worker/verifier sessions need scoped runtime inputs.
---

# Run task

1. Read `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, and relevant repository instructions.
2. Require a concrete requested outcome, at least one acceptance criterion, an implementation scope, an explicit test expectation, and one registered repository.
3. Run `npm run run-task -- --request <text> --repository <name> --acceptance <criterion> --scope <path> --test-scope <test-path>`; repeat scope or verification options as needed. Omit test scope only when verifier-only evidence is intentional, or select another explicit `--test-policy`.
4. Stop on dirty-repository or validation failure. Never reset, stash, clean, or delete work to bypass protection.
5. Read the emitted manifest and worker input. Record `worker-started` with `npm run record-result -- --run-id <id> --repository <name> --stage worker-started`, then use the current host's adapter to launch a fresh worker only in the emitted worktree.
6. Require the worker to follow the task's test expectation and save its result at the specified path, then record `worker-result` with the same command. Stop if validation, scope, required-test, identity, commit, or worktree checks fail.
7. Launch a different fresh, read-only verifier with the emitted verifier input. Do not let it repair findings.
8. Require the verifier to save its result at the specified path, then record `verifier-result`. Report the resulting manifest status and evidence to the human.

Do not push, open a pull request, merge, deploy, mutate activity state, or clean runtime evidence unless the human separately authorizes the corresponding workflow.
