# Test-harness roles

This folder gives each actor in the template test harness its own dedicated
design. Every role spec follows the same shape: identity, purpose, inputs,
outputs, tools/capabilities, must / must-not, protocol, interfaces, and stop
conditions.

## The actors

| Role | Kind | Level | Spec |
| --- | --- | --- | --- |
| human-simulator | LLM agent | maintainer test | [human-simulator.md](./human-simulator.md) |
| coordinator-under-test | LLM agent | product | [coordinator-under-test.md](./coordinator-under-test.md) |
| grader | mechanical (deterministic) | maintainer test | [grader.md](./grader.md) |

The **harness** (orchestrator) is not an actor; it is a mechanical driver
specified in [../harness-and-evaluation.md](../harness-and-evaluation.md).
An optional **efficiency-analyst** is described inside the grader spec; it is not
a standing role and never runs in-loop.

## Why only these three, and why the split matters

- The **human-simulator** is an agent because it must produce natural, adaptive
  conversation. It is deliberately ignorant of internals.
- The **coordinator-under-test** is the real product agent. The harness only
  configures and observes it; it is never told it is under test.
- The **grader** is mechanical, not an agent, on purpose: what files were read,
  what state changed, and how much was consumed are **ground-truth traces** from
  the runner. Using an LLM to judge those would be non-deterministic and would
  itself read files and spend tokens, contaminating the very things being judged.

Judgment is therefore split by trustworthiness:

- conversational quality → human-simulator (from words it saw);
- state, file-access discipline, and cost → grader (from exact traces).

## Interaction contract

The human-simulator and grader are host-neutral; the coordinator-under-test is
launched under the run's chosen host adapter (`codex` / `claude-code` /
`cursor-agent`). See [../host-matrix.md](../host-matrix.md).

```
harness (runs under one product host per run)
  ├─ assembles template → isolates workspace → applies case fixtures
  ├─ spawns coordinator-under-test  (cwd = instantiated workspace)
  ├─ spawns human-simulator         (given only the case `human:` block)
  │     human-simulator ──natural prompt──▶ coordinator-under-test
  │                     ◀──plain reply─────
  │     (coordinator may spawn worker/verifier in full-execution cases)
  ├─ records: transcript + file-access trace + state + (optional) usage/timing
  └─ invokes grader (given only the case `grader:` block + the recorded traces)
        grader → per-dimension results + overall verdict → run summary
```

No role reads another role's private block: the human-simulator never sees the
`grader:` block, and the coordinator never sees any test metadata.
