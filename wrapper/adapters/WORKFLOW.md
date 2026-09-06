# Workflow

This file owns the authority order and the conversational contract for a Context
Circuit workspace.

## Authority order

1. Host and system instructions.
2. `AGENTS.md` safety spine and `WORKFLOW.md` (this file).
3. `workspace.yaml` identity and `wrapper/contracts/invariants.yaml` ownership.
4. Product Knowledge in `context/` (accepted project facts).
5. Intents in `intent/<id>/` (the approved decision: goal, criteria, scope, tier).
6. Plan files in `plans/<plan-id>/` (a derivation of an approved intent; status is a
   projection).
7. Runtime execution evidence in `.runtime/` (never overrides plan status).

## The two human gates

Everything the human decides is one of two gates; everything between is mechanical.

- **Gate 1 — the intent.** The human approves *what "correct" means*, from the plain
  ask alone. Approval freezes the criteria (`contract_digest`) and confirms the
  coordinator understood the ask, which is what lets the tracer read the real code
  next; on approval a read-only tracer maps the code and the coordinator runs a
  feasibility check before any plan is written.
- **Gate 2 — delivery.** The human authorizes the irreversible act (pull request,
  merge, push, deploy). Never implied by verification or completion.

## Conversational actions

| You say | Action |
| --- | --- |
| What is this workspace? | Read-only orientation. |
| Gather context about X. | Write live context files with provenance; keep INDEX.md consistent. |
| Connect / clone / initialize the R repository. | Register and bind; clone/init only when explicitly asked. |
| Work on / refine X with me, or `/cc-pair`. | Explore-tier direct collaboration in one connected repository; live human supervision, no intent, plan, tracer, or verifier — promotable in place. |
| Shape what I want to build (X). | Author an **intent**: goal, non-goals, constraints, outcome-level acceptance criteria, a coarse optional scope, tier — from the plain ask, without reading the code. |
| Approve the intent. | **Gate 1**: the single upstream human approval; freezes the criteria and spawns the tracer to read the real code, then the feasibility check runs. |
| Create a plan for F. | Derive a grounded, readable plan from an approved intent and its trace manifest. Carries no separate approval. |
| Review plan X. | Non-executing discussion; may update draft content. |
| Execute plan X. | Execute a plan authorized by its approved intent; a criteria drift after approval re-enters Gate 1, not run. |
| Execute plans X…Z / run the ready stack. | Run a set of intent-authorized plans in dependency order: the runtime detects which are ready (dependencies verified, paths free) and selects each base; the coordinator may overlap provably-independent ready plans up to a fan-out width (the lease arbitrates races), each still one worker and one independent verifier. Adds no authority; marks nothing done or delivered. |
| What happened with X? | Summarize execution evidence. |
| Repair the failed X verification. | Another worker attempt if allowed. |
| Accept the result / mark X complete. | Accept the candidate if asked; mark the plan done only on an explicit mark-done. Standard and Critical both use explicit mark-done. Explore is planless. When the plan affected Product Knowledge, mark-done then updates live context files. |
| Archive / restore intent or plan X. | Move out of / into the active area; no status change. |
| Open a pull request for X / deliver X. | **Gate 2**: separate delivery; one pull request per covering tip (a repository may have several sibling stacks); source = covering execution branch; target = base branch. Delivery does not spawn a verifier, does not mark the plan done, and does not start a knowledge update. |

## Rules

- Distinguish inspect from mutate, the intent gate from delivery, repository change
  from delivery.
- Never infer a consequential action from "okay" or "looks good".
- The human gate is on the **intent** (Gate 1); a plan derived from an approved
  intent executes with no separate plan approval and no automated scope gate. A
  criteria change after approval re-enters Gate 1; scope-safety is settled at
  delivery (Gate 2).
- At Standard and Critical, one worker and one independent read-only verifier handle
  each execution; three worker failures stop it and preserve all evidence. Explore
  (direct collaboration) has one worker and the human as live oracle and is never
  verified (INV-PAIR-01 / INV-ASSURE-01).
- Verification produces `verified` evidence bound to the candidate. Completion is
  an explicit mark-done at Standard and Critical; Explore is planless and
  verification, candidate acceptance, and delivery never complete anything.
- Delivery (pull request, merge, push, deploy) and cleanup are separate explicit
  actions (Gate 2); a pull request targets the recorded `base_branch`, never
  `default_branch`.
- Failed or interrupted work is preserved, never silently cleaned up.
