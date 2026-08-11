---
name: initialize-workspace
description: Inspect, configure, and validate a new Kao Delivery Workspace wrapper. Use when a human creates a wrapper from the template or asks to register repositories, choose solo or team mode, configure optional activity capabilities, or reconcile initialization files safely.
---

# Initialize workspace

1. Read `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, canonical context files,
   and Git status. Record the initial
   wrapper status, inspect available repositories and project documents, then ask
   questions only if required.
2. Resolve only material unknowns with the human: workspace boundary, `solo` or
   `team` mode, repositories and roles, ignored-clone or submodule mode, default
   branches, domain agents, human gates, context sources, and optional activity
   capabilities and required/optional/manual lifecycle actions. Require stable
   action IDs and exact fallback descriptions; do not infer tool availability
   from configuration. Recommend team mode and ignored clones. Record unresolved facts
   explicitly instead of guessing.
3. Preserve existing and unrecorded work. Never reset, clean, overwrite
   team-authored content, replace an existing clone, or silently convert a clone
   to a submodule. Confirm exact paths before cloning or registering submodules.
4. Populate `workspace.yaml`, `WORKFLOW.md`, root instructions, domain agents,
   and canonical context using reviewable edits. Keep provider credentials,
   tokens, MCP configuration, and personal workflow details out of the wrapper.
5. For every configured repository, require an accessible Git root at the exact
   configured path. An ignored clone receives an exact managed `.gitignore`
   entry. A submodule must already be registered at that path in `.gitmodules`
   and must not be ignored. Never ignore `repositories/` wholesale.
6. If the human requests inspection, validation, or a report without changes,
   run only `node .agents/bin/kao.mjs initialize-workspace --check-only` and the read-only
   validators in step 7; never run the applying command. Otherwise, run the
   check-only command first, inspect its findings, then run
   `node .agents/bin/kao.mjs initialize-workspace` to reconcile the exact ignored-clone block.
   Stop on unexpected pre-existing wrapper changes, validation, accessibility,
   mode, or submodule mismatch. Do not confuse the skill's own reviewable edits
   with changes that existed before initialization began.
7. Run `node .agents/bin/kao.mjs validate --check-paths --check-documents`, `git diff --check`,
   and the initialization check again. Present the resulting workspace mode,
   repositories, remotes, default and current branches, repository cleanliness,
   repository instruction discovery, warnings, and wrapper changes for human
   approval.

Reruns must be idempotent: propose reviewable changes, retain team-authored
content, and leave the managed ignore block unchanged when configuration has not
changed. Do not commit, push, open a pull request, or mutate an activity system
unless the human separately authorizes it.
