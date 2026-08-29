# Architecture

Context Circuit is a universal project workspace (the wrapper) with a compact
conversational entry, human-facing plan artifacts, and private resumable runtime
evidence. The coordinator interprets a request and selects one bounded action;
there is no separately exposed router layer, and conversation keywords are not
the source of lifecycle policy. The source repository separates the shipped
product layer under `wrapper/` from the blank mutable seed under `template/`. An
instantiated workspace keeps the shipped layer, identity, context, sources,
plans, runtime, and registered repositories distinct.

The invariant catalog (`wrapper/contracts/invariants.yaml`) is the
one-owner-per-rule authority; skills and the coordinator are thin
natural-language adapters. The core execution model: one worker executes all
tasks of one approved plan in dependency order across mapped repositories,
committing each repository before an independent read-only verifier checks the
latest commits; repairs add new commits and execution stops after three worker
failures. Workers are isolated by an atomic exclusive-create ownership lock
(`locks/<plan-id>/owner.yaml`) and per-repository worktrees; a live lock is never
silently stolen. Runtime records are filesystem evidence, not a database or
scheduler.
