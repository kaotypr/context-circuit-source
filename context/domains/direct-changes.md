# Changing something directly

The lane a person takes when their request is already the specification.

## Why the path exists to be skipped

The intent gate protects against one specific risk: an agent's reading of a
request drifting from what the person meant. Writing the outcome down and having
it approved is how that drift is caught before any code is touched.

Where the request already fixes the outcome — rename this flag, raise this
timeout, correct this message — nothing was derived, so there is nothing derived
to approve. The four steps then buy protection against a risk that is not
present, and the cost is real: a candidate released with only the planned path
made small changes annoying enough to avoid, which is how a product stops being
used for the work it should be easiest at.

## Two lanes, one difference

| | Planned path | Direct lane |
| --- | --- | --- |
| What approves it | The person approves the written intent | The request itself, because nothing was derived |
| Records | An intent and its plans | None; no ID is allocated |
| Where it runs | A prepared worktree on its own branch | The bound checkout, on its current branch |
| What it leaves | Committed work a dependent plan can start from | An uncommitted diff, unless a commit is authorized |
| Knowledge | Reconciled at completion | Reconciled the same way, asked for by repository |

Every other gate is identical. Committing in a bound checkout, pushing, opening
a pull request, merging, deploying, publishing, and deleting each still need a
person, exactly as they do on the planned path.

## The person routes, not the agent

The agent may offer the lane where a request is already its own specification. It
never takes it unasked.

That asymmetry is deliberate, and it is the same failure worktrees already had:
an agent weighing ceremony against convenience mid-task reliably chooses less
ceremony, which is why [preparing a worktree](worktrees-and-reuse.md) stopped
being a judgment call. Leaving the routing with the person keeps the bypass a
decision someone made rather than a habit an agent formed.

## Where it stops

The lane fits a change whose outcome is settled. It stops as soon as that is no
longer true:

- the change reaches a second repository, or must be coordinated with planned
  work;
- settling it needs a decision the person has not made;
- the fix reveals the real problem is somewhere else;
- the reading needed to make it safely is large enough to have been a plan.

Stopping means saying what was found and offering the intent, with correct work
left in place. A change that grew past the lane and was carried through anyway is
an unrecorded feature nobody approved — the exact outcome the gate exists to
prevent, reached by the route that was supposed to be cheap.

## The reconcile half

Completion names the knowledge a plan made due. A direct change has no record to
name it, so the repository is asked instead: the catalog entries claiming it come
back together with the same outstanding act. The obligation never depended on a
record existing — it depends on code having changed — and this is the one trace a
direct change leaves besides its diff.

Owner:

- `context-circuit-source@product/skills/cc-direct/SKILL.md` — the procedure, its
  refusals, and the reconcile step.
- `context-circuit-source@product/AGENTS.md.in` — the resident routing rule that
  keeps the choice with the person.
- `context-circuit-source@internal/workspace/records.go` `CatalogEntriesFor` —
  the entries returned for a repository with no record to scope them.
