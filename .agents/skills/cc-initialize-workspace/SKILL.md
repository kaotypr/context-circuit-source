---
name: cc-initialize-workspace
description: Initialize a solo or team workspace with repositories or a valid no-repository project path.
---

# Initialize a workspace

Use this skill when the user asks to start, set up, initialize, or orient a
Context Circuit workspace.

## Read

- `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, `context/index.md`,
  `context/WORKSPACE.md`, and `context/PROJECT.md`;
- only the relevant Product Knowledge; and
- repository metadata needed to confirm paths and available branches.

Do not scan `sources/` during ordinary initialization. Read a source only when
the user explicitly chooses source intake as part of the next action.

## Ask and record

Ask one compact set of core identity questions:

1. Is the workspace for a solo developer or a team?
2. Which repositories or project items are known? Zero repositories is valid.
3. For each repository, what is its role and which branch is the default active
   branch? Recommend `development` only when that branch exists; let the user
   choose another branch.

Record the confirmed result in `workspace.yaml` and summarize accepted
workspace identity in `context/WORKSPACE.md`. Do not write workspace identity
into `context/PROJECT.md`. Repository entries contain a path, mode, role,
agent, and `default_branch`. Keep repository paths specific and
credential-free.

For a no-repository workspace, offer an Idea Brief conversation, selected
source intake, both, or deferment. This is a first-class project state.

## Do not ask

Initialization does not ask about delivery behavior, commits, pushes, merges,
publication, deployment, or external activity integrations. Those belong to
optional later configuration.

## Output and gate

Output a compact proposed identity and the next optional path. Ask the user to
confirm the proposed workspace identity before treating it as accepted
Product Knowledge or writing it as the workspace baseline in
`workspace.yaml` and `context/WORKSPACE.md`. Preserve uncertainty and
unanswered choices explicitly.

## Next action

After confirmation, continue with the user's chosen next action: `cc-idea-brief`,
request-scoped source intake, `cc-create-prd`, or ordinary orientation. Do not
force an unused artifact step.
