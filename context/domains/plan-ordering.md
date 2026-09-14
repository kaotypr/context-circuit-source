# Deriving stacked plan order

Running several plans of one intent without guessing what may run beside what.

## Derivation, not execution

Ordering reads recorded dependencies and completion dates and reports the waves
and their concurrency, each plan's start reference per repository, the
integration merges a fan-in needs, which plans in one wave share a repository, a
strategy summary with the cost of each shape, the completed and blocked plans,
and how many remain.

It **runs nothing, merges nothing, and reserves nothing**, and says so, because
a report that looks like a plan of record invites being treated as one.

## Two shapes with their costs

| | Wall clock | Merges | Fits |
| --- | --- | --- | --- |
| Waves | Independent plans overlap | One integration merge per fan-in | Several repositories with independent plans |
| Linear chain | Strictly serial | None | A single repository with any fan-in |

The recommendation is mechanical and reported with its reason: a parallel peak
of one recommends the chain because it is the same work; a single repository
with any fan-in recommends the chain because waves buy wall clock and cost
merges exactly where they are most likely to conflict; otherwise waves, because
independent plans in separate repositories overlap without a merge at all.

Present the recommendation with its cost, honor an explicit choice without
re-asking, and confirm the shape once before starting.

## The run

Confirm once, then run to completion without further prompting: prepare each
worktree from the reported start, perform any reported integration merge with
ordinary Git, dispatch a worker per plan, wait, inspect real diffs rather than
workers' claims, run the repositories' ordinary checks, record progress, and
mark work complete only when it actually landed and its checks passed.

That last condition is load-bearing. Ordering reads completion to release the
next wave, so unfinished work holds its dependents automatically and there is no
separate readiness state to maintain — which is why completion must never be
recorded optimistically. Recompute the order after each wave rather than
trusting the first result; the world changed during the wave.

## Conflicts and stops

A conflict from an integration merge inside the run is resolved directly: both
sides belong to the same approved intent and both their records and diffs are
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

## The coarseness the report admits

Shared-repository reporting reflects repositories, not paths, because plan
records name repositories. It cannot tell neighbouring modules from the same
file, so it is a signal for choosing a shape rather than a collision predictor.
Start references also name a predecessor's branch: once that work is merged and
its branch is gone the reference needs adjusting by hand, since the derivation
reads records and records do not know what a provider merged.

Owner: `context-circuit-source@internal/workspace/order.go`.
