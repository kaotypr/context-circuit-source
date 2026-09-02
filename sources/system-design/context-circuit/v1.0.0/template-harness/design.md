# Shared role tiering and three-host scenario matrix

## Capability

The template harness runs one scenario library against Codex CLI, Claude Code,
and Cursor Agent CLI with the same scenario inputs and the same deterministic
grader. Every generated run has an isolated workspace, a host-local role
configuration, a host-specific driver, and a result keyed by
`(scenario, host)`.

The harness-wide role configuration is authored once. It is grouped by host so
each adapter can name models that actually exist on that host. Codex uses
`gpt-5.6-luna` for both worker and verifier at `medium` effort. Claude Code and
Cursor Agent have separate groups and may use different host-native model IDs,
while sharing the same role shape and default effort.

## Problem

Role tiering currently lives as an optional flat `setup.role_tiering` block in
individual case files. That makes the behavior scenario-specific, prevents a
single host matrix from having one authoritative preference, and requires the
grader to hard-code model names that only make sense on one host. The harness
also has no common runner that schedules the full case library across all host
adapters.

## Fixed decisions

1. `template-harness/fixtures/role-tiering.yaml` is the single tracked harness
   fixture. It contains `hosts.codex`, `hosts.claude-code`, and
   `hosts.cursor-agent`, each with `worker` and `verifier` entries.
2. `run-scenario.sh` copies that fixture into the disposable workspace as
   `role-tiering.local.yaml`. The copied file remains host-local and ignored;
   the tracked fixture is source-only and never enters the released template.
3. Case files contain project and conversation fixtures only. They do not own
   model, effort, or escalation settings.
4. Grading resolves configured model expectations from the selected host group
   instead of embedding one model ID in a case. Host adapters that cannot apply
   per-child effort may record it as evidence while degrading to model-only,
   as defined by `docs/role-tiering.md`.
5. `run-matrix.sh` enumerates all case directories and runs them in three
   concurrent host lanes. A lane processes its scenarios sequentially, giving
   each `(scenario, host)` pair a separate run directory. The runner completes
   all lanes, writes a keyed summary, and returns failure only after every lane
   has reported.
6. Each host binding owns only transport differences: starting/resuming the
   coordinator, running the human simulator, and recording the evidence that
   host can expose. It does not add routing, authorization, verification, or
   completion policy.

## Shape

```text
template-harness/
├── fixtures/
│   └── role-tiering.yaml          # one tracked matrix input
├── scenarios/                     # shared human + grader cases
└── human/
    ├── run-scenario.sh            # one isolated (case, host) run
    ├── run-matrix.sh              # three parallel host lanes
    └── drivers/
        ├── codex.sh
        ├── claude-code.sh
        └── cursor-agent.sh
```

```mermaid
flowchart LR
    F[role-tiering.yaml] --> P[run-scenario.sh]
    C[case.yaml] --> P
    P --> W[disposable workspace]
    F --> L[role-tiering.local.yaml]
    L --> W
    W --> D{host driver}
    D --> CX[Codex coordinator]
    D --> CL[Claude Code coordinator]
    D --> CU[Cursor coordinator]
    CX --> R[worker/verifier roles]
    CL --> R
    CU --> R
    R --> G[deterministic grade]
    G --> S[(scenario, host) summary]
```

The matrix is a fan-out over the host dimension, not a shared execution. The
three host lanes may be active simultaneously, but no workspace, runtime
directory, transcript, or provider session is shared between runs.

## Configuration contract

The tracked fixture uses the same host-grouped shape as the product's
host-local file:

```yaml
schema_version: 1
hosts:
  codex:
    worker:
      model: gpt-5.6-luna
      effort: medium
      escalate_on_repair: false
    verifier:
      model: gpt-5.6-luna
      effort: medium
      escalate_on_repair: false
  claude-code:
    worker: { model: sonnet, effort: medium, escalate_on_repair: false }
    verifier: { model: haiku, effort: medium, escalate_on_repair: false }
  cursor-agent:
    worker: { model: gpt-5, effort: medium, escalate_on_repair: false }
    verifier: { model: gpt-5, effort: medium, escalate_on_repair: false }
```

The concrete IDs are host configuration, not a cross-host equivalence claim.
Changing a host's model changes cost and speed only. It cannot authorize a
route, remove the verifier, alter the repair limit, or mark a plan complete.

Scenario 14 is the deliberate Explore/direct-collaboration case: its shared
fixture still has a verifier entry for matrix consistency, but the case itself
must not launch a verifier and is graded as human-supervised rather than
independently verified.

## Matrix execution

The default live command is:

```text
sh template-harness/human/run-matrix.sh
```

It launches one lane for each host. Each lane calls
`run-scenario.sh --host <host> --live <case>` and keeps the generated run
artifacts under `template-harness/human/.out/`. The deterministic smoke mode is:

```text
sh template-harness/human/run-matrix.sh --prepare-only
```

`--case <id>` may be repeated to limit a matrix during development. The live
matrix is opt-in and remains outside the credential-free acceptance gate. A
missing CLI, unavailable child capability, or host-specific driver failure is
reported in that pair's log and summary; another host lane continues running.

## Evidence and grading

The deterministic grader continues to own the semantic verdict. The host
driver supplies only transcript, optional trace, optional telemetry, and
bounded role evidence. A case that checks tiering uses a configured form such
as `role_evidence_configured: 0001-alpha:worker`; the grader reads the selected
host group from the generated local file and compares the observed role model
to that value. This keeps one case valid across the matrix without moving host
policy into the grader or runtime.

Codex supplies bounded child evidence from its local session metadata. Claude
Code and Cursor Agent use the evidence surfaces their CLIs expose; when a host
does not expose a dependable file trace, dimension C degrades to warning-only.
The harness never copies prompts, provider payloads, credentials, or auth state.

## Non-goals and boundaries

- Do not copy the tracked fixture into the released template; only the
  disposable harness workspace receives `role-tiering.local.yaml`.
- Do not let a scenario override the shared role configuration.
- Do not use the matrix runner as a second product lifecycle or approval gate.
- Do not make host availability authorize execution or replace a required
  verifier; unavailable child capability remains `host-blocked`.
- Do not treat equal model/effort settings as loss of verifier independence.

The design introduces no new invariant. It composes the existing host evidence,
runtime opacity, verifier independence, isolation, and direct-collaboration
contracts; implementation details belong to the harness scripts and the
host-adapter surfaces they already exercise.
