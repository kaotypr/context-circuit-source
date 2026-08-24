# Host matrix

The product is tri-host: Codex CLI, Claude Code, and Cursor Agent are transports
for the same coordinator/worker/verifier contract. So the harness must be able to
run the same scenario library against each host adapter. This document defines how.

## 1. Two host axes

Keep these separate:

- **Product host under test** — the adapter the coordinator-under-test runs on:
  `codex`, `claude-code`, or `cursor-agent`. This is what a run validates.
- **Driver host** — what powers the rig (spawns sub-agents, records traces). It is
  an implementation detail of the test rig, not of the product.

A run names the product host: `--host codex | claude-code | cursor-agent`. The
default is the driver host, so "run it here" tests the host you are on.

## 2. What stays host-neutral

The bulk of the harness does not change per host:

- the **scenario library** (`test/template-runtime/scenarios/`) — natural prompts
  and human-visible expectations;
- the **human-simulator** persona behavior — it is an ignorant user regardless of
  host;
- the **grader**'s product invariants — state post-conditions, transcript checks,
  and the access-discipline policy assert the same product behavior everywhere;
- the assembled template artifact — one product, many transports.

These are authored once and reused for every host.

## 3. What is host-specific

A thin per-host binding covers only the parts that genuinely differ:

| Concern | Per-host binding |
| --- | --- |
| Spawning the coordinator-under-test | each host's native child/sub-agent primitive |
| Coordinator launching worker/verifier | each adapter maps the native child to the worker or independent verifier packet |
| human-simulator definition location | a host-native, source-only spawnable definition (see §4) |
| File-access trace capture | from that host's tool/child log, where exposed (§5) |
| Usage/timing telemetry | from that host's runner, where exposed (§5) |

## 4. human-simulator definition per host

The human-simulator is one host-neutral persona expressed through each host's
native mechanism, kept source-only and never shipped:

- `claude-code`: `.claude/agents/cc-human-simulator.md` (Claude Code sub-agent).
- `codex`: the Codex-native maintainer agent/prompt definition (host-dependent
  location under the source checkout's Codex config).
- `cursor-agent`: the Cursor-native maintainer agent definition (host-dependent
  location under the source checkout's Cursor config).

All three carry the same persona contract from
[roles/human-simulator.md](./roles/human-simulator.md); only the host packaging
differs. The exact Codex/Cursor paths are fixed during implementation against the
current host CLIs; the contract does not change.

## 5. Trace-capture availability and degradation

Grader dimensions C (access-discipline) and D (efficiency) depend on the host
exposing a per-child tool-call and usage trace:

- where the host exposes it, C is a **hard gate** and D is recorded;
- where it does not, C degrades to best-effort transcript inference (**warning
  only**) and D is reported **unavailable**.

A run records, per host, whether the trace was available, so a "pass" is never
mistaken for "audited" when the host could not provide the trace. Dimensions A
(state) and B (transcript) are always available and always hard gates.

## 6. Host-specific expectations

Some product behavior is defined precisely because the host varies; these are
asserted identically on all three adapters:

- **host-blocked** — if the host cannot create an independent read-only verifier
  child, the coordinator reports blocked and never self-verifies (AGENTS.md host
  contract; INV-VERIFY-02).
- **offline / provider-unavailable** — a provider failure falls back to
  filesystem-only work and never blocks deterministic evidence.
- **native child mapping** — a native child maps only to the single worker or the
  independent verifier packet; no host gains a fourth role.
- **permission mode is evidence, not authorization** — host identity, version,
  capability, and permission mode never authorize a route, role, or gate
  (INV-HOST-01).

These are exercised by a dedicated scenario:

- `09-verifier-unavailable-host-blocked` — the run forces the verifier child to be
  unavailable and asserts `host-blocked` with no self-verification and all
  evidence preserved. It must pass identically on `codex`, `claude-code`, and
  `cursor-agent`.

## 7. Runner and summary

- `test/template-runtime/human/run-scenario.sh --host <h>` runs one case on one
  host; an opt-in matrix runner iterates {scenario × host}.
- The run summary is keyed by (scenario, host) and records the per-host verdict,
  the trace-availability flag, and the mapped acceptance criteria and invariants.

## 8. Boundaries

Every per-host binding, human-simulator definition, grader, harness, and
generated result is source-only maintainer material. The release manifest
excludes `test/`, and host config directories (`.claude/`, and the source Codex
and Cursor config) are never part of `context-circuit-template`.
