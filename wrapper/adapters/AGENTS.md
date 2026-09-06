# Context Circuit workspace

This is a Context Circuit universal project workspace: an agent-oriented place to
hold Product Knowledge, create grounded plans, and execute them safely across one
or more Git repositories. Talk to it in ordinary language.

The root conversational session is the coordinator. Before its first
user-facing reply, it reads and follows `agents/coordinator.md`, which owns
conversation routing and reporting language for every host.

## How it works

A planned request flows as a normal conversation: the coordinator reads the relevant
Product Knowledge and plan material, drafts or reviews a readable plan, you
approve it in conversation, the workspace prepares isolated repository worktrees,
one worker implements the whole plan in that plan's single repository and commits it, an independent
verifier checks the latest commits, and the worker repairs failures with new
commits. A plan becomes done only on an explicit mark-done at Standard and
Critical. Explore is planless and human-supervised.

Branches, worktrees, runtime records, and verifier setup are hidden. The plan,
the worker's handoff, the verifier's result, and your authority over completion
are not hidden.

For small, live changes, the user may instead work directly with the coordinator
and one worker. That mode uses a separate working copy, has no plan or independent
verifier, and is always described as human-supervised rather than verified.

## Safety spine

- Read first: `WORKFLOW.md`, `workspace.yaml`, and `context/INDEX.md`. Read only
  the plan and context units an action needs.
- Rule ownership is indexed in `wrapper/contracts/invariants.yaml`. One rule has
  one owner; do not add parallel policy to a skill or role file.
- The human gate is on the intent (Gate 1), not the plan: on approval a read-only
  tracer reads the real code and the coordinator runs a feasibility check, then a
  plan derives from the approved intent and executes with no separate plan gate and
  no automated scope gate (scope-safety is settled at delivery, Gate 2). Intent
  approval and delivery are explicit conversational actions, never confirmation
  cards or hidden tokens.
- During plan execution, one worker writes and one independent read-only verifier checks. If the verifier
  child cannot be created, the result is `host-blocked` — never self-verify.
- During direct collaboration, one worker writes under live human supervision;
  no verifier is created and the result never gains verified status (INV-PAIR-01).
- `sources/` is passive: read only exact request-named files; never scan all
  sources or sibling workspaces to fill a context gap.
- `plans/archive/` is outside normal context; read it only via explicit restore.
- Pull requests, merge, push, deployment, archive, and cleanup are separate
  explicit human actions. A plan becomes done only on an explicit mark-done at
  Standard and Critical. Verification, candidate acceptance, and delivery never
  imply completion.
- Credentials stay in host Git config or the host agent; never in workspace files
  or runtime records.

## Runtime

`wrapper/runtime/engine.sh` is a small host-neutral deterministic library for
workspace, Git, and execution-state operations. A role that invokes the engine
**must not read** `wrapper/runtime/engine.sh` or any runtime implementation file:
invoke actions as `sh wrapper/runtime/engine.sh <action> <args>` and consume their
printed results. The invoking skill and the execution brief carry everything needed
to drive every action by construction, and because the runtime holds no prompts,
Product Knowledge interpretation, or routing policy (INV-RUNTIME-01) there is
nothing in it to interpret — only actions to call. This is the coordinator-side
corollary of INV-RUNTIME-01; every skill that drives the engine references this
statement rather than restating it.

## Host adapters

Codex, Claude Code, and Cursor Agent are transports. Host identity, version,
capability, permission mode, and provider status are bounded provider-neutral
`host_evidence` only; they never authorize a route, role, verification, or
completion. A native child maps to the single worker (for execution or direct
collaboration) or the independent verifier (for execution only).

## Per-role model & effort

The coordinator may run the worker, verifier, and tracer at a per-role `(model, effort)`
from an optional host-local, per-user, gitignored `role-tiering.local.yaml`,
grouped by host so each host names the models available on it. Read that file
from the workspace root (next to `repositories.local.yaml`) before spawning;
never look for it inside the child's working copy. A missing file in an
isolated working copy is not an absent config. This is bounded
host evidence (`host_evidence`): it changes cost and speed, never a route, role,
lease, verification, completion, verifier independence, or the failure limit, and
the runtime never learns it (INV-RUNTIME-01). Reading the config does not apply
it — the host adapter sets the model on the child spawn; absent any config, the
adapter defaults apply. Full rules: `docs/role-tiering.md`.

### Applying a configured tier on Codex

When this host is Codex and a role has a configured model or effort, launch the
child with `spawn_agent` using that exact `model` and `reasoning_effort`. Read
`role-tiering.local.yaml` from the workspace root in the coordinator session
first; keep that session at the workspace root. Only the child's working
directory is the isolated copy. A model
or effort override requires `fork_turns: "none"`; provide the complete role,
scope, working-copy path, and task in the spawn prompt instead of relying on
forked conversation context. End the Codex `task_name` with `_worker`,
`_verifier`, or `_tracer` so bounded host evidence can identify the role without
retaining a provider prompt. For direct collaboration, apply the `worker` tier to
its one worker and never launch a verifier. For tracing, apply the `tracer` tier
when that host group has a tracer entry.

The coordinator on any host may keep one resumable root session, represented by
a session or thread id, for the human conversation. That root is not a worker,
verifier, or tracer. Every required child role must use the host's native
child-agent primitive and be attached to that coordinator: Codex uses
`spawn_agent`, Claude Code uses `Task`/subagent, and Cursor Agent uses its native
child feature when available. A separate top-level task, peer thread, or resumed
root session does not satisfy the child requirement. If the native child cannot
be created, the route is `host-blocked` and remains read-only — do not replace
the child with another conversation session or self-verification.
