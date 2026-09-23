# Deriving stacked plan order

Running a selected set of plans without guessing what may run beside what.

## Derivation, not execution

Ordering reads recorded dependencies and completion dates for all active plans,
an intent's linked plans, or exact IDs selected with repeatable `--plan`. It reports the waves
and their concurrency, each plan's start reference per repository, the
integration merges a fan-in needs, which plans in one wave share a repository, a
strategy summary with the cost of each shape, the completed and blocked plans,
and how many remain. The result names its selected plan IDs. An unfinished
dependency outside that selection blocks its dependents; it is not added.

It **runs nothing, merges nothing, and reserves nothing**, and says so, because
a report that looks like a plan of record invites being treated as one.

## Two shapes with their costs

| | Wall clock | Merges | Fits |
| --- | --- | --- | --- |
| Waves | Independent plans overlap | One integration merge per fan-in | Several repositories with independent plans |
| Linear chain | Strictly serial | None | A single repository with any fan-in |

The recommendation is mechanical and reported with its reason:

- A parallel peak of one recommends the chain, because it is the same work.
- A single repository with any fan-in recommends the chain, because waves buy
  wall clock and cost merges exactly where they are most likely to conflict.
- Anything else recommends waves, because independent plans in separate
  repositories overlap without a merge at all.

Present the recommendation with its cost and honor an explicit choice. The
separate request to execute the presented plans authorizes the run; use the
recommended shape when no shape was chosen.

## The run

After the person requests execution, prepare each worktree from the reported
start, perform any reported integration merge with ordinary Git, dispatch the
workers, and wait. Completion is a later person-requested act.

A wave is judged from what each worker reports: the files it changed, the checks
it ran, and their real outcome. Re-reading the diff and re-running those checks
as a routine audit repeats the expensive half of the work, which is what makes a
delegated wave slower than doing it directly.

The trade is deliberate and worth naming. The coordinator is trusting a report it
did not reproduce, so the worker is told that a check it did not run is one
nobody ran, and that a failure is reported plainly rather than worked around.

Read the diff where integration needs it — a merge to resolve, or a report
naming a conflict, a failure, or an assumption — not to satisfy yourself that
work already reported was really done.

Ordering reads previously recorded completion to judge dependencies outside the
run. Within a run, follow the selected waves and actual Git ancestry; do not
stamp completion merely to advance the next wave.

## Conflicts and stops

A conflict from an integration merge inside the run is resolved directly: both
sides belong to the selected run and both their records and diffs are
available. The resolution preserves both behaviors, then runs the repository's
checks — a resolved merge is not trusted until they pass — and is recorded so it
can be audited. A conflict against anything outside the run, or one where
preserving both sides is impossible, is a stop.

Every stop preserves all work: failing checks after implementation or after a
resolved merge; a conflict outside the run or one needing a decision; a worker
that cannot complete or returns blocked; worktree preparation that refuses;
skipped dependency reuse whose fallback setup fails; or an order reporting a
cycle, an unknown repository, or a broken link. Completed work is never unwound,
and the report names what finished, what is stuck, and what is held behind it.

## Where a chain ends, and why delivery reads that

The derivation also marks, per plan, the repositories whose chain ends there. A
chain end's branch already contains the plans it was prepared from, so it is the
only branch that needs a request: delivering its predecessors as well would offer
the same commits twice.

The mark is per repository rather than per plan, because a plan can end a chain in
one repository and sit mid-chain in another — it is delivered where nothing
follows it and skipped where a successor's branch already carries it. One
repository can also hold two chain ends and take two requests.

Reading the last wave instead is the mistake this replaces. Waves answer
readiness, not completion: a plan nothing depends on sits in the first wave and
still ends its chain, so anything reasoning from waves drops it in silence.

## The coarseness the report admits

Shared-repository reporting reflects repositories, not paths, because plan
records name repositories. It cannot tell neighbouring modules from the same
file, so it is a signal for choosing a shape rather than a collision predictor.
Start references also name a predecessor's branch: once that work is merged and
its branch is gone the reference needs adjusting by hand, since the derivation
reads records and records do not know what a provider merged.

Owner:

- `context-circuit-source@internal/workspace/order.go` — the derivation, its
  waves, start references, and integration merges.
