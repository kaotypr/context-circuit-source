<!-- context-circuit:title -->
<h1 align="center">Context Circuit</h1>
<!-- /context-circuit:title -->

<!-- context-circuit:badges -->
<!-- /context-circuit:badges -->

<!-- context-circuit:purpose -->
<!-- /context-circuit:purpose -->

## What this repository is

This is a **Context Circuit workspace**: the shared place a project keeps its
context, its goals, and its plans, separate from the code they describe. It is
an ordinary Git repository. Review it, branch it, and share it the way you
already share project documentation.

The code lives elsewhere. A workspace binds to one or more repositories and
records where each one is on the machine you are working from, so the same
workspace serves everybody without pinning anyone's directory layout.

## What Context Circuit is

Coding agents are good at changing code. The harder problem is carrying the
right project context across sessions, repositories, people, and machines.

Context Circuit gives an agent a shared place to learn the project, prepare a
goal before coding starts, coordinate work across repositories, update what the
project knows when the code changes, and keep what the team learned afterward.
The decisions that matter stay with you: a plan is approved by a person, not
inferred by an agent.

It is two products — this workspace template, and a native CLI that handles the
records, bindings, working copies, and diagnostics. The badges above report the
versions *this* workspace received, not whatever was published since.

## Where to look

| Path | What it holds |
| --- | --- |
| `context/` | What the project knows, catalogued by `context/INDEX.md` |
| `intent/` | Goals, with the reasoning and approval behind each one |
| `plans/` | Standalone or intent-linked plan folders, with optional supporting design files; legacy single-file plans remain readable |
| `sources/` | Reference material the project did not write |
| `workspace.yaml` | Repositories, relationships, and this workspace's name |
| `members.yaml` | Who works here |
| `.context-circuit/` | Documentation and version stamps the CLI maintains |

Start with `AGENTS.md` for how an agent is expected to work here, and
`.context-circuit/docs/` for the workspace contract and the command surface.

## This file is yours

Context Circuit wrote this README once, when the workspace was initialized. It
belongs to this project now. Edit it, replace it, or delete sections freely — a
later template update will not write over it.
