# Cursor Agent adapter

@AGENTS.md

Follow the shared `AGENTS.md` and `WORKFLOW.md` contract. Cursor Agent is a thin
transport. The root Cursor session is the coordinator (see `.context-circuit/agents/coordinator.md`)
for all normal conversation: orientation, context gathering, planning, review,
approval interpretation, direct collaboration, completion, archive/restore, and
delivery discussion.

A Cursor Task/subagent maps only to the single bounded worker (`.context-circuit/agents/worker.md`)
for an execution or direct-collaboration session, the independent read-only
verifier (`.context-circuit/agents/verifier.md`) for one execution, or the planner
(`.context-circuit/agents/planner.md`) after intent approval.
Record provider-neutral `host_evidence` for the child; a host permission flag is
an observation, not authorization.

## Applying per-role model & effort

When `role-tiering.local.yaml` configures a `(model, effort)` for a role (see
`.context-circuit/docs/role-tiering.md`), read the file from the workspace root (next to
`repositories.local.yaml`), never from a Task working directory, then read the
group for this host from `hosts.cursor-agent` and apply it — which on Cursor
means **setting it on the spawn, not just recording it**. Open the path with
Read; gitignore is not a read block, and a search that skips ignored files is
not an absent config. A missing file in an isolated working copy is not an
absent config. Concretely:

- **Model — set it on the spawn.** Pass the role's configured model id to the
  worker's, verifier's, or planner's `Task` spawn as its `model` parameter. The
  value must be a slug the host lists on Task (for example `composer-2.5-fast`,
  `cursor-grok-4.6-high`). A child launched without `model`, or with
  `model: inherit`, inherits the coordinator's session model, so an inherit or
  unset spawn is the whole tier being silently ignored.
- **The yaml is the user's explicit model request.** Cursor's Task tool prefers
  `inherit` unless another listed model was requested. `role-tiering.local.yaml`
  *is* that standing request for child roles. When `hosts.cursor-agent` names a
  listed slug for `worker`, `verifier`, or `planner`, pass that slug. Do not default to inherit because the root session already has a model.
- **Invalid or display names.** Names such as `Auto` are not Task slugs; do not
  pass them. If the configured id is not on the host's Task allow-list, do not
  omit `model` (omitting inherits). Report the mismatch, or use a listed slug
  the user already named for that role — never silently inherit.
- **Effort — not per-child on this host.** The `Task` spawn exposes no per-child
  effort control. Some slugs encode effort in the id (`cursor-grok-4.6-high`).
  Record the configured `effort` as host evidence, but do not expect a separate
  per-child effort switch — per-role tiering degrades to **model-only** on
  Cursor. (The coordinator's own effort remains the user's session control.)

Set the worker's `model` from the `worker` entry, the verifier's from the
`verifier` entry, and the planner's from the `planner` entry before launching.

Model and effort stay bounded, provider-neutral `host_evidence` (INV-HOST-01):
they change cost and speed, never a route, role, lease, verification, or
completion, and the runtime never learns them (INV-RUNTIME-01). Record the
`(model, effort)` each role ran at with `attempt-evidence-record`; do NOT record
inference wall-clock — an attempt already carries `started_at` and `checked_at`
timestamps whose span is its duration, and the most precise per-role figure is in
the host's sub-agent transcript, so no coordinator stopwatch is needed.

Cursor permission prompts, authentication, memory, MCP settings, and transcripts
are host-local. They never replace a human approval or completion gate and never
enter workspace state. If Task/subagent creation is unavailable, report
`host-blocked` and keep the route read-only; never self-verify.
