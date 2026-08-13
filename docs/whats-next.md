# Read-only next-work recommendation

`whats-next` recommends one executable task, or the smallest concrete action
needed to make known work executable. It never claims a task, edits a plan,
creates runtime records, branches, or worktrees, and never invokes an activity
provider. Explicit execution remains a separate `run-task` request.

## Local approved plans

Run the deterministic selector from the wrapper root:

```bash
node .agents/bin/cc.mjs whats-next
```

It validates every real directory under `context/plans/`. Plan work items keep
live state out of plan Markdown. A root item can be ready when its plan is
approved, its explicit repository key resolves to a registered repository, the plan contains
acceptance criteria, and repository access is available. A dependent item is
blocked until its prerequisite completion can be verified. Draft plans yield
an approval action rather than implementation work.
The descriptive `area` never participates in repository resolution. Invalid or
legacy-ambiguous plans are skipped with an actionable warning rather than
guessed from area prefixes.

## Read-only work-state projection

Plan intent is reconciled with four source classes, in descending precedence:

1. Git-tracked, structurally valid completed-work contributions;
2. validated closeout records associated to the same run, work ID, and repository;
3. validated plan-linked runtime manifests whose task brief confirms the run,
   work ID, sole repository, plan reference, and current approval version/digest;
4. configured read-only activity facts.

Precedence determines the projected state only when the evidence agrees. Any
different reported state is a contradiction and returns a `reconcile` action
with every state source cited; it never starts implementation or mutates a
source. A contribution cannot establish completion by itself: its required run
reference must resolve to that validated plan runtime, and both the manifest and
matching closeout record must associate the exact contribution path and outcome.
Untracked or ambiguously associated contributions are not durable work-state
evidence.

Runtime `passed` evidence becomes a `review` action. Closing or prepared
closeout evidence becomes a `closeout` action. Closed/merged and cancelled or
abandoned outcomes are excluded. Malformed, symlinked, planless, and unrelated
runtime material is ignored safely; malformed relevant evidence and symlinked
traversal entries are reported in `warnings`.

The JSON output matches
`.agents/contracts/whats-next-result.schema.json`. Reconciliation actions rank
first because execution would risk duplication. Otherwise ranking is stable:
urgent work, actionable in-progress or closeout work, review/verification/CI
failures, highest-priority ready source tasks, then dependency-ready
approved-plan items. Priority breaks ties within a class; candidate ID is the
final stable tie-break.

## Deterministic fake activity source

Tests and manual proofs may supply a JSON file matching
`.agents/contracts/fake-activity-source.schema.json`:

```bash
node .agents/bin/cc.mjs validate --schema fake-activity-source /tmp/activity.json
node .agents/bin/cc.mjs whats-next --activity-fixture /tmp/activity.json
```

The fixture may add read-only candidates or facts for a plan `work_id`, such as
completed dependencies, current ownership, urgency, or priority. It must never
contain credentials. It is not stored, claimed, synchronized, or treated as an
activity-provider SDK.

## Manual host proof

In Codex, invoke `$cc-whats-next`. In Claude Code, invoke `/cc-whats-next`. Use the
same wrapper state and, when testing activity ordering, the same fixture path.
Confirm both hosts delegate to `.agents/skills/cc-whats-next/SKILL.md`, produce the
same recommendation and alternatives, cite the same readiness evidence, and
say that no state changed. Before and after each run, compare `git status
--short`, `git worktree list --porcelain`, and the relevant plan files. No task
claim, runtime directory, branch, worktree, commit, push, or external write may
appear.
