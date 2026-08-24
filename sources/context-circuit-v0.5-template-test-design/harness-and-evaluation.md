# Harness, isolation, and evaluation

This document defines how one scenario runs, how the workspace is isolated, how a
run is graded, and the boundaries that keep the harness maintainer-only.

## 1. Run lifecycle

For one scenario case:

1. **Assemble the template.** Build `context-circuit-template` with
   `scripts/release-artifact.sh` into a temporary directory. This is the same
   artifact a user would receive.
2. **Instantiate an isolated workspace.** Copy the assembled artifact into a
   fresh temp directory. This is the coordinator's working directory. It must not
   import the source repository's Product Knowledge, plans, `.runtime/`,
   bindings, or implementation state — only the shipped template. Immediately
   record a read-only **baseline snapshot** of the pristine tree (at least
   `context/` and `workspace.yaml`) before the coordinator runs, so the grader can
   later prove Product Knowledge was not silently changed and workspace identity
   was not fabricated (dimension A). The runtime exposes no read for these two
   post-conditions; they are a diff against this baseline.
3. **Apply case setup.** Create any `setup.repositories` fixture git repos and
   `setup.sources` files. Case 01 creates none.
4. **Spawn the coordinator sub-agent** under the run's chosen product host
   (`--host codex | claude-code | cursor-agent`, default = the driver host), with
   its working directory set to the instantiated workspace. It loads that
   workspace's `AGENTS.md`, `WORKFLOW.md`, skills, and role deltas as a real
   product would. The scenario library is host-neutral; only the launch mechanism
   and trace capture are per-host — see [host-matrix.md](./host-matrix.md).
5. **Spawn the human-simulator** with only the case's `human:` block.
6. **Drive the conversation** turn by turn (§3), capturing a transcript, the
   **file-access trace** (which files each actor opened, via which tool, from the
   runner's tool log), and — where the runner exposes it — per-action
   **usage/timing telemetry**.
7. **Grade** with the human-simulator's conversational verdict plus the
   deterministic grader over the workspace state, transcript, file-access trace,
   and telemetry (§4; dimensions defined in
   [roles/grader.md](./roles/grader.md)).
8. **Record the result** and dispose of the temp workspace (§6).

## 2. Isolation

- Every run uses its own temp assembly and temp workspace; nothing touches the
  source checkout or its `.runtime/`.
- Generated workspaces, fixture repos, worktrees, commits, and transcripts are
  disposable test state and live outside the source-side `.runtime/`.
- The coordinator may create branches, worktrees, and runtime records only inside
  the isolated workspace.
- Fixture git repositories are local, credential-free, and created per run; the
  harness never uses a real remote and never pushes.

## 3. Turn protocol

```
human-simulator ──prompt──▶ coordinator-under-test
                ◀─reply────
```

- The human-simulator sends one natural prompt, waits for the coordinator's full
  reply, then chooses the next prompt using the case reaction rules.
- The coordinator behaves normally; in `full-execution` cases it may spawn the
  single worker and the independent verifier as its own sub-agents.
- Every prompt and reply is appended to a plain-text transcript for grading and
  human review.
- The loop ends when the case turns are exhausted, a hard `visible_expectation`
  is violated, or a safety stop triggers.

## 4. Evaluation

Grading has two independent parts; both must pass.

### 4.1 Conversational verdict (human-simulator)

Judged from the transcript alone, at the level of product invariants rather than
exact wording. The human-simulator reports, per `visible_expectation`, a
pass/fail with the sentence of evidence that satisfied or violated it. It focuses
on observable product behavior: did it refuse the draft run, did it avoid
inventing a repository, did it keep jargon out, did it decline to claim work was
done.

### 4.2 Deterministic grader dimensions

`grade.sh` judges four dimensions from ground-truth traces (full detail in
[roles/grader.md](./roles/grader.md)):

- **A. State post-conditions** (hard gate) — workspace/plan/execution state via
  `engine.sh` reads and file checks (e.g. "no plan advanced past draft", "commits
  recorded before verification", "done only after a verifier pass").
- **B. Transcript checks** (hard gate) — `forbids_regex` (no internals leaked to a
  lay user) and `requires_any` (a natural intent phrase).
- **C. Access-discipline audit** (hard gate) — the file-access trace vs a
  per-action `access_policy` of `required` / `allowed` / `forbidden` files; this
  proves the coordinator read the right files and only the necessary files
  (archived exclusion, sources-passive, never loading the runtime as a
  substitute).
- **D. Efficiency ledger** (soft, warning-only) — per-action turns, tokens,
  tool-call breakdown, wall-clock, and the process-overhead ratio, compared to
  per-action `budgets`.

A case passes only when the conversational verdict and every **hard-gate**
dimension (A, B, C) pass; dimension D annotates the result and never fails a run.
Each finding links back to the design acceptance criterion and the
`wrapper/contracts/invariants.yaml` rule it protects.

### 4.3 Non-determinism policy

Conversations vary. Grading must key on **product invariants**, never on exact
phrasing:

- assert behaviors and side effects (refused, asked one question, nothing ran,
  status unchanged), not sentences;
- allow paraphrase in `requires_any` by matching intent phrases, not full lines;
- treat a genuine invariant violation as a hard failure; treat wording drift as a
  pass;
- allow a small bounded retry for transient model/tool errors, but never retry
  past an invariant violation.

## 5. Where results go

- Transcript, conversational verdict, and grader report are written under a
  disposable per-run directory (for example `test/template-runtime/human/.out/`,
  git-ignored), never under the source `.runtime/`.
- A run summary lists each case, its mapped acceptance criteria, and pass/fail.

## 6. Boundaries and safety

- The human-simulator role, grader, harness, scenarios, and all generated state
  are source-only. The release manifest excludes `test/`, and `.claude/` is never
  shipped, so nothing here can enter `context-circuit-template`.
- The product retains exactly three roles. The human-simulator is a test actor,
  not a fourth product role, and no `cc-*` product skill is added for it (the
  assembler's skill allowlist would reject one).
- No credentials, remotes, or real user repositories are used; every fixture is
  local and disposable.
- The harness preserves the source checkout untouched and cleans up its own temp
  directories.

## 7. Relationship to the deterministic laboratory and CI

- `test/template-runtime/test-template-runtime.sh` stays in `test/acceptance.sh`:
  it is fast, deterministic, and needs no live model.
- The human-simulator suite is **agent-driven** (it needs a live model and
  sub-agent spawning), so it is a separate, opt-in maintainer suite — invoked
  explicitly, not part of the deterministic `sh test/acceptance.sh` run, and not
  a required CI gate unless a live-agent runner is available.
- Both layers share the same assembled artifact and the same source of truth for
  "correct": the product design and `wrapper/contracts/invariants.yaml`.

## 8. Suggested implementation order

1. Add the human-simulator role for the driver host (start with
   `.claude/agents/cc-human-simulator.md`) and the coordinator-under-test spawn
   contract.
2. Scaffold `test/template-runtime/scenarios/` with case 01 and its `case.yaml`.
3. Implement `test/template-runtime/human/run-scenario.sh --host <h>` (assemble,
   isolate, drive) and `grade.sh` (dimensions A–D).
4. Add cases 02–09 (09 is the host-blocked case).
5. Add the per-host bindings for `codex` and `cursor-agent`
   (see [host-matrix.md](./host-matrix.md)).
6. Add an opt-in matrix runner that iterates {scenario × host} and prints the
   summary keyed by (scenario, host) with the trace-availability flag.
