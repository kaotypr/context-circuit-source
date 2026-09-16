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
Use `agent settings` to resolve the tiering in force: the shared
`.context-circuit/role-tiering.yaml`, with any host/role this machine overrode in
`.context-circuit/role-tiering.local.yaml` replacing it. The report names those
overrides. The initial model and effort for every role are `inherit`, as chosen by
the maintainer. Honor explicit user settings; never infer model quality order or
silently escalate.

Initialization installs role definitions for every host. Run setup again after
settings change, or in a clone, where the ignored definitions never arrived:

```sh
context-circuit-cli --workspace <root> agent setup
```

With no `--host` it covers every host, since one workspace is opened in several.

This installs native role definitions without replacing custom edits. Reload the
host's agent definitions if required. Do not claim a definition was loaded merely
because a file was written. If native roles are unavailable, use the host's live
spawn tool with the full resolved role prompt and settings, where supported.

For each bounded task, obtain its dispatch specification:

```sh
context-circuit-cli --workspace <root> --json agent dispatch \
  --host <host> --role <explorer|planner|worker|reviewer> \
  --plan <plan-id> --path <actual-working-directory>
```

The quoted record is the assignment. A planner plans the whole intent and refuses
`--task`; a worker owning its whole plan needs none, and adding one only restates
the plan it was already handed. Add `--task '<the assignment>'` where no record
assigns the role — an explorer's question, a reviewer's criteria — or to name the
single slice a worker owns when several share a plan, which is said nowhere else.

Dispatch a planner against `--intent <intent-id>`, never `--plan`: the planner
returns the plan shape, so numbering a plan first settles the split it was asked
to propose, and an intent that turns out to hold several plans has no single ID
to pass. Create the plan records afterwards, from what it returns. Omit both when
the task is not record work. Add `--shared` only when several workers use one
worktree.

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

**Pass `prompt` unmodified.** It is the whole brief: working directory, ownership,
the quoted intent or plan body, the task, and what to return. Appending a line is
allowed; rewriting, summarizing, or retyping it from memory is not. A shortened
brief looks equivalent and is not — the quoted record is the copy that stops the
agent from searching the repository for a file it was never given, and dropping it
costs far more time than the paragraph saved. If a detail is missing, put it in
`--task` and dispatch again rather than editing the returned text.

- Codex: run `agent dispatch` and the spawn in one exec cell, so the brief goes
  from one to the other without being retyped or spending context on the way:

  ```js
  const out = await tools.exec_command({cmd: "<the dispatch command above>"});
  const spec = JSON.parse(out.output);
  const {agent_id} = await tools.multi_agent_v1__spawn_agent({
    agent_type: spec.agent_type, message: spec.prompt,
  });
  ```

  The installed role definition carries the resolved model and effort, so pass
  neither here: a pinned setting there cannot be overridden on the spawn, and an
  `inherit` role is one meant to take the parent's model. Add `model` and
  `reasoning_effort` only on the fallback path where no native definition
  exists, omitting either whose setting is `inherit`. Wait on `agent_id` with
  `multi_agent_v1__wait_agent`, and release it with
  `multi_agent_v1__close_agent` once the result is in hand, since a completed
  agent stays open and counts against the concurrency limit until closed.

  A wait is a poll loop by construction: the exec cell parks long before
  `timeout_ms` elapses, and every resume is a full round-trip carrying the
  whole session, so the cost of delegating grows with how long the delegated work
  takes. Make each round-trip cover as much of the wait as the host allows — a
  `timeout_ms` measured in minutes, and the longest yield permitted on the resume
  rather than the default. For tools where full-history inheritance prevents
  overrides, choose a fresh context.
- Claude Code: invoke the registered `cc-<role>` with the host's Agent tool.
  Model and effort are set in its native definition. Read-only roles have only
  Read/Glob/Grep; the coordinator supplies diff text because they cannot run Git.
- Cursor: invoke `cc-<role>` with its Task tool. The native definition sets model,
  optional model effort parameter, and read-only access. A read-only role there
  is under the same constraint as on Claude Code: supply diff text rather than
  expecting it to run Git.

Check requested settings against the host's actual available models/efforts. A
configured pair is a request, not proof of the model that ran. Report any rejected
or substituted setting; do not silently downgrade. If the host cannot apply a pin,
resolve the mismatch with the user. Inherited defaults need no extra confirmation.

## When to dispatch a planner or an explorer

Delegating a plan's code investigation is a judgment about size, not a step.
Dispatch a planner when the reading would be substantially larger than the plan
it produces: an intent spanning several repositories, code this session has not
read, or a read that would take a serious part of the remaining context.
Investigate in the session when it already holds that code, since a planner would
re-derive what is present and add only latency.

One planner runs per intent; splitting it per repository destroys the
cross-repository order it exists to produce. Dispatch it against the intent
rather than a plan: the plan shape is what it returns, so create the plan records
from its answer rather than numbering one first and handing it over.
`.agents/skills/cc-plan/SKILL.md` covers what to write from that answer.

Dispatch explorers when several independent codebase questions stand between this
session and a grounded intent or plan. Each is narrow and read-only, they run in
parallel, and their reading never enters this session. Neither role can run
anything, so a check either one names is one it read, not one it saw pass.

## Two kinds of parallel work

Parallelism happens on two axes and the correct brief is opposite in each. One
worker per plan owns a whole worktree: separate plans get separate working copies
and branches, so they cannot touch each other's files, and a sibling's work
arrives later through an integration merge. Several workers inside one plan's
worktree genuinely share files and must preserve each other's edits.

A plan naming several repositories has a worktree and a branch in each, so it
takes one sole-owner worker per repository rather than one for the plan. Each is
briefed with the whole plan and sees only its own working directory.

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
then integrate in dependency order and report what they reported. Nothing appends
a progress log to the plan; completion is a person's request and carries the
result. Preserve interrupted work. Do not create automatic repair loops, nested
agent chains, execution records, or delivery side effects.

**Do not re-run a worker's checks or re-read its diff to satisfy yourself.** The
worker reports the commands it ran and their real outcome; read that report and
integrate from it. Repeating the reading and the test run is the expensive half
of the work done twice, and it is why delegation saves nothing. Read the diff
when integration needs it — resolving a merge, or a report that names a conflict,
a failure, or an assumption — not as a routine audit of work already reported. A
worker that reports a failure is reporting, not failing: carry it forward. If you
want a second pair of eyes on the change, that is independent review, which is
read-only, separately dispatched, and only on the user's request.
