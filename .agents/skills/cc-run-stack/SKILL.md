---
name: cc-run-stack
description: Execute a named set of already-approved plans in one run (a plan stack), in dependency order, with one worker and one independent verifier per plan; adds no authority and marks nothing done or delivered.
---

## When to use

On a request to run a *set* of approved plans — "execute plans 0001 through 0010",
"run the ready stack", "build all the approved plans in one go". This adds no new
authority over the single-plan `cc-execute`: every plan in the set is separately
approved, executed by one worker, checked by one independent verifier, and left for
the human to complete and deliver. Concurrency only decides order and overlap among
plans the human already approved (INV-CONCURRENCY-01/02).

Refuse to run any plan in the set that is not approved, and say so plainly, without
blocking the rest.

## Resolving the set

Run every runtime action from the workspace directory. Refer to each plan by its
full id (`NNNN-slug`); the human usually names them by number ("0001 through 0010")
— read the active plan index to map each number to its full id and pass the whole
set to `run-stack-ready`. The runtime returns a correct partition even when the set
is large; if a single runtime call ever returns something that contradicts the
recorded plan states, re-check the call itself before concluding the workspace is
broken.

## Runtime actions — invoke, never read the engine

You **must not read** `wrapper/runtime/engine.sh` or any runtime implementation
file (`wrapper/adapters/AGENTS.md` → Runtime owns this boundary); everything the
loop needs is below. Invoke each as `sh wrapper/runtime/engine.sh <action> <args>`
from the workspace directory (its root is the current directory) — use these
actions as a tool.

- `run-stack-ready <plan-id> ...` — partition the set; prints one `<plan>: <bucket>` line each.
- `lease-acquire . <repo> <plan-id> "<path> ..."` — reserve the plan's bounded paths in a repository.
- `lease-release . <repo> <plan-id>` — release on delivery (not merely verification).
- `execution-begin . <plan-id> <owner>` — begin the base-aware execution; prints `execution_id`.
- `attempt-begin <execution-dir>` — start a worker attempt.
- `worker-commit-record <execution-dir> <repo> implementation|repair` — capture the worker's commit.
- `worker-handoff-record <execution-dir> <handoff-file>` — store the worker handoff.
- `verifier-prepare <execution-dir>` — confirm a worker commit exists to verify.
- `verifier-result-record <execution-dir> <attempt> passed|failed|blocked` — record the independent verifier's outcome.
- `attempt-evidence-record <execution-dir> <attempt> <key=value> ...` — record bounded per-attempt host evidence (inference wall-clock and the `(model, effort)` each role ran at). Evidence only; never a gate.
- `repair-allowed <execution-dir>` — whether another repair attempt remains.

The execution directory is `.runtime/executions/<plan-id>/<execution-id>/`; the
assigned worktree(s) are named in each execution's per-repository record.

## The loop

Repeat until no plan in the set is runnable:

1. Ask the runtime to partition the whole set with `run-stack-ready`. It buckets
   each plan as **verified** (already built and checked — not runnable again),
   **ready** (approved, never run, dependencies verified, paths free),
   **waiting** (dependencies not yet verified, or a needed path is leased),
   **failed**, **blocked**, or **refused** (not approved). A plan whose latest
   execution is terminal or in progress is never runnable.

2. If nothing is **ready**, stop the loop.

3. **Fan out over the ready bucket.** Launch up to a **fan-out width** of the
   ready plans as concurrent worker → verifier pipelines, then re-partition when
   any finishes. The safety is already built — INV-CONCURRENCY-01's atomic path
   lease and INV-CONCURRENCY-02's per-plan base make independent-plan overlap
   safe — so this is a coordinator policy, not a new rule; running plans together
   only removes avoidable wall-clock.

   - **Fan-out width is your policy, bounded by the host.** Start conservative
     (2–3 concurrent pipelines); the ceiling is host concurrency and your own
     ability to track several pipelines cleanly, never a contract limit. **Width
     1 is exactly serial id-order behavior** — the compatibility floor; use it for
     a single-repository stack or whenever tracking overlap would be error-prone.
   - **The lease is the race arbiter, so you need not prove disjointness first.**
     Readiness is evaluated at partition time, so two ready plans could target
     overlapping paths; step 3a `lease-acquire` is an atomic exclusive-create, so
     the first to acquire proceeds and the other gets `LEASE_CONFLICT` and stays
     **waiting**. A naive "attempt to lease all ready" is already correct; proving
     independence up front only avoids wasted `execution-begin` work.
   - **Keep the in-flight pipelines from cross-contaminating.** Each pipeline owns
     its own execution directory and records; drive each by its own
     `<execution-dir>`. Runtime records are atomic and per-execution, which
     contains the blast radius.

   For each ready plan you launch this round:

   a. **Acquire path leases.** For every repository the plan touches, acquire a
      lease on its bounded paths with `lease-acquire`. If `lease-acquire` reports a
      `LEASE_CONFLICT`, the region is held by an unrelated plan — leave this plan
      **waiting** and move on; it will become ready when the holder is delivered.

   b. **Begin the execution (base-aware).** Run `execution-begin`. For a plan with
      same-repo predecessors the runtime selects the base for you — the single
      predecessor branch, or a runtime-authored integration merge of two or more —
      and records `based_on`; a plan with no dependency starts from the anchor tip
      exactly as under `cc-execute`. Do **not** acquire leases inside execution;
      this loop owns them. If `execution-begin` reports `status: blocked`
      (`BASE_UNBUILDABLE`) or otherwise fails, the plan is **blocked**, not a worker
      failure: preserve its evidence, run no worker, and hold its descendants.

   c. **One worker.** Assemble the worker brief with `worker-brief-assemble` (it
      carries the repository-grounding directive discovered from the worktree; add
      only a one-line task focus and deliver it verbatim — INV-GROUND-01/03), then
      launch exactly one worker (`agents/worker.md`) with that brief and its
      assigned worktree(s), at the worker's configured `(model, effort)`. It reads
      and honors the repository's own agent guidance, implements every task in
      dependency order inside the assigned worktree and declared paths only, and
      commits each repository. Record each commit with `worker-commit-record` and
      the handoff with `worker-handoff-record`.

   d. **One independent verifier.** After `verifier-prepare`, launch exactly one
      independent, read-only verifier (`agents/verifier.md`), at the verifier's
      configured `(model, effort)`, over the latest commit of every affected
      repository. Record its outcome with `verifier-result-record`, and record the
      observed inference wall-clock and the `(model, effort)` each role ran at with
      `attempt-evidence-record`. If the host cannot create an independent verifier
      child, the result is `host-blocked` — never self-verify.

   e. **Repair within the limit.** On a verifier failure, pass the evidence back to
      the same worker within the same execution: check `repair-allowed`, begin a new
      attempt, let the worker create a new commit, and verify again. When the
      worker's role has `escalate_on_repair: true`, raise the repair attempt's
      `(model, effort)` above its configured start (`docs/role-tiering.md`);
      escalation changes only which model runs the attempt, never the accounting.
      The failure counter increments on each rejection; at three failures the plan
      is **failed** and execution stops with all evidence preserved.

   f. **On verified:** keep the plan's leases held (they hold until delivery, not
      merely verification) so an unrelated plan cannot grab the same region and
      diverge. The next partition will show the plan's descendants as ready.

   g. **On failed or blocked:** preserve every branch, worktree, commit, and record.
      Only this plan's descendants are held; unrelated verified and ready plans are
      untouched — a failure prunes one branch of the graph, not the whole run.

## Report

When nothing is runnable, report in plain project language, by effect: which plans
were built and independently checked, which need attention and why (failed after
its checks, or blocked because a prerequisite did not pass), and which were skipped
because they were not approved. Never expose internal mechanism — no runtime file
names, no `cc/...` branches, no worktrees, no fan-out width, and no `(model,
effort)` values. Do not use a worker's own claim as verifier evidence.

Concurrent progress interleaves — several plans build and check at once — so
narrate interleaved **effects** ("one plan is built and being checked while
another is still building"), never the mechanism behind the overlap. A plan left
waiting for a busy path region reads simply as "waiting on another plan's area";
never mention leases or conflicts. `docs/terminology.md` is the internal→user
mapping.

## Boundaries

Never auto-approve, auto-complete, merge, push, publish, deploy, or discard. Running
a stack produces `verified` evidence per plan; it never marks a plan `done` and
never delivers. Completion and delivery remain separate, explicit, per-plan human
actions. A drift guard (`delivery-drift` / `delivery-rebase`) applies only when the
human later delivers a plan whose base has diverged (see `cc-deliver`).
