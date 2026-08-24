# Context Circuit v0.5 terminology

This file settles the terms used by the design. The purpose is to prevent the
source repository, distributable product template, instantiated project
workspace, Product Knowledge, plans, and runtime state from being described as
if they were the same thing.

## 1. Product names and identities

| Term | Meaning | Do not confuse it with |
| --- | --- | --- |
| Context Circuit | The product concept and design: a universal project workspace that helps agents understand project knowledge, create grounded plans, and execute them safely. | The current source repository or one instantiated project. |
| Cc | Prompt-friendly shorthand for Context Circuit. It is an abbreviation, not a separate product or lifecycle. | A generic context window, compiler, or connector. |
| context-circuit-source | The current maintainer source repository. It contains the Context Circuit implementation, template assembly, tests, source design, and maintainer instructions. | An end-user project workspace. |
| context-circuit-template | The product repository/template distributed for use with projects. It contains the universal project workspace structure and the Context Circuit runtime, contracts, skills, and starter files. | The maintainer source repository. |
| Universal project workspace | The project-oriented product provided by context-circuit-template. It can be initialized for any project and can connect one or more Git repositories. | A single repository checkout. |
| Project workspace | An instantiated universal project workspace configured for one specific project or product area. | The template before initialization. |
| Wrapper | Accepted synonym for the universal project workspace — the Context Circuit product a user works in. Use it interchangeably with the workspace product. | The maintainer source repository; the `wrapper/` implementation directory taken alone. |

The v0.5 product is context-circuit-template as a repository template and
universal project workspace. context-circuit-source is where that product is
designed and maintained.

In a prompt or short conversation, Cc may replace Context Circuit after the
meaning is clear. For example, “Cc, review plan 0001-billing-v2” means “Context
Circuit, review plan 0001-billing-v2.” Formal design, file, and release names use
the full names where ambiguity would matter.

## 2. Source and template terms

| Term | Meaning |
| --- | --- |
| Source repository | The maintainer repository containing the product source and release assembly. |
| Product template | The distributable repository/template that becomes a project workspace. |
| Cc source | Short form of context-circuit-source. It means the maintainer source repository together with its self-hosted Cc development workspace. |
| Cc template | Short form of context-circuit-template. It means the distributable product repository/template and the universal project workspace it provides. |
| Template seed | Blank files and directories copied into a new project workspace. |
| Template-owned file | A file the product may replace during an upgrade. |
| Workspace-owned file | Project identity, Product Knowledge, plans, local bindings, repositories, and runtime evidence that an upgrade must preserve. |
| Source design | Maintainer design material in the source repository, not project Product Knowledge. |
| Source artifact | A named file in a project workspace's sources/ directory or a maintainer source file explicitly selected for reading. |
| Repository source | A connected Git repository used as evidence or implementation scope. |

The word source is always qualified when ambiguity matters: source repository,
source design, source artifact, or repository source.

## 3. Workspace terms

| Term | Meaning |
| --- | --- |
| Workspace identity | Portable project name, purpose, and logical repository identities. |
| Product Knowledge | Agent-oriented, accepted understanding of the product and project. It is stored in structured, indexed, human-readable knowledge units. |
| Context unit | One Product Knowledge unit with stable ID, summary, scope, retrieval metadata, facts, constraints, decisions, and provenance. |
| Context index | Agent retrieval catalog mapping concepts, aliases, domains, repositories, decisions, and constraints to context units. |
| Local binding | Host-local mapping from a logical repository ID to a concrete checkout path. |
| Connected repository | A repository registered in workspace identity and resolved through a local binding. |
| Workspace repository | The optional Git repository at the workspace root, addressed by the reserved logical ID `workspace`. It has its own local anchor branch. |
| Default branch | Optional portable repository-identity hint used when cloning or setting up a repository. It is not a user's active branch and does not by itself authorize or define execution. |
| Anchor branch | Required user-local branch setting for an executable repository binding. It is the user's actual active branch, such as `development` or `kao/development/v0.5`, and the source branch whose tip becomes the runtime worktree base. |
| Anchor commit | The execution-start commit captured from an anchor branch. Worktree branches are created from this immutable execution snapshot. |
| Execution branch | The isolated `cc/<plan-id>/<repository-id>` branch created from an anchor commit for implementation and verification. |
| Pull-request target | The recorded repository `anchor_branch` used as the default target when a human explicitly requests a pull request from an implemented plan. |
| Plan ID | Stable four-digit sequence plus lowercase kebab-case short-title identifier, such as `0001-billing-v2`, used to mention, locate, and reference a plan. The ID is not the plan title or status. |
| Plan short title | The concise lowercase kebab-case suffix after a plan's four-digit sequence. It describes the requested outcome without replacing the canonical Plan ID. |
| Active plan index | `plans/INDEX.md`, the catalog of plans available for normal agent operations. It excludes archived plans. |
| Archived plan | A plan directory moved under `plans/.archived/`. Its plan files and status are preserved, but the normal agent must not read it until an explicit restore action returns it to `plans/`. |
| Plan | Human-reviewed intent for a project outcome, including request coverage, context grounding, repository map, tasks, acceptance, verification, risks, and open questions. |
| Task | Detailed, bounded unit of a plan assigned to one or more repositories. |
| Execution | One runtime attempt to implement an approved plan with one worker and an independent verifier. |
| Implementation completion record | Durable evidence that a human marked a plan done, linked to the verified execution, plan revision, repository commits, and verifier result. |
| Knowledge reconciliation | Post-implementation comparison of the completed changes against relevant Product Knowledge to find durable additions, changes, invalidated claims, or no update needed. |
| Context update proposal | A reviewable, evidence-linked proposal to add, change, retire, or mark a Product Knowledge unit stale. |
| Root agent / coordinator | The conversational agent responsible for workspace interaction, context gathering, planning, review, lifecycle interpretation, and delivery discussion. It is not a child execution role. |
| Planner | Informal name for the root agent's planning capability. It is not a separate agent or child role. |
| Worker | The sole writer role for one plan execution. |
| Verifier | Independent read-only role that checks the worker's latest commits. |
| Repair attempt | A new worker commit followed by a new verifier check after a failed verification. |
| Completion | Human-controlled change of plan status after verified execution. |
| Archive | Explicit file-organization action that moves a plan out of the active plan area without changing plan status, validating execution, or implying completion. |
| Restore | Explicit file-organization action that moves an archived plan back into the active plan area without changing its preserved plan status. |
| Delivery | Separate pull-request creation, merge, push, publication, deployment, or integration action. |

## 4. Runtime terms

| Term | Meaning |
| --- | --- |
| Context Circuit runtime | Deterministic internal support for workspace validation, repository bindings, worktrees, branches, execution records, commits, verification evidence, and recovery. |
| Host adapter | Provider-specific integration that interprets conversation and launches worker/verifier children. |
| Execution brief | Generated, bounded input describing the approved plan, context references, repositories, worktrees, tasks, evidence, and stop conditions. |
| Execution evidence | Runtime records, commits, handoffs, and verifier results describing what happened. |
| Plan status | Human-owned draft, approved, or done state. |
| Execution status | Runtime state such as running, verifying, repairing, verified, failed, or blocked. |

## 5. The term "wrapper"

Wrapper is an accepted synonym for the universal project workspace — the Context
Circuit product. "The wrapper" and "the workspace product" name the same thing;
use either according to what reads clearly.

The physical directory wrapper/ holds the product's shipped layer (runtime,
contracts, adapters) inside context-circuit-source and context-circuit-template.
The term names the product; the directory names one part of it. When referring
to the directory, write the path explicitly so the two are not confused. For
example:

- the runtime implementation at wrapper/runtime/engine.sh;
- the template-owned wrapper/ directory is replaced during an upgrade;
- the wrapper (equivalently, the project workspace) creates a plan;
- the Context Circuit runtime prepares worktrees.

## 6. Authority by term

- context-circuit-source owns implementation and release assembly.
- context-circuit-template owns the distributable workspace product.
- workspace files own project identity and accepted Product Knowledge.
- plan files own plan intent and human plan status.
- `plans/INDEX.md` owns active plan discovery and mention lookup, not plan
  status.
- `plans/.archived/` owns archived plan placement, not a new lifecycle status.
- the Context Circuit runtime owns execution evidence and repository safety.
- the host adapter owns provider-native child creation.
- the worker owns implementation commits.
- the verifier owns independent verification results.

## 7. User-facing translation of internal terms

The terms above are the design and runtime vocabulary. They are not the words a
coordinator uses with an ordinary user. The coordinator reports actions and
state by their effect, in plain project language, and never exposes internal
mechanism unless the user explicitly asks for diagnostics (doc 02 §1). This
table is the single canonical mapping; the shipped projection is
`docs/terminology.md`, and the coordinator role references that projection
instead of restating the list.

| Internal term or artifact | Never say to a lay user | Say instead, by effect |
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
