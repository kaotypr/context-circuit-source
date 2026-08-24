# Terminology

Two vocabularies live in a Context Circuit workspace, and they must not be
confused:

- **Context Circuit terms** — the product's own words for how the workspace
  works (below). These are stable across every project.
- **Your project's terms** — the domain vocabulary of the specific project this
  workspace serves. Those live in `context/TERMINOLOGY.md` as accepted Product
  Knowledge, not here.

This page is the shipped glossary of Context Circuit terms and, most
importantly, how the coordinator translates them into plain language for a user.

## Context Circuit terms

| Term | Meaning |
| --- | --- |
| Wrapper | Accepted synonym for the universal project workspace — the Context Circuit product a user works in. The directory `wrapper/` holds its shipped layer. |
| Product Knowledge | Accepted, agent-oriented understanding of the project, stored as indexed, human-readable units under `context/`. |
| Context unit | One Product Knowledge unit with a stable ID, summary, scope, facts, decisions, and provenance. |
| Context index | The retrieval catalog (`context/INDEX.md`) mapping concepts and aliases to context units. |
| Plan | Human-reviewed intent for an outcome: coverage, grounding, repository map, tasks, acceptance, verification, risks, open questions. |
| Plan status | The human-owned `draft`, `approved`, or `done` state of a plan. |
| Execution | One runtime attempt to implement an approved plan, with one worker and one independent verifier. |
| Worker | The single writer role for one execution. |
| Verifier | The independent, read-only role that checks the worker's latest commits. |
| Repair attempt | A new worker commit plus a new independent check after a failed verification. |
| Completion | The human decision to mark a plan done after a verified execution. |
| Archive / restore | Setting a plan aside, or bringing it back, without changing its status. |
| Delivery | Opening a pull request, merging, pushing, or publishing — always a separate, explicit action. |
| Connected repository | A repository registered in the workspace and resolved to a local checkout. |
| Host-blocked | A state where the environment cannot run a required step (for example, an independent check), so the coordinator reports it and preserves the work rather than faking the step. |

The exact meanings and authority of these terms are settled by the runtime
contracts under `wrapper/contracts/`; this page is the plain-language reference.

## Talking to a user: say the effect, not the mechanism

The coordinator reports actions and state by their effect, in ordinary project
language. It never exposes internal file names or mechanism unless the user
explicitly asks for diagnostics.

| Internal term or artifact | Never say to a user | Say instead, by effect |
| --- | --- | --- |
| Worktree | "worktree" | describe the effect ("a separate working copy of your project"), or say nothing about the mechanism |
| Anchor branch | "anchor branch" | the branch's plain name — "I'll work from `develop`" |
| Local binding / binding | "binding" | "I've connected your `<name>` project" |
| Execution branch (`cc/<plan>/<repo>`) | the `cc/...` branch name | "the changes for `<plan title>`" |
| Default branch (as jargon) | "default branch" | the branch name in plain terms |
| Verifier | "the verifier" | "an independent check" — "I had it checked independently" |
| Worker | "the worker" | "I made the changes" |
| Host-blocked | "host-blocked" | "I can't run an independent check in this environment right now" |
| Workspace and runtime files (`workspace.yaml`, `repositories.local.yaml`, `plan.yaml`, `plans/INDEX.md`, `engine.sh`) | any internal file name | the thing by its effect — "your plan", "your connected projects" |
| Delivery / delivery boundary | "delivery boundary" | "opening a pull request", named as a separate step |
| Archive / restore (as file moves) | the file-move mechanics | "set aside" / "bring back" the plan |

Reveal these mechanics only when the user explicitly asks for diagnostics.
