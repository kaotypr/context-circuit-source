# Context Circuit v0.5 — Implementation Report

Status: implementation of the accepted v0.5 design
Authoritative source design: `sources/system-design/context-circuit/v0.5/`
Scope: complete refactor with no backward compatibility for the old lifecycle,
routing model, context packets, leases, gates, confirmation cards, or
planner/verifier split.

This document describes what was built, where it lives, and how it maps to the
design. It is maintainer material, not workspace Product Knowledge.

---

## 1. Summary

The old implementation was a ~123 KB / 2,523-line sourced shell library whose
central abstraction was a two-stage keyword router (`cc_probe`/`cc_action`/
`cc_route`) plus machinery for context packets and byte budgets, leases,
transaction graphs, confirmation cards, receipts/evidence layers, and
planner/verifier delegation. All of that is old-design and was removed.

v0.5 replaces it with:

- a small, host-neutral, deterministic runtime (`wrapper/runtime/engine.sh`,
  ~44 KB POSIX sh) that owns only workspace/Git/execution-state operations;
- a one-rule-one-owner contract set (`wrapper/contracts/`);
- a blank universal-project-workspace template seed (`template/`);
- seven conversational skills and three thin role deltas;
- a full semantic acceptance suite (`test/`) plus a template-runtime laboratory;
- updated release assembly and CI.

The complete suite (`sh test/acceptance.sh`) passes: 11 behavior suites plus the
AC-01…AC-28 criteria cross-check.

---

## 2. Environment constraints honored

- `/bin/sh` is **dash**. The runtime and all tests are strict POSIX sh — no
  `[[ ]]`, arrays, `local`, or `BASH_SOURCE`. The CLI/self-invoke guard uses
  `case "$0" in *engine.sh) cc_main "$@" ;; esac`.
- Tooling available: `git`, `jq`. **No `yq`, shellcheck, or bats.** The runtime
  parses only the known plan.yaml field shape with `awk`; it is not a general
  YAML parser.
- Tests **source** `engine.sh` (via `test/lib/assert.sh`) and call `cc_*`
  functions directly, and also exercise the thin CLI.

---

## 3. Runtime engine (`wrapper/runtime/engine.sh`)

A single-file POSIX-sh library, sourced by adapters/tests and also usable as a
thin CLI: `sh engine.sh <command> [args]`. Output is deterministic `key: value`
lines with UPPER_SNAKE reason codes on failure; there is no prose, routing, or
cards.

### What it owns

| Area | Functions (illustrative) |
| --- | --- |
| Safe ids / paths / digests / atomic writes | `cc_safe_id`, `cc_safe_slug`, `cc_safe_relative`, `cc_plan_id_valid`, `cc_digest`, `cc_atomic_write` |
| YAML readers (controlled subset) | `cc_scalar`, `cc_list_ids`, `cc_inline_list`, `cc_task_list` (inline **and** block lists), `cc_task_ids`, `cc_plan_repositories`, `cc_plan_affected_repositories`, `cc_plan_repo_paths` |
| Workspace | `cc_workspace_validate`, `cc_workspace_init` |
| Repository bindings & worktrees | `cc_repository_register`, `cc_binding_field`, `cc_repo_resolve`, `cc_repo_anchor_commit`, `cc_repo_clean`, `cc_repository_preflight`, `cc_worktree_prepare` |
| Plan lifecycle & index | `cc_plan_validate`, `cc_plan_allocate_id`, `cc_plan_index_upsert`/`_remove`, `cc_plan_status`, `cc_plan_approve` |
| Archive / restore | `cc_plan_archive`, `cc_plan_restore` (org-lock, exact path + index-collision checks, move rollback on index failure) |
| Locking | `cc_lock_acquire`, `cc_lock_owner`, `cc_lock_release` (one active writer) |
| Execution records | `cc_execution_begin`, `cc_execution_next_id`, `cc_exec_set`, `cc_attempt_begin`, `cc_worker_commit_record`, `cc_worker_handoff_record` |
| Verification & repair | `cc_verifier_prepare`, `cc_verifier_result_record` (read-only enforcement), `cc_repair_allowed`, `cc_execution_status` |
| Completion & knowledge handoff | `cc_completion_ready`, `cc_plan_complete`, `cc_context_impact_record` |
| Delivery (read-only report) | `cc_delivery_targets` |
| Recovery | `cc_recovery_inspect` (snapshot/worktree/owner resume-eligibility) |

### Key deterministic behaviors

- Only an **approved** plan can begin execution; `execution-begin` is the sole
  path that creates `execution.yaml`.
- One branch (`cc/<plan-id>/<repository-id>`) and one worktree per affected
  repository, created from the captured `anchor_branch` tip; the anchor checkout
  is never written.
- Worker commits are captured from each worktree **before** verification; a
  no-op commit (HEAD unchanged since the last recorded revision) is rejected, so
  every repair yields a new commit.
- Verifier is read-only: an explicit product-write flag and any changed branch
  tip since the worker commit are rejected.
- `worker_failures` increments on each verifier rejection (including the first)
  and stops at exactly 3. `blocked` and `waived` are non-passing but are **not**
  worker faults.
- Completion is human-only and requires the latest execution to be `verified`;
  verification alone never sets plan status to `done`.
- Archive/restore never inspect plan or execution status; they perform exact
  path + index-collision checks under a short-lived org-lock and roll back the
  move if the index update fails.
- Records are written atomically; a partial record grants nothing.

### What it deliberately excludes

No provider launch, model prompts, Product Knowledge interpretation,
plan-writing intelligence, conversational routing, confirmation cards, product
test-runner logic, or automatic pull-request/merge/push/publish/deploy/
cleanup/completion.

### Runtime state layout

```
.runtime/
├── executions/<plan-id>/<execution-id>/
│   ├── execution.yaml
│   ├── snapshot/{PLAN.md,plan.yaml,tasks/}   # immutable plan snapshot
│   ├── handoff.md
│   ├── completion.yaml
│   ├── context-impact.yaml
│   ├── attempts/NNN/{worker.yaml,verifier.yaml}
│   └── repositories/<repository-id>.yaml     # worktree, branch, anchor, base,
│                                             # latest_commit, allowed_paths
└── locks/<plan-id>.lock/owner.yaml
```

---

## 4. Contracts (`wrapper/contracts/`)

- `invariants.yaml` — the one-rule-one-owner map. 31 invariants (INV-KNOWLEDGE,
  INV-PLAN, INV-APPROVE, INV-EXEC, INV-VERIFY, INV-REPAIR, INV-PRESERVE,
  INV-COMPLETE, INV-ARCHIVE, INV-REPO, INV-DELIVER, INV-SEC, INV-RUNTIME,
  INV-OWN, INV-HOST) plus an `owners:` concern→file map.
- `manifest.yaml` — runtime version, artifact kind, accepted schema versions,
  runtime exclusions, release boundary, upgrade boundary, host-evidence rules.
- `schemas/` (11 files): `workspace`, `repositories-local`, `plan`, `task`,
  `execution`, `worker-handoff`, `verifier-result`, `completion`,
  `context-impact`, `context-proposal`, `context-index`.

Removed old-design contracts: `routes.yaml`, `context-sets.yaml`, `tier0.yaml`,
and schemas `lease`, `delegation`, `session`, `stack`, `child-start`,
`context-receipt`, `handoff`, `archive`.

---

## 5. Skills, roles, adapters

Skills (`.agents/skills/`): `cc-workspace` (init/orient/connect/clone/init),
`cc-plan` (create + review), `cc-execute` (approve + execute + repair
coordination), `cc-verify` (independent verifier), `cc-complete` (completion +
reconciliation), `cc-archive` (archive/restore), `cc-deliver` (PR/merge
boundary). Removed: `cc-entry`, `cc-gates`, `cc-next`, `cc-upgrade`.

Role deltas (`agents/`): `coordinator.md`, `writer.md`, `verifier.md`. Removed
aliases: `repository-worker.md`, `reviewer.md`.

Adapters (`wrapper/adapters/`, shipped to the workspace root): `AGENTS.md`,
`WORKFLOW.md`, `CLAUDE.md`, `README.md`.

Docs (`docs/`): rewritten `getting-started`, `planning`, `plan-review`,
`product-knowledge`, and `templates/{plan.md,plan.yaml,task.md}`. Old-design
docs (routing, gates, delivery-policies, context-loading, etc.) removed.

---

## 6. Template and assembly

`template/` is the blank universal-project-workspace seed: `workspace.yaml`
(uninitialized), `.gitignore` (ignores `/repositories/`,
`/repositories.local.yaml`, `/.runtime/`), the `context/` structure (INDEX,
WORKSPACE, PROJECT, ARCHITECTURE, CONVENTIONS, DECISIONS, SOURCES, sources.yaml,
domains/roles/proposals READMEs), `sources/README.md`, and
`plans/{INDEX.md,README.md,.archived/}`.

`scripts/release-artifact.sh` + `scripts/release-manifest.txt` assemble the
distributable `context-circuit-template` from template-owned trees plus the
seed, excluding source design, tests, scripts, runtime, adapters source, and
maintainer plans. `scripts/build-dist.sh` wraps it (default `v0.5.0`).

---

## 7. Tests (`test/`)

Runner: `sh test/acceptance.sh`. Shared helpers in `test/lib/assert.sh` and
`test/lib/fixture.sh` (builds git-repo workspaces).

| Suite | Covers |
| --- | --- |
| `contracts/` | invariants owner map, manifest, schema presence, skills/adapters/roles present, old-design artifacts absent |
| `runtime/` | unapproved-block, undeclared repo, worktree isolation, commit-before-verify, verifier write rejection, three-failure stop, resume-eligibility, archive semantics, host-neutrality |
| `repositories/` | bindings, personal/team anchor, default `repositories/` dest, git-init+initial-commit, anchor-not-inferred, `workspace` root repo, multi-repo preflight, dirty rejection, deterministic registration |
| `plans/` | stable/sequential/never-reused IDs, validation, block-form plans, index maintenance, approval gate, compound approve-and-execute |
| `execution/` | one branch/worktree per repo, anchor untouched, one worker across repos, commit-before-verify, blocked (not a failure), waived, repair new commit, one active writer, allowed_paths |
| `completion/` | completion blocked before pass, verified≠done, records commits, context-impact without touching PK, failed execution cannot complete |
| `archive/` | archive any status without validation, index removal, active-execution continuation from snapshot, restore, collision safety |
| `delivery/` | source/target report, no delivery action in runtime, skill boundary |
| `security/` | safe paths/ids/plan-ids, symlink rejection, credential absence, gitignore |
| `release/` | assembly, exclusion boundary, only v0.5 skills ship, build-dist |
| `template-runtime/` | assembles the template, runs full lifecycle through the artifact's own engine, proves isolation from source state |

`test/acceptance/criteria-map.yaml` maps design acceptance criteria AC-01…AC-28
to the runnable suite that demonstrates each; the runner verifies the map is
credential-free, complete, and references real suites.

---

## 8. Design coverage

All 28 acceptance criteria (design §21) are demonstrated. Everything the runtime
owns is implemented deterministically and tested. Behaviors the design assigns
to the coordinator/worker/verifier — context gathering and reconciliation
acceptance, provider-native child launch, and actually opening pull requests —
are specified in the skills and role deltas, with the runtime providing the
deterministic evidence those actions rely on (`context-impact-record`,
`verifier-result-record`, `delivery-targets`, execution/completion records).

---

## 9. Notes and boundaries

- Source-side maintainer Product Knowledge (`context/`) and plans
  (`plans/context-circuit-plans/`) from the old design were preserved as
  workspace data, not deleted; they are source-only and excluded from release.
- The engine parses inline and block-list YAML for the plan/task shapes the
  design uses; it is intentionally not a general YAML parser.
- Delivery, merge, push, publication, deployment, archive, and cleanup remain
  separate explicit human actions; the runtime performs none of them.
- No credentials, provider payloads, or transcripts are stored anywhere in
  shipped or workspace files.
