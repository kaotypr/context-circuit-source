# Conventions

- Keep canonical behavior under `.agents/` and host adapters thin.
- Use UTC identifiers, full Git object IDs, and repository-scoped worktrees.
- Keep repository-specific conventions in each repository.
- Keep plans and activity providers optional for explicit scoped work.
- Preserve dirty, unpushed, unrecorded, failed, and blocked work; never reset,
  stash, clean, or implicitly delete it as workflow recovery.
- Recommend ignored clones with exact `.gitignore` entries; preserve registered
  submodules and never ignore `repositories/` wholesale.
- Keep the unused template neutral: no `repositories/` directory, repository
  registration, or technology-specific agent. Track other intentional empty
  directories with `.gitkeep`.
- Keep contributions append-only and cite them in canonical context updates.
- Keep secrets out of configuration, runtime evidence, plans, and contributions.
- Require schema-valid compact worker/verifier results and independent evidence.
- Require humans for plan approval, task selection, code merge, and closeout
  outcome. An explicit request for known work satisfies task selection.
