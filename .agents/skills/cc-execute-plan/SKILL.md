---
name: cc-execute-plan
description: Execute one explicitly approved numbered plan as a continuous dependency-ordered runtime with cumulative repository review and human closeout.
---

# Execute plan

1. Read the wrapper instructions, workspace configuration, root `plans/`
   roadmap, exact repository collection, approved numbered plan, and repository
   instructions. Use only the exact approved plan reference, version, and digest.
2. Prepare with `node .agents/bin/cc.mjs execute-plan --plan <reference> --version <version> --approved-digest <digest>`. The command derives the complete authoritative graph and creates one branch and isolated worktree per affected repository.
3. Launch fresh workers and independent read-only verifiers only from emitted task inputs. Same-repository tasks remain serialized; independent repositories proceed only when declared dependencies pass. Record every result through `record-result`.
4. Use bounded repair with fresh roles for failed or blocked tasks. After every task passes, run the holistic verifier from the emitted plan-verifier input and record `plan-verifier-result`.
5. Prepare exactly one review handoff per affected repository. Publish only after explicit authorization; never merge or deploy automatically.
6. After human merge, record exact merge confirmation for every repository, then use `finish-work --cleanup --refresh` to fast-forward the configured integration target and preserve the refreshed commit.
7. For an unmerged material revision, explicitly reapprove the plan and resume the same run with `execute-plan --resume-run`. Revision records are append-only; changed tasks and dependency closures are re-executed and stale final/review/merge evidence is invalidated.

`run-task` is not a plan-level step. It remains a rare, manually invoked escape
hatch for one exact task on one exact approved plan and never advances plan
lifecycle, review, or closeout state.
