# Read-only next-work recommendation

`whats-next` recommends one executable task, or the smallest concrete action
needed to make known work executable. It never claims a task, edits a plan,
creates runtime records, branches, or worktrees, and never invokes an activity
provider. Explicit execution remains a separate `run-task` request.

## Local approved plans

Run the deterministic selector from the wrapper root:

```bash
node .agents/bin/kao.mjs whats-next
```

It validates every real directory under `context/plans/`. Plan work items keep
live state out of plan Markdown. A root item can be ready when its plan is
approved, its `area` resolves to a registered repository, the plan contains
acceptance criteria, and repository access is available. A dependent item is
blocked until its prerequisite completion can be verified. Draft plans yield
an approval action rather than implementation work.

The JSON output matches
`.agents/contracts/whats-next-result.schema.json`. Ranking is stable: urgent
work, actionable in-progress work, review/verification/CI failures,
highest-priority ready source tasks, then dependency-ready approved-plan items.
Priority breaks ties within a class; candidate ID is the final stable tie-break.

## Deterministic fake activity source

Tests and manual proofs may supply a JSON file matching
`.agents/contracts/fake-activity-source.schema.json`:

```bash
node .agents/bin/kao.mjs validate --schema fake-activity-source /tmp/activity.json
node .agents/bin/kao.mjs whats-next --activity-fixture /tmp/activity.json
```

The fixture may add read-only candidates or facts for a plan `work_id`, such as
completed dependencies, current ownership, urgency, or priority. It must never
contain credentials. It is not stored, claimed, synchronized, or treated as an
activity-provider SDK.

## Manual host proof

In Codex, invoke `$whats-next`. In Claude Code, invoke `/whats-next`. Use the
same wrapper state and, when testing activity ordering, the same fixture path.
Confirm both hosts delegate to `.agents/skills/whats-next/SKILL.md`, produce the
same recommendation and alternatives, cite the same readiness evidence, and
say that no state changed. Before and after each run, compare `git status
--short`, `git worktree list --porcelain`, and the relevant plan files. No task
claim, runtime directory, branch, worktree, commit, push, or external write may
appear.
