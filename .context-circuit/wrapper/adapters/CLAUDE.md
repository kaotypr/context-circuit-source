# Claude Code adapter

@AGENTS.md

Follow the shared `AGENTS.md` and `WORKFLOW.md` contract. Claude Code is a thin
transport. The root Claude session is the coordinator (see `.context-circuit/agents/coordinator.md`)
for all normal conversation: orientation, context gathering, planning, review,
approval interpretation, direct collaboration, completion, archive/restore, and
delivery discussion.

A Claude Task/subagent maps only to the single bounded worker (`.context-circuit/agents/worker.md`)
for an execution or direct-collaboration session, the independent read-only
verifier (`.context-circuit/agents/verifier.md`) for one execution, or the planner
(`.context-circuit/agents/planner.md`) after intent approval.
Record provider-neutral `host_evidence` for the child; a host permission flag is
an observation, not authorization.

## Applying per-role model & effort

Owner: `.context-circuit/docs/role-tiering.md` (how to obtain the config, by
invoking the engine verb, from the workspace root). Read that file. Do not
restate the rule here.

When the config carries a `(model, effort)` for a role, obtain it by invoking
the engine verb, per the owner doc — never from a Task working directory —
then read the group for this host from `hosts.claude-code` and apply it —
which on Claude Code means **setting it on the spawn, not just recording it**.
A missing file in an isolated working copy is not an absent config. Concretely:

- **Model — set it on the spawn.** Pass the role's configured model id to the
  worker's or verifier's `Task`/subagent spawn as its `model` parameter. A child
  launched without a `model` inherits the coordinator's session model, so an unset
  spawn is the whole tier being silently ignored. Set the worker's `model` from the
  `worker` entry, the verifier's from the `verifier` entry, and the planner's from
  the `planner` entry before launching.
- **Effort — session-level on this host; not per-child.** The `Task`/subagent
  spawn exposes no per-child effort control, so a worker and a verifier launched
  in the same session cannot run at different efforts here. Record the configured
  `effort` as host evidence, but do not expect a per-child effort switch — per-role
  tiering degrades to **model-only** on Claude Code. (The coordinator's own effort
  remains the user's session control.)

Model and effort stay bounded, provider-neutral `host_evidence` (INV-HOST-01):
they change cost and speed, never a route, role, lease, verification, or
completion, and the runtime never learns them (INV-RUNTIME-01). Record the
`(model, effort)` each role ran at with `attempt-evidence-record`; do NOT record
inference wall-clock — an attempt already carries `started_at` and `checked_at`
timestamps whose span is its duration, and the most precise per-role figure is in
the host's sub-agent transcript, so no coordinator stopwatch is needed.

Claude permission prompts, print mode, authentication, memory, and MCP settings
are host-local. They never replace a human approval or completion gate and never
enter workspace state. If Task/subagent creation is unavailable, report
`host-blocked` and keep the route read-only; never self-verify.

## Host-native routes

Claude discovers Context Circuit child roles, standing rules, and slash skills
through the committed `.claude/` tree. Those files are thin routes, not a second
policy:

- Agents: `.claude/agents/{worker,verifier,planner}.md` → Read
  `.context-circuit/agents/<role>.md`
- Rules: `.claude/rules/` → owning invariants (role-tiering spawn, commit
  convention) and this file where the host requires it
- Skills: `.claude/skills/cc-*/SKILL.md` are thin routes. Read
  `.agents/skills/<name>/SKILL.md` — that file is the owner (INV-SKILL-01).
  Do not copy the skill body into the Claude tree. The coordinator still
  resolves skills by path; the Claude skills tree is how `/cc-*` is discovered.

`.claude/` is committed workspace integration. Personal Claude state
(`~/.claude/`, `.claude/settings.local.json`, transcripts, credentials) stays
host-local and never enters workspace state. Native folders grant no route,
role, or authority that the read-as-procedure path does not already carry.
