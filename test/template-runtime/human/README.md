# Human-simulated template test harness (maintainer-only)

A **semantic, conversational** test at the human/product boundary: a driver
spawns the product coordinator inside a fresh workspace assembled from the
released template, an ignorant human-simulator talks to it in plain language, and
a deterministic grader judges the result. It complements — does not replace — the
deterministic `../test-template-runtime.sh` laboratory that exercises the runtime
boundary directly.

This suite is **opt-in and source-only**. It is NOT part of `sh test/acceptance.sh`
(it needs a live model and sub-agent spawning), and the release manifest excludes
all of `test/`, so nothing here ships in `context-circuit-template`.

Design set:
`sources/context-circuit-v0.5-template-test-design/`.

## Files

- `run-scenario.sh` — assemble → isolate → apply setup → snapshot baseline →
  (optionally) drive → grade. The deterministic prep runs anywhere; the
  conversation itself is delegated to a per-host driver.
- `grade.sh` — deterministic grader over ground-truth artifacts (dimensions A–D).
- `.out/` — disposable per-run results (git-ignored).

The human-simulator persona for the Claude Code driver lives at
`.claude/agents/cc-human-simulator.md` (source-only; the assembler's skill
allowlist and the `test/`/`.claude/` exclusions keep it out of the template).

## Running

Prepare a run (no live model needed):

```
sh test/template-runtime/human/run-scenario.sh --host claude-code 01-new-project-simple-idea
```

This prints a run directory under `.out/` containing the isolated `workspace/`,
the pristine `baseline/`, and a `run.yaml` manifest, with `status: prepared`.

### Safety / sandboxing (important)

The `claude-code` driver runs the real coordinator headless with
`--permission-mode bypassPermissions`, which has **no filesystem confinement on
its own**. To observe the coordinator faithfully it keeps the full toolset
(including Bash), so an over-eager coordinator can `git init` / write **outside**
the disposable workspace (observed once: a repo created in the user's home dir).

The driver is therefore **safe by default**:

- it runs each `claude` invocation inside a **bubblewrap** sandbox (whole FS
  read-only except the workspace, the run dir, `/tmp`, and `~/.claude`), and
- if bubblewrap cannot confine writes on this host (e.g. unprivileged user
  namespaces are disabled — `bwrap: setting up uid map: Permission denied`), the
  driver **refuses to run** unless you opt in to the risk, one of:
  - `CC_ALLOW_UNSANDBOXED=1` in the environment (per-run), or
  - a local, **gitignored** marker so unsandboxed becomes the default on THIS
    machine without shipping the risk: `touch drivers/.allow-unsandboxed`.
  Prefer running the whole harness inside a throwaway **container/VM** instead.
- a post-run **isolation guard** flags any repository binding whose path escapes
  the workspace and records `isolation_violation:` in `run.yaml`.

### Driving the conversation

`run-scenario.sh` delegates the live conversation to a driver (`--driver CMD` or
`$CC_HUMAN_DRIVER`). The driver is invoked with these environment variables and
must, for the chosen host:

1. spawn the product coordinator with its working directory set to `$CC_WORKSPACE`
   (it loads that workspace's own `AGENTS.md`/`CLAUDE.md`/`WORKFLOW.md`/skills);
2. spawn the human-simulator (`$CC_HUMAN_SIM`) given ONLY the case `human:` block;
3. run the turns and append every prompt/reply to `$CC_TRANSCRIPT`;
4. optionally record `$CC_TRACE` and `$CC_TELEMETRY` (see below).

| Variable | Meaning |
| --- | --- |
| `CC_RUN_DIR` | the run directory |
| `CC_WORKSPACE` | instantiated workspace (coordinator cwd) |
| `CC_BASELINE` | pristine snapshot for dimension-A diffs |
| `CC_CASE_FILE` / `CC_CASE_ID` | the case being run |
| `CC_HOST` / `CC_MODE` | product host; `conversation-only` or `full-execution` |
| `CC_HUMAN_SIM` | human-simulator definition path |
| `CC_TRANSCRIPT` | append prompts + replies here (required) |
| `CC_TRACE` | optional TSV `action<TAB>tool<TAB>path` |
| `CC_TELEMETRY` | optional TSV `action<TAB>turns<TAB>tokens` |

On Claude Code the driver is a Claude session using the Task tool to spawn the
two sub-agents; the Codex and Cursor bindings are deferred (host-matrix §4).

### Grading

Automatic after a driver run, or manually against any run directory that has a
transcript:

```
sh test/template-runtime/human/grade.sh .out/<run-id>
```

## Grader dimensions

- **A. State post-conditions** (hard gate) — workspace/plan/execution state via
  the artifact's `engine.sh` reads and file checks; two checks
  (`product_knowledge_unchanged_silently`, `workspace_identity_not_fabricated`)
  are diffed against the captured baseline, since the runtime exposes no read for
  them.
- **B. Transcript checks** (hard gate) — `forbids_regex` (no internals leaked to a
  lay user) and `requires_any` (a natural intent phrase). Note "verifier" is
  intentionally NOT forbidden — the product surfaces it by design.
- **C. Access-discipline audit** (**session-level**; forbidden = hard gate,
  required = advisory) — the file-access trace vs the `access_policy`. Because
  which turn a read lands on is non-deterministic across model runs, C is
  evaluated over the whole session, not per turn: a `forbidden` path read **at any
  point** fails (evaluated as the intersection of the occurring actions' forbidden
  sets); a `required` file miss is **advisory (warning), never a failure**, since a
  capable model can reach a correct outcome via a different read path. Occurrence
  of an action is judged from workspace state (e.g. a `plan.yaml` exists), never
  from turn timing; the trace's action label is retained only for human reading.
  **Requires the host to expose a file-access trace; where it does not, C degrades
  to warning-only** and the verdict rests on A and B (host-matrix §5). The
  `claude-code` driver captures the trace from the coordinator's `Read`/`Grep`/
  `Glob` tool calls (stream-json).
- **D. Efficiency ledger** (soft, warning-only) — per-action turns/tokens vs
  `budgets`; never fails a run.

A run passes only when every hard-gate dimension (A, B, and C when enforced)
passes.

## Boundaries

Every generated workspace, fixture, worktree, transcript, and result is
disposable test state under `.out/` and never touches the source `.runtime/`.
Only `01-new-project-simple-idea` is scaffolded; cases 02–09 and the
non-empty-fixture setup path (`setup.repositories`/`setup.sources`) are follow-on
work.
