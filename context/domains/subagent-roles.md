# Subagent roles and host settings

## Roles are capabilities, not authority

| Role | Responsibility | Access |
| --- | --- | --- |
| explorer | Answer a specific codebase question with evidence | Read-only |
| planner | Investigate an approved intent and return a grounded plan | Read-only; the coordinator writes the plan |
| worker | Implement assigned work and run normal checks | Repository edits within its assignment |
| reviewer | Independently examine a diff on request | Read-only; no fixes, no follow-up |

The coordinator — the main session — owns task selection, preparation, waiting,
and integration. Delegation is a tool, not a ritual: a small task may stay in
the main session, and no tier or risk classification compels a spawn. A role
never carries ownership, consequence, or permission to proceed.

Which delegation repays its cost follows from shape. One planner runs per
intent, so nothing overlaps and dispatching one always costs latency; what it
buys is a large read kept out of the coordinator, which pays only when that read
would dwarf the plan it produces. A session that already holds the code gains
nothing, because a planner would re-derive what is present. Splitting a planner
per repository destroys the cross-repository order it exists to produce.
Explorers are the opposite shape — narrow, many, genuinely parallel — so several
independent questions are answered at once and none of that reading enters the
session.

A planner answers under fixed headings: verdict, approach, tasks and order,
risks and checks, and plan shape. Approach, tasks and order, and risks and
checks transcribe into the record.
Verdict and plan shape are answers to the coordinator rather than plan content:
both are spent the moment the records are created, and a plan recording that it
is a single plan states what its own existence already says. `not-feasible` is a complete answer that
ends planning with no record written, and a verdict that the outcome or success
criteria must change returns for renewed approval instead. Because read-only
roles run nothing, a check either role names is one it read, never one it saw
pass.

Plan shape being the planner's answer is why a planner is dispatched against an
approved intent and refuses a plan ID: numbering a plan first settles the split
it was asked to propose, and an intent that turns out to hold several plans has
no single ID to hand it. The same reasoning refuses an unapproved intent, since
planning one plans an outcome nobody agreed to.

## A brief is the whole context

A subagent starts with nothing, so the brief is read top to bottom and is
composed in that order: working directory, ownership, the record, the task, and
last what to return. Ownership reaches every role, because a directory named
without a word on what may be done in it leaves how to read it open, and a host
offering GUI automation and web search will occasionally answer that badly: read
and search the files directly, drive no other application, consult nothing
external. Every record it names is quoted in full, because an agent
sent to open a file searches the repository to find it and arrives having read
far more than the file. The quoted record is also the assignment, so the task
section appears only where one is given: a planner takes the whole intent and
refuses a task, a worker owning its whole plan needs none, and what remains is
a role no record assigns or one slice of a plan several workers share. For the same reason the returned prompt is launched
unmodified — a coordinator that retypes it in its own words drops the quoted
record, and the agent then re-derives everything the coordinator already knew.

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

## Shared agreement, local answer

The tiering is shared because a roster should agree on which role earns the
strongest settings. What a member can actually run is not shared: one runs a
different host, another pays for a different model. So a machine-local file
overrides the shared one, and the override is per host/role rather than
whole-file — a member who raises their reviewer keeps tracking every later change
the team makes to the other three, which a forked copy would silently stop doing.

Its host must already appear in the shared file, since a local answer responds to
a host the workspace supports rather than introducing one nobody else has. Both
files face the same validation, so a local file cannot hold a setting the shared
one would have been refused. The settings report names which pairs this machine
overrode: a merged view that cannot say which half supplied a value claims more
than it established.

## Native host mapping

Setup materializes four native role definitions under each host's agents
directory; those generated files are local and ignored while the settings that
produce them are shared. One workspace is opened in several hosts — an intent
written in one, a plan grounded in another, execution in a third — so setup
covers every host unless `--host` narrows it, and initialization runs it rather
than leaving a step between a new workspace and its first delegation. Seeding
those files from the template instead would ship content no inventory records,
which a later version reads as a customization and refuses to replace. A clone
still carries none of them, because they are ignored; dispatch reporting an
absent definition is what makes that visible. Codex expresses model, reasoning effort, and a sandbox
mode for read-only roles; Claude Code expresses model, effort, and a read-only
tool allowlist; Cursor expresses model with an optional effort parameter and a
read-only flag, and requires an explicit model whenever an effort is explicit.
Codex omits an inherited setting, because its configuration has no word for
inheriting; Claude Code and Cursor both accept `inherit` as a model value and
receive it written out, which says the same thing in a file a person may read. A
written setting is not an applied one: a Codex agent spawned against a read-only
role definition was observed running at full access, with that file's model,
effort, and instructions all in force and its sandbox mode alone dropped. Nor
would the sandbox have covered much, since it governs file writes rather than a
host's GUI automation or web search. So the brief states the conduct itself and
treats host enforcement as a bonus. A customized native file is preserved and
reported rather than overwritten. Because read-only roles cannot run Git, the
coordinator supplies diff text to them directly.

**Writing a file is not loading it.** The host may need a restart, and neither
the executable nor the dispatching skill may claim a definition was loaded
merely because it was written.

## Dispatch

The executable resolves settings, composes a brief, and returns the invocation
with a flag stating that a launch is still required. It has not launched
anything, and a specification is never evidence that an agent ran or completed.
Because role files are host-local, gitignored, and written only by setup, a
specification can name an agent type the host never registered; dispatch
therefore reports the definition's path, whether it is present, and the setup
run that writes it. Reporting rather than refusing keeps the prompt usable by a
live spawn tool where native roles are unavailable.
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

## Integration reads the report, not the diff

A worker reports the files it changed, the checks it ran, and their real outcome,
and the coordinator integrates from that report. Re-reading the diff and re-running
the checks as a routine audit repeats the expensive half of the work and is what
makes delegation cost more than doing it directly. The coordinator reads the diff
where integration needs it — a merge to resolve, or a report naming a conflict, a
failure, or an assumption — and a report of failing checks is carried into the
plan rather than repaired in a loop. Wanting a second pair of eyes on the change
is independent review, which is read-only, separately dispatched, and requested by
the user.

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

Owner:

- `context-circuit-source@internal/workspace/agents.go` — the roles, the
  settings that reach a host, and the composed brief.
- `context-circuit-source@template/role-tiering.yaml` — the shipped reasoning
  about which role earns which settings.
- `context-circuit-source@product/skills/cc-dispatch/SKILL.md` — the dispatch
  procedure a coordinator follows.
