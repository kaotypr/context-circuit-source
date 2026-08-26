# Workflow

This file owns the authority order and the conversational contract for a Context
Circuit workspace.

## Authority order

1. Host and system instructions.
2. `AGENTS.md` safety spine and `WORKFLOW.md` (this file).
3. `workspace.yaml` identity and `wrapper/contracts/invariants.yaml` ownership.
4. Product Knowledge in `context/` (accepted project facts).
5. Plan files in `plans/<plan-id>/` (plan intent and human plan status).
6. Runtime execution evidence in `.runtime/` (never overrides plan status).

## Conversational actions

| You say | Action |
| --- | --- |
| What is this workspace? | Read-only orientation. |
| Gather context about X. | Propose a context update with provenance. |
| Connect / clone / initialize the R repository. | Register and bind; clone/init only when explicitly asked. |
| Create a plan for F. | Draft a grounded readable plan. Does not approve or execute. |
| Review plan X. | Non-executing discussion; may update draft content. |
| Approve plan X. | Explicit approval gate: draft → approved. |
| Execute plan X. | Execute only if already approved. |
| Approve plan X and execute it. | Approve, then execute if preflight passes. |
| Execute plans X…Z / run the ready stack. | Run a set of already-approved plans in dependency order: the runtime detects which are ready (dependencies verified, paths free) and selects each base; each plan is still one worker and one independent verifier. Adds no authority; marks nothing done or delivered. |
| What happened with X? | Summarize execution evidence. |
| Repair the failed X verification. | Another worker attempt if allowed. |
| Mark X complete. | Human-controlled completion; only when verified. |
| Review / accept context updates for X. | Discuss / accept a knowledge proposal. |
| Archive / restore plan X. | Move out of / into the active plan area; no status change. |
| Open a pull request for X. | Separate delivery; source = execution branch, target = anchor branch. |

## Rules

- Distinguish inspect from mutate, approval from execution, repository change from
  delivery.
- Never infer a consequential action from "okay" or "looks good".
- Only an approved plan executes; refuse to execute a draft plainly.
- One worker, one independent read-only verifier; three worker failures stop
  execution and preserve all evidence.
- Verification produces `verified` evidence; only an explicit human request marks
  a plan `done`, and only when verified.
- Delivery (pull request, merge, push, publish, deploy) and cleanup are separate
  explicit actions; a pull request targets the recorded `anchor_branch`, never
  `default_branch`.
- Failed or interrupted work is preserved, never silently cleaned up.
