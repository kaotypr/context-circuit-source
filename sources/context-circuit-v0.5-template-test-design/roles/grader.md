# Role: grader

## Identity

- Kind: mechanical, deterministic — **not** an LLM agent.
- Level: maintainer test actor (source-only; never shipped).
- Lives at: `test/template-runtime/human/grade.sh`.

The grader is mechanical on purpose. What files were read, what state changed,
and how much was consumed are ground-truth traces from the runner. Judging them
with an LLM would be non-deterministic and would itself read files and spend
tokens, contaminating the very things being judged.

## Purpose

Decide a scenario run's pass/fail from **facts** — the workspace state, the
file-access trace, the transcript, and (optionally) usage/timing — and map every
result back to the product acceptance criterion and invariant it protects.

## Inputs

- The instantiated workspace path (read-only).
- The recorded transcript (all prompts and replies).
- The **file-access trace**: which files each actor opened and via which tool,
  from the runner's tool log (availability varies by host, see below).
- Optional **usage/timing telemetry**: per-action tokens, tool-call counts, and
  wall-clock from the runner.
- The product host the run used (`codex` / `claude-code` / `cursor-agent`) and a
  flag for whether the trace was available on that host.
- The case's `grader:` block (acceptance-criteria and invariant mapping,
  post-conditions, transcript checks, access policy, budgets). It never receives
  the `human:` block's private reasoning.

## Host-neutral

The grader is host-neutral: its product invariants (dimensions A–C) assert the
same behavior on every adapter. Only trace availability differs by host — where a
host does not expose the file-access/usage trace, dimension C degrades to
best-effort transcript inference (warning-only) and dimension D is unavailable.
See [../host-matrix.md](../host-matrix.md) §5.

## Outputs

- A per-dimension result (below), an overall verdict, and for each finding the
  design acceptance criterion (`AC-nn`) and `wrapper/contracts/invariants.yaml`
  rule it maps to.
- A line in the run summary.

## Dimensions

### A. State post-conditions — hard gate

Inspects the workspace using the runtime's read-only operations and file checks:
repositories registered count, plans created and their status, presence/absence
of execution records, commit-before-verify ordering, and `done` only after a
verifier pass. Two post-conditions have no runtime read and are diffed against
the pristine baseline snapshot the harness captured at instantiation (see
[../harness-and-evaluation.md](../harness-and-evaluation.md) §1): Product
Knowledge was not silently changed, and workspace identity was not fabricated (a
real identity appears only when the human supplied one; the shipped placeholder
is a pass).

### B. Transcript checks — hard gate

Over the recorded conversation: `forbids_regex` (no `cc_*`/`engine.sh`/`worktree`/
`plan.yaml` leaked to a lay user) and `requires_any` (a natural intent phrase,
matched by intent, not exact wording).

### C. Access-discipline audit — `forbidden` is the hard gate; `required` is advisory

The core v0.5 check: did the coordinator read **only the necessary** files and
never the internals it must not? Compares the file-access trace to a per-action
`access_policy`, evaluated at the **session** level (which turn a read lands on is
non-deterministic across model runs, so per-turn attribution is unreliable):

- `forbidden` — files it must not read; a match **anywhere in the session fails**
  (the hard gate), e.g. `plans/.archived/**` (archived exclusion),
  `wrapper/runtime/engine.sh` (don't load the runtime as a substitute), `sources/**`
  unless explicitly named (INV-SEC-02). Evaluated as the **intersection** of the
  forbidden sets of the actions that actually occurred, so a path forbidden only
  for some actions (e.g. `wrapper/contracts/**` for orient but not create-plan) is
  not flagged session-wide.
- `required` — files the action ought to read; a miss is **advisory (warning), not
  a failure**, because a capable model can reach a correct outcome via a different
  read path (state, tool invocation, already-loaded context). Enforced only for an
  action that actually occurred (occurrence judged from workspace state).
- `allowed` — files it may read.

This dimension proves the v0.5 promise that motivated the refactor: never
reconstruct the lifecycle by loading internals. It requires the runner to expose
the file-access trace; where it does not, the whole dimension degrades to
best-effort transcript inference and is reported as warning-only.

### D. Efficiency ledger — soft, warning-only

Records, per product action, the turns, tokens/context, tool-call breakdown, and
wall-clock from telemetry, plus a derived **process-overhead ratio** (cost spent
operating Cc — runtime calls, contracts, records — versus doing project work —
Product Knowledge, repo code, plan/task content). Compared to per-action
`budgets` as **warnings**, never hard failures, because these numbers are
non-deterministic and model-dependent. Overhead attribution starts as a
heuristic (classify by which files were touched); see the analyst note below.

## Verdict

A run passes only when every **hard-gate** dimension passes: A (state), B
(transcript), and C's **forbidden**-read check. C's `required`-read misses and
dimension D never fail a run; they annotate the result and feed trend tracking
across template versions. Each failing finding names the violated `AC-nn` and
invariant.

## Case `grader:` schema

```yaml
grader:
  acceptance_criteria: [ AC-nn, ... ]  # §21 criteria this case demonstrates
  invariants: [ INV-..., ... ]         # invariants.yaml rules it protects
  post_conditions: [ ... ]            # dimension A
  transcript_checks: [ ... ]          # dimension B
  access_policy:                       # dimension C, keyed by action
    <action>:
      required:  [ ... ]
      allowed:   [ ... ]
      forbidden: [ ... ]
  budgets:                             # dimension D (soft), keyed by action
    <action>: { max_turns: N, max_tokens: M }
```

## Must

- Use only ground-truth traces and read-only inspection; be fully deterministic.
- Key on invariants and side effects, never on exact wording (§4.3 policy).
- Map every finding to an acceptance criterion and invariant owner.

## Must not

- Use an LLM to reach its verdict (the optional analyst is out-of-loop and
  advisory only, never part of the pass/fail path).
- Read internal state during the conversation or in any way that could influence
  the coordinator's behavior.
- Fail a run on efficiency numbers alone.

## Optional efficiency-analyst (not a standing role)

If heuristic overhead attribution proves too crude, a thin analyst may run
**once, out-of-loop, after** the conversation, reading only the finished
transcript and ledger to label ambiguous cost as overhead-vs-work. It is
advisory: it enriches dimension D and never changes the hard-gate verdict. It is
never in-loop, so it cannot perturb the run.

## Interfaces

- ← harness: workspace path, transcript, file-access trace, optional telemetry,
  and the case `grader:` block.
- → run summary: per-dimension results, overall verdict, and invariant mapping.
