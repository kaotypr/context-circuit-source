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
context-circuit-cli --workspace <root> agent setup --host <codex|claude-code|cursor>
```

This installs native role definitions without replacing custom edits. Reload the
host's agent definitions if required. Do not claim a definition was loaded merely
because a file was written. If native roles are unavailable, use the host's live
spawn tool with the full resolved role prompt and settings, where supported.

For each bounded task, obtain its dispatch specification:

```sh
context-circuit-cli --workspace <root> --json agent dispatch \
  --host <host> --role <explorer|planner|worker|reviewer> \
  --plan <plan-id> --path <actual-working-directory> \
  --task '<bounded assignment and references>'
```

Omit `--plan` when the task is not plan work. Add `--shared` only when several
workers use one worktree.

Use `--review-requested` only to represent the actual user request for review.
It supplies no human approval by itself. Include the intended base/head revisions
and complete diff for a reviewer, along with success criteria and relevant code.

Check `definition_installed` on the returned specification. When it is false the
host has no native role of that name, and `setup_required` names the run that
writes one. Install it, then reload the host's definitions if required; where the
current session cannot pick up a newly written definition, fall back to the live
spawn tool with the returned prompt rather than abandoning the dispatch.

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

## Two kinds of parallel work

Parallelism happens on two axes and the correct brief is opposite in each. One
worker per plan owns a whole worktree: separate plans get separate working copies
and branches, so they cannot touch each other's files, and a sibling's work
arrives later through an integration merge. Several workers inside one plan's
worktree genuinely share files and must preserve each other's edits.

Pass `--shared` only for the second case. The default brief tells a worker it is
the sole owner of its directory; `--shared` replaces that with shared-ownership
language. Claiming the wrong one either invites a worker to guess at edits it
cannot see, or lets it overwrite edits it can.

Pass `--plan` for plan work. The CLI then states the plan's repositories, its
dependencies, that their completed work is already in this branch's ancestry so
it need not be reimplemented, and that a concurrent sibling plan is invisible and
must not be guessed at. A worker reporting an interface it assumes a sibling may
also be changing is supplying information, not failing; carry it to the
integration merge rather than stopping the run.

Give workers explicit ownership and the actual prepared worktree paths. Tell them
they share the codebase and must preserve other agents' edits. Wait for results,
inspect changes, integrate in dependency order, and update the existing plan with
progress and limitations. Preserve interrupted work. Do not create automatic
repair loops, nested agent chains, execution records, or delivery side effects.
