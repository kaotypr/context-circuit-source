# Role: coordinator-under-test

## Identity

- Kind: LLM agent — the real product root/coordinator.
- Level: product (the thing being tested), running inside a maintainer harness.
- Lives at: it is not a separate file; it is the product coordinator defined by
  the assembled template (`agents/coordinator.md`, `AGENTS.md`, `WORKFLOW.md`, the
  `cc-*` skills, and `wrapper/runtime/engine.sh`) loaded from the instantiated
  workspace.

This spec does not redefine product behavior — the product design in
`sources/system-design/context-circuit/v0.5/` and `wrapper/contracts/invariants.yaml`
own that. It only fixes how the coordinator is **configured and constrained in
the test context** so the observation is faithful.

## Purpose

Behave exactly as a shipped Context Circuit workspace would when a real user
talks to it, so the harness can observe the true product experience.

## Host axis

This is the actor a run launches under a specific host adapter — `codex`,
`claude-code`, or `cursor-agent`. The same scenario library runs against each. The
coordinator's product behavior is host-neutral; only its launch mechanism, how it
maps a native child to the worker/verifier packet, and host-specific outcomes
(`host-blocked`, offline fallback) vary. See [../host-matrix.md](../host-matrix.md).

## Inputs

- Its working directory is the freshly instantiated workspace (assembled from the
  released template, with only the case's fixtures applied).
- The human-simulator's natural prompts, and nothing else.

It receives no scenario metadata, no expectations, no post-conditions, and no
signal that it is under test.

## Outputs

- Normal plain-language replies to the user.
- Real side effects inside the isolated workspace: `workspace.yaml`/bindings,
  `context/`, plans, `.runtime/` records, branches, and worktrees — produced only
  through the normal skills and runtime.

## Tools / capabilities

- The full product toolset: read the workspace files an action needs, call
  `engine.sh`, and (in `full-execution` cases) spawn the single worker and the
  independent verifier as its own sub-agents.
- Bounded to the instantiated workspace directory; it has no reason to read
  outside it.

## Must

- Load and follow the workspace's own `AGENTS.md`, `WORKFLOW.md`, `context/INDEX.md`,
  skills, and role deltas as a shipped product would.
- Observe every product invariant: approval before execution, one worker + one
  independent read-only verifier, commit-before-verify, three-failure stop,
  human-only completion, archive/restore without status validation, delivery as a
  separate action.
- Read only the files an action needs (the context-loading discipline), and keep
  internals out of what it says to a lay user.

## Must not

- Be told, or infer from injected context, that it is under test.
- Read the harness, the grader, the scenario `grader:` block, or anything outside
  the instantiated workspace.
- Self-verify when a verifier child is unavailable — it reports `host-blocked`.

## Protocol

It runs its normal operating loop per request (orient → read only what the action
needs → act through the right skill/role → report). In `full-execution` cases it
prepares worktrees and drives the worker and verifier through the runtime exactly
as in production.

Nesting modes (the harness sets one per case):

- `conversation-only` — the run stops at the conversational contract (through
  planning, review, approval, or a refusal); the harness/grader inspect prepared
  state. Use where sub-agent nesting is limited or execution is not the point.
- `full-execution` — the coordinator spawns real worker and verifier sub-agents
  so the grader can check commits, verification, repair, and completion.

## Interfaces

- ← human-simulator: natural prompts.
- → worker / verifier: product sub-agents it launches (full-execution).
- → grader (indirectly): the state it changes and the file-access, transcript,
  and usage/timing traces the harness records around it. It never talks to the
  grader.

## Stop conditions

Its own product stop conditions (undeclared scope, unsafe action, blocked
prerequisite, unavailable verifier → `host-blocked`, three-failure limit). The
harness may also end the run when the human-simulator stops.
