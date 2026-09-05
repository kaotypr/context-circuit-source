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
tasks of one intent-authorized plan in dependency order in that plan's single
repository (exactly one repository per plan),
committing that repository before an independent read-only verifier checks the
latest commits; repairs add new commits and execution stops after three worker
failures. An intent whose scope covers two repositories derives at least two
plans. Workers are isolated by an atomic exclusive-create ownership lock
(`locks/<plan-id>/owner.yaml`) and per-repository worktrees; a live lock is never
silently stolen. Runtime records are filesystem evidence, not a database or
scheduler.

Assurance is a single consequence-tier ladder — Explore, Standard, Critical
(`wrapper/contracts/invariants.yaml`, INV-ASSURE-01) — declared on the intent from
transparent risk signals and raisable by the human. `cc-pair` (direct
collaboration) is the **Explore tier** of that ladder, not a separate mode: live
human-supervised work in one connected repository with a coordinator and one
worker, no independent verifier, and output never labeled "verified". Standard and
Critical run the execution model above with an independent read-only verifier
bound to the candidate. Explore is planless until an explicit **promote** step
attaches an intent and raises the tier, at which point the verifier appears and
the session becomes a candidate-bearing change without restart. The runtime is
model-blind: it never selects a tier, only classifies risk signals and enforces
the tier floor. Direct collaboration therefore does not alter or bypass any plan
lifecycle gate — it is the ladder's bottom rung, with a ramp upward.
