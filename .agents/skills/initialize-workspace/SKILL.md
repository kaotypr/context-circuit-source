---
name: initialize-workspace
description: Inspect, configure, and validate a new Context Circuit wrapper. Use when a human creates a wrapper from the template or asks to register repositories, choose solo or team mode, configure optional activity capabilities, or reconcile initialization files safely.
---

# Initialize workspace

1. Read `AGENTS.md`, `WORKFLOW.md`, the neutral `workspace.yaml`, canonical
   context, and available Git state. Determine whether this is a new wrapper or
   an initialized wrapper. Do not treat the absence of `.git` or `repositories/`
   as an error in a fresh release archive.
2. Ask only for material choices: workspace name, `solo` or `team`, wrapper
   branch, repositories and roles, repository source, paths, branches, domain
   agents, context sources, human gates, and optional activity behavior.
   Recommend team mode, ignored repositories, `main`, and no activity provider.
   For each repository classify exactly one source:
   - `new`: create a new local Git repository and empty base commit;
   - `clone`: clone a credential-free URL into a new ignored path;
   - `existing`: register an existing local Git root without changing it;
   - `submodule`: add a credential-free URL as a tracked Git submodule.
3. Inspect every proposed path before mutation. Never replace an existing path,
   relocate a repository, reset, stash, clean, overwrite authored content, or
   convert repository mode silently. Never accept embedded credentials or ignore
   `repositories/` wholesale.
4. Normalize the approved configuration and actions into an ignored
   `.runtime/bootstrap/request.json` document that validates against the
   `workspace-bootstrap-request` contract. For a new or unborn wrapper,
   require explicit authorization for its configured initial commit. For every
   `new` repository, require explicit authorization for its empty base commit.
   Existing, cloned, and submodule repositories must not authorize an artificial
   commit. Use an explicitly supplied commit author only when the human provides
   both name and email; otherwise use their configured Git identity.
5. Present the exact action plan before running it: files to configure, wrapper
   Git initialization, repository creation/clone/registration, ignore entries,
   domain-agent files, initial commits, and validation. After explicit approval,
   run `node .agents/bin/cc.mjs initialize-workspace --bootstrap .runtime/bootstrap/request.json`.
   The command makes the configured wrapper state its first commit; it does not
   commit changes in an already initialized wrapper.
6. For an already initialized workspace that needs no bootstrap, first run
   `node .agents/bin/cc.mjs initialize-workspace --check-only`, make only the
   approved reviewable configuration edits, then run
   `node .agents/bin/cc.mjs initialize-workspace`. An inspection-only request
   must remain read-only and stop after check-only validation.
7. Run `node .agents/bin/cc.mjs validate --check-paths --check-documents`,
   `git diff --check`, and initialization check-only again. Present the wrapper
   commit when created, repository source and mode, remotes, default/current
   branches, cleanliness, instructions, warnings, and remaining wrapper changes.

Reruns must be idempotent: propose reviewable changes, retain team-authored
content, and leave the managed ignore block unchanged when configuration has not
changed. Initial-commit authorization does not authorize push, remote creation,
pull requests, later commits, or activity mutation.
