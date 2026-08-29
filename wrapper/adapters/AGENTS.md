# Context Circuit workspace

This is a Context Circuit universal project workspace: an agent-oriented place to
hold Product Knowledge, create grounded plans, and execute them safely across one
or more Git repositories. Talk to it in ordinary language.

## How it works

A request flows as a normal conversation: the coordinator reads the relevant
Product Knowledge and plan material, drafts or reviews a readable plan, you
approve it in conversation, the workspace prepares isolated repository worktrees,
one worker implements the whole plan and commits each repository, an independent
verifier checks the latest commits, the worker repairs failures with new commits,
and you decide when the plan is complete.

Branches, worktrees, runtime records, and verifier setup are hidden. The plan,
the worker's handoff, the verifier's result, and your authority over completion
are not hidden.

## Safety spine

- Read first: `WORKFLOW.md`, `workspace.yaml`, and `context/INDEX.md`. Read only
  the plan and context units an action needs.
- Rule ownership is indexed in `wrapper/contracts/invariants.yaml`. One rule has
  one owner; do not add parallel policy to a skill or role file.
- Only an approved plan may execute. Approval and execution are explicit
  conversational actions, never confirmation cards or hidden tokens.
- One worker writes; one independent read-only verifier checks. If the verifier
  child cannot be created, the result is `host-blocked` — never self-verify.
- `sources/` is passive: read only exact request-named files; never scan all
  sources or sibling workspaces to fill a context gap.
- `plans/archive/` is outside normal context; read it only via explicit restore.
- Completion, pull requests, merge, push, deployment, archive, and
  cleanup are separate explicit human actions. Verification never implies any.
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
completion. A native child maps to the single worker or the independent verifier.

## Per-role model & effort

The coordinator may run the worker and verifier at a per-role `(model, effort)`
from an optional host-local, per-user, gitignored `role-tiering.local.yaml`,
grouped by host so each host names the models available on it. This is bounded
host evidence (`host_evidence`): it changes cost and speed, never a route, role,
lease, verification, completion, verifier independence, or the failure limit, and
the runtime never learns it (INV-RUNTIME-01). Reading the config does not apply
it — the host adapter sets the model on the child spawn; absent any config, the
adapter defaults apply. Full rules: `docs/role-tiering.md`.
