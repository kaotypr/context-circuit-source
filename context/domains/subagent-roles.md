# Subagent roles and host settings

## Roles are capabilities, not authority

| Role | Responsibility | Access |
| --- | --- | --- |
| explorer | Answer a specific codebase question with evidence | Read-only |
| planner | Investigate an approved intent and return a grounded plan | Read-only; the coordinator writes the plan |
| worker | Implement assigned work and run normal checks | Repository edits within the assigned task |
| reviewer | Independently examine a diff on request | Read-only; no fixes, no follow-up |

The coordinator — the main session — owns task selection, preparation, waiting,
and integration. Delegation is a tool, not a ritual: a small task may stay in
the main session, and no tier or risk classification compels a spawn. A role
never carries ownership, consequence, or permission to proceed.

## Settings without a model catalog

Role tiering is a concrete model-and-effort pair per role and host, every pair
shipping as inherit so the host chooses its own defaults. **No model catalog and
no cost ladder ships**, for two structural reasons: any list of model names goes
stale within months, and a hardcoded quality ordering invites silent escalation
by an agent that "knows" one model outranks another.

What ships instead is reasoning a reader can apply to whatever their host
offers: a planner runs once per intent and everything downstream inherits its
mistakes, and a reviewer has to catch what another model already convinced
itself was fine, so those two repay the strongest settings available; a worker
is where token volume goes and an explorer is high-count and shallow; and where
a host offers different model families, a reviewer drawn from a different family
than the worker has different blind spots. Then check what the host actually
supports.

A configured pair is a **request**, not proof of what ran. A rejected or
substituted setting is reported, never silently downgraded.

## Native host mapping

Setup materializes four native role definitions under the host's agents
directory; those generated files are local and ignored while the settings that
produce them are shared. Codex expresses model, reasoning effort, and a sandbox
mode for read-only roles; Claude Code expresses model, effort, and a read-only
tool allowlist; Cursor expresses model with an optional effort parameter and a
read-only flag, and requires an explicit model whenever an effort is explicit.
Inherit omits the override entirely. A customized native file is preserved and
reported rather than overwritten. Because read-only roles cannot run Git, the
coordinator supplies diff text to them directly.

**Writing a file is not loading it.** The host may need a restart, and neither
the executable nor the dispatching skill may claim a definition was loaded
merely because it was written.

## Dispatch

The executable resolves settings, composes a brief, and returns the invocation
with a flag stating that a launch is still required. It has not launched
anything, and a specification is never evidence that an agent ran or completed.
The dispatching skill calls the host's own subagent tool, waits, and integrates.
Prefer a fresh context, and always use a fresh independent one for a reviewer.

### Two axes of parallelism

Getting this backwards is expensive in opposite directions. One worker per plan
**owns** its whole worktree: separate plans get separate working copies and
branches, cannot touch each other's files, and a sibling's work arrives later
through an integration merge. Several workers inside one plan's worktree
genuinely **share** files and must preserve each other's edits. A shared flag
selects the second; claiming the wrong one either invites a worker to guess at
edits it cannot see, or lets it overwrite edits it can.

Plan context adds the repositories, the dependencies, the fact that those
dependencies' work is already in the branch's ancestry and need not be
reimplemented, and the fact that a concurrent sibling is invisible and must not
be guessed at. A worker reporting an interface it assumes a sibling may also be
changing is supplying information, not failing; that goes to the integration
merge rather than stopping the run.

## Independent review

Review is manually requested, read-only, and usually follows pull-request
creation; it can also run afterwards as an audit. The reviewer receives the
requested diff, exact base and head revisions, relevant surrounding code, and
the success criteria, and returns findings with locations and limitations. It
modifies no code, dispatches no fixes, and posts nothing externally unless
asked. Tests that change files are implementation, not read-only review.

What was removed is the compulsion around it: review never starts during
execution, is never triggered by a classification, never blocks a pull request,
delivery, or completion, and never feeds an automatic repair loop. Where
independence is unavailable the agent says so and offers an ordinary review,
never calling the implementing session's own inspection independent.

Owner: `context-circuit-source@internal/workspace/agents.go`;
`context-circuit-source@template/role-tiering.yaml` carries the shipped
reasoning; dispatch procedure is
`context-circuit-source@product/skills/cc-dispatch/SKILL.md`.
