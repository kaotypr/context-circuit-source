# Workspace member

Runs the product against a real project, alone or on a roster.

## [Grounding an intent and its open questions](../domains/intent-grounding.md)

1. **As a member, I want the outcome written down before any code is
   investigated,** so that I am approving what the agent understood rather than
   what I happened to type.
2. **As a member, I want unsettled decisions to reach me as numbered
   questions,** so that I can answer three of five and leave the rest open
   without losing my place.

## [Authorization boundaries](../domains/authorization.md)

1. **As a member, I want a request to execute to also cover the setup it
   obviously needs,** so that I am not asked permission four times for one
   decision I already made.
2. **As a member, I want push, pull request, merge, and deletion to stay mine,**
   so that nothing reaches anyone else's view of the project without me.
3. **As a member, I want a gate to name what it just made due,** so that a
   command reporting success does not read as a request finished.

## [Worktrees and environment reuse](../domains/worktrees-and-reuse.md)

1. **As a member, I want isolation that does not cost a dependency install,** so
   that I use it on the risky work instead of avoiding it exactly then.
2. **As a member, I want my dirty and untracked files left alone,** so that
   letting the product work in my checkout is not a gamble.

## [Deriving stacked plan order](../domains/plan-ordering.md)

1. **As a member, I want the order of several plans derived from what I
   recorded,** so that sequencing is not re-argued from memory each session.
2. **As a member, I want the recommendation to come with its reason and its
   cost,** so that overriding it is an informed choice rather than a guess.

## [Record identity and allocation bands](../domains/record-ids-and-bands.md)

1. **As a member on a roster, I want my record numbers not to collide with a
   teammate's,** so that two people can work offline and merge without
   renumbering anything.

## [Durable notes and the retrieval catalog](../domains/knowledge-notes.md)

1. **As a member, I want to retrieve the few notes that bear on my task,** so
   that starting work does not mean reading everything the project knows.
2. **As a member, I want most completions to record that nothing durable
   changed,** so that the knowledge circuit stays honest instead of producing
   an edit per plan.

Owner:

- `context-circuit-source@internal/workspace/` — the bookkeeping this actor
  relies on and the refusals that protect their checkout.
- `context-circuit-source@product/AGENTS.md.in` — the gates that keep the two
  decisions theirs.
