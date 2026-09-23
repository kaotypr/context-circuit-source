# Subagent roles and host settings

What a delegated agent may do, and what it is handed to do it with.

## Roles are capabilities, not authority

| Role | Responsibility | Access |
| --- | --- | --- |
| explorer | Answer a specific codebase question with evidence | Read-only |
| planner | Investigate an approved intent or specified standalone outcome and return a grounded plan | Read-only; the coordinator writes the plan |
| worker | Implement assigned work and run normal checks | Repository edits within its assignment |
| reviewer | Independently examine a diff on request | Read-only; no fixes, no follow-up |

The coordinator — the main session — owns task selection, preparation, waiting,
and integration. Delegation is a tool, not a ritual: a small task may stay in
the main session, and no tier or risk classification compels a spawn. A role
never carries ownership, consequence, or permission to proceed.

Which delegation repays its cost follows from shape, and the two read-only roles
have opposite ones:

- **A planner** runs once per planning request, so nothing overlaps and dispatching one
  always costs latency. What it buys is a large read kept out of the
  coordinator, which pays only when that read would dwarf the plan it produces.
  A session already holding the code gains nothing, because a planner re-derives
  what is present, and splitting one per repository destroys the
  cross-repository order it exists to produce.
- **Explorers** are narrow, many, and genuinely parallel. Several independent
  questions are answered at once, and none of that reading enters the session.

A planner answers under fixed headings, and where each one lands is fixed too:

| Heading | Where it goes |
| --- | --- |
| Approach | transcribed into the plan record |
| Tasks and order | transcribed into the plan record |
| Risks and checks | transcribed into the plan record |
| Verdict | an answer to the coordinator, spent once the records exist |
| Plan shape | an answer to the coordinator, spent once the records exist |

A plan recording that it is a single plan states what its own existence already
says. `not-feasible` is a complete answer that ends planning with no record
written, and a verdict that the outcome or success criteria must change returns
for renewed approval instead. Because read-only roles run nothing, a check
either role names is one it read, never one it saw pass.

Plan shape being the planner's answer is why a planner is dispatched against an
approved intent or an explicit task specifying a standalone outcome, and refuses
a plan ID: numbering a plan first settles the split it was asked to propose.
An unapproved intent is refused.

## A brief is the whole context

A subagent starts with nothing, so the brief is read top to bottom and composed
in that order:

1. **The working directory** it stands in.
2. **Ownership** — what may be done there. It reaches every role, because a
   directory named without a word on what may be done in it leaves how to read
   it open, and a host offering GUI automation and web search will occasionally
   answer that badly: read and search the files directly, drive no other
   application, consult nothing external.
3. **The record body**, quoted for record-based work, so an agent can read its
   assignment without searching the repository. A folder plan's supporting
   files are instead listed by workspace-relative path for its repository.
4. **The task**, where one is given. The quoted record is itself the assignment,
   so a planner takes the whole intent and refuses a task and a worker owning
   its whole plan needs none; what remains is a role no record assigns, or one
   slice of a plan several workers share.
5. **The language boundary**, where a record is written in one. A subagent that
   receives a record in another language has to be told which of it is prose that
   stays in that language and which is structure that stays English, or it
   answers in whichever the record led with.
6. **What the role must return.**

The returned prompt is launched unmodified for the same reason: a coordinator
that retypes it in its own words drops the quoted record, and the agent then
re-derives everything the coordinator already knew.

## Settings without a model catalog

Role tiering is a concrete model-and-effort pair per role and host, every pair
shipping as inherit so the host chooses its own defaults. **No model catalog and
no cost ladder ships**, for two structural reasons: any list of model names goes
stale within months, and a hardcoded quality ordering invites silent escalation
by an agent that "knows" one model outranks another.

What ships instead is reasoning a reader can apply to whatever their host
offers:

- A planner runs once per intent and everything downstream inherits its
  mistakes; a reviewer has to catch what another model already convinced itself
  was fine. Those two repay the strongest settings available.
- A worker is where token volume goes, and an explorer is high-count and
  shallow.
- Where a host offers several model families, a reviewer drawn from a different
  family than the worker has different blind spots.

Then check what the host actually supports.

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
than leaving a step between a new workspace and its first delegation.

Seeding those files from the template instead would ship content no inventory
records, which a later version reads as a customization and refuses to replace.
A clone still carries none of them, because they are ignored, and dispatch
reporting an absent definition is what makes that visible.

Each host expresses the same pair in its own words:

| Host | Expresses |
| --- | --- |
| Codex | Model, reasoning effort, and a sandbox mode for read-only roles. An inherited setting is omitted, because its configuration has no word for inheriting |
| Claude Code | Model, effort, and a read-only tool allowlist. `inherit` is accepted as a model value and written out |
| Cursor | Model with an optional effort parameter and a read-only flag. An explicit effort requires an explicit model, and `inherit` is written out |

Writing `inherit` out rather than omitting it says the same thing in a file a
person may read.

A written setting is not an applied one. A Codex agent spawned against a
read-only role definition was observed running at full access, with that file's
model, effort, and instructions all in force and its sandbox mode alone dropped.
Nor would the sandbox have covered much, since it governs file writes rather than
a host's GUI automation or web search. So the brief states the conduct itself and
treats host enforcement as a bonus.

A customized native file is preserved and reported rather than overwritten.
Because read-only roles cannot run Git, the coordinator supplies diff text to
them directly.

**Writing a file is not loading it.** The host may need a restart, and neither
the CLI nor the dispatching skill may claim a definition was loaded merely
because it was written.

## Dispatch

The CLI resolves settings, composes a brief, and returns the invocation with a
flag stating that a launch is still required. It has not launched anything, and
a specification is never evidence that an agent ran or completed.

Because role files are host-local, gitignored, and written only by setup, a
specification can name an agent type the host never registered. Dispatch
therefore reports the definition's path, whether it is present, and the setup
run that writes it; reporting rather than refusing keeps the prompt usable by a
live spawn tool where native roles are unavailable. The dispatching skill calls
the host's own subagent tool, waits, and integrates. Prefer a fresh context, and
always use a fresh independent one for a reviewer.

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
be guessed at. A multi-repository plan assigns each worker a repository; the
brief lists only its assigned and shared required files and the absolute
workspace root they resolve against.

That ancestry claim is only true because a worker commits: the brief ends by
requiring the work be committed on the branch it was given, so what a worker
leaves behind is something the next plan can start from and an integration merge
can merge. Pushing, opening a pull request, and merging into a base branch remain
outside the role. A worker reporting an interface it assumes a sibling may also
be changing is supplying information, not failing; that goes to the integration
merge rather than stopping the run.

## Integration reads the report, not the diff

A worker reports the files it changed, the checks it ran, and their real outcome,
and the coordinator integrates from that report. Re-reading the diff and
re-running the checks as a routine audit repeats the expensive half of the work,
which is what makes delegation cost more than doing it directly.

The coordinator reads the diff where integration needs it — a merge to resolve,
or a report naming a conflict, a failure, or an assumption — and a report of
failing checks is carried into the plan rather than repaired in a loop. Wanting a
second pair of eyes on the change is independent review, which is read-only,
separately dispatched, and requested by the user.

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
