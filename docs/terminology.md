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
| Worker | The single role that implements an approved plan and commits its changes for one execution. |
| Verifier | The independent, read-only role that checks the worker's latest commits. |
| Repair attempt | A new worker commit plus a new independent check after a failed verification. |
| Completion | The human decision to mark a plan done after a verified execution. |
| Archive / restore | Setting a plan aside, or bringing it back, without changing its status. |
| Delivery | Opening a pull request, merging, or pushing — always a separate, explicit action. "Publish" is not a delivery word; it names the external surface. |
| Connected repository | A repository registered in the workspace and resolved to a local checkout. |
| Host-blocked | A state where the environment cannot run a required step (for example, an independent check), so the coordinator reports it and preserves the work rather than faking the step. |
| Plan stack | A named set of approved plans executed in one run; the runtime orders and overlaps them safely without changing any gate. |
| Plan dependency | Inter-plan ordering (`plan_dependencies`), distinct from a task's intra-plan `depends_on`; declaring it makes the plan `schema_version: 2`. |
| Path lease | A `(repository, path region)` reservation extending the one-worker lock, so plans touching overlapping paths serialize while disjoint ones run together. |
| Fan-out width | How many provably-independent ready plans a plan stack runs at once — a coordinator policy bounded by the host, not a gate; width 1 is plain serial order, and the path lease arbitrates any race. |
| Role tiering | The per-role `(model, effort)` the coordinator runs the worker and verifier at, from a host-local config with adapter defaults; bounded host evidence that changes cost and speed, never meaning, and never independence or the failure limit. |
| Execution base | The commit a plan's work is built on: the anchor tip, a predecessor's branch (stack), or a runtime-authored integration merge of several predecessors. |
| Drift guard | Rebasing a plan onto the current branch tip and re-checking it before a pull request, when a sibling already merged. |
| Repository grounding | The worker reading and honoring the target repository's own agent guidance, discovered from the working copy. |
| Grounding manifest | The discovered record of a repository's agent guidance (files, skills, prepared environment) for one execution. |
| Worker brief | The instructions handed to the worker for one execution, assembled from the grounding manifest and the plan. |
| System design | A structured write-up of the shape of a larger change, authored as source material; a source, not a lifecycle stage. |
| Publication | A pipeline you declare once and trigger by hand (`cc-publish`) to publish workspace data to an external system — a tracker, chat, or docs space; separate from the core workflow, one folder per publication under `publication/`. |
| Publish / `cc-publish` | Sending workspace data outward to an external system, on request. Reserved for the external surface — git delivery is "push" or "open a pull request", never "publish". |
| Publication kind | What a publication sends: `plan` (a plan and its tasks → a tracker) or `thread` (a plan's open questions → a chat discussion). |

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
| Path lease / region | "lease", "path region" | say nothing about the mechanism; if relevant, "I ran those together" or "I did those one after another" |
| Execution base / integration base / `refs/cc-base/...` | "base", "integration merge", the ref name | "built on top of the earlier work", by effect |
| Grounding manifest / worker brief | those terms | "I followed your project's own contributor guidance" |
| Drift guard / rebase-and-re-verify | "drift", "rebase" | "I brought it up to date with your branch and checked it again" |
| Fan-out width / concurrent pipelines | "fan-out", "width N", "leases" | "I'm working on several at once" / a waiting plan is "waiting on another plan's area" |
| Role tiering / `(model, effort)` | the model id or effort level | say nothing about the model; describe only the effect, "I ran the harder step with more effort" |

Reveal these mechanics only when the user explicitly asks for diagnostics.
