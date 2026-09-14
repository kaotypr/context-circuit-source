---
name: cc-dispatch
description: Dispatch Context Circuit planner, explorer, worker, or explicitly requested reviewer subagents with the configured host model and effort settings.
---

# Dispatch workspace subagents

The coordinator owns task selection and integration. Use subagents for bounded
work that benefits from delegation, including parallel tasks with distinct file
ownership. A small task may stay in the main session. Do not launch review unless
the user requests it. Main-session implementation checks remain ordinary tests.

Read `.context-circuit/docs/agents.md` for the current host's settings mechanism.
Use `agent settings` to resolve `.context-circuit/role-tiering.yaml`. The initial
model and effort for every role are `inherit`, as chosen by the maintainer. Honor
explicit user settings; never infer model quality order or silently escalate.

Before the first dispatch, or after settings change, run:

```sh
context-circuit --workspace <root> agent setup --host <codex|claude-code|cursor>
```

This installs native role definitions without replacing custom edits. Reload the
host's agent definitions if required. Do not claim a definition was loaded merely
because a file was written. If native roles are unavailable, use the host's live
spawn tool with the full resolved role prompt and settings, where supported.

For each bounded task, obtain its dispatch specification:

```sh
context-circuit --workspace <root> --json agent dispatch \
  --host <host> --role <explorer|planner|worker|reviewer> \
  --path <actual-working-directory> --task '<bounded assignment and references>'
```

Use `--review-requested` only to represent the actual user request for review.
It supplies no human approval by itself. Include the intended base/head revisions
and complete diff for a reviewer, along with success criteria and relevant code.

**Complete the dispatch:** call the host's subagent tool with the returned prompt,
role, working directory, and supported model/effort settings. The CLI returns
`launch_required: true`; it has not launched or completed an agent. Prefer a fresh
context, and always use a fresh independent reviewer context. Never reuse the
implementer's session as its independent reviewer.

- Codex: use the registered `cc_<role>` agent or the available spawn tool's role,
  model and reasoning-effort fields. For tools where full-history inheritance
  prevents overrides, choose a fresh context. Omit overrides for `inherit`.
- Claude Code: invoke the registered `cc-<role>` with the host's Agent tool.
  Model and effort are set in its native definition. Read-only roles have only
  Read/Glob/Grep; the coordinator supplies diff text because they cannot run Git.
- Cursor: invoke `cc-<role>` with its Task tool. The native definition sets model,
  optional model effort parameter, and read-only access.

Check requested settings against the host's actual available models/efforts. A
configured pair is a request, not proof of the model that ran. Report any rejected
or substituted setting; do not silently downgrade. If the host cannot apply a pin,
resolve the mismatch with the user. Inherited defaults need no extra confirmation.

Give workers explicit ownership and the actual prepared worktree paths. Tell them
they share the codebase and must preserve other agents' edits. Wait for results,
inspect changes, integrate in dependency order, and update the existing plan with
progress and limitations. Preserve interrupted work. Do not create automatic
repair loops, nested agent chains, execution records, or delivery side effects.
